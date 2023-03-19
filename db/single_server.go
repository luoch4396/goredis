package db

import (
	"goredis/config"
	"goredis/pkg/log"
	"sync/atomic"
)

type SingleServer struct {
	//db数组
	dbs []*atomic.Value
	//角色（主、从）
	role int32
}

// NewSingleServer 创建一个简单的单机redis服务
func NewSingleServer() *SingleServer {
	var singleServer = &SingleServer{}
	configs := &config.Configs
	if configs.Server.Databases == 0 {
		//redis 默认16个数据库
		configs.Server.Databases = 16
	}
	for i := range singleServer.dbs {
		var singleDB = newDB()
		singleDB.index = i
		var oneDb = &atomic.Value{}
		oneDb.Store(singleDB)
		singleServer.dbs[i] = oneDb
	}
	log.Info("create default 16 databases success!")
	return singleServer
}
