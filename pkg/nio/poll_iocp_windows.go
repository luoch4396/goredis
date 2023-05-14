package nio

import (
	"goredis/pkg/log"
	"net"
	"runtime"
	"time"
)

const (
	// EPOLLLT .
	EPOLLLT = 0

	// EPOLLET .
	EPOLLET = 1
)

type poller struct {
	g *Engine

	index int

	ReadBuffer []byte

	pollType   string
	isListener bool
	listener   net.Listener
	shutdown   bool

	chStop chan struct{}
}

func (p *poller) accept() error {
	conn, err := p.listener.Accept()
	if err != nil {
		return err
	}

	c := conn
	o := p.g.pollers[c.Hash()%len(p.g.pollers)]
	o.addConn(c)

	return nil
}

func (p *poller) readConn(c *Conn) {
	for {
		buffer := p.g.borrow(c)
		_, err := c.read(buffer)
		p.g.payback(c, buffer)
		if err != nil {
			c.Close()
			return
		}
	}
}

func (p *poller) addConn(c *Conn, virtualUDPConn ...interface{}) error {
	c.p = p
	p.g.mux.Lock()
	p.g.connsStd[c] = struct{}{}
	p.g.mux.Unlock()
	// should not call onOpen for udp server conn
	if c.typ != ConnTypeUDPServer {
		p.g.onOpen(c)
	}
	// should not read udp client from reading udp server conn
	if c.typ != ConnTypeUDPClientFromRead {
		go p.readConn(c)
	}

	return nil
}

func (p *poller) deleteConn(c *Conn) {
	p.g.mux.Lock()
	delete(p.g.connsStd, c)
	p.g.mux.Unlock()
	// should not call onClose for udp server conn
	if c.typ != ConnTypeUDPServer {
		p.g.onClose(c, c.closeErr)
	}
}

func (p *poller) start() {
	if p.g.lockListener {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
	}
	defer p.g.Done()

	log.Debugf("Nio_Iocp_Poller[%v][%v_%v] start", p.g.Name, p.pollType, p.index)

	if p.isListener {
		var err error
		p.shutdown = false
		for !p.shutdown {
			err = p.accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					log.Errorf("Nio_Iocp_Poller[%v][%v_%v] Accept failed: temporary error, retrying...", p.g.Name, p.pollType, p.index)
					time.Sleep(time.Second / 20)
				} else {
					if !p.shutdown {
						log.Errorf("Nio_Iocp_Poller[%v][%v_%v] Accept failed: %v, exit...", p.g.Name, p.pollType, p.index, err)
					}
					break
				}
			}

		}
	}
	<-p.chStop
}

func (p *poller) stop() {
	log.Debugf("Nio_Iocp_Poller[%v][%v_%v] stop...", p.g.Name, p.pollType, p.index)
	p.shutdown = true
	if p.isListener {
		p.listener.Close()
	}
	close(p.chStop)
}

func newPoller(g *Engine, isListener bool, index int) (*poller, error) {
	p := &poller{
		g:          g,
		index:      index,
		isListener: isListener,
		chStop:     make(chan struct{}),
	}

	if isListener {
		var err error
		var addr = g.addrs[index%len(g.addrs)]
		p.listener, err = g.listen(g.network, addr)
		if err != nil {
			return nil, err
		}
		p.pollType = "Poller-Listener"
	} else {
		p.pollType = "Poller"
	}

	return p, nil
}
