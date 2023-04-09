//go:build netbsd || freebsd || dragonfly || linux

package nio

import "syscall"

func SetKeepAlive(flag bool, fd, sec int) error {
	//禁用TCP_KEEPALIVE
	if !flag {
		if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 0); err != nil {
			return err
		}
		return nil
	}
	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 1); err != nil {
		return err
	}
	if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPINTVL, sec); err != nil {
		return err
	}
	return syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPIDLE, sec)
}
