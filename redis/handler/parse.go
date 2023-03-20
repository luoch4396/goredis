package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-netty/go-netty"
	"goredis/interface/tcp"
	"goredis/pkg/log"
	"goredis/redis/request"
	"goredis/redis/strategies"
	"io"
	"runtime/debug"
	"strconv"
)

// 根据RESP解析为统一格式返回
func parse(message netty.Message, ch chan<- *tcp.Request) {
	defer func() {
		//错误恢复
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
			//TODO rdb操作
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
			operator := strategies.Operator{
				ParseStrategy: &strategies.BulkStringsStrategy{},
			}
			operatorRequest := operator.DoStrategy(reader, lineBytes)
			if operatorRequest != nil {
				//log.Debug("解析后: " + string(operatorRequest.Data.RequestInfo()))
				ch <- operatorRequest
				if operatorRequest.Error != nil {
					close(ch)
					return
				}
			}
		case '*':
			operator := strategies.Operator{
				ParseStrategy: &strategies.ArrayStrategy{},
			}
			operatorRequest := operator.DoStrategy(reader, lineBytes)
			if operatorRequest != nil {
				//log.Debug("解析后: " + string(operatorRequest.Data.RequestInfo()))
				ch <- operatorRequest
				if operatorRequest.Error != nil {
					close(ch)
					return
				}
			}
		//整型（Integers）： 响应的首字节是 ":"
		case ':':
			value, err := strconv.ParseInt(string(lineBytes[1:]), 10, 64)
			if err != nil {
				log.Errorf("", err)
				continue
			}
			ch <- &tcp.Request{
				Data: request.NewIntRequest(value),
			}
			//
		default:
			var args = bytes.Split(lineBytes, []byte{' '})
			println("收到数据4", args)
		}
	}
}
