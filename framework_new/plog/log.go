package plog

import (
	"fmt"
	"time"
	"universal/framework_new/common/base"
)

const (
	LOG_TRACE  = 0x01
	LOG_DEBUG  = 0x02
	LOG_WARN   = 0x04
	LOG_INFO   = 0x08
	LOG_ERROR  = 0x10
	LOG_FATAL  = 0x20
	LOG_DEFAUL = LOG_INFO | LOG_ERROR | LOG_FATAL
	LOG_ALL    = 0xff
)

func LevelToString(level uint32) string {
	switch level {
	case LOG_TRACE:
		return "TRACE"
	case LOG_DEBUG:
		return "DEBUG"
	case LOG_WARN:
		return "WARN"
	case LOG_INFO:
		return "INFO"
	case LOG_ERROR:
		return "ERROR"
	case LOG_FATAL:
		return "FATAL"
	}
	return ""
}

var (
	logger *Logger
	stdout *Stdout
)

func init() {
	logger = NewLogger(LOG_ALL, 0, "", NewWriter(0, "", "log"))
	stdout = NewStdout()
}

func Init(level uint32, id int32, name, path string) {
	logger = NewLogger(level, id, name, NewWriter(id, name, path))
}

func Gout(format string, args ...interface{}) {
	// 获取调用堆栈
	msg := fmt.Sprintf(format, args...)
	msg = fmt.Sprintf("[%s][%s] %s", time.Now().Format("2006-01-02 15:04:05.000"), logger.ToString(), msg)
	stdout.Write(base.StringToBytes(msg))
}

func Trace(format string, args ...interface{}) {
	logger.Trace(1, format, args...)
}

func Debug(format string, args ...interface{}) {
	logger.Debug(1, format, args...)
}

func Warn(format string, args ...interface{}) {
	logger.Warn(1, format, args...)
}

func Info(format string, args ...interface{}) {
	logger.Info(1, format, args...)
}

func Error(format string, args ...interface{}) {
	logger.Error(1, format, args...)
}

func Fatal(format string, args ...interface{}) {
	logger.Fatal(1, format, args...)
}

func TraceSkip(skip int, format string, args ...interface{}) {
	logger.Trace(skip+1, format, args...)
}

func DebugSkip(skip int, format string, args ...interface{}) {
	logger.Debug(skip+1, format, args...)
}

func WarnSkip(skip int, format string, args ...interface{}) {
	logger.Warn(skip+1, format, args...)
}

func InfoSkip(skip int, format string, args ...interface{}) {
	logger.Info(skip+1, format, args...)
}

func ErrorSkip(skip int, format string, args ...interface{}) {
	logger.Error(skip+1, format, args...)
}

func FatalSkip(skip int, format string, args ...interface{}) {
	logger.Fatal(skip+1, format, args...)
}
