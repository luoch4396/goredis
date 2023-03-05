package config

import (
	"goredis/pkg/log"
	"io"
	"os"
)

var GlobalProperties *GlobalServerProperties

// GlobalServerProperties 定义redis 配置
type GlobalServerProperties struct {
	Bind        string `cfg:"bind"`
	Port        int    `cfg:"port"`
	MaxConn     int    `cfg:"max-conn"`
	Databases   int    `cfg:"databases"`
	ReplTimeout int    `cfg:"repl-timeout"`
}

func NewRedisProperties(globalConfigFileName string) {
	file, err := os.Open(globalConfigFileName)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error("error closing redis config file, program exit", err)
			return
		}
	}(file)
	GlobalProperties = parse(file)
}

func parse(src io.Reader) *GlobalServerProperties {
	return nil
}
