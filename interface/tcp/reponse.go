package tcp

// ResponseInfo 返回数据模型定义
type ResponseInfo interface {
	ResponseInfo() []byte
}

// Response 返回模型定义
type Response struct {
	Data  ResponseInfo
	Error error
}
