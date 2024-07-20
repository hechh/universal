package plog

import (
	"fmt"
	"path"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"universal/framework/util"
)

type IFormat interface {
	GetTime() time.Time
	ToString() string
}

type IWriter interface {
	Close()
	Write(IFormat)
}

type Logger struct {
	sync.RWMutex
	level  uint32
	prefix string
	w      IWriter
}

func NewLogger(level uint32, name string, w IWriter) *Logger {
	return &Logger{
		level:  level,
		prefix: name,
		w:      w,
	}
}

func (d *Logger) Close() {
	d.w.Close()
}

func (d *Logger) SetLevel(val uint32) {
	atomic.StoreUint32(&d.level, val)
}

func (d *Logger) Trace(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level)&LOG_TRACE != 0 {
		d.output(skip+1, LOG_TRACE, fmt.Sprintf(format, args...))
	}
}

func (d *Logger) Debug(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level)&LOG_DEBUG != 0 {
		d.output(skip+1, LOG_DEBUG, fmt.Sprintf(format, args...))
	}
}

func (d *Logger) Warn(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level)&LOG_WARN != 0 {
		d.output(skip+1, LOG_WARN, fmt.Sprintf(format, args...))
	}
}

func (d *Logger) Info(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level)&LOG_INFO != 0 {
		d.output(skip+1, LOG_INFO, fmt.Sprintf(format, args...))
	}
}

func (d *Logger) Error(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level)&LOG_ERROR != 0 {
		d.output(skip+1, LOG_ERROR, fmt.Sprintf(format, args...))
	}
}

func (d *Logger) Fatal(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level)&LOG_FATAL != 0 {
		d.output(skip+1, LOG_FATAL, fmt.Sprintf(format, args...))
	}
}

func (d *Logger) output(skip int, level uint32, msg string) {
	// 获取调用堆栈
	pc, file, line, _ := runtime.Caller(skip + 1)
	funcName := path.Base(runtime.FuncForPC(pc).Name())

	// 写入数据
	d.w.Write(&MetaData{
		tt:     util.GetNowTime(),
		prefix: d.prefix,
		file:   file,
		line:   line,
		fname:  funcName,
		level:  levelToString(level),
		msg:    msg,
	})
}

func levelToString(level uint32) string {
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
