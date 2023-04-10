//go:build darwin || netbsd || freebsd || openbsd || dragonfly

package nio

import (
	"syscall"
)

type kqueuePoll struct {
	trigger int
	fd      int
}

func (p *kqueuePoll) accept() error {
	return nil
}

// Trigger implements Poll.
func (p *kqueuePoll) Trigger() error {
	_, err := syscall.Kevent(p.fd, []syscall.Kevent_t{{
		Ident:  0,
		Filter: syscall.EVFILT_USER,
		Fflags: syscall.NOTE_TRIGGER,
	}}, nil, nil)
	return err
}
