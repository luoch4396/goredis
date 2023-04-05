package db

import "goredis/interface/tcp"

func Set(db *DB, args [][]byte) *tcp.Response {
	return &tcp.Response{}
}

func SetEX(db *DB, args [][]byte) *tcp.Response {
	return &tcp.Response{}
}

func SetNX(db *DB, args [][]byte) *tcp.Response {
	return &tcp.Response{}
}
