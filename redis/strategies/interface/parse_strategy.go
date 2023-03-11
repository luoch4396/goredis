package _interface

import (
	"bufio"
	"goredis/interface/tcp"
)

// ParseStrategy 解析策略接口
type ParseStrategy interface {
	Do(reader *bufio.Reader, lineBytes []byte) *tcp.Request
}
