package actor

import (
	"universal/common/pb"
	"universal/framework/actor/domain"
	"universal/framework/actor/internal/manager"
)

func SetPacketHandle(h domain.PacketHandle) {
	manager.SetPacketHandle(h)
}

func Send(key string, pac *pb.Packet) {
	manager.GetSession(key).Send(pac)
}
