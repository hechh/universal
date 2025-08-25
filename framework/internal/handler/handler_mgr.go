package handler

import (
	"universal/common/pb"
	"universal/framework/define"
	"universal/library/util"
)

var (
	val2str = make(map[uint32]string)
	actors  = make(map[uint32]define.IHandler)
)

func RegisterHandler[S any, T any, R any](actorFunc string, h Handler[S, T, R]) {
	val := util.String2Int(actorFunc)
	val2str[val] = actorFunc
	actors[val] = h
}

func RegisterNotify[S any, T any](actorFunc string, h Notify[S, T]) {
	val := util.String2Int(actorFunc)
	val2str[val] = actorFunc
	actors[val] = h
}

func RegisterTrigger[S any](actorFunc string, h Trigger[S]) {
	val := util.String2Int(actorFunc)
	val2str[val] = actorFunc
	actors[val] = h
}

func Has(head *pb.Head) bool {
	_, ok := actors[head.ActorFunc]
	return ok
}

func GetActorFunc(val uint32) string {
	return val2str[val]
}

func Call(sendrsp define.SendRspFunc, s interface{}, head *pb.Head, args ...interface{}) func() {
	return actors[head.ActorFunc].Call(sendrsp, s, head, args...)
}

func Rpc(sendrsp define.SendRspFunc, s interface{}, head *pb.Head, buf []byte) func() {
	return actors[head.ActorFunc].Rpc(sendrsp, s, head, buf)
}
