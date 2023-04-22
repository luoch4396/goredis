package db

import (
	"goredis/interface/redis"
	"goredis/interface/tcp"
	"goredis/pkg/errors"
	"goredis/pkg/log"
	"goredis/pkg/utils"
	"goredis/redis/config"
	"strings"
	"sync/atomic"
	"time"
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
	if config.GetDatabases() == 0 {
		//redis 默认16个数据库
		config.SetDatabases(16)
	}
	for i := range singleServer.dbs {
		var singleDB = NewDB()
		singleDB.Index = i
		var oneDb = &atomic.Value{}
		oneDb.Store(singleDB)
		singleServer.dbs[i] = oneDb
	}
	if config.GetServerType() == "" {
		config.SetServerType("single")
	}
	config.SetStartUpTime(time.Now())
	log.Info("create default 16 databases success!")
	return singleServer
}

func (server *SingleServer) Exec(client redis.ClientConn, cmdBytes [][]byte) (rep tcp.Info) {
	//TODO:错误恢复 移动至协程池？
	defer func() {
		if err := recover(); err != nil {
			log.MakeErrorLog(err)
			rep = &errors.UnknownError{}
		}
	}()
	//认证
	if !isAuthed(client) {
		return errors.NewStandardError("please check your password")
	}
	cmdName := strings.ToLower(utils.BytesToString(cmdBytes[0]))
	switch cmdName {
	case "ping":
		return DoPingCmd(cmdBytes[1:])
	case "auth":
		return DoAuthCmd(client, cmdBytes[1:])
	case "info":
		return DoInfoCmd(cmdBytes)
	case "client":
		return DoInfoCmd(cmdBytes)
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
		return nil, errors.NewStandardError("index is out of range")
	}
	return server.dbs[dbIndex].Load().(*DB), nil
}

func isAuthed(client redis.ClientConn) bool {
	pwd := config.GetPassword()
	if pwd == "" {
		return true
	}
	return client.GetPassword() == pwd
}
