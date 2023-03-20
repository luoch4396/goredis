package request

import (
	"goredis/pkg/utils"
	"strconv"
)

// BulkRequest 二进制指令
type BulkRequest struct {
	Arg []byte
}

var (
	nullBulkBytes       = []byte("$-1\r\n")
	CRLF                = "\r\n"
	emptyMultiBulkBytes = []byte("*0\r\n")
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

// MultiBulkRequest 多行二进制指令
type MultiBulkRequest struct {
	Args [][]byte
}

func NewMultiBulkRequest(args [][]byte) *MultiBulkRequest {
	return &MultiBulkRequest{
		Args: args,
	}
}

func (r *MultiBulkRequest) RequestInfo() []byte {
	return nil
}

// StatusRequest 状态指令
type StatusRequest struct {
	Status string
}

func NewStatusRequest(status string) *StatusRequest {
	return &StatusRequest{
		Status: status,
	}
}

func (r *StatusRequest) RequestInfo() []byte {
	return nil
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

func (r *IntRequest) RequestInfo() []byte {
	return []byte(":" + strconv.FormatInt(r.Code, 10) + CRLF)
}

// EmptyBulkRequest 空字符串
type EmptyBulkRequest struct{}

func (r *EmptyBulkRequest) RequestInfo() []byte {
	return nullBulkBytes
}

func NewNullBulkRequest() *EmptyBulkRequest {
	return &EmptyBulkRequest{}
}

// EmptyMultiBulkRequest 多行空字符串
type EmptyMultiBulkRequest struct{}

func (r *EmptyMultiBulkRequest) RequestInfo() []byte {
	return emptyMultiBulkBytes
}

func NewEmptyMultiBulkRequest() *EmptyBulkRequest {
	return &EmptyBulkRequest{}
}

// StandardErrorRequest 错误指令
type StandardErrorRequest struct {
	Status string
}

func NewStandardErrorRequest(status string) *StandardErrorRequest {
	return &StandardErrorRequest{
		Status: status,
	}
}

func (r *StandardErrorRequest) RequestInfo() []byte {
	return []byte(utils.NewStringBuilder("-", r.Status, CRLF))
}

func (r *StandardErrorRequest) Error() string {
	return r.Status
}
