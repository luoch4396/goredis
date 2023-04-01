package db

import (
	"goredis/interface/tcp"
	"goredis/pkg/errors"
	"goredis/redis/conn"
	"goredis/redis/exchange"
)

// CmdStrategy 命令解析策略接口
type CmdStrategy interface {
	Do(conn *conn.ClientConn, args [][]byte) tcp.Info
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

func (operator *CmdOperator) DoCmdStrategy(conn *conn.ClientConn, args [][]byte) tcp.Info {
	return operator.CmdStrategy.Do(conn, args)
}

// PingStrategy ping策略
type PingStrategy struct{}

func (*PingStrategy) Do(_ *conn.ClientConn, args [][]byte) tcp.Info {
	if len(args) == 0 {
		return &exchange.PongResponse{}
	} else if len(args) == 1 {
		return exchange.NewStatusInfo(string(args[0]))
	} else {
		return errors.NewStandardError("ERR wrong number of arguments for 'ping' command")
	}
}

// AuthStrategy 认证策略
type AuthStrategy struct{}

func (*AuthStrategy) Do(conn *conn.ClientConn, args [][]byte) tcp.Info {
	return nil
}

// InfoStrategy redis服务信息策略
type InfoStrategy struct{}

func (*InfoStrategy) Do(conn *conn.ClientConn, args [][]byte) tcp.Info {
	return nil
}
