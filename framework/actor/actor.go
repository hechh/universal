package actor

import (
	"reflect"
	"time"
	"universal/common/pb"
	"universal/framework/define"
	"universal/framework/handler"
	"universal/library/async"
	"universal/library/mlog"
	"universal/library/uerror"
)

type Actor struct {
	*async.Async
	name  string
	self  define.IActor
	funcs map[string]define.IHandler
}

func (a *Actor) GetActorName() string {
	return a.name
}

func (a *Actor) Register(ac define.IActor, _ ...int) {
	a.Async = async.New()
	a.name = parseName(reflect.TypeOf(ac))
	a.self = ac
	a.funcs = handler.GetActor(self.Type, a.name)
}

func (d *Actor) SendMsg(h *pb.Head, args ...interface{}) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "接口%s.%s未注册", h.ActorName, h.FuncName)
	}
	d.Push(mm.Call(sendrsp, d.self, h, args...))
	return nil
}

func (d *Actor) Send(h *pb.Head, buf []byte) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "接口%s.%s未注册", h.ActorName, h.FuncName)
	}
	d.Push(mm.Rpc(sendrsp, d.self, h, buf))
	return nil
}

func (d *Actor) RegisterTimer(id *uint64, h *pb.Head, ttl time.Duration, times int32) error {
	return t.Register(id, func() {
		if err := d.SendMsg(h); err != nil {
			mlog.Errorf("Actor定时器转发失败: %v", err)
		}
	}, ttl, times)
}
