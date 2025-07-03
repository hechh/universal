package actor

import (
	"reflect"
	"strings"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/framework/internal/funcs"
	"universal/library/async"
	"universal/library/uerror"
)

type Actor struct {
	*async.Async
	name  string
	rval  reflect.Value
	funcs map[string]*funcs.Method
}

func (a *Actor) GetActorName() string {
	return a.name
}

func (a *Actor) Register(ac domain.IActor, _ ...int) {
	a.Async = async.NewAsync()
	a.rval = reflect.ValueOf(ac)
	a.name = parseName(a.rval.Elem().Type())
}

func (d *Actor) ParseFunc(tt interface{}) {
	switch vv := tt.(type) {
	case map[string]*funcs.Method:
		d.funcs = vv
	case reflect.Type:
		d.funcs = make(map[string]*funcs.Method)
		for i := 0; i < vv.NumMethod(); i++ {
			m := vv.Method(i)
			if ff := funcs.NewMethod(m); ff != nil {
				d.funcs[m.Name] = ff
			}
		}
	default:
		panic("注册参数错误，必须是方法列表或reflect.Type")
	}
}

func (d *Actor) SendMsg(h *pb.Head, args ...interface{}) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.N(1, -1, "%v", h)
	}
	d.Push(mm.Call(d.rval, h, args...))
	return nil
}

func (d *Actor) Send(h *pb.Head, buf []byte) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.N(1, -1, "%v", h)
	}
	d.Push(mm.Rpc(d.rval, h, buf))
	return nil
}

func parseName(rr reflect.Type) string {
	name := rr.String()
	if index := strings.Index(name, "."); index > -1 {
		name = name[index+1:]
	}
	return name
}
