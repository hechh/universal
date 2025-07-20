package test

import (
	"fmt"
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/library/util"
	"poker_server/server/room/internal/internal/rummy"
	"poker_server/server/room/internal/module/card"
	"sort"
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestWild(t *testing.T) {
	//fmt.Println(448>>7, 1<<3-1)
	for i := 1; i < 2; i++ {
		tmp := make([]uint32, 0, 4)
		for _, j := range []uint32{8, 7, 9} {
			tmp = append(tmp, (1<<(16+i))|j)
		}

		tmp[1] = card.Card(tmp[1]).AddWild()
		sort.Slice(tmp, func(j, k int) bool {
			return card.Card(tmp[j]).Rank() < card.Card(tmp[k]).Rank()
			//return tmp[j] < tmp[k]
		})
		cardType, _ := rummy.GetCardType(tmp)
		fmt.Printf("tmp : %v  cardtype:%v \n", card.CardList(tmp), cardType)
	}
}

func TestMatchId(t *testing.T) {
	types := util.DestructRoomId(30098325510)
	fmt.Printf("types:game_type:%v coin_type %v", types.GetGameType(), types.GetCoinType())
}

// TestScore 测试计分
func TestScoreHard(t *testing.T) {
	checkScore := make([][]uint32, 0, 7)

	for i := 3; i < 4; i++ {
		tmp := make([]uint32, 0, 4)
		for _, j := range []uint32{9, 12, 13, 14, 14} {
			tmp = append(tmp, (1<<(16+i))|j)
		}

		fmt.Printf("tmp : %v \n", card.CardList(tmp))
		checkScore = append(checkScore, tmp)
	}

	for i := 2; i < 3; i++ {
		tmp := make([]uint32, 0, 4)
		for _, j := range []uint32{10} {
			tmp = append(tmp, (1<<(16+i))|j)
		}

		fmt.Printf("tmp : %v \n", card.CardList(tmp))
		checkScore = append(checkScore, tmp)
	}

	for i := 1; i < 2; i++ {
		tmp := make([]uint32, 0, 4)
		for _, j := range []uint32{7, 8, 11} {
			tmp = append(tmp, (1<<(16+i))|j)
		}

		fmt.Printf("tmp : %v \n", card.CardList(tmp))
		checkScore = append(checkScore, tmp)
	}

	for i := 2; i < 3; i++ {
		tmp := make([]uint32, 0, 4)
		for _, j := range []uint32{6, 10} {
			tmp = append(tmp, (1<<(16+i))|j)
		}

		fmt.Printf("tmp : %v \n", card.CardList(tmp))
		checkScore = append(checkScore, tmp)
	}

	for i := 2; i < 3; i++ {
		tmp := make([]uint32, 0, 4)
		for _, j := range []uint32{10} {
			tmp = append(tmp, (1<<(16+i))|j)
		}

		fmt.Printf("tmp : %v \n", card.CardList(tmp))
		checkScore = append(checkScore, tmp)
	}

	score := rummy.GetCardValue(checkScore)
	fmt.Printf("score: %v  \n", score)
}

// TestScore 测试计分
func TestScore(t *testing.T) {
	checkScore := make([][]uint32, 0, 7)

	for i := 0; i < 1; i++ {
		tmp := make([]uint32, 0, 4)
		for j := uint32(4); j <= 6; j++ { //4-6
			tmp = append(tmp, (1<<(16+i))|j)
		}

		fmt.Printf("tmp : %v \n", card.CardList(tmp))
		checkScore = append(checkScore, tmp)
	}

	for i := 1; i < 2; i++ {
		tmp := make([]uint32, 0, 4)
		for j := uint32(9); j <= 10; j++ { //4-6
			tmp = append(tmp, (1<<(16+i))|j)
		}

		tmp = append(tmp, (1<<(31))|15) //joker
		fmt.Printf("tmp : %v \n", card.CardList(tmp))
		checkScore = append(checkScore, tmp)
	}

	for i := 0; i < 1; i++ {
		tmp := make([]uint32, 0, 4)
		for j := uint32(4); j <= 5; j++ { //4-6
			tmp = append(tmp, (1<<(16+i))|j)
		}

		fmt.Printf("tmp : %v \n", card.CardList(tmp))
		checkScore = append(checkScore, tmp)
	}

	for i := 3; i < 4; i++ {
		tmp := make([]uint32, 0, 4)
		for _, j := range []uint32{3, 7, 8, 14} {
			tmp = append(tmp, (1<<(16+i))|j)
		}

		fmt.Printf("tmp : %v \n", card.CardList(tmp))
		checkScore = append(checkScore, tmp)
	}

	for i := 2; i < 3; i++ {
		tmp := make([]uint32, 0, 4)
		for _, j := range []uint32{3} {
			tmp = append(tmp, (1<<(16+i))|j)
		}

		fmt.Printf("tmp : %v \n", card.CardList(tmp))
		checkScore = append(checkScore, tmp)
	}

	score := rummy.GetCardValue(checkScore)
	fmt.Printf("score: %v  \n", score)
}

func TestSeq(t *testing.T) {
	cards := make([]uint32, 0, 4)
	cards = append(cards, (1<<(16+1))|13) //k
	cards = append(cards, (1<<(16+1))|14) //j
	cards = append(cards, (1<<(31))|15)   //joker

	cardType, _ := rummy.GetCardType(cards)
	fmt.Printf("cardType: %v  %v \n", cardType, card.CardList(cards))
}

func TestPureSeq(t *testing.T) {
	cards := make([]uint32, 0, 4)
	cards = append(cards, (1<<(16+1))|7)
	for i := uint32(1); i <= 1; i++ {
		for j := uint32(12); j <= 13; j++ { //ace
			cards = append(cards, (1<<(16+i))|j)
		}
	}

	cards[0] = card.Card(cards[0]).AddWild()

	cardType, _ := rummy.GetCardType(cards)
	fmt.Printf("cardType: %v  %v \n", cardType, card.CardList(cards))
}

func TestAceSet(t *testing.T) {
	cards := make([]uint32, 0, 4)
	for i := uint32(1); i <= 3; i++ {
		for j := uint32(14); j <= 14; j++ { //ace
			cards = append(cards, (1<<(16+i))|j)
		}
	}
	cards = append(cards, (1<<(31))|15) //joker

	fmt.Printf("cards: %v \n", card.CardList(cards))
	cardType, _ := rummy.GetCardType(cards)
	fmt.Printf("cardType: %v \n", cardType)
}

func TestCopy(t *testing.T) {
	src := []uint32{1, 2, 3, 5, 6, 7, 8, 9}
	dst := make([]uint32, len(src))
	copy(dst, src)

	fmt.Printf("<UNK>: %v,%v", src, dst)
}

func TestFuck(t *testing.T) {
	data := []byte{
		10, 169, 1, 8, 254, 85, 18, 163, 1, 104, 101, 97, 100, 58, 110, 111, 100, 101, 95, 116,
		121, 112, 101, 58, 78, 111, 100, 101, 84, 121, 112, 101, 82, 111, 111, 109, 32, 110, 111,
		100, 101, 95, 105, 100, 58, 49, 32, 97, 99, 116, 111, 114, 95, 110, 97, 109, 101, 58, 34,
		82, 117, 109, 109, 121, 71, 97, 109, 101, 34, 32, 102, 117, 110, 99, 95, 110, 97, 109, 101,
		58, 34, 70, 105, 120, 67, 97, 114, 100, 34, 32, 97, 99, 116, 111, 114, 95, 105, 100, 58, 51,
		48, 48, 56, 49, 53, 52, 56, 50, 57, 48, 32, 114, 111, 117, 116, 101, 114, 58, 123, 103, 97,
		116, 101, 58, 49, 32, 114, 111, 111, 109, 58, 49, 32, 103, 97, 109, 101, 58, 49, 125, 44,
		32, 102, 114, 97, 117, 100, 32, 99, 97, 114, 100, 103, 114, 111, 117, 112, 32, 108, 111,
		115, 101, 32, 116, 104, 105, 115, 32, 103, 97, 109, 101, 16, 32, 26, 13, 10, 9, 138, 128,
		8, 139, 128, 8, 140, 128, 8, 16, 2, 26, 15, 10, 11, 137, 128, 32, 138, 128, 32, 143, 128,
		128, 128, 8, 16, 3, 26, 15, 10, 11, 140, 128, 32, 141, 128, 32, 130, 128, 160, 128, 8, 16,
		3, 26, 16, 10, 12, 140, 128, 4, 134, 128, 16, 136, 128, 16, 136, 128, 16, 16, 1,
	}
	pb := &pb.RummyFixCardRsp{}
	if err := proto.Unmarshal(data, pb); err != nil {
		fmt.Printf("解析protobuf失败: %v", err)
	}
	fmt.Printf("解析protobuf: %v", pb)
}

func TestCardSet(t *testing.T) {
	cards := []uint32{}
	for i := uint32(1); i <= 3; i++ {
		for j := uint32(3); j <= 3; j++ {
			cards = append(cards, (1<<(16+i))|j)
		}
	}
	//cards = append(cards, (1<<(16+1))|3)
	card_type, _ := rummy.GetCardType(cards)
	fmt.Printf("card_type: %v cards %v \n", card_type, card.CardList(cards))
}

func TestCard(t *testing.T) {
	mlog.Init("mason", 1, "debug", "./log")
	rcg := []*pb.RummyCardGroup{
		{
			Cards: []uint32{65544, 65545, 65546},
		},
		{
			Cards: []uint32{2147614722, 524301, 524302},
		},
		{
			Cards: []uint32{131077, 131082, 131084},
		},
		{
			Cards: []uint32{2147614722, 2148007938},
		},
		{
			Cards: []uint32{262151, 262152, 131075},
		},
	}
	rcg = rummy.DelCardRCG(rcg, 65542)
	for i := range rcg {
		mlog.Infof("%v", len(rcg[i].Cards))
	}
	reala, score := rummy.CheckRCG([]uint32{65544, 131084, 2147614722, 2148007938, 65545, 65546, 524291, 524301, 524302, 131077, 131082, 262151, 262152, 131075}, rcg)

	fmt.Printf("real: %v score %v \n", reala, score)
}
