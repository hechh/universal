package domain

import (
	"corps/common"
	"corps/pb"
)

const (
	PERCENT = 100
)

type IReward interface {
	Probability(isAdd bool, prob map[common.ItemIndex]int64, items ...*pb.PBAddItemData) []*pb.PBAddItemData // 返回概率加成奖励
}
