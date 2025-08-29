package handler

import (
	"strings"
	"universal/common/pb"
	"universal/framework/define"
	"universal/framework/internal/entity"
	"universal/library/util"
)

var (
	val2str = make(map[uint32]string)
	apis    = make(map[pb.NodeType]map[string]map[string]define.IHandler)
)

func GetActor(nt pb.NodeType, actorName string) map[string]define.IHandler {
	mm, ok := apis[nt]
	if !ok {
		mm = make(map[string]map[string]define.IHandler)
		apis[nt] = mm
	}
	hs, ok := mm[actorName]
	if !ok {
		hs = make(map[string]define.IHandler)
		mm[actorName] = hs
	}
	return hs
}

func GetHandler(nt pb.NodeType, actorName, funcName string) define.IHandler {
	return GetActor(nt, actorName)[funcName]
}

func RegisterTrigger[S any](nt pb.NodeType, actorFunc string, h entity.TriggerHandler[S]) {
	pos := strings.Index(actorFunc, ".")
	actorName, funcName := actorFunc[:pos], actorFunc[pos+1:]
	hs := GetActor(nt, actorName)
	hs[funcName] = h
	val2str[util.String2Int(actorFunc)] = actorFunc
}

func RegisterEvent[S any, T any](nt pb.NodeType, actorFunc string, h entity.EventHandler[S, T]) {
	pos := strings.Index(actorFunc, ".")
	actorName, funcName := actorFunc[:pos], actorFunc[pos+1:]
	hs := GetActor(nt, actorName)
	hs[funcName] = h
	val2str[util.String2Int(actorFunc)] = actorFunc
}

func RegisterCmd[S any, T any, R any](nt pb.NodeType, actorFunc string, h entity.CmdHandler[S, T, R]) {
	pos := strings.Index(actorFunc, ".")
	actorName, funcName := actorFunc[:pos], actorFunc[pos+1:]
	hs := GetActor(nt, actorName)
	hs[funcName] = h
	val2str[util.String2Int(actorFunc)] = actorFunc
}

func RegisterGob1[S any, T any](nt pb.NodeType, actorFunc string, h entity.Gob1Handler[S, T]) {
	pos := strings.Index(actorFunc, ".")
	actorName, funcName := actorFunc[:pos], actorFunc[pos+1:]
	hs := GetActor(nt, actorName)
	hs[funcName] = h
	val2str[util.String2Int(actorFunc)] = actorFunc
}

func RegisterGob2[S any, T any, U any](nt pb.NodeType, actorFunc string, h entity.Gob2Handler[S, T, U]) {
	pos := strings.Index(actorFunc, ".")
	actorName, funcName := actorFunc[:pos], actorFunc[pos+1:]
	hs := GetActor(nt, actorName)
	hs[funcName] = h
	val2str[util.String2Int(actorFunc)] = actorFunc
}

func RegisterGob3[S any, T any, U any, A any](nt pb.NodeType, actorFunc string, h entity.Gob3Handler[S, T, U, A]) {
	pos := strings.Index(actorFunc, ".")
	actorName, funcName := actorFunc[:pos], actorFunc[pos+1:]
	hs := GetActor(nt, actorName)
	hs[funcName] = h
	val2str[util.String2Int(actorFunc)] = actorFunc
}

func GenRouterId(id uint64, tt uint64) uint64 {
	return (id << 8) | tt
}

func ParseRouterId(routerId uint64) uint64 {
	return (routerId >> 8)
}

func GetActorFunc(id uint32) (actorName string, funcName string) {
	str, ok := val2str[id]
	if !ok {
		return
	}
	pos := strings.Index(str, ".")
	actorName = str[:pos]
	funcName = str[pos+1:]
	return
}

func GetActorFuncId(val string) uint32 {
	id := util.String2Int(val)
	val2str[id] = val
	return id
}
