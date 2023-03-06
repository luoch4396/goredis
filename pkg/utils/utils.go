package utils

import (
	"os"
	"strings"
)

// FileIsExist 判断文件是否存在
func FileIsExist(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

// NewStringBuilder 字符串拼接
func NewStringBuilder(strs ...string) string {
	var stringBuilder strings.Builder
	//stringBuilder.Grow(n * len(str))
	for i := 0; i < len(strs); i++ {
		stringBuilder.WriteString(strs[i])
	}
	return stringBuilder.String()
}
