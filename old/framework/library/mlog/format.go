package mlog

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"time"
	"universal/framework/library/util"
)

var (
	builders = util.NewSyncPool[bytes.Buffer]()
)

type Format struct {
	t   time.Time
	msg string
}

func (d *Format) GetTime() time.Time {
	return d.t
}

func (d *Format) GetString() string {
	return d.msg
}

func NewFormat(skip int, level uint32, format string, args ...interface{}) IFormat {
	// 获取调用堆栈
	pc, file, line, _ := runtime.Caller(skip + 1)
	fname := path.Base(runtime.FuncForPC(pc).Name())
	tt := time.Now()
	builder := builders.Get().(*bytes.Buffer)
	defer builders.Put(builder)
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
	return &Format{tt, builder.String()}
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
