package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var (
	SingleServer = "single"
	//ClusterMode  = "cluster"
)

// Config 定义redis 配置
type Config struct {
	Log struct {
		FileName string `yaml:"file-name"`
		FilePath string `yaml:"file-path"`
		LogLevel string `yaml:"log-level"`
	}
	Server struct {
		Address    string `yaml:"address"`
		Port       int    `yaml:"port"`
		MaxConn    int    `yaml:"max-conn"`
		Databases  int    `yaml:"databases"`
		Password   string `yaml:"password"`
		ServerType string `yaml:"Server-type"`
	}
	Pools []Pool `yaml:"pools"`
}

type Pool struct {
	size int `yaml:"size"`
}

var Configs Config

func NewConfig(globalConfigFileName string) {
	//配置文件不存在，获取默认配置
	configFileName := os.Getenv("CONFIG")
	if configFileName != "" {
		globalConfigFileName = configFileName
	}
	file, err := os.ReadFile(globalConfigFileName)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &Configs)
	if err != nil {
		panic(err)
	}
}

func GetDatabases() int {
	return Configs.Server.Databases
}

func SetDatabases(databases int) {
	Configs.Server.Databases = databases
}

func GetMaxConn() int {
	return Configs.Server.MaxConn
}

func GetServerType() string {
	return Configs.Server.ServerType
}

func SetServerType(serverType string) {
	Configs.Server.ServerType = serverType
}

func GetPassword() string {
	return Configs.Server.Password
}
