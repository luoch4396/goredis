package nio

import (
	"goredis/pkg/log"
	"goredis/pkg/pool/bytepool"
	"time"
)

// TcpConfigs TCP配置
type TcpConfigs struct {
	MaxOpenFiles int
	NumPolls     int
	Lb           *LoadBalancer
	KeepAlive    *time.Duration
	NoDelay      bool
	RecvBuffer   int
	SendBuffer   int
	Log          log.FormatterLogger
	Allocator    bytepool.Allocator
}

//var (
//	MaxOpenFiles = 1024 * 1024 * 2
//)
