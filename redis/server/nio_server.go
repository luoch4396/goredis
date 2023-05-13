package server

import (
	"goredis/interface/redis"
	"goredis/pkg/log"
	"goredis/pkg/nio"
	"goredis/redis/handler"
)

func NewNioServer(config *Config, server redis.Server) {
	engine := nio.NewEngine(nio.Config{
		Network:            "tcp",
		Addrs:              []string{config.Address},
		MaxWriteBufferSize: 128 * 1024 * 1024,
	})

	//engine.BeforeWrite(func(c *nio.Conn) {
	//	c.SetWriteDeadline(time.Now().Add(time.Second * 5))
	//})

	engine.OnOpen(func(c *nio.Conn) {
		log.Debugf("OnOpen:", c.RemoteAddr().String())
	})

	engine.OnClose(func(c *nio.Conn, err error) {
		log.Errorf("handle message with errors: %s, channel will be closed: %s", err, c.RemoteAddr())
	})

	engine.OnData(func(c *nio.Conn, data []byte) {
		handler.Handle(c, data, server)
	})

	err := engine.Start()
	if err != nil {
		log.Fatalf("goredis nio server start failed: %v\n", err)
	}

	<-make(chan int)
}
