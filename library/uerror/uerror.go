package uerror

import (
	"fmt"
	"runtime"

	"github.com/spf13/cast"
)

type ICode interface {
	String() string
}

type UError struct {
	file     string      // 文件名
	line     int         // 文件行号
	funcname string      // 函数名
	code     interface{} // 错误码
	msg      string      // 错误
	err      error       // 错误
}

func New(depth int, code interface{}, format string, args ...interface{}) *UError {
	// 获取调用堆栈
	pc, file, line, _ := runtime.Caller(depth)
	funcName := runtime.FuncForPC(pc).Name()
	// 返回错误
	return &UError{
		file:     file,
		line:     line,
		funcname: funcName,
		code:     code,
		msg:      fmt.Sprintf(format, args...),
	}
}

func Error(depth int, code interface{}, err error) *UError {
	// 获取调用堆栈
	pc, file, line, _ := runtime.Caller(depth)
	funcName := runtime.FuncForPC(pc).Name()
	// 返回错误
	return &UError{
		file:     file,
		line:     line,
		funcname: funcName,
		code:     code,
		err:      err,
	}
}

func (e *UError) GetCode() int32 {
	return cast.ToInt32(e.code)
}

func (e *UError) GetMsg() string {
	if e.err != nil {
		return e.err.Error()
	}
	return e.msg
}

func (e *UError) Error() string {
	switch vv := e.code.(type) {
	case ICode:
		if e.err != nil {
			return fmt.Sprintf("[%s] file:%s, line:%d, fname:%s, error:%v", vv.String(), e.file, e.line, e.funcname, e.err)
		}
		return fmt.Sprintf("[%s] file:%s, line:%d, fname:%s, error:%s", vv.String(), e.file, e.line, e.funcname, e.msg)
	case uint32, int32, int64, uint64, int, uint:
		if e.err != nil {
			return fmt.Sprintf("[%d] file:%s, line:%d, fname:%s, error:%v", vv, e.file, e.line, e.funcname, e.err)
		}
		return fmt.Sprintf("[%d] file:%s, line:%d, fname:%s, error:%s", vv, e.file, e.line, e.funcname, e.msg)
	}
	if e.err != nil {
		return fmt.Sprintf("[%v] file:%s, line:%d, fname:%s, error:%v", e.code, e.file, e.line, e.funcname, e.err)
	}
	return fmt.Sprintf("[%v] file:%s, line:%d, fname:%s, error:%s", e.code, e.file, e.line, e.funcname, e.msg)
}
