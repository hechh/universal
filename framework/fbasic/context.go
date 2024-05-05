package fbasic

import (
	"universal/common/pb"
)

type IData interface {
	//	ToBytes() ([]byte, error)
}

type Context struct {
	*pb.PacketHead                  // rpc请求头
	readyOnlys     map[string]IData // 零时缓存
}

func NewContext(head *pb.PacketHead, datas map[string]IData) *Context {
	return &Context{
		PacketHead: head,
		readyOnlys: datas,
	}
}

func (d *Context) GetValue(key string) IData {
	return d.readyOnlys[key]
}

func ApiCodeToClusterType(val int32) (typ pb.ClusterType) {
	switch val / 1000000 {
	case 1:
		typ = pb.ClusterType_GATE
	case 2:
		typ = pb.ClusterType_GAME
	default:
		typ = pb.ClusterType_NONE
	}
	return
}
