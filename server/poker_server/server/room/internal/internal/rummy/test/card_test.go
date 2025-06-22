package test

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/rummy"
	"poker_server/server/room/internal/module/card"
	"testing"
)

// TestScore 测试计分
func TestScore(t *testing.T) {
	checkScore := make([][]uint32, 0, 7)

	for i := 0; i < 1; i++ {
		tmp := make([]uint32, 0, 4)
		for j := uint32(8); j <= 11; j++ { //8-j
			tmp = append(tmp, (1<<(16+i))|j)
		}

		fmt.Printf("tmp : %v \n", card.CardList(tmp))
		checkScore = append(checkScore, tmp)
	}

	tmp := make([]uint32, 0, 4)
	for i := 0; i < 3; i++ {
		tmp = append(tmp, (1<<(16+i))|7)
	}
	fmt.Printf("tmp : %v \n", card.CardList(tmp))
	checkScore = append(checkScore, tmp)

	//for i := 0; i < 1; i++ {
	//	tmp := make([]uint32, 0, 4)
	//
	//	for j := uint32(10); j <= 12; j++ { //8-j
	//		tmp = append(tmp, (1<<(16+i))|j)
	//	}
	//	fmt.Printf("tmp : %v \n", card.CardList(tmp))
	//	checkScore = append(checkScore, tmp)
	//}

	tmp = make([]uint32, 0, 4)
	tmp = append(tmp, (1<<(16+2))|2)
	tmp = append(tmp, (1<<(16+2))|9)
	tmp = append(tmp, (1<<(16+2))|10)
	//tmp[2] = card.Card(tmp[2]).AddWild()
	checkScore = append(checkScore, tmp)
	fmt.Printf("tmp : %v \n", card.CardList(tmp))

	tmp = make([]uint32, 0, 4)
	tmp = append(tmp, (1<<(16))|12)
	tmp = append(tmp, (1<<(16))|13)
	tmp = append(tmp, (1<<(16+2))|12)
	//tmp[2] = card.Card(tmp[2]).AddWild()
	checkScore = append(checkScore, tmp)
	fmt.Printf("tmp : %v \n", card.CardList(tmp))

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
	mlog.Init("mason", "debug", "./log")
	rcg := []*pb.RummyCardGroup{
		{
			Cards: []uint32{262155, 65540, 65542, 65546, 65547, 524294, 524295},
		},
		{
			Cards: []uint32{524296, 524297, 131075},
		},
		{
			Cards: []uint32{131077, 131082, 131084},
		},
		{
			Cards: []uint32{65542},
		},
	}
	rcg = rummy.DelCardRCG(rcg, 65542)
	for i := range rcg {
		mlog.Infof("%v", len(rcg[i].Cards))
	}
	rummy.CheckRCG([]uint32{262155, 65540, 65542, 65546, 65547, 524294, 524295, 524296, 524297, 131075, 131077, 131082, 131084, 65542}, rcg)
}
