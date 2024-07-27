package uerror

import (
	"fmt"
	"runtime"
)

type UError struct {
	file     string // 文件名
	line     int    // 文件行号
	funcName string // 函数名
	code     int32  // 错误码
	errMsg   string // 错误
}

func NewUError(skip int, code int32, format string, msgs ...interface{}) *UError {
	// 获取调用堆栈
	pc, file, line, _ := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()
	// 返回错误
	return &UError{
		code:     code,
		file:     file,
		line:     line,
		funcName: funcName,
		errMsg:   fmt.Sprintf(format, msgs...),
	}
}

func GetCodeMsg(err error) (code int32, errmsg string) {
	switch vv := err.(type) {
	case *UError:
		code, errmsg = vv.Code(), vv.Error()
	case error:
		code, errmsg = -1, err.Error()
	}
	return
}

func (d *UError) Append(format string, args ...interface{}) {
	if len(d.errMsg) > 0 {
		d.errMsg += fmt.Sprintf(" | "+format, args...)
	} else {
		d.errMsg = fmt.Sprintf(format, args...)
	}
}

func (d *UError) Code() int32 {
	return d.code
}

func (d *UError) Error() string {
	return d.errMsg
}

func (d *UError) ToString() string {
	return fmt.Sprintf("%s:%d\t%s\t%d: %s", d.file, d.line, d.funcName, d.code, d.errMsg)
}
