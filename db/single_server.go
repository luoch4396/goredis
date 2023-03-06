package db

import (
	"goredis/config"
	"sync/atomic"
)

type SingleServer struct {
	dbs  []*atomic.Value //db
	role int32
}

// NewSingleServer 创建一个单机redis服务
func NewSingleServer() *SingleServer {
	var singleServer = &SingleServer{}
	if config.GlobalProperties.Databases == 0 {
		config.GlobalProperties.Databases = 8
	}
	for i := range singleServer.dbs {
		//singleDB := makeDB()
		//singleDB.index = i
		var holder = &atomic.Value{}
		//holder.Store(singleDB)
		singleServer.dbs[i] = holder
	}
	return singleServer
}
