package repository

import (
	"corps/common"
	"corps/pb"
	"corps/server/game/module/reward/internal/base"
)

type Common struct{}

func (d *Common) Probability(isAdd bool, probs map[common.ItemIndex]int64, items ...*pb.PBAddItemData) (rets []*pb.PBAddItemData) {
	// 分类统计合并
	tmps := map[common.ItemIndex]*pb.PBAddItemData{}
	for _, item := range items {
		index := common.ItemIndex{Kind: item.Kind, Id: item.Id}
		if _, ok := tmps[index]; !ok {
			tmps[index] = &pb.PBAddItemData{
				Kind:      item.Kind,
				Id:        item.Id,
				DoingType: item.DoingType,
			}
		}
		tmps[index].Count += item.Count
	}
	// 计算概率权重
	for index, item := range tmps {
		if isAdd {
			item.Count = base.AddProb(probs[index], item.Count)
		} else {
			item.Count = base.SubProb(probs[index], item.Count)
		}
		rets = append(rets, item)
	}
	return
}
