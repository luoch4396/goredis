package exchange

import (
	"bytes"
	"goredis/pkg/utils"
	"strconv"
)

// MultiBulkRequest 多行二进制指令
type MultiBulkRequest struct {
	Args [][]byte
}

func NewMultiBulkRequest(args [][]byte) *MultiBulkRequest {
	return &MultiBulkRequest{
		Args: args,
	}
}

func (r *MultiBulkRequest) Info() []byte {
	argLen := len(r.Args)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argLen) + CRLF)
	for _, arg := range r.Args {
		if arg == nil {
			buf.WriteString("$-1" + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + utils.BytesToString(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

// IntRequest 解析为数值
type IntRequest struct {
	Code int64
}

func NewIntRequest(code int64) *IntRequest {
	return &IntRequest{
		Code: code,
	}
}

func (r *IntRequest) Info() []byte {
	return utils.StringToBytes(":" + strconv.FormatInt(r.Code, 10) + CRLF)
}

// EmptyBulkRequest 空字符串
type EmptyBulkRequest struct{}

func (r *EmptyBulkRequest) Info() []byte {
	return nullBulkBytes
}

func NewNullBulkRequest() *EmptyBulkRequest {
	return &EmptyBulkRequest{}
}

// EmptyMultiBulkRequest 多行空字符串
type EmptyMultiBulkRequest struct{}

func (r *EmptyMultiBulkRequest) Info() []byte {
	return emptyMultiBulkBytes
}

func NewEmptyMultiBulkRequest() *EmptyBulkRequest {
	return &EmptyBulkRequest{}
}
