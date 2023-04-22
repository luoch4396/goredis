package bytepool

import (
	"strings"
	"testing"
)

func BenchmarkRingBufferAsyncRead(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	rb := NewRingBuffer(32 * 1024)
	data := []byte(strings.Repeat("bbbbbaaaaa", 512))
	buf := make([]byte, 512)

	go func() {
		for {
			rb.Read(buf)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Write(data)
	}
}

func BenchmarkRingBufferAsyncWrite(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	rb := NewRingBuffer(32 * 1024)
	data := []byte(strings.Repeat("bbbbbaaaaa", 512))
	_ = make([]byte, 512)

	go func() {
		for {
			rb.Write(data)
		}
	}()

	b.ResetTimer()
	//for i := 0; i < b.N; i++ {
	//	rb.Read(buf)
	//}
}

func BenchmarkBytePoolAsyncRead(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	bs := Malloc(32 * 1024)
	data := strings.Repeat("bbbbbaaaaa", 512)

	go func() {
		for {
			AppendString(bs, data)
		}
	}()

	b.ResetTimer()
	//for i := 0; i < b.N; i++ {
	//	bs.Write(data)
	//}
}

func BenchmarkBytePoolAsyncWrite(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	bs := Malloc(32 * 1024)
	data := strings.Repeat("bbbbbaaaaa", 512)

	go func() {
		for {
			AppendString(bs, data)
		}
	}()

	b.ResetTimer()
	//for i := 0; i < b.N; i++ {
	//	bs.Read(data)
	//}
}
