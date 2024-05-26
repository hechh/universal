package ulog

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

type LogType uint32

const (
	DEBUG LogType = iota
	INFO  LogType = iota
	ERROR LogType = iota
	FATAL LogType = iota
)

var (
	ulogObj *ULog
)

func LogType2String(level uint32) string {
	switch LogType(level) {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return "UNKNOWN"
}

func init() {
	ulogObj = &ULog{level: uint32(DEBUG), w: os.Stdout}
}

func InitULog(level uint32, filename string) error {
	ulogObj = &ULog{level: level, w: NewWriter(filename)}
	return nil
}

func Debug(skip int, format string, args ...interface{}) {
	ulogObj.Debug(skip+1, format, args...)
}

func Info(skip int, format string, args ...interface{}) {
	ulogObj.Info(skip+1, format, args...)
}

func Error(skip int, format string, args ...interface{}) {
	ulogObj.Error(skip+1, format, args...)
}

func Fatal(skip int, format string, args ...interface{}) {
	ulogObj.Fatal(skip+1, format, args...)
}

type ULog struct {
	sync.RWMutex
	level uint32
	w     IWriter
}

func (d *ULog) Close() {
	d.w.Close()
}

func (d *ULog) SetLevel(val uint32) {
	atomic.StoreUint32(&d.level, val)
}

func (d *ULog) Debug(skip int, format string, args ...interface{}) {
	d.Output(skip+1, uint32(DEBUG), format, args...)
}

func (d *ULog) Info(skip int, format string, args ...interface{}) {
	d.Output(skip+1, uint32(INFO), format, args...)
}

func (d *ULog) Error(skip int, format string, args ...interface{}) {
	d.Output(skip+1, uint32(ERROR), format, args...)
}

func (d *ULog) Fatal(skip int, format string, args ...interface{}) {
	d.Output(skip+1, uint32(FATAL), format, args...)
}

func (d *ULog) Output(skip int, level uint32, format string, args ...interface{}) {
	if atomic.LoadUint32(&d.level) > level {
		return
	}
	// 获取调用堆栈
	pc, file, line, _ := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()
	// 格式化输出
	msg := fmt.Sprintf("[%s] %s:%d\t%s\t%s\n", LogType2String(level), file, line, funcName, format)
	msg = fmt.Sprintf(msg, args...)
	// 标准输出
	os.Stdout.Write(strToBytes(msg))
	// 日志文件
	d.w.Write(strToBytes(msg))
}

func strToBytes(str string) []byte {
	s := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	b := &reflect.SliceHeader{Data: s.Data, Len: s.Len, Cap: s.Len}
	return *(*[]byte)(unsafe.Pointer(b))
}
