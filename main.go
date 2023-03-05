package main

import "goredis/server"

var banner = `
                  ___     
   ________  ____/ (_)____
  / ___/ _ \/ __  / / ___/
 / /  /  __/ /_/ / (__  ) 
/_/   \___/\__,_/_/____/
`

func main() {
	//打印banner
	print(banner)
	//开启netty tcp服务器
	server.NewNettyServer()

}
