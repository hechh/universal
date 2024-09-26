package service

import (
	"corps/common"
	"corps/pb"
	"corps/server/game/module/reward/internal/base"
	"corps/server/game/module/reward/internal/manager"
)

func AddProbReward(probs map[common.ItemIndex]int64, items ...*pb.PBAddItemData) (rets []*pb.PBAddItemData) {
	if len(probs) <= 0 {
		return
	}
	// 分类
	classType, _ := manager.GetClass(items...)
	for typ, vals := range classType {
		tmps := manager.GetIReward(typ).Probability(true, probs, vals...)
		if len(tmps) > 0 {
			rets = append(rets, tmps...)
		}
	}
	return
}

func SubProbReward(probs map[common.ItemIndex]int64, items ...*common.ItemInfo) (rets []*common.ItemInfo) {
	tmps := map[common.ItemIndex]int64{}
	for _, item := range items {
		index := base.ToItemIndexByItem(item)
		if _, ok := probs[index]; !ok {
			rets = append(rets, item)
			continue
		}
		tmps[index] += item.Count
	}

	// 计算减免
	for index, count := range tmps {
		newItem := &common.ItemInfo{
			Kind:  index.Kind,
			Id:    index.Id,
			Count: base.SubProb(probs[index], count),
		}
		if newItem.Count > 0 {
			rets = append(rets, newItem)
		}
	}
	return
}
