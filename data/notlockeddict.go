package data

import (
	_interface "goredis/data/interface"
	"goredis/pkg/errors"
	"goredis/pkg/hasher"
	"goredis/pkg/log"
	"goredis/pkg/utils"
	"sync"
	"sync/atomic"
)

// SpinDict TODO:尝试使用比锁跟轻量级的实现，目前并未实现
type SpinDict struct {
	spinDictShard []*spinDictShard
	count         int32
	shardCount    int
}

type spinDictShard struct {
	table map[string]interface{}
	lock  sync.Locker
}

func NewSpinDict(shardCount int) *SpinDict {
	shardCount = computeCapacity(shardCount)
	var table = make([]*spinDictShard, shardCount)
	for i := 0; i < shardCount; i++ {
		table[i] = &spinDictShard{
			table: make(map[string]interface{}),
			lock:  utils.NewLightLock(16),
		}
	}
	var spinDict = &SpinDict{
		count:         0,
		spinDictShard: table,
		shardCount:    shardCount,
	}
	return spinDict
}

func (dict *SpinDict) spread(hashCode uint32) uint32 {
	var tableSize = uint32(len(dict.spinDictShard))
	return (tableSize - 1) & hashCode
}

func (dict *SpinDict) getShard(index uint32) *spinDictShard {
	return dict.spinDictShard[index]
}

func (dict *SpinDict) addCount() int32 {
	return atomic.AddInt32(&dict.count, 1)
}

func (dict *SpinDict) Put(key string, value interface{}) (result bool) {
	_, err := errors.CheckIsNotNull(dict)
	if err != nil {
		log.Errorf("", err)
		return false
	}
	var hashCode = hasher.Sum32([]byte(key))
	var index = dict.spread(hashCode)
	var s = dict.getShard(index)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.table[key]; ok {
		s.table[key] = value
		return true
	}
	dict.addCount()
	s.table[key] = value
	return false
}

func (dict *SpinDict) Size() int {
	return int(atomic.LoadInt32(&dict.count))
}

func (dict *SpinDict) Get(key string) (value interface{}, exists bool) {
	_, err := errors.CheckIsNotNull(dict)
	if err != nil {
		log.Errorf("", err)
		return nil, false
	}
	var hashCode = hasher.Sum32([]byte(key))
	var index = dict.spread(hashCode)
	var s = dict.getShard(index)
	//TODO:读是否可以不加锁？
	s.lock.Lock()
	defer s.lock.Unlock()
	value, exists = s.table[key]
	return
}

func (dict *SpinDict) PutIfAbsent(key string, value interface{}) (result bool) {
	_, err := errors.CheckIsNotNull(dict)
	if err != nil {
		log.Errorf("", err)
		return false
	}
	var hashCode = hasher.Sum32([]byte(key))
	var index = dict.spread(hashCode)
	var s = dict.getShard(index)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.table[key]; ok {
		return true
	}
	s.table[key] = value
	dict.addCount()
	return false
}
func (dict *SpinDict) PutIfPresent(key string, value interface{}) (result bool) {
	_, err := errors.CheckIsNotNull(dict)
	if err != nil {
		log.Errorf("", err)
		return false
	}
	var hashCode = hasher.Sum32([]byte(key))
	var index = dict.spread(hashCode)
	var s = dict.getShard(index)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.table[key]; ok {
		s.table[key] = value
		return false
	}
	return true
}

func (dict *SpinDict) Remove(key string) (result bool) {
	return false
}

func (dict *SpinDict) ForEach(consumer _interface.DictConsumer) {
	_, err := errors.CheckIsNotNull(dict)
	if err != nil {
		log.Errorf("", err)
		return
	}
	for _, shard := range dict.spinDictShard {
		//shard.lock.Lock()
		res := func() bool {
			//defer shard.lock.Unlock()
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

func (dict *SpinDict) Clear() {
	*dict = *NewSpinDict(dict.shardCount)
}
