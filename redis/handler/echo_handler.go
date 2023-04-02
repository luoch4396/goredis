package handler

import (
	"github.com/go-netty/go-netty"
	"goredis/pkg/errors"
	"goredis/redis/config"
	"sync/atomic"
)

type EchoHandler struct {
	maxChan int32
}

// HandleActive 开启连接
func (l EchoHandler) HandleActive(ctx netty.ActiveContext) {
	atomic.AddInt32(&l.maxChan, 1)
	if l.maxChan > int32(config.GetMaxConn()) {
		ctx.Channel().Close(errors.NewStandardError("the number of connections more than max-conn"))
	}
}

// HandleInactive 关闭连接
func (l EchoHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	atomic.AddInt32(&l.maxChan, -1)
}
