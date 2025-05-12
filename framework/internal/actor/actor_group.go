package actor

import (
	"reflect"
	"strings"
	"sync"
	"universal/framework/domain"
	"universal/framework/library/uerror"

	"github.com/golang/protobuf/proto"
)

var (
	headType  = reflect.TypeOf((*domain.IHead)(nil)).Elem()
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	protoType = reflect.TypeOf((*proto.Message)(nil)).Elem()
)

type FuncInfo struct {
	reflect.Method
	hasHead    bool
	hasError   bool
	isVariadic bool
	isProto    bool
}

type ActorGroup struct {
	name   string
	mutex  sync.RWMutex
	actors map[uint64]domain.IActor
	funcs  map[string]*FuncInfo
}

func (d *ActorGroup) GetActor(id uint64) domain.IActor {
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

func (d *ActorGroup) AddActor(id uint64, act domain.IActor) {
	act.ParseFunc(d.funcs)
	d.mutex.Lock()
	d.actors[id] = act
	d.mutex.Unlock()
}

func (d *ActorGroup) GetActorName() string {
	return d.name
}

func (d *ActorGroup) Register(ac domain.IActor) {
	rtype := reflect.TypeOf(ac)
	d.name = parseName(rtype)
	d.actors = make(map[uint64]domain.IActor)
}

func (d *ActorGroup) ParseFunc(rr interface{}) {
	switch vv := rr.(type) {
	case map[string]*FuncInfo:
		d.funcs = vv
	case reflect.Type:
		d.funcs = parseFuncs(vv)
	default:
		panic("注册参数错误，必须是方法列表或reflect.Type")
	}
}

func (d *ActorGroup) Send(h domain.IHead, args ...interface{}) error {
	if _, ok := d.funcs[h.GetFuncName()]; !ok {
		return uerror.New(1, -1, "%s.%s未实现", h.GetActorName(), h.GetFuncName())
	}
	if act := d.GetActor(h.GetRouteId()); act != nil {
		return act.Send(h, args...)
	}
	return uerror.New(1, -1, "Actor不存在: %d", h.GetRouteId())
}

func (d *ActorGroup) SendRpc(h domain.IHead, buf []byte) error {
	if _, ok := d.funcs[h.GetFuncName()]; !ok {
		return uerror.New(1, -1, "%s.%s未实现", h.GetActorName(), h.GetFuncName())
	}
	if act := d.GetActor(h.GetRouteId()); act != nil {
		return act.SendRpc(h, buf)
	}
	return uerror.New(1, -1, "Actor不存在: %d", h.GetRouteId())
}

func parseName(rr reflect.Type) string {
	name := rr.String()
	if index := strings.Index(name, "."); index > -1 {
		name = name[index+1:]
	}
	return name
}

func parseFuncs(vv reflect.Type) map[string]*FuncInfo {
	funcs := make(map[string]*FuncInfo)
	for i := 0; i < vv.NumMethod(); i++ {
		m := vv.Method(i)
		hasError := false
		if m.Type.NumOut() > 0 && m.Type.Out(0).Implements(errorType) {
			hasError = true
		}
		hasHead := false
		pos := 1
		if m.Type.NumIn() > 1 && m.Type.In(1).Implements(headType) {
			hasHead = true
			pos++
		}
		isProto := true
		for i := pos; i < m.Type.NumIn(); i++ {
			if m.Type.In(i).Implements(protoType) {
				continue
			}
			isProto = false
			break
		}
		funcs[m.Name] = &FuncInfo{
			Method:     m,
			isVariadic: m.Type.IsVariadic(),
			hasHead:    hasHead,
			hasError:   hasError,
			isProto:    isProto,
		}
	}
	return funcs
}
