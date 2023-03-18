package log

import (
	"fmt"
	"goredis/pkg/utils"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// log日志默认配置
var logger = Logger{
	level: INFO,
	w:     os.Stderr,
}

func NewLoggerBuilder() Builder {
	return &LoggerBuilder{}
}

// Build 初始化日志参数
func (builder *LoggerBuilder) Build() *Logger {
	builder.logger.stdLog = log.New(logger.w, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	logger = builder.logger
	return &logger
}

func (builder *LoggerBuilder) BuildStdOut(w io.Writer) *LoggerBuilder {
	builder.logger.w = w
	return builder
}

func (builder *LoggerBuilder) BuildLevel(lv string) *LoggerBuilder {
	switch lv {
	case "DEBUG":
		builder.logger.level = DEBUG
	case "INFO":
		builder.logger.level = INFO
	case "WARNING":
		builder.logger.level = WARNING
	case "ERROR":
		builder.logger.level = ERROR
	case "FATAL":
		builder.logger.level = FATAL
	}
	return builder
}

func (builder *LoggerBuilder) BuildFile(settings *FileSettings) *LoggerBuilder {
	fileName := fmt.Sprintf("%s-%s.%s",
		settings.FileName,
		time.Now().Format("2006-01-02"),
		"logs")
	logFile, err := utils.CreateIfNotExist(fileName, settings.Path)
	if err != nil {
		err = fmt.Errorf("logger.WithFile error: %s", err)
		panic(err)
	}
	logger.w = io.MultiWriter(logger.w, logFile)
	return builder
}

func Debug(message string) {
	logger.Debug(message)
}

func Info(message string) {
	logger.Info(message)
}

func Warning(message string) {
	logger.Warning(message)
}

func Error(message string) {
	logger.Error(message)
}

func Fatal(message string) {
	logger.Fatal(message)
}

func Debugf(message string, v ...interface{}) {
	logger.Debug(message, v...)
}

func Infof(message string, v ...interface{}) {
	logger.Info(message, v...)
}

func Warningf(message string, v ...interface{}) {
	logger.Warning(message, v...)
}

func Errorf(message string, v ...interface{}) {
	logger.Error(message, v...)
}

func Fatalf(message string, v ...interface{}) {
	logger.Fatal(message, v...)
}

func (dl *Logger) Debug(message string, v ...interface{}) {
	dl.basePrintLog(DEBUG, &message, v...)
}

func (dl *Logger) Info(message string, v ...interface{}) {
	dl.basePrintLog(INFO, &message, v...)
}

func (dl *Logger) Warning(message string, v ...interface{}) {
	dl.basePrintLog(WARNING, &message, v...)
}

func (dl *Logger) Error(message string, v ...interface{}) {
	dl.basePrintLog(ERROR, &message, v...)
}

func (dl *Logger) Fatal(message string, v ...interface{}) {
	dl.basePrintLog(FATAL, &message, v...)
}

func (dl *Logger) basePrintLog(logLevel Level, message *string, v ...interface{}) {
	if dl.level > logLevel {
		return
	}
	builder := &strings.Builder{}
	_, err := builder.WriteString(levels[logLevel])
	if err != nil {
		panic(err)
		return
	}
	if message != nil {
		builder.WriteString(fmt.Sprintf(*message, v...))
	} else {
		builder.WriteString(fmt.Sprint(v...))
	}
	err = dl.stdLog.Output(4, builder.String())
	if err != nil {
		return
	}
	if logLevel == FATAL {
		//出现严重错误，程序退出
		os.Exit(1)
	}
}
