package safe

import (
	"runtime/debug"
	"universal/library/mlog"
)

/*
var catch func(string, ...interface{})

func Init(cb func(string, ...interface{})) {
	catch = cb
}
*/

func Recover(f func()) {
	defer func() {
		if err := recover(); err != nil {
			mlog.Fatalf("%v stack: %v", err, string(debug.Stack()))
		}
	}()
	f()
}

func Go(f func()) {
	go func() {
		Recover(f)
	}()
}
