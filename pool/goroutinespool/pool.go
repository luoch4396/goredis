package goroutinespool

import (
	"sync"
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
	size uint32
}

func newGoPool(config *Config) {

}
