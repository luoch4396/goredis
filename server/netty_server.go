package server

import (
	"encoding/binary"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec/format"
	"github.com/go-netty/go-netty/codec/frame"
	"github.com/go-netty/go-netty/utils"
	"time"
)

type Config struct {
	Address  string        `config:"address"`
	MaxConns uint16        `config:"max-conns"`
	Timeout  time.Duration `config:"timeout"`
}

// NewNettyServer 实现一个netty server
func NewNettyServer(config *Config) {
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
	err := bootstrap.Listen(config.Address).Sync()
	utils.Assert(err)
}
