package redis

import (
	"github.com/go-netty/go-netty"
	"goredis/db"
	"goredis/pkg/log"
	"goredis/redis/handler"
	"time"
)

type Config struct {
	Address  string        `config:"address"`
	MaxConns uint32        `config:"max-conns"`
	Timeout  time.Duration `config:"timeout"`
}

// NewRedisServer 实现一个netty redis
func NewRedisServer(config *Config) {
	// setup child pipeline initializer.
	var childInitializer = func(channel netty.Channel) {
		channel.Pipeline().
			AddLast(handler.EchoHandler{}).
			AddLast(handler.RedisCodec())

	}
	// new bootstrap
	var bootstrap = netty.NewBootstrap(netty.WithChildInitializer(childInitializer))
	// setup bootstrap & startup redis.
	println("redis服务启动地址:", config.Address)
	var listener = bootstrap.Listen(config.Address)
	var err = listener.Sync()
	if err != nil {
		var err = listener.Close()
		if err != nil {
			log.Error("", err)
			return
		}
	}
	//todo 后期增加cluster 模式，现在仅有单机模式
	db.NewSingleServer()
}
