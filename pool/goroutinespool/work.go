package goroutinespool

import (
	"goredis/pkg/log"
	"runtime"
	"time"
)

type goWorker struct {
	// pool who owns this worker.
	pool *Pool

	// task is a job should be done.
	task chan func()

	// recycleTime will be updated when putting a worker back into queue.
	recycleTime time.Time
}

func (gw *goWorker) run() {
	gw.pool.addRunning(1)
	go func() {
		defer func() {
			gw.pool.addRunning(-1)
			gw.pool.workerCache.Put(gw)
			if p := recover(); p != nil {
				if ph := gw.pool.options.PanicHandler; ph != nil {
					ph(p)
				} else {
					log.Errorf("worker exits from a panic: %v", p)
					var buf [4096]byte
					n := runtime.Stack(buf[:], false)
					log.Errorf("worker exits from panic: %s", string(buf[:n]))
				}
			}
			// Call Signal() here in case there are goroutines waiting for available workers.
			gw.pool.cond.Signal()
		}()

		for f := range gw.task {
			if f == nil {
				return
			}
			f()
			if ok := gw.pool.revertWorker(gw); !ok {
				return
			}
		}
	}()
}
