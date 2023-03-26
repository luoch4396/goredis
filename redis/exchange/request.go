package exchange

import (
	"bytes"
	"goredis/pkg/utils"
	"strconv"
)

// BulkRequest 二进制指令
type BulkRequest struct {
	Arg []byte
}

func NewBulkRequest(arg []byte) *BulkRequest {
	return &BulkRequest{
		Arg: arg,
	}
}

func (r *BulkRequest) Info() []byte {
	if r.Arg == nil {
		return nullBulkBytes
	}
	//for example:
	//$(命令长度)
	//命令具体参数：QUIT
	return []byte(utils.NewStringBuilder("$", strconv.Itoa(len(r.Arg)), CRLF, string(r.Arg), CRLF))
}

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
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
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
	return []byte(":" + strconv.FormatInt(r.Code, 10) + CRLF)
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
