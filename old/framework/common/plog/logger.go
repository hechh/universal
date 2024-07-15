package plog

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

type IWriter interface {
	Close()
	Write(time.Time, []byte)
}

type Logger struct {
	serverId   int32
	serverName string
	level      uint32
	w          IWriter
}

func NewLogger(level uint32, id int32, name string, w IWriter) *Logger {
	return &Logger{
		serverId:   id,
		serverName: name,
		level:      level,
		w:          w,
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

func (d *Logger) ToString() string {
	return fmt.Sprintf("%s%02d", d.serverName, d.serverId)
}

func (d *Logger) output(skip int, level uint32, msg string) {
	// 获取调用堆栈
	pc, file, line, _ := runtime.Caller(skip + 1)
	funcName := path.Base(runtime.FuncForPC(pc).Name())
	// 格式化输出
	tt := time.Now()
	var builder bytes.Buffer
	builder.WriteString("[")
	builder.WriteString(tt.Format("2006-01-02 15:04:05.999"))
	builder.WriteString("] [")
	builder.WriteString(d.serverName)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(int(d.serverId)))
	builder.WriteString("] [")
	builder.WriteString(levelToString(level))
	builder.WriteString("]\t")
	builder.WriteString(file)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(line))
	builder.WriteString("\t")
	builder.WriteString(funcName)
	builder.WriteString("\t")
	builder.WriteString(msg)
	builder.WriteString("\n")
	// 日志文件
	d.w.Write(tt, builder.Bytes())
}
