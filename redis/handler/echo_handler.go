package handler

import (
	"fmt"
	"github.com/go-netty/go-netty"
)

type EchoHandler struct {
	role string
}

// HandleActive 开启处理器
func (l EchoHandler) HandleActive(ctx netty.ActiveContext) {
	var channel = ctx.Channel()
	if !channel.IsActive() {
		return
	}
}

// HandleInactive 关闭处理器
func (l EchoHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	fmt.Println(l.role, "->", "inactive:", ctx.Channel().RemoteAddr(), ex)
}
