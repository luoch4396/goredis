package db

import "goredis/interface/data"

// DB 定义redis的db
type DB struct {
	index uint16

	data data.Dict

	ttlMap data.Dict

	versionMap data.Dict
}
