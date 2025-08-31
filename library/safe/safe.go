package safe

import (
	"runtime/debug"
	"universal/library/mlog"
)

func Go(f func()) {
	go func() {
		Recover(f)
	}()
}

func Recover(f func()) {
	defer func() {
		if err := recover(); err != nil {
			mlog.Fatalf("%v stack: %v", err, string(debug.Stack()))
		}
	}()
	f()
}
