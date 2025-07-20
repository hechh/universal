package rummy

import (
	"math"
	"poker_server/common/pb"
	"poker_server/library/mlog"
)

func (d *RummyGame) Stop() {
	mlog.Infof("Rummy Extra Game关闭: roomId:%d GameType: %v", d.GetId(), d.Data.RoomCfg.GameType)
	d.Actor.Stop()
}

func (d *RummyGame) ExtReset() {
	d.RuntimeGC()
	// reset game public data
	d.Data.Stage = pb.GameState_RummyExt_STAGE_INIT
	// 回收match per数据
	d.Data.Match.Match++
	d.Data.Match.StartTime = 0
	d.Data.Match.EndTime = 0

	d.Walk(0, func(player *pb.RummyRoomPlayer) bool {
		player.Private.Reset()
		if player.State != pb.RummyPlayState_Rummy_PLAYSTATE_ELIMINATED { //淘汰玩家可在观众席
			player.State = pb.RummyPlayState_Rummy_PLAYSTATE_PLAY
		}
		player.Private.OutCards = player.Private.OutCards[:0]
		return true
	})
}

// CostTicket 支付入场费
func (d *RummyGame) CostTicket(player *pb.RummyRoomPlayer) {
	player.Coin -= uint64(d.Data.RoomCfg.BaseScore)
	d.Data.Common.PrizePool += d.Data.RoomCfg.BaseScore

	switch d.Data.RoomCfg.GameType {
	case pb.GameType_GameTypeDR:
		player.Total = int64(d.Data.RoomCfg.Deals) * 80
	default: // pool rummy玩法
		player.Total = 0
	}
}

// pool rummy 淘汰赛 and deal rummy
func (d *RummyGame) giveUpGame(playerId uint64) {
	player := d.Data.Common.Players[playerId]
	if d.Data.Status == pb.RoomStatus_RoomStatusPlaying && player.State == pb.RummyPlayState_Rummy_PLAYSTATE_PLAY {
		d.giveupCard(player, false)

		ntf := &pb.RummyOprCardNtf{ //广播玩家操作
			RoomId:    d.GetRoomId(),
			PlayerId:  playerId,
			OprType:   pb.RummyOprType_Rummy_OPR_TYPE_GIVEUP,
			OprCard:   player.Private.PrevCard,
			ShowCard:  d.Data.Common.ShowCard,
			ShowCard2: d.Data.Common.ShowCard2,
			DrawCount: uint32(len(d.Data.Private.Cards)) - d.Data.Private.CardIdx,
			ScorePoll: d.Data.Common.OprPlayer.ScorePool,
		}
		err := d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_RummyOprCard, ntf)
		mlog.Infof("player : %v giveUpGame NotifyToClient %v", playerId, err)
	}
}

// SetMatchPlayerLose 扩展玩法 设置玩家判负
func (d *RummyGame) SetMatchPlayerLose(playerID uint64, isFold, isFake bool) {
	player := d.Data.Common.Players[playerID]
	if isFold { //玩家弃牌
		player.State = pb.RummyPlayState_Rummy_PLAYSTATE_GIVEUP
		if isFake {
			player.Private.Score = 80 //炸胡
		} else {
			player.Private.Score = d.getGiveUpCard(player)
		}
		mlog.Infof("Rummy Ext SetPlayerLose Fold PlayerId:%d Score:%d", playerID, player.Private.Score)
	} else {
		player.State = pb.RummyPlayState_Rummy_PLAYSTATE_LOSE
		_, Score := CheckRCG(player.Private.HandCards, player.Private.CardGroup)
		player.Private.Score = Score
		mlog.Infof("Rummy Ext SetPlayerLose PlayerId:%d Score:%d", playerID, player.Private.Score)
	}
}

// GetContinues 判断继续游戏玩家数
func (d *RummyGame) GetContinues() (normal int) {
	diff := int64(0)
	d.Walk(0, func(player *pb.RummyRoomPlayer) bool {
		if player.State != pb.RummyPlayState_Rummy_PLAYSTATE_ELIMINATED {
			normal++
			diff ^= player.Total
		}
		return true
	})

	switch d.Data.RoomCfg.GameType {
	case pb.GameType_GameTypeDR:
		if d.Data.Match.Match >= d.Data.RoomCfg.Deals {
			return 1
		}
	default:
		if normal == 2 && diff == 0 {
			return math.MaxInt
		}
	}
	return
}

func (d *RummyGame) GetExtSettlePlayerInfo(player *pb.RummyRoomPlayer, amount int64) *pb.RummySettlePlayerInfo {
	return &pb.RummySettlePlayerInfo{
		PicUrl:    player.PicUrl,
		NickName:  player.PlayerNick,
		PlayerId:  player.PlayerId,
		CardGroup: player.Private.CardGroup,
		State:     player.State,
		Coin:      amount,
		BonusCash: int64(player.Coin),
		Total:     player.Total,
		GhostCard: d.Data.Common.GhostCard,
	}
}

func (d *RummyGame) GetPlayerSettleInfo(player *pb.RummyRoomPlayer) *pb.RummySettlePlayerInfo {
	return &pb.RummySettlePlayerInfo{
		PlayerId:  player.PlayerId,
		CardGroup: player.Private.CardGroup,
		Score:     player.Private.Score,
		State:     player.State,
		NickName:  player.PlayerNick,
		BonusCash: player.JoinCoin,
		GhostCard: d.Data.Common.GhostCard,
		Total:     player.Total,
	}
}
