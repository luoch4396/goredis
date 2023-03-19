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

// ParseStrategy 行解析策略接口
type ParseStrategy interface {
	Do(reader *bufio.Reader, lineBytes []byte) *tcp.Request
}

type Operator struct {
	ParseStrategy ParseStrategy
}

func (operator *Operator) DoStrategy(reader *bufio.Reader, lineBytes []byte) *tcp.Request {
	return operator.ParseStrategy.Do(reader, lineBytes)
}

// BulkStringsStrategy 解析多行字符串
type BulkStringsStrategy struct{}

func (*BulkStringsStrategy) Do(reader *bufio.Reader, lineBytes []byte) *tcp.Request {
	strLen, err := strconv.ParseInt(string(lineBytes[1:]), 10, 64)
	var redisRequest = &tcp.Request{}
	if err != nil || strLen < -1 {
		err := errors.New(utils.NewStringBuilder("parse error: illegal bulk string lineBytes: ", string(lineBytes)))
		redisRequest.Error = err
		return redisRequest
	} else if strLen == -1 {
		return nil
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

// ArrayStrategy 解析数组
type ArrayStrategy struct{}

func (*ArrayStrategy) Do(reader *bufio.Reader, lineBytes []byte) *tcp.Request {
	nStrs, err := strconv.ParseInt(string(lineBytes[1:]), 10, 64)
	if err != nil || nStrs < 0 {
		//protocolError(ch, "illegal array header "+string(header[1:]))
		return nil
	} else if nStrs == 0 {
		//ch <- &Payload{
		//Data: protocol.MakeEmptyMultiBulkReply(),
		//}
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
			//protocolError(ch, "illegal bulk string header "+string(line))
			break
		}
		strLen, err := strconv.ParseInt(string(line[1:length-2]), 10, 64)
		if err != nil || strLen < -1 {
			//protocolError(ch, "illegal bulk string length "+string(line))
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
	//ch <- &Payload{
	//	Data: protocol.MakeMultiBulkReply(lines),
	//}
	return nil
}

type StatusStrategy struct{}

func (*StatusStrategy) Do(reader *bufio.Reader, lineBytes []byte) *tcp.Request {
	return nil
}
