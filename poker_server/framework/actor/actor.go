package actor

import (
	"poker_server/common/pb"
	"poker_server/framework/domain"
	"poker_server/framework/internal/method"
	"poker_server/library/async"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"reflect"
	"time"
)

type Actor struct {
	*async.Async
	name  string
	rval  reflect.Value
	funcs map[string]*method.Method
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
	case map[string]*method.Method:
		d.funcs = vv
	case reflect.Type:
		d.funcs = make(map[string]*method.Method)
		for i := 0; i < vv.NumMethod(); i++ {
			m := vv.Method(i)
			if ff := method.NewMethod(m); ff != nil {
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
		return uerror.New(1, -1, "%s.%s未实现, head:%v", h.ActorName, h.FuncName, h)
	}
	d.Push(mm.Call(d.rval, h, args...))
	return nil
}

func (d *Actor) Send(h *pb.Head, buf []byte) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "%s.%s未实现, head:%v", h.ActorName, h.FuncName, h)
	}
	d.Push(mm.Rpc(d.rval, h, buf))
	return nil
}

func (d *Actor) RegisterTimer(h *pb.Head, ttl time.Duration, times int32) error {
	return t.Register(d.GetIdPointer(), func() {
		if err := d.SendMsg(h); err != nil {
			mlog.Errorf("定时器发送消息失败: head:%v, error:%v", h, err)
		}
	}, ttl, times)
}
