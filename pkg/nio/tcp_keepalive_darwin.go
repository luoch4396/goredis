package nio

import "syscall"

func SetKeepAlive(fd, sec int, flag bool) error {
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
	//不支持开启TCP_KEEPALIVE
	switch err := syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, 0x101, sec); err {
	case nil, syscall.ENOPROTOOPT:
	default:
		return err
	}
	return syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPALIVE, sec)
}
