package manager

import (
	"fmt"
	"reflect"
	"universal/common/pb"
	"universal/framework/basic"
	"universal/framework/packet/domain"
	"universal/framework/packet/internal/repository"

	"google.golang.org/protobuf/proto"
)

type index struct {
	ActorName string
	FuncName  string
}

type ActorMgr struct {
	req  reflect.Type
	rsp  reflect.Type
	apis map[index]*domain.IPacket
}

func newActorMgr(req, rsp proto.Message) *ActorMgr {
	return &ActorMgr{
		req: reflect.TypeOf(req).Elem(),
		rsp: reflect.TypeOf(rsp).Elem(),
	}
}

func (d *ActorMgr) Register(h interface{}) {
	for _, attr := range repository.NewStructPacket(st) {
		index := IndexPacket{attr.GetStructName(), attr.GetFuncName()}
		if _, ok := packetPool[index]; ok {
			panic(fmt.Sprintf("%s.%s() has already registered", attr.GetStructName(), attr.GetFuncName()))
		}
		packetPool[index] = attr
	}
}

func (d *ActorMgr) Call(ctx *basic.Context, pac *pb.Packet) (*pb.Packet, error) {
	return nil, nil
}
