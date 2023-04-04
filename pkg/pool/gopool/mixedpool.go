package gopool

import (
	"goredis/pkg/log"
	"runtime"
	"sync/atomic"
	"unsafe"
)

// MixedPool .
type MixedPool struct {
	*FixedNoOrderPool
	cuncurrent int32
	nativeSize int32
	call       func(f func())
}

func (mp *MixedPool) callWithRecover(f func()) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Errorf("taskpool call failed: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
		atomic.AddInt32(&mp.cuncurrent, -1)
	}()
	f()
}

func (mp *MixedPool) callWitoutRecover(f func()) {
	defer atomic.AddInt32(&mp.cuncurrent, -1)
	f()
}

// Go .
func (mp *MixedPool) Go(f func()) {
	if atomic.AddInt32(&mp.cuncurrent, 1) <= mp.nativeSize {
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
		atomic.AddInt32(&mp.cuncurrent, -1)
		mp.FixedNoOrderPool.Go(f)
	}
}

// Stop .
func (mp *MixedPool) Stop() {
	close(mp.chTask)
}

// NewMixedPool .
func NewMixedPool(nativeSize int, fixedSize int, bufferSize int, v ...interface{}) *MixedPool {
	mp := &MixedPool{
		FixedNoOrderPool: NewFixedNoOrderPool(fixedSize, bufferSize),
		nativeSize:       int32(nativeSize),
	}
	mp.call = mp.callWithRecover
	if len(v) > 0 {
		if withoutRecover, ok := v[0].(bool); ok && withoutRecover {
			mp.call = mp.callWitoutRecover
		}
	}
	return mp
}
