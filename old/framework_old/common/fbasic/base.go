package fbasic

import (
	"hash/crc32"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"universal/framework/common/plog"
	"unsafe"
)

func SafeRecover(f func()) {
	defer func() {
		if err := recover(); err != nil {
			plog.ErrorSkip(1, "panic: %v, stack: %s", err, string(debug.Stack()))
		}
	}()

	f()
}

func SafeGo(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				plog.ErrorSkip(1, "panic: %v, stack: %s", err, string(debug.Stack()))
			}
		}()
		f()
	}()
}

func StrToBytes(str string) []byte {
	if len(str) == 0 {
		return nil
	}
	s := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	b := &reflect.SliceHeader{Data: s.Data, Len: s.Len, Cap: s.Len}
	return *(*[]byte)(unsafe.Pointer(b))
}

func BytesToStr(bts []byte) string {
	if len(bts) == 0 {
		return ""
	}
	b := *(*reflect.SliceHeader)(unsafe.Pointer(&bts))
	s := &reflect.StringHeader{Data: b.Data, Len: b.Len}
	return *(*string)(unsafe.Pointer(s))
}

func GetCrc32(str string) uint32 {
	return crc32.ChecksumIEEE(StrToBytes(str))
}

func GetFuncName(fun interface{}) string {
	var name string
	switch vv := fun.(type) {
	case reflect.Value:
		name = runtime.FuncForPC(vv.Pointer()).Name()
	case uintptr:
		name = runtime.FuncForPC(vv).Name()
	default:
		name = runtime.FuncForPC(reflect.ValueOf(vv).Pointer()).Name()
	}
	return strings.Split(name, ".")[1]
}
