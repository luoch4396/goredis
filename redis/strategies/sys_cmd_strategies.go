package strategies

import (
	"goredis/interface/tcp"
	"goredis/redis"
	"goredis/redis/exchange"
)

// CmdStrategy 命令解析策略接口
type CmdStrategy interface {
	Do(conn redis.ClientConn, args [][]byte) tcp.ResponseInfo
}

type CmdOperator struct {
	CmdStrategy CmdStrategy
}

func (operator *CmdOperator) DoCmdStrategy(conn redis.ClientConn, args [][]byte) tcp.ResponseInfo {
	return operator.CmdStrategy.Do(conn, args)
}

// PingStrategy ping策略
type PingStrategy struct{}

func (*PingStrategy) Do(conn redis.ClientConn, args [][]byte) tcp.ResponseInfo {
	if len(args) == 0 {
		return &exchange.PongResponse{}
	}
	//else if len(args) == 1 {
	//return &exchange.NewStatusResponse(string(args[0]))
	//} else {
	//return &exchange.NewErrResponse("ERR wrong number of arguments for 'ping' command")
	//}
	return nil
}

// AuthStrategy 认证策略
type AuthStrategy struct{}

// InfoStrategy redis服务信息策略
type InfoStrategy struct{}
