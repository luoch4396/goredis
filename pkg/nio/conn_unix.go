package nio

import "syscall"

type Conn struct {
	//多路复用器
	p *poll
	//描述符
	fd int
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
