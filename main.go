package main

import (
	"goredis/config"
	"goredis/pkg/log"
	"goredis/pkg/utils"
	"goredis/redis"
	"os"
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
	//创建配置文件解析器
	config.NewConfig("redis.yaml")
	configs := config.Configs
	//日志
	fs := &log.FileSettings{
		Path:     configs.Log.FilePath,
		FileName: configs.Log.FileName,
	}
	//初始化日志模块
	log.NewLoggerBuilder().
		BuildStdOut(os.Stdout).
		BuildLevel(configs.Log.LogLevel).
		BuildFile(fs).
		Build()
	//开启tcp服务
	redis.NewRedisServer(&redis.Config{
		Address: utils.NewStringBuilder(configs.Server.Address,
			":", strconv.FormatInt(int64(configs.Server.Port), 10)),
	})
}
