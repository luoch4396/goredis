//go:build darwin || netbsd || freebsd || openbsd || dragonfly

package nio

import (
	"goredis/pkg/log"
	"net"
	"os"
	"sync"
	"syscall"
	"time"
)

type poll struct {
	fd           int
	eventFd      int
	shutdown     bool
	listener     net.Listener
	isListener   bool
	unixSockAddr string
	ReadBuffer   []byte
	pollType     string
	eventList    []syscall.Kevent_t
	mux          sync.Mutex
	index        int
}

func (p *poll) addConn(c *Conn) {
	fd := c.fd
	//if fd >= len(p.g.connsUnix) {
	//c.closeWithError(fmt.Errorf("too many open files, fd[%d] >= MaxOpenFiles[%d]", fd, len(p.g.connsUnix)))
	//return
	//}
	c.p = p
	//p.g.connsUnix[fd] = c
	p.modRead(fd)
}

func (p *poll) getConn(fd int) *Conn {
	//return p.g.connsUnix[fd]
	return nil
}

//删除连接
func (p *poll) deleteConn(c *Conn) {
	if c == nil {
		return
	}
}

func (p *poll) accept() error {
	return nil
}

func (p *poll) trigger() error {
	_, err := syscall.Kevent(p.fd, []syscall.Kevent_t{{
		Ident:  0,
		Filter: syscall.EVFILT_USER,
		Fflags: syscall.NOTE_TRIGGER,
	}}, nil, nil)
	return err
}

func (p *poll) start() {
	if p.isListener {
		p.acceptLoop()
	} else {
		defer syscall.Close(p.eventFd)
		p.readWriteLoop()
	}
}

func (p *poll) acceptLoop() {
	p.shutdown = false
	for !p.shutdown {
		conn, err := p.listener.Accept()
		if err == nil {
			_, err := NewConn(conn)
			if err != nil {
				conn.Close()
				continue
			}
			//p.g.pollers[c.Hash()%len(p.g.pollers)].addConn(c)
		} else {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				log.Errorf("kqueuePoll[%v][%v_%v] Accept failed: temporary error, will be retrying ...", p.pollType, p.index)
				time.Sleep(time.Second / 10)
			} else {
				if !p.shutdown {
					log.Errorf("kqueuePoll[%v][%v_%v] Accept failed: %v, will be exited ...", p.pollType, p.index, err)
				}
				break
			}
		}
	}
}

func (p *poll) readWriteLoop() {
	var events = make([]syscall.Kevent_t, 1024)
	var changes []syscall.Kevent_t

	p.shutdown = false
	for !p.shutdown {
		p.mux.Lock()
		changes = p.eventList
		p.eventList = nil
		p.mux.Unlock()
		n, err := syscall.Kevent(p.eventFd, changes, events, nil)
		if err != nil && err != syscall.EINTR {
			return
		}

		for i := 0; i < n; i++ {
			switch int(events[i].Ident) {
			case p.eventFd:
			default:
				p.readWrite(&events[i])
			}
		}
	}
}

func (p *poll) readWrite(ev *syscall.Kevent_t) {
	if ev.Flags&syscall.EV_DELETE > 0 {
		return
	}
	fd := int(ev.Ident)
	c := p.getConn(fd)
	if c != nil {
		if ev.Filter&syscall.EVFILT_READ == syscall.EVFILT_READ {
			//if p.g.onRead == nil {
			//	for {
			//		rc, n, err := c.ReadAndGetConn(buffer)
			//		if n > 0 {
			//			p.g.onData(rc, buffer[:n])
			//		}
			//		p.g.payback(c, buffer)
			//		if err == syscall.EINTR {
			//			continue
			//		}
			//		if err == syscall.EAGAIN {
			//			return
			//		}
			//		if (err != nil || n == 0) && ev.Flags&syscall.EV_DELETE == 0 {
			//			c.closeWithError(err)
			//		}
			//		if n < len(buffer) {
			//			break
			//		}
			//	}
			//} else {
			//p.g.onRead(c)
			//}
		}

		if ev.Filter&syscall.EVFILT_WRITE == syscall.EVFILT_WRITE {
			//c.flush()
		}
	} else {
		syscall.Close(fd)
		// p.deleteEvent(fd)
	}
}

func (p *poll) deleteEvent(fd int) {
	p.mux.Lock()
	p.eventList = append(p.eventList, syscall.Kevent_t{Ident: uint64(fd), Flags: syscall.EV_DELETE, Filter: syscall.EVFILT_READ})
	p.mux.Unlock()
	p.trigger()
}

func (p *poll) stop() {
	log.Debugf("kqueuePoll[%v][%v_%v] stop...", p.pollType, p.index)
	p.shutdown = true
	if p.listener != nil {
		p.listener.Close()
		if p.unixSockAddr != "" {
			os.Remove(p.unixSockAddr)
		}
	}
	p.trigger()
}

func (p *poll) modRead(fd int) {
	p.mux.Lock()
	p.eventList = append(p.eventList, syscall.Kevent_t{Ident: uint64(fd), Flags: syscall.EV_ADD, Filter: syscall.EVFILT_READ})
	p.mux.Unlock()
	p.trigger()
}

func (p *poll) modWrite(fd int) {
	p.mux.Lock()
	p.eventList = append(p.eventList, syscall.Kevent_t{Ident: uint64(fd), Flags: syscall.EV_ADD, Filter: syscall.EVFILT_WRITE})
	p.mux.Unlock()
	p.trigger()
}
