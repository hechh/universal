package safe

var panicNotify func(interface{})

func SetPanicNotify(f func(interface{})) {
	panicNotify = f
}

func SafeRecover(f func()) {
	defer func() {
		if err := recover(); err != nil {
			if panicNotify != nil {
				panicNotify(err)
			}
		}
	}()
	f()
}

func SafeGo(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if panicNotify != nil {
					panicNotify(err)
				}
			}
		}()

		f()
	}()
}
