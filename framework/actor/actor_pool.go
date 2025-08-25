package actor

import (
	"reflect"
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/define"
	"universal/framework/internal/handler"
	"universal/library/async"
	"universal/library/mlog"
	"universal/library/templ"
	"universal/library/uerror"
)

type ActorPool struct {
	size  int
	pool  []*async.Async
	id    uint64
	name  string
	rval  define.IActor
	funcs map[string]define.IHandler
}

func (d *ActorPool) GetIdPointer() *uint64 {
	return &d.id
}

func (d *ActorPool) GetId() uint64 {
	return atomic.LoadUint64(&d.id)
}

func (d *ActorPool) SetId(id uint64) {
	atomic.StoreUint64(&d.id, id)
}

func (d *ActorPool) Start() {
	for _, act := range d.pool {
		act.Start()
	}
}

func (d *ActorPool) Stop() {
	atomic.StoreUint64(&d.id, 0)
	for _, act := range d.pool {
		act.Stop()
	}
}

func (d *ActorPool) GetActorName() string {
	return d.name
}

func (d *ActorPool) Register(ac define.IActor, sizes ...int) {
	d.size = templ.Index[int](sizes, 0, 10)
	d.pool = make([]*async.Async, d.size)
	for i := 0; i < d.size; i++ {
		d.pool[i] = async.New()
	}
	d.name = parseName(reflect.TypeOf(ac))
	d.rval = ac
	d.funcs = handler.GetActor(self.Type, d.name)
}

func (d *ActorPool) SendMsg(h *pb.Head, args ...interface{}) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "接口%s.%s未注册", h.ActorName, h.FuncName)
	}
	d.pool[h.ActorId%uint64(d.size)].Push(mm.Call(sendrsp, d.rval, h, args...))
	return nil
}

func (d *ActorPool) Send(h *pb.Head, buf []byte) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "接口%s.%s未注册", h.ActorName, h.FuncName)
	}
	d.pool[h.ActorId%uint64(d.size)].Push(mm.Rpc(sendrsp, d.rval, h, buf))
	return nil
}

func (d *ActorPool) RegisterTimer(id *uint64, h *pb.Head, ttl time.Duration, times int32) error {
	return t.Register(id, func() {
		if err := d.SendMsg(h); err != nil {
			mlog.Errorf("Actor定时器转发失败: %v", err)
		}
	}, ttl, times)
}
