package redis

import "goredis/interface/tcp"

type Server interface {
	// Exec 执行命令
	Exec(client ClientConn, cmd [][]byte) (rep tcp.Info)
	// Close 关闭
	Close()
}
