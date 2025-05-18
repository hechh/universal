package actor

import (
	"poker_server/common/pb"
	"poker_server/framework/domain"
	"poker_server/framework/library/async"
	"poker_server/framework/library/uerror"
	"reflect"
	"strings"
)

type Actor struct {
	*async.Async
	id    uint64
	name  string
	rval  reflect.Value
	funcs map[string]*FuncInfo
}

func (a *Actor) GetId() uint64 {
	return a.id
}

func (a *Actor) SetId(id uint64) {
	a.id = id
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
			d.funcs[m.Name] = parseFunc(m)
		}
	default:
		panic("注册参数错误，必须是方法列表或reflect.Type")
	}
}

func (d *Actor) SendMsg(h *pb.Head, args ...interface{}) error {
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

func (d *Actor) Send(h *pb.Head, buf []byte) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "%s.%s未实现", h.ActorName, h.FuncName)
	}
	// 发送事件
	if mm.isNotify {
		d.Push(mm.handleRpcNotify(d.rval, h, buf))
	} else if mm.isCmd {
		d.Push(mm.handleRpcCmd(d.rval, h, buf))
	} else {
		d.Push(mm.handleRpcGob(d.rval, h, buf))
	}
	return nil
}

func parseName(rr reflect.Type) string {
	name := rr.String()
	if index := strings.Index(name, "."); index > -1 {
		name = name[index+1:]
	}
	return name
}
