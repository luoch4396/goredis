package data

import (
	"goredis/pkg/pool/gopool"
	"sync"
	"testing"
)

var BenchTimes = 10000

func BenchmarkConcurrentDictGetByPool(b *testing.B) {
	dict := NewConcurrentDict(16)
	p := gopool.NewMixedPool(32, 4, 1024)
	b.ReportAllocs()
	b.ResetTimer()
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(BenchTimes)
		for j := 0; j < BenchTimes; j++ {
			p.Go(func() {
				dict.Get("1")
				dict.Put("1", "1")
				dict.Get("1")
				dict.Put("2", "1")
				dict.Get("2")
				wg.Done()
			})
		}
		wg.Wait()
	}
}
