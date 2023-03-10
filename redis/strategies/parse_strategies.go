package strategies

import (
	"bufio"
	"errors"
	"goredis/interface/tcp"
	"goredis/pkg/utils"
	"goredis/redis/request"
	"io"
	"strconv"
)

// ParseStrategy 解析策略接口
type ParseStrategy interface {
	do(reader *bufio.Reader, lineBytes []byte) *tcp.Request
}

type Operator struct {
	ParseStrategy ParseStrategy
}

func (operator *Operator) DoStrategy(reader *bufio.Reader, lineBytes []byte) *tcp.Request {
	return operator.ParseStrategy.do(reader, lineBytes)
}

// BulkStringsStrategy 解析多行字符串
type BulkStringsStrategy struct{}

func (*BulkStringsStrategy) do(reader *bufio.Reader, lineBytes []byte) *tcp.Request {
	strLen, err := strconv.ParseInt(string(lineBytes[1:]), 10, 64)
	var redisRequest = &tcp.Request{}
	if err != nil || strLen < -1 {
		err := errors.New(utils.NewStringBuilder("protocol error: illegal bulk string header: ", string(lineBytes)))
		redisRequest.Error = err
		return redisRequest
	} else if strLen == -1 {
		//TODO 返回为空
		return redisRequest
	}
	body := make([]byte, strLen+2)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		redisRequest.Error = err
		return redisRequest
	}
	redisRequest.Data = request.NewBulkRequest(body[:len(body)-2])
	return redisRequest
}
