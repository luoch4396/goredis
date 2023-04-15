package utils

import (
	"strings"
)

// NewStringBuilder 字符串拼接
func NewStringBuilder[T []byte | string](t ...T) string {
	if t == nil || len(t) == 0 {
		return ""
	}
	builder := &strings.Builder{}
	//stringBuilder.Grow(n * len(str))
	for i := 0; i < len(t); i++ {
		builder.WriteString(string(t[i]))
	}
	return builder.String()
}
