package server

import (
	"encoding/binary"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec/format"
	"github.com/go-netty/go-netty/codec/frame"
	"goredis/db"
	"goredis/pkg/log"
	"time"
)

type Config struct {
	Address  string        `config:"address"`
	MaxConns uint32        `config:"max-conns"`
	Timeout  time.Duration `config:"timeout"`
}

// NewRedisServer 实现一个netty server
func NewRedisServer(config *Config) {
	// setup child pipeline initializer.
	childInitializer := func(channel netty.Channel) {
		channel.Pipeline().
			AddLast(frame.LengthFieldCodec(binary.LittleEndian, 1024, 0, 2, 0, 2)).
			AddLast(format.TextCodec()).
			AddLast(EchoHandler{"Server"})
	}
	// new bootstrap
	var bootstrap = netty.NewBootstrap(netty.WithChildInitializer(childInitializer))
	// setup bootstrap & startup server.
	var listener = bootstrap.Listen(config.Address)
	err := listener.Sync()
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
