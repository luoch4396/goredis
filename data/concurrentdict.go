package data

import (
	"goredis/data/interface"
	"goredis/pkg/errors"
	"goredis/pkg/log"
	"goredis/pkg/utils/hasher"
	"math"
	"sync"
	"sync/atomic"
)

// ConcurrentDict 适用读多写少场景
type ConcurrentDict struct {
	dictShard  []*dictShard
	count      int32
	shardCount int
}

// 使用数据分片的方式和分段读写锁提高性能
type dictShard struct {
	table         map[string]interface{}
	readWriteLock sync.RWMutex
}

func NewConcurrentDict(shardCount int) *ConcurrentDict {
	shardCount = computeCapacity(shardCount)
	var table = make([]*dictShard, shardCount)
	for i := 0; i < shardCount; i++ {
		table[i] = &dictShard{
			table: make(map[string]interface{}),
		}
	}
	var concurrentDict = &ConcurrentDict{
		count:      0,
		dictShard:  table,
		shardCount: shardCount,
	}
	return concurrentDict
}

func computeCapacity(param int) (size int) {
	if param <= 16 {
		return 2 << 3
	}
	var n = param - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return math.MaxInt32
	}
	return n + 1
}

func (dict *ConcurrentDict) spread(hashCode uint32) uint32 {
	var tableSize = uint32(len(dict.dictShard))
	return (tableSize - 1) & hashCode
}

func (dict *ConcurrentDict) getShard(index uint32) *dictShard {
	return dict.dictShard[index]
}

func (dict *ConcurrentDict) addCount() int32 {
	return atomic.AddInt32(&dict.count, 1)
}

func (dict *ConcurrentDict) Put(key string, value interface{}) (result bool) {
	_, err := errors.CheckIsNotNull(dict)
	if err != nil {
		log.Errorf("", err)
		return false
	}
	var hashCode = hasher.Sum32([]byte(key))
	var index = dict.spread(hashCode)
	var s = dict.getShard(index)
	s.readWriteLock.Lock()
	defer s.readWriteLock.Unlock()
	if _, ok := s.table[key]; ok {
		s.table[key] = value
		return true
	}
	dict.addCount()
	s.table[key] = value
	return false
}

func (dict *ConcurrentDict) Size() int {
	return int(atomic.LoadInt32(&dict.count))
}

func (dict *ConcurrentDict) Get(key string) (value interface{}, exists bool) {
	_, err := errors.CheckIsNotNull(dict)
	if err != nil {
		log.Errorf("", err)
		return nil, false
	}
	var hashCode = hasher.Sum32([]byte(key))
	var index = dict.spread(hashCode)
	var s = dict.getShard(index)
	s.readWriteLock.RLock()
	defer s.readWriteLock.RUnlock()
	value, exists = s.table[key]
	return
}

func (dict *ConcurrentDict) PutIfAbsent(key string, value interface{}) (result bool) {
	_, err := errors.CheckIsNotNull(dict)
	if err != nil {
		log.Errorf("", err)
		return false
	}
	var hashCode = hasher.Sum32([]byte(key))
	var index = dict.spread(hashCode)
	var s = dict.getShard(index)
	s.readWriteLock.Lock()
	defer s.readWriteLock.Unlock()
	if _, ok := s.table[key]; ok {
		return true
	}
	s.table[key] = value
	dict.addCount()
	return false
}
func (dict *ConcurrentDict) PutIfPresent(key string, value interface{}) (result bool) {
	_, err := errors.CheckIsNotNull(dict)
	if err != nil {
		log.Errorf("", err)
		return false
	}
	var hashCode = hasher.Sum32([]byte(key))
	var index = dict.spread(hashCode)
	var s = dict.getShard(index)
	s.readWriteLock.Lock()
	defer s.readWriteLock.Unlock()
	if _, ok := s.table[key]; ok {
		s.table[key] = value
		return false
	}
	return true
}

func (dict *ConcurrentDict) Remove(key string) (result bool) {
	return false
}

func (dict *ConcurrentDict) ForEach(consumer _interface.DictConsumer) {
	_, err := errors.CheckIsNotNull(dict)
	if err != nil {
		log.Errorf("", err)
		return
	}
	for _, shard := range dict.dictShard {
		shard.readWriteLock.RLock()
		res := func() bool {
			defer shard.readWriteLock.RUnlock()
			for key, value := range shard.table {
				var isNext = consumer(key, value)
				if !isNext {
					return false
				}
			}
			return true
		}
		if !res() {
			break
		}
	}
}

func (dict *ConcurrentDict) Clear() {
	*dict = *NewConcurrentDict(dict.shardCount)
}
