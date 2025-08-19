package test

import (
	"encoding/json"
	"fmt"
	"poker_server/common/pb"
	"poker_server/framework/mock"
	"poker_server/server/room/internal/module/card"
	"testing"
)

func TestCardParser(t *testing.T) {
	fmt.Printf("%v", card.Card(uint32(2147549186)).String())
}

func TestJsonParser(t *testing.T) {
	jsonStr := `{"resp_msg":{"code":-1,"message":"currency code does not exist"},"resp_data":[]}`
	rsp := &pb.HttpTransferInRsp{}
	json.Unmarshal([]byte(jsonStr), rsp)
	fmt.Printf("%v", rsp)
}

// 测试换桌
func TestRummy(t *testing.T) {
	uid1 := uint64(1000222)
	uid2 := uint64(1000223)
	uid3 := uint64(1000224)
	uid4 := uint64(1000225)
	uid5 := uint64(1000226)
	uid6 := uint64(1000227)
	roomId := uint64(30081548290)
	//roomId := uint64(34376515597)
	// -----------匹配加入------------ todo 重连
	t.Run("matchRoom1", func(t *testing.T) {
		err := mock.Request(uid1, roomId, pb.CMD_RUMMY_MATCH_REQ, &pb.RummyMatchReq{
			CfgId:    9,
			GameType: pb.GameType(8),
			CoinType: pb.CoinType(1),
		})
		t.Logf("%v", err)
	})

	// -----------分支玩法退出房间------ todo playermgr
	t.Run("quitRoomExt", func(t *testing.T) {
		err := mock.Request(uid1, roomId, pb.CMD_RUMMY_GIVE_UP_REQ, &pb.RummyGiveUpReq{
			RoomId: roomId,
		})
		t.Logf("%v", err)
	})

	t.Run("cancelmatch", func(t *testing.T) {
		err := mock.Request(uid1, roomId, pb.CMD_RUMMY_CANCEL_MATCH_REQ, &pb.RummyCancelMatchReq{
			CfgId:    9,
			GameType: pb.GameType(8),
			CoinType: pb.CoinType(1),
		})
		t.Logf("%v", err)
	})

	// -----------加入房间------------
	t.Run("JoinRoom1", func(t *testing.T) {
		err := mock.Request(uid1, roomId, pb.CMD_RUMMY_JOIN_ROOM_REQ, &pb.RummyJoinRoomReq{RoomId: roomId})
		t.Logf("%v", err)
	})

	// ----------玩家准备------------
	t.Run("ReadyRoom1", func(t *testing.T) {
		mock.Request(uid1, roomId, pb.CMD_RUMMY_READY_ROOM_REQ, &pb.RummyReadyRoomReq{RoomId: roomId})
	})

	// ----------玩家退出------------
	t.Run("QuitRoom1", func(t *testing.T) {
		mock.Request(uid1, roomId, pb.CMD_RUMMY_QUIT_ROOM_REQ, &pb.RummyQuitRoomReq{RoomId: roomId})
	})

	// -----------加入房间------------
	t.Run("JoinRoom2", func(t *testing.T) {
		mock.Request(uid2, roomId, pb.CMD_RUMMY_JOIN_ROOM_REQ, &pb.RummyJoinRoomReq{RoomId: roomId})
	})

	// ----------玩家准备------------
	t.Run("ReadyRoom2", func(t *testing.T) {
		mock.Request(uid2, roomId, pb.CMD_RUMMY_READY_ROOM_REQ, &pb.RummyReadyRoomReq{RoomId: roomId})
	})

	// ----------玩家准备------------
	t.Run("QuitRoom2", func(t *testing.T) {
		mock.Request(uid2, roomId, pb.CMD_RUMMY_QUIT_ROOM_REQ, &pb.RummyQuitRoomReq{RoomId: roomId})
	})

	// -----------加入房间------------
	t.Run("JoinRoom3", func(t *testing.T) {
		mock.Request(uid3, roomId, pb.CMD_RUMMY_JOIN_ROOM_REQ, &pb.RummyJoinRoomReq{RoomId: roomId})
	})

	// ----------玩家准备------------
	t.Run("ReadyRoom3", func(t *testing.T) {
		mock.Request(uid3, roomId, pb.CMD_RUMMY_READY_ROOM_REQ, &pb.RummyReadyRoomReq{RoomId: roomId})
	})

	// ----------玩家准备------------
	t.Run("QuitRoom3", func(t *testing.T) {
		mock.Request(uid3, roomId, pb.CMD_RUMMY_QUIT_ROOM_REQ, &pb.RummyQuitRoomReq{RoomId: roomId})
	})

	// -----------加入房间------------
	t.Run("JoinRoom4", func(t *testing.T) {
		mock.Request(uid4, roomId, pb.CMD_RUMMY_JOIN_ROOM_REQ, &pb.RummyJoinRoomReq{RoomId: roomId})
	})

	// ----------玩家准备------------
	t.Run("ReadyRoom4", func(t *testing.T) {
		mock.Request(uid4, roomId, pb.CMD_RUMMY_READY_ROOM_REQ, &pb.RummyReadyRoomReq{RoomId: roomId})
	})

	// ----------玩家准备------------
	t.Run("QuitRoom4", func(t *testing.T) {
		mock.Request(uid4, roomId, pb.CMD_RUMMY_QUIT_ROOM_REQ, &pb.RummyQuitRoomReq{RoomId: roomId})
	})

	// -----------加入房间------------
	t.Run("JoinRoom5", func(t *testing.T) {
		mock.Request(uid5, roomId, pb.CMD_RUMMY_JOIN_ROOM_REQ, &pb.RummyJoinRoomReq{RoomId: roomId})
	})

	// ----------玩家准备------------
	t.Run("ReadyRoom5", func(t *testing.T) {
		mock.Request(uid5, roomId, pb.CMD_RUMMY_READY_ROOM_REQ, &pb.RummyReadyRoomReq{RoomId: roomId})
	})

	// ----------玩家准备------------
	t.Run("QuitRoom5", func(t *testing.T) {
		mock.Request(uid5, roomId, pb.CMD_RUMMY_QUIT_ROOM_REQ, &pb.RummyQuitRoomReq{RoomId: roomId})
	})

	// -----------加入房间------------
	t.Run("JoinRoom6", func(t *testing.T) {
		mock.Request(uid6, roomId, pb.CMD_RUMMY_JOIN_ROOM_REQ, &pb.RummyJoinRoomReq{RoomId: roomId})
	})

	// ----------玩家准备------------
	t.Run("ReadyRoom6", func(t *testing.T) {
		mock.Request(uid6, roomId, pb.CMD_RUMMY_READY_ROOM_REQ, &pb.RummyReadyRoomReq{RoomId: roomId})
	})

	// ----------玩家准备------------
	t.Run("QuitRoom6", func(t *testing.T) {
		mock.Request(uid6, roomId, pb.CMD_RUMMY_QUIT_ROOM_REQ, &pb.RummyQuitRoomReq{RoomId: roomId})
	})
}

func TestTexas(t *testing.T) {
	uid1, uid2, uid3 := uint64(144), uint64(145), uint64(100002)
	roomId := uint64(4311744517)
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

	// -----------换房间------------
	t.Run("ChangeRoom1", func(t *testing.T) {
		mock.Request(uid1, roomId, pb.CMD_TEXAS_CHANGE_ROOM_REQ, &pb.TexasChangeRoomReq{RoomId: roomId})
	})
	t.Run("ChangeRoom2", func(t *testing.T) {
		mock.Request(uid2, roomId, pb.CMD_TEXAS_CHANGE_ROOM_REQ, &pb.TexasChangeRoomReq{RoomId: roomId})
	})

	// -----------退出房间------------
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
