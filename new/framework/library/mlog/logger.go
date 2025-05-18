package mlog

import (
	"poker_server/common/yaml"
	"strings"
	"sync/atomic"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 分装zap日志库接口使用
type Logger struct {
	logger *zap.SugaredLogger
	level  int32
}

func NewLogger(cfg *yaml.NodeConfig) (*Logger, error) {
	var config zap.Config
	switch cfg.Env {
	case "release":
		config = zap.NewProductionConfig()
	default:
		config = zap.NewDevelopmentConfig()
	}
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	config.OutputPaths = []string{cfg.LogPath, "stdout"}
	config.ErrorOutputPaths = []string{cfg.LogPath, "stderr"}
	//config.OutputPaths = append(config.OutputPaths, cfg.LogPath)
	//config.ErrorOutputPaths = []string{cfg.LogPath}

	// 创建日志
	logger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(2))
	if err != nil {
		return nil, err
	}
	return &Logger{logger: logger.Sugar(), level: str2Level[strings.ToLower(cfg.LogLevel)]}, nil
}

func (d *Logger) Debugf(format string, args ...interface{}) {
	if atomic.LoadInt32(&d.level) <= int32(zapcore.DebugLevel)+1 {
		d.logger.Debugf(format, args...)
	}
}

func (d *Logger) Warnf(format string, args ...interface{}) {
	if atomic.LoadInt32(&d.level) <= int32(zapcore.WarnLevel)+1 {
		d.logger.Warnf(format, args...)
	}
}

func (d *Logger) Infof(format string, args ...interface{}) {
	if atomic.LoadInt32(&d.level) <= int32(zapcore.InfoLevel)+1 {
		d.logger.Infof(format, args...)
	}
}

func (d *Logger) Errorf(format string, args ...interface{}) {
	if atomic.LoadInt32(&d.level) <= int32(zapcore.ErrorLevel)+1 {
		d.logger.Errorf(format, args...)
	}
}

func (d *Logger) Fatalf(format string, args ...interface{}) {
	if atomic.LoadInt32(&d.level) <= int32(zapcore.FatalLevel)+1 {
		d.logger.Fatalf(format, args...)
	}
}
