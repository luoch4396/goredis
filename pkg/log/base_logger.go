package log

import (
	"io"
	"log"
)

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
	stdLog *log.Logger
	level  Level
	w      io.Writer
}

type FileSettings struct {
	fileName string
	path     string
}

type LoggerBuilder struct {
}
