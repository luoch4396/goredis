package db

import (
	"fmt"
	"goredis/config"
	"goredis/interface/tcp"
	"goredis/pkg/log"
	"goredis/redis"
	"goredis/redis/exchange"
	"runtime/debug"
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

func (server *SingleServer) Exec(client redis.ClientConn, cmd [][]byte) (rep tcp.Info) {
	defer func() {
		if err := recover(); err != nil {
			log.Warning(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
			rep = &exchange.BulkRequest{}
		}
	}()
	return nil
}

func (server *SingleServer) Close() {

}
