package define

import (
	"reflect"
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

const (
	HEAD_FLAG           = 1 << 0
	REQ_FLAG            = 1 << 1
	RSP_FLAG            = 1 << 2
	BYTES_FLAG          = 1 << 3
	INTERFACE_FLAG      = 1 << 4
	GOB_FLAG            = 1 << 5
	HEAD_REQ_RSP_MASK   = HEAD_FLAG | REQ_FLAG | RSP_FLAG
	HEAD_REQ_MASK       = HEAD_FLAG | REQ_FLAG
	HEAD_RSP_MASK       = HEAD_FLAG | RSP_FLAG
	REQ_RSP_MASK        = REQ_FLAG | RSP_FLAG
	HEAD_BYTES_MASK     = HEAD_FLAG | BYTES_FLAG
	HEAD_INTERFACE_MASK = HEAD_FLAG | INTERFACE_FLAG
)

var (
	HeadType      = reflect.TypeOf((*pb.Head)(nil))
	ReqType       = reflect.TypeOf((*proto.Message)(nil)).Elem()
	RspType       = reflect.TypeOf((*IRspProto)(nil)).Elem()
	BytesType     = reflect.TypeOf((*[]byte)(nil)).Elem()
	ErrorType     = reflect.TypeOf((*error)(nil)).Elem()
	InterfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
	NilError      = reflect.ValueOf((*error)(nil)).Elem()
)

const (
	ActorTypeUID    = 0
	ActorTypeRoomID = 1
)

func ToActorId(id uint64, actorType uint64) uint64 {
	return uint64(id<<8) | uint64(actorType&0xFF)
}

func UidToActorId(uid uint64) uint64 {
	return ToActorId(uid, ActorTypeUID)
}

func RoomIdToActorId(roomId uint64) uint64 {
	return ToActorId(roomId, ActorTypeRoomID)
}
