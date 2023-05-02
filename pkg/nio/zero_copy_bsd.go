//go:build darwin || dragonfly || freebsd || netbsd || openbsd

package nio

import "syscall"

//no support zero copy
func setZeroCopy(fd int) error {
	return syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, SO_ZEROCOPY, 1)
}
