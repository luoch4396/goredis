//go:build windows

package nio

import "net"

type poll struct {
	listener net.Listener
}

func (p *poll) accept() error {
	return nil
}
