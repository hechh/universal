package actor

import (
	"reflect"
	"universal/framework/define"
	"universal/library/async"
)

type Actor struct {
	*async.Async
	rValue reflect.Value
}

func (d *Actor) RegisterValue(ac interface{}) {
	d.rValue = reflect.ValueOf(ac)
}

func (d *Actor) GetValue() reflect.Value {
	return d.rValue
}

func (d *Actor) Send(header define.IHeader, fname string, args ...interface{}) {
	d.Push(func() {
		m, _ := GetMethod(fname)
		ins := make([]reflect.Value, m.Type.NumIn())
		ins[0] = d.rValue
		for i := 1; i < m.Type.NumIn(); i++ {
			ins[i] = reflect.ValueOf(args[i-1])
		}

		// 无返回值
		m.Func.Call(ins)
	})
}
