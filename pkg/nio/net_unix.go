package nio

import (
	"errors"
	"net"
	"syscall"
)

func init() {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err == nil {
		if n := int(limit.Max); n > 0 && n < MaxOpenFiles {
			MaxOpenFiles = n
		}
	}
}

func dupStdConn(conn net.Conn) (*Conn, error) {
	sc, ok := conn.(interface {
		SyscallConn() (syscall.RawConn, error)
	})
	if !ok {
		return nil, errors.New("RawConn Unsupported")
	}
	rc, err := sc.SyscallConn()
	if err != nil {
		return nil, errors.New("RawConn Unsupported")
	}

	var newFd int
	errCtrl := rc.Control(func(fd uintptr) {
		newFd, err = syscall.Dup(int(fd))
	})

	if errCtrl != nil {
		return nil, errCtrl
	}

	if err != nil {
		return nil, err
	}

	lAddr := conn.LocalAddr()
	rAddr := conn.RemoteAddr()

	conn.Close()

	c := &Conn{
		fd:    newFd,
		lAddr: lAddr,
		rAddr: rAddr,
	}

	switch conn.(type) {
	case *net.TCPConn:
		c.typ = ConnTypeTCP
	case *net.UnixConn:
		c.typ = ConnTypeUnix
	default:
	}

	return c, nil
}
