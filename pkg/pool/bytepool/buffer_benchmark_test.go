package bytepool

import (
	"strings"
	"testing"
)

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
