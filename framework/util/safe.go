package util

func SafeRecover(cb func(error), f func()) {
	defer func() {
		if err := recover(); err != nil {
			cb(err.(error))
		}
	}()
	f()
}

func SafeGo(cb func(err error), f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				cb(err.(error))
				//plog.Fatal("error: %v, stack: %s", err, string(debug.Stack()))
			}
		}()
		f()
	}()
}
