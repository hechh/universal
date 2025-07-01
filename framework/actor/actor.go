package actor

import (
	"reflect"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/library/async"
	"universal/library/uerror"
)

type Actor struct {
	*async.Async
	name  string
	rval  reflect.Value
	funcs map[string]*FuncInfo
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
	case map[string]*FuncInfo:
		d.funcs = vv
	case reflect.Type:
		d.funcs = make(map[string]*FuncInfo)
		for i := 0; i < vv.NumMethod(); i++ {
			m := vv.Method(i)
			d.funcs[m.Name] = parseFuncInfo(m)
		}
	default:
		panic("注册参数错误，必须是方法列表或reflect.Type")
	}
}

func (d *Actor) SendMsg(h *pb.Head, args ...interface{}) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.N(1, -1, "%s.%s未实现, head:%v", h.ActorName, h.FuncName, h)
	}
	switch mm.flag {
	case CMD_HANDLER:
		d.Push(mm.localCmd(d.rval, h, args...))
	default:
		d.Push(mm.localProto(d.rval, h, args...))
	}
	return nil
}
