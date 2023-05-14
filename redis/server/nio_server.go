package server

import (
	"goredis/db"
	"goredis/interface/redis"
	"goredis/pkg/log"
	"goredis/pkg/nio"
	"goredis/pkg/pool/gopool"
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

func NewNioServer(config *Config, server redis.Server) {
	engine := nio.NewEngine(nio.Config{
		Network:            "tcp",
		Addrs:              []string{config.Address},
		MaxWriteBufferSize: 128 * 1024 * 1024,
		LockListener:       true,
	})

	//engine.BeforeWrite(func(c *nio.Conn) {
	//	c.SetWriteDeadline(time.Now().Add(time.Second * 60))
	//})

	engine.OnOpen(func(c *nio.Conn) {
		//log.Debugf("OnOpen:", c.RemoteAddr().String())
	})

	engine.OnClose(func(c *nio.Conn, err error) {
		//log.Warningf("`errors: %s, channel will be closed: %s", err, c.RemoteAddr().String())
	})

	engine.OnData(func(c *nio.Conn, data []byte) {
		gopool.HandleGo(func() {
			handler.Handle(c, data, server)
		})
	})

	err := engine.Start()
	if err != nil {
		log.Fatalf("goredis nio server start failed: %v\n", err)
	}

	<-make(chan int)
}
