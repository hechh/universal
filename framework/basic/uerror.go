package basic

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

func NewUError(skip int, code pb.ErrorCode, msg interface{}) error {
	var errMsg string
	switch v := msg.(type) {
	case *UError:
		return v
	case nil:
		return nil
	case string:
		errMsg = v
	case error:
		errMsg = v.Error()
	}
	pc, file, line, _ := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()
	return &UError{
		code:     int32(code),
		file:     file,
		line:     line,
		funcName: funcName,
		errMsg:   errMsg,
	}
}

func (d *UError) GetCode() int32 {
	return d.code
}

func (d *UError) GetErrMsg() string {
	return d.errMsg
}

func (d *UError) Error() string {
	return fmt.Sprintf("%s:%d %s code: %d, errmsg: %s", d.file, d.line, d.funcName, d.code, d.errMsg)
}
