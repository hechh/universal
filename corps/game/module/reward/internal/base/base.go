package base

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/pb"
	"math"
)

func ToItemIndex(item *pb.PBAddItem) common.ItemIndex {
	switch item.Kind {
	case uint32(cfgEnum.ESystemType_Equipment):
		return common.ItemIndex{Kind: item.Kind, Id: item.Id, Quality: item.Params[0], Star: item.Params[1]}
	case uint32(cfgEnum.ESystemType_Hero):
		return common.ItemIndex{Kind: item.Kind, Id: item.Id, Star: item.Params[0]}
	}
	return common.ItemIndex{Kind: item.Kind, Id: item.Id}
}

func ToItemIndexByData(item *pb.PBAddItemData) common.ItemIndex {
	switch item.Kind {
	case uint32(cfgEnum.ESystemType_Equipment):
		return common.ItemIndex{Kind: item.Kind, Id: item.Equipment.Id, Quality: item.Equipment.Quality, Star: item.Equipment.Star}
	case uint32(cfgEnum.ESystemType_Hero):
		return common.ItemIndex{Kind: item.Kind, Id: item.Id, Star: item.Params[0]}
	}
	return common.ItemIndex{Kind: item.Kind, Id: item.Id}
}

func ToItemIndexByItem(item *common.ItemInfo) common.ItemIndex {
	switch item.Kind {
	case uint32(cfgEnum.ESystemType_Equipment):
		return common.ItemIndex{Kind: item.Kind, Id: item.Id, Quality: item.Params[0], Star: item.Params[1]}
	case uint32(cfgEnum.ESystemType_Hero):
		return common.ItemIndex{Kind: item.Kind, Id: item.Id, Star: item.Params[0]}
	}
	return common.ItemIndex{Kind: item.Kind, Id: item.Id}
}

func AddProb(prob int64, count int64) int64 {
	return count + int64(math.Ceil(float64(prob*count)/base.MIL_PERCENT))
}

func SubProb(prob int64, count int64) int64 {
	diff := count - int64(math.Ceil(float64(prob*count)/base.MIL_PERCENT))
	if diff <= 0 {
		diff = 0
	}
	return diff
}
