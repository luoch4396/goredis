package utils

import (
	"fmt"
	"os"
)

func isNotExistMkDir(src string) error {
	if DirIsExist(src) {
		return mkDir(src)
	}
	return nil

}

func mkDir(src string) error {
	return os.MkdirAll(src, os.ModePerm)
}

func CreateIfNotExist(fileName, dir string) (*os.File, error) {
	if CheckPermission(dir) {
		return nil, fmt.Errorf("permission denied dir: %s", dir)
	}

	if err := isNotExistMkDir(dir); err != nil {
		return nil, fmt.Errorf("error during make dir %s, err: %s", dir, err)
	}

	f, err := os.OpenFile(dir+string(os.PathSeparator)+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to open file, err: %s", err)
	}

	return f, nil
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// DirIsExist 判断文件是否存在
func DirIsExist(dir string) bool {
	_, err := os.Stat(dir)
	return os.IsNotExist(err)
}
