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

type ActorRpc struct {
	actors map[string]*RpcInfo
	values map[uint32]*RpcInfo
}

type RpcInfo struct {
	cmd       pb.CMD
	nodeType  pb.NodeType
	idType    uint32
	actorFunc string
	id        uint32
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
		cmd:       util.Index[pb.CMD](ccs, 0, pb.CMD_CMD_NONE),
		nodeType:  nt,
		idType:    idType,
		actorFunc: actorFunc,
		id:        crc32.ChecksumIEEE([]byte(actorFunc)),
	}
	cmds[item.cmd] = item
	vals.actors[item.actorFunc] = item
	vals.values[item.id] = item
}

func NewNodeRouterByCmd(cmd pb.CMD, actorId uint64) *pb.NodeRouter {
	api, ok := cmds[cmd-cmd%2]
	if !ok {
		return nil
	}
	return &pb.NodeRouter{
		NodeType:  api.nodeType,
		ActorFunc: api.id,
		ActorId:   actorId<<8 | uint64(api.idType&0xFF),
	}
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
		ActorFunc: api.id,
		ActorId:   actorId<<8 | uint64(api.idType&0xFF),
	}
}

func ParseNodeRouter(head *pb.Head, actorFuncs ...string) {
	if head.Dst == nil {
		return
	}
	vals, ok := rpcs[head.Dst.NodeType]
	if !ok {
		return
	}
	var api *RpcInfo
	if head.Dst.ActorFunc > 0 {
		api, ok = vals.values[head.Dst.ActorFunc]
	} else {
		api, ok = vals.actors[util.Index[string](actorFuncs, 0, "")]
	}
	if ok {
		pos := strings.Index(api.actorFunc, ".")
		head.ActorName = api.actorFunc[:pos]
		head.FuncName = api.actorFunc[pos+1:]
		head.ActorId = (head.Dst.ActorId >> 8)
	}
}
