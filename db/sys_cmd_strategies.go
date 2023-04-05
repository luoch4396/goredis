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
	if config.GetPassword() == "" {
		return errors.NewStandardError("Client send AUTH, but server config-password is null")
	}
	passwd := string(args[0])
	conn.SetPassword(passwd)
	if config.GetPassword() != passwd {
		return errors.NewStandardError("invalid password")
	}
	return &exchange.OkResponse{}
}

// InfoStrategy redis服务信息策略
type InfoStrategy struct{}

func (*InfoStrategy) Do(conn redis.ClientConn, args [][]byte) tcp.Info {
	return nil
}
