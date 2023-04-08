//go:build windows

package nio

type Conn struct {
	p *poll
}

func (c *Conn) Write(b []byte) (int, error) {
	return 0, nil
}
