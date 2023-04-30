package nio

import (
	"bytes"
	"errors"
	"goredis/pkg/utils/timer"
	"io"
	"net"
	"sync"
	"time"
)

// Conn form net.Conn. all func is supported tcp conn and fit in redis tcp server
type Conn struct {
	p *poll

	hash int

	mux sync.Mutex

	conn net.Conn

	rTimer *timer.Item

	typ      ConnType
	closed   bool
	closeErr error

	ReadBuffer []byte

	// user session.
	session interface{}

	execList []func()

	cache *bytes.Buffer

	DataHandler func(c *Conn, data []byte)
}

// Hash returns a hashcode.
func (c *Conn) Hash() int {
	return c.hash
}

// Read wraps net.Conn.Read.
func (c *Conn) Read(b []byte) (int, error) {
	if c.closeErr != nil {
		return 0, c.closeErr
	}

	var reader io.Reader = c.conn
	if c.cache != nil {
		reader = c.cache
	}
	nread, err := reader.Read(b)
	if c.closeErr == nil {
		c.closeErr = err
	}
	return nread, err
}

func (c *Conn) read(b []byte) (int, error) {
	return c.readTCP(b)
}

func (c *Conn) readTCP(b []byte) (int, error) {
	g := c.p.g
	g.beforeRead(c)
	nread, err := c.conn.Read(b)
	if c.closeErr == nil {
		c.closeErr = err
	}
	if g.onRead != nil {
		if nread > 0 {
			if c.cache == nil {
				c.cache = bytes.NewBuffer(nil)
			}
			c.cache.Write(b[:nread])
		}
		g.onRead(c)
		return nread, nil
	} else if nread > 0 {
		g.onData(c, b[:nread])
	}
	return nread, err
}

// Write from net.Conn.Write.
func (c *Conn) Write(b []byte) (int, error) {
	//TCP
	nwrite, err := c.writeTCP(b)
	return nwrite, err
}

func (c *Conn) writeTCP(b []byte) (int, error) {
	c.p.g.beforeWrite(c)

	nwrite, err := c.conn.Write(b)
	if err != nil {
		if c.closeErr == nil {
			c.closeErr = err
		}
		c.Close()
	}

	return nwrite, err
}

// Writev from buffers.WriteTo/syscall.Writev.
func (c *Conn) Writev(in [][]byte) (int, error) {
	var total = 0
	for _, b := range in {
		nwrite, err := c.Write(b)
		if nwrite > 0 {
			total += nwrite
		}
		if err != nil {
			if c.closeErr == nil {
				c.closeErr = err
			}
			c.Close()
			return total, err
		}
	}
	return total, nil
}

// Close from net.Conn.Close.
func (c *Conn) Close() error {
	var err error
	c.mux.Lock()
	if !c.closed {
		c.closed = true

		if c.rTimer != nil {
			c.rTimer.Stop()
			c.rTimer = nil
		}
		err = c.conn.Close()
		c.mux.Unlock()
		if c.p.g != nil {
			c.p.deleteConn(c)
		}
		return err
	}
	c.mux.Unlock()
	return err
}

func (c *Conn) CloseWithError(err error) error {
	if c.closeErr == nil {
		c.closeErr = err
	}
	return c.Close()
}

// LocalAddr from net.Conn.LocalAddr.
func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr from net.Conn.RemoteAddr.
func (c *Conn) RemoteAddr() net.Addr {
	//tcp
	return c.conn.RemoteAddr()
}

// SetDeadline from net.Conn.SetDeadline.
func (c *Conn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

// SetReadDeadline from net.Conn.SetReadDeadline.
func (c *Conn) SetReadDeadline(t time.Time) error {
	if t.IsZero() {
		t = time.Now().Add(timer.TimeForever)
	}
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline from net.Conn.SetWriteDeadline.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	if t.IsZero() {
		t = time.Now().Add(timer.TimeForever)
	}
	return c.conn.SetWriteDeadline(t)
}

// SetNoDelay from net.Conn.SetNoDelay.
func (c *Conn) SetNoDelay(noDelay bool) error {
	conn, ok := c.conn.(*net.TCPConn)
	if ok {
		return conn.SetNoDelay(noDelay)
	}
	return nil
}

// SetReadBuffer from net.Conn.SetReadBuffer.
func (c *Conn) SetReadBuffer(bytes int) error {
	conn, ok := c.conn.(*net.TCPConn)
	if ok {
		return conn.SetReadBuffer(bytes)
	}
	return nil
}

// SetWriteBuffer from net.Conn.SetWriteBuffer.
func (c *Conn) SetWriteBuffer(bytes int) error {
	conn, ok := c.conn.(*net.TCPConn)
	if ok {
		return conn.SetWriteBuffer(bytes)
	}
	return nil
}

// SetKeepAlive from net.Conn.SetKeepAlive.
func (c *Conn) SetKeepAlive(keepalive bool, d time.Duration) error {
	conn, ok := c.conn.(*net.TCPConn)
	if ok {
		if keepalive && d != 0 {
			err := conn.SetKeepAlive(keepalive)
			err = conn.SetKeepAlivePeriod(d)
			return err
		}
	}
	return nil
}

// SetLinger from net.Conn.SetLinger.
func (c *Conn) SetLinger(onOff int32, linger int32) error {
	conn, ok := c.conn.(*net.TCPConn)
	if ok {
		return conn.SetLinger(int(linger))
	}
	return nil
}

// Session returns user session.
func (c *Conn) Session() interface{} {
	return c.session
}

// SetSession sets user session.
func (c *Conn) SetSession(session interface{}) {
	c.session = session
}

func newConn(conn net.Conn) *Conn {
	c := &Conn{}
	addr := conn.LocalAddr().String()

	c.conn = conn
	c.typ = ConnTypeTCP

	for _, ch := range addr {
		c.hash = 31*c.hash + int(ch)
	}
	if c.hash < 0 {
		c.hash = -c.hash
	}

	return c
}

func NewConn(conn net.Conn) (*Conn, error) {
	if conn == nil {
		return nil, errors.New("invalid conn: nil")
	}
	c, ok := conn.(*Conn)
	if !ok {
		c = newConn(conn)
	}
	return c, nil
}
