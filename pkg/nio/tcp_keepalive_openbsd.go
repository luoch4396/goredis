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
	//只设置开启
	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 1); err != nil {
		return err
	}
	return nil
}
