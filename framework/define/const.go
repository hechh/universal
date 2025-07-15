package define

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
