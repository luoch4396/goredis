package goroutinespool

import (
	"goredis/pkg/utils"
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	cpus = runtime.NumCPU()
)

type Pool struct {
	//协程池的总大小
	size uint32
	//跑了多少个协程
	running int32
	//池持有的锁
	lock sync.Locker
	//工作队列
	workers workQueue
	//协程池的状态
	state int32
	//池持有的锁的条件
	cond *sync.Cond
	//工作func缓存
	workerCache sync.Pool
	//协程执行的等待时间
	waiting int32
	//配置
	config *Config
}

type Config struct {
	//是否预分配内存
	PreAllocated bool
	//协程池的最大数量
	size uint16
}

func newGoPool(config *Config) {
	size := config.size
	if size <= 0 {
		size = uint16(cpus * 16)
	}

	pool := &Pool{
		size:   uint32(size),
		lock:   utils.NewLightLock(5),
		config: config,
	}
	pool.workerCache.New = func() interface{} {
		return &goWorker{
			pool: pool,
			task: make(chan func(), workerChanCap),
		}
	}
	if pool.config.PreAllocated {
		pool.workers = newWorkerArray(loopQueueType, size)
	} else {
		pool.workers = newWorkerArray(stackType, 0)
	}

	pool.cond = sync.NewCond(pool.lock)
	return pool, nil
}

//增加一个任务数量
func (p *Pool) addRunning(num int) {
	atomic.AddInt32(&p.running, int32(num))
}
