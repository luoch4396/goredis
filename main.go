package main

import (
	"fmt"
	"goredis/config"
	"goredis/pkg/log"
	"goredis/server"
)

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
	//初始化日志模块
	log.NewLog4j()
	//创建配置文件解析器
	config.NewRedisProperties("redis.properties")
	//开启tcp服务
	server.NewNettyServer(&server.Config{
		Address: fmt.Sprintf("%s:%d", config.GlobalProperties.Address, config.GlobalProperties.Port),
	})

}
