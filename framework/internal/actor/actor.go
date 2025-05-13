package actor

import (
	"reflect"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/framework/library/async"
	"universal/framework/library/uerror"
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
		d.funcs = parseFuncs(vv)
	default:
		panic("注册参数错误，必须是方法列表或reflect.Type")
	}
}

func (d *Actor) Send(h *pb.Head, args ...interface{}) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "%s.%s未实现", h.ActorName, h.FuncName)
	}
	if mm.Type.IsVariadic() {
		d.Push(mm.handleVariadic(d.rval, h, args...))
	} else {
		d.Push(mm.handle(d.rval, h, args...))
	}
	return nil
}

func (d *Actor) SendRpc(h *pb.Head, buf []byte) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "%s.%s未实现", h.ActorName, h.FuncName)
	}
	// 发送事件
	if mm.isProto {
		d.Push(mm.handleRpcProto(d.rval, h, buf))
	} else {
		d.Push(mm.handleRpcGob(d.rval, h, buf))
	}
	return nil
}
