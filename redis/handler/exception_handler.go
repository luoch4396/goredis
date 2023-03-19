package handler

import (
	"github.com/go-netty/go-netty"
	"goredis/pkg/log"
)

type ExceptionHandler struct {
	//role string
}

func (l EchoHandler) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	log.Errorf("", ex)

	//TODO 所有异常都走到这里？
	ctx.Channel().Close(ex)
}
