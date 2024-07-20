package plog

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	builders = sync.Pool{New: func() interface{} { return new(strings.Builder) }}
)

type MetaData struct {
	tt     time.Time
	prefix string
	file   string
	line   int
	fname  string
	level  string
	msg    string
}

func (d *MetaData) GetTime() time.Time {
	return d.tt
}

func (d *MetaData) ToString() string {
	builder := builders.Get().(*strings.Builder)
	defer builders.Put(builder)

	builder.Reset()
	// 格式化输出
	builder.WriteString("[")
	builder.WriteString(d.tt.Format("2006-01-02 15:04:05.999"))
	builder.WriteString("]	")
	builder.WriteString("[")
	builder.WriteString(d.prefix)
	builder.WriteString("]	")
	builder.WriteString("[")
	builder.WriteString(d.level)
	builder.WriteString("]	")
	builder.WriteString(d.file)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(d.line))
	builder.WriteString("	")
	builder.WriteString(d.fname)
	builder.WriteString("	")
	builder.WriteString(d.msg)
	builder.WriteString("\n")
	return builder.String()
}
