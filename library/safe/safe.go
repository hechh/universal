package safe

import "runtime/debug"

var catch func(string, ...interface{})

func Set(cb func(string, ...interface{})) {
	catch = cb
}

func Recover(f func()) {
	defer func() {
		if err := recover(); err != nil {
			catch("%v stack: %v", err, string(debug.Stack()))
		}
	}()
	f()
}

func Go(f func()) {
	go func() {
		Recover(f)
	}()
}
