package util

func SafeRecover(cb func(interface{}), f func()) {
	defer func() {
		if err := recover(); err != nil {
			if cb != nil {
				cb(err)
			}
		}
	}()
	f()
}

func SafeGo(cb func(interface{}), f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if cb != nil {
					cb(err)
				}
				//plog.Fatal("error: %v, stack: %s", err, string(debug.Stack()))
			}
		}()
		f()
	}()
}
