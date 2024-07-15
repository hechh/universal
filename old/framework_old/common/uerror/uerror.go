package uerror

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cast"
)

type UError struct {
	code     int32  // 错误码
	file     string // 文件名
	line     int    // 文件行号
	funcName string // 函数名
	errMsg   string // 错误
}

func NewUErrorf(skip int, code interface{}, format string, msgs ...interface{}) *UError {
	// 获取调用堆栈
	pc, file, line, _ := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()
	// 返回错误
	return &UError{
		code:     cast.ToInt32(code),
		file:     file,
		line:     line,
		funcName: funcName,
		errMsg:   fmt.Sprintf(format, msgs...),
	}
}

func NewUError(skip int, code interface{}, msgs ...interface{}) *UError {
	// 获取错误信息
	str := []string{}
	for _, msg := range msgs {
		str = append(str, fmt.Sprintf("%v", msg))
	}
	// 获取调用堆栈
	pc, file, line, _ := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()
	// 返回错误
	return &UError{
		code:     cast.ToInt32(code),
		file:     file,
		line:     line,
		funcName: funcName,
		errMsg:   strings.Join(str, "|"),
	}
}

func (d *UError) Append(format string, args ...interface{}) {
	if len(d.errMsg) > 0 {
		d.errMsg += fmt.Sprintf("|"+format, args...)
	} else {
		d.errMsg = fmt.Sprintf(format, args...)
	}
}

func (d *UError) GetCode() int32 {
	return d.code
}

func (d *UError) GetErrMsg() string {
	return d.errMsg
}

func (d *UError) Error() string {
	return fmt.Sprintf("%s:%d\t%s\tErrorCode(%d): %s\n", d.file, d.line, d.funcName, d.code, d.errMsg)
}

func GetCodeMsg(err error) (code int32, errmsg string) {
	switch vv := err.(type) {
	case *UError:
		code, errmsg = vv.GetCode(), vv.GetErrMsg()
	case error:
		code, errmsg = -1, err.Error()
	}
	return
}
