package uerror

import (
	"fmt"
	"runtime"
)

type UError struct {
	filename string // 文件名
	line     int    // 文件行号
	funcname string // 函数名
	code     int32  // 错误码
	msg      string // 错误
}

func New(depth int, code int32, format string, args ...interface{}) *UError {
	// 获取调用堆栈
	pc, file, line, _ := runtime.Caller(depth)
	funcName := runtime.FuncForPC(pc).Name()
	// 返回错误
	return &UError{
		filename: file,
		line:     line,
		funcname: funcName,
		code:     code,
		msg:      fmt.Sprintf(format, args...),
	}
}

func (e *UError) Error() string {
	return fmt.Sprintf("%s:%d\t%s\t%s", e.filename, e.line, e.funcname, e.msg)
}
