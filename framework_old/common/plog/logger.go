package plog

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

type IWriter interface {
	Close() error
	Write([]byte) (int, error)
}

type Logger struct {
	sync.RWMutex
	serverId   int
	serverName string
	level      uint32
	w          IWriter
}

func NewLogger(srvName string, opts ...OpOption) *Logger {
	op := Op{}
	op.applyOpts(opts...)
	// 设置默认值
	if op.level <= 0 {
		op.level = LOG_ALL
	}
	if len(op.path) <= 0 {
		op.path = "./log"
	}
	return &Logger{
		serverId:   op.serverId,
		serverName: srvName,
		level:      op.level,
		w:          NewWriter(op.path, srvName),
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
		d.output(skip+1, LOG_TRACE, format, args...)
	}
}

func (d *Logger) Debug(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level)&LOG_DEBUG != 0 {
		d.output(skip+1, LOG_DEBUG, format, args...)
	}
}

func (d *Logger) Warn(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level)&LOG_WARN != 0 {
		d.output(skip+1, LOG_WARN, format, args...)
	}
}

func (d *Logger) Info(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level)&LOG_INFO != 0 {
		d.output(skip+1, LOG_INFO, format, args...)
	}
}

func (d *Logger) Error(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level)&LOG_ERROR != 0 {
		d.output(skip+1, LOG_ERROR, format, args...)
	}
}

func (d *Logger) Fatal(skip int, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level)&LOG_FATAL != 0 {
		d.output(skip+1, LOG_FATAL, format, args...)
	}
}

func (d *Logger) output(skip int, level uint32, format string, args ...interface{}) {
	// 获取调用堆栈
	pc, file, line, _ := runtime.Caller(skip + 1)
	funcName := path.Base(runtime.FuncForPC(pc).Name())
	// 格式化输出
	msg := fmt.Sprintf("[%s:%d][%s] %s:%d\t%s\t%s\n", d.serverName, d.serverId, levelToString(level), file, line, funcName, format)
	msg = fmt.Sprintf(msg, args...)
	// 日志文件
	s := *(*reflect.StringHeader)(unsafe.Pointer(&msg))
	b := &reflect.SliceHeader{Data: s.Data, Len: s.Len, Cap: s.Len}
	d.w.Write(*(*[]byte)(unsafe.Pointer(b)))
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
