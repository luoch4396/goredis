package server

import (
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/transport/tcp"
	"goredis/db"
	"goredis/interface/redis"
	"goredis/pkg/log"
	"goredis/redis/handler"
	"time"
)

type Config struct {
	Address  string        `config:"address"`
	MaxConns uint32        `config:"max-conns"`
	Timeout  time.Duration `config:"timeout"`
}

// NewRedisDB 创建redis数据库，通过处理器去调用db执行操作
func NewRedisDB() redis.Server {
	//todo 后期增加cluster 模式，现在仅有单机模式
	return db.NewSingleServer()
}

// NewRedisServer 实现一个netty redis
func NewRedisServer(config *Config, server redis.Server) {
	var childInitializer = func(channel netty.Channel) {
		channel.Pipeline().
			AddLast(handler.EchoHandler{}).
			AddLast(handler.ExceptionHandler{}).
			AddLast(handler.NewRedisCodec(server))
	}
	//TODO 需要控制TCP连接数
	var bootstrap = netty.NewBootstrap(netty.WithChildInitializer(childInitializer))
	log.Info("start goredis server success: " + config.Address + ", start listening...")
	err := bootstrap.Listen(config.Address, tcp.WithOptions(newTcpOp())).Sync()
	if err != nil {
		panic(err)
		return
	}
}

// TCP配置初始化 TODO 改为配置化
func newTcpOp() *tcp.Options {
	return &tcp.Options{
		Timeout:         time.Second * 5,
		KeepAlive:       true,
		KeepAlivePeriod: time.Second * 60,
		Linger:          0,
		NoDelay:         true,
		SockBuf:         2048,
	}
}
