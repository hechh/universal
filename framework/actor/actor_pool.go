package actor

import (
	"reflect"
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/define"
	"universal/framework/internal/funcs"
	"universal/library/async"
	"universal/library/mlog"
	"universal/library/uerror"
	"universal/library/util"
)

type ActorPool struct {
	size  int
	pool  []*async.Async
	id    uint64
	name  string
	rval  reflect.Value
	funcs map[string]define.IFuncs
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

func (a *ActorPool) GetActorType() uint32 {
	return 0
}

func (d *ActorPool) GetActorName() string {
	return d.name
}

func (d *ActorPool) Register(ac define.IActor, sizes ...int) {
	d.size = util.Index[int](sizes, 0, 10)
	d.pool = make([]*async.Async, d.size)
	for i := 0; i < d.size; i++ {
		d.pool[i] = async.NewAsync()
	}
	d.name = parseName(d.rval.Elem().Type())
	d.rval = reflect.ValueOf(ac)
}

func (d *ActorPool) ParseFunc(tt interface{}) {
	switch vv := tt.(type) {
	case map[string]define.IFuncs:
		d.funcs = vv
	case reflect.Type:
		d.funcs = make(map[string]define.IFuncs)
		for i := 0; i < vv.NumMethod(); i++ {
			m := vv.Method(i)
			if ff := funcs.NewMethod(d, m); ff != nil {
				d.funcs[m.Name] = ff
				apis[GetCrc32(d.name+"."+m.Name)] = ff
			}
		}
	default:
		panic("注册参数错误，必须是方法列表或reflect.Type")
	}
}

func (d *ActorPool) SendMsg(h *pb.Head, args ...interface{}) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "请求不存在(%s.%s)", h.ActorName, h.FuncName)
	}
	d.pool[h.ActorId%uint64(d.size)].Push(mm.Call(d.rval, h, args...))
	return nil
}

func (d *ActorPool) Send(h *pb.Head, buf []byte) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "请求不存在(%s.%s)", h.ActorName, h.FuncName)
	}
	d.pool[h.ActorId%uint64(d.size)].Push(mm.Rpc(d.rval, h, buf))
	return nil
}

func (d *ActorPool) RegisterTimer(id *uint64, h *pb.Head, ttl time.Duration, times int32) error {
	return t.Register(id, func() {
		if err := d.SendMsg(h); err != nil {
			mlog.Errorf("Actor定时器转发失败: %v", err)
		}
	}, ttl, times)
}
