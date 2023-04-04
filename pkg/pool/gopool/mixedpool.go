package gopool

import (
	"goredis/pkg/log"
	"runtime"
	"sync/atomic"
	"unsafe"
)

var cpus = runtime.NumCPU()

// MixedPool .
type MixedPool struct {
	*FixedNoOrderPool
	parallelism      int32
	totalParallelism int32
	call             func(f func())
	panicHandler     func(interface{})
}

func (mp *MixedPool) callWithRecover(f func()) {
	defer func() {
		if err := recover(); err != nil {
			if mp.panicHandler != nil {
				mp.panicHandler(err)
			} else {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				log.Errorf("taskpool call failed: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
			}
		}
		atomic.AddInt32(&mp.parallelism, -1)
	}()
	//执行函数
	f()
}

// Go .
func (mp *MixedPool) Go(f func()) {
	if atomic.AddInt32(&mp.parallelism, 1) <= mp.totalParallelism {
		go func() {
			mp.call(f)
			for len(mp.chTask) > 0 {
				select {
				case f = <-mp.chTask:
					mp.call(f)
				default:
					return
				}
			}
		}()
	} else {
		atomic.AddInt32(&mp.parallelism, -1)
		mp.FixedNoOrderPool.Go(f)
	}
}

// Stop .
func (mp *MixedPool) Stop() {
	close(mp.chTask)
}

// NewMixedPool .
func NewMixedPool(totalParallelism int, fixedSize int, bufferSize int) *MixedPool {
	if totalParallelism <= 1 {
		totalParallelism = cpus
	}
	mp := &MixedPool{
		FixedNoOrderPool: NewFixedNoOrderPool(fixedSize, bufferSize),
		totalParallelism: int32(totalParallelism),
	}
	mp.call = mp.callWithRecover
	return mp
}
