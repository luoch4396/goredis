package utils

import (
	"strings"
)

// NewStringBuilder 字符串拼接
func NewStringBuilder[T []byte | string](strs ...T) string {
	if strs == nil || len(strs) == 0 {
		return ""
	}
	var stringBuilder strings.Builder
	//stringBuilder.Grow(n * len(str))
	for i := 0; i < len(strs); i++ {
		stringBuilder.Write([]byte(strs[i]))
	}
	return stringBuilder.String()
}
