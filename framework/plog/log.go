package plog

const (
	LOG_TRACE = 0x01
	LOG_DEBUG = 0x02
	LOG_WARN  = 0x04
	LOG_INFO  = 0x08
	LOG_ERROR = 0x10
	LOG_FATAL = 0x20
	LOG_ALL   = 0xff
)

var (
	logger *Logger
)

func init() {
	logger = NewLogger(LOG_ALL, "plog", NewWriter("./log", "output"))
}

func Close() {
	logger.Close()
}

func Init(level uint32, path, name string) {
	logger = NewLogger(level, name, NewWriter(path, name))
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
