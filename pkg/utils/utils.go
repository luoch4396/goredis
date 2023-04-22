package utils

import (
	"reflect"
	"strings"
	"unsafe"
)

// NewStringBuilder 字符串拼接
func NewStringBuilder(t ...string) string {
	if t == nil || len(t) == 0 {
		return ""
	}
	builder := &strings.Builder{}
	//stringBuilder.Grow(n * len(str))
	for i := 0; i < len(t); i++ {
		builder.WriteString(t[i])
	}
	return builder.String()
}

// BytesToString no memory allocation api
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToBytes(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}
