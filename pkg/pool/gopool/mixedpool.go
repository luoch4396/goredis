package gopool

import (
	"goredis/pkg/log"
	"runtime"
	"sync/atomic"
	"unsafe"
)

type MixedPool struct {
	*FixedNoOrderPool
	parallelism      int32
	totalParallelism int32
	call             func(f func())
}

func (mp *MixedPool) call0(f func()) {
	defer func() {
		if err := recover(); err != nil {
			//自定义panicHandler
			if mp.panicHandler != nil {
				mp.panicHandler(err)
			} else {
				//default panicHandler
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				log.Errorf("gopool call failed: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
			}
		}
		atomic.AddInt32(&mp.parallelism, -1)
	}()
	//执行函数
	f()
}

func (mp *MixedPool) Go(f func()) {
	//并行度判断
	if atomic.AddInt32(&mp.parallelism, 1) <= mp.totalParallelism {
		go func() {
			//执行任务
			mp.call(f)
			//执行任务队列里的任务
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
		//返还任务队列
		mp.FixedNoOrderPool.Go(f)
	}
}

func (mp *MixedPool) Stop() {
	close(mp.chTask)
}

var goPool *MixedPool
var cpus = runtime.NumCPU()

func init() {
	goPool = NewMixedPool(32, 4, 1024)
}

func Go(f func()) {
	goPool.Go(f)
}

func NewMixedPool(totalParallelism int, fixedSize int, bufferSize int) *MixedPool {
	if totalParallelism <= 1 {
		totalParallelism = cpus * 2
	}
	mp := &MixedPool{
		FixedNoOrderPool: NewFixedNoOrderPool(fixedSize, bufferSize),
		totalParallelism: int32(totalParallelism),
	}
	mp.call = mp.call0
	return mp
}
