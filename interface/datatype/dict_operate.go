package datatype

type Consumer func(key string, val interface{}) bool

// Dict redis 基本数据结构dict
type Dict interface {
	Get(key string) (valve interface{}, exists bool)

	Len() int

	Put(key string, valve interface{}) (result int)

	PutIfAbsent(key string, valve interface{}) (result int)

	PutIfPresent(key string, valve interface{}) (result int)

	Remove(key string) (result int)

	ForEach(consumer Consumer)

	RandomKeys(limit int) []string

	RandomDistinctKeys(limit int) []string

	Clear()
}
