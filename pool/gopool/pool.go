package gopool

import (
	"context"
	"sync"
	"sync/atomic"
)

type Pool interface {
	// Name returns the corresponding pool name.
	Name() string
	// SetCap sets the goroutine capacity of the pool.
	SetCap(cap int32)
	// Go executes f.
	Go(f func())
	// CtxGo executes f and accepts the context.
	CtxGo(ctx context.Context, f func())
	// SetPanicHandler sets the panic handler.
	SetPanicHandler(f func(context.Context, interface{}))
	// WorkerCount returns the number of running workers
	WorkerCount() int32
}

var taskPool sync.Pool

func init() {
	taskPool.New = newTask
}

type task struct {
	ctx context.Context
	f   func()

	next *task
}

func (t *task) zero() {
	t.ctx = nil
	t.f = nil
	t.next = nil
}

func (t *task) Recycle() {
	t.zero()
	taskPool.Put(t)
}

func newTask() interface{} {
	return &task{}
}

type taskList struct {
	sync.Mutex
	taskHead *task
	taskTail *task
}

type pool struct {
	// The name of the pool
	name string

	// capacity of the pool, the maximum number of goroutines that are actually working
	cap int32
	// Configuration information
	config *Config
	// 所有需要执行的任务形成一个链表结构
	taskHead  *task
	taskTail  *task
	taskLock  sync.Mutex
	taskCount int32

	// 正在跑的任务数量
	workerCount int32

	// 任务的错误处理器
	panicHandler func(context.Context, interface{})
}

// NewPool creates a new pool with the given name, cap and config.
func NewPool(name string, cap int32, config *Config) Pool {
	p := &pool{
		name:   name,
		cap:    cap,
		config: config,
	}
	return p
}

func (p *pool) Name() string {
	return p.name
}

func (p *pool) SetCap(cap int32) {
	atomic.StoreInt32(&p.cap, cap)
}

func (p *pool) Go(f func()) {
	p.CtxGo(context.Background(), f)
}

func (p *pool) CtxGo(ctx context.Context, f func()) {
	t := taskPool.Get().(*task)
	t.ctx = ctx
	t.f = f
	p.taskLock.Lock()
	if p.taskHead == nil {
		p.taskHead = t
		p.taskTail = t
	} else {
		p.taskTail.next = t
		p.taskTail = t
	}
	p.taskLock.Unlock()
	atomic.AddInt32(&p.taskCount, 1)
	// The following two conditions are met:
	// 1. the number of tasks is greater than the threshold.
	// 2. The current number of workers is less than the upper limit p.cap.
	// or there are currently no workers.
	if (atomic.LoadInt32(&p.taskCount) >= p.config.ScaleThreshold && p.WorkerCount() < atomic.LoadInt32(&p.cap)) || p.WorkerCount() == 0 {
		p.incWorkerCount()
		w := workerPool.Get().(*worker)
		w.pool = p
		w.run()
	}
}

// SetPanicHandler the func here will be called after the panic has been recovered.
func (p *pool) SetPanicHandler(f func(context.Context, interface{})) {
	p.panicHandler = f
}

func (p *pool) WorkerCount() int32 {
	return atomic.LoadInt32(&p.workerCount)
}

func (p *pool) incWorkerCount() {
	atomic.AddInt32(&p.workerCount, 1)
}

func (p *pool) decWorkerCount() {
	atomic.AddInt32(&p.workerCount, -1)
}
