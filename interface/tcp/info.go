package tcp

type Info interface {
	Info() []byte
}

// ErrorInfo 错误
type ErrorInfo interface {
	Error() string
	Info() []byte
}

// Error 错误模型
type Error struct {
	Data  Info
	Error error
}

// Response 返回模型定义
type Response struct {
	Data  Info
	Error error
}

// Request 请求模型定义
type Request struct {
	Data  Info
	Error error
}

// Requests 多个请求
type Requests struct {
	Data  []Request
	Error error
}
