package db

import (
	"fmt"
	"goredis/redis/config"
	"os"
	"runtime"
	"testing"
	"time"
)

func BenchmarkGetServerInfo(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getServerInfo()
	}
}

func BenchmarkGetServerInfoNoCache(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		startUpTime := time.Since(config.GetStartUpTime()) / time.Second
		fmt.Sprintf("# Server\r\n"+
			//redis版本信息
			"redis_version:%s\r\n"+
			"redis_mode:%s\r\n"+
			"os:%s %s\r\n"+
			"sys_bits:%d\r\n"+
			"go_version:%s\r\n"+
			"process_id:%d\r\n"+
			"redis_port:%d\r\n"+
			"uptime_in_seconds:%d\r\n"+
			"uptime_in_days:%d\r\n",
			config.Version,
			config.GetServerType(),
			runtime.GOOS, runtime.GOARCH,
			32<<(^uint(0)>>63),
			runtime.Version(),
			os.Getpid(),
			config.GetPort(),
			startUpTime,
			startUpTime/time.Duration(3600*24))
	}
}
