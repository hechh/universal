package mlog

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"sync"
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
	Write(time.Time, []byte) error
	Close() error
}

type Logger struct {
	level   int32
	writers []IWriter
	pool    sync.Pool
}

func NewLogger(level int32, writers ...IWriter) *Logger {
	return &Logger{
		level:   level,
		writers: writers,
		pool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, 1024))
			},
		},
	}
}

func (l *Logger) Close() error {
	var err error
	for _, writer := range l.writers {
		if e := writer.Close(); e != nil {
			err = e
		}
	}
	return err
}

func (l *Logger) SetLevel(level int32) {
	atomic.StoreInt32(&l.level, level)
}

func (d *Logger) Trace(skip int, format string, args ...interface{}) {
	if atomic.LoadInt32(&d.level) <= LOG_TRACE {
		d.output(skip+1, LOG_TRACE, format, args...)
	}
}

func (d *Logger) Debug(skip int, format string, args ...interface{}) {
	if atomic.LoadInt32(&d.level) <= LOG_DEBUG {
		d.output(skip+1, LOG_DEBUG, format, args...)
	}
}

func (d *Logger) Warn(skip int, format string, args ...interface{}) {
	if atomic.LoadInt32(&d.level) <= LOG_WARN {
		d.output(skip+1, LOG_WARN, format, args...)
	}
}

func (d *Logger) Info(skip int, format string, args ...interface{}) {
	if atomic.LoadInt32(&d.level) <= LOG_INFO {
		d.output(skip+1, LOG_INFO, format, args...)
	}
}

func (d *Logger) Error(skip int, format string, args ...interface{}) {
	if atomic.LoadInt32(&d.level) <= LOG_ERROR {
		d.output(skip+1, LOG_ERROR, format, args...)
	}
}

func (d *Logger) Fatal(skip int, format string, args ...interface{}) {
	if atomic.LoadInt32(&d.level) <= LOG_FATAL {
		d.output(skip+1, LOG_FATAL, format, args...)
	}
}

func (d *Logger) output(skip int, level int32, format string, args ...interface{}) {
	// 获取调用堆栈
	pc, file, line, _ := runtime.Caller(skip + 1)
	fname := path.Base(runtime.FuncForPC(pc).Name())
	tt := time.Now()

	builder := d.pool.Get().(*bytes.Buffer)
	defer d.pool.Put(builder)
	builder.Reset()

	// 格式化输出
	builder.WriteString(tt.Format("2006-01-02 15:04:05.000"))
	builder.WriteByte(' ')
	builder.WriteString(levelToString(level))
	builder.WriteByte(' ')
	builder.WriteString(file)
	builder.WriteByte(':')
	builder.WriteString(strconv.Itoa(line))
	builder.WriteByte(' ')
	builder.WriteString(fname)
	builder.WriteByte(' ')
	builder.WriteString(fmt.Sprintf(format, args...))
	builder.WriteByte('\n')
	for _, ww := range d.writers {
		ww.Write(tt, builder.Bytes())
	}
}

func levelToString(level int32) string {
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
