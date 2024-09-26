package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common/cfgData"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/module/entry"
)

type (
	PlayerItemKindLootGroupFun struct {
		*PlayerBagFun
		emKind cfgEnum.ESystemType
	}
)

func (this *PlayerItemKindLootGroupFun) Init(pFun *PlayerBagFun, emKind cfgEnum.ESystemType) {
	this.PlayerBagFun = pFun
	this.emKind = emKind
}

func (this *PlayerItemKindLootGroupFun) AddItem(head *pb.RpcHead, pbItem *pb.PBAddItemData) cfgEnum.ErrorCode {
	plog.Reward("head: %v, items: %v", head, pbItem)
	return plog.Print(this.AccountId, cfgData.GetLootGroupErrorCode(pbItem.Id), pbItem)
}

// 加掉落组
func (this *PlayerItemKindLootGroupFun) GetPbItem(itemId uint32, itemCount int64, emDoing pb.EmDoingType, params ...uint32) (arrPbItems []*pb.PBAddItemData) {
	arrPbItems = make([]*pb.PBAddItemData, 0)
	if itemCount <= 0 || itemId <= 0 {
		return
	}
	uRate := uint32(base.MIL_PERCENT)
	if len(params) > 0 {
		uRate = params[0]
	}
	//通过概率算修正
	if uRate != base.MIL_PERCENT {
		uTmpCount := int64(0)
		for i := int64(0); i < itemCount; i++ {
			if base.IsRadio(uRate) {
				uTmpCount++
			}
		}
		itemCount = uTmpCount
		if itemCount < 0 {
			return
		}
	}

	// 额外掉落词条(不能循环额外掉落)
	if emDoing != pb.EmDoingType_EDT_Entry {
		vals := entry.KeyValueToMap(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_ExtraLootDropProb), itemId)...)
		for lootGroupId, prob := range vals {
			if base.IsRadio(prob) {
				arrPbItems = append(arrPbItems, this.GetPbItem(lootGroupId, 1, pb.EmDoingType_EDT_Entry)...)
			}
		}
	}

	// 掉落组表掉落
	cfgLootGroup := cfgData.GetCfgLootGroup(itemId)
	if cfgLootGroup == nil || len(cfgLootGroup.LootPool) <= 0 {
		return
	}
	lootPool := make([]uint32, 0)
	lootPool = append(lootPool, cfgLootGroup.LootPool...)

	//需要增加抽奖up池子
	arrUpLoot := this.getPlayerSystemDrawFun().GetUpLoot(itemId)
	if len(arrUpLoot) > 0 {
		for _, lootid := range arrUpLoot {
			if !base.ArrayContainsValue(lootPool, lootid) {
				lootPool = append(lootPool, lootid)
			}
		}
	}

	switch cfgEnum.ELootRandType(cfgLootGroup.RandType) {
	case cfgEnum.ELootRandType_All:
		//获取词条
		vals := entry.KeyValueToMap(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_LootGroupLootRate), itemId)...)

		listLootId := []uint32{}
		mapQualityCount := make(map[uint32]map[uint32]uint32)
		for uCountIndex := int64(0); uCountIndex < itemCount; uCountIndex++ {
			if len(cfgLootGroup.MapProfHeroId) > 0 {
				lootPool = cfgLootGroup.LootPool
				//需要获取 全局前10个蓝色和全局前10个紫色，保两轮
				uPorf := this.getPlayerHeroFun().GetNextGlobalRandHeroProf(cfgLootGroup.ProfQuality, mapQualityCount)
				if uPorf >= 0 {
					plog.Trace(" rand hero uid:%d quality:%d prof:%d", this.AccountId, cfgLootGroup.ProfQuality, uPorf)
					if _, ok := cfgLootGroup.MapProfHeroId[uint32(uPorf)]; ok {
						lootPool = cfgLootGroup.MapProfHeroId[uint32(uPorf)]
					}
				}
			}
			wrand := base.NewWeightedRandom()
			for i := 0; i < len(lootPool); i++ {
				cfgLoot := cfgData.GetCfgLoot(lootPool[i])
				if cfgLoot == nil {
					plog.Info("AddLootGroup error lootid:%d", lootPool[i])
					continue
				}
				uGlobalRate := cfgLoot.LootGlobalProp
				if _, ok := vals[cfgLoot.Id]; ok {
					uGlobalRate += vals[cfgLoot.Id] / 10
				}

				wrand.Add(cfgLoot.Id, uGlobalRate)
			}
			listLootId = append(listLootId, wrand.GetRandomKey())
		}
		for _, lootId := range listLootId {
			cfgLoot := cfgData.GetCfgLoot(lootId)
			if cfgLoot == nil {
				plog.Info("AddLootGroup error lootid:%d", lootId)
				break
			}

			pbTmpItems := this.getLoot(cfgLoot, 1, emDoing)
			if len(pbTmpItems) <= 0 {
				continue
			}

			arrPbItems = append(arrPbItems, pbTmpItems...)
		}

	case cfgEnum.ELootRandType_Single:
		//获取词条
		vals := entry.KeyValueToMap(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_LootGroupLootRate), itemId)...)
		for i := 0; i < len(lootPool); i++ {
			cfgLoot := cfgData.GetCfgLoot(lootPool[i])
			if cfgLoot == nil {
				plog.Info("AddLootGroup error lootid:%d", lootPool[i])
				continue
			}

			uTmpAdd := uint32(0)
			for uCountIndex := int64(0); uCountIndex < itemCount; uCountIndex++ {
				if cfgLoot.LootSingleRate >= base.MIL_PERCENT || base.IsRadio(cfgLoot.LootSingleRate+vals[cfgLoot.Id]) {
					uTmpAdd++
				}
			}

			if uTmpAdd == 0 {
				continue
			}

			pbTmpItems := this.getLoot(cfgLoot, uTmpAdd, emDoing)
			if len(pbTmpItems) <= 0 {
				continue
			}

			arrPbItems = append(arrPbItems, pbTmpItems...)
		}
	case cfgEnum.ELootRandType_Group: //掉落组
		for i := 0; i < len(lootPool); i++ {
			pbTmpItems := this.GetPbItem(lootPool[i], itemCount, emDoing, base.MIL_PERCENT)
			if len(pbTmpItems) <= 0 {
				continue
			}

			arrPbItems = append(arrPbItems, pbTmpItems...)
		}
	}
	return
}

func (this *PlayerBagFun) getLoot(cfgLoot *cfgData.LootCfg, uCount uint32, emDoing pb.EmDoingType) (arrPbItems []*pb.PBAddItemData) {
	arrPbItems = make([]*pb.PBAddItemData, 0)
	if cfgLoot == nil {
		return
	}

	uAddCount := uint32(0)
	for i := uint32(0); i < uCount; i++ {
		uAddCount += base.RandRange(cfgLoot.MinCount, cfgLoot.MaxCount)
	}

	if uAddCount <= 0 {
		return
	}

	return this.getPlayerItemFun(cfgLoot.LootKind).GetPbItem(cfgLoot.LootId, int64(uAddCount), emDoing, cfgLoot.LootParam...)
}
