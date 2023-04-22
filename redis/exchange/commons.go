package exchange

import (
	"goredis/pkg/utils"
	"strconv"
)

var (
	nullBulkBytes       = utils.StringToBytes("$-1\r\n")
	CRLF                = "\r\n"
	emptyMultiBulkBytes = utils.StringToBytes("*0\r\n")
)

// StatusInfo 状态指令
type StatusInfo struct {
	Status string
}

func NewStatusInfo(status string) *StatusInfo {
	return &StatusInfo{
		Status: status,
	}
}

func (r *StatusInfo) Info() []byte {
	return utils.StringToBytes("+" + r.Status + CRLF)
}

// BulkInfo 二进制指令
type BulkInfo struct {
	Arg []byte
}

func NewBulkInfo(arg []byte) *BulkInfo {
	return &BulkInfo{
		Arg: arg,
	}
}

func (r *BulkInfo) Info() []byte {
	if r.Arg == nil {
		return nullBulkBytes
	}
	//for example:
	//$(命令长度)
	//命令具体参数：QUIT
	return utils.StringToBytes(utils.NewStringBuilder(
		"$", strconv.Itoa(len(r.Arg)), CRLF, utils.BytesToString(r.Arg), CRLF))
}
