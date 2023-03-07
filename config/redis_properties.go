package config

import (
	"bufio"
	"goredis/pkg/log"
	"goredis/pkg/utils"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	GlobalProperties  *GlobalServerProperties
	defaultProperties = &GlobalServerProperties{
		Address:   "127.0.0.1",
		Port:      6379,
		MaxConns:  100,
		Databases: 5,
	}
)

// GlobalServerProperties 定义redis 配置
type GlobalServerProperties struct {
	Address     string `config:"address"`
	Port        int    `config:"port"`
	MaxConns    int    `config:"max-conns"`
	Databases   int    `config:"databases"`
	CopyTimeout int    `config:"copy-timeout"`
}

func NewConfig(globalConfigFileName string) {
	//配置文件不存在，获取默认配置
	if !utils.FileIsExist(globalConfigFileName) {
		GlobalProperties = defaultProperties
		return
	}
	file, err := os.Open(globalConfigFileName)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error("exception closing redis config file, program exit", err)
			return
		}
	}(file)
	GlobalProperties = parse(file)
}

//解析配置文件
func parse(src io.Reader) *GlobalServerProperties {
	config := &GlobalServerProperties{}
	// read config file
	rawMap := make(map[string]string)
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && strings.TrimLeft(line, " ")[0] == '#' {
			continue
		}
		pivot := strings.IndexAny(line, "=")
		if pivot > 0 && pivot < len(line)-1 { // separator found
			key := line[0:pivot]
			value := strings.Trim(line[pivot+1:], " ")
			rawMap[strings.ToLower(key)] = value
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(``, err)
	}

	//使用反射去解析具体配置参数
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	n := t.Elem().NumField()
	for i := 0; i < n; i++ {
		field := t.Elem().Field(i)
		fieldVal := v.Elem().Field(i)
		key, ok := field.Tag.Lookup("config")
		if !ok || strings.TrimLeft(key, " ") == "" {
			key = field.Name
		}
		value, ok := rawMap[strings.ToLower(key)]
		if ok {
			// fill config
			switch field.Type.Kind() {
			case reflect.String:
				fieldVal.SetString(value)
			case reflect.Int:
				intValue, err := strconv.ParseInt(value, 10, 64)
				if err == nil {
					fieldVal.SetInt(intValue)
				}
			case reflect.Bool:
				boolValue := "yes" == value
				fieldVal.SetBool(boolValue)
			}
		}
	}
	//log.Info(fmt.Sprint(config.Port))
	//todo: 开启服务的时候打印一下配置信息
	//fmt.Print(config.Port)
	return config
}
