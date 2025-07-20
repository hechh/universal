package rummy

import (
	"math"
	"math/bits"
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/library/random"
	"poker_server/server/room/internal/module/card"
	"sort"
)

// 初始化牌堆
func GetRummyOneCards() (cards []uint32) {
	for i := uint32(0); i <= 3; i++ {
		for j := uint32(2); j <= 14; j++ {
			cards = append(cards, (1<<(16+i))|j)
		}
	}
	cards = append(cards, (1<<(31))|15) //joker
	return
}

func GetRummyCardsByCfg(decks, joker uint32) (cards []uint32) {
	for time := uint32(0); time < decks; time++ {
		for i := uint32(0); i <= 3; i++ {
			for j := uint32(2); j <= 14; j++ {
				cards = append(cards, (1<<(16+i))|j)
			}
		}
	}
	for i := uint32(0); i < joker; i++ {
		if i%2 == 0 {
			cards = append(cards, (1<<(31))|15) //small joker
		} else {
			cards = append(cards, (1<<(31))|16) //big joker
		}
	}
	return
}

// 洗牌
func Shuffle(cards []uint32, times int) {
	// 洗牌
	for j := 0; j < times; j++ {
		for i := 0; i < len(cards); i++ {
			pos := random.Intn(len(cards))
			cards[i], cards[pos] = cards[pos], cards[i]
		}
	}
}

// 手牌初始化
func NewCardGroup(hands []uint32) (rcg []*pb.RummyCardGroup) {
	tmp := make(map[uint32]*pb.RummyCardGroup, 5)
	// 按花色初始分组
	for _, cardItem := range hands {
		if !card.Card(cardItem).IsWild() && tmp[uint32(card.Card(cardItem).Color())] == nil {
			tmp[uint32(card.Card(cardItem).Color())] = &pb.RummyCardGroup{
				Cards:     make([]uint32, 0, 13),
				GroupType: pb.RummyGroupType_Rummy_GROUP_TYPE_INVALID, //默认不确定类型
			}
		}

		//癞子独立牌堆
		if card.Card(cardItem).IsWild() {
			if tmp[0] == nil {
				tmp[0] = &pb.RummyCardGroup{
					Cards:     make([]uint32, 0, 13),
					GroupType: pb.RummyGroupType_Rummy_GROUP_TYPE_INVALID, //默认不确定类型
				}
			}
			tmp[0].Cards = append(tmp[0].Cards, cardItem)
		} else {
			tmp[uint32(card.Card(cardItem).Color())].Cards = append(tmp[uint32(card.Card(cardItem).Color())].Cards, cardItem)
		}
	}

	for _, item := range tmp {
		rcg = append(rcg, item)
	}

	// 按牌值升序
	for _, cards := range rcg {
		//默认升序
		sort.Slice(cards.Cards, func(i, j int) bool {
			return cards.Cards[i] < cards.Cards[j]
		})
		cards.GroupType, _ = GetCardType(cards.Cards)
	}
	return
}

// 手牌 rcg 追加一张牌 如果是赖子单独开辟一个group 如果满足最大容量或者是常规牌追加在最后一个元素堆
func AddCardRCG(rcg []*pb.RummyCardGroup, card uint32) []*pb.RummyCardGroup {
	if len(rcg) < 7 {
		rcg = append(rcg, &pb.RummyCardGroup{Cards: []uint32{card}})
		return rcg
	}
	//默认
	rcg[len(rcg)-1].Cards = append(rcg[len(rcg)-1].Cards, card)
	return rcg
}

// 手牌自定义牌序 rcg 删除一张牌
func DelCardRCG(rcg []*pb.RummyCardGroup, cardItem uint32) []*pb.RummyCardGroup {
	// 倒序处理牌序
	for key := len(rcg) - 1; key >= 0; key-- {
		for i := len(rcg[key].Cards) - 1; i >= 0; i-- {
			if rcg[key].Cards[i] == cardItem {
				rcg[key].Cards = append(rcg[key].Cards[:i], rcg[key].Cards[i+1:]...)

				if len(rcg[key].Cards) == 0 { //回收切片
					rcg = append(rcg[:key], rcg[key+1:]...)
				}
				return rcg
			}
		}
	}
	return rcg
}

// CheckRCG 对比上传 手牌和牌序数据有效性 返回计分
func CheckRCG(hardcards []uint32, rcg []*pb.RummyCardGroup) (result bool, score int64) {
	checkCards := make([]uint32, 0, len(hardcards))
	checkScore := make([][]uint32, 0, len(rcg))

	sort.Slice(hardcards, func(i, j int) bool {
		return hardcards[i] < hardcards[j]
	})

	for i := 0; i < len(rcg); i++ {
		if len(rcg[i].Cards) > 0 {
			checkCards = append(checkCards, rcg[i].Cards...)

			cards := make([]uint32, 0, len(rcg[i].Cards))
			cards = append(cards, rcg[i].Cards...)

			sort.Slice(cards, func(j, k int) bool {
				//return card.Card(cards[j]).Rank() < card.Card(cards[k]).Rank()
				return cards[j] < cards[k]
			})
			checkScore = append(checkScore, cards)
		}
	}

	score = GetCardValue(checkScore)

	sort.Slice(checkCards, func(i, j int) bool {
		return checkCards[i] < checkCards[j]
	})

	if len(checkCards) != len(hardcards) {
		mlog.Infof("错误: 用户手牌数据:%v \n 用户传入卡组数据 %v", card.CardList(hardcards), card.CardList(checkCards))
		result = false
		return
	}

	for i := 0; i < len(hardcards); i++ {
		if hardcards[i] != checkCards[i] {
			mlog.Infof("错误: 用户手牌数据:%v \n 用户传入卡组数据 %v", card.CardList(hardcards), card.CardList(checkCards))
			result = false
			return
		}
	}

	result = true
	return
}

// GetCardValue 获取手牌积分 首顺0分 满足首顺情况下赖顺0分 满足前二后顺子刻子0分
func GetCardValue(playerPrivate [][]uint32) (score int64) {
	isFirst := false
	secondCount := int64(0)
	isSecond := false
	otherCount := int64(0)

	for _, set := range playerPrivate {
		cardType, scoreItem := GetCardType(set)
		println(scoreItem)
		switch cardType {
		case pb.RummyGroupType_Rummy_GROUP_TYPE_INVALID:
			otherCount += int64(scoreItem)
		case pb.RummyGroupType_Rummy_GROUP_TYPE_PURE_SEQ:
			if isFirst { //两个纯顺
				isSecond = true
			}
			isFirst = true
		case pb.RummyGroupType_Rummy_GROUP_TYPE_SEQ:
			isSecond = true
			secondCount += int64(scoreItem)
		case pb.RummyGroupType_Rummy_GROUP_TYPE_SET:
			secondCount += int64(scoreItem)
		}
	}

	score += otherCount
	if !(isFirst && isSecond) {
		score += secondCount
	}

	if score >= 80 { // 封顶80分
		score = 80
	}
	return
}

// rummy比牌算法
func GetCardType(vals []uint32) (cardType pb.RummyGroupType, score uint32) {
	lenC := len(vals)
	cardType = pb.RummyGroupType_Rummy_GROUP_TYPE_INVALID

	isWild := false
	numWild := 0
	bit := uint32(0)
	color := pb.ColorType(0)

	realBit := uint32(0)
	realColor := pb.ColorType(0)
	minCardRank := math.MaxInt32     //除百搭以外最小牌
	minRealCardRank := math.MaxInt32 //除百搭以外最小牌

	hasAce := false
	for _, cardItem := range vals { // todo double A 10分

		scoreItem := uint32(card.Card(cardItem).Rank())

		if card.Card(cardItem).IsWild() { //joker todo wild
			isWild = true
			numWild++

			if minRealCardRank > int(scoreItem-1) {
				minRealCardRank = int(scoreItem - 1)
			}
		} else {
			if card.Card(cardItem).Rank() == pb.RankType_RANK_A {
				bit |= 1
				hasAce = true
				realBit |= 1
			}
			bit |= card.Card(cardItem).Bit()
			color |= card.Card(cardItem).Color()
			if scoreItem >= 10 {
				score += 10
			} else {
				score += scoreItem
			}

			if minCardRank > int(scoreItem-1) {
				minCardRank = int(scoreItem - 1)
			}

			if minRealCardRank > int(scoreItem-1) {
				minRealCardRank = int(scoreItem - 1)
			}
		}
		realBit |= card.Card(cardItem).Bit()
		realColor |= card.Card(cardItem).Color()
	}

	if lenC < 3 {
		return
	} //必须统计分数后再返回

	if realColor == pb.ColorType_Club || realColor == pb.ColorType_Diamond || realColor == pb.ColorType_Spade || realColor == pb.ColorType_Heart {
		if minRealCardRank < 12 { //k A Joker 不能为顺子
			//判断非赖顺
			if (realBit>>minRealCardRank) == (1<<lenC)-1 || (hasAce && realBit&^(1<<13) == (1<<lenC)-1) {
				// must first
				cardType = pb.RummyGroupType_Rummy_GROUP_TYPE_PURE_SEQ
				return
			}
		}
	}

	//如果带A 判断两次顺 结果取 1一次 || 14一次
	if color == pb.ColorType_Club || color == pb.ColorType_Diamond || color == pb.ColorType_Spade || color == pb.ColorType_Heart {
		//同花 判断首顺 赖顺
		if isWild {
			if lenC-numWild <= 1 {
				// second
				cardType = pb.RummyGroupType_Rummy_GROUP_TYPE_SEQ
			} else {
				unions := bits.OnesCount32(bit)
				gap := countZerosBetweenOnes(bit >> minCardRank)
				if hasAce {
					unions -= 1
					gap = min(gap, countZerosBetweenOnes(bit&^(1<<13)))
				}

				if unions != lenC-numWild { // 检查同花重复牌
					cardType = pb.RummyGroupType_Rummy_GROUP_TYPE_INVALID
					return
				}

				if gap <= numWild { // 间隙0小于癞子数
					cardType = pb.RummyGroupType_Rummy_GROUP_TYPE_SEQ
				}
			}
		}
	} else { //非同花
		if hasAce {
			bit -= 1
		}
		// 牌型4张以内 不是纯癞子 牌型是刻子 牌数-癞子数和花色数相等
		if lenC <= 4 && minCardRank != math.MaxInt32 && bit == 1<<minCardRank && lenC-numWild == bits.OnesCount32(uint32(color)) {
			cardType = pb.RummyGroupType_Rummy_GROUP_TYPE_SET
		} else if isWild && numWild == lenC { //3-5张牌 纯赖顺
			cardType = pb.RummyGroupType_Rummy_GROUP_TYPE_SEQ
		}
	}
	return
}

func countZerosBetweenOnes(n uint32) (result int) {
	count := 0
	inSequence := false // 是否已经遇到第一个 1

	for i := 0; i < 14; i++ {
		if (n & (1 << i)) != 0 { // 当前位是 1
			if inSequence && count > 0 {
				result += count // 记录 0 的数量
			}
			inSequence = true
			count = 0 // 重置计数
		} else if inSequence { // 当前位是 0，且之前遇到了 1
			count++
		}
	}

	return
}
