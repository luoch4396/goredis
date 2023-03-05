package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	logFile            *os.File
	defaultPrefix      = ""
	defaultCallerDepth = 2
	logger             *log.Logger
	syn                sync.Mutex
	logPrefix          = ""
	levelFlags         = []string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}
)

const (
	DEBUG   = 0
	INFO    = 1
	WARNING = 2
	ERROR   = 3
	FATAL   = 4
)

type logLevel int

// NewLog4j 初始化日志参数
func NewLog4j() {

}

func Debug(message string, v ...interface{}) {
	basePrintLog(DEBUG, message, v)
}

func Info(message string, v ...interface{}) {
	basePrintLog(INFO, message, v)
}

func Warning(message string, v ...interface{}) {
	basePrintLog(WARNING, message, v)
}

func Error(message string, v ...interface{}) {
	basePrintLog(ERROR, message, v)
}

func Fatal(message string, v ...interface{}) {
	basePrintLog(FATAL, message, v)
}

func basePrintLog(logLevel logLevel, message string, v ...interface{}) {
	syn.Lock()
	defer syn.Unlock()
	setPrefix(logLevel)
	logger.Println(fmt.Sprintf(message, v...))
}

func setPrefix(level logLevel) {
	//TODO 这边据说runtime.Caller获取行号有性能问题
	_, file, line, ok := runtime.Caller(defaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d] ", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s] ", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}
