package server

import (
	"fmt"
	"github.com/go-netty/go-netty"
)

type EchoHandler struct {
	role string
}

func (l EchoHandler) HandleActive(ctx netty.ActiveContext) {
	fmt.Println(l.role, "->", "active:", ctx.Channel().RemoteAddr())

	ctx.Write("Hello I'm " + l.role)
	ctx.HandleActive()
}

func (l EchoHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	fmt.Println(l.role, "->", "handle read:", message)
	ctx.HandleRead(message)
}

func (l EchoHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	fmt.Println(l.role, "->", "inactive:", ctx.Channel().RemoteAddr(), ex)
	ctx.HandleInactive(ex)
}
