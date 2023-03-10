package tcp

// RequestInfo 请求数据定义
type RequestInfo interface {
	RequestInfo() []byte
}

// Request 请求模型定义
type Request struct {
	Data  RequestInfo
	Error error
}

type Requests struct {
	Datas []Request
}
