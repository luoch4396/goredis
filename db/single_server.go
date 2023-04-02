package db

import (
	"goredis/interface/redis"
	"goredis/interface/tcp"
	"goredis/pkg/errors"
	"goredis/pkg/log"
	"goredis/redis/config"
	"strings"
	"sync/atomic"
)

var (
	//定义策略
	pingStrategy = NewCmdOperator(&PingStrategy{})
	authStrategy = NewCmdOperator(&AuthStrategy{})
	infoStrategy = NewCmdOperator(&InfoStrategy{})
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
		var singleDB = NewDB()
		singleDB.Index = i
		var oneDb = &atomic.Value{}
		oneDb.Store(singleDB)
		singleServer.dbs[i] = oneDb
	}
	log.Info("create default 16 databases success!")
	return singleServer
}

func (server *SingleServer) Exec(client redis.ClientConn, cmdBytes [][]byte) (rep tcp.Info) {
	//TODO:错误恢复 移动至协程池？
	//defer func() {
	//	if err := recover(); err != nil {
	//		log.Warning(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
	//		rep = &exchange.UnknownErrResponse{}
	//	}
	//}()
	cmdName := strings.ToLower(string(cmdBytes[0]))
	switch cmdName {
	case "ping":
		return pingStrategy.DoCmdStrategy(&client, cmdBytes[1:])
	case "auth":
		return authStrategy.DoCmdStrategy(&client, cmdBytes[1:])
	case "info":
		return infoStrategy.DoCmdStrategy(&client, cmdBytes)
	}
	if !isAuthenticated(&client) {
		return errors.NewStandardError("NOAUTH Authentication required")
	}
	dbIndex := client.GetDBIndex()
	_, errRep := server.selectDB(dbIndex)
	if errRep != nil {
		return errRep
	}
	return nil
	//selectedDB.Exec(client, cmdBytes)
}

func (server *SingleServer) Close() {

}

func (server *SingleServer) selectDB(dbIndex int) (*DB, *errors.StandardError) {
	//验证
	if dbIndex >= len(server.dbs) || dbIndex < 0 {
		return nil, errors.NewStandardError("DB index is out of range")
	}
	return server.dbs[dbIndex].Load().(*DB), nil
}

func isAuthenticated(client *redis.ClientConn) bool {
	//if config.Properties.RequirePass == "" {
	//	return true
	//}
	//return client.GetPassword() == config.Configs.Server.PassWord
	return true
}
