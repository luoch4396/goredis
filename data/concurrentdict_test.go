package data

import (
	"goredis/pool"
	"testing"
)

type test1 struct {
	a string
	b string
	c int32
	d int64
	e interface{}
}

var _ = pool.GetInstance(100)

func BenchmarkConcurrentDictPutByPool(b *testing.B) {
	b.ResetTimer()
	test1 := test1{
		a: "1",
		b: "2",
	}
	dict := NewConcurrentDict(16)
	for i := 0; i < b.N; i++ {
		_ = pool.Async(func() {
			dict.Put("1", test1)
		})
	}
}

func BenchmarkConcurrentDictPut(b *testing.B) {
	b.ResetTimer()
	test1 := test1{
		a: "1",
		b: "2",
	}
	dict := NewConcurrentDict(16)
	for i := 0; i < b.N; i++ {
		go dict.Put("1", test1)
	}
}

func BenchmarkConcurrentDictGetByPool(b *testing.B) {
	b.ResetTimer()
	dict := NewConcurrentDict(16)
	for i := 0; i < b.N; i++ {
		_ = pool.Async(func() {
			dict.Get("1")
		})
	}
}

func BenchmarkConcurrentDictGet(b *testing.B) {
	b.ResetTimer()
	dict := NewConcurrentDict(16)
	for i := 0; i < b.N; i++ {
		go dict.Get("1")
	}
}
