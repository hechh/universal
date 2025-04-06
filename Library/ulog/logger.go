package ulog

import (
	"fmt"
	"hego/framework/basic"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	builders = sync.Pool{New: func() interface{} { return new(strings.Builder) }}
)

type IWriter interface {
	Close()
	Write(time.Time, string)
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
	tt := basic.GetNowTime()
	builder := builders.Get().(*strings.Builder)
	builder.Reset()
	// 格式化输出
	builder.WriteString("[")
	builder.WriteString(tt.Format("2006-01-02 15:04:05.000"))
	builder.WriteString("]	")
	builder.WriteString("[")
	builder.WriteString(d.prefix)
	builder.WriteString("]	")
	builder.WriteString("[")
	builder.WriteString(levelToString(level))
	builder.WriteString("]	")
	builder.WriteString(file)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(line))
	builder.WriteString("	")
	builder.WriteString(funcName)
	builder.WriteString("	")
	builder.WriteString(msg)
	builder.WriteString("\n")
	// 写入数据
	d.w.Write(tt, builder.String())
	builders.Put(builder)
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
