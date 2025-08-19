package mlog

import "strings"

const (
	LOG_TRACE = 1
	LOG_DEBUG = 2
	LOG_WARN  = 3
	LOG_INFO  = 4
	LOG_ERROR = 5
	LOG_FATAL = 6
)

var (
	logger = NewLogger(LOG_DEBUG, &StdWriter{})
)

func StringToLevel(str string) int32 {
	switch strings.ToUpper(str) {
	case "TRACE":
		return LOG_TRACE
	case "DEBUG":
		return LOG_DEBUG
	case "WARN":
		return LOG_WARN
	case "INFO":
		return LOG_INFO
	case "ERROR":
		return LOG_ERROR
	case "FATAL":
		return LOG_FATAL
	}
	return LOG_WARN
}

func Init(logpath, logname, level string) {
	logger.Close()
	logger = NewLogger(StringToLevel(level), NewLogWriter(logpath, logname, 1024))
}

func Close() error {
	return logger.Close()
}

func SetLevel(level int32) {
	logger.SetLevel(level)
}

func Tracef(format string, args ...interface{}) {
	logger.Trace(-1, format, args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debug(-1, format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warn(-1, format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Info(-1, format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Error(-1, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatal(-1, format, args...)
}

// --------Debugf---------
func Trace(skip int, format string, args ...interface{}) {
	logger.Trace(skip+1, format, args...)
}

func Debug(skip int, format string, args ...interface{}) {
	logger.Debug(skip+1, format, args...)
}

func Warn(skip int, format string, args ...interface{}) {
	logger.Warn(skip+1, format, args...)
}

func Info(skip int, format string, args ...interface{}) {
	logger.Info(skip+1, format, args...)
}

func Error(skip int, format string, args ...interface{}) {
	logger.Error(skip+1, format, args...)
}

func Fatal(skip int, format string, args ...interface{}) {
	logger.Fatal(skip+1, format, args...)
}
