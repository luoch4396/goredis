package utils

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type LightLocker struct {
	maxBackoff uint16
}

type lightLocker uint32

func (ll *lightLocker) Lock() {
	backoff := 1
	for !atomic.CompareAndSwapUint32((*uint32)(ll), 0, 1) {
		for i := 0; i < backoff; i++ {
			//让出cpu
			runtime.Gosched()
		}
		if backoff < maxBackoff {
			backoff = backoff - 1
		}
	}
}

func (ll *lightLocker) Unlock() {
	atomic.StoreUint32((*uint32)(ll), 0)
}

func NewLightLock(locker LightLocker) sync.Locker {
	return new(lightLocker)
}
