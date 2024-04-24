package manager

import (
	"fmt"
	"universal/framework/packet/domain"
	"universal/framework/packet/internal/repository"

	"google.golang.org/protobuf/proto"
)

var (
	apiPool = make(map[int32]domain.IPacket)
)

func RegisterApi(apiCode int32, h domain.ApiFunc, req, rsp proto.Message) {
	attr := repository.NewApiPacket(h, req, rsp)
	if _, ok := apiPool[apiCode]; ok {
		panic(fmt.Sprintf("ApiCode(%d) has already registered", apiCode))
	}
	apiPool[apiCode] = attr
}

func RegisterActor(apiCode int32, h interface{}) {

}

/*
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
*/
