package goroutinespool

import (
	"goredis/pkg/log"
	"runtime"
	"time"
)

type goWorker struct {
	//协程池
	pool *Pool
	//chan
	task chan func()
	//重新放回队列的时间
	recycleTime time.Time
}

func (gw *goWorker) run() {
	pool := gw.pool

	pool.addRunning(1)
	go func() {
		defer func() {
			pool.addRunning(-1)
			pool.workerCache.Put(gw)
			if p := recover(); p != nil {
				var buf [4096]byte
				//创建栈
				n := runtime.Stack(buf[:], false)
				log.Infof("worker run with any error: %s", string(buf[:n]))
			}
			pool.cond.Signal()
		}()

		for f := range gw.task {
			if f == nil {
				return
			}
			f()
			if ok := pool.revertWorker(gw); !ok {
				return
			}
		}
	}()
}
