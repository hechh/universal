package base

import "runtime/debug"

func SafeRecover(notify func(string, ...interface{}), f func()) {
	defer func() {
		if err := recover(); err != nil {
			notify("error: %v, stack: %s", err, string(debug.Stack()))
		}
	}()
	f()
}

func SafeGo(notify func(string, ...interface{}), f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				notify("error: %v, stack: %s", err, string(debug.Stack()))
			}
		}()
		f()
	}()
}
