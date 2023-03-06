package main

import (
	"goredis/config"
	"goredis/pkg/log"
	"goredis/pkg/utils"
	"goredis/server"
	"strconv"
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
	config.NewConfig("redis.properties")
	//开启tcp服务
	server.NewRedisServer(&server.Config{
		Address: utils.NewStringBuilder(config.GlobalProperties.Address,
			":", strconv.FormatInt(config.GlobalProperties.Port, 10)),
	})
}
