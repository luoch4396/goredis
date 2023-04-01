package redis

import (
	"github.com/go-netty/go-netty"
	"goredis/pkg/utils"
	"goredis/pool/gopool"
	"sync"
	"time"
)

// 对象池
var connPool = sync.Pool{
	New: func() interface{} {
		return &ClientConn{}
	},
}

// Builder ClientConn 建造者
type Builder interface {
	Build() *ClientConn
	BuildChannel(channel netty.Channel) *ClientConnBuilder
	BuildPwd(pwd string) *ClientConnBuilder
	BuildDBIndex(dbIndex int) *ClientConnBuilder
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
	//等待客户端关闭
	waitFinished sync.WaitGroup
}

func NewClientConnBuilder() Builder {
	var build = &ClientConnBuilder{}
	//尝试从连接池获取连接
	c, ok := connPool.Get().(*ClientConn)
	if !ok {
		build.conn = &ClientConn{}
		return build
	}
	build.conn = c
	return build
}

// ClientConnBuilder 建造者
type ClientConnBuilder struct {
	conn *ClientConn
}

func (builder *ClientConnBuilder) Build() *ClientConn {
	builder.conn.lock = utils.NewLightLock(16)
	return builder.conn
}

func (builder *ClientConnBuilder) BuildChannel(channel netty.Channel) *ClientConnBuilder {
	builder.conn.channel = channel
	return builder
}

func (builder *ClientConnBuilder) BuildPwd(pwd string) *ClientConnBuilder {
	builder.conn.password = pwd
	return builder
}

func (builder *ClientConnBuilder) BuildDBIndex(dbIndex int) *ClientConnBuilder {
	builder.conn.selectedDB = dbIndex
	return builder
}

// Write 发送返回数据
func (conn *ClientConn) Write(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	conn.waitFinished.Add(1)
	defer func() {
		conn.waitFinished.Done()
	}()
	return conn.channel.Write(b)
}

func (conn *ClientConn) Name() string {
	if conn.channel != nil {
		return conn.channel.RemoteAddr()
	}
	return ""
}

// Close 关闭redis客户端连接
func (conn *ClientConn) Close() {
	//_ = conn.channel.Close()
	conn.password = ""
	conn.selectedDB = 0
	connPool.Put(conn)
}

func waitTimeout(conn *ClientConn, timeout time.Duration) bool {
	c := make(chan struct{}, 1)
	gopool.Go(func() {
		defer close(c)
		conn.waitFinished.Wait()
		c <- struct{}{}
	})
	select {
	case <-c:
		//正常
		return false
	case <-time.After(timeout):
		//超时
		return true
	}
}
