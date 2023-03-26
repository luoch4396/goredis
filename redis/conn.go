package redis

import (
	"github.com/go-netty/go-netty"
	"io"
	"sync"
)

//对象池
var connPool = sync.Pool{
	New: func() interface{} {
		return &ClientConn{}
	},
}

// Builder ClientConn 建造者
type Builder interface {
	Build() *ClientConn
	BuildChannel(channel netty.Channel) *ClientConnBuilder
	BuildStdOut(w io.Writer) *ClientConnBuilder
}

type ClientConnBuilder struct {
	logger *ClientConn
}

type ClientConn struct {
	//连接
	channel netty.Channel
	//锁
	lock sync.Locker
	//密码
	password string
	// 被选中的db
	selectedDB int
	//当前可能的角色
	role string
}

func (conn *ClientConn) Close() {

}
