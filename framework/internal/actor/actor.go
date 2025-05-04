package actor

import (
	"reflect"
	"strings"
	"universal/framework/define"
	"universal/library/async"
	"universal/library/baselib/uerror"
)

type Actor struct {
	*async.Async
	name    string
	rValue  reflect.Value
	methods map[string]reflect.Method
}

func (d *Actor) GetName() string {
	return d.name
}

func (d *Actor) Register(ac define.IActor, tt interface{}) error {
	// 初始化
	if ac != nil {
		vv := reflect.ValueOf(ac)
		name := vv.Elem().Type().String()
		if index := strings.Index(name, "."); index != -1 {
			name = name[index+1:]
		}
		d.Async = async.NewAsync()
		d.name = name
		d.rValue = vv
	}

	// 注册方法
	if tt != nil {
		switch vv := tt.(type) {
		case map[string]reflect.Method:
			d.methods = vv
		case reflect.Type:
			d.methods = make(map[string]reflect.Method)
			for i := 0; i < vv.NumMethod(); i++ {
				methond := vv.Method(i)
				d.methods[methond.Name] = methond
			}
		default:
			return uerror.New(1, -1, "传入必须是Actor的 reflect.Type或方法列表")
		}
	}
	return nil
}

func (d *Actor) Send(header define.IHeader, args ...interface{}) error {
	m, ok := d.methods[header.GetFuncName()]
	if !ok {
		return uerror.New(1, -1, "方法不存在: %s.%s", header.GetActorName(), header.GetFuncName())
	}

	d.Push(func() {
		ins := make([]reflect.Value, m.Type.NumIn())
		ins[0] = d.rValue

		for i := 1; i < m.Type.NumIn(); i++ {
			ins[i] = reflect.ValueOf(args[i-1])
		}

		// 无返回值
		m.Func.Call(ins)
	})
	return nil
}
