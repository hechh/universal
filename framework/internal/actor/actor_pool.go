package actor

import (
	"reflect"
	"universal/framework/domain"
	"universal/framework/library/async"
)

type ActorPool struct {
	pool  []*async.Async
	rval  reflect.Value
	name  string
	funcs map[string]*FuncInfo
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

func (d *ActorPool) Register(ac domain.IActor) {
	d.rval = reflect.ValueOf(ac)
	d.name = parseName(d.rval.Elem().Type())
}

func (d *ActorPool) ParseFunc(tt interface{}) {
	switch vv := tt.(type) {
	case map[string]*FuncInfo:
		d.funcs = vv
	case reflect.Type:
		d.funcs = parseFuncs(vv)
	default:
		panic("注册参数错误，必须是方法列表或reflect.Type")
	}
}
