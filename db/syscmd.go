package db

import (
	"goredis/interface/redis"
	"goredis/interface/tcp"
	"goredis/pkg/errors"
	"goredis/pkg/monitor"
	"goredis/pkg/utils"
	"goredis/redis/config"
	"goredis/redis/exchange"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// DoPingCmd ping/pong
func DoPingCmd(args [][]byte) tcp.Info {
	if len(args) == 0 {
		return &exchange.PongResponse{}
	} else if len(args) == 1 {
		return exchange.NewStatusInfo(utils.BytesToString(args[0]))
	} else {
		return errors.NewStandardError("wrong number of arguments for 'ping' command")
	}
}

// DoAuthCmd auth
func DoAuthCmd(conn redis.ClientConn, args [][]byte) tcp.Info {
	if len(args) != 1 {
		return errors.NewStandardError("Wrong number of arguments for 'auth' command")
	}
	ok := &exchange.OkResponse{}
	if config.GetPassword() == "" {
		//服务端无密码，不予认证
		return ok
	}
	passwd := utils.BytesToString(args[0])
	conn.SetPassword(passwd)
	if config.GetPassword() != passwd {
		return errors.NewStandardError("invalid password")
	}
	return ok
}

//DoInfoCmd info
func DoInfoCmd(args [][]byte) tcp.Info {
	//多个一起返回
	if len(args) == 1 {
		infoCommandList := [...]string{"server", "client", "cpu", "memory"}
		builder := &strings.Builder{}
		for _, s := range infoCommandList {
			builder.WriteString(GetCustomizeRedisInfo(s))
		}
		return exchange.NewBulkInfo(utils.StringToBytes(builder.String()))
	} else if len(args) == 2 {
		section := strings.ToLower(utils.BytesToString(args[1]))
		switch section {
		//服务器信息
		case "server":
			return exchange.NewBulkInfo(utils.StringToBytes(GetCustomizeRedisInfo("server")))
		//客户端列表
		case "list":
			return exchange.NewBulkInfo(utils.StringToBytes(GetCustomizeRedisInfo("list")))
		//cpu
		case "cpu":
			return exchange.NewBulkInfo(utils.StringToBytes(GetCustomizeRedisInfo("cpu")))
		//内存
		case "memory":
			return exchange.NewBulkInfo(utils.StringToBytes(GetCustomizeRedisInfo("memory")))
		//统计
		case "stats":
			return exchange.NewBulkInfo(utils.StringToBytes(GetCustomizeRedisInfo("stats")))
		default:
			return exchange.NewNullBulkRequest()
		}
	}
	return errors.NewStandardError("ERR wrong number of arguments for 'info' command")
}

// GetCustomizeRedisInfo 返回redis信息
func GetCustomizeRedisInfo(redisType string) string {
	switch redisType {
	case "server":
		return getServerInfo()
	case "cpu":
		return getCpuInfo()
	case "memory":
		return getMemoryInfo()
	}

	return ""
}

//server
func getServerInfo() string {
	startUpTime := time.Since(config.GetStartUpTime()) / time.Second
	builder := &strings.Builder{}
	builder.Grow(256)
	builder.WriteString("redis_version:" + config.Version + "\n")
	builder.WriteString("redis_mode:" + config.GetServerType() + "\n")
	builder.WriteString("os:" + runtime.GOOS + " " + runtime.GOARCH + "\n")
	builder.WriteString("go_version:" + monitor.GetGOVersion() + "\n")
	builder.WriteString("pid:" + monitor.GetPid() + "\n")
	builder.WriteString("tcp_port:" + strconv.Itoa(config.GetPort()) + "\n")
	builder.WriteString("uptime_in_seconds:" + strconv.FormatInt(int64(startUpTime), 10) + "\n")
	builder.WriteString("uptime_in_days:" + strconv.FormatInt(int64(startUpTime/time.Duration(3600*24)), 10) + "\n")
	return builder.String()
}

//cpu
func getCpuInfo() string {
	metricsMap := monitor.GetMetricsInfo(monitor.Goroutines)
	builder := &strings.Builder{}
	builder.Grow(128)
	builder.WriteString("total_cpus:" + strconv.Itoa(runtime.NumCPU()) + "\n")
	builder.WriteString("runtime_threads:" + strconv.Itoa(monitor.GetRuntimeNumThreads()) + "\n")
	builder.WriteString("runtime_goroutines:" + metricsMap[monitor.Goroutines] + "\n")
	builder.WriteString("used_cpu:" + "" + "\n")
	return builder.String()
}

//memory
func getMemoryInfo() string {
	v := monitor.GetGCInfo()
	metricsMap := monitor.GetMetricsInfo(monitor.MemoryClassesTotalBytes, monitor.MemoryClassesHeapObjectsBytes,
		monitor.MemoryClassesHeapUnusedBytes,
	)
	builder := &strings.Builder{}
	builder.Grow(256)
	builder.WriteString("total_memory:" + metricsMap[monitor.MemoryClassesTotalBytes] + "\n")
	builder.WriteString("used_memory:" + "" + "\n")
	builder.WriteString("used_memory_heap:" + metricsMap[monitor.MemoryClassesHeapObjectsBytes] + "\n")
	builder.WriteString("unused_memory_heap:" + metricsMap[monitor.MemoryClassesHeapUnusedBytes] + "\n")
	builder.WriteString("gc_count:" + strconv.FormatInt(v.NumGC, 10) + "\n")
	builder.WriteString("gc_pause_total:" + strconv.FormatInt(int64(v.PauseTotal), 10) + "\n")
	builder.WriteString("memory_clear_strategy:" + "this goredis version is not supported" + "\n")
	return builder.String()
}
