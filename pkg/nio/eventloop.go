package nio

// EventLoop 定义回调事件
type EventLoop interface {
	onOpen(c *Conn)
	onClose(c *Conn) error
	onRead(c *Conn) []byte
	onData(c *Conn)
}
