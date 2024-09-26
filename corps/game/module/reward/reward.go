package reward

import (
	"corps/base/cfgEnum"
	"corps/common"
	"corps/pb"
	"corps/server/game/module/reward/internal/base"
	"corps/server/game/module/reward/internal/manager"
	"corps/server/game/module/reward/internal/repository"
	"corps/server/game/module/reward/internal/service"
)

// 概率加成
func AddProbReward(prob int64, items ...*pb.PBAddItemData) []*pb.PBAddItemData {
	probs := map[common.ItemIndex]int64{}
	for _, item := range items {
		probs[base.ToItemIndexByData(item)] = prob
	}
	return service.AddProbReward(probs, items...)
}

// 概率加成 百分比
func AddProbCommonItem(prob uint32, items ...*common.ItemInfo) []*common.ItemInfo {
	arrInfo := make([]*common.ItemInfo, 0)
	for _, info := range items {
		uNewCount := info.Count * int64(prob) / 10000
		if uNewCount <= 0 {
			continue
		}
		tmpInfo := *info
		tmpInfo.Count = uNewCount
		arrInfo = append(arrInfo, &tmpInfo)

	}

	return arrInfo
}

// 概率减免
func SubProbReward(probs map[common.ItemIndex]int64, items ...*common.ItemInfo) []*common.ItemInfo {
	if len(probs) <= 0 {
		return items
	}
	return service.SubProbReward(probs, items...)
}

// 类型支持
func init() {
	manager.Register(&repository.Common{},
		uint32(cfgEnum.ESystemType_Item),
		uint32(cfgEnum.ESystemType_Hero),
		uint32(cfgEnum.ESystemType_Head),
		uint32(cfgEnum.ESystemType_HeadIcon),
		uint32(cfgEnum.ESystemType_Adverting),
		uint32(cfgEnum.ESystemType_Crystal),
		uint32(cfgEnum.ESystemType_CrystalRobot),
	)

	manager.Register(&repository.Equipment{}, uint32(cfgEnum.ESystemType_Equipment))
}

type itemInfo struct {
	Count uint32
	Item  *pb.PBAddItemData
}

func MergeBoxItems(args ...*pb.PBAddItemData) (rets []*pb.PBAddItemData) {
	tmps := map[common.ItemIndex]*itemInfo{}
	for _, item := range args {
		if item.Kind == uint32(cfgEnum.ESystemType_Equipment) {
			rets = append(rets, item)
			continue
		}
		// 获取下标
		var index common.ItemIndex
		switch item.Kind {
		case uint32(cfgEnum.ESystemType_Hero), uint32(cfgEnum.ESystemType_Crystal):
			index = common.ItemIndex{Kind: item.Kind, Id: item.Id, Star: item.Params[0]}
		default:
			index = common.ItemIndex{Kind: item.Kind, Id: item.Id}
		}
		// 合并统计
		val, ok := tmps[index]
		if !ok {
			tmps[index] = &itemInfo{Item: item, Count: 1}
			continue
		}
		if val.Count >= 10 {
			rets = append(rets, val.Item)
			val.Item = item
			val.Count = 1
			continue
		}
		val.Count++
		val.Item.Count += item.Count
	}
	for _, item := range tmps {
		if item.Count > 0 {
			rets = append(rets, item.Item)
		}
	}
	return
}
