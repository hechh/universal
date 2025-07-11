package actor

import (
	"reflect"
	"sync"
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/framework/internal/funcs"
	"universal/library/mlog"
	"universal/library/uerror"
)

type ActorMgr struct {
	id     uint64
	name   string
	mutex  sync.RWMutex
	actors map[uint64]domain.IActor
	funcs  map[string]*funcs.Method
}

func (d *ActorMgr) GetCount() int {
	return len(d.actors)
}

func (d *ActorMgr) GetActor(id uint64) domain.IActor {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.actors[id]
}

func (d *ActorMgr) DelActor(id uint64) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.actors, id)
}

func (d *ActorMgr) AddActor(act domain.IActor) {
	act.ParseFunc(d.funcs)
	id := act.GetId()
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.actors[id] = act
}

func (d *ActorMgr) GetIdPointer() *uint64 {
	return &d.id
}

func (d *ActorMgr) GetId() uint64 {
	return atomic.LoadUint64(&d.id)
}

func (d *ActorMgr) SetId(id uint64) {
	atomic.StoreUint64(&d.id, id)
}

func (d *ActorMgr) Start() {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	for _, act := range d.actors {
		act.Start()
	}
}

func (d *ActorMgr) Stop() {
	atomic.StoreUint64(&d.id, 0)
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	for _, act := range d.actors {
		act.Stop()
	}
}

func (a *ActorMgr) GetActorType() uint32 {
	return 0
}

func (d *ActorMgr) GetActorName() string {
	return d.name
}

func (d *ActorMgr) Register(ac domain.IActor, _ ...int) {
	rtype := reflect.TypeOf(ac)
	d.name = parseName(rtype)
	d.actors = make(map[uint64]domain.IActor)
}

func (d *ActorMgr) ParseFunc(rr interface{}) {
	switch vv := rr.(type) {
	case map[string]*funcs.Method:
		d.funcs = vv
	case reflect.Type:
		d.funcs = make(map[string]*funcs.Method)
		for i := 0; i < vv.NumMethod(); i++ {
			m := vv.Method(i)
			if ff := funcs.NewMethod(d, m); ff != nil {
				d.funcs[m.Name] = ff
			}
		}
	default:
		panic("注册参数错误，必须是方法列表或reflect.Type")
	}
}

func (d *ActorMgr) SendMsg(h *pb.Head, args ...interface{}) error {
	if _, ok := d.funcs[h.FuncName]; !ok {
		return uerror.New(1, -1, "%v", h)
	}
	switch h.SendType {
	case pb.SendType_POINT:
		if act := d.GetActor(h.ActorId); act != nil {
			return act.SendMsg(h, args...)
		} else {
			return uerror.New(1, -1, "Actor不存在: %v", h)
		}
	case pb.SendType_BROADCAST:
		d.mutex.RLock()
		defer d.mutex.RUnlock()
		for _, act := range d.actors {
			if err := act.SendMsg(h, args...); err != nil {
				mlog.Errorf("head:%v, error:%v", h, err)
			}
		}
	default:
		return uerror.New(1, -1, "%v", h)
	}
	return nil
}

func (d *ActorMgr) Send(h *pb.Head, buf []byte) error {
	if _, ok := d.funcs[h.FuncName]; !ok {
		return uerror.New(1, -1, "%v", h)
	}
	switch h.SendType {
	case pb.SendType_POINT:
		if act := d.GetActor(h.ActorId); act != nil {
			return act.Send(h, buf)
		} else {
			return uerror.New(1, -1, "Actor不存在: %v", h)
		}
	case pb.SendType_BROADCAST:
		d.mutex.RLock()
		defer d.mutex.RUnlock()
		for _, act := range d.actors {
			if err := act.Send(h, buf); err != nil {
				mlog.Errorf("head:%v, error:%v", h, err)
			}
		}
	default:
		return uerror.New(1, -1, "%v", h)
	}
	return nil
}

func (d *ActorMgr) RegisterTimer(id *uint64, h *pb.Head, ttl time.Duration, times int32) error {
	return t.Register(id, func() {
		if err := d.SendMsg(h); err != nil {
			mlog.Errorf("Actor定时器转发失败: %v", err)
		}
	}, ttl, times)
}
