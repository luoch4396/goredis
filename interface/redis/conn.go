package redis

// ClientConn redis 客户端连接
type ClientConn interface {
	// Write 向客户端发送返回数据
	Write([]byte) bool
	// Close 执行关闭操作，归还给pool
	Close() error
	// Name 名称
	Name() string
}
