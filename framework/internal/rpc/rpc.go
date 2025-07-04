package rpc

import (
	"hash/crc32"
	"strings"
	"universal/common/pb"
	"universal/library/util"
)

var (
	cmds = make(map[pb.CMD]*RpcInfo)
	rpcs = make(map[pb.NodeType]*ActorRpc)
)

type RpcInfo struct {
	Cmd       pb.CMD
	NodeType  pb.NodeType
	IdType    uint32
	ActorFunc string
	Id        uint32
}

type ActorRpc struct {
	actors map[string]*RpcInfo
	values map[uint32]*RpcInfo
}

func Register(nt pb.NodeType, idType uint32, actorFunc string, ccs ...pb.CMD) {
	vals, ok := rpcs[nt]
	if !ok {
		vals = &ActorRpc{
			actors: make(map[string]*RpcInfo),
			values: make(map[uint32]*RpcInfo),
		}
		rpcs[nt] = vals
	}
	item := &RpcInfo{
		Cmd:       util.Index[pb.CMD](ccs, 0, pb.CMD_CMD_NONE),
		NodeType:  nt,
		IdType:    idType,
		ActorFunc: actorFunc,
		Id:        crc32.ChecksumIEEE([]byte(actorFunc)),
	}
	cmds[item.Cmd] = item
	vals.actors[item.ActorFunc] = item
	vals.values[item.Id] = item
}

func NewNodeRouter(nt pb.NodeType, actorFunc string, actorId uint64) *pb.NodeRouter {
	vals, ok := rpcs[nt]
	if !ok {
		return nil
	}
	api, ok := vals.actors[actorFunc]
	if !ok {
		return nil
	}
	return &pb.NodeRouter{
		NodeType:  nt,
		ActorFunc: api.Id,
		ActorId:   actorId<<8 | uint64(api.IdType&0xFF),
	}
}

func ParseNodeRouter(head *pb.Head, dst *pb.NodeRouter) {
	if dst == nil {
		return
	}
	vals, ok := rpcs[dst.NodeType]
	if !ok {
		return
	}
	api, ok := vals.values[dst.ActorFunc]
	if !ok {
		return
	}
	pos := strings.Index(api.ActorFunc, ".")
	head.ActorName = api.ActorFunc[:pos]
	head.FuncName = api.ActorFunc[pos+1:]
	head.ActorId = (dst.ActorId >> 8)
}
