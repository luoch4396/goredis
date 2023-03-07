package db

import (
	data2 "goredis/data"
	"goredis/interface/data"
)

const (
	//redis hash槽数量
	dataDictSize = 2 << 15
	ttlDictSize  = 2 << 10
	lockerSize   = 2 << 10
)

// DB 定义redis的db
type DB struct {
	index      int
	data       data.Dict
	ttlMap     data.Dict
	versionMap data.Dict
}

func newDB() *DB {
	db := &DB{
		data:       data2.NewConcurrentDict(dataDictSize),
		ttlMap:     data2.NewConcurrentDict(ttlDictSize),
		versionMap: data2.NewConcurrentDict(dataDictSize),
		//locker:     lock.NewLocker(lockerSize),
	}
	return db
}
