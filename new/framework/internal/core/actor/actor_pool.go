package actor

import (
	"poker_server/common/pb"
	"poker_server/framework/domain"
	"poker_server/framework/library/async"
	"poker_server/framework/library/uerror"
	"reflect"
)

type ActorPool struct {
	id       uint64
	pool     []*async.Async
	poolSize int
	rval     reflect.Value
	name     string
	funcs    map[string]*FuncInfo
}

func (d *ActorPool) GetId() uint64 {
	return d.id
}

func (d *ActorPool) SetId(id uint64) {
	d.id = id
}

func (d *ActorPool) Start() {
	for _, async := range d.pool {
		async.Start()
	}
}

func (d *ActorPool) Stop() {
	for _, async := range d.pool {
		async.Stop()
	}
}

func (d *ActorPool) GetActorName() string {
	return d.name
}

func (d *ActorPool) Register(ac domain.IActor, sizes ...int) {
	if len(sizes) <= 0 {
		panic("ActorPool注册参数错误，必须指定协程池大小")
	}
	d.poolSize = sizes[0]
	d.pool = make([]*async.Async, d.poolSize)
	for i := 0; i < d.poolSize; i++ {
		d.pool[i] = async.NewAsync()
	}
	d.rval = reflect.ValueOf(ac)
	d.name = parseName(d.rval.Elem().Type())
}

func (d *ActorPool) ParseFunc(tt interface{}) {
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

func (d *ActorPool) SendMsg(h *pb.Head, args ...interface{}) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "%s.%s未实现", h.ActorName, h.FuncName)
	}
	if mm.Type.IsVariadic() {
		d.pool[h.Id%uint64(d.poolSize)].Push(mm.handleVariadic(d.rval, h, args...))
	} else {
		d.pool[h.Id%uint64(d.poolSize)].Push(mm.handle(d.rval, h, args...))
	}
	return nil
}

func (d *ActorPool) Send(h *pb.Head, buf []byte) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "%s.%s未实现", h.ActorName, h.FuncName)
	}
	if mm.isNotify {
		d.pool[h.Id%uint64(d.poolSize)].Push(mm.handleRpcNotify(d.rval, h, buf))
	} else if mm.isCmd {
		d.pool[h.Id%uint64(d.poolSize)].Push(mm.handleRpcCmd(d.rval, h, buf))
	} else {
		d.pool[h.Id%uint64(d.poolSize)].Push(mm.handleRpcGob(d.rval, h, buf))
	}
	return nil
}
