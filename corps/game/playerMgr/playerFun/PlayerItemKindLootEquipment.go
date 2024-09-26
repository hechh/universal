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
	PlayerItemKindLootEquipmentFun struct {
		*PlayerBagFun
		emKind cfgEnum.ESystemType
	}
)

func (this *PlayerItemKindLootEquipmentFun) Init(pFun *PlayerBagFun, emKind cfgEnum.ESystemType) {
	this.PlayerBagFun = pFun
	this.emKind = emKind
}

func (this *PlayerItemKindLootEquipmentFun) GetPbItem(itemId uint32, itemCount int64, emDoing pb.EmDoingType, params ...uint32) (arrPbItems []*pb.PBAddItemData) {
	arrPbItems = make([]*pb.PBAddItemData, 0)
	if itemCount <= 0 || itemId <= 0 {
		return
	}

	cfgLootEquip := cfgData.GetCfgLootEquipment(itemId)
	if cfgLootEquip == nil {
		return
	}

	for i := uint32(0); i < uint32(itemCount); i++ {
		pEquipment := this.getLootEquipmentPB(cfgLootEquip, emDoing)
		if pEquipment == nil {
			plog.Print(this.AccountId, cfgData.GetLootEquipmentErrorCode(itemId), cfgLootEquip)
			continue
		}

		arrPbItems = append(arrPbItems, &pb.PBAddItemData{
			Kind:      uint32(cfgEnum.ESystemType_Equipment),
			DoingType: emDoing,
			Equipment: pEquipment,
		})
	}

	return arrPbItems
}

func (this *PlayerItemKindLootEquipmentFun) AddItem(head *pb.RpcHead, pbItem *pb.PBAddItemData) cfgEnum.ErrorCode {
	plog.Reward("head: %v, items: %v", head, pbItem)
	return plog.Print(this.AccountId, cfgData.GetLootEquipmentErrorCode(pbItem.Id), pbItem)
}

func (this *PlayerItemKindLootEquipmentFun) getLootEquipmentPB(cfgLootEquipment *cfgData.LootEquipmentCfg, emDoing pb.EmDoingType) *pb.PBEquipment {
	if cfgLootEquipment == nil {
		return nil
	}

	//随机ID
	uEquipId := base.RandArrayKey(cfgLootEquipment.EquipmentList)

	//如果是战斗获得，取配置
	uStar := uint32(1)
	mapRandQuality := make(map[uint32]int32)
	for key, value := range cfgLootEquipment.MapQualityRate {
		mapRandQuality[key] += int32(value)
	}

	if emDoing == pb.EmDoingType_EDT_Offline || emDoing == pb.EmDoingType_EDT_BattleHook || emDoing == pb.EmDoingType_EDT_BattleNormal {
		mapTecEffect := this.getPlayerSystemHookTechFun().GetHookTechEffect(cfgEnum.TechEffectType_AddHookEquipMaxStar)
		if len(mapTecEffect) > 0 {
			uStar = uint32(mapTecEffect[0])
		}

		//挂机科技增加
		mapAddQualityProb := this.getPlayerSystemHookTechFun().GetHookTechEffect(cfgEnum.TechEffectType_AddEquipQualityRate)
		for key, value := range mapAddQualityProb {
			mapRandQuality[key] = value
		}

		//需要判断词条 增加装备品质概率 词条加成 todo 品质 概率万分比
		mapEntryProb := entry.KeyValueToMap(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_BattleEnd_EquipmentDropProb), uint32(cfgEnum.EBattleType_BattleNormal))...)
		for key, value := range mapEntryProb {
			if _, ok := mapRandQuality[key]; ok {
				mapRandQuality[key] += mapRandQuality[key] * int32(value) / base.MIL_PERCENT
			}
		}

	} else {
		wrandStar := base.NewWeightedRandom()
		for key, value := range cfgLootEquipment.MapStarRate {
			wrandStar.Add(key, value)
		}
		uStar = wrandStar.GetRandomKey()
	}

	wrandQuality := base.NewWeightedRandom()
	for key, value := range mapRandQuality {
		if value > 0 {
			wrandQuality.Add(key, uint32(value))
		}
	}

	uQuality := wrandQuality.GetRandomKey()
	pbEquip := this.GetPlayerEquipmentFun().GetNewEquipment(uEquipId, uQuality, uStar, emDoing)
	return pbEquip
}
