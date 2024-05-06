package fbasic

import (
	"fmt"
	"runtime"
	"universal/common/pb"
)

type UError struct {
	code     int32  // 错误码
	file     string // 文件名
	line     int    // 文件行号
	funcName string // 函数名
	errMsg   string // 错误
}

func NewUError(skip int, code pb.ErrorCode, msgs ...interface{}) *UError {
	pc, file, line, _ := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()
	ret := &UError{
		code:     int32(code),
		file:     file,
		line:     line,
		funcName: funcName,
	}
	for _, msg := range msgs {
		switch v := msg.(type) {
		case string:
			ret.errMsg += fmt.Sprintf(" | %s", v)
		case *UError:
			ret.errMsg += fmt.Sprintf(" | %s", v.Error())
		case error:
			ret.errMsg += fmt.Sprintf(" | %s", v.Error())
		default:
			ret.errMsg += (" | " + fmt.Sprint(v))
		}
	}
	return ret
}

func GetCodeMsg(err error) (code int32, errmsg string) {
	switch vv := err.(type) {
	case *UError:
		code = vv.GetCode()
		errmsg = vv.GetErrMsg()
	case nil:
		code, errmsg = int32(pb.ErrorCode_Success), ""
	default:
		code = -1
		errmsg = err.Error()
	}
	return
}

func (d *UError) Append(format string, args ...interface{}) {
	d.errMsg += (" | " + fmt.Sprintf(format, args...))
}

func (d *UError) GetCode() int32 {
	return d.code
}

func (d *UError) GetErrMsg() string {
	return d.Error()
}

func (d *UError) Error() string {
	ctype := pb.ErrorCode(d.code)
	return fmt.Sprintf("%s:%d %s \n\t%s(%d): %s\n", d.file, d.line, d.funcName, ctype.String(), d.code, d.errMsg)
}
