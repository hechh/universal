package cmd

import (
	"hash/crc32"
	"strings"
	"universal/common/pb"
	"universal/library/util"
)

var (
	names = make(map[string]uint32)
	cmds  = make(map[pb.CMD]*ApiInfo)
	apis  = make(map[uint32]*ApiInfo)
)

type ApiInfo struct {
	cmd       pb.CMD
	nodeType  pb.NodeType
	actorType uint32
	actorFunc string
	id        uint32
}

func get(actorFunc string) uint32 {
	if _, ok := names[actorFunc]; ok {
		names[actorFunc] = crc32.ChecksumIEEE([]byte(actorFunc))
	}
	return names[actorFunc]
}

func Get(val uint32) *ApiInfo {
	return apis[val]
}

func GetByCmd(cmd pb.CMD) *ApiInfo {
	return cmds[cmd]
}

func Register(nt pb.NodeType, actorType uint32, actorFunc string, cmdargs ...pb.CMD) {
	item := &ApiInfo{
		cmd:       util.Index[pb.CMD](cmdargs, 0, pb.CMD_CMD_NONE),
		nodeType:  nt,
		actorType: actorType,
		actorFunc: actorFunc,
		id:        get(actorFunc),
	}
	if item.cmd != pb.CMD_CMD_NONE {
		cmds[item.cmd] = item
	}
	apis[names[actorFunc]] = item
}

func NewNodeRouter(api *ApiInfo, actorId uint64) *pb.NodeRouter {
	return &pb.NodeRouter{NodeType: api.nodeType, ActorFunc: api.id, ActorId: actorId}
}

func ParseNodeRouter(head *pb.Head, actorFuncs ...string) {
	var ok bool
	var api *ApiInfo
	if head.Dst.ActorFunc > 0 {
		api, ok = apis[head.Dst.ActorFunc]
	} else if len(actorFuncs) > 0 {
		api, ok = apis[get(actorFuncs[0])]
	}
	if ok {
		pos := strings.Index(api.actorFunc, ".")
		head.ActorName = api.actorFunc[:pos]
		head.FuncName = api.actorFunc[pos+1:]
		head.ActorId = (head.Dst.ActorId >> 8)
	}
}
