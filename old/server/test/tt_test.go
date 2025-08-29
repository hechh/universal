package test

import (
	"hash/crc32"
	"testing"
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

func TestPacket(t *testing.T) {
	actorFunc := "PlayerMgr.LoginRequest"
	t.Log(crc32.ChecksumIEEE([]byte(actorFunc)))

	bb := &pb.NodeRouter{
		Type:      pb.NodeType_NodeTypeGate,
		Id:        1,
		ActorFunc: crc32.ChecksumIEEE([]byte(actorFunc)),
		ActorId:   12,
	}
	bufb, _ := proto.Marshal(bb)
	t.Log(len(bufb))
}
