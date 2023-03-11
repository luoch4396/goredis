package _interface

// FormatterLogger 日志模块定义
type FormatterLogger interface {
	Debug(message string, v ...interface{})
	Info(message string, v ...interface{})
	Warning(message string, v ...interface{})
	Error(message string, v ...interface{})
	Fatal(message string, v ...interface{})
}
