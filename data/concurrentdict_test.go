package data

import (
	"goredis/pool/gopool"
	"testing"
)

type test1 struct {
	a string
	b string
	c int32
	d int64
	e interface{}
}

func BenchmarkConcurrentDictPutByPool(b *testing.B) {
	test1 := test1{
		a: "1",
		b: "2",
	}
	dict := NewConcurrentDict(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopool.Go(func() {
			dict.Put("1", test1)
		})
	}
}

func BenchmarkConcurrentDictPut(b *testing.B) {
	test1 := test1{
		a: "1",
		b: "2",
	}
	dict := NewConcurrentDict(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go dict.Put("1", test1)
	}
}

func BenchmarkConcurrentDictGetByPool(b *testing.B) {
	dict := NewConcurrentDict(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gopool.Go(func() {
			dict.Get("1")
			dict.Put("1", "1")
			dict.Get("1")
			dict.Put("2", "1")
			dict.Get("2")
		})
	}
}

func BenchmarkConcurrentDictGet(b *testing.B) {
	dict := NewConcurrentDict(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			dict.Get("1")
			dict.Put("1", "1")
			dict.Get("1")
			dict.Put("2", "1")
			dict.Get("2")
		}()
	}
}
