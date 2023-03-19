package handler

import (
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	"goredis/interface/tcp"
	"goredis/pkg/log"
	"goredis/pool"
	"strings"
)

func RedisCodec() codec.Codec {
	return &readHandler{}
}

type readHandler struct{}

func (*readHandler) CodecName() string {
	return "read-handler"
}

func (*readHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	ch := make(chan *tcp.Request)
	//TODO 配置化协程池大小
	pools, err := pool.GetInstance(1000)
	if err != nil {
		log.Errorf("run parse message with any error, func exit: ", err)
		return
	}
	err = pools.Pool.Submit(func() {
		parse(message, ch)
	})
	if err != nil {
		log.Errorf("run parse message with any error, func exit: ", err)
		return
	}
	//ctx.Write(message)
}

func (*readHandler) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	switch s := message.(type) {
	case string:
		ctx.HandleWrite(strings.NewReader(s))
	default:
		ctx.HandleWrite(message)
	}
}
