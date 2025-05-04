package actor

import (
	"reflect"
	"strings"
	"sync"
	"universal/framework/define"
	"universal/library/baselib/uerror"
)

type ActorGroup struct {
	name    string
	methods map[string]reflect.Method
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

func (d *ActorGroup) Init(ac define.IActor) {
	rType := reflect.TypeOf(ac)
	name := rType.String()
	if index := strings.Index(name, "."); index != -1 {
		name = name[index+1:]
	}

	// 初始化
	d.name = name
	d.actors = make(map[uint64]define.IActor)
}

func (d *ActorGroup) Register(ac define.IActor, _ interface{}) error {
	rType := reflect.TypeOf(ac)
	name := rType.String()
	if index := strings.Index(name, "."); index != -1 {
		name = name[index+1:]
	}

	d.name = name
	d.methods = make(map[string]reflect.Method)
	d.actors = make(map[uint64]define.IActor)

	for i := 0; i < rType.NumMethod(); i++ {
		methond := rType.Method(i)
		d.methods[methond.Name] = methond
	}
	return nil
}

func (d *ActorGroup) Send(header define.IHeader, args ...interface{}) error {
	if _, ok := d.methods[header.GetFuncName()]; !ok {
		return uerror.New(1, -1, "方法不存在: %s.%s", header.GetActorName(), header.GetFuncName())
	}

	// 获取Actor
	if act := d.GetActor(header.GetRouteId()); act != nil {
		return act.Send(header, args...)
	}
	return uerror.New(1, -1, "Actor不存在: %d", header.GetRouteId())
}
