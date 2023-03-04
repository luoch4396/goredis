package main

import "goredis/server"

func main() {

	//开启netty tcp服务器
	server.NewNettyServer()

}
