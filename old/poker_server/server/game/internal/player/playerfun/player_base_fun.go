package playerfun

import (
	"poker_server/common/pb"
	"poker_server/common/room_util"
	"poker_server/framework"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"poker_server/library/util"
	"poker_server/server/game/internal/player/domain"
	"poker_server/server/game/internal/player/request"

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
		types := util.DestructRoomId(d.data.RoomInfo.RoomId)
		rsp.GameType = types.GetGameType()
		rsp.CoinType = types.GetCoinType()
		rsp.MatchType = types.GetMatchType()
	}
	return nil
}

func (d *PlayerBaseFun) LoadComplate() error {
	return nil
}

func (d *PlayerBaseFun) SetPlayerInfo(info *pb.PlayerInfo) {
	if info == nil {
		return
	}
	d.data.PlayerInfo = info
}

func (d *PlayerBaseFun) GetPlayerInfo() *pb.PlayerInfo {
	return d.data.PlayerInfo
}

func (d *PlayerBaseFun) GetRoomInfo() *pb.PlayerRoomInfo {
	return d.data.RoomInfo
}

func (d *PlayerBaseFun) ResetRoomInfo() {
	d.data.RoomInfo = nil
	d.Change()
}

func (d *PlayerBaseFun) GetRummyRealRoomId(head *pb.Head, roomId uint64) (uint64, error) {
	types := util.DestructRoomId(roomId)
	if d.data.RoomInfo == nil || types.GetGameType() == d.data.RoomInfo.GameType && d.data.RoomInfo.RoomId == roomId {
		return roomId, nil
	}

	if d.data.RoomInfo.GameType != types.GetGameType() {
		return 0, uerror.New(1, pb.ErrorCode_GAME_PLAYER_IN_OTHER_GAME, "玩家已在%v,无法加入Rummy房间", d.data.RoomInfo)
	}

	// 断线重连
	newHead := &pb.Head{
		Uid: head.Uid,
		Src: framework.NewSrcRouter(head.Uid, "Player"),
		Dst: framework.NewMatchRouter(uint64(types.GetGameType())<<32|uint64(types.GetCoinType()), "MatchRummyRoom", "Query"),
	}
	newRsp := &pb.HasRoomRsp{}
	newReq := &pb.HasRoomReq{RoomId: d.data.RoomInfo.RoomId}
	if err := cluster.Request(newHead, roomId, newReq); err != nil {
		return 0, err
	}
	if newRsp.IsExist {
		return d.data.RoomInfo.RoomId, nil
	}
	return roomId, nil
}

func (d *PlayerBaseFun) RummyJoinRoom(roomId uint64) {
	types := util.DestructRoomId(roomId)
	d.data.RoomInfo = &pb.PlayerRoomInfo{GameType: types.GetGameType(), RoomId: roomId}
	d.Change()
	return
}

func (d *PlayerBaseFun) RummyQuitRoom() {
	d.data.RoomInfo = nil
	d.Change()
	return
}

func (d *PlayerBaseFun) Close(uid uint64) {
	roomInfo := d.GetRoomInfo()
	if roomInfo == nil || roomInfo.RoomId <= 0 {
		return
	}

	switch roomInfo.GameType {
	case pb.GameType_GameTypeNormal:
		head := &pb.Head{
			Uid: uid,
			Src: framework.NewSrcRouter(uid, "Player"),
			Dst: framework.NewRoomRouter(roomInfo.RoomId, "TexasGameMgr", "QuitRoomReq"),
		}
		req := &pb.TexasQuitRoomReq{RoomId: roomInfo.RoomId}
		rsp := &pb.TexasQuitRoomRsp{}
		if err := d.TexasQuitRoom(head, req, rsp); err != nil {
			mlog.Errorf("玩家关闭错误 %v", err)
		}
	}

	d.GetBagFunc().ChargeTransOut(&pb.Head{Uid: uid}, roomInfo.RoomId)
}

// ----------------------------德州扑克--------------------
// 加入房间请求
func (d *PlayerBaseFun) TexasJoinRoom(head *pb.Head, req *pb.TexasJoinRoomReq, rsp *pb.TexasJoinRoomRsp, isMatch bool) error {
	// 断线重连
	if roomInfo := d.GetRoomInfo(); roomInfo != nil && req.RoomId > 0 {
		matchType, gameType, coinType := room_util.TexasRoomIdTo(req.RoomId)
		mType, gType, cType := room_util.TexasRoomIdTo(roomInfo.RoomId)
		if gType != gameType || cType != coinType || mType != matchType {
			return uerror.New(1, pb.ErrorCode_GAME_PLAYER_IN_OTHER_GAME, "玩家不在德州游戏中%v", roomInfo)
		}
		joinReq := &pb.TexasJoinRoomReq{
			RoomId:     roomInfo.RoomId,
			TableId:    roomInfo.TableId,
			MatchType:  matchType,
			PlayerInfo: d.GetPlayerInfo(),
		}
		if err := request.TexasJoinReq(head, joinReq, rsp); err != nil {
			mlog.Errorf("断线重连失败, 房间不存在：%v", err)
			d.data.RoomInfo = nil
		} else {
			req.RoomId = roomInfo.RoomId
			req.MatchType = mType
			req.CoinType = cType
			return nil
		}
	}
	// 匹配房间
	if isMatch && req.TableId > 0 {
		matchReq := &pb.TexasMatchRoomReq{TableId: req.TableId}
		matchRsp := &pb.TexasMatchRoomRsp{}
		if err := request.TexasMatchRoomReq(head, matchReq, matchRsp); err != nil {
			return err
		}
		req.RoomId = matchRsp.RoomId
		req.TableId = matchReq.TableId
	}
	// 加入房间
	if req.RoomId > 0 {
		matchType, gameType, coinType := room_util.TexasRoomIdTo(req.RoomId)
		req.MatchType = matchType
		req.CoinType = coinType
		req.PlayerInfo = d.GetPlayerInfo()
		if err := request.TexasJoinReq(head, req, rsp); err != nil {
			return err
		}
		if matchType == pb.MatchType_MatchTypeNone {
			d.GetBagFunc().SubProp(req.CoinType, req.BuyInChips)
		}
		d.data.RoomInfo = &pb.PlayerRoomInfo{GameType: gameType, RoomId: req.RoomId, TableId: req.TableId}
		d.Change()
	}
	return nil
}

// 退出房间请求
func (d *PlayerBaseFun) TexasQuitRoom(head *pb.Head, req *pb.TexasQuitRoomReq, rsp *pb.TexasQuitRoomRsp) error {
	if err := request.TexasQuitReq(head, req, rsp); err != nil {
		return err
	}
	matchType, _, _ := room_util.TexasRoomIdTo(req.RoomId)
	if matchType == pb.MatchType_MatchTypeNone {
		d.GetBagFunc().AddProp(rsp.CoinType, rsp.Chip)
	}
	roomInfo := d.GetRoomInfo()
	if roomInfo != nil && roomInfo.RoomId == req.RoomId {
		d.data.RoomInfo = nil
		d.Change()
	}
	return nil
}

// 德州换房间
func (d *PlayerBaseFun) TexasChangeRoom(head *pb.Head, req *pb.TexasChangeRoomReq, rsp *pb.TexasChangeRoomRsp) error {
	roomInfo := d.GetRoomInfo()
	if roomInfo == nil || roomInfo.RoomId != req.RoomId {
		return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "参数错误%v", req)
	}
	// 匹配房间
	matchReq := &pb.TexasMatchRoomReq{TableId: roomInfo.TableId, RoomId: roomInfo.RoomId}
	matchRsp := &pb.TexasMatchRoomRsp{}
	if err := request.TexasMatchRoomReq(head, matchReq, matchRsp); err != nil {
		return err
	}

	// 退出房间
	quitReq := &pb.TexasQuitRoomReq{RoomId: req.RoomId}
	quitRsp := &pb.TexasQuitRoomRsp{}
	if err := d.TexasQuitRoom(head, quitReq, quitRsp); err != nil {
		return err
	}

	// 加入房间
	joinReq := &pb.TexasJoinRoomReq{
		RoomId:     matchRsp.RoomId,
		TableId:    matchRsp.TableId,
		PlayerInfo: d.GetPlayerInfo(),
		CoinType:   quitRsp.CoinType,
		BuyInChips: quitRsp.Chip,
	}
	joinRsp := &pb.TexasJoinRoomRsp{}
	if err := d.TexasJoinRoom(head, joinReq, joinRsp, false); err != nil {
		return err
	}
	rsp.Duration = joinRsp.Duration
	rsp.TableInfo = joinRsp.TableInfo
	rsp.PlayerInfo = joinRsp.PlayerInfo
	rsp.RoomInfo = joinRsp.RoomInfo
	return nil
}

func (d *PlayerBaseFun) MatchStart(timeout int64, gameType pb.GameType, coinType pb.CoinType) {
	d.data.MatchInfo = &pb.PlayerMatchInfo{
		TimeOut:  timeout,
		GameType: gameType,
		CoinType: coinType,
	}
	d.Change()
	return
}

func (d *PlayerBaseFun) MatchStop() {
	d.data.MatchInfo = nil
	d.Change()
	return
}

func (d *PlayerBaseFun) GetMatchInfo() *pb.PlayerMatchInfo {
	return d.data.MatchInfo
}
