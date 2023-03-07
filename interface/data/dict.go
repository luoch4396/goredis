package data

type DictConsumer func(key string, val interface{}) bool

// Dict redis-字典数据结构定义
type Dict interface {
	Get(key string) (valve interface{}, exists bool)
	Size() int
	Put(key string, valve interface{}) (result bool)
	PutIfAbsent(key string, valve interface{}) (result bool)
	PutIfPresent(key string, valve interface{}) (result bool)
	Remove(key string) (result bool)
	ForEach(consumer DictConsumer)
	Clear()
}
