package funcs

import (
	"hash/crc32"
	"strings"
	"universal/common/pb"
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

func getValue(actorFunc string) uint32 {
	if _, ok := names[actorFunc]; ok {
		names[actorFunc] = crc32.ChecksumIEEE([]byte(actorFunc))
	}
	return names[actorFunc]
}

func Register(nt pb.NodeType, actorFunc string, actorType uint32, cmdargs ...pb.CMD) {
	iid := getValue(actorFunc)
	item, ok := apis[iid]
	if !ok {
		item = &ApiInfo{
			nodeType:  nt,
			actorType: actorType,
			actorFunc: actorFunc,
			id:        iid,
		}
		apis[iid] = item
	}
	if len(cmdargs) > 0 {
		item.cmd = cmdargs[0]
	}
}

func NewNodeRouterByCmd(cmd pb.CMD, actorId uint64) *pb.NodeRouter {
	api, ok := cmds[cmd]
	if !ok {
		return nil
	}
	return &pb.NodeRouter{NodeType: api.nodeType, ActorFunc: api.id, ActorId: actorId}
}

func ParseNodeRouter(head *pb.Head, actorFuncs ...string) {
	var ok bool
	var api *ApiInfo
	if head.Dst.ActorFunc > 0 {
		api, ok = apis[head.Dst.ActorFunc]
	} else if len(actorFuncs) > 0 {
		api, ok = apis[getValue(actorFuncs[0])]
	}
	if ok {
		pos := strings.Index(api.actorFunc, ".")
		head.ActorName = api.actorFunc[:pos]
		head.FuncName = api.actorFunc[pos+1:]
		head.ActorId = (head.Dst.ActorId >> 8)
	}
}
