package pb

import reflect "reflect"

var (
	types = make(map[string]reflect.Type)
)

func init() {
	types["ActorRequest"] = reflect.TypeOf((*ActorRequest)(nil)).Elem()
	types["ActorResponse"] = reflect.TypeOf((*ActorResponse)(nil)).Elem()
	types["GameLoginRequest"] = reflect.TypeOf((*GameLoginRequest)(nil)).Elem()
	types["GameLoginResponse"] = reflect.TypeOf((*GameLoginResponse)(nil)).Elem()
	types["GateLoginRequest"] = reflect.TypeOf((*GateLoginRequest)(nil)).Elem()
	types["GateLoginResponse"] = reflect.TypeOf((*GateLoginResponse)(nil)).Elem()
	types["LimitUpStockConfig"] = reflect.TypeOf((*LimitUpStockConfig)(nil)).Elem()
	types["LimitUpStockConfigAry"] = reflect.TypeOf((*LimitUpStockConfigAry)(nil)).Elem()
	types["Packet"] = reflect.TypeOf((*Packet)(nil)).Elem()
	types["PacketHead"] = reflect.TypeOf((*PacketHead)(nil)).Elem()
	types["RpcHead"] = reflect.TypeOf((*RpcHead)(nil)).Elem()
	types["ServerNode"] = reflect.TypeOf((*ServerNode)(nil)).Elem()

}
