package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

// Config 定义redis 配置
type Config struct {
	Log struct {
		FileName string `yaml:"file-name"`
		FilePath string `yaml:"file-path"`
		LogLevel string `yaml:"log-level"`
	}
	Server struct {
		Address   string `yaml:"address"`
		Port      int    `yaml:"port"`
		MaxConn   int    `yaml:"max-conn"`
		Databases int    `yaml:"databases"`
	}
}

var Configs Config

func NewConfig(globalConfigFileName string) {
	//配置文件不存在，获取默认配置
	configFileName := os.Getenv("CONFIG")
	if configFileName != "" {
		globalConfigFileName = configFileName
	}
	file, err := ioutil.ReadFile(globalConfigFileName)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &Configs)
}
