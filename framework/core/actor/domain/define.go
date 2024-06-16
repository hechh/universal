package domain

import "universal/common/pb"

type ISend interface {
	GetName() string
	GetID() uint64
	Start()
	Stop()
	Register(IActor, uint64)
	SendMsg(*pb.RpcHead, ...interface{})
	Send(*pb.RpcHead, []byte)
}

type IActor interface {
	ISend
	Init()
}
