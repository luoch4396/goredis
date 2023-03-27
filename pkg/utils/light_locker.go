package utils

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type lightLocker struct {
	maxBackoff int
	locker     uint32
}

// Lock 自旋
func (ll *lightLocker) Lock() {
	backoff := 1
	for !atomic.CompareAndSwapUint32(&ll.locker, 0, 1) {
		for i := 0; i < backoff; i++ {
			//让出cpu
			runtime.Gosched()
		}
		if backoff < ll.maxBackoff {
			backoff <<= 1
		}
	}
}

func (ll *lightLocker) Unlock() {
	atomic.StoreUint32(&ll.locker, 0)
}

// TryLock 使用cas实现一个轻量级锁
func (ll *lightLocker) TryLock() bool {
	return atomic.CompareAndSwapUint32(&ll.locker, 0, 1)
}

func NewLightLock(maxBackOff int) sync.Locker {
	if maxBackOff < 1 {
		maxBackOff = 1
	}
	return &lightLocker{
		maxBackoff: maxBackOff,
		locker:     0,
	}
}
