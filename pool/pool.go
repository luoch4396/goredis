package pool

import (
	"github.com/panjf2000/ants/v2"
	"sync"
)

type AntsPool struct {
	size int
	pool *ants.Pool
}

var (
	instance *AntsPool
	once     sync.Once
)

// GetInstance 获取协程池
func GetInstance(size int) error {
	if size <= 0 {
		size = 1
	}
	var err error
	once.Do(func() {
		if instance == nil {
			pool, err := ants.NewPool(size)
			if err != nil {
				return
			}
			instance = &AntsPool{
				pool: pool,
				size: size,
			}
		}
	})
	return err
}

// Submit 提交任务
func Submit(task func()) error {
	if instance == nil {
		err := GetInstance(1000)
		if err != nil {
			return err
		}
	}
	return instance.pool.Submit(task)
}
