package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-netty/go-netty"
	"goredis/interface/tcp"
	"goredis/pkg/errors"
	"goredis/pkg/log"
	"goredis/redis/exchange"
	"goredis/redis/strategies"
	"io"
	"strconv"
)

// 根据RESP解析为统一格式返回
func parseStreaming(message netty.Message, ch chan<- *tcp.Request) {
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
			operator := strategies.ParseOperator{
				ParseStrategy: &strategies.BulkStringsStrategy{},
			}
			err := operator.DoParseStrategy(reader, lineBytes, ch)
			if err != nil {
				close(ch)
				return
			}
		case '*':
			operator := strategies.ParseOperator{
				ParseStrategy: &strategies.ArrayStrategy{},
			}
			err := operator.DoParseStrategy(reader, lineBytes, ch)
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
