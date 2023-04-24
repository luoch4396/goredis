package nio

import (
	"net"
)

// Config TCP配置
type Config struct {
	Name string

	// Network is the listening protocol, used with Addrs together.
	// tcp* supported only by now, there's no plan for other protocol such as udp,
	// because it's too easy to write udp server/client.
	Network string

	// Addrs is the listening addr list for a nbio server.
	// if it is empty, no listener created, then the Engine is used for client by default.
	Addrs []string

	// NPoller represents poller goroutine num, it's set to runtime.NumCPU() by default.
	NPoller int

	// ReadBufferSize represents buffer size for reading, it's set to 16k by default.
	ReadBufferSize int

	// MaxWriteBufferSize represents max write buffer size for Conn, it's set to 1m by default.
	// if the connection's Send-Q is full and the data cached by nbio is
	// more than MaxWriteBufferSize, the connection would be closed by nbio.
	MaxWriteBufferSize int

	// MaxConnReadTimesPerEventLoop represents max read times in one poller loop for one fd
	MaxConnReadTimesPerEventLoop int

	// LockListener represents listener's goroutine to lock thread or not, it's set to false by default.
	LockListener bool

	// LockPoller represents poller's goroutine to lock thread or not, it's set to false by default.
	LockPoller bool

	// EpollMod sets the epoll mod, EPOLLLT by default.
	EpollMod uint32

	// TimerExecute sets the executor for timer callbacks.
	TimerExecute func(f func())

	// Listen is used to create listener for Engine.
	Listen func(network, addr string) (net.Listener, error)

	MaxOpenFiles int64
}
