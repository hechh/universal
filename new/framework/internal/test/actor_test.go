package test

import (
	"fmt"
	"poker_server/common/pb"
	"poker_server/framework/internal/core/actor"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
)

type Player struct {
	actor.Actor
}

func (p *Player) HeartRequest(head *pb.Head, req *pb.GateHeartRequest, rsp *pb.GateHeartResponse) error {
	fmt.Println("-------->", head, req)
	return nil
}

func TestActor(t *testing.T) {
	player := &Player{}
	player.Actor.Register(player)
	player.Actor.ParseFunc(reflect.TypeOf(player))
	player.Start()
	actor.Register(player)

	head := &pb.Head{
		SendType:    pb.SendType_POINT,
		DstNodeType: pb.NodeType_Gate,
		DstNodeId:   1,
		IdType:      pb.IdType_UID,
		Id:          1000001,
		RouteId:     1000001,
		Cmd:         uint32(pb.CMD_GATE_HEART_REQUEST),
		ActorName:   "Player",
		FuncName:    "HeartRequest",
	}
	buf, _ := proto.Marshal(&pb.GateHeartRequest{})

	actor.Send(head, buf)

	select {}
}
