package actor

import (
	"reflect"
	"sync"
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/define"
	"universal/framework/handler"
	"universal/library/mlog"
	"universal/library/uerror"
)

type ActorMgr struct {
	id     uint64
	name   string
	mutex  sync.RWMutex
	actors map[uint64]define.IActor
	funcs  map[string]define.IHandler
}

func (d *ActorMgr) GetCount() int {
	return len(d.actors)
}

func (d *ActorMgr) GetActor(id uint64) define.IActor {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.actors[id]
}

func (d *ActorMgr) DelActor(id uint64) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.actors, id)
}

func (d *ActorMgr) AddActor(act define.IActor) {
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

func (d *ActorMgr) GetActorName() string {
	return d.name
}

func (d *ActorMgr) Register(ac define.IActor, _ ...int) {
	d.actors = make(map[uint64]define.IActor)
	d.name = parseName(reflect.TypeOf(ac))
	d.funcs = handler.GetActor(self.Type, d.name)
}

func (d *ActorMgr) SendMsg(h *pb.Head, args ...interface{}) error {
	if _, ok := d.funcs[h.FuncName]; !ok {
		return uerror.New(1, -1, "接口%s.%s未注册", h.ActorName, h.FuncName)
	}
	switch h.SendType {
	case pb.SendType_POINT:
		if act := d.GetActor(h.ActorId); act != nil {
			return act.SendMsg(h, args...)
		} else {
			return uerror.New(1, -1, "ActorId(%d)不存在", h.ActorId)
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
		return uerror.New(1, -1, "发送类型不支持%v", h.SendType)
	}
	return nil
}

func (d *ActorMgr) Send(h *pb.Head, buf []byte) error {
	if _, ok := d.funcs[h.FuncName]; !ok {
		return uerror.New(1, -1, "接口%s.%s未注册", h.ActorName, h.FuncName)
	}
	switch h.SendType {
	case pb.SendType_POINT:
		if act := d.GetActor(h.ActorId); act != nil {
			return act.Send(h, buf)
		} else {
			return uerror.New(1, -1, "ActorId(%d)不存在", h.ActorId)
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
		return uerror.New(1, -1, "发送类型不支持%v", h.SendType)
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
