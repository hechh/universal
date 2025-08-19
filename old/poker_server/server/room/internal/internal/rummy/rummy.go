package rummy

import (
	"poker_server/common/config/repository/rummy_config"
	"poker_server/common/config/repository/rummy_machine_config"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"poker_server/library/util"
	"poker_server/server/room/internal/module/card"
	"poker_server/server/room/internal/module/machine"
	"time"

	"github.com/golang/protobuf/proto"
)

// RummyGame 游戏内存定义 游戏数据分
type RummyGame struct {
	actor.Actor
	// 状态机
	Data       *pb.RummyRoomData
	machine    *machine.Machine //状态机
	isChange   bool             // 是否有数据变更
	IsFinish   bool             // 关闭协程
	finishTime int64            // 关闭时间
}

// NewRummyGame 初始化游戏对象
func NewRummyGame(data *pb.RummyRoomData) *RummyGame {
	if data.RoomCfg == nil {
		rummyCfg := rummy_config.MGetID(data.GameId)
		if rummyCfg == nil {
			mlog.Infof("rummy配置不存在: %d", data.GameId)
			return nil
		}
		data.RoomCfg = rummyCfg
	}

	if data.MachineCfg == nil {
		machineCfg := rummy_machine_config.MGetGameType(data.RoomCfg.GameType)
		if machineCfg == nil {
			mlog.Infof("rummy状态机配置不存在: %d", data.GameId)
			return nil
		}
		data.MachineCfg = machineCfg
	}

	if data.Common == nil {
		data.Common = &pb.RummyCommon{
			CreateTime: time.Now().UnixMilli(),
			PlayerIds:  make([]uint64, 0, data.RoomCfg.MaxPlayerCount),
			OprPlayer: &pb.RummyOprPlayer{
				Round: 1,
			},
		}
		data.Match = &pb.MatchParam{
			Match: 1,
		}
	}

	if data.Private == nil {
		data.Private = &pb.RummyPrivate{}
	}

	if data.Common.Players == nil {
		data.Common.Players = make(map[uint64]*pb.RummyRoomPlayer, data.RoomCfg.MaxPlayerCount)
	}

	if data.Common.Seats == nil {
		data.Common.Seats = make(map[uint32]uint64, data.RoomCfg.MaxPlayerCount)
	}

	if data.Common.EmptySeats == nil {
		data.Common.EmptySeats = make([]uint32, 0, data.RoomCfg.MaxPlayerCount)
		for i := uint32(1); i <= data.RoomCfg.MaxPlayerCount; i++ {
			data.Common.EmptySeats = append(data.Common.EmptySeats, i)
		}
	}

	ret := &RummyGame{Data: data}
	nowMs := time.Now().UnixMilli()

	if data.RoomCfg.GameType == pb.GameType_GameTypePR {
		ret.machine = machine.NewMachine(nowMs, pb.GameState_Rummy_STAGE_INIT, ret)
	} else {
		ret.machine = machine.NewMachine(nowMs, pb.GameState_RummyExt_STAGE_INIT, ret)
	}
	ret.Actor.Register(ret)
	ret.Actor.SetId(data.RoomId)
	ret.Start()

	// 启动定时器
	head := &pb.Head{SendType: pb.SendType_POINT, ActorName: "RummyGame", FuncName: "OnTick"}
	err := ret.RegisterTimer(head, 50*time.Millisecond, -1)
	if err != nil {
		mlog.Debug(head, "register timer err: %v", err)
	}

	return ret
}

func (d *RummyGame) GetMachine() *machine.Machine {
	return d.machine
}

func (d *RummyGame) Finish() {
	d.Change()
	d.IsFinish = true

	d.Walk(0, func(player *pb.RummyRoomPlayer) bool { // 15分钟检测房间没开始gc房间
		player.Health = pb.RummyPlayHealth_Rummy_KICK
		return true
	})
	nowMs := time.Now().UnixMilli()
	d.DelExitPlayer(nowMs)
}

func (d *RummyGame) Change() {
	d.isChange = true
}

func (d *RummyGame) Save() {
	if !d.isChange {
		return
	}

	cfg := rummy_config.MGetID(d.Data.GameId)
	head := &pb.Head{
		Src: framework.NewSrcRouter(d.GetRoomId(), d.GetActorName()),
		Dst: framework.NewMatchRouter(uint64(cfg.GameType)<<32|uint64(cfg.CoinType), "MatchRummyRoom", "Update"),
	}
	if err := cluster.Send(head, d.Data); err != nil {
		mlog.Errorf("Rummy房间数据保存失败: %v, error: %v", d.Data, err)
	} else {
		d.isChange = false
	}
}

func (d *RummyGame) OnTick() {
	nowMs := time.Now().UnixMilli()
	if d.machine == nil {
		types := util.DestructRoomId(d.GetRoomId())
		if types.GameType == pb.GameType_GameTypePR {
			d.machine = machine.NewMachine(nowMs, pb.GameState_Rummy_STAGE_INIT, d)
		} else { // ext
			d.machine = machine.NewMachine(nowMs, pb.GameState_RummyExt_STAGE_INIT, d)
		}
	} else {
		// 执行状态机
		d.machine.Handle(nowMs, d)
	}

	// 定时保存数据
	d.Save()
	if d.IsFinish && nowMs >= d.finishTime {
		actor.SendMsg(&pb.Head{ActorName: "RummyGameMgr", FuncName: "Remove"}, d.GetRoomId())
		d.finishTime = nowMs + 3000
	}
}

func (d *RummyGame) GetRoomId() uint64 {
	return d.Data.RoomId
}

func (d *RummyGame) RuntimeGC() {
	d.Data.Common.OprPlayer.Reset()
	d.Data.Common.OprPlayer.Round = 1
	d.Data.Private.Reset()
	//游戏过程记录
	d.Data.Common.GameFinish = false
	d.Data.Common.WinnerId = 0
}

func (d *RummyGame) Reset() {
	d.RuntimeGC()
	// reset game public data
	d.Data.Stage = pb.GameState_Rummy_STAGE_INIT

	// 回收match per数据
	d.Data.Match.Match++
	d.Data.Match.StartTime = 0
	d.Data.Match.EndTime = 0

	d.Walk(0, func(player *pb.RummyRoomPlayer) bool {
		player.Private.Reset()
		player.State = pb.RummyPlayState_Rummy_PLAYSTATE_PLAY
		player.Private.OutCards = player.Private.OutCards[:0]
		return true
	})
}

func (d *RummyGame) GetPlayerMap() map[uint64]*pb.RummyRoomPlayer {
	return d.Data.Common.Players
}

func (d *RummyGame) GetCurState() pb.GameState {
	return d.machine.GetCurState()
}

// GetOrNewPlayer 获取大厅玩家数据
func (d *RummyGame) GetOrNewPlayer(playerId uint64, req *pb.RummyJoinRoomReq) (ret *pb.RummyRoomPlayer) {
	now := time.Now().UnixMilli()

	ret = &pb.RummyRoomPlayer{
		PlayerNick:   req.PlayerInfo.NickName,
		PicUrl:       req.PlayerInfo.Avatar,
		PlayerId:     playerId,
		State:        pb.RummyPlayState_Rummy_PLAYSTATE_INIT,
		Health:       pb.RummyPlayHealth_Rummy_NORMAL,
		Coin:         uint64(max(req.Coin, 0)),
		Seat:         d.Data.Common.EmptySeats[0],
		ReadyTimeout: now + d.Data.MachineCfg.ReadyDuration*2000,
		Private: &pb.RummyPlayerRumtime{
			HandCards: make([]uint32, 0, 14),
			CardGroup: make([]*pb.RummyCardGroup, 0, 7),
			OutCards:  make([]uint32, 0, 16),
		},
		Total: d.Data.RoomCfg.MinBuyIn,
	}
	ret.JoinCoin = int64(ret.Coin)
	d.Data.Common.Seats[ret.Seat] = playerId
	d.Data.Common.EmptySeats = d.Data.Common.EmptySeats[1:]
	return
}

// UpReadyPlayer 非游戏状态 将就绪玩家转移至游玩队列
func (d *RummyGame) UpReadyPlayer() {
	d.Change()
	if int32(d.GetCurState()) < int32(pb.GameState_Rummy_STAGE_START) || (d.Data.RoomCfg.GameType != pb.GameType_GameTypePR && int32(d.GetCurState()) < int32(pb.GameState_RummyExt_STAGE_START)) {
		players := d.GetPlayerMap()
		newPlayerIds := make([]uint64, 0, len(players))
		rechargePlayerIds := make([]uint64, 0, len(players))
		for i := range players {
			player := players[i]
			if player.State == pb.RummyPlayState_Rummy_PLAYSTATE_READY {
				player.State = pb.RummyPlayState_Rummy_PLAYSTATE_PLAY
			}

			if int64(player.Coin) < d.Data.RoomCfg.MinBuyIn { // ban recharge users
				player.State = pb.RummyPlayState_Rummy_PLAYSTATE_READY
				rechargePlayerIds = append(rechargePlayerIds, player.PlayerId)
			}

			if player.State == pb.RummyPlayState_Rummy_PLAYSTATE_PLAY {
				newPlayerIds = append(newPlayerIds, player.PlayerId)
			}
		}

		if len(rechargePlayerIds) > 0 {
			// ntf
			ntf := &pb.RummyRechargeNtf{
				RoomId:   d.GetRoomId(),
				PlayerId: rechargePlayerIds,
			}
			err := d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_Recharge, ntf)
			mlog.Debugf("recharge PlayerIds NotifyToClient: %v", err)
		}
		// 就绪玩家 快速入桌
		d.Data.Common.PlayerIds = newPlayerIds
	}

	mlog.Infof("%v 当前可游戏玩家总数: %v", d.GetPlayerCount(), d.Data.Common.Players)
}

// GetPlayerUidList 获取所有人id包括观战
func (d *RummyGame) GetPlayerUidList() (uids []uint64) {
	for i := range d.Data.Common.Players {
		uids = append(uids, d.Data.Common.Players[i].PlayerId)
	}
	return
}

func (d *RummyGame) NewRummyEventNotify(notifyId pb.RummyEventType, msg proto.Message) *pb.RummyEventNotify {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return nil
	}
	return &pb.RummyEventNotify{Event: notifyId, Content: buf}
}

// GetPlayerCount 获取就绪玩家
func (d *RummyGame) GetPlayerCount() int {
	return len(d.Data.Common.PlayerIds)
}

func (d *RummyGame) GetNextState() pb.GameState {
	if d.Data.RoomCfg.GameType == pb.GameType_GameTypePR {
		switch d.GetCurState() {
		case pb.GameState_Rummy_STAGE_INIT:
			return pb.GameState_Rummy_STAGE_READY_START
		case pb.GameState_Rummy_STAGE_READY_START:
			return pb.GameState_Rummy_STAGE_START
		case pb.GameState_Rummy_STAGE_START:
			return pb.GameState_Rummy_STAGE_ZHUANG
		case pb.GameState_Rummy_STAGE_ZHUANG:
			return pb.GameState_Rummy_STAGE_DEAL
		case pb.GameState_Rummy_STAGE_DEAL:
			return pb.GameState_Rummy_STAGE_PLAYING
		case pb.GameState_Rummy_STAGE_PLAYING:
			return pb.GameState_Rummy_STAGE_FIX_CARD
		case pb.GameState_Rummy_STAGE_SHUFFLE:
			return pb.GameState_Rummy_STAGE_PLAYING
		case pb.GameState_Rummy_STAGE_FIX_CARD:
			return pb.GameState_Rummy_STAGE_SETTLE
		case pb.GameState_Rummy_STAGE_SETTLE:
			return pb.GameState_Rummy_STAGE_INIT
		default:
			return pb.GameState_Rummy_STAGE_INIT
		}
	} else {
		switch d.GetCurState() {
		case pb.GameState_RummyExt_STAGE_INIT:
			return pb.GameState_RummyExt_STAGE_READY_START
		case pb.GameState_RummyExt_STAGE_READY_START:
			return pb.GameState_RummyExt_STAGE_START
		case pb.GameState_RummyExt_STAGE_START:
			return pb.GameState_RummyExt_STAGE_ZHUANG
		case pb.GameState_RummyExt_STAGE_ZHUANG:
			return pb.GameState_RummyExt_STAGE_DEAL
		case pb.GameState_RummyExt_STAGE_DEAL:
			return pb.GameState_RummyExt_STAGE_PLAYING
		case pb.GameState_RummyExt_STAGE_PLAYING:
			return pb.GameState_RummyExt_STAGE_FIX_CARD
		case pb.GameState_RummyExt_STAGE_SHUFFLE:
			return pb.GameState_RummyExt_STAGE_PLAYING
		case pb.GameState_RummyExt_STAGE_FIX_CARD:
			return pb.GameState_RummyExt_STAGE_SETTLE
		case pb.GameState_RummyExt_STAGE_SETTLE:
			return pb.GameState_RummyExt_STAGE_INIT
		case pb.GameState_RummyExt_STAGE_FIN_SETTLE:
			return pb.GameState_RummyExt_STAGE_INIT
		default:
			return pb.GameState_RummyExt_STAGE_INIT
		}
	}
}

func (d *RummyGame) GetMinStartPlayers() int {
	return int(d.Data.GetRoomCfg().MinPlayerCount)
}

func (d *RummyGame) DelExitPlayer(nowMs int64) {
	dels := make([]uint64, 0, d.Data.RoomCfg.MaxPlayerCount)
	for playerKey := range d.Data.Common.Players {
		// 玩家已删除 或者 玩家长时间不准备
		if d.Data.Common.Players[playerKey].Health != pb.RummyPlayHealth_Rummy_NORMAL || (d.Data.Common.Players[playerKey].State == pb.RummyPlayState_Rummy_PLAYSTATE_INIT && d.Data.Common.Players[playerKey].ReadyTimeout <= nowMs) {
			mlog.Infof("不健康玩家:%v", playerKey)
			dels = append(dels, playerKey)
		}
	}

	for i := 0; i < len(dels); i++ {
		d.delTargetPlayer(dels[i])
	}
}

func (d *RummyGame) delTargetPlayer(playerID uint64) {
	d.Change()
	playerMap := d.GetPlayerMap()
	player, ok := playerMap[playerID]

	if ok && player.Health != pb.RummyPlayHealth_Rummy_QUIT { // 主动退出 清桌不提示消息
		ntf := &pb.RummyQuitRoomNtf{
			RoomId:    d.GetRoomId(),
			PlayerId:  playerID,
			LeaveType: pb.RummyLeaveType_Rummy_LEAVE_TYPE_KICK,
		}
		err := d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_RummyQuitRoom, ntf)
		mlog.Debugf("delTargetPlayer NotifyToClient: %v", err)
		d.Data.Common.EmptySeats = append(d.Data.Common.EmptySeats, player.Seat)

		// 还原道具到背包
		kickReq := &pb.RummyKickPlayerReq{
			RoomId: d.GetRoomId(),
			Charge: int64(player.Coin),
		}

		// 被动移出房间 还原背包货币 并退回大厅
		newHead := &pb.Head{
			Uid: playerID,
			Src: framework.NewActorRouter(d),
			Dst: framework.NewGameRouter(player.PlayerId, "PlayerMgr", "RummyKickPlayerReq"),
		}
		if err = cluster.Send(newHead, kickReq); err != nil { //异步补余额
			mlog.Infof("RummyGame RewardReq Error: %v", err)
		}
	}
	delete(playerMap, playerID)
	delete(d.Data.Common.Seats, player.Seat)
	newPlayerIds := make([]uint64, 0, d.Data.GetRoomCfg().MaxPlayerCount)
	for _, item := range d.Data.Common.PlayerIds {
		if item != playerID {
			newPlayerIds = append(newPlayerIds, item)
		}
	}
	d.Data.Common.PlayerIds = newPlayerIds
	return
}

func (d *RummyGame) FlushExpireTime(nowMs int64) {
	d.Data.Stage = d.GetCurState()
	d.Data.Common.TimeOut = nowMs + d.GetCurStateTTL()
	d.Data.Common.TotalTime = d.GetCurStateTTL()
	d.machine.SetCurStateStartTime(nowMs)
}

func (d *RummyGame) GetCurStateTTL() int64 {

	switch d.GetCurState() {
	case pb.GameState_Rummy_STAGE_READY_START:
		fallthrough
	case pb.GameState_RummyExt_STAGE_READY_START:
		return d.Data.MachineCfg.ReadyDuration * 1000
	case pb.GameState_Rummy_STAGE_START:
		fallthrough
	case pb.GameState_RummyExt_STAGE_START:
		return d.Data.MachineCfg.StartDuration * 1000
	case pb.GameState_Rummy_STAGE_ZHUANG:
		fallthrough
	case pb.GameState_RummyExt_STAGE_ZHUANG:
		return d.Data.MachineCfg.HostDuration * 1000
	case pb.GameState_Rummy_STAGE_DEAL:
		fallthrough
	case pb.GameState_RummyExt_STAGE_DEAL:
		return d.Data.MachineCfg.DealDuration * 1000
	case pb.GameState_Rummy_STAGE_PLAYING:
		fallthrough
	case pb.GameState_RummyExt_STAGE_PLAYING:
		return d.Data.MachineCfg.PlayDuration * 1000
	case pb.GameState_Rummy_STAGE_FIX_CARD:
		fallthrough
	case pb.GameState_RummyExt_STAGE_FIX_CARD:
		return d.Data.MachineCfg.FixDuration * 1000
	case pb.GameState_Rummy_STAGE_SETTLE:
		fallthrough
	case pb.GameState_RummyExt_STAGE_SETTLE:
		return d.Data.MachineCfg.SettleDuration * 1000
	case pb.GameState_Rummy_STAGE_SHUFFLE:
		fallthrough
	case pb.GameState_RummyExt_STAGE_SHUFFLE:
		return d.Data.MachineCfg.ShuffleDuration * 1000
	case pb.GameState_RummyExt_STAGE_FIN_SETTLE:
		return d.Data.MachineCfg.FinSettleDuration * 1000
	default:
		return d.Data.MachineCfg.ReadyDuration * 1000
	}
}

func (d *RummyGame) Walk(pos int, f func(player *pb.RummyRoomPlayer) bool) {
	size := len(d.Data.Common.PlayerIds)
	for i := pos; i < pos+size; i++ {
		if !f(d.Data.Common.Players[d.Data.Common.PlayerIds[i%size]]) {
			return
		}
	}
}

// SetPlayerLose 设置玩家判负
func (d *RummyGame) SetPlayerLose(playerID uint64, isFold, isFake bool) {
	player := d.Data.Common.Players[playerID]
	if isFold { //玩家弃牌
		player.State = pb.RummyPlayState_Rummy_PLAYSTATE_GIVEUP
		if isFake {
			player.Private.Score = 80 //炸胡
		} else {
			player.Private.Score = d.getGiveUpCard(player)
		}
		d.Data.Common.OprPlayer.ScorePool += player.Private.Score
		coin := d.Data.RoomCfg.BaseScore * player.Private.Score
		player.Coin -= uint64(coin)
		mlog.Infof("SetPlayerLose Fold PlayerId:%d Score:%d", playerID, player.Private.Score)
		return
	}
	// 修改奖池
	player.State = pb.RummyPlayState_Rummy_PLAYSTATE_LOSE
	_, Score := CheckRCG(player.Private.HandCards, player.Private.CardGroup)
	player.Private.Score = Score
	d.Data.Common.OprPlayer.ScorePool += Score
	mlog.Infof("SetPlayerLose PlayerId:%d Score:%d", playerID, player.Private.Score)
}

// NotifyToClient 组播 根据uid进行组播
func (d *RummyGame) NotifyToClient(uids []uint64, notifyId pb.RummyEventType, msg proto.Message) error {
	ntf := d.NewRummyEventNotify(notifyId, msg)

	newHead := &pb.Head{
		Src: framework.NewActorRouter(d),
		Cmd: uint32(pb.CMD_RUMMY_EVENT_NOTIFY),
	}
	return cluster.SendToClient(newHead, ntf, uids...)
}

func (d *RummyGame) CheckNormalPlayer() pb.GameState {
	normalCount := 0
	var normalPlayer *pb.RummyRoomPlayer
	d.Walk(0, func(player *pb.RummyRoomPlayer) bool { //可用玩家检查
		if player.State == pb.RummyPlayState_Rummy_PLAYSTATE_PLAY {
			normalCount++
			normalPlayer = player
		}
		return true
	})

	if normalCount <= 1 {
		if normalCount == 0 {
			d.Data.Common.GameFinish = true
			mlog.Infof("没有正常玩家 游戏结束 游戏信息： %v", d.Data)
			return pb.GameState_Rummy_STAGE_INIT
		}
		normalPlayer.State = pb.RummyPlayState_Rummy_PLAYSTATE_WIN
		d.Data.Common.WinnerId = normalPlayer.PlayerId
		mlog.Infof("只剩下小于一个正常玩家  赢家id: %v 回合信息: %v", d.Data.Common.WinnerId, d.Data.Common.OprPlayer)
		return pb.GameState_Rummy_STAGE_SETTLE
	}
	return d.GetCurState()
}

// SetNextPlayer 以庄家开始 打牌玩家轮转
func (d *RummyGame) SetNextPlayer() pb.GameState {
	d.Change()
	// 回收吃牌状态
	d.Data.Common.OprPlayer.IsEat = false
	mlog.Infof("当前回合 :%v", d.Data.Common.OprPlayer.Round)
	// 洗牌
	if d.Data.Private.CardIdx == uint32(len(d.Data.Private.Cards)) {
		return pb.GameState_Rummy_STAGE_SHUFFLE
	}

	size := uint32(len(d.Data.Common.PlayerIds))
	if d.Data.Common.OprPlayer.Round == 1 {
		for i := uint32(0); i < size; i++ { //根据定庄初始化第一个玩家
			if d.Data.Common.PlayerIds[i] == d.Data.Common.ZhuangId {
				d.Data.Common.OprPlayer.OprCurId = (i + 1) % size
				d.Data.Common.OprPlayer.CurOprPlayerId = d.Data.Common.PlayerIds[d.Data.Common.OprPlayer.OprCurId]
				d.Data.Common.OprPlayer.PrevOutUid = 0
				break
			}
		}
	} else {
		if state := d.CheckNormalPlayer(); state != d.GetCurState() {
			return state
		}
		d.Data.Common.OprPlayer.PrevOutUid = d.Data.Common.OprPlayer.CurOprPlayerId
		for {
			d.Data.Common.OprPlayer.OprCurId = (d.Data.Common.OprPlayer.OprCurId + 1) % size
			d.Data.Common.OprPlayer.CurOprPlayerId = d.Data.Common.PlayerIds[d.Data.Common.OprPlayer.OprCurId]
			curPlayer := d.GetCurPlayer()
			if curPlayer != nil && curPlayer.State == pb.RummyPlayState_Rummy_PLAYSTATE_PLAY {
				//用户有数据 没投降
				break
			}
		}

		mlog.Infof("出牌堆： %v 回合信息:%v", card.CardList(d.Data.Common.OprPlayer.OutCards), d.Data.Common.OprPlayer)
	}
	d.Data.Common.OprPlayer.Round++
	d.Data.Common.OprPlayer.Step = pb.RummyRoundStep_PLAYSTEPDRAW
	ntf := &pb.RummyPlayCardNtf{
		RoomId:    d.GetRoomId(),
		PlayerId:  d.Data.Common.OprPlayer.CurOprPlayerId,
		PlayStep:  int32(d.Data.Common.OprPlayer.Step),
		TimeOut:   d.Data.Common.TimeOut,
		TotalTime: d.GetCurStateTTL(),
	}
	err := d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_RummyPlayCard, ntf)
	mlog.Debugf("setNextPlayer NotifyToClient all err: %v", err)
	return d.GetCurState()
}

// 计算投降分数
func (d *RummyGame) getGiveUpCard(player *pb.RummyRoomPlayer) (score int64) {
	//计算积分
	if player.Private.PrevCard == 0 {
		// 没抽牌投降
		score = d.Data.RoomCfg.FirstDrop
	} else {
		score = d.Data.RoomCfg.MiddleDrop
	}
	return
}

func (d *RummyGame) GetCurPlayer() *pb.RummyRoomPlayer {
	return d.Data.Common.Players[d.Data.Common.OprPlayer.CurOprPlayerId]
}

// OnPlayerTimeout 玩家回合超时 标记退出
func (d *RummyGame) OnPlayerTimeout() {

	if d.machine.GetCurState() == pb.GameState_Rummy_STAGE_PLAYING {
		roomData := d.Data
		player := d.GetCurPlayer()
		oprType := pb.RummyOprType_Rummy_OPR_TYPE_LOSE
		var opcard uint32
		opcard = d.lostCard(player)
		player.TimeoutCount++

		// RummyOprCardNtf  玩家操作
		ntf := &pb.RummyOprCardNtf{
			RoomId:    d.GetRoomId(),
			PlayerId:  player.PlayerId,
			OprType:   oprType,
			OprCard:   opcard,
			ShowCard:  roomData.Common.ShowCard,
			ShowCard2: roomData.Common.ShowCard2,
			DrawCount: uint32(len(roomData.Private.Cards)) - roomData.Private.CardIdx,
			GaveScore: d.getGiveUpCard(player),
			ScorePoll: d.Data.Common.OprPlayer.ScorePool,
		}
		err := d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_RummyOprCard, ntf)

		mlog.Infof("玩家:%v 超时触发操作 %v 操作牌 %v 明牌1: %v 明牌2:%v ", player.PlayerId, oprType, card.Card(opcard).String(), card.Card(roomData.Common.ShowCard).String(), card.Card(roomData.Common.ShowCard2).String())

		//连续超时 按放弃游戏处理 连续超时两次 被t出 游戏开始自动被清理
		if player.TimeoutCount >= 2 {
			//设置阶段 和玩家已摸牌 按顶格处理
			roomData.Common.OprPlayer.Step = pb.RummyRoundStep_PLAYSTEPNONE
			d.giveupCard(player, false)
			player.Health = pb.RummyPlayHealth_Rummy_TIMEOUT

			ntf2 := &pb.RummyPlayerActionNtf{
				RoomId:    d.GetRoomId(),
				PlayerId:  player.PlayerId,
				Type:      pb.RummyPlayerActionType_Rummy_PLAYER_ACTION_GIVEUP,
				OprCard:   opcard,
				ScorePoll: d.Data.Common.OprPlayer.ScorePool,
			}
			err = d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_Timeout, ntf2)
			mlog.Infof("玩家:%v 超时触发投降操作 err : %v", player.PlayerId, err)
			player.TimeoutCount = 0
		}
	}

}

// giveupCard 投降
// 1.正常摸牌阶段投降; 20-40分
// 2.游戏中途退出游戏投降; 20-40分
// 3.炸胡投降; 80分
func (d *RummyGame) giveupCard(player *pb.RummyRoomPlayer, isFake bool) {
	// 更新投降状态
	if player.State == pb.RummyPlayState_Rummy_PLAYSTATE_PLAY {
		switch d.Data.RoomCfg.GameType {
		case pb.GameType_GameTypePR:
			d.SetPlayerLose(player.PlayerId, true, isFake)
		default: //分支玩法
			d.SetMatchPlayerLose(player.PlayerId, true, isFake)
		}
	}

	// 弃牌之前，自动把摸的最后一张牌，打出
	if len(player.Private.HandCards) == 14 {
		d.outCard(player, player.Private.PrevCard)

		oprType := pb.RummyOprType_Rummy_OPR_TYPE_LOSE
		ntf := &pb.RummyOprCardNtf{
			RoomId:    d.GetRoomId(),
			PlayerId:  player.PlayerId,
			OprType:   oprType,
			OprCard:   player.Private.PrevCard,
			ShowCard:  d.Data.Common.ShowCard,
			ShowCard2: d.Data.Common.ShowCard2,
			ScorePoll: d.Data.Common.OprPlayer.ScorePool, //推送变动后奖池
		}
		player.Private.CardGroup = DelCardRCG(player.Private.CardGroup, player.Private.PrevCard)
		// 玩家投降消息
		err := d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_RummyOprCard, ntf)
		mlog.Infof("err : %v", err)
	}

	if player.PlayerId == d.GetCurPlayer().PlayerId {
		d.Data.Common.OprPlayer.Step = pb.RummyRoundStep_PLAYSTEPSTEAFINISH
	}

}

// grabCard 抓牌
func (d *RummyGame) grabCard(player *pb.RummyRoomPlayer) uint32 {
	cardItem := d.Data.Private.Cards[d.Data.Private.CardIdx]
	d.Data.Private.CardIdx++

	player.Private.HandCards = append(player.Private.HandCards, cardItem)
	player.Private.IsGrabCard = true
	player.Private.PrevCard = cardItem
	d.Data.Common.OprPlayer.Step = pb.RummyRoundStep_PLAYSTEPDRAWDECLARE
	mlog.Infof("抽牌操作...用户: %v 当前手牌 %v", player.PlayerId, card.CardList(player.Private.HandCards))
	return cardItem
}

// chiCard 吃牌 抓明牌堆
func (d *RummyGame) chiCard(player *pb.RummyRoomPlayer) uint32 {
	cardItem := d.Data.Common.OprPlayer.OutCards[len(d.Data.Common.OprPlayer.OutCards)-1]
	d.Data.Common.OprPlayer.OutCards = d.Data.Common.OprPlayer.OutCards[:len(d.Data.Common.OprPlayer.OutCards)-1]

	///更新明牌和明牌堆
	if len(d.Data.Common.OprPlayer.OutCards) < 1 {
		d.Data.Common.ShowCard = 0
		d.Data.Common.ShowCard2 = 0
	} else {
		d.Data.Common.ShowCard = d.Data.Common.OprPlayer.OutCards[len(d.Data.Common.OprPlayer.OutCards)-1]
		if len(d.Data.Common.OprPlayer.OutCards) >= 2 {
			d.Data.Common.ShowCard2 = d.Data.Common.OprPlayer.OutCards[len(d.Data.Common.OprPlayer.OutCards)-2]
		}
	}

	///更新手牌
	player.Private.HandCards = append(player.Private.HandCards, cardItem)
	player.Private.IsGrabCard = true
	player.Private.PrevCard = cardItem
	d.Data.Common.OprPlayer.Step = pb.RummyRoundStep_PLAYSTEPDRAWDECLARE
	d.Data.Common.OprPlayer.IsEat = true
	mlog.Infof("吃牌操作...用户: %v 当前手牌 %v", player.PlayerId, card.CardList(player.Private.HandCards))

	return cardItem

}

// outCard 出牌
func (d *RummyGame) outCard(player *pb.RummyRoomPlayer, cardItem uint32) uint32 {
	//多副牌 倒序自动出牌
	for i := len(player.Private.HandCards) - 1; i >= 0; i-- {
		if player.Private.HandCards[i] == cardItem {
			player.Private.HandCards = append(player.Private.HandCards[:i], player.Private.HandCards[i+1:]...)

			d.Data.Common.OprPlayer.PrevOutUid = player.PlayerId
			d.Data.Common.OprPlayer.OutCards = append(d.Data.Common.OprPlayer.OutCards, cardItem)
			d.Data.Common.ShowCard = d.Data.Common.OprPlayer.OutCards[len(d.Data.Common.OprPlayer.OutCards)-1]
			if len(d.Data.Common.OprPlayer.OutCards) >= 2 {
				d.Data.Common.ShowCard2 = d.Data.Common.OprPlayer.OutCards[len(d.Data.Common.OprPlayer.OutCards)-2]
			}
			mlog.Infof("出牌操作...用户: %v 当前手牌 %v", player.PlayerId, card.CardList(player.Private.HandCards))
			return cardItem
		}
	}
	return 0
}

// lostCard 多副牌弃牌可能改变卡序
func (d *RummyGame) lostCard(player *pb.RummyRoomPlayer) (cardItem uint32) {
	switch d.Data.Common.OprPlayer.Step {
	case pb.RummyRoundStep_PLAYSTEPDRAW:
		cardItem = d.grabCard(player)
	case pb.RummyRoundStep_PLAYSTEPDRAWDECLARE:
		if d.Data.Common.OprPlayer.IsEat { // 归还明牌  抽暗牌打出
			d.outCard(player, player.Private.PrevCard)
			player.Private.CardGroup = DelCardRCG(player.Private.CardGroup, player.Private.PrevCard)
			cardItem = d.grabCard(player)
			player.Private.PrevCard = cardItem
		} else { //玩家抽牌后打出
			cardItem = player.Private.PrevCard
			player.Private.CardGroup = DelCardRCG(player.Private.CardGroup, cardItem)
		}
	}
	d.outCard(player, cardItem)
	return
}

func (d *RummyGame) HiddenPlayerPriData(common *pb.RummyCommon, playerId uint64) *pb.RummyCommon {
	newCommon := proto.Clone(common).(*pb.RummyCommon)
	for key := range newCommon.Players {
		if key != playerId {
			newCommon.Players[key].Private.Reset()
		}
	}
	return newCommon
}

func (d *RummyGame) GetSettlePlayerInfo(player *pb.RummyRoomPlayer) *pb.RummySettlePlayerInfo {
	return &pb.RummySettlePlayerInfo{
		PicUrl:    player.PicUrl,
		NickName:  player.PlayerNick,
		PlayerId:  player.PlayerId,
		CardGroup: player.Private.CardGroup,
		Score:     player.Private.Score,
		State:     player.State,
		Coin:      -d.Data.RoomCfg.BaseScore * player.Private.Score,
		BonusCash: int64(player.Coin),
		Total:     player.Total,
	}
}
