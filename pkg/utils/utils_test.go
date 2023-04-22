package utils

import "testing"

func BenchmarkStringToBytes(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	s := "11111"
	for i := 0; i < b.N; i++ {
		a([]byte(s))
	}
}

func BenchmarkDefaultStringToBytes(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	s := "11111"
	for i := 0; i < b.N; i++ {
		a(StringToBytes(s))
	}
}

func a(b []byte) {
	//test
}
