package base

import (
	"fmt"
	"poker_server/common/pb"
	"strings"
)

var (
	color_name = map[pb.TexasColor]string{
		pb.TexasColor_TC_DIAMOND: "♦",
		pb.TexasColor_TC_CLUB:    "♣",
		pb.TexasColor_TC_HEART:   "♥",
		pb.TexasColor_TC_SPADE:   "♠",
	}
	rank_name = map[pb.TexasRank]string{
		pb.TexasRank_TR_1:  "A",
		pb.TexasRank_TR_2:  "2",
		pb.TexasRank_TR_3:  "3",
		pb.TexasRank_TR_4:  "4",
		pb.TexasRank_TR_5:  "5",
		pb.TexasRank_TR_6:  "6",
		pb.TexasRank_TR_7:  "7",
		pb.TexasRank_TR_8:  "8",
		pb.TexasRank_TR_9:  "9",
		pb.TexasRank_TR_10: "10",
		pb.TexasRank_TR_J:  "J",
		pb.TexasRank_TR_Q:  "Q",
		pb.TexasRank_TR_K:  "K",
		pb.TexasRank_TR_A:  "A",
	}
)

type Card uint32

type CardList []uint32

func (d Card) Color() pb.TexasColor {
	return pb.TexasColor(d >> 16)
}

func (d Card) Rank() pb.TexasRank {
	return pb.TexasRank(d & 0x0F)
}

func (d Card) Bit() uint32 {
	return 1 << (d.Rank() - 1)
}

func (t Card) Value() uint32 {
	return uint32(t)
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", rank_name[c.Rank()], color_name[c.Color()])
}

func (d CardList) String() string {
	strs := []string{}
	for _, v := range d {
		strs = append(strs, Card(v).String())
	}
	return strings.Join(strs, ",")
}
