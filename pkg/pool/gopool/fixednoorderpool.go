package gopool

import (
	"goredis/pkg/log"
	"runtime"
	"unsafe"
)

// FixedNoOrderPool .
type FixedNoOrderPool struct {
	chTask       chan func()
	panicHandler func(interface{})
}

func (np *FixedNoOrderPool) taskLoop() {
	for f := range np.chTask {
		call(f)
	}
}

func call(f func()) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Errorf("taskpool call failed: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	f()
}

// Go .
func (np *FixedNoOrderPool) Go(f func()) {
	np.chTask <- f
}

// GoByIndex .
func (np *FixedNoOrderPool) GoByIndex(index int, f func()) {
	np.Go(f)
}

// Stop .
func (np *FixedNoOrderPool) Stop() {
	close(np.chTask)
}

func (np *FixedNoOrderPool) SetPanicHandler(f func(interface{})) {
	np.panicHandler = f
}

// NewFixedNoOrderPool .
func NewFixedNoOrderPool(size int, bufferSize int) *FixedNoOrderPool {
	np := &FixedNoOrderPool{
		chTask: make(chan func(), bufferSize),
	}

	for i := 0; i < size; i++ {
		go np.taskLoop()
	}

	return np
}
