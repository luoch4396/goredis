package utils

import (
	"strings"
)

// NewStringBuilder 字符串拼接
func NewStringBuilder(strs ...string) string {
	var stringBuilder strings.Builder
	//stringBuilder.Grow(n * len(str))
	for i := 0; i < len(strs); i++ {
		stringBuilder.WriteString(strs[i])
	}
	return stringBuilder.String()
}

func NewStringBuilder0(bytes ...[]byte) string {
	var stringBuilder strings.Builder
	for i := 0; i < len(bytes); i++ {
		stringBuilder.Write(bytes[i])
	}
	return stringBuilder.String()
}
