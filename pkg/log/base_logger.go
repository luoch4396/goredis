package log

import (
	"io"
	"log"
)

// FormatterLogger 日志模块定义
type FormatterLogger interface {
	Debug(message string, v ...interface{})
	Info(message string, v ...interface{})
	Warning(message string, v ...interface{})
	Error(message string, v ...interface{})
	Fatal(message string, v ...interface{})
}

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var levels = []string{
	"[DEBUG] ",
	"[INFO] ",
	"[WARNING] ",
	"[ERROR] ",
	"[FATAL] ",
}

// Level 日志级别
type Level int

type Logger struct {
	stdLog *log.Logger //日志
	level  Level       //日志级别
	w      io.Writer   //日志输出
}

type FileSettings struct {
	FileName string
	Path     string
}

type Builder interface {
	Build() *Logger
	BuildLevel(lv string) *LoggerBuilder
	BuildStdOut(w io.Writer) *LoggerBuilder
	BuildFile(settings *FileSettings) *LoggerBuilder
}

type LoggerBuilder struct {
	logger Logger
}
