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

// BytesToString non memory copy api convert []byte to string
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes non memory copy api convert string to []byte
func StringToBytes(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
