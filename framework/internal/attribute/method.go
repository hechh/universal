package attribute

import (
	"universal/common/pb"
	"universal/framework/define"
)

type Method struct {
	funcs map[string]define.IHandler
}

func NewMethod() *Method {
	return &Method{
		funcs: make(map[string]define.IHandler),
	}
}

func (m *Method) Register(name string, h define.IHandler) {
	m.funcs[name] = h
}

func (m *Method) Call(sendrsp define.SendRspFunc, s interface{}, head *pb.Head, args ...interface{}) func() {
	return m.funcs[head.FuncName].Call(sendrsp, s, head, args...)
}

func (m *Method) Rpc(sendrsp define.SendRspFunc, s interface{}, head *pb.Head, buf []byte) func() {
	return m.funcs[head.FuncName].Rpc(sendrsp, s, head, buf)
}
