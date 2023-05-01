//go:build linux || darwin || netbsd || freebsd || openbsd || dragonfly

package nio

import (
	"errors"
	"goredis/pkg/pool/bytepool"
	"goredis/pkg/utils/timer"
	"net"
	"sync"
	"syscall"
	"time"
)

// Conn implements net.Conn.
type Conn struct {
	mux sync.Mutex

	p *poll

	fd int

	rTimer *timer.Item
	wTimer *timer.Item

	writeBuffer []byte

	connType ConnType
	closed   bool
	isWAdded bool
	closeErr error

	lAddr net.Addr
	rAddr net.Addr

	ReadBuffer []byte

	session interface{}

	chWaitWrite chan struct{}

	execList []func()

	DataHandler func(c *Conn, data []byte)
}

// Hash returns a hash code.
func (c *Conn) Hash() int {
	return c.fd
}

// Read implements Read.
func (c *Conn) Read(b []byte) (int, error) {
	// use lock to prevent multiple conn data confusion when fd is reused on unix.
	c.mux.Lock()
	if c.closed {
		c.mux.Unlock()
		return 0, net.ErrClosed
	}

	_, n, err := c.doRead(b)
	c.mux.Unlock()
	if err == nil {
		c.p.g.afterRead(c)
	}

	return n, err
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
		c.p.g.afterRead(c)
	}

	return dstConn, n, err
}

func (c *Conn) doRead(b []byte) (*Conn, int, error) {
	return c.readStream(b)
}

func (c *Conn) readStream(b []byte) (*Conn, int, error) {
	nRead, err := syscall.Read(c.fd, b)
	return c, nRead, err
}

// Write implements Write.
func (c *Conn) Write(b []byte) (int, error) {
	c.p.g.beforeWrite(c)

	c.mux.Lock()
	if c.closed {
		c.mux.Unlock()
		return -1, net.ErrClosed
	}

	n, err := c.write(b)
	if err != nil && !errors.Is(err, syscall.EINTR) && !errors.Is(err, syscall.EAGAIN) {
		c.closed = true
		c.mux.Unlock()
		c.closeWithErrorWithoutLock(err)
		return n, err
	}

	if len(c.writeBuffer) == 0 {
		if c.wTimer != nil {
			c.wTimer.Stop()
			c.wTimer = nil
		}
	} else {
		c.modWrite()
	}

	c.mux.Unlock()
	return n, err
}

// Writev implements Writev.
func (c *Conn) Writev(in [][]byte) (int, error) {
	c.p.g.beforeWrite(c)

	c.mux.Lock()
	if c.closed {
		c.mux.Unlock()

		return 0, net.ErrClosed
	}

	var n int
	var err error
	switch len(in) {
	case 1:
		n, err = c.write(in[0])
	default:
		n, err = c.writev(in)
	}
	if err != nil && !errors.Is(err, syscall.EINTR) && !errors.Is(err, syscall.EAGAIN) {
		c.closed = true
		c.mux.Unlock()
		c.closeWithErrorWithoutLock(err)
		return n, err
	}
	if len(c.writeBuffer) == 0 {
		if c.wTimer != nil {
			c.wTimer.Stop()
			c.wTimer = nil
		}
	} else {
		c.modWrite()
	}

	c.mux.Unlock()
	return n, err
}

func (c *Conn) writeStream(b []byte) (int, error) {
	return syscall.Write(c.fd, b)
}

func (c *Conn) Close() error {
	return c.closeWithError(nil)
}

// CloseWithError .
func (c *Conn) CloseWithError(err error) error {
	return c.closeWithError(err)
}

// LocalAddr implements LocalAddr.
func (c *Conn) LocalAddr() net.Addr {
	return c.lAddr
}

// RemoteAddr implements RemoteAddr.
func (c *Conn) RemoteAddr() net.Addr {
	return c.rAddr
}

// SetDeadline implements SetDeadline.
func (c *Conn) SetDeadline(t time.Time) error {
	c.mux.Lock()
	if !c.closed {
		if !t.IsZero() {
			g := c.p.g
			now := time.Now()
			if c.rTimer == nil {
				c.rTimer = g.AfterFunc(t.Sub(now), func() { c.closeWithError(getError("read timeout")) })
			} else {
				c.rTimer.Reset(t.Sub(now))
			}
			if c.wTimer == nil {
				c.wTimer = g.AfterFunc(t.Sub(now), func() { c.closeWithError(getError("write timeout")) })
			} else {
				c.wTimer.Reset(t.Sub(now))
			}
		} else {
			if c.rTimer != nil {
				c.rTimer.Stop()
				c.rTimer = nil
			}
			if c.wTimer != nil {
				c.wTimer.Stop()
				c.wTimer = nil
			}
		}
	}
	c.mux.Unlock()
	return nil
}

func (c *Conn) setDeadline(timer **timer.Item, returnErr error, t time.Time) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.closed {
		return nil
	}
	if !t.IsZero() {
		now := time.Now()
		if *timer == nil {
			*timer = c.p.g.AfterFunc(t.Sub(now), func() { c.closeWithError(returnErr) })
		} else {
			(*timer).Reset(t.Sub(now))
		}
	} else if *timer != nil {
		(*timer).Stop()
		*timer = nil
	}
	return nil
}

// SetReadDeadline implements SetReadDeadline.
func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.setDeadline(&c.rTimer, getError("kqueue read timeout"), t)
}

// SetWriteDeadline implements SetWriteDeadline.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.setDeadline(&c.wTimer, getError("kqueue write timeout"), t)
}

// SetNoDelay implements SetNoDelay.
func (c *Conn) SetNoDelay(noDelay bool) error {
	return SetNoDelay(c.fd, noDelay)
}

// SetReadBuffer implements SetReadBuffer.
func (c *Conn) SetReadBuffer(bytes int) error {
	return SetReadBuffer(c.fd, bytes)
}

// SetWriteBuffer implements SetWriteBuffer.
func (c *Conn) SetWriteBuffer(bytes int) error {
	return SetWriteBuffer(c.fd, bytes)
}

func (c *Conn) SetKeepAlive(keepalive bool, d time.Duration) error {
	if keepalive {
		return SetKeepAlive(c.fd, int(d.Seconds()), keepalive)
	} else {
		return SetKeepAlive(c.fd, 0, keepalive)
	}
}

func (c *Conn) SetLinger(onOff int32, linger int32) error {
	return SetLinger(c.fd, onOff, linger)
}

// Session returns user session.
func (c *Conn) Session() interface{} {
	return c.session
}

// SetSession sets user session.
func (c *Conn) SetSession(session interface{}) {
	c.session = session
}

func (c *Conn) modWrite() {
	if !c.closed && !c.isWAdded {
		c.isWAdded = true
		c.p.modWrite(c.fd)
	}
}

func (c *Conn) resetRead() {
	if !c.closed && c.isWAdded {
		c.isWAdded = false
		p := c.p
		p.deleteEvent(c.fd)
		p.addRead(c.fd)
	}
}

func (c *Conn) write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	if c.overflow(len(b)) {
		return -1, syscall.EINVAL
	}

	if len(c.writeBuffer) == 0 {
		n, err := c.doWrite(b)
		if err != nil && !errors.Is(err, syscall.EINTR) && !errors.Is(err, syscall.EAGAIN) {
			return n, err
		}
		if n < 0 {
			n = 0
		}
		left := len(b) - n
		if left > 0 && c.connType == ConnTypeTCP {
			c.writeBuffer = bytepool.Malloc(left)
			copy(c.writeBuffer, b[n:])
			c.modWrite()
		}
		return len(b), nil
	}
	c.writeBuffer = bytepool.Append(c.writeBuffer, b...)

	return len(b), nil
}

func (c *Conn) flush() error {
	c.mux.Lock()
	if c.closed {
		c.mux.Unlock()
		return net.ErrClosed
	}

	if len(c.writeBuffer) == 0 {
		c.mux.Unlock()
		return nil
	}

	old := c.writeBuffer

	n, err := c.doWrite(old)
	if err != nil && !errors.Is(err, syscall.EINTR) && !errors.Is(err, syscall.EAGAIN) {
		c.closed = true
		c.mux.Unlock()
		c.closeWithErrorWithoutLock(err)
		return err
	}
	if n < 0 {
		n = 0
	}
	left := len(old) - n
	if left > 0 {
		if n > 0 {
			c.writeBuffer = bytepool.Malloc(left)
			copy(c.writeBuffer, old[n:])
			bytepool.Free(old)
		}
	} else {
		bytepool.Free(old)
		c.writeBuffer = nil
		if c.wTimer != nil {
			c.wTimer.Stop()
			c.wTimer = nil
		}
		c.resetRead()
		if c.chWaitWrite != nil {
			select {
			case c.chWaitWrite <- struct{}{}:
			default:
			}
		}
	}

	c.mux.Unlock()
	return nil
}

func (c *Conn) writev(in [][]byte) (int, error) {
	size := 0
	for _, v := range in {
		size += len(v)
	}
	if c.overflow(size) {
		return -1, syscall.EINVAL
	}
	if len(c.writeBuffer) > 0 {
		for _, v := range in {
			c.writeBuffer = bytepool.Append(c.writeBuffer, v...)
		}
		return size, nil
	}

	if len(in) > 1 && size <= 65536 {
		b := bytepool.Malloc(size)
		copied := 0
		for _, v := range in {
			copy(b[copied:], v)
			copied += len(v)
		}
		n, err := c.write(b)
		bytepool.Free(b)
		return n, err
	}

	nwrite := 0
	for _, b := range in {
		n, err := c.write(b)
		if n > 0 {
			nwrite += n
		}
		if err != nil {
			return nwrite, err
		}
	}
	return nwrite, nil
}

func (c *Conn) doWrite(b []byte) (int, error) {
	return c.writeStream(b)
}

func (c *Conn) overflow(n int) bool {
	g := c.p.g
	return g.maxWriteBufferSize > 0 && (len(c.writeBuffer)+n > g.maxWriteBufferSize)
}

func (c *Conn) closeWithError(err error) error {
	c.mux.Lock()
	if !c.closed {
		c.closed = true

		if c.wTimer != nil {
			c.wTimer.Stop()
			c.wTimer = nil
		}
		if c.rTimer != nil {
			c.rTimer.Stop()
			c.rTimer = nil
		}

		c.mux.Unlock()
		return c.closeWithErrorWithoutLock(err)
	}
	c.mux.Unlock()
	return nil
}

func (c *Conn) closeWithErrorWithoutLock(err error) error {
	c.closeErr = err

	if c.writeBuffer != nil {
		bytepool.Free(c.writeBuffer)
		c.writeBuffer = nil
	}

	if c.chWaitWrite != nil {
		select {
		case c.chWaitWrite <- struct{}{}:
		default:
		}
	}

	if c.p.g != nil {
		c.p.deleteConn(c)
	}

	return syscall.Close(c.fd)
}

func NewConn(conn net.Conn) (*Conn, error) {
	if conn == nil {
		return nil, errors.New("invalid conn: nil")
	}
	c, ok := conn.(*Conn)
	if !ok {
		var err error
		c, err = dupStdConn(conn)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}
