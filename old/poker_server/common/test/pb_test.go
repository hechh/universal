package test

import (
	"poker_server/common/pb"
	"poker_server/framework/mock"
	"poker_server/library/util"
	"testing"

	"github.com/golang/protobuf/proto"
)

func BenchmarkUnmarshal1(b *testing.B) {
	buf, _ := proto.Marshal(&pb.Packet{
		Head: &pb.Head{ActorName: "afdasfa", FuncName: "afaeasdfa", ActorId: 1234134, Cmd: 123413, Seq: 123},
		Body: []byte("asdfasdfasdfklaj;lkjl;aksjdf"),
	})

	for i := 0; i < b.N; i++ {
		msg := pb.Packet{}
		proto.Unmarshal(buf, &msg)
	}
}

var (
	pools = util.Pool[pb.Packet]()
)

func get() *pb.Packet {
	return pools.Get().(*pb.Packet)
}

func put(pac *pb.Packet) {
	pools.Put(pac)
}

func BenchmarkUnmarshal2(b *testing.B) {
	buf, _ := proto.Marshal(&pb.Packet{
		Head: &pb.Head{ActorName: "afdasfa", FuncName: "afaeasdfa", ActorId: 1234134, Cmd: 123413, Seq: 123},
		Body: []byte("asdfasdfasdfklaj;lkjl;aksjdf"),
	})

	for i := 0; i < b.N; i++ {
		msg := get()
		proto.Unmarshal(buf, msg)
		put(msg)
	}
}

func TestTexas(t *testing.T) {
	mock.Request(0, 0, pb.CMD_TEXAS_JOIN_ROOM_REQ, &pb.TexasJoinRoomReq{TableId: 1})
}
