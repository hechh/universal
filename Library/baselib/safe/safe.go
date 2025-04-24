package safe

func SafeRecover(cb func(interface{}), f func()) {
	defer func() {
		if err := recover(); err != nil {
			cb(err)
		}
	}()
	f()
}

func SafeGo(cb func(interface{}), f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				cb(err)
			}
		}()
		f()
	}()
}
