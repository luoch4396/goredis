package handler

import (
	"github.com/go-netty/go-netty"
	"sync/atomic"
)

type EchoHandler struct {
	//role string
}

var (
	maxChan uint32 = 0
)

// HandleActive 开启连接
func (l EchoHandler) HandleActive(ctx netty.ActiveContext) {
	//TODO 控制连接上限，以及连接池管理？
	atomic.AddUint32(&maxChan, 1)
	//var channel = ctx.Channel()
	//if !channel.IsActive() {
	//	return
	//}
}

// HandleInactive 关闭连接
func (l EchoHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	//fmt.Println(l.role, "->", "inactive:", ctx.Channel().RemoteAddr(), ex)
	ctx.Channel().Close(ex)
}
