package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	"github.com/panjf2000/ants/v2"
	"goredis/interface/tcp"
	"goredis/pkg/log"
	"goredis/redis/strategies"
	"io"
	"runtime/debug"
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
	ch := make(chan *tcp.Request)
	parse := func() {
		parse(message, ch)
	}
	pool, err := ants.NewPool(5000)
	if err != nil {
		log.Errorf("run parse message with any error, func exit", err)
		return
	}
	err = pool.Submit(parse)
	if err != nil {
		log.Errorf("run parse message with any error, func exit", err)
		return
	}
	//ctx.Write(message)
}

func (*redisCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	switch s := message.(type) {
	case string:
		ctx.HandleWrite(strings.NewReader(s))
	default:
		ctx.HandleWrite(message)
	}
}

// 根据RESP解析为统一格式返回
func parse(message netty.Message, ch chan<- *tcp.Request) {
	defer func() {
		if err := recover(); err != nil {
			err = fmt.Errorf(string(debug.Stack()), err)
		}
	}()
	t, ok := message.(io.Reader)
	if !ok {
		var err = fmt.Errorf("message codec produce any error")
		ch <- &tcp.Request{
			Error: err,
		}
		close(ch)
		return
	}
	var reader = bufio.NewReader(t)
	for {
		lineBytes, err := reader.ReadBytes('\n')
		if err != nil {
			ch <- &tcp.Request{
				Error: err,
			}
			close(ch)
			return
		}
		var length = len(lineBytes)
		if length <= 2 || lineBytes[length-2] != '\r' {
			continue
		}
		//需要先把字符串的'\r', '\n'去掉
		lineBytes = bytes.TrimSuffix(lineBytes, []byte{'\r', '\n'})
		switch lineBytes[0] {
		//单行字符串（Simple Strings）： 响应的首字节是 "+"
		case '+':
			var content = string(lineBytes[1:])
			println("收到数据1", content)
		//错误（Errors）： 响应的首字节是 "-"
		case '-':
			value, err := strconv.ParseInt(string(lineBytes[1:]), 10, 64)
			println("收到数据2", value)
			if err != nil {
				_ = fmt.Errorf(string(debug.Stack()), err)
				continue
			}
		//多行字符串（Bulk Strings）： 响应的首字节是"\$"
		case '$':
			//println("解析前:", string(lineBytes))
			var operator = strategies.Operator{
				ParseStrategy: &strategies.BulkStringsStrategy{},
			}
			var operatorRequest = operator.DoStrategy(reader, lineBytes)
			//println("解析后:", string(operatorResponse.Data.RequestInfo()))
			ch <- operatorRequest
		case '*':
			//err = parseArray(lineBytes, reader)
			if err != nil {
				//return requests, err
			}
		//整型（Integers）： 响应的首字节是 ":"
		case ':':
			value, err := strconv.ParseInt(string(lineBytes[1:]), 10, 64)
			if err != nil {
				_ = fmt.Errorf(string(debug.Stack()), err)
				continue
			}
			println(value)
			//
		default:
			var args = bytes.Split(lineBytes, []byte{' '})
			println("收到数据4", args)
		}
	}
}
