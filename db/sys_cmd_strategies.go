package db

import (
	"goredis/interface/redis"
	"goredis/interface/tcp"
	"goredis/pkg/errors"
	"goredis/redis/config"
	"goredis/redis/exchange"
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

func (*InfoStrategy) Do(conn redis.ClientConn, args [][]byte) tcp.Info {
	return nil
}

// GetCustomizeRedisInfo 返回redis service 信息
func GetCustomizeRedisInfo() {

}
