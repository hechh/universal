package entry

import (
	"corps/base/cfgEnum"
	"corps/common"
	"corps/pb"
	"corps/server/game/module/entry/internal/condition"
	"corps/server/game/module/entry/internal/entity"
	"corps/server/game/module/entry/internal/manager"
)

func Init() {}

func init() {
	// 注册条件类型
	manager.RegisterCondition(&condition.EmptyCondition{}, uint32(cfgEnum.EEntryParamTypeCondition_None))
	manager.RegisterCondition(&condition.BaseCondition{}, uint32(cfgEnum.EEntryParamTypeCondition_Times))
	manager.RegisterCondition(&condition.DayCondition{}, uint32(cfgEnum.EEntryParamTypeCondition_Day))

	// 注册效果参数类型
	manager.RegisterEntity(entity.NewEmptyEntity, uint32(cfgEnum.EEntryParamTypeEffect_None))
	manager.RegisterEntity(entity.NewOneEntity, uint32(cfgEnum.EEntryParamTypeEffect_One))
	manager.RegisterEntity(entity.NewTwoEntity, uint32(cfgEnum.EEntryParamTypeEffect_Two))
	manager.RegisterEntity(entity.NewTwoArrayEntity, uint32(cfgEnum.EEntryParamTypeEffect_TwoArray))
	manager.RegisterEntity(entity.NewThreeArrayEntity, uint32(cfgEnum.EEntryParamTypeEffect_ThreeArray))
}

func ToItem(vals ...*pb.EntryEffectValue) (rets map[common.ItemIndex]int64) {
	rets = make(map[common.ItemIndex]int64)
	for _, val := range vals {
		rets[common.ItemIndex{Kind: val.List[0], Id: val.List[1]}] += int64(val.List[2])
	}
	return
}

func KeyValueToMap(vals ...*pb.EntryEffectValue) map[uint32]uint32 {
	ret := make(map[uint32]uint32)
	for _, val := range vals {
		if len(val.List) != 2 {
			continue
		}
		ret[val.List[0]] += val.List[1]
	}
	return ret
}
func KeyValueToDMap(vals ...*pb.EntryEffectValue) map[uint32]map[uint32]uint32 {
	ret := make(map[uint32]map[uint32]uint32)
	for _, val := range vals {
		if len(val.List) != 3 {
			continue
		}

		if _, ok := ret[val.List[0]]; !ok {
			ret[val.List[0]] = make(map[uint32]uint32)
		}
		ret[val.List[0]][val.List[1]] += val.List[2]
	}
	return ret
}

func ToValue(vals ...*pb.EntryEffectValue) uint32 {
	if len(vals) <= 0 {
		return 0
	}
	return vals[0].List[0]
}
