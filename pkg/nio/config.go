package nio

// Config TCP配置
type Config struct {
	MaxOpenFiles int
	KeepAlive    bool
}

var (
	MaxOpenFiles = 1024 * 1024 * 2
)
