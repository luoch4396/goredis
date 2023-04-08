//go:build windows

package nio

import "net"

type pollWin struct {
	listener net.Listener
}

func (p *pollWin) accept() error {
	return nil
}
