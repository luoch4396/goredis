package gopool

import (
	"context"
)

var (
	//这里只实现一个pool，用于控制内部协程上限，目前是所有的，没有做pool缓存
	defaultPool Pool
)

func init() {
	defaultPool = NewPool("pool.DefaultPool", 1000, NewConfig())
}

// Go is an alternative to the go keyword, which is able to recover panic.
//
//	pool.Go(func(arg interface{}){
//	    ...
//	}(nil))
func Go(f func()) {
	CtxGo(context.Background(), f)
}

// CtxGo context.Context 使用协程上下文
// CtxGo is preferred than Go.
func CtxGo(ctx context.Context, f func()) {
	defaultPool.CtxGo(ctx, f)
}

// SetCap is not recommended to be called, this func changes the global base's capacity which will affect other callers.
func SetCap(cap int32) {
	defaultPool.SetCap(cap)
}

// SetPanicHandler sets the panic handler for the global base.
func SetPanicHandler(f func(context.Context, interface{})) {
	defaultPool.SetPanicHandler(f)
}

// WorkerCount returns the number of global default base's running workers
func WorkerCount() int32 {
	return defaultPool.WorkerCount()
}
