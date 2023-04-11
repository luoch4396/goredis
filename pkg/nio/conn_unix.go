package nio

import (
	"errors"
	"net"
	"sync"
	"syscall"
)

type Conn struct {
	mux sync.Mutex
	//多路复用器
	p *poll
	//描述符
	fd int
	//连接类型 tcp、unix、udp
	connType ConnType
	closed   bool
	closeErr error
}

func NewConn(conn net.Conn) (*Conn, error) {
	return nil, nil
}

func (c *Conn) Write(b []byte) (int, error) {
	return 0, nil
}

func (c *Conn) readStream(b []byte) (*Conn, int, error) {
	n, err := syscall.Read(c.fd, b)
	return c, n, err
}

func (c *Conn) writeStream(b []byte) (int, error) {
	return syscall.Write(c.fd, b)
}

// ReadAndGetConn .
func (c *Conn) ReadAndGetConn(b []byte) (*Conn, int, error) {
	// use lock to prevent multiple conn data confusion when fd is reused on unix.
	c.mux.Lock()
	if c.closed {
		c.mux.Unlock()
		return c, 0, net.ErrClosed
	}

	dstConn, n, err := c.doRead(b)
	c.mux.Unlock()
	if err == nil {
		//c.p.g.afterRead(c)
	}

	return dstConn, n, err
}

func (c *Conn) doRead(b []byte) (*Conn, int, error) {
	switch c.connType {
	case ConnTypeTCP, ConnTypeUnix:
		return c.readStream(b)
	default:
	}
	return c, 0, errors.New("invalid udp conn for reading")
}
