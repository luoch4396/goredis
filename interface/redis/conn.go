package redis

// ClientConn redis 客户端连接
type ClientConn interface {
	// Write 向redis db写数据
	Write([]byte) (int, error)
	// Close 执行关闭操作，归还给pool
	Close() error
	// Name 名称
	Name() string
}
