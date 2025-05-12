package async

import (
	"runtime/debug"
)

func SafeRecover(cb func(string, ...interface{}), f func()) {
	defer func() {
		if err := recover(); err != nil {
			if cb == nil {
				return
			}
			cb("%v stack: %v", err, string(debug.Stack()))
		}
	}()
	f()
}

func SafeGo(cb func(string, ...interface{}), f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if cb == nil {
					return
				}
				cb("%v stack: %v", err, string(debug.Stack()))
			}
		}()
		f()
	}()
}
