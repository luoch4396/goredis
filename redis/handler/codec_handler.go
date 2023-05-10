package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	"goredis/interface/redis"
	"goredis/interface/tcp"
	"goredis/pkg/errors"
	"goredis/pkg/log"
	"goredis/pkg/pool/gopool"
	"goredis/pkg/utils"
	"goredis/redis/conn"
	"goredis/redis/exchange"
	"io"
	"strconv"
	"strings"
	"sync"
)

var (
	//定义解析策略
	arrayStrategy       = NewParseOperator(&ArrayStrategy{})
	bulkStringsStrategy = NewParseOperator(&BulkStringsStrategy{})
	//暂不支持的命令/未知命令
	unknownOperation = utils.StringToBytes("-ERR unknown\r\n")
)

func NewRedisCodec(server redis.Server) codec.Codec {
	return &codecHandler{
		server: server,
	}
}

type codecHandler struct {
	activeConn sync.Map
	server     redis.Server
}

func (*codecHandler) CodecName() string {
	return "codec-handler"
}

func (c *codecHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	//包装连接对象
	client := conn.NewClientConnBuilder().BuildChannel(ctx.Channel()).Build()
	c.activeConn.Store(client, struct{}{})
	//命令异步处理
	ch := make(chan *tcp.Request)
	var parseStreamingFunc = func() {
		parseStreaming(message, ch)
	}
	gopool.Go(parseStreamingFunc)
	//循环结果
	for req := range ch {
		err := req.Error
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF ||
				strings.Contains(err.Error(), "use of closed network connection") {
				log.Errorf("handle message with errors: %s, channel will be closed: %s", err.Error(),
					ctx.Channel().RemoteAddr())
				ctx.Channel().Close(err)
				return
			}
			errRep := errors.NewStandardError(err.Error())
			ctx.Write(errRep.Info())
			continue
		}
		if req.Data == nil {
			log.Error("empty commands")
			continue
		}
		r, ok := req.Data.(*exchange.MultiBulkRequest)
		if !ok {
			log.Error("error from multi bulk exchange")
			continue
		}
		//处理解析后的命令
		result := c.server.Exec(client, r.Args)
		if result != nil {
			ctx.Write(result.Info())
		} else {
			ctx.Write(unknownOperation)
		}
	}
}

func (*codecHandler) HandleWrite(ctx netty.OutboundContext, message netty.Message) {

}

// 根据RESP解析为统一格式返回
func parseStreaming(message netty.Message, ch chan<- *tcp.Request) {
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
			var content = utils.BytesToString(lineBytes[1:])
			ch <- &tcp.Request{
				Data: exchange.NewStatusInfo(content),
			}
			//TODO rdb操作
		//错误（Errors）： 响应的首字节是 "-"
		case '-':
			ch <- &tcp.Request{
				Data: errors.NewStandardError(utils.BytesToString(lineBytes[1:])),
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
			value, err := strconv.ParseInt(utils.BytesToString(lineBytes[1:]), 10, 64)
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
