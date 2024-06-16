package test

import (
	"fmt"
	"testing"
	"universal/framework/common/fbasic"
	"universal/framework/common/timer"
)

func Print() {
	fmt.Println("--------")
}

func TestNode(t *testing.T) {

	list := new(timer.NodeList)
	list.Insert(timer.NewTask(Print, 5))
	list.Insert(timer.NewTask(Print, 4))
	list.Insert(timer.NewTask(Print, 5))
	list.Insert(timer.NewTask(Print, 4))
	list.Insert(timer.NewTask(Print, 3))
	list.Insert(timer.NewTask(Print, 3))
	list.Insert(timer.NewTask(Print, 3))
	list.Insert(timer.NewTask(Print, 3))
	list.Insert(timer.NewTask(Print, 2))
	list.Insert(timer.NewTask(Print, 1))
	list.Insert(timer.NewTask(Print, 1))
	list.Insert(timer.NewTask(Print, 1230))
	list.Insert(timer.NewTask(Print, 9))
	list.Insert(timer.NewTask(Print, 8))
	list.Insert(timer.NewTask(Print, 6))

	list.Print()
	t.Log("--------expire----------")
	rets := list.Pop(fbasic.GetNow(), 3)
	for _, item := range rets {
		item.Print()
	}
}
