package handler

import (
	"bytes"
	"goredis/interface/redis"
	"goredis/interface/tcp"
	"goredis/pkg/errors"
	"goredis/pkg/log"
	"goredis/pkg/nio"
	"goredis/pkg/pool/gopool"
	"goredis/pkg/utils"
	"goredis/redis/conn"
	"goredis/redis/exchange"
	"io"
	"strconv"
	"strings"
)

func Handle(c *nio.Conn, data []byte, server redis.Server) {
	client := conn.NewClientConnBuilder().Build()
	//命令异步处理
	ch := make(chan *tcp.Request)
	var parseStreamingFunc = func() {
		noCopyParse(data, ch)
	}
	gopool.Go(parseStreamingFunc)
	//循环结果
	for req := range ch {
		err := req.Error
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF ||
				strings.Contains(err.Error(), "use of closed network connection") {
				log.Errorf("handle message with errors: %s, channel will be closed: %s", err.Error(),
					c.RemoteAddr())
				c.Close()
				return
			}
			errRep := errors.NewStandardError(err.Error())
			c.Write(errRep.Info())
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
		result := server.Exec(client, r.Args)
		if result != nil {
			c.Write(result.Info())
		} else {
			c.Write(unknownOperation)
		}
	}
}

// 根据RESP解析为统一格式返回
func noCopyParse(data []byte, ch chan<- *tcp.Request) {
	//var reader = bufio.NewReader(bytes.NewReader(data))
	strs := strings.Split(utils.BytesToString(data), "\n")
	log.Debugf(utils.BytesToString(data))
	for i := 0; i < len(strs); i++ {
		var str = strs[i]
		if str == "" {
			continue
		}
		//需要先把字符串的'\r', '\n'去掉
		str = strings.TrimSuffix(str, "\r")
		str = strings.TrimSuffix(str, "\n")
		//log.Debug(str)
		switch str[0] {
		//单行字符串（Simple Strings）： 响应的首字节是 "+"
		case '+':
			ch <- &tcp.Request{
				Data: exchange.NewStatusInfo(str),
			}
			//TODO rdb操作
		//错误（Errors）： 响应的首字节是 "-"
		case '-':
			ch <- &tcp.Request{
				Data: errors.NewStandardError(str),
			}
		//多行字符串（Bulk Strings）： 响应的首字节是"\$"
		case '$':
			//log.Debugf(str)
			//err := bulkStringsStrategy.DoParseStrategy(reader, lineBytes, ch)
			//if err != nil {
			//	close(ch)
			//	return
			//}
		case '*':
			//err := arrayStrategy.DoParseStrategy(reader, lineBytes, ch)
			//if err != nil {
			//	close(ch)
			//	return
			//}
			//log.Debugf(str)
		//整型（Integers）： 响应的首字节是 ":"
		case ':':
			value, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				log.Errorf("", err)
				continue
			}
			ch <- &tcp.Request{
				Data: exchange.NewIntRequest(value),
			}
		default:
			var args = bytes.Split(utils.StringToBytes(str), []byte{' '})
			ch <- &tcp.Request{
				Data: exchange.NewMultiBulkRequest(args),
			}
		}
	}
	//for {
	//	lineBytes, err := reader.ReadBytes('\n')
	//	if err != nil {
	//		ch <- &tcp.Request{
	//			Error: err,
	//		}
	//		close(ch)
	//		return
	//	}
	//	var length = len(lineBytes)
	//	if length <= 2 || lineBytes[length-2] != '\r' {
	//		continue
	//	}
	//	//需要先把字符串的'\r', '\n'去掉
	//	lineBytes = bytes.TrimSuffix(lineBytes, []byte{'\r', '\n'})
	//	switch lineBytes[0] {
	//	//单行字符串（Simple Strings）： 响应的首字节是 "+"
	//	case '+':
	//		var content = utils.BytesToString(lineBytes[1:])
	//		ch <- &tcp.Request{
	//			Data: exchange.NewStatusInfo(content),
	//		}
	//		//TODO rdb操作
	//	//错误（Errors）： 响应的首字节是 "-"
	//	case '-':
	//		ch <- &tcp.Request{
	//			Data: errors.NewStandardError(utils.BytesToString(lineBytes[1:])),
	//		}
	//	//多行字符串（Bulk Strings）： 响应的首字节是"\$"
	//	case '$':
	//		err := bulkStringsStrategy.DoParseStrategy(reader, lineBytes, ch)
	//		if err != nil {
	//			close(ch)
	//			return
	//		}
	//	case '*':
	//		err := arrayStrategy.DoParseStrategy(reader, lineBytes, ch)
	//		if err != nil {
	//			close(ch)
	//			return
	//		}
	//	//整型（Integers）： 响应的首字节是 ":"
	//	case ':':
	//		value, err := strconv.ParseInt(utils.BytesToString(lineBytes[1:]), 10, 64)
	//		if err != nil {
	//			log.Errorf("", err)
	//			continue
	//		}
	//		ch <- &tcp.Request{
	//			Data: exchange.NewIntRequest(value),
	//		}
	//		//
	//	default:
	//		var args = bytes.Split(lineBytes, []byte{' '})
	//		ch <- &tcp.Request{
	//			Data: exchange.NewMultiBulkRequest(args),
	//		}
	//	}
	//}
}
