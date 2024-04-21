package manager

import (
	"fmt"
	"universal/framework/packet/domain"
	"universal/framework/packet/internal/repository"

	"google.golang.org/protobuf/proto"
)

var (
	packetPool = make(map[IndexPacket]domain.IPacket)
)

type IndexPacket struct {
	ActorName string
	FuncName  string
}

func RegisterCMD(h domain.CmdFunc, req, rsp proto.Message) {
	attr := repository.NewCmdPacket(h, req, rsp)
	index := IndexPacket{"", attr.GetFuncName()}
	if _, ok := packetPool[index]; ok {
		panic(fmt.Sprintf("CMD(%s) has already registered", attr.GetFuncName()))
	}
	packetPool[index] = attr
}

func RegiserFunc(h interface{}) {
	attr := repository.NewFuncPacket(h)
	index := IndexPacket{"", attr.GetFuncName()}
	if _, ok := packetPool[index]; ok {
		panic(fmt.Sprintf("Func(%s) has already registered", attr.GetFuncName()))
	}
	packetPool[index] = attr
}

func RegisterStruct(st interface{}) {
	for _, attr := range repository.NewStructPacket(st) {
		index := IndexPacket{attr.GetStructName(), attr.GetFuncName()}
		if _, ok := packetPool[index]; ok {
			panic(fmt.Sprintf("%s.%s() has already registered", attr.GetStructName(), attr.GetFuncName()))
		}
		packetPool[index] = attr
	}
}

func GetIPacket(actorName, fname string) domain.IPacket {
	if val, ok := packetPool[IndexPacket{actorName, fname}]; ok {
		return val
	}
	return repository.NewEmptyPacket(actorName, fname)
}
