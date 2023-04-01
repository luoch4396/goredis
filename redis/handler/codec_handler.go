package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	"goredis/interface/tcp"
	"goredis/pkg/errors"
	"goredis/pkg/log"
	"goredis/pkg/pool/gopool"
	"goredis/redis/exchange"
	"io"
	"strconv"
	"strings"
)

var (
	//定义解析策略
	arrayStrategy       = NewParseOperator(&ArrayStrategy{})
	bulkStringsStrategy = NewParseOperator(&BulkStringsStrategy{})
	//暂不支持的命令/未知命令
	unknownOperation = []byte("-ERR unknown\r\n")
)

func RedisCodec() codec.Codec {
	return &codecHandler{}
}

type codecHandler struct{}

func (*codecHandler) CodecName() string {
	return "codec-handler"
}

func (*codecHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	var handleReadFunc = func() {
		handleRead(ctx, message)
	}
	gopool.Go(handleReadFunc)
}

func handleRead(ctx netty.InboundContext, message netty.Message) {
	ch := make(chan *tcp.Request)
	var parseStreamingFunc = func() {
		parseStreaming(message, ch)
	}
	gopool.Go(parseStreamingFunc)
	//循环结果
	for req := range ch {
		if req.Error != nil {
			if req.Error == io.EOF || req.Error == io.ErrUnexpectedEOF ||
				strings.Contains(req.Error.Error(), "use a closed network channel") {
				log.Errorf("handle message with errors, channel will be closed: " + ctx.Channel().RemoteAddr())
				ctx.Channel().Close(req.Error)
				return
			}
			errRep := errors.NewStandardError(req.Error.Error())
			ctx.Write(errRep.Info())
			continue
		}
		if req.Data == nil {
			log.Error("empty commands")
			continue
		}
		_, ok := req.Data.(*exchange.MultiBulkRequest)
		if !ok {
			log.Error("error from multi bulk exchange")
			continue
		}
		//命令处理
		//result := h.db.Exec(message, r.Args)
		//if result != nil {
		//	ctx.Write(result.ToBytes())
		//} else {
		//	ctx.Write(unknownOperation)
		//}
	}
}

func (*codecHandler) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	//switch s := message.(type) {
	//case string:
	//	ctx.HandleWrite(strings.NewReader(s))
	//default:
	//	ctx.HandleWrite(message)
	//}
}

// 根据RESP解析为统一格式返回
func parseStreaming(message netty.Message, ch chan<- *tcp.Request) {
	//TODO:错误恢复 移动至协程池？
	//defer func() {
	//	//错误恢复
	//	if err := recover(); err != nil {
	//		err = fmt.Errorf(string(debug.Stack()), err)
	//	}
	//}()
	t, ok := message.(io.Reader)
	if !ok {
		var err = fmt.Errorf("message codec produce any errors")
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
			ch <- &tcp.Request{
				Data: exchange.NewStatusInfo(content),
			}
			//TODO rdb操作
		//错误（Errors）： 响应的首字节是 "-"
		case '-':
			ch <- &tcp.Request{
				Data: errors.NewStandardError(string(lineBytes[1:])),
			}
		//多行字符串（Bulk Strings）： 响应的首字节是"\$"
		case '$':
			err := bulkStringsStrategy.DoParseStrategy(reader, lineBytes, ch)
			if err != nil {
				close(ch)
				return
			}
		case '*':
			err := arrayStrategy.DoParseStrategy(reader, lineBytes, ch)
			if err != nil {
				close(ch)
				return
			}
		//整型（Integers）： 响应的首字节是 ":"
		case ':':
			value, err := strconv.ParseInt(string(lineBytes[1:]), 10, 64)
			if err != nil {
				log.Errorf("", err)
				continue
			}
			ch <- &tcp.Request{
				Data: exchange.NewIntRequest(value),
			}
			//
		default:
			var args = bytes.Split(lineBytes, []byte{' '})
			ch <- &tcp.Request{
				Data: exchange.NewMultiBulkRequest(args),
			}
		}
	}
}
