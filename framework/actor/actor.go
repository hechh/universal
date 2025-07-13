package actor

import (
	"reflect"
	"time"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/framework/internal/funcs"
	"universal/library/async"
	"universal/library/mlog"
	"universal/library/uerror"
)

type Actor struct {
	*async.Async
	name  string
	rval  reflect.Value
	funcs map[string]domain.IFuncs
}

func (a *Actor) GetActorType() uint32 {
	return 0
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
	case map[string]domain.IFuncs:
		d.funcs = vv
	case reflect.Type:
		d.funcs = make(map[string]domain.IFuncs)
		for i := 0; i < vv.NumMethod(); i++ {
			m := vv.Method(i)
			ff := funcs.NewMethod(d, m)
			if ff != nil {
				d.funcs[m.Name] = ff
				apis[GetCrc32(d.name+"."+m.Name)] = ff
			}
		}
	default:
		panic("注册参数错误，必须是方法列表或reflect.Type")
	}
}

func (d *Actor) SendMsg(h *pb.Head, args ...interface{}) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "%v", h)
	}
	d.Push(mm.Call(d.rval, h, args...))
	return nil
}

func (d *Actor) Send(h *pb.Head, buf []byte) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "%v", h)
	}
	d.Push(mm.Rpc(d.rval, h, buf))
	return nil
}

func (d *Actor) RegisterTimer(id *uint64, h *pb.Head, ttl time.Duration, times int32) error {
	return t.Register(id, func() {
		if err := d.SendMsg(h); err != nil {
			mlog.Errorf("Actor定时器转发失败: %v", err)
		}
	}, ttl, times)
}
