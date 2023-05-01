//go:build linux || darwin || netbsd || freebsd || openbsd || dragonfly

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
	case "unix", "tcp", "tcp4", "tcp6":
		for i := range g.addrs {
			ln, err := newPoll(g, true, i)
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
		p, err := newPoll(g, false, i)
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
		g.polls[i].ReadBuffer = make([]byte, g.readBufferSize)
		g.Add(1)
		go g.polls[i].start()
	}

	for _, l := range g.listeners {
		g.Add(1)
		go l.start()
	}

	g.Timer.Start()

	if len(g.addrs) == 0 {
		log.Infof("NIO-SERVER[%v] start", g.Name)
	} else {
		log.Infof("NIO-SERVER[%v] start listen on: [\"%v@%v\"]", g.Name, g.network, strings.Join(g.addrs, `", "`))
	}
	return nil
}

// NewEngine is a factory impl.
func NewEngine(conf Config) *Engine {
	cpuNum := runtime.NumCPU()
	if conf.Name == "" {
		conf.Name = "NIO"
	}
	if conf.NPoller <= 0 {
		conf.NPoller = cpuNum
	}
	if conf.ReadBufferSize <= 0 {
		conf.ReadBufferSize = DefaultReadBufferSize
	}
	if conf.MaxConnReadTimesPerEventLoop <= 0 {
		conf.MaxConnReadTimesPerEventLoop = DefaultMaxConnReadTimesPerEventLoop
	}
	if conf.Listen == nil {
		conf.Listen = net.Listen
	}

	g := &Engine{
		Timer:                        timer.New(conf.Name, conf.TimerExecute),
		Name:                         conf.Name,
		network:                      conf.Network,
		addrs:                        conf.Addrs,
		listen:                       conf.Listen,
		pollerNum:                    conf.NPoller,
		readBufferSize:               conf.ReadBufferSize,
		maxWriteBufferSize:           conf.MaxWriteBufferSize,
		maxConnReadTimesPerEventLoop: conf.MaxConnReadTimesPerEventLoop,
		epollMod:                     conf.EpollMod,
		lockListener:                 conf.LockListener,
		lockPoller:                   conf.LockPoller,
		listeners:                    make([]*poll, len(conf.Addrs))[0:0],
		polls:                        make([]*poll, conf.NPoller),
		connsUnix:                    make([]*Conn, MaxOpenFiles),
	}

	g.initHandlers()

	return g
}
