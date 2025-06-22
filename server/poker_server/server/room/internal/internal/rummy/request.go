package rummy

import (
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"poker_server/library/util"
	"poker_server/server/room/internal/module/card"
)

// GetRummyOutCards 获取所有出牌列表
func (d *RummyGame) GetRummyOutCards(head *pb.Head, req *pb.RummyGetOutCardsReq, rsp *pb.RummyGetOutCardsRsp) error {
	if uint32(d.Data.Stage) < uint32(pb.GameState_Rummy_STAGE_DEAL) {
		return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_NOT_PLAY, head, "invaild req with wrong time")
	}
	rsp.RoomId = d.GetRoomId()
	rsp.OutCards = d.Data.Common.OprPlayer.OutCards
	return nil
}

// SaveCardGroup 保存玩家手牌
func (d *RummyGame) SaveCardGroup(head *pb.Head, req *pb.RummySaveCardGroupReq, rsp *pb.RummySaveCardGroupRsp) error {
	player, ok := d.Data.Common.Players[head.Uid]
	if !ok || player.State != pb.RummyPlayState_Rummy_PLAYSTATE_PLAY {
		return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_NOT_PLAY, head, "player not exist")
	}
	isReal, Score := CheckRCG(player.Private.HandCards, req.GetGroups())
	if !isReal { ///数据有效性检查
		mlog.Infof("用户手牌: %v \n", card.CardList(player.Private.HandCards))
		save := req.GetGroups()
		for i := range save {
			mlog.Infof("用户提交卡组 %v: %v \n", i, card.CardList(save[i].Cards))
		}

		return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_DATA_INVAILD, head, "invaild player data")
	}
	d.Change()

	player.Private.CardGroup = req.GetGroups()
	for i := range player.Private.CardGroup {
		player.Private.CardGroup[i].GroupType, _ = GetCardType(player.Private.CardGroup[i].Cards)
	}

	player.Private.Score = Score
	rsp.RoomId = d.GetRoomId()
	rsp.Score = player.Private.Score
	rsp.Groups = player.Private.CardGroup

	// 确认阶段 最后保存一次手牌
	if req.LastReq && d.machine.GetCurState() == pb.GameState_Rummy_STAGE_FIX_CARD && player.State == pb.RummyPlayState_Rummy_PLAYSTATE_PLAY {
		d.SetPlayerLose(player.PlayerId, false)

		playerIds := make([]uint64, 0, len(d.Data.Common.PlayerIds))
		d.Walk(0, func(playerItem *pb.RummyRoomPlayer) bool {
			if playerItem.State == pb.RummyPlayState_Rummy_PLAYSTATE_PLAY {
				playerIds = append(playerIds, playerItem.PlayerId)
			}
			return true
		})

		ntf := &pb.RummyFixCardPlayersNtf{
			RoomId:      d.GetRoomId(),
			TimeOut:     d.Data.Common.TimeOut,
			WinId:       d.Data.Common.WinnerId,
			Players:     playerIds,
			CurPlayerId: player.PlayerId,
			TotalTime:   d.Data.Common.TotalTime,
		}
		err := d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_RummyFixCardPlayers, ntf)
		mlog.Infof("SaveCardGroup NotifyToClient %v", err)
	}
	return nil
}

// FixCard 声明操作 处理玩家胡牌
func (d *RummyGame) FixCard(head *pb.Head, req *pb.RummyFixCardReq, rsp *pb.RummyFixCardRsp) error {
	if head.Uid != d.GetCurPlayer().PlayerId {
		return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_STEP_INVAILD, head, "player is not current")
	}

	player := d.Data.Common.Players[head.Uid]
	if d.Data.Common.OprPlayer.Step != pb.RummyRoundStep_PLAYSTEPDRAWDECLARE {
		return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_ACTION_INVAILD, head, "player action is not allow")
	}

	carditem := req.GetOprCard()
	carditem = d.outCard(player, carditem)
	if carditem == 0 {
		return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_DATA_INVAILD, head, "noexist handcard")
	}

	isReal, Score := CheckRCG(player.Private.HandCards, req.GetGroups())
	if !isReal {
		return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_DATA_INVAILD, head, "invaild player data")
	}

	d.Change()
	cardGroup := req.GetGroups()
	for i := 0; i < len(cardGroup); i++ {
		cardGroup[i].GroupType, _ = GetCardType(cardGroup[i].Cards)
	}
	player.Private.Score = Score
	player.Private.CardGroup = cardGroup

	rsp.Score = Score
	rsp.Groups = player.Private.CardGroup

	if Score == 0 {
		player.State = pb.RummyPlayState_Rummy_PLAYSTATE_WIN
		d.Data.Common.OprPlayer.Step = pb.RummyRoundStep_PLAYSTEPNONE
		d.Data.Common.WinnerId = player.PlayerId
		d.Data.Common.FixCard = carditem
		d.Data.Common.GameFinish = true //游戏结束

		ntf2 := &pb.RummyOprCardNtf{
			RoomId:    d.GetRoomId(),
			PlayerId:  player.PlayerId,
			OprType:   pb.RummyOprType_Rummy_OPR_TYPE_FIX_FINISH,
			ShowCard:  d.Data.Common.ShowCard,
			ShowCard2: d.Data.Common.ShowCard2,
			OprCard:   carditem,
		}

		//广播胡牌前最后一张牌
		err := d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_CustomFixCard, ntf2)
		mlog.Infof("SaveCardGroup NotifyToClient %v", err)
	} else {
		//判定直接投降 直接结算
		d.giveupCard(player)
		d.Data.Common.OprPlayer.Step = pb.RummyRoundStep_PLAYSTEPSTEAFINISH

		//广播炸胡
		ntf2 := &pb.RummyPlayerActionNtf{
			RoomId:   d.GetRoomId(),
			PlayerId: player.PlayerId,
			Type:     pb.RummyPlayerActionType_Rummy_PLAYER_ACTION_FINISH_ERR,
			OprCard:  carditem,
		}
		err := d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_CustomFake, ntf2)
		mlog.Infof("SaveCardGroup NotifyToClient %v", err)
		return uerror.NEW(pb.ErrorCode_FAILED, head, "fraud cardgroup lose this game")
	}
	return nil
}

// OprCardReq 处理玩家操作牌: 摸牌，吃牌，弃牌，完成牌等四种操作
func (d *RummyGame) OprCardReq(head *pb.Head, req *pb.RummyOprCardReq, rsp *pb.RummyOprCardRsp) error {
	if d.machine.GetCurState() != pb.GameState_Rummy_STAGE_PLAYING {
		return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_ACTION_INVAILD, head, "this stage cant play")
	}

	mlog.Infof("OprCardReq head:%v req:%v debug : %v", head, req, d.GetCurPlayer())
	if head.Uid != d.GetCurPlayer().PlayerId && req.GetOprType() != pb.RummyOprType_Rummy_OPR_TYPE_GIVEUP {
		return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_STEP_INVAILD, head, "player is not current")
	}

	player, ok := d.Data.Common.Players[head.Uid]
	if !ok {
		return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_NOT_IN_ROOM, head, "player is not exist")
	}

	//这个阶段只能 抽 吃 投
	if d.Data.Common.OprPlayer.Step == pb.RummyRoundStep_PLAYSTEPDRAW && (int32(req.GetOprType()) > int32(pb.RummyOprType_Rummy_OPR_TYPE_GIVEUP) || req.GetOprType() == pb.RummyOprType_Rummy_OPR_TYPE_UNDEFINE) {
		return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_ACTION_INVAILD, head, "invalid operation error")
	} else if d.Data.Common.OprPlayer.Step == pb.RummyRoundStep_PLAYSTEPDRAWDECLARE && (req.GetOprType() != pb.RummyOprType_Rummy_OPR_TYPE_FIX && req.GetOprType() != pb.RummyOprType_Rummy_OPR_TYPE_OUT && req.GetOprType() != pb.RummyOprType_Rummy_OPR_TYPE_GIVEUP) {
		//这个阶段只能 打牌 胡牌 投降
		return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_ACTION_INVAILD, head, "DECLARE step invalid operation error")
	}

	// 要出的牌
	cardItem := req.GetOprCard()
	switch req.GetOprType() {
	case pb.RummyOprType_Rummy_OPR_TYPE_GRAB: //抽牌
		cardItem = d.grabCard(player)
	case pb.RummyOprType_Rummy_OPR_TYPE_CHI: //吃牌
		cardItem = d.chiCard(player)
	case pb.RummyOprType_Rummy_OPR_TYPE_OUT: //出牌
		cardItem = d.outCard(player, cardItem)
		if cardItem == 0 { //用户手牌没变 或者打吃的牌
			return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_DATA_INVAILD, head, "invalid player hardcards Data")
		}

		// 禁止摸什么牌打什么牌  手上有两张一样的牌除外
		if d.Data.Common.OprPlayer.IsEat && cardItem == player.Private.PrevCard {
			same := 0
			for i := range player.Private.HandCards {
				if cardItem == player.Private.HandCards[i] {
					//打的牌手牌数量
					same++
				}
			}
			if same <= 1 {
				return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_DATA_INVAILD, head, "invalid player hardcards Data")
			}
		}
		d.Data.Common.OprPlayer.Step = pb.RummyRoundStep_PLAYSTEPSTEAFINISH
	case pb.RummyOprType_Rummy_OPR_TYPE_GIVEUP: // 投降
		d.giveupCard(player)
	}

	if req.GetOprType() == pb.RummyOprType_Rummy_OPR_TYPE_GRAB || req.GetOprType() == pb.RummyOprType_Rummy_OPR_TYPE_CHI {
		player.Private.CardGroup = AddCardRCG(player.Private.CardGroup, cardItem) // 更新手牌
	}
	if req.GetOprType() == pb.RummyOprType_Rummy_OPR_TYPE_OUT { //出牌后更新牌序
		isReal, score := CheckRCG(player.Private.HandCards, req.GetGroups())
		if !isReal {
			return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_DATA_INVAILD, head, "invaild player data")
		}
		player.Private.Score = score
		player.Private.CardGroup = req.GetGroups()
		for i := 0; i < len(player.Private.CardGroup); i++ {
			player.Private.CardGroup[i].GroupType, _ = GetCardType(player.Private.CardGroup[i].Cards)
		}
	}

	d.Change()

	player.TimeoutCount = 0 // 连续超时次数清空
	rsp.OprType = req.OprType
	rsp.OprCard = cardItem
	rsp.PlayStep = int32(d.Data.Common.OprPlayer.Step)
	rsp.IsEat = d.Data.Common.OprPlayer.IsEat

	rsp.Groups = player.Private.CardGroup
	ntf := &pb.RummyOprCardNtf{ //广播玩家操作
		RoomId:    d.GetRoomId(),
		PlayerId:  player.PlayerId,
		OprType:   req.OprType,
		OprCard:   req.OprCard,
		ShowCard:  d.Data.Common.ShowCard,
		ShowCard2: d.Data.Common.ShowCard2,
		DrawCount: uint32(len(d.Data.Private.Cards)) - d.Data.Private.CardIdx,
		ScorePoll: d.Data.Common.OprPlayer.ScorePool,
	}
	err := d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_RummyOprCard, ntf)
	mlog.Infof("SaveCardGroup NotifyToClient %v", err)
	return nil
}

// ReadyRoomReq 玩家就绪
func (d *RummyGame) ReadyRoomReq(head *pb.Head, req *pb.RummyReadyRoomReq, rsp *pb.RummyReadyRoomRsp) error {
	isPlaying := util.SliceIsset[uint64](head.Uid, d.Data.Common.PlayerIds)

	if !isPlaying {
		// 检查房间是否存在这个玩家
		if _, ok := d.Data.Common.Players[head.Uid]; !ok {
			return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_NOT_IN_ROOM, head, "player is not exist")
		}
		d.Change()
		player := d.Data.Common.Players[head.Uid]
		player.State = pb.RummyPlayState_Rummy_PLAYSTATE_READY
		d.UpReadyPlayer()
	}

	rsp.Stage = d.Data.Stage
	rsp.RoomInfo = &pb.RummyRoomPubData{
		RoomId:   d.GetRoomId(),
		Stage:    d.Data.Stage,
		GameId:   d.Data.GameId,
		RoomName: d.Data.RoomName,
		Common:   d.HiddenPlayerPriData(d.Data.Common, head.Uid),
		Match:    d.Data.Match,
		RoomCfg:  d.Data.RoomCfg,
	}
	return nil
}

func (d *RummyGame) JoinRoomReq(head *pb.Head, req *pb.RummyJoinRoomReq, rsp *pb.RummyJoinRoomRsp) error {
	// 检查房间是否存在这个玩家
	if _, ok := d.Data.Common.Players[head.Uid]; !ok {
		if len(d.Data.Common.EmptySeats) <= 0 {
			return uerror.NEW(pb.ErrorCode_RUMMY_ROOM_FULL, head, "rummy room full already")
		}
		player := d.GetOrNewPlayer(head.Uid)
		d.Data.Common.Players[head.Uid] = player
		d.Change()
	} else {
		d.Data.Common.Players[head.Uid].Health = pb.RummyPlayHealth_Rummy_NORMAL
	}

	// 入桌消息
	ntf := &pb.RummyEnterDeskNtf{
		RoomId: d.GetRoomId(),
		Player: d.Data.Common.Players[head.Uid],
	}
	err := d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_RummyEnterDesk, ntf)
	mlog.Infof("JoinRoomReq ntf send: %v ntf : %v", err, ntf)

	roomInfo := &pb.RummyRoomPubData{
		RoomId:   d.GetRoomId(),
		Stage:    d.Data.Stage,
		GameId:   d.Data.GameId,
		RoomName: d.Data.RoomName,
		Common:   d.HiddenPlayerPriData(d.Data.Common, head.Uid),
		Match:    d.Data.Match,
		RoomCfg:  d.Data.RoomCfg,
	}

	if req.IsChange {
		changeRsp := &pb.RummyChangeRoomRsp{
			RoomId:    d.GetRoomId(),
			RoomInfo:  roomInfo,
			GaveScore: d.getGiveUpCard(d.Data.Common.Players[head.Uid]),
		}
		return framework.SendToClient(head, changeRsp)
	} else {
		// 返回数据
		rsp.RoomId = d.GetRoomId()
		rsp.RoomInfo = roomInfo
		rsp.GaveScore = d.getGiveUpCard(d.Data.Common.Players[head.Uid])
		return nil
	}
}

func (d *RummyGame) QuitRoomReq(head *pb.Head, req *pb.RummyQuitRoomReq, rsp *pb.RummyQuitRoomRsp) error {
	player, ok := d.Data.Common.Players[head.Uid]
	if !ok || player == nil {
		return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_NOT_IN_ROOM, head, "玩家不在房间内: %d", head.Uid)
	}

	for i := 0; i < len(d.Data.Common.PlayerIds); i++ {
		if uint32(d.machine.GetCurState()) >= uint32(pb.GameState_Rummy_STAGE_START) && d.Data.Common.PlayerIds[i] == head.Uid && player.State == pb.RummyPlayState_Rummy_PLAYSTATE_PLAY { //非观战
			return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_GAMEING, head, "玩家正在游戏中: %d", head.Uid)
		}
	}
	d.Change()

	// 立即释放座位
	d.Data.Common.EmptySeats = append(d.Data.Common.EmptySeats, d.Data.Common.Players[head.Uid].Seat)
	player.Health = pb.RummyPlayHealth_Rummy_QUIT

	//推送退出消息
	ntf := &pb.RummyQuitRoomNtf{
		RoomId:    d.GetRoomId(),
		PlayerId:  head.Uid,
		LeaveType: pb.RummyLeaveType_Rummy_LEAVE_TYPE_QUIT,
	}
	err := d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_RummyQuitRoom, ntf)
	mlog.Infof("QuitRoomReq ntf send: %v", err)

	rsp.RoomId = d.GetRoomId()
	if req.IsChange {
		changeReq := &pb.RummyJoinRoomReq{
			RoomId: d.GetRoomId(),
		}
		return framework.Send(framework.SwapToGame(head, head.Uid, "Player", "RummyChangeNewRoomReq"), changeReq)
	}
	// 删除房间信息
	return nil
}
