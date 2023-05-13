package nio

import (
	"context"
	"goredis/pkg/log"
	"goredis/pkg/utils/timer"
	"net"
	"runtime"
	"sync"
	"time"
	"unsafe"
)

const (
	// DefaultReadBufferSize .
	DefaultReadBufferSize = 1024 * 32

	// DefaultMaxWriteBufferSize .
	DefaultMaxWriteBufferSize = 1024 * 1024

	// DefaultMaxConnReadTimesPerEventLoop .
	DefaultMaxConnReadTimesPerEventLoop = 3

	// DefaultUDPReadTimeout .
	DefaultUDPReadTimeout = 120 * time.Second
)

var (
	// MaxOpenFiles .
	MaxOpenFiles = 1024 * 1024 * 2
)

// Engine is a manager of poller.
type Engine struct {
	*timer.Timer
	sync.WaitGroup

	Name string

	Execute      func(f func())
	TimerExecute func(f func())

	mux sync.Mutex

	wgConn sync.WaitGroup

	network   string
	addrs     []string
	listen    func(network, addr string) (net.Listener, error)
	listenUDP func(network string, laddr *net.UDPAddr) (*net.UDPConn, error)

	pollerNum                    int
	readBufferSize               int
	maxWriteBufferSize           int
	maxConnReadTimesPerEventLoop int
	udpReadTimeout               time.Duration
	epollMod                     uint32
	lockListener                 bool
	lockPoller                   bool

	connsStd  map[*Conn]struct{}
	connsUnix []*Conn

	listeners []*poller
	pollers   []*poller

	onOpen            func(c *Conn)
	onClose           func(c *Conn, err error)
	onRead            func(c *Conn)
	onData            func(c *Conn, data []byte)
	onReadBufferAlloc func(c *Conn) []byte
	onReadBufferFree  func(c *Conn, buffer []byte)
	// onWriteBufferFree func(c *Conn, buffer []byte)
	beforeRead  func(c *Conn)
	afterRead   func(c *Conn)
	beforeWrite func(c *Conn)
	onStop      func()
}

// Stop closes listeners/pollers/conns/timer.
func (g *Engine) Stop() {
	for _, l := range g.listeners {
		l.stop()
	}

	g.mux.Lock()
	conns := g.connsStd
	g.connsStd = map[*Conn]struct{}{}
	connsUnix := g.connsUnix
	g.mux.Unlock()

	g.wgConn.Done()
	for c := range conns {
		if c != nil {
			cc := c
			g.Async(func() {
				cc.Close()
			})
		}
	}
	for _, c := range connsUnix {
		if c != nil {
			cc := c
			g.Async(func() {
				cc.Close()
			})
		}
	}

	g.wgConn.Wait()
	time.Sleep(time.Second / 5)

	g.onStop()

	g.Timer.Stop()

	for i := 0; i < g.pollerNum; i++ {
		g.pollers[i].stop()
	}

	g.Wait()
	log.Infof("Nio_Server[%v] stop", g.Name)
}

// Shutdown stops Engine gracefully with context.
func (g *Engine) Shutdown(ctx context.Context) error {
	ch := make(chan struct{})
	go func() {
		g.Stop()
		close(ch)
	}()

	select {
	case <-ch:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// AddConn adds conn to a poller.
func (g *Engine) AddConn(conn net.Conn) (*Conn, error) {
	c, err := NBConn(conn)
	if err != nil {
		return nil, err
	}

	p := g.pollers[c.Hash()%len(g.pollers)]
	p.addConn(c)
	return c, nil
}

// OnOpen registers callback for new connection.
func (g *Engine) OnOpen(h func(c *Conn)) {
	if h == nil {
		panic("invalid nil handler")
	}
	g.onOpen = func(c *Conn) {
		g.wgConn.Add(1)
		h(c)
	}
}

// OnClose registers callback for disconnected.
func (g *Engine) OnClose(h func(c *Conn, err error)) {
	if h == nil {
		panic("invalid nil handler")
	}
	g.onClose = func(c *Conn, err error) {
		defer g.wgConn.Done()
		h(c, err)
	}
}

// OnRead registers callback for reading event.
//func (g *Engine) OnRead(h func(c *Conn)) {
//	g.onRead = h
//}

// OnData registers callback for data.
func (g *Engine) OnData(h func(c *Conn, data []byte)) {
	if h == nil {
		panic("invalid nil handler")
	}
	g.onData = h
}

// OnReadBufferAlloc registers callback for memory allocating.
func (g *Engine) OnReadBufferAlloc(h func(c *Conn) []byte) {
	if h == nil {
		panic("invalid nil handler")
	}
	g.onReadBufferAlloc = h
}

// OnReadBufferFree registers callback for memory release.
func (g *Engine) OnReadBufferFree(h func(c *Conn, b []byte)) {
	if h == nil {
		panic("invalid nil handler")
	}
	g.onReadBufferFree = h
}

// OnWriteBufferRelease registers callback for write buffer memory release.
// func (g *Engine) OnWriteBufferRelease(h func(c *Conn, b []byte)) {
// 	if h == nil {
// 		panic("invalid nil handler")
// 	}
// 	g.onWriteBufferFree = h
// }

// BeforeRead registers callback before syscall.Read
// the handler would be called on windows.
func (g *Engine) BeforeRead(h func(c *Conn)) {
	if h == nil {
		panic("invalid nil handler")
	}
	g.beforeRead = h
}

// AfterRead registers callback after syscall.Read
// the handler would be called on *nix.
func (g *Engine) AfterRead(h func(c *Conn)) {
	if h == nil {
		panic("invalid nil handler")
	}
	g.afterRead = h
}

// BeforeWrite registers callback befor syscall.Write and syscall.Writev
// the handler would be called on windows.
func (g *Engine) BeforeWrite(h func(c *Conn)) {
	if h == nil {
		panic("invalid nil handler")
	}
	g.beforeWrite = h
}

// OnStop registers callback before Engine is stopped.
func (g *Engine) OnStop(h func()) {
	if h == nil {
		panic("invalid nil handler")
	}
	g.onStop = h
}

// PollerBuffer returns Poller's buffer by Conn, can be used on linux/bsd.
func (g *Engine) PollerBuffer(c *Conn) []byte {
	return c.p.ReadBuffer
}

func (g *Engine) initHandlers() {
	g.wgConn.Add(1)
	g.OnOpen(func(c *Conn) {})
	g.OnClose(func(c *Conn, err error) {})
	// g.OnRead(func(c *Conn, b []byte) ([]byte, error) {
	// 	n, err := c.Read(b)
	// 	if n > 0 {
	// 		return b[:n], err
	// 	}
	// 	return nil, err
	// })
	g.OnData(func(c *Conn, data []byte) {})
	g.OnReadBufferAlloc(g.PollerBuffer)
	g.OnReadBufferFree(func(c *Conn, buffer []byte) {})
	// g.OnWriteBufferRelease(func(c *Conn, buffer []byte) {})
	g.BeforeRead(func(c *Conn) {})
	g.AfterRead(func(c *Conn) {})
	g.BeforeWrite(func(c *Conn) {})
	g.OnStop(func() {})

	if g.Execute == nil {
		g.Execute = func(f func()) {
			defer func() {
				if err := recover(); err != nil {
					const size = 64 << 10
					buf := make([]byte, size)
					buf = buf[:runtime.Stack(buf, false)]
					log.Errorf("execute failed: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
				}
			}()
			f()
		}
	}
}

func (g *Engine) borrow(c *Conn) []byte {
	return g.onReadBufferAlloc(c)
}

func (g *Engine) payback(c *Conn, buffer []byte) {
	g.onReadBufferFree(c, buffer)
}
