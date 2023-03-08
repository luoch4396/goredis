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
	//fmt.Println(l.role, "->", "active:", ctx.Channel().RemoteAddr())
	//
	//ctx.Write([]byte("-ERR unknown\r\n"))
}

// HandleRead 读取数据处理器
func (l EchoHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	var channel = ctx.Channel()
	if !channel.IsActive() {
		return
	}
	//ctx.HandleRead(message)
}

// HandleWrite 发送数据处理器
func (l EchoHandler) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	fmt.Println(l.role, "->", "handle read:", message)
	//ctx.HandleRead(message)
	//ctx.Write([]byte("-ERR unknown\r\n"))
}

// HandleInactive 关闭处理器
func (l EchoHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	fmt.Println(l.role, "->", "inactive:", ctx.Channel().RemoteAddr(), ex)
	//ctx.HandleInactive(ex)
	//ctx.Write([]byte("-ERR unknown\r\n"))
}
