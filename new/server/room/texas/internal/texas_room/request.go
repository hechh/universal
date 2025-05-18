package texas_room

import (
	"poker_server/common/config/repository/machine_config"
	"poker_server/common/config/repository/texas_config"
	"poker_server/common/dao/repository/redis/room_router_data"
	"poker_server/common/dao/repository/redis/texas_room_data"
	"poker_server/common/dao/repository/redis/texas_room_player"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/library/mlog"
	"poker_server/server/room/texas/internal/base"
	"poker_server/server/room/texas/internal/machine"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
)

func SendPackageByPlayerIDRummy(uid uint64, cmd pb.CMD, data proto.Message) error {
	if data == nil {
		mlog.Errorf("数据为空, uid:%d, cmd:%d", uid, cmd)
		return nil
	}
	head := &pb.Head{
		Cmd: uint32(cmd),
		Id:  uid,
	}
	return framework.SendMsgToClient(head, data)
}

func BroadcastPackageByPlayerIDRummy(uids []uint64, cmd pb.CMD, data proto.Message) error {
	if data == nil {
		return nil
	}
	buf, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	head := &pb.Head{Cmd: uint32(cmd)}
	for _, uid := range uids {
		head.Id = uid
		framework.SendToClient(head, buf)
	}
	return nil
}

// 加载房间数据
func (d *TexasRoom) LoadRoomReq(uid uint64, req *pb.TexasJoinRoomReq, rsp *pb.TexasJoinRoomRsp) {
	// 从redis加载
	data, err := texas_room_data.Get(req.RoomId)
	if err != nil {
		atomic.AddInt32(&d.isFinished, 1)
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_REDIS_LOAD_ROOM, err.Error())
		SendPackageByPlayerIDRummy(uid, pb.CMD_TEXAS_JOIN_ROOM_RSP, rsp)
		mlog.Errorf("redis加载房间失败, uid:%d, roomId:%d, err:%s", uid, req.RoomId, err.Error())
		return
	}
	if data == nil {
		atomic.AddInt32(&d.isFinished, 1)
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_ROOM_NOT_FOUND, "房间不存在")
		SendPackageByPlayerIDRummy(uid, pb.CMD_TEXAS_JOIN_ROOM_RSP, rsp)
		mlog.Errorf("房间不存在, uid:%d, roomId:%d", uid, req.RoomId)
		return
	}
	d.data = data

	if d.data.Table.ChairInfo == nil {
		d.data.Table.ChairInfo = make(map[uint32]uint64)
	}
	if d.data.Table.Players == nil {
		d.data.Table.Players = make(map[uint64]*pb.TexasPlayerData)
	}
	if d.data.Table.GameData.PotPool == nil {
		d.data.Table.GameData.PotPool = &pb.TexasPotPoolData{}
	}

	// 加载配置
	d.texasCfg = texas_config.MGetRoomStageCoinType(int32(data.BaseInfo.RoomStage), int32(data.BaseInfo.CoinType))
	if d.texasCfg == nil {
		atomic.AddInt32(&d.isFinished, 1)
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_TEXAS_ROOM_CFG_NOT_FOUND, "房间配置不存在")
		SendPackageByPlayerIDRummy(uid, pb.CMD_TEXAS_JOIN_ROOM_RSP, rsp)
		mlog.Errorf("房间配置不存在, uid:%d, roomId:%d", uid, req.RoomId)
		return
	}

	d.machineCfg = machine_config.MGetGameId(data.BaseInfo.GameType)
	if d.machineCfg == nil {
		atomic.AddInt32(&d.isFinished, 1)
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_TEXAS_MACHINE_CFG_NOT_FOUND, "状态机配置不存在")
		SendPackageByPlayerIDRummy(uid, pb.CMD_TEXAS_JOIN_ROOM_RSP, rsp)
		mlog.Errorf("状态机配置不存在, uid:%d, roomId:%d", uid, req.RoomId)
		return
	}

	// 设置玩家所在房间
	if err := texas_room_player.HSet(uid, req.RoomId); err != nil {
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_REDIS_SAVE, err.Error())
		SendPackageByPlayerIDRummy(uid, pb.CMD_TEXAS_JOIN_ROOM_RSP, rsp)
		mlog.Errorf("redis保存失败, uid:%d, roomId:%d, err:%s", uid, req.RoomId, err.Error())
		return
	}

	d.JoinRoomReq(uid, req, rsp)
}

// 加入房间
func (d *TexasRoom) JoinRoomReq(uid uint64, req *pb.TexasJoinRoomReq, rsp *pb.TexasJoinRoomRsp) {
	if d.data == nil || d.data.Table == nil {
		mlog.Errorf("房间数据未加载")
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_ROOM_NOT_LOAD, "房间数据未加载")
		d.SendToClient(uid, pb.CMD_TEXAS_JOIN_ROOM_RSP, rsp)
		return
	}

	// 判断玩家是否已经加入房间
	if _, ok := d.data.Table.Players[uid]; ok {
		mlog.Debugf("断线重连玩家重新加入房间")
	} else {
		usr := d.GetOrNewPlayer(uid)
		base.SetPlayerState(usr, pb.TexasPlayerState_TPS_JOIN_ROOM)
		d.Change()
	}

	now := time.Now().UnixMilli()
	if d.machine == nil {
		d.machine = machine.NewTexasMachine(now, d.data.Table.CurState, d)
	}
	d.isStart = 1

	// 设置数据
	rsp.RoomInfo = (d.data.BaseInfo)
	rsp.RoomInfo.RoomId = (d.data.RoomId)
	rsp.TableInfo = (d.data.Table)
	rsp.Duration = (d.GetStateStartTime() + d.GetCurStateTTL() - now)
	mlog.Infof("玩家加入房间成功, uid:%d, roomId:%d, rsp:%v", uid, d.data.RoomId, rsp)
	d.SendToClient(uid, pb.CMD_TEXAS_JOIN_ROOM_RSP, rsp)
}

// 买入请求
func (d *TexasRoom) BuyInReq(uid uint64, req *pb.TexasBuyInReq, rsp *pb.TexasBuyInRsp) {
	usr, ok := d.data.Table.Players[uid]
	if !ok {
		mlog.Errorf("玩家未加入房间")
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_PLAYER_NOT_IN_ROOM, "玩家未加入房间")
		d.SendToClient(uid, pb.CMD_TEXAS_BUY_IN_RSP, rsp)
		return
	}
	usr.Chips += req.Chip
	d.Change()

	rsp.RoomId = (d.data.RoomId)
	rsp.CoinType = (int32(d.data.BaseInfo.CoinType))
	rsp.Chip = (usr.Chips)
	d.SendToClient(uid, pb.CMD_TEXAS_BUY_IN_RSP, rsp)
}

// 坐下请求
func (d *TexasRoom) SitDownReq(uid uint64, req *pb.TexasSitDownReq, rsp *pb.TexasSitDownRsp) {
	usr, ok := d.data.Table.Players[uid]
	if !ok || usr == nil {
		mlog.Errorf("玩家未加入房间")
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_PLAYER_NOT_IN_ROOM, "玩家未加入房间")
		d.SendToClient(uid, pb.CMD_TEXAS_SIT_DOWN_RSP, rsp)
		return
	}

	// 最大桌位数
	if len(d.data.Table.ChairInfo) >= int(d.data.BaseInfo.MaxPlayerCount) {
		mlog.Errorf("牌桌已满")
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_TABLE_FULL, "牌桌已满")
		d.SendToClient(uid, pb.CMD_TEXAS_SIT_DOWN_RSP, rsp)
		return
	}

	// 筹码不足
	if usr.Chips < d.data.BaseInfo.BigBlind {
		mlog.Errorf("筹码不足")
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_CHIPS_NOT_ENOUGH, "筹码不足")
		d.SendToClient(uid, pb.CMD_TEXAS_SIT_DOWN_RSP, rsp)
		return
	}

	// 判断位置是否被占用
	if cur, ok := d.data.Table.ChairInfo[req.ChairId]; ok {
		if uid == cur {
			mlog.Errorf("玩家已在该位置, uid: %d, req: %v", uid, req)
			rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_PLAYER_IN_CHAIR, "玩家已在该位置")
			d.SendToClient(uid, pb.CMD_TEXAS_SIT_DOWN_RSP, rsp)
			return
		}
		mlog.Errorf("位置已被占用, uid:%d, req:%v", uid, req)
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_CHAIR_OCCUPIED, "位置已被占用")
		d.SendToClient(uid, pb.CMD_TEXAS_SIT_DOWN_RSP, rsp)
		return
	}

	// 加入牌桌
	d.data.Table.ChairInfo[req.ChairId] = uid
	usr.ChairId = req.ChairId
	base.SetPlayerState(usr, pb.TexasPlayerState_TPS_JOIN_TABLE)
	d.SendToClient(uid, pb.CMD_TEXAS_SIT_DOWN_RSP, rsp)
	d.Change()

	// 广播消息
	d.NotifyToClient(pb.TexasEventType_EVENT_SIT_DOWN, &pb.TexasPlayerEventNotify{
		RoomId:  (d.data.RoomId),
		ChairId: req.ChairId,
		Player:  (usr),
	})
}

// 站起请求
func (d *TexasRoom) StandUpReq(uid uint64, req *pb.TexasStandUpReq, rsp *pb.TexasStandUpRsp) {
	usr, ok := d.data.Table.Players[uid]
	if !ok || usr == nil {
		mlog.Errorf("玩家未加入房间")
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_PLAYER_NOT_IN_ROOM, "玩家未加入房间")
		d.SendToClient(uid, pb.CMD_TEXAS_STAND_UP_RSP, rsp)
		return
	}

	// 请求参数错误
	if usr.ChairId > 0 && usr.ChairId != req.ChairId {
		mlog.Errorf("请求参数错误")
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_REQUEST_PARAM_ERROR, "请求参数错误")
		d.SendToClient(uid, pb.CMD_TEXAS_STAND_UP_RSP, rsp)
		return
	}

	// 玩家已经站起
	if base.IsPlayerState(usr, pb.TexasPlayerState_TPS_QUIT_TABLE) {
		mlog.Errorf("玩家已经站起")
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_PLAYER_HAS_STAND_UP, "玩家已经站起")
		d.SendToClient(uid, pb.CMD_TEXAS_STAND_UP_RSP, rsp)
		return
	}

	if !usr.GameInfo.InPlaying || d.data.Table.CurState == pb.TexasGameState_TGS_INIT || d.data.Table.CurState == pb.TexasGameState_TGS_END {
		usr.ChairId = 0
		delete(d.data.Table.ChairInfo, req.ChairId)
		d.NotifyToClient(pb.TexasEventType_EVENT_STAND_UP, &pb.TexasPlayerEventNotify{
			RoomId:  (d.data.RoomId),
			ChairId: req.ChairId,
		})
	} else {
		base.SetPlayerState(usr, pb.TexasPlayerState_TPS_QUIT_TABLE)
	}
	d.Change()
}

// 下注请求
func (d *TexasRoom) DoBetReq(uid uint64, req *pb.TexasDoBetReq, rsp *pb.TexasDoBetRsp) {
	usr := d.GetCursor()
	if usr == nil {
		mlog.Errorf("玩家未加入房间")
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_PLAYER_NOT_IN_ROOM, "玩家未加入房间")
		d.SendToClient(uid, pb.CMD_TEXAS_DO_BET_RSP, rsp)
		return
	}

	// 判断请求是否合法
	if usr.Uid != uid || usr.ChairId != req.ChairId {
		mlog.Errorf("请求参数错误")
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_REQUEST_PARAM_ERROR, "请求参数错误")
		d.SendToClient(uid, pb.CMD_TEXAS_DO_BET_RSP, rsp)
		return
	}

	// 玩家已经弃牌或者梭哈了
	if usr.GameInfo.Operate == pb.TexasOperateType_TOT_FOLD || usr.GameInfo.Operate == pb.TexasOperateType_TOT_ALL_IN {
		mlog.Errorf("玩家已经弃牌或者梭哈了")
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_PLAYER_HAS_FOLD_OR_ALLIN, "玩家已经弃牌或者梭哈了")
		d.SendToClient(uid, pb.CMD_TEXAS_DO_BET_RSP, rsp)
		return
	}

	// 判断筹码是否足够
	if req.OperateType != int32(pb.TexasOperateType_TOT_CHECK) && req.OperateType != int32(pb.TexasOperateType_TOT_FOLD) {
		if usr.Chips < req.Chip || usr.Chips == req.Chip {
			req.OperateType = int32(pb.TexasOperateType_TOT_ALL_IN)
			req.Chip = usr.Chips
		} else if req.Chip > d.data.Table.GameData.MaxBetChips-usr.GameInfo.BetChips {
			req.OperateType = (int32(pb.TexasOperateType_TOT_RAISE))
		} else if req.Chip == d.data.Table.GameData.MaxBetChips-usr.GameInfo.BetChips {
			req.OperateType = (int32(pb.TexasOperateType_TOT_CALL))
		}
	}

	switch pb.TexasOperateType(req.OperateType) {
	case pb.TexasOperateType_TOT_CALL:
		if req.Chip != d.data.Table.GameData.MaxBetChips-usr.GameInfo.BetChips {
			mlog.Errorf("请求参数错误")
			rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_OPERATE_NOT_SUPPORTED, "下注操筹码不足")
			d.SendToClient(uid, pb.CMD_TEXAS_DO_BET_RSP, rsp)
			return
		}
	case pb.TexasOperateType_TOT_RAISE:
		if req.Chip < usr.Chips {
			if req.Chip < d.data.Table.GameData.MinRaise+d.data.Table.GameData.MaxBetChips-usr.GameInfo.BetChips {
				mlog.Errorf("筹码不足: req: %v", req)
				rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_CHIPS_NOT_ENOUGH, "下注操筹码不足")
				d.SendToClient(uid, pb.CMD_TEXAS_DO_BET_RSP, rsp)
				return
			}
			d.data.Table.GameData.MinRaise = req.Chip - d.data.Table.GameData.MaxBetChips + usr.GameInfo.BetChips
		} else {
			minRaise := req.Chip - d.data.Table.GameData.MaxBetChips + usr.GameInfo.BetChips
			if minRaise > d.data.Table.GameData.MinRaise {
				d.data.Table.GameData.MinRaise = minRaise
			}
		}
	}

	usr.Chips -= req.Chip

	// 下注操作
	d.Operate(usr, pb.TexasOperateType(req.OperateType), req.Chip)
	d.Change()

	// 返回客户端
	rsp.Round = (d.data.Table.Round)
	rsp.ChairId = req.ChairId
	rsp.OpType = req.OperateType
	rsp.Chip = req.Chip
	rsp.BankRoll = (usr.Chips)
	rsp.TotalBet = (usr.GameInfo.BetChips)
	rsp.RoomId = (d.data.RoomId)
	d.SendToClient(uid, pb.CMD_TEXAS_DO_BET_RSP, rsp)
}

// 获取房间信息
func (d *TexasRoom) RoomInfoReq(uid uint64, req *pb.TexasRoomInfoReq, rsp *pb.TexasRoomInfoRsp) {
	rsp.Base = (d.data.BaseInfo)
	rsp.TableInfo = (d.data.Table)
	d.SendToClient(uid, pb.CMD_TEXAS_ROOM_INFO_RSP, rsp)
}

// 离开房间请求
func (d *TexasRoom) QuitRoomReq(uid uint64, req *pb.TexasQuitRoomReq, rsp *pb.TexasQuitRoomRsp) {
	usr, ok := d.data.Table.Players[uid]
	if !ok || usr == nil {
		mlog.Errorf("玩家未加入房间, uid:%d, rqe:%v", uid, req)
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_PLAYER_NOT_IN_ROOM, "玩家未加入房间")
		d.SendToClient(uid, pb.CMD_TEXAS_QUIT_ROOM_RSP, rsp)
		return
	}

	// 玩家已经站起
	if usr.ChairId > 0 {
		mlog.Errorf("玩家未站起来, uid:%d, req:%v", uid, req)
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_PLAYER_IN_CHAIR, "玩家未站起来")
		d.SendToClient(uid, pb.CMD_TEXAS_QUIT_ROOM_RSP, rsp)
		return
	}

	// 删除路由
	if err := room_router_data.HDel(uid); err != nil {
		mlog.Errorf("玩家路由删除失败, uid:%d, req:%v", uid, req)
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_REDIS_SAVE, err.Error())
		d.SendToClient(uid, pb.CMD_TEXAS_QUIT_ROOM_RSP, rsp)
		return
	}
	// 删除房间信息
	if err := texas_room_player.HDel(uid); err != nil {
		mlog.Errorf("玩家断线重连信息删除失败, uid:%d, req:%v", uid, req)
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_REDIS_SAVE, err.Error())
		d.SendToClient(uid, pb.CMD_TEXAS_QUIT_ROOM_RSP, rsp)
		return
	}

	rsp.RoomId = (d.data.RoomId)
	d.SendToClient(uid, pb.CMD_TEXAS_QUIT_ROOM_RSP, rsp)
}
