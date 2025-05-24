package mlog

import (
	"sync/atomic"
	"time"
)

const (
	LOG_DEBUG = 1
	LOG_TRACE = 2
	LOG_WARN  = 3
	LOG_INFO  = 4
	LOG_ERROR = 5
	LOG_FATAL = 6
)

type IWriter interface {
	Write(IFormat) error
	Close() error
}

type IFormat interface {
	GetTime() time.Time
	GetString() string
}

type FormatFunc func(int, uint32, string, ...interface{}) IFormat

type Logger struct {
	level uint32
	f     FormatFunc
	w     IWriter
}

func NewLogger(level uint32, f FormatFunc, w IWriter) *Logger {
	return &Logger{level, f, w}
}

func (d *Logger) Close() error {
	return d.w.Close()
}

func (d *Logger) SetLevel(level uint32) {
	atomic.StoreUint32(&d.level, level)
}

func (d *Logger) Trace(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level) <= LOG_TRACE {
		d.w.Write(d.f(skip+1, LOG_TRACE, format, args...))
	}
}

func (d *Logger) Debug(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level) <= LOG_DEBUG {
		d.w.Write(d.f(skip+1, LOG_DEBUG, format, args...))
	}
}

func (d *Logger) Warn(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level) <= LOG_WARN {
		d.w.Write(d.f(skip+1, LOG_WARN, format, args...))
	}
}

func (d *Logger) Info(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level) <= LOG_INFO {
		d.w.Write(d.f(skip+1, LOG_INFO, format, args...))
	}
}

func (d *Logger) Error(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level) <= LOG_ERROR {
		d.w.Write(d.f(skip+1, LOG_ERROR, format, args...))
	}
}

func (d *Logger) Fatal(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level) <= LOG_FATAL {
		d.w.Write(d.f(skip+1, LOG_FATAL, format, args...))
	}
}
