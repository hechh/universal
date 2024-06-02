package pb

import (
	"reflect"

	"google.golang.org/protobuf/proto"
)

var (
	types = make(map[string]reflect.Type)
)

func NewType(name string) proto.Message {
	if vv, ok := types[name]; ok {
		return reflect.New(vv).Interface().(proto.Message)
	}
	return nil
}

func registerType(val interface{}) {
	ttt := reflect.TypeOf(val).Elem()
	types[ttt.Name()] = ttt
}

func init() {
	registerType((*ActorRequest)(nil))
	registerType((*ActorResponse)(nil))
	registerType((*GameLoginRequest)(nil))
	registerType((*GameLoginResponse)(nil))
	registerType((*GateLoginRequest)(nil))
	registerType((*GateLoginResponse)(nil))
	registerType((*LimitUpStockConfig)(nil))
	registerType((*LimitUpStockConfigAry)(nil))
	registerType((*Packet)(nil))
	registerType((*PacketHead)(nil))
	registerType((*RpcHead)(nil))
	registerType((*ServerNode)(nil))

}
