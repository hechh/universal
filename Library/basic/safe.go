package basic

import (
	"os"
	"os/signal"
	"syscall"
)

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
				//ulog.Fatal("error: %v, stack: %s", err, string(debug.Stack()))
			}
		}()
		f()
	}()
}

// 阻塞接受信号
func SignalHandle(f func(os.Signal), sigs ...os.Signal) {
	ch := make(chan os.Signal, 0)
	args := append([]os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL}, sigs...)
	signal.Notify(ch, args...)
	for item := range ch {
		f(item)
		os.Exit(0)
	}
}
