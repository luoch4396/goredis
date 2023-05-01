package nio

import (
	"goredis/pkg/errors"
	"goredis/pkg/log"
	"net"
	"time"
)

// ConnType .
type ConnType = int8

const (
	// ConnTypeTCP tcp socket
	ConnTypeTCP ConnType = iota + 1
	// ConnTypeUnix unix socket
	ConnTypeUnix
)

// 封装错误
func getError(errStr string) error {
	return errors.NewStandardError(errStr)
}

// Type .
func (c *Conn) Type() ConnType {
	return c.connType
}

// IsTCP .
func (c *Conn) IsTCP() bool {
	return c.connType == ConnTypeTCP
}

// IsUnix FOR UNIX CONN
func (c *Conn) IsUnix() bool {
	return c.connType == ConnTypeUnix
}

// OnData registers callback for data.
func (c *Conn) OnData(h func(conn *Conn, data []byte)) {
	c.DataHandler = h
}

// Lock .
func (c *Conn) Lock() {
	c.mux.Lock()
}

// Unlock .
func (c *Conn) Unlock() {
	c.mux.Unlock()
}

// IsClosed .
func (c *Conn) IsClosed() (bool, error) {
	return c.closed, c.closeErr
}

// ExecuteLen .
func (c *Conn) ExecuteLen() int {
	c.mux.Lock()
	n := len(c.execList)
	c.mux.Unlock()
	return n
}

// Execute .
func (c *Conn) Execute(f func()) bool {
	c.mux.Lock()
	if c.closed {
		c.mux.Unlock()
		return false
	}

	isHead := len(c.execList) == 0
	c.execList = append(c.execList, f)
	c.mux.Unlock()

	if isHead {
		c.p.g.Execute(func() {
			i := 0
			for {
				func() {
					defer func() {
						if err := recover(); err != nil {
							log.MakeErrorLog(err)
						}
					}()
					f()
				}()

				c.mux.Lock()
				i++
				if len(c.execList) == i {
					c.execList = c.execList[0:0]
					c.mux.Unlock()
					return
				}
				f = c.execList[i]
				c.mux.Unlock()
			}
		})
	}

	return true
}

// MustExecute .
func (c *Conn) MustExecute(f func()) {
	c.mux.Lock()
	isHead := len(c.execList) == 0
	c.execList = append(c.execList, f)
	c.mux.Unlock()

	if isHead {
		c.p.g.Execute(func() {
			i := 0
			for {
				f()

				c.mux.Lock()
				i++
				if len(c.execList) == i {
					c.execList = c.execList[0:0]
					c.mux.Unlock()
					return
				}
				f = c.execList[i]
				c.mux.Unlock()
			}
		})
	}
}

// Dial from net.Dial.
func Dial(network string, address string) (*Conn, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewConn(conn)
}

// DialTimeout wraps net.DialTimeout.
func DialTimeout(network string, address string, timeout time.Duration) (*Conn, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}
	return NewConn(conn)
}
