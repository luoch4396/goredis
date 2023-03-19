package request

import (
	"goredis/pkg/utils"
	"strconv"
)

type BulkRequest struct {
	Arg []byte
}

var (
	nullBulkBytes = []byte("$-1\r\n")
	CRLF          = "\r\n"
)

func NewBulkRequest(arg []byte) *BulkRequest {
	return &BulkRequest{
		Arg: arg,
	}
}

func (r *BulkRequest) RequestInfo() []byte {
	if r.Arg == nil {
		return nullBulkBytes
	}
	//for example:
	//$(命令长度)
	//命令具体参数：QUIT
	return []byte(utils.NewStringBuilder("$", strconv.Itoa(len(r.Arg)), CRLF, string(r.Arg), CRLF))
}
