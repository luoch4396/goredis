package monitor

import (
	"strconv"
	"syscall"
)

var Cache = make(map[string]string, 8)

func init() {
	//pid
	Cache["pid"] = strconv.Itoa(syscall.Getpid())
}
