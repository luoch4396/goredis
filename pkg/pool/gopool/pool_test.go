package gopool

import (
	"sync"
	"testing"
	"time"
)

const testLoopNum = 1024
const sleepTime = time.Nanosecond * 10

func BenchmarkGo(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(testLoopNum)
		for j := 0; j < testLoopNum; j++ {
			go call(func() {
				time.Sleep(sleepTime)
				wg.Done()
			}, nil)
		}
		wg.Wait()
	}
}

func BenchmarkFixedNoOrderPool(b *testing.B) {
	p := NewFixedNoOrderPool(32, 1024)
	defer p.Stop()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(testLoopNum)
		for j := 0; j < testLoopNum; j++ {
			p.Go(func() {
				time.Sleep(sleepTime)
				wg.Done()
			})
		}
		wg.Wait()
	}
}

func BenchmarkMixedPool(b *testing.B) {
	p := NewMixedPool(32, 4, 1024)
	defer p.Stop()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(testLoopNum)
		for j := 0; j < testLoopNum; j++ {
			p.Go(func() {
				time.Sleep(sleepTime)
				wg.Done()
			})
		}
		wg.Wait()
	}
}
