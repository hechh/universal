package mlog

import (
	"fmt"
	"poker_server/common/yaml"
)

var log ILog

const (
	DEBUG = iota
	WARN
	INFO
	ERROR
	FATAL
)

var (
	level2Str = map[int32]string{
		DEBUG: "debug",
		WARN:  "warn",
		INFO:  "info",
		ERROR: "error",
		FATAL: "fatal",
	}
	str2Level = map[string]int32{
		"debug": DEBUG,
		"warn":  WARN,
		"info":  INFO,
		"error": ERROR,
		"fatal": FATAL,
	}
)

type ILog interface {
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

func init() {
	logger, err := NewLogger(&yaml.NodeConfig{
		Env:      "develop",
		LogLevel: "debug",
		LogPath:  "./log",
	})
	if err != nil {
		panic(fmt.Sprintf("日志库初始化失败: %v", err))
	}
	log = logger
}

func Init(cfg *yaml.NodeConfig) error {
	logger, err := NewLogger(cfg)
	if err != nil {
		return fmt.Errorf("日志库初始化失败: %v", err)
	}
	log = logger
	return nil
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
