package nio

import (
	"goredis/pkg/log"
	"goredis/pkg/utils/timer"
	"net"
	"runtime"
	"strings"
)

// Start init and start pollers.
func (g *Engine) Start() error {
	switch g.network {
	case "tcp", "tcp4", "tcp6":
		for i := range g.addrs {
			ln, err := newPoller(g, true, i)
			if err != nil {
				for j := 0; j < i; j++ {
					g.listeners[j].stop()
				}
				return err
			}
			g.addrs[i] = ln.listener.Addr().String()
			g.listeners = append(g.listeners, ln)
		}
	}

	for i := 0; i < g.pollerNum; i++ {
		p, err := newPoller(g, false, i)
		if err != nil {
			for j := 0; j < len(g.listeners); j++ {
				g.listeners[j].stop()
			}

			for j := 0; j < i; j++ {
				g.polls[j].stop()
			}
			return err
		}
		g.polls[i] = p
	}

	for i := 0; i < g.pollerNum; i++ {
		g.Add(1)
		go g.polls[i].start()
	}
	for _, l := range g.listeners {
		g.Add(1)
		go l.start()
	}

	g.Timer.Start()

	if len(g.addrs) == 0 {
		log.Infof("NIO-SERVER-ENGINE %v start", g.Name)
	} else {
		log.Infof("NIO-SERVER-ENGINE[%v] start listen on: [\"%v@%v\"]", g.Name, g.network, strings.Join(g.addrs, `", "`))
	}
	return nil
}

// NewEngine is a factory impl.
func NewEngine(conf Config) *Engine {
	cpuNum := runtime.NumCPU()
	if conf.Name == "" {
		conf.Name = "NIO-IOCP"
	}
	if conf.NPoller <= 0 {
		conf.NPoller = cpuNum
	}
	if conf.ReadBufferSize <= 0 {
		conf.ReadBufferSize = DefaultReadBufferSize
	}
	if conf.Listen == nil {
		conf.Listen = net.Listen
	}

	g := &Engine{
		Timer:              timer.New(conf.Name, conf.TimerExecute),
		Name:               conf.Name,
		network:            conf.Network,
		addrs:              conf.Addrs,
		listen:             conf.Listen,
		pollerNum:          conf.NPoller,
		readBufferSize:     conf.ReadBufferSize,
		maxWriteBufferSize: conf.MaxWriteBufferSize,
		lockListener:       conf.LockListener,
		lockPoller:         conf.LockPoller,
		listeners:          make([]*poll, len(conf.Addrs))[0:0],
		polls:              make([]*poll, conf.NPoller),
		connsStd:           map[*Conn]struct{}{},
	}

	g.initHandlers()

	g.OnReadBufferAlloc(func(c *Conn) []byte {
		if c.ReadBuffer == nil {
			c.ReadBuffer = make([]byte, g.readBufferSize)
		}
		return c.ReadBuffer
	})

	return g
}
