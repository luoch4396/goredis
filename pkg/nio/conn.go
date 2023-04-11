package nio

// ConnType 连接类型
type ConnType = int8

const (
	// ConnTypeTCP tcp连接
	ConnTypeTCP ConnType = iota + 1
	// ConnTypeUnix 进程连接
	ConnTypeUnix
)

func IsTcpConn() {

}
