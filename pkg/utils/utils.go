package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-netty/go-netty"
	"io"
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

func NewStringBuilder0(bytes ...[]byte) string {
	var stringBuilder strings.Builder
	for i := 0; i < len(bytes); i++ {
		stringBuilder.Write(bytes[i])
	}
	return stringBuilder.String()
}

func ParseMessage(message netty.Message) ([]byte, error) {
	switch r := message.(type) {
	case io.Reader:
		var reader = bufio.NewReader(r)
		line, err := reader.ReadBytes('\n')
		return line, err
	case []byte:
		return r, nil
	case [][]byte:
		buffer := bytes.NewBuffer(nil)
		for _, b := range r {
			buffer.Write(b)
		}
		return buffer.Bytes(), nil
	case string:
		return []byte(r), nil
	default:
		return nil, fmt.Errorf("unrecognized type: %T", message)
	}
}
