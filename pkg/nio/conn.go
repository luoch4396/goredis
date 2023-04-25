package nio

import (
	"goredis/pkg/errors"
	"goredis/pkg/log"
)

// ConnType .
type ConnType = int8

const (
	// ConnTypeTCP tcp socket
	ConnTypeTCP ConnType = iota + 1
	xx
	xxx
	xxxx
	// ConnTypeUnix unix socket
	ConnTypeUnix
)

// 封装错误
func getError(errStr string) error {
	return errors.NewStandardError(errStr)
}

// Type .
func (c *Conn) Type() ConnType {
	return c.typ
}

// IsTCP .
func (c *Conn) IsTCP() bool {
	return c.typ == ConnTypeTCP
}

// IsUnix FOR UNIX CONN
func (c *Conn) IsUnix() bool {
	return c.typ == ConnTypeUnix
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
