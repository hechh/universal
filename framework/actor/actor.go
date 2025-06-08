package actor

import (
	"reflect"
	"strings"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/library/async"
	"universal/library/uerror"
)

type Actor struct {
	*async.Async
	name  string
	rval  reflect.Value
	funcs map[string]*FuncInfo
}

func (a *Actor) GetActorName() string {
	return a.name
}

func (a *Actor) Register(ac domain.IActor, _ ...int) {
	a.Async = async.NewAsync()
	a.rval = reflect.ValueOf(ac)
	a.name = parseName(a.rval.Elem().Type())
}

func (d *Actor) ParseFunc(tt interface{}) {
	switch vv := tt.(type) {
	case map[string]*FuncInfo:
		d.funcs = vv
	case reflect.Type:
		d.funcs = make(map[string]*FuncInfo)
		for i := 0; i < vv.NumMethod(); i++ {
			m := vv.Method(i)
			d.funcs[m.Name] = parseFuncInfo(m)
		}
	default:
		panic("注册参数错误，必须是方法列表或reflect.Type")
	}
}

func (d *Actor) SendMsg(h *pb.Head, args ...interface{}) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "%s.%s未实现, head:%v", h.ActorName, h.FuncName, h)
	}
	if mm.flag&CMD_FLAG == CMD_FLAG {
		d.Push(mm.localCmd(d.rval, h, args...))
	} else {
		d.Push(mm.localProto(d.rval, h, args...))
	}
	return nil
}

func (d *Actor) Send(h *pb.Head, buf []byte) error {
	mm, ok := d.funcs[h.FuncName]
	if !ok {
		return uerror.New(1, -1, "%s.%s未实现, head:%v", h.ActorName, h.FuncName, h)
	}
	// 发送事件
	if mm.flag&CMD_FLAG == CMD_FLAG {
		d.Push(mm.rpcCmd(d.rval, h, buf))
	} else if mm.flag&NOTIFY_FLAG > 0 {
		d.Push(mm.rpcNotify(d.rval, h, buf))
	} else if mm.flag&BYTES_FLAG > 0 {
		d.Push(mm.localProto(d.rval, h, buf))
	} else {
		d.Push(mm.rpcGob(d.rval, h, buf))
	}
	return nil
}

/*
func (d *Actor) RegisterTimer(h *pb.Head, ttl time.Duration, times int32) error {
	return timer.Register(d.GetIdPointer(), func() {
		if err := d.SendMsg(h); err != nil {
			mlog.Errorf("定时器发送消息失败: head:%v, error:%v", h, err)
		}
	}, ttl, times)
}
*/

func parseName(rr reflect.Type) string {
	name := rr.String()
	if index := strings.Index(name, "."); index > -1 {
		name = name[index+1:]
	}
	return name
}
