package actor

import (
	"reflect"
	"strings"
	"universal/framework/define"
	"universal/library/async"
	"universal/library/baselib/uerror"
	"universal/library/encode"
	"universal/library/mlog"

	"github.com/golang/protobuf/proto"
)

var (
	messageType = reflect.TypeOf((*proto.Message)(nil)).Elem()
	errorType   = reflect.TypeOf((*error)(nil)).Elem()
	ctxType     = reflect.TypeOf((*define.IContext)(nil)).Elem()
)

type MethodInfo struct {
	method     reflect.Method
	hasContext bool
	isHandler  bool
}

type Actor struct {
	*async.Async
	name    string
	rValue  reflect.Value
	methods map[string]*MethodInfo
}

func (d *Actor) GetName() string {
	return d.name
}

func (d *Actor) Register(ac define.IActor, tt interface{}) error {
	// 初始化
	if ac != nil {
		vv := reflect.ValueOf(ac)
		d.Async = async.NewAsync()
		d.name = parseName(vv.Elem().Type())
		d.rValue = vv
	}
	// 注册方法
	switch vv := tt.(type) {
	case map[string]*MethodInfo:
		d.methods = vv
	case reflect.Type:
		d.methods = parseMethod(vv)
	default:
		return uerror.New(1, -1, "传入必须是Actor的 reflect.Type或方法列表")
	}
	return nil
}

func (d *Actor) Send(ctx define.IContext, args ...interface{}) error {
	m, ok := d.methods[ctx.GetFuncName()]
	if !ok {
		return uerror.New(1, -1, "方法不存在: %s.%s", ctx.GetActorName(), ctx.GetFuncName())
	}
	d.Push(func() {
		if m.hasContext {
			ins := make([]reflect.Value, m.method.Type.NumIn())
			ins[0] = d.rValue
			ins[1] = reflect.ValueOf(ctx)
			for i := 2; i < m.method.Type.NumIn(); i++ {
				ins[i] = reflect.ValueOf(args[i-2])
			}
			// 无返回值
			m.method.Func.Call(ins)
		} else {
			ins := make([]reflect.Value, m.method.Type.NumIn())
			ins[0] = d.rValue
			for i := 1; i < m.method.Type.NumIn(); i++ {
				ins[i] = reflect.ValueOf(args[i-1])
			}
			// 无返回值
			m.method.Func.Call(ins)
		}
	})
	return nil
}

func (d *Actor) SendFrom(head define.IContext, buf []byte) error {
	m, ok := d.methods[head.GetFuncName()]
	if !ok {
		return uerror.New(1, -1, "方法不存在: %s.%s", head.GetActorName(), head.GetFuncName())
	}
	d.Push(func() {
		if m.hasContext {
			if m.isHandler {
				ins := make([]reflect.Value, m.method.Type.NumIn())
				ins[0] = d.rValue
				ins[1] = reflect.ValueOf(head)
				ins[2] = reflect.New(m.method.Type.In(2))
				ins[3] = reflect.New(m.method.Type.In(3))
				m.method.Func.Call(ins)
			} else {
				ins, err := encode.Decode(buf, m.method, 2)
				if err != nil {
					mlog.Error("参数解析输错: head:%v, error:%v", head, err)
					return
				}
				ins[0] = d.rValue
				ins[1] = reflect.ValueOf(head)
				m.method.Func.Call(ins)
			}
		} else {
			ins, err := encode.Decode(buf, m.method, 1)
			if err != nil {
				mlog.Error("参数解析输错: head:%v, error:%v", head, err)
				return
			}
			ins[0] = d.rValue
			m.method.Func.Call(ins)
		}
	})
	return nil
}

func parseMethod(m reflect.Type) (ret map[string]*MethodInfo) {
	ret = make(map[string]*MethodInfo)
	for i := 0; i < m.NumMethod(); i++ {
		mm := m.Method(i)
		hasContext := false
		if mm.Type.NumIn() > 1 {
			hasContext = mm.Type.In(1).Implements(ctxType)
		}
		if mm.Type.NumIn() == 4 && mm.Type.NumOut() == 0 &&
			mm.Type.In(1).Implements(ctxType) &&
			mm.Type.In(2).Implements(messageType) &&
			mm.Type.In(3).Implements(messageType) {
			ret[mm.Name] = &MethodInfo{method: mm, isHandler: true, hasContext: hasContext}
		} else {
			ret[mm.Name] = &MethodInfo{method: mm, isHandler: false, hasContext: hasContext}
		}
	}
	return
}

func parseName(rr reflect.Type) string {
	name := rr.String()
	if index := strings.Index(name, "."); index > -1 {
		name = name[index+1:]
	}
	return name
}
