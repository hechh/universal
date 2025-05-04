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
)

type MethodInfo struct {
	method    reflect.Method
	isHandler bool
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

func (d *Actor) Send(header define.IHeader, args ...interface{}) error {
	m, ok := d.methods[header.GetFuncName()]
	if !ok {
		return uerror.New(1, -1, "方法不存在: %s.%s", header.GetActorName(), header.GetFuncName())
	}
	d.Push(func() {
		ins := make([]reflect.Value, m.method.Type.NumIn())
		ins[0] = d.rValue
		for i := 1; i < m.method.Type.NumIn(); i++ {
			ins[i] = reflect.ValueOf(args[i-1])
		}

		// 无返回值
		if !m.isHandler {
			m.method.Func.Call(ins)
		} else {
			result := m.method.Func.Call(ins)
			switch val := result[0].Interface().(type) {
			case error:
				mlog.Error("接口调用失败，head:%v, error:%v", header, val)
			default:
				mlog.Error("接口调用成功")
			}
		}
	})
	return nil
}

func (d *Actor) SendFrom(head define.IHeader, buf []byte) error {
	m, ok := d.methods[head.GetFuncName()]
	if !ok {
		return uerror.New(1, -1, "方法不存在: %s.%s", head.GetActorName(), head.GetFuncName())
	}
	if !m.isHandler {
		d.Push(otherFunc(m.method, d.rValue, head, buf))
	} else {
		d.Push(handleFunc(m.method, d.rValue, head, buf))
	}
	return nil
}

func otherFunc(m reflect.Method, rValue reflect.Value, head define.IHeader, buf []byte) func() {
	return func() {
		ins, err := encode.Decode(buf, m)
		if err != nil {
			mlog.Error("参数解析输错: head:%v, error:%v", head, err)
			return
		}
		ins[0] = rValue

		m.Func.Call(ins)
	}
}

func handleFunc(m reflect.Method, rValue reflect.Value, head define.IHeader, buf []byte) func() {
	return func() {
		ins := make([]reflect.Value, m.Type.NumIn())
		ins[0] = rValue
		ins[1] = reflect.New(m.Type.In(1))
		ins[2] = reflect.New(m.Type.In(2))

		// 解析参数
		if err := proto.Unmarshal(buf, ins[1].Interface().(proto.Message)); err != nil {
			mlog.Error("head:%v, error: %v", head, err)
		}

		// 调用接口
		result := m.Func.Call(ins)
		switch val := result[0].Interface().(type) {
		case error:
			mlog.Error("接口调用失败，head:%v, error:%v", head, val)
		default:
			mlog.Error("接口调用成功")
		}
	}
}

func parseMethod(m reflect.Type) (ret map[string]*MethodInfo) {
	ret = make(map[string]*MethodInfo)
	for i := 0; i < m.NumMethod(); i++ {
		mm := m.Method(i)
		if mm.Type.NumIn() == 3 && mm.Type.NumOut() == 1 &&
			mm.Type.In(1).Implements(messageType) &&
			mm.Type.In(2).Implements(messageType) &&
			mm.Type.Out(0).Implements(errorType) {
			ret[mm.Name] = &MethodInfo{method: mm, isHandler: true}
		} else {
			ret[mm.Name] = &MethodInfo{method: mm, isHandler: false}
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
