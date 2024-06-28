package test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := time.NewTimer(2 * time.Second)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("----------定时器触发01-----------")
		<-timer.C
		fmt.Println("----------定时器触发02-----------")
	}()

	time.Sleep(time.Millisecond)
	timer.Stop()
	fmt.Println("定时器已取消")
	wg.Wait()
}

func TestTimer02(t *testing.T) {
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	// 创建一个定时器，在2秒后触发
	timer := time.NewTimer(0)
	t.Log("---------1----------")
	<-timer.C
	/*
		go func() {
			timer.Stop()
		}()
	*/
	t.Log("---------2----------")
	//wg.Done()
}
