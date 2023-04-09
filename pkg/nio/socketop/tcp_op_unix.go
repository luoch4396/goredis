package socketop

import "syscall"

func SetReadBuffer(fd, bytes int) error {
	//接收
	return syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, bytes)
}

func SetWriteBuffer(fd, bytes int) error {
	//发送
	return syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_SNDBUF, bytes)
}

func SetNoDelay(fd int, noDelay bool) error {
	//开启stream
	if noDelay {
		return syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1)
	}
	return syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 0)
}

// SetLinger implements SetLinger.
func etLinger(fd int, onOff, linger int32) error {
	return syscall.SetsockoptLinger(fd, syscall.SOL_SOCKET, syscall.SO_LINGER, &syscall.Linger{
		Onoff:  onOff,  // 1
		Linger: linger, // 0
	})
}
