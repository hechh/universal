package test

import (
	"back-end/common/client_pb/room"
	"back-end/common/client_pb/rummy"
	"back-end/driver/structs/proto/gate_rpc"
	"context"
	"testing"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

// 1746838592048652300

func TestRequest(t *testing.T) {
	// 测试加入房间
	TexasJoinRoomReq(1746838592048652300, 10000051)
	TexasJoinRoomReq(1746838592048652300, 10000052)
	// 测试买入
	BuyInReq(1746838592048652300, 10000051)
	BuyInReq(1746838592048652300, 10000052)
	// 测试坐下
	SitDownReq(1746838592048652300, 10000051, 1)
	SitDownReq(1746838592048652300, 10000052, 2)
	// 测试站起
	StandUpReq(1746838592048652300, 10000051, 1)
	StandUpReq(1746838592048652300, 10000052, 2)
}

func TexasJoinRoomReq(roomId, uid uint64) {
	// 请替换为实际的 gRPC 服务地址
	conn, _ := grpc.Dial("localhost:13001", grpc.WithInsecure())
	defer conn.Close()
	client := gate_rpc.NewMessageHandlerClient(conn)
	// 构建请求消息
	request := &room.TexasJoinRoomReq{
		RoomId: proto.Uint64(roomId),
	}
	buf, _ := proto.Marshal(request)
	// 发送消息
	client.HandleClientMessage(context.Background(), &gate_rpc.ClientMessage{
		PlayerId:    uid,
		RequestData: buf,
		Header:      &gate_rpc.Header{MsgId: uint32(rummy.RummyMsgID_TexasJoinRoomReq)},
	})
}

func BuyInReq(roomId, uid uint64) {
	// 请替换为实际的 gRPC 服务地址
	conn, _ := grpc.Dial("localhost:13001", grpc.WithInsecure())
	defer conn.Close()
	client := gate_rpc.NewMessageHandlerClient(conn)
	// 构建请求消息
	request := &room.TexasBuyInReq{
		RoomId:   proto.Uint64(roomId),
		Chip:     proto.Int64(100000),
		CoinType: proto.Int32(1),
	}
	buf, _ := proto.Marshal(request)
	// 发送消息
	client.HandleClientMessage(context.Background(), &gate_rpc.ClientMessage{
		PlayerId:    uid,
		RequestData: buf,
		Header:      &gate_rpc.Header{MsgId: uint32(rummy.RummyMsgID_TexasBuyInReq)},
	})
}

func SitDownReq(roomId, uid uint64, chairId uint32) {
	// 请替换为实际的 gRPC 服务地址
	conn, _ := grpc.Dial("localhost:13001", grpc.WithInsecure())
	defer conn.Close()
	client := gate_rpc.NewMessageHandlerClient(conn)
	// 构建请求消息
	request := &room.TexasSitDownReq{
		RoomId:  proto.Uint64(roomId),
		ChairId: proto.Uint32(chairId),
	}
	buf, _ := proto.Marshal(request)
	// 发送消息
	client.HandleClientMessage(context.Background(), &gate_rpc.ClientMessage{
		PlayerId:    uid,
		RequestData: buf,
		Header:      &gate_rpc.Header{MsgId: uint32(rummy.RummyMsgID_TexasSitDownReq)},
	})
}

func StandUpReq(roomId, uid uint64, chairId uint32) {
	// 请替换为实际的 gRPC 服务地址
	conn, _ := grpc.Dial("localhost:13001", grpc.WithInsecure())
	defer conn.Close()
	client := gate_rpc.NewMessageHandlerClient(conn)
	// 构建请求消息
	request := &room.TexasStandUpReq{
		RoomId:  proto.Uint64(roomId),
		ChairId: proto.Uint32(chairId),
	}
	buf, _ := proto.Marshal(request)
	// 发送消息
	client.HandleClientMessage(context.Background(), &gate_rpc.ClientMessage{
		PlayerId:    uid,
		RequestData: buf,
		Header:      &gate_rpc.Header{MsgId: uint32(rummy.RummyMsgID_TexasStandUpRsp)},
	})
}
