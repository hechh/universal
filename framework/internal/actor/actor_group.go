package actor

import (
	"reflect"
	"sync"
	"universal/framework/define"
	"universal/library/baselib/uerror"
)

type ActorGroup struct {
	name    string
	methods map[string]*MethodInfo
	mutex   sync.RWMutex
	actors  map[uint64]define.IActor
}

func (d *ActorGroup) GetActor(id uint64) define.IActor {
	d.mutex.RLock()
	actor, ok := d.actors[id]
	d.mutex.RUnlock()
	if ok {
		return actor
	}
	return nil
}

func (d *ActorGroup) DelActor(id uint64) {
	d.mutex.Lock()
	delete(d.actors, id)
	d.mutex.Unlock()
}

func (d *ActorGroup) AddActor(id uint64, act define.IActor) {
	act.Register(nil, d.methods)

	d.mutex.Lock()
	d.actors[id] = act
	d.mutex.Unlock()
}

func (d *ActorGroup) GetName() string {
	return d.name
}

func (d *ActorGroup) Register(ac define.IActor, rr interface{}) error {
	if ac != nil {
		rType := reflect.TypeOf(ac)
		d.name = parseName(rType)
		d.actors = make(map[uint64]define.IActor)
	}

	switch vv := rr.(type) {
	case map[string]*MethodInfo:
		d.methods = vv
	case reflect.Type:
		d.methods = parseMethod(vv)
	default:
		return uerror.New(1, -1, "传入必须是Actor的 reflect.Type或方法列表")
	}
	return nil
}

func (d *ActorGroup) Send(header define.IContext, args ...interface{}) error {
	if _, ok := d.methods[header.GetFuncName()]; !ok {
		return uerror.New(1, -1, "方法不存在: %s.%s", header.GetActorName(), header.GetFuncName())
	}
	// 获取Actor
	if act := d.GetActor(header.GetRouteId()); act != nil {
		return act.Send(header, args...)
	}
	return uerror.New(1, -1, "Actor不存在: %d", header.GetRouteId())
}

func (d *ActorGroup) SendFrom(head define.IContext, buf []byte) error {
	if _, ok := d.methods[head.GetFuncName()]; !ok {
		return uerror.New(1, -1, "方法不存在: %s.%s", head.GetActorName(), head.GetFuncName())
	}

	if act := d.GetActor(head.GetRouteId()); act != nil {
		return act.SendFrom(head, buf)
	}
	return uerror.New(1, -1, "Actor不存在: %d", head.GetRouteId())
}
