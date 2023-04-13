package db

import (
	"fmt"
	"goredis/interface/redis"
	"goredis/interface/tcp"
	"goredis/pkg/errors"
	"goredis/redis/config"
	"goredis/redis/exchange"
	"os"
	"runtime"
	"strings"
	"time"
)

// CmdStrategy 命令解析策略接口
type CmdStrategy interface {
	Do(conn redis.ClientConn, args [][]byte) tcp.Info
}

type CmdOperator struct {
	CmdStrategy CmdStrategy
}

// NewCmdOperator 创建一个命令策略
func NewCmdOperator(strategy CmdStrategy) *CmdOperator {
	return &CmdOperator{
		CmdStrategy: strategy,
	}
}

func (operator *CmdOperator) DoCmdStrategy(conn redis.ClientConn, args [][]byte) tcp.Info {
	return operator.CmdStrategy.Do(conn, args)
}

// PingStrategy ping策略
type PingStrategy struct{}

func (*PingStrategy) Do(_ redis.ClientConn, args [][]byte) tcp.Info {
	if len(args) == 0 {
		return &exchange.PongResponse{}
	} else if len(args) == 1 {
		return exchange.NewStatusInfo(string(args[0]))
	} else {
		return errors.NewStandardError("wrong number of arguments for 'ping' command")
	}
}

// AuthStrategy 认证策略
type AuthStrategy struct{}

func (*AuthStrategy) Do(conn redis.ClientConn, args [][]byte) tcp.Info {
	if len(args) != 1 {
		return errors.NewStandardError("Wrong number of arguments for 'auth' command")
	}
	ok := &exchange.OkResponse{}
	if config.GetPassword() == "" {
		//服务端无密码，不予认证
		return ok
	}
	passwd := string(args[0])
	conn.SetPassword(passwd)
	if config.GetPassword() != passwd {
		return errors.NewStandardError("invalid password")
	}
	return ok
}

// InfoStrategy redis服务信息策略
type InfoStrategy struct{}

func (*InfoStrategy) Do(_ redis.ClientConn, args [][]byte) tcp.Info {
	if len(args) == 1 {
		infoCommandList := [...]string{"server", "client"}
		var allSection []byte
		for _, s := range infoCommandList {
			allSection = append(allSection, GetCustomizeRedisInfo(s)...)
		}
		return exchange.NewBulkInfo(allSection)
	} else if len(args) == 2 {
		section := strings.ToLower(string(args[1]))
		switch section {
		case "server":
			rep := GetCustomizeRedisInfo("server")
			return exchange.NewBulkInfo(rep)
		case "client":
			return exchange.NewBulkInfo(GetCustomizeRedisInfo("client"))
		default:
			return exchange.NewNullBulkRequest()
		}
	}
	return errors.NewStandardError("ERR wrong number of arguments for 'info' command")
}

// GetCustomizeRedisInfo 返回redis service 信息
func GetCustomizeRedisInfo(redisType string) []byte {
	startUpTime := time.Since(config.GetStartUpTime()) / time.Second
	switch redisType {
	case "server":
		s := fmt.Sprintf("# Server\r\n"+
			"goredis_version:%s\r\n"+
			"goredis_mode:%s\r\n"+
			"os:%s %s\r\n"+
			"arch_bits:%d\r\n"+
			"go_version:%s\r\n"+
			"process_id:%d\r\n"+
			"redis_port:%d\r\n"+
			"uptime_in_seconds:%d\r\n"+
			"uptime_in_days:%d\r\n",
			config.Version,
			config.GetServerType(),
			runtime.GOOS, runtime.GOARCH,
			32<<(^uint(0)>>63),
			runtime.Version(),
			os.Getpid(),
			config.GetPort(),
			startUpTime,
			startUpTime/time.Duration(3600*24))
		return []byte(s)
	case "client":
		s := fmt.Sprintf("# Clients\r\n"+
			"connected_clients:%d\r\n",
			//TODO 连接数
			1,
			//1,
			//2,
			//3,
		)
		return []byte(s)
	}

	return []byte("")
}
