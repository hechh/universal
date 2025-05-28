package mlog

import (
	"fmt"
	"os"
	"path"
	"universal/library/mlog/zap"
)

const (
	DEBUG = iota
	WARN
	INFO
	ERROR
	FATAL
)

var log ILog

var str2Level = map[string]int32{
	"debug": DEBUG,
	"warn":  WARN,
	"info":  INFO,
	"error": ERROR,
	"fatal": FATAL,
}

type ILog interface {
	GetLevel() int32
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

func InitDefault() {
	val := zap.NewLogger("develop", 0, "./default.log")
	switch vv := val.(type) {
	case error:
		panic(fmt.Sprintf("日志库初始化失败: %v", vv))
	case ILog:
		log = vv
	}
}

func Init(env string, level string, logfile string) error {
	os.MkdirAll(path.Dir(logfile), os.FileMode(0777))
	switch vv := zap.NewLogger(env, str2Level[level], logfile).(type) {
	case error:
		return vv
	case ILog:
		log = vv
	}
	return nil
}

func Debugf(format string, args ...interface{}) {
	if log.GetLevel() <= DEBUG {
		log.Debugf(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if log.GetLevel() <= WARN {
		log.Warnf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if log.GetLevel() <= INFO {
		log.Infof(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if log.GetLevel() <= ERROR {
		log.Errorf(format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if log.GetLevel() <= FATAL {
		log.Fatalf(format, args...)
	}
}
