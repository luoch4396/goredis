package monitor

import (
	"runtime"
	"runtime/debug"
	"strconv"
	"syscall"
)

var cache = make(map[string]string, 8)

func init() {
	//pid
	cache["pid"] = strconv.Itoa(syscall.Getpid())
}

// GetRuntimeNumThreads go-runtime的线程数
func GetRuntimeNumThreads() int {
	n, _ := runtime.ThreadCreateProfile(nil)
	return n
}

// GetPid 进程id
func GetPid() string {
	return cache["pid"]
}

// GetGOVersion 获取go版本号
func GetGOVersion() string {
	return runtime.Version()
}

// GetGCInfo 获取GC信息
func GetGCInfo() *debug.GCStats {
	var d debug.GCStats
	debug.ReadGCStats(&d)
	return &d
}
