package handler

import (
	"bytes"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	Utils "goredis/pkg/utils"
	"strconv"
	"strings"
)

func RedisCodec() codec.Codec {
	return &redisCodec{}
}

type redisCodec struct{}

func (*redisCodec) CodecName() string {
	return "redis-handler"
}

func (*redisCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {
	//go parse(ctx, message)
	parse(ctx, message)
	ctx.Write(message)
}

func (*redisCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	switch s := message.(type) {
	case string:
		ctx.HandleWrite(strings.NewReader(s))
	default:
		ctx.HandleWrite(message)
	}
}

func parse(ctx netty.InboundContext, message netty.Message) {
	var textBytes, err = Utils.ParseMessage(message)
	if err != nil {
		return
	}
	var length = len(textBytes)
	var content1 = string(textBytes[1:])
	println("收到数据", content1)
	for {
		if length <= 2 || textBytes[length-2] != '\r' {
			continue
		}
		textBytes = bytes.TrimSuffix(textBytes, []byte{'\r', '\n'})
		switch textBytes[0] {
		case '+':
			var content = string(textBytes[1:])
			println("收到数据1", content)
		case '-':
			value, err := strconv.ParseInt(string(textBytes[1:]), 10, 64)
			println("收到数据2", value)
			if err != nil {
				continue
			}

		case '$':
			strLen, err := strconv.ParseInt(string(textBytes[1:]), 10, 64)
			println("收到数据3", strLen)
			if err != nil {
				continue
			}
		default:
			args := bytes.Split(textBytes, []byte{' '})
			println("收到数据4", args)
		}

	}

	//RedisUtils.NewStringBuilder0(textBytes)
}
