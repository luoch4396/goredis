package db

import (
	data2 "goredis/data"
	"goredis/data/interface"
)

const (
	//redis hash槽数量
	dataDictSize = 2 << 15
	ttlDictSize  = 2 << 10
	lockerSize   = 2 << 10
)

// DB 定义redis的db
type DB struct {
	Index      int
	data       _interface.Dict
	ttlMap     _interface.Dict
	versionMap _interface.Dict
}

func NewDB() *DB {
	db := &DB{
		data:       data2.NewConcurrentDict(dataDictSize),
		ttlMap:     data2.NewConcurrentDict(ttlDictSize),
		versionMap: data2.NewConcurrentDict(dataDictSize),
		//locker:     lock.NewLocker(lockerSize),
	}
	return db
}
