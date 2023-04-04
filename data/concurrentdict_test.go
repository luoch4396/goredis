package data

import (
	"sync"
	"testing"
)

type test1 struct {
	a string
	b string
	c int32
	d int64
	e interface{}
}

var BenchTimes = 10000

func BenchmarkConcurrentDictGet(b *testing.B) {
	dict := NewConcurrentDict(16)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(BenchTimes)
		for j := 0; j < BenchTimes; j++ {
			go func() {
				dict.Get("1")
				dict.Put("1", "1")
				dict.Get("1")
				dict.Put("2", "1")
				dict.Get("2")
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
