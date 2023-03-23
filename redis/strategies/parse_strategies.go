package strategies

import (
	"bufio"
	"goredis/interface/tcp"
	"goredis/pkg/errors"
	"goredis/pkg/utils"
	"goredis/redis/request"
	"io"
	"strconv"
)

// ParseStrategy 行解析策略接口
type ParseStrategy interface {
	Do(reader *bufio.Reader, lineBytes []byte, ch chan<- *tcp.Request) error
}

type Operator struct {
	ParseStrategy ParseStrategy
}

func (operator *Operator) DoStrategy(reader *bufio.Reader, lineBytes []byte, ch chan<- *tcp.Request) error {
	return operator.ParseStrategy.Do(reader, lineBytes, ch)
}

// BulkStringsStrategy 解析多行字符串
type BulkStringsStrategy struct{}

func (*BulkStringsStrategy) Do(reader *bufio.Reader, lineBytes []byte, ch chan<- *tcp.Request) error {
	strLen, err := strconv.ParseInt(string(lineBytes[1:]), 10, 64)
	if err != nil || strLen < -1 {
		handleParseError(utils.NewStringBuilder("illegal bulk strings lineBytes ", string(lineBytes[1:])), ch)
		return nil
	} else if strLen == -1 {
		ch <- &tcp.Request{
			Data: request.NewEmptyMultiBulkRequest(),
		}
	}
	body := make([]byte, strLen+2)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return err
	}
	ch <- &tcp.Request{
		Data: request.NewBulkRequest(body[:len(body)-2]),
	}
	return nil
}

// ArrayStrategy 解析数组
type ArrayStrategy struct{}

func (*ArrayStrategy) Do(reader *bufio.Reader, lineBytes []byte, ch chan<- *tcp.Request) error {
	nStrs, err := strconv.ParseInt(string(lineBytes[1:]), 10, 64)
	if err != nil || nStrs < 0 {
		handleParseError(utils.NewStringBuilder("illegal bulk strings lineBytes ", string(lineBytes[1:])), ch)
		return nil
	} else if nStrs == 0 {
		ch <- &tcp.Request{
			Data: request.NewEmptyMultiBulkRequest(),
		}
		return nil
	}
	lines := make([][]byte, 0, nStrs)
	for i := int64(0); i < nStrs; i++ {
		var line []byte
		line, err = reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		length := len(line)
		if length < 4 || line[length-2] != '\r' || line[0] != '$' {
			handleParseError(utils.NewStringBuilder("illegal bulk strings lineBytes ", string(line)), ch)
			break
		}
		strLen, err := strconv.ParseInt(string(line[1:length-2]), 10, 64)
		if err != nil || strLen < -1 {
			handleParseError(utils.NewStringBuilder("illegal bulk strings length ", string(line)), ch)
			break
		} else if strLen == -1 {
			lines = append(lines, []byte{})
		} else {
			body := make([]byte, strLen+2)
			_, err := io.ReadFull(reader, body)
			if err != nil {
				return err
			}
			lines = append(lines, body[:len(body)-2])
		}
	}
	//解析为多行请求
	ch <- &tcp.Request{
		Data: request.NewMultiBulkRequest(lines),
	}
	return nil
}

// RDBBulkStringsStrategy 解析rdb多行请求
type RDBBulkStringsStrategy struct{}

func (*RDBBulkStringsStrategy) Do(reader *bufio.Reader, lineBytes []byte, ch chan<- *tcp.Request) error {
	return nil
}

// 封装解析异常处理，并返回
func handleParseError(msg string, ch chan<- *tcp.Request) {
	err := errors.NewParseError(msg)
	ch <- &tcp.Request{
		Error: err,
	}
}
