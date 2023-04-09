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
	//开启TCP_KEEPALIVE
	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 1); err != nil {
		return err
	}
	//低版本 osx 不支持开启TCP_KEEPALIVE 返回错误
	switch err := syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, 0x101, sec); err {
	case nil, syscall.ENOPROTOOPT:
	default:
		return err
	}
	return syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPALIVE, sec)
}
