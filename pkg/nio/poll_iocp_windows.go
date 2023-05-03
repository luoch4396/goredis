package nio

import (
	"goredis/pkg/log"
	"net"
	"runtime"
	"time"
)

type poll struct {
	g *Engine

	index int

	ReadBuffer []byte

	pollType   string
	isListener bool
	listener   net.Listener
	shutdown   bool

	chStop chan struct{}
}

func (p *poll) accept() error {
	conn, err := p.listener.Accept()
	if err != nil {
		return err
	}

	c := newConn(conn)
	o := p.g.polls[c.Hash()%len(p.g.polls)]
	o.addConn(c)

	return nil
}

func (p *poll) readConn(c *Conn) {
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

func (p *poll) addConn(c *Conn) error {
	c.p = p
	p.g.mux.Lock()
	p.g.connsStd[c] = struct{}{}
	p.g.mux.Unlock()
	p.g.onOpen(c)
	go p.readConn(c)
	return nil
}

func (p *poll) deleteConn(c *Conn) {
	p.g.mux.Lock()
	delete(p.g.connsStd, c)
	p.g.mux.Unlock()
	//关闭tcp连接
	p.g.onClose(c, c.closeErr)
}

func (p *poll) start() {
	if p.g.lockListener {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
	}
	defer p.g.Done()

	log.Debugf("iocp-poll[%v][%v_%v] start", p.g.Name, p.pollType, p.index)
	defer log.Debugf("iocp-poll[%v][%v_%v] stopped", p.g.Name, p.pollType, p.index)

	if p.isListener {
		var err error
		p.shutdown = false
		for !p.shutdown {
			err = p.accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					log.Errorf("iocp-poll[%v][%v_%v] Accept failed: temporary error, retrying...", p.g.Name, p.pollType, p.index)
					time.Sleep(time.Second / 20)
				} else {
					if !p.shutdown {
						log.Errorf("iocp-poll[%v][%v_%v] Accept failed: %v, exit...", p.g.Name, p.pollType, p.index, err)
					}
					break
				}
			}

		}
	}
	<-p.chStop
}

func (p *poll) stop() {
	log.Debugf("iocp-poll[%v][%v_%v] stop...", p.g.Name, p.pollType, p.index)
	p.shutdown = true
	if p.isListener {
		p.listener.Close()
	}
	close(p.chStop)
}

func newPoller(g *Engine, isListener bool, index int) (*poll, error) {
	p := &poll{
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
		p.pollType = "POLL-LISTENER"
	} else {
		p.pollType = "POLL-IOCP"
	}

	return p, nil
}
