package redis

import "goredis/interface/tcp"

type DB interface {
	// Exec 执行命令
	Exec(client ClientConn, cmd [][]byte) tcp.Info
	// Close 关闭
	Close()
}
