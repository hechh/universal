package mlog

var logger *Logger = NewLogger(LOG_DEBUG, &StdWriter{})

func Init(logpath, logname string, level int32) {
	logger.Close()
	logger = NewLogger(level, NewLogWriter(logpath, logname, 1024))
}

func Close() error {
	return logger.Close()
}

func SetLevel(level int32) {
	logger.SetLevel(level)
}

func Debug(skip int, format string, args ...interface{}) {
	logger.Debug(skip+1, format, args...)
}

func Trace(skip int, format string, args ...interface{}) {
	logger.Trace(skip+1, format, args...)
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

// --------Debugf---------
func Debugf(format string, args ...interface{}) {
	logger.Debug(1, format, args...)
}

func Tracef(format string, args ...interface{}) {
	logger.Trace(1, format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warn(1, format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Info(1, format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Error(1, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatal(1, format, args...)
}
