package handler

import (
	"strings"
	"universal/common/pb"
	"universal/framework/define"
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

func Register0[S any](nt pb.NodeType, actorFunc string, h ZeroProto[S]) {
	pos := strings.Index(actorFunc, ".")
	actorName, funcName := actorFunc[:pos], actorFunc[pos+1:]
	hs := GetActor(nt, actorName)
	hs[funcName] = h
	val2str[util.String2Int(actorFunc)] = actorFunc
}

func Register1[S any, T any](nt pb.NodeType, actorFunc string, h OneProto[S, T]) {
	pos := strings.Index(actorFunc, ".")
	actorName, funcName := actorFunc[:pos], actorFunc[pos+1:]
	hs := GetActor(nt, actorName)
	hs[funcName] = h
	val2str[util.String2Int(actorFunc)] = actorFunc
}

func Register2[S any, T any, R any](nt pb.NodeType, actorFunc string, h TwoProto[S, T, R]) {
	pos := strings.Index(actorFunc, ".")
	actorName, funcName := actorFunc[:pos], actorFunc[pos+1:]
	hs := GetActor(nt, actorName)
	hs[funcName] = h
	val2str[util.String2Int(actorFunc)] = actorFunc
}
