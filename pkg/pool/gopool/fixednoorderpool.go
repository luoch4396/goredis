package gopool

import (
	"goredis/pkg/log"
)

type FixedNoOrderPool struct {
	chTask       chan func()
	panicHandler func(interface{})
}

func (np *FixedNoOrderPool) taskLoop() {
	for f := range np.chTask {
		call(f, np.panicHandler)
	}
}

func call(f func(), panicHandler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			if panicHandler != nil {
				panicHandler(err)
			} else {
				log.MakeErrorLog(err)
			}

		}
	}()
	//执行
	f()
}

func (np *FixedNoOrderPool) Go(f func()) {
	np.chTask <- f
}

func (np *FixedNoOrderPool) Stop() {
	close(np.chTask)
}

func (np *FixedNoOrderPool) SetPanicHandler(f func(interface{})) {
	np.panicHandler = f
}

func NewFixedNoOrderPool(size int, bufferSize int) *FixedNoOrderPool {
	np := &FixedNoOrderPool{
		chTask: make(chan func(), bufferSize),
	}

	for i := 0; i < size; i++ {
		go np.taskLoop()
	}

	return np
}
