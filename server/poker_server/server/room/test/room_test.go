package test

import (
	"fmt"
	"poker_server/common/pb"
	"poker_server/framework/mock"
	"poker_server/server/room/internal/module/card"
	"testing"
)

func TestCardParser(t *testing.T) {
	fmt.Printf("%v", card.Card(uint32(2147549186)).String())
}

// 测试换桌
func TestRummy(t *testing.T) {
	uid1 := uint64(144)
	roomId := uint64(30081548290)
	// -----------加入房间------------
	t.Run("JoinRoom1", func(t *testing.T) {
		mock.Request(uid1, roomId, pb.CMD_RUMMY_JOIN_ROOM_REQ, &pb.RummyJoinRoomReq{RoomId: roomId})
	})

	// ----------玩家准备------------
	t.Run("ReadyRoom1", func(t *testing.T) {
		mock.Request(uid1, roomId, pb.CMD_RUMMY_READY_ROOM_REQ, &pb.RummyReadyRoomReq{RoomId: roomId})
	})

	t.Run("ChangeRoom1", func(t *testing.T) {
		mock.Request(uid1, roomId, pb.CMD_RUMMY_CHANGE_ROOM_REQ, &pb.RummyJoinRoomReq{RoomId: roomId})
	})
}

func TestTexas(t *testing.T) {
	uid1, uid2, uid3 := uint64(144), uint64(145), uint64(100002)
	roomId := uint64(4311744514)
	// -----------加入房间------------
	t.Run("JoinRoom1", func(t *testing.T) {
		mock.Request(uid1, roomId, pb.CMD_TEXAS_JOIN_ROOM_REQ, &pb.TexasJoinRoomReq{RoomId: roomId})
	})
	t.Run("JoinRoom2", func(t *testing.T) {
		mock.Request(uid2, roomId, pb.CMD_TEXAS_JOIN_ROOM_REQ, &pb.TexasJoinRoomReq{RoomId: roomId})
	})
	t.Run("JoinRoom3", func(t *testing.T) {
		mock.Request(uid3, roomId, pb.CMD_TEXAS_JOIN_ROOM_REQ, &pb.TexasJoinRoomReq{RoomId: roomId})
	})

	t.Run("QuitRoom1", func(t *testing.T) {
		mock.Request(uid1, roomId, pb.CMD_TEXAS_QUIT_ROOM_REQ, &pb.TexasQuitRoomReq{RoomId: roomId})
	})
	t.Run("QuitRoom2", func(t *testing.T) {
		mock.Request(uid2, roomId, pb.CMD_TEXAS_QUIT_ROOM_REQ, &pb.TexasQuitRoomReq{RoomId: roomId})
	})

	// -----------买入------------
	t.Run("BuyIn1", func(t *testing.T) {
		mock.Request(uid1, roomId, pb.CMD_TEXAS_BUY_IN_REQ, &pb.TexasBuyInReq{RoomId: roomId, Chip: 100000, CoinType: 1})
	})
	t.Run("BuyIn2", func(t *testing.T) {
		mock.Request(uid2, roomId, pb.CMD_TEXAS_BUY_IN_REQ, &pb.TexasBuyInReq{RoomId: roomId, Chip: 100000, CoinType: 1})
	})
	t.Run("BuyIn3", func(t *testing.T) {
		mock.Request(uid3, roomId, pb.CMD_TEXAS_BUY_IN_REQ, &pb.TexasBuyInReq{RoomId: roomId, Chip: 100000, CoinType: 1})
	})

	// -----------坐下------------
	t.Run("SitDown1", func(t *testing.T) {
		mock.Request(uid1, roomId, pb.CMD_TEXAS_SIT_DOWN_REQ, &pb.TexasSitDownReq{RoomId: roomId, ChairId: 1})
	})
	t.Run("SitDown2", func(t *testing.T) {
		mock.Request(uid2, roomId, pb.CMD_TEXAS_SIT_DOWN_REQ, &pb.TexasSitDownReq{RoomId: roomId, ChairId: 2})
	})
	t.Run("SitDown3", func(t *testing.T) {
		mock.Request(uid3, roomId, pb.CMD_TEXAS_SIT_DOWN_REQ, &pb.TexasSitDownReq{RoomId: roomId, ChairId: 3})
	})

	// ------------站起------------
	t.Run("StandUp1", func(t *testing.T) {
		mock.Request(uid1, roomId, pb.CMD_TEXAS_STAND_UP_REQ, &pb.TexasStandUpReq{RoomId: roomId, ChairId: 1})
	})
	t.Run("StandUp2", func(t *testing.T) {
		mock.Request(uid2, roomId, pb.CMD_TEXAS_STAND_UP_REQ, &pb.TexasStandUpReq{RoomId: roomId, ChairId: 2})
	})
	t.Run("StandUp3", func(t *testing.T) {
		mock.Request(uid3, roomId, pb.CMD_TEXAS_STAND_UP_REQ, &pb.TexasStandUpReq{RoomId: roomId, ChairId: 3})
	})

	// ------------下注请求------------
	t.Run("Bet1", func(t *testing.T) {
		mock.Request(uid1, roomId, pb.CMD_TEXAS_DO_BET_REQ, &pb.TexasDoBetReq{Chip: 800, ChairId: 1, RoomId: roomId, OperateType: int32(pb.OperateType_BET)})
	})
	t.Run("Bet2", func(t *testing.T) {
		mock.Request(uid2, roomId, pb.CMD_TEXAS_DO_BET_REQ, &pb.TexasDoBetReq{Chip: 800, ChairId: 2, RoomId: roomId, OperateType: int32(pb.OperateType_BET)})
	})
	t.Run("Bet3", func(t *testing.T) {
		mock.Request(uid3, roomId, pb.CMD_TEXAS_DO_BET_REQ, &pb.TexasDoBetReq{Chip: 800, ChairId: 3, RoomId: roomId, OperateType: int32(pb.OperateType_BET)})
	})
}
