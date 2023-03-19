package pool

import (
	"github.com/panjf2000/ants/v2"
	"goredis/pkg/utils"
)

type AntsPool struct {
	size int
	Pool *ants.Pool
}

var (
	instance *AntsPool
	lock     = utils.NewLightLock(16)
)

func GetInstance(size int) (*AntsPool, error) {
	if size <= 0 {
		size = 1
	}
	if instance != nil {
		return instance, nil
	} else {
		lock.Lock()
		defer lock.Unlock()
		//双重判断
		if instance != nil {
			return instance, nil
		} else {
			pool, err := ants.NewPool(size)
			if err != nil {
				return nil, err
			}
			instance = &AntsPool{
				Pool: pool,
				size: size,
			}
			return instance, nil
		}
	}
}
