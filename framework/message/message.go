package message

import (
	"hash/crc32"
	"universal/common/pb"
)

var (
	cmds   = make(map[pb.CMD]*ApiInfo)
	actors = make(map[pb.NodeType]map[string]*ApiInfo)
	values = make(map[pb.NodeType]map[uint32]*ApiInfo)
)

type ApiInfo struct {
	Cmd       pb.CMD
	NodeType  pb.NodeType
	IdType    uint32
	ActorName string
	FuncName  string
	ActorFunc uint32
}

func register[T comparable](key T, item *ApiInfo, tmps map[pb.NodeType]map[T]*ApiInfo) {
	if vals, ok := tmps[item.NodeType]; !ok {
		vals = make(map[T]*ApiInfo)
		vals[key] = item
		tmps[item.NodeType] = vals
	} else {
		vals[key] = item
	}
}

func Register(cmd pb.CMD, nt pb.NodeType, routerType uint32, actorName, funcName string) {
	key := actorName + "." + funcName
	item := &ApiInfo{
		Cmd:       cmd,
		NodeType:  nt,
		IdType:    routerType,
		ActorName: actorName,
		FuncName:  funcName,
		ActorFunc: crc32.ChecksumIEEE([]byte(key)),
	}
	if cmd > pb.CMD_CMD_NONE {
		cmds[cmd] = item
	}
	register[string](key, item, actors)
	register[uint32](item.ActorFunc, item, values)
}

func to(otherid uint64, idType uint32) uint64 {
	return (otherid << 8) | uint64(idType&0xFF)
}

func parse(actorId uint64) uint64 {
	return actorId >> 8
}

func get(nt pb.NodeType, actorFunc string) *ApiInfo {
	if vals, ok := actors[nt]; ok {
		return vals[actorFunc]
	}
	return nil
}

func getByValue(nt pb.NodeType, actorFunc uint32) *ApiInfo {
	if vals, ok := values[nt]; ok {
		return vals[actorFunc]
	}
	return nil
}

func NewNRouter(nt pb.NodeType, actorFunc string, actorId uint64) *pb.NodeRouter {
	if api := get(nt, actorFunc); api != nil {
		return &pb.NodeRouter{
			NodeType:  nt,
			ActorFunc: api.ActorFunc,
			ActorId:   to(actorId, api.IdType),
		}
	}
	return nil
}

func ParseNRouter(head *pb.Head, dst *pb.NodeRouter) {
	if dst != nil {
		if api := getByValue(dst.NodeType, dst.ActorFunc); api != nil {
			head.ActorName = api.ActorName
			head.FuncName = api.FuncName
			head.ActorId = parse(head.Dst.ActorId)
		}
	}
}
