package data

import (
	"goredis/pool"
	"testing"
)

type test2 struct {
	a string
	b string
	c int32
	d int64
	e interface{}
}

var err = pool.GetInstance(100)

func BenchmarkSpinDictPutByPool(b *testing.B) {
	test2 := test2{
		a: "1",
		b: "2",
	}
	dict := NewSpinDict(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pool.Async(func() {
			dict.Put("1", test2)
		})
	}
}

func BenchmarkSpinDictPut(b *testing.B) {
	test2 := test2{
		a: "1",
		b: "2",
	}
	dict := NewSpinDict(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go dict.Put("1", test2)
	}
}

func BenchmarkSpinDictGetByPool(b *testing.B) {
	dict := NewSpinDict(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pool.Async(func() {
			dict.Get("1")
		})
	}
}

func BenchmarkSpinDictGet(b *testing.B) {
	dict := NewSpinDict(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go dict.Get("1")
	}
}
