package playerfun

import (
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/library/uerror"
	"poker_server/library/util"
	"poker_server/server/game/internal/player/domain"

	"github.com/golang/protobuf/proto"
)

type PlayerBaseFun struct {
	*PlayerFun
	data *pb.PlayerDataBase
}

func NewPlayerBaseFun(fun *PlayerFun) domain.IPlayerFun {
	return &PlayerBaseFun{PlayerFun: fun}
}

func (d *PlayerBaseFun) Load(msg *pb.PlayerData) error {
	if msg == nil {
		return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "玩家基础数据为空")
	}
	if msg == nil || msg.Base == nil {
		return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "玩家基础数据为空")
	}
	d.data = msg.Base
	return nil
}

func (d *PlayerBaseFun) Save(data *pb.PlayerData) error {
	if data == nil {
		return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "玩家数据为空")
	}
	buf, _ := proto.Marshal(d.data)
	newBase := &pb.PlayerDataBase{}
	proto.Unmarshal(buf, newBase)
	data.Base = newBase
	return nil
}

func (d *PlayerBaseFun) QueryPlayerData(head *pb.Head, req *pb.QueryPlayerDataReq, rsp *pb.QueryPlayerDataRsp) error {
	rsp.Data = d.data
	if d.data.RoomInfo != nil {
		// todo more 游戏类型生成器规则
		types := util.DestructRoomId(d.data.RoomInfo.RoomId)
		rsp.GameType = types.GetGameType()
		rsp.CoinType = types.GetCoinType()
		rsp.MatchType = types.GetMatchType()
	}
	return nil
}

// todo 换桌 打到match 选举一个新桌

func (d *PlayerBaseFun) LoadComplate() error {
	return nil
}

func (d *PlayerBaseFun) TexasJoinRoom(head *pb.Head, req *pb.TexasJoinRoomReq, rsp *pb.TexasJoinRoomRsp) error {
	gameType := pb.GameType((req.RoomId >> 32) & 0xFF)
	if d.data.RoomInfo == nil {
		req.PlayerInfo = d.data.PlayerInfo
		d.data.RoomInfo = &pb.PlayerRoomInfo{GameType: gameType, RoomId: req.RoomId}
		d.Change()
		return framework.Send(framework.SwapToRoom(head, req.RoomId, "TexasGameMgr", "JoinRoomReq"), req)
	}
	if d.data.RoomInfo.GameType != gameType {
		return uerror.NEW(pb.ErrorCode_GAME_PLAYER_IN_OTHER_GAME, head, "玩家已在其他游戏中,无法加入德州扑克房间")
	}
	if d.data.RoomInfo.RoomId > 0 && d.data.RoomInfo.RoomId != req.RoomId {
		matchType, gameType, coinType := pb.MatchType((req.RoomId>>40)&0xFF), pb.GameType((req.RoomId>>32)&0xFF), pb.CoinType((req.RoomId>>24)&0xFF)
		newrsp := &pb.HasRoomRsp{}
		newreq := &pb.HasRoomReq{RoomId: d.data.RoomInfo.RoomId}
		var dst *pb.NodeRouter
		switch matchType {
		case pb.MatchType_MatchTypeNone:
			dst = framework.NewMatchRouter(uint64(matchType)<<32|uint64(gameType)<<16|uint64(coinType), "MatchTexasRoom", "HasRoomReq")
		case pb.MatchType_MatchTypeSNG:
			dst = framework.NewMatchRouter(uint64(pb.DataType_DataTypeSngRoom), "SngRoomMgr", "HasRoomReq")
		}
		if err := framework.Request(&pb.Head{Dst: dst}, newreq, newrsp); err != nil {
			return err
		}
		if newrsp.IsExist {
			req.RoomId = d.data.RoomInfo.RoomId
		} else {
			d.data.RoomInfo.RoomId = req.RoomId
			d.Change()
		}
	}
	req.PlayerInfo = d.data.PlayerInfo
	return framework.Send(framework.SwapToRoom(head, req.RoomId, "TexasGameMgr", "JoinRoomReq"), req)
}

func (d *PlayerBaseFun) TexasQuitRoom(head *pb.Head, req *pb.TexasQuitRoomReq, rsp *pb.TexasQuitRoomRsp) error {
	gameType := pb.GameType((req.RoomId >> 32) & 0xFF)
	if d.data.RoomInfo == nil {
		return uerror.NEW(pb.ErrorCode_GAME_PLAYER_NOT_IN_ROOM, head, "玩家不在房间中")
	}

	if d.data.RoomInfo.GameType != gameType {
		return uerror.NEW(pb.ErrorCode_PARAM_INVALID, head, "玩家已在其他游戏中,无法退出德州扑克房间")
	}
	d.data.RoomInfo = nil
	d.Change()
	return nil
	//return framework.Send(framework.SwapToRoom(head, req.RoomId, "TexasGameMgr", "QuitRoomReq"), req)
}

// RummyJoinRoomReq
func (d *PlayerBaseFun) RummyJoinRoom(head *pb.Head, req *pb.RummyJoinRoomReq, rsp *pb.RummyJoinRoomRsp) error {
	gameType := pb.GameType((req.RoomId >> 32) & 0xFF)
	if d.data.RoomInfo == nil {
		d.data.RoomInfo = &pb.PlayerRoomInfo{GameType: gameType, RoomId: req.RoomId}
		d.Change()
		return framework.Send(framework.SwapToRoom(head, req.RoomId, "RummyGameMgr", "JoinRoomReq"), req)
	}

	if d.data.RoomInfo.GameType != gameType {
		return uerror.NEW(pb.ErrorCode_GAME_PLAYER_IN_OTHER_GAME, head, "玩家已在其他游戏中,无法加入Rummy房间")
	}

	if d.data.RoomInfo.RoomId > 0 && d.data.RoomInfo.RoomId != req.RoomId {
		req.RoomId = d.data.RoomInfo.RoomId
		d.Change()
	}
	req.PlayerInfo = d.data.PlayerInfo
	return framework.Send(framework.SwapToRoom(head, req.RoomId, "RummyGameMgr", "JoinRoomReq"), req)
}

func (d *PlayerBaseFun) RummyQuitRoom(head *pb.Head, req *pb.RummyQuitRoomReq, rsp *pb.RummyQuitRoomRsp) error {
	gameType := pb.GameType((req.RoomId >> 32) & 0xFF)
	if d.data.RoomInfo == nil {
		return uerror.NEW(pb.ErrorCode_GAME_PLAYER_NOT_IN_ROOM, head, "玩家不在房间中")
	}

	if d.data.RoomInfo.GameType != gameType {
		return uerror.NEW(pb.ErrorCode_PARAM_INVALID, head, "玩家已在其他游戏中,无法退出Rummy房间")
	}
	d.data.RoomInfo = nil
	d.Change()
	return framework.Send(framework.SwapToRoom(head, req.RoomId, "RummyGame", "QuitRoomReq"), req)
}
