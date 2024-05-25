package ulog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"universal/framework/common/fbasic"
)

type LogType uint32

const (
	DEBUG LogType = iota
	INFO  LogType = iota
	ERROR LogType = iota
	FATAL LogType = iota
)

var (
	ulogObj ULog
)

type ULog struct {
	sync.RWMutex
	level uint32
	dir   string
	name  string
	value atomic.Value
}

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

func InitULog(level uint32, dir string, file string) error {
	ulogObj.level = level
	ulogObj.dir = dir
	ulogObj.name = file
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

func (d *ULog) Close() {
	if fb, ok := d.value.Load().(*os.File); ok && fb != nil {
		fb.Close()
	}
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
	os.Stdout.Write(fbasic.StrToBytes(msg))
	// 日志文件
	d.write(fbasic.StrToBytes(msg))
}

func (d *ULog) write(data []byte) (n int, err error) {
	var fb *os.File
	var ok bool
	// 获取文件句柄
	if fb, ok = d.value.Load().(*os.File); ok && fb != nil {
		// 判断句柄是否有效
		if !d.check(fb) {
			// 重新打开文件
			if fb, err = d.getFile(); err != nil {
				return
			}
		}
	} else {
		// 重新打开文件
		if fb, err = d.getFile(); err != nil {
			return
		}
	}
	// 写入数据
	return fb.Write(data)
}

// 检查句柄是否有效
func (d *ULog) check(fb *os.File) bool {
	var err error
	var st1, st2 os.FileInfo
	// 文件被删除
	if st1, err = fb.Stat(); err != nil {
		return false
	}
	// 判断是否需要更新文件
	if st2, err = os.Stat(filepath.Join(d.getDir(), d.getFileName())); err != nil {
		return false
	}
	return os.SameFile(st1, st2)
}

func (d *ULog) getFile() (fb *os.File, err error) {
	var ok bool
	d.Lock()
	defer d.Unlock()
	if fb, ok = d.value.Load().(*os.File); ok && fb != nil {
		return
	}
	if fb, err = d.newFile(); err != nil {
		return nil, err
	}
	return
}

func (d *ULog) newFile() (fb *os.File, err error) {
	// 判断路径是否存在
	path := d.getDir()
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.FileMode(0755)); err != nil {
			return nil, err
		}
	}
	// 创建文件
	fileName := d.getFileName()
	if fb, err = os.OpenFile(filepath.Join(path, fileName), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644); err != nil {
		return nil, err
	}
	return
}

func (d *ULog) getDir() string {
	if len(d.dir) <= 0 {
		return "./"
	}
	return d.dir
}

func (d *ULog) getFileName() string {
	tt := time.Now()
	return fmt.Sprintf("%s-%04d%02d%02d.log", d.name, tt.Year(), tt.Month(), tt.Day())
}
