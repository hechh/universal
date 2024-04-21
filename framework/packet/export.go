package actor

import (
	"universal/framework/packet/domain"
	"universal/framework/packet/internal/manager"

	"google.golang.org/protobuf/proto"
)

func RegisterCMD(h domain.CmdFunc, req, rsp proto.Message) {
	manager.RegisterCMD(h, req, rsp)
}

func RegiserFunc(h interface{}) {
	manager.RegiserFunc(h)
}

func RegisterStruct(st interface{}) {
	manager.RegisterStruct(st)
}

func GetIPacket(actorName, fname string) domain.IPacket {
	return manager.GetIPacket(actorName, fname)
}
