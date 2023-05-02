//go:build linux || darwin || netbsd || freebsd || openbsd || dragonfly

package nio

import (
	"goredis/pkg/utils"
	"syscall"
)

func SetReadBuffer(fd, bytes int) error {
	return syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, bytes)
}

func SetWriteBuffer(fd, bytes int) error {
	return syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_SNDBUF, bytes)
}

func SetNoDelay(fd int, noDelay bool) error {
	//开启stream
	return syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, utils.BoolToInt(noDelay))
}

func SetLinger(fd int, onOff, linger int32) error {
	return syscall.SetsockoptLinger(fd, syscall.SOL_SOCKET, syscall.SO_LINGER, &syscall.Linger{
		Onoff:  onOff,  // 1
		Linger: linger, // 0
	})
}
