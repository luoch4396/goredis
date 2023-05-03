package nio

import "syscall"

// no support zero copy?
func setZeroCopy(fd int) error {
	return syscall.EINVAL
}
