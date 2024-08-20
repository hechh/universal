package test

import (
	"fmt"
	"testing"
	"time"
	"universal/framework/basic/timer"
)

func TestTimer(t *testing.T) {
	tt := timer.NewTimer()
	times := 1000000
	for i := 1; i <= times; i++ {
		tt.Insert(func() {
			if i == times {
				fmt.Println("-------> ", i)
			}
		}, (1+time.Duration(i)%2)*time.Second, false)
	}
	time.Sleep(8 * time.Second)
	tt.Stop()
}

func TestTimer01(t *testing.T) {
	tt := timer.NewTimer()
	times := 3
	for i := 1; i <= times; i++ {
		tt.Insert(func() {
			fmt.Println("-------> ", i, time.Now().Unix())
		}, (time.Duration(2*i))*time.Second, false)
	}
	time.Sleep(8 * time.Second)
	tt.Stop()
}
