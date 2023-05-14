package conn

import (
	"goredis/pkg/nio"
	"goredis/pkg/utils"
	"sync"
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
	BuildChannel(channel *nio.Conn) *ClientConnBuilder
	BuildPwd(pwd string) *ClientConnBuilder
	BuildDBIndex(dbIndex int) *ClientConnBuilder
}

type ClientConn struct {
	//连接
	conn *nio.Conn
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

func (builder *ClientConnBuilder) BuildChannel(channel *nio.Conn) *ClientConnBuilder {
	builder.conn.conn = channel
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
	return true
}

func (conn *ClientConn) Name() string {
	return ""
}

// Close 关闭redis客户端连接
func (conn *ClientConn) Close() error {
	//正常关闭客户端，并回收连接
	err := conn.conn.Close()
	conn.password = ""
	conn.selectedDB = 0
	connPool.Put(conn)
	return err
}

func (conn *ClientConn) GetDBIndex() int {
	return conn.selectedDB
}

func (conn *ClientConn) SetPassword(password string) {
	conn.password = password
}

func (conn *ClientConn) GetPassword() string {
	return conn.password
}
