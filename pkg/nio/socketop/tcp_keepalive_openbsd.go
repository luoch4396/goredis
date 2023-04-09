package socketop

import "syscall"

func SetKeepAlive(flag bool, fd, sec int) error {
	//禁用TCP_KEEPALIVE
	if !flag {
		if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 0); err != nil {
			return err
		}
	}
	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 1); err != nil {
		return err
	}
	//不知道这个系统为什么有KEEPALIVE，但是没有时间设置？
	return nil
}
