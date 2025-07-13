package cmd

import (
	"universal/common/pb"
	"universal/library/util"
)

var (
	cmds = make(map[pb.CMD]*CmdInfo)
)

type CmdInfo struct {
	cmd       pb.CMD
	nodeType  pb.NodeType
	actorType uint32
	actorName string
	funcName  string
	uniqueId  uint32
}

func Register(cmd pb.CMD, nt pb.NodeType, actorType uint32, actorName, funcName string) {
	id := util.GetCrc32(actorName + "." + funcName)
	cmds[cmd] = &CmdInfo{cmd, nt, actorType, actorName, funcName, id}
}

func NewNodeRouter(cmd pb.CMD, actorId uint64) *pb.NodeRouter {
	info, ok := cmds[cmd]
	if !ok {
		return nil
	}
	return &pb.NodeRouter{
		NodeType:  info.nodeType,
		ActorFunc: info.uniqueId,
		ActorId:   uint64(actorId<<8) | uint64(info.actorType&0xFF),
	}
}
