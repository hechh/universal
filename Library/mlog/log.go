package mlog

var (
	logger = NewLogger(LOG_DEBUG, NewFormat, &StdWriter{})
)

func Init(logPath, logName string, level uint32) {
	w := NewWriter(logPath, logName, 1024)
	logger = NewLogger(level, NewFormat, w)
}

func Close() error {
	return logger.Close()
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
