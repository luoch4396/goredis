package db

import (
	"goredis/config"
	"sync/atomic"
)

var (
	//暂不支持的命令行操作
	unknownOperation = []byte("-ERR unknown\r\n")
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
	if config.GlobalProperties.Databases == 0 {
		//redis 默认16个数据库
		config.GlobalProperties.Databases = 16
	}
	for i := range singleServer.dbs {
		var singleDB = newDB()
		singleDB.index = i
		var oneDb = &atomic.Value{}
		oneDb.Store(singleDB)
		singleServer.dbs[i] = oneDb
	}
	return singleServer
}
