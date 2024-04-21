package basic

import (
	"log"
	"runtime/debug"
)

func SafeRecover(f func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("stack: \n%s", string(debug.Stack()))
		}
	}()

	f()
}

func SafeGo(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("stack: \n%s", string(debug.Stack()))
			}
		}()
		f()
	}()
}
