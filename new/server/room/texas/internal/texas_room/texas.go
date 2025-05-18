package texas_room

import (
	"poker_server/common/pb"
	"poker_server/framework/library/mlog"
	"poker_server/server/room/texas/domain"
	"poker_server/server/room/texas/internal/base"
	"sort"
)

func (d *TexasRoom) GetCompare() domain.ICompare {
	return base.NewBackTrackComapre()
}

func (d *TexasRoom) GetCursor() *pb.TexasPlayerData {
	gameData := d.data.Table.GameData
	llen := len(gameData.UidList)
	return d.data.Table.Players[gameData.UidList[gameData.UidCursor%uint32(llen)]]
}

func (d *TexasRoom) GetDealer() *pb.TexasPlayerData {
	gameData := d.data.Table.GameData
	ll := len(gameData.UidList)
	usr := d.data.Table.Players[d.data.Table.ChairInfo[gameData.DealerChairId]]
	return d.data.Table.Players[gameData.UidList[int(usr.GameInfo.Position)%ll]]
}

func (d *TexasRoom) GetSmall() *pb.TexasPlayerData {
	gameData := d.data.Table.GameData
	ll := len(gameData.UidList)
	usr := d.data.Table.Players[d.data.Table.ChairInfo[gameData.SmallChairId]]
	return d.data.Table.Players[gameData.UidList[int(usr.GameInfo.Position)%ll]]
}

func (d *TexasRoom) GetBig() *pb.TexasPlayerData {
	gameData := d.data.Table.GameData
	ll := len(gameData.UidList)
	usr := d.data.Table.Players[d.data.Table.ChairInfo[gameData.BigChairId]]
	return d.data.Table.Players[gameData.UidList[int(usr.GameInfo.Position)%ll]]
}

func (d *TexasRoom) Walk(pos int, f func(*pb.TexasPlayerData) bool) {
	size := len(d.data.Table.GameData.UidList)
	for i := pos; i < pos+size; i++ {
		if !f(d.data.Table.Players[d.data.Table.GameData.UidList[i%size]]) {
			return
		}
	}
}

func (d *TexasRoom) WalkPrev(pos int, f func(*pb.TexasPlayerData) bool) {
	size := len(d.data.Table.GameData.UidList)
	for i := pos + size - 1; i >= pos; i-- {
		if !f(d.data.Table.Players[d.data.Table.GameData.UidList[i%size]]) {
			return
		}
	}
}

func (d *TexasRoom) GetPrev(pos int, state pb.TexasGameState, opr ...pb.TexasOperateType) *pb.TexasPlayerData {
	size := len(d.data.Table.GameData.UidList)
	for i := pos + size - 1; i > pos; i-- {
		flag := false
		usr := d.data.Table.Players[d.data.Table.GameData.UidList[i%size]]
		if usr.GameInfo.GameState != state {
			continue
		}
		for _, opr := range opr {
			flag = flag || (usr.GameInfo.Operate == opr)
		}
		if flag {
			continue
		}
		return usr
	}
	return nil
}

func (d *TexasRoom) GetNext(pos int, state pb.TexasGameState, oprs ...pb.TexasOperateType) *pb.TexasPlayerData {
	size := len(d.data.Table.GameData.UidList)
	for i := pos + 1; i < pos+size; i++ {
		flag := false
		usr := d.data.Table.Players[d.data.Table.GameData.UidList[i%size]]
		if usr.GameInfo.GameState != state {
			continue
		}
		for _, opr := range oprs {
			flag = flag || (usr.GameInfo.Operate == opr)
		}
		if flag {
			continue
		}
		return usr
	}
	return nil
}

func (d *TexasRoom) GetPlayers(state pb.TexasGameState) (users []*pb.TexasPlayerData) {
	for _, uid := range d.data.Table.GameData.UidList {
		usr := d.data.Table.Players[uid]
		if usr.GameInfo.GameState == state {
			users = append(users, usr)
		}
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].GameInfo.BetChips < users[j].GameInfo.BetChips
	})
	return
}

func (d *TexasRoom) Operate(usr *pb.TexasPlayerData, opr pb.TexasOperateType, val int64) {
	switch opr {
	case pb.TexasOperateType_TOT_CHECK, pb.TexasOperateType_TOT_FOLD:
		usr.IsChange = true
		usr.GameInfo.Operate = opr
	case pb.TexasOperateType_TOT_BET_BIG_BLIND, pb.TexasOperateType_TOT_BET_SMALL_BLIND:
		usr.IsChange = false
		usr.GameInfo.Operate = opr
		usr.GameInfo.BetChips += val
		usr.Chips -= val
		d.data.Table.GameData.MinRaise = val
		d.data.Table.GameData.PotPool.TotalBetChips += val
	default:
		usr.IsChange = true
		usr.GameInfo.Operate = opr
		usr.GameInfo.BetChips += val
		d.data.Table.GameData.PotPool.TotalBetChips += val
	}
	if d.data.Table.GameData.MaxBetChips < usr.GameInfo.BetChips {
		d.data.Table.GameData.MaxBetChips = usr.GameInfo.BetChips
	}
}

func (d *TexasRoom) UpdateMain(users ...*pb.TexasPlayerData) {
	if len(users) > 0 {
		game := d.data.Table.GameData
		game.PotPool.PotList = append(game.PotPool.PotList, &pb.TexasPotData{
			PotType: 0,
			Chips:   game.PotPool.BetChips,
			UidList: playerToUid(users...),
		})
		game.PotPool.BetChips = 0
	}
}

func (d *TexasRoom) UpdateSide(users ...*pb.TexasPlayerData) {
	count := len(users)
	prevBetChips := int64(0)
	game := d.data.Table.GameData

	for i, usr := range users {
		switch usr.GameInfo.Operate {
		case pb.TexasOperateType_TOT_FOLD:
			game.PotPool.BetChips += (usr.GameInfo.BetChips - prevBetChips)

		case pb.TexasOperateType_TOT_ALL_IN:
			add := int64(count-i) * (usr.GameInfo.BetChips - prevBetChips)
			if add > 0 {
				game.PotPool.PotList = append(game.PotPool.PotList, &pb.TexasPotData{
					PotType: 1,
					Chips:   game.PotPool.BetChips + add,
					UidList: playerToUid(users[i:]...),
				})
				game.PotPool.BetChips = 0
			}
			prevBetChips = usr.GameInfo.BetChips

		default:
			game.PotPool.BetChips += (usr.GameInfo.BetChips - prevBetChips)

		}
		usr.GameInfo.BetChips = 0
	}
}

func playerToUid(users ...*pb.TexasPlayerData) (rets []uint64) {
	for _, usr := range users {
		rets = append(rets, usr.GetUid())
	}
	return
}

// 计算获胜者筹码
func (d *TexasRoom) Reward(count int, f func(uint64, int64)) (total int64) {
	// 默认值
	rate := int64(0)
	maxLimit := int64(0)
	isNoFlopNoDrop := true

	// 读取配置
	isNoFlopNoDrop = d.texasCfg.IsNoFlopNoDrop
	rate = d.texasCfg.RakeRate
	if count > len(d.texasCfg.RakeUpLimit) {
		maxLimit = d.texasCfg.RakeUpLimit[len(d.texasCfg.RakeUpLimit)-1]
	} else {
		maxLimit = d.texasCfg.RakeUpLimit[count-1]
	}

	// 计算抽水
	// 没有进入翻牌阶段,不需要抽取服务费
	if !(isNoFlopNoDrop && len(d.data.Table.GameData.PublicCardList) <= 0) {
		for _, pot := range d.data.Table.GameData.PotPool.PotList {
			if len(pot.UidList) == 1 {
				continue
			}
			diff := (pot.Chips * rate / 10000)
			if diff > maxLimit {
				diff = maxLimit
			}
			total += diff
			pot.Chips -= diff
		}
	}
	tmps := map[uint64]int64{}
	for _, pot := range d.data.Table.GameData.PotPool.PotList {
		// 获取成员
		users := []*pb.TexasPlayerData{}
		for _, uid := range pot.UidList {
			users = append(users, d.data.Table.Players[uid])
		}
		// 比较牌力大小
		sort.Slice(users, func(i, j int) bool {
			return users[i].GameInfo.BestCardValue < users[j].GameInfo.BestCardValue
		})
		// 获取获胜者
		j := -1
		maxValue := users[0].GameInfo.BestCardValue
		for _, usr := range users {
			if usr.GameInfo.BestCardValue == maxValue {
				j++
				users[j] = usr
				continue
			}
			break
		}
		users = users[:j+1]
		// 统计获胜筹码
		for _, usr := range users {
			tmps[usr.Uid] += pot.Chips / int64(len(users))
		}
	}
	for uid, chips := range tmps {
		f(uid, chips)
	}
	return
}

// 洗牌
func (d *TexasRoom) Shuffle(times int) {
	// 初始化牌堆
	gameData := d.data.Table.GameData
	gameData.CardCursor = 0
	for i := uint32(0); i <= 3; i++ {
		for j := uint32(2); j <= 14; j++ {
			gameData.CardList = append(gameData.CardList, (1<<(16+i))|j)
		}
	}
	// 洗牌
	for j := 0; j < times; j++ {
		for i := 0; i < len(gameData.CardList); i++ {
			pos := base.Intn(len(gameData.CardList))
			gameData.CardList[i], gameData.CardList[pos] = gameData.CardList[pos], gameData.CardList[i]
		}
	}
	mlog.Debugf("洗牌次数: %d, gameData:%v", times, gameData)
}

// 发牌
func (d *TexasRoom) Deal(count uint32, f func(uint32, uint32)) {
	table := d.data.Table
	for i := table.GameData.CardCursor; i < table.GameData.CardCursor+count; i++ {
		if f != nil {
			f(i, table.GameData.CardList[i])
		}
	}
	table.GameData.CardCursor += count
}

func (d *TexasRoom) UpdateBest(usr *pb.TexasPlayerData, cmp domain.ICompare, publics []uint32) {
	tmps := append([]uint32{}, publics...)
	tmps = append(tmps, usr.GameInfo.HandCardList...)
	sort.Slice(tmps, func(i, j int) bool {
		return base.Card(tmps[i]).Rank() > base.Card(tmps[j]).Rank()
	})

	usr.GameInfo.BestCardType, usr.GameInfo.BestCardValue, usr.GameInfo.BestCardList = cmp.Get(tmps...)
	if base.Card(tmps[0]).Rank() == pb.TexasRank_TR_A {
		tmps = tmps[1:]
		tmps = append(tmps, (tmps[0]&0xFFFF0000)|uint32(pb.TexasRank_TR_1))
		bestType, value, cards := cmp.Get(tmps...)
		if usr.GameInfo.BestCardValue < value {
			usr.GameInfo.BestCardType = bestType
			usr.GameInfo.BestCardValue = value
			usr.GameInfo.BestCardList = cards
		}
	}
}
