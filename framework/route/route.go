package route

import (
	"hash/crc32"
	"universal/common/pb"
)

var (
	cmds   = make(map[pb.CMD]*RouteInfo)
	actors = make(map[pb.NodeType]map[string]*RouteInfo)
	values = make(map[pb.NodeType]map[uint32]*RouteInfo)
)

type RouteInfo struct {
	Cmd       pb.CMD
	NodeType  pb.NodeType
	IdType    uint32
	ActorName string
	FuncName  string
	ActorFunc uint32
}

func Register(nt pb.NodeType, idType uint32, actorName, funcName string) {
	RegisterCMD(pb.CMD_CMD_NONE, nt, idType, actorName, funcName)
}

func RegisterCMD(cmd pb.CMD, nt pb.NodeType, idType uint32, actorName, funcName string) {
	key := actorName + "." + funcName
	item := &RouteInfo{
		Cmd:       cmd,
		NodeType:  nt,
		IdType:    idType,
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

func register[T comparable](key T, item *RouteInfo, tmps map[pb.NodeType]map[T]*RouteInfo) {
	if vals, ok := tmps[item.NodeType]; !ok {
		vals = make(map[T]*RouteInfo)
		vals[key] = item
		tmps[item.NodeType] = vals
	} else {
		vals[key] = item
	}
}

func NewNRouter(nt pb.NodeType, actorFunc string, actorId uint64) *pb.NodeRouter {
	if api := get[string](nt, actorFunc, actors); api != nil {
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
		if api := get[uint32](dst.NodeType, dst.ActorFunc, values); api != nil {
			head.ActorName = api.ActorName
			head.FuncName = api.FuncName
			head.ActorId = parse(dst.ActorId)
		}
	}
}

func get[T comparable](nt pb.NodeType, actorFunc T, tmps map[pb.NodeType]map[T]*RouteInfo) *RouteInfo {
	if vals, ok := tmps[nt]; ok {
		return vals[actorFunc]
	}
	return nil
}

func to(otherid uint64, idType uint32) uint64 {
	return (otherid << 8) | uint64(idType&0xFF)
}

func parse(actorId uint64) uint64 {
	return actorId >> 8
}
