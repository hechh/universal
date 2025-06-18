package uerror

import (
	"fmt"
	"runtime"

	"github.com/spf13/cast"
)

type IErrorCode interface {
	String() string
	Number() int32
}

type ErrorCode int32

func (d ErrorCode) String() string {
	return cast.ToString(d)
}

func (d ErrorCode) Number() int32 {
	return int32(d)
}

type UError struct {
	file  string
	fname string
	line  int
	code  IErrorCode
	msg   string
}

func (ue *UError) Error() string {
	return fmt.Sprintf("[%d]\t%s:%d %s\terror:%s", ue.code, ue.file, ue.line, ue.fname, ue.msg)
}

func (ue *UError) GetCode() int32 {
	return ue.code.Number()
}

func (ue *UError) GetMsg() string {
	return ue.msg
}

func E(depth int, code int32, err error) *UError {
	if vv, ok := err.(*UError); ok {
		return vv
	}
	pc, file, line, _ := runtime.Caller(depth)
	fname := runtime.FuncForPC(pc).Name()
	return &UError{file: file, line: line, fname: fname, code: ErrorCode(code), msg: err.Error()}
}

func N(depth int, code int32, format string, args ...interface{}) *UError {
	pc, file, line, _ := runtime.Caller(depth)
	fname := runtime.FuncForPC(pc).Name()
	return &UError{file: file, line: line, fname: fname, code: ErrorCode(code), msg: fmt.Sprintf(format, args...)}
}
