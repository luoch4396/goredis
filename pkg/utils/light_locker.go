package utils

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type LightLocker struct {
	MaxBackoff int
	Locker     uint32
}

// Lock 自旋
func (ll *LightLocker) Lock() {
	backoff := 1
	for !atomic.CompareAndSwapUint32(&ll.Locker, 0, 1) {
		for i := 0; i < backoff; i++ {
			//让出cpu
			runtime.Gosched()
		}
		if backoff < ll.MaxBackoff {
			backoff += 1
		}
	}
}

func (ll *LightLocker) Unlock() {
	atomic.StoreUint32(&ll.Locker, 0)
}

// TryLock 使用cas实现一个轻量级锁
func (ll *LightLocker) TryLock() bool {
	return atomic.CompareAndSwapUint32(&ll.Locker, 0, 1)
}

func NewLightLock(maxBackOff int) sync.Locker {
	if maxBackOff < 1 {
		maxBackOff = 1
	}
	return &LightLocker{
		MaxBackoff: maxBackOff,
		Locker:     0,
	}
}
