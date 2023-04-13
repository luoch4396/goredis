package handler

import (
	"github.com/go-netty/go-netty"
	"goredis/pkg/errors"
	"goredis/pkg/log"
)

type ExceptionHandler struct {
	//role string
}

func (l EchoHandler) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	log.Errorf("", ex)
	ctx.Channel().Write(errors.NewStandardError(ex.Error()))
	ctx.Channel().Close(ex)
}
