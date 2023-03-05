package utils

import "os"

// FileIsExist 判断文件是否存在
func FileIsExist(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}
