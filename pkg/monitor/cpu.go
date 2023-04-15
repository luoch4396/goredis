package monitor

import "runtime"

// GetRuntimeNumThreads go-runtime的线程数
func GetRuntimeNumThreads() int {
	n, _ := runtime.ThreadCreateProfile(nil)
	return n
}
