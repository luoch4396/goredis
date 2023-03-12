package main

import (
	"bytes"
	"goredis/config"
	"goredis/pkg/log"
	"goredis/pkg/utils"
	"goredis/redis"
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
	//创建日志建造者
	builder := &log.LoggerBuilder{}
	buf := new(bytes.Buffer)
	fs := &log.FileSettings{
		Path:     config.GlobalProperties.LogFilePath,
		FileName: config.GlobalProperties.LogFileName,
	}
	//初始化日志模块
	builder.
		BuildOutput(buf).
		BuildLevel(config.GlobalProperties.LogLevel).
		BuildFile(fs).
		Build()
	//创建配置文件解析器
	config.NewConfig("redis.properties")
	//开启tcp服务
	redis.NewRedisServer(&redis.Config{
		Address: utils.NewStringBuilder(config.GlobalProperties.Address,
			":", strconv.FormatInt(int64(config.GlobalProperties.Port), 10)),
	})
}
