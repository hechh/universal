package rpc

import (
	"strings"
	"universal/common/pb"
	"universal/library/util"
)

var (
	cmds  = make(map[pb.CMD]*CmdInfo)
	types = make(map[pb.NodeType]map[string]uint32)
)

type CmdInfo struct {
	cmd       pb.CMD
	nodeType  pb.NodeType
	actorFunc string
}

func Register(nt pb.NodeType, actorName string, actorType uint32) {
	if _, ok := types[nt]; !ok {
		types[nt] = make(map[string]uint32)
	}
	types[nt][actorName] = actorType
}

func RegisterCmd(nt pb.NodeType, actorFunc string, cc pb.CMD) {
	cmds[cc] = &CmdInfo{
		cmd:       cc,
		nodeType:  nt,
		actorFunc: actorFunc,
	}
}

func NewNodeRouterByCmd(cmd pb.CMD, actorId uint64) *pb.NodeRouter {
	api, ok := cmds[cmd-cmd%2]
	if !ok {
		return nil
	}
	vals, ok := types[api.nodeType]
	if !ok {
		return nil
	}
	actorType, ok := vals[prefixString(api.actorFunc, strings.Index(api.actorFunc, "."))]
	if !ok {
		return nil
	}
	return &pb.NodeRouter{
		NodeType:  api.nodeType,
		ActorFunc: api.actorFunc,
		ActorId:   actorId<<8 | uint64(actorType&0xFF),
	}
}

func NewNodeRouter(nt pb.NodeType, actorFunc string, actorId uint64) *pb.NodeRouter {
	vals, ok := types[nt]
	if !ok {
		return nil
	}
	actorType, ok := vals[prefixString(actorFunc, strings.Index(actorFunc, "."))]
	if !ok {
		return nil
	}
	return &pb.NodeRouter{
		NodeType:  nt,
		ActorFunc: actorFunc,
		ActorId:   actorId<<8 | uint64(actorType&0xFF),
	}
}

func ParseNodeRouter(head *pb.Head, actorFuncs ...string) {
	if head.Dst == nil {
		return
	}
	actorFunc := util.Index[string](actorFuncs, 0, head.Dst.ActorFunc)
	pos := strings.Index(actorFunc, ".")
	head.ActorName = actorFunc[:pos]
	head.FuncName = actorFunc[pos+1:]
	head.ActorId = (head.Dst.ActorId >> 8)
}

func prefixString(str string, pos int) string {
	if pos < 0 || pos > len(str) {
		return str
	}
	return str[:pos]
}
