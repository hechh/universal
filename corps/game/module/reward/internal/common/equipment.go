package common

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common/cfgData"
	"corps/framework/plog"
	"corps/pb"
)

type Equipment struct{}

type equipmentIndex struct {
	Id      uint32
	Quality uint32
	Star    uint32
}

func (d *Equipment) Handle(prob int64, items ...*pb.PBAddItemData) (rets []*pb.PBAddItemData) {
	// 分类
	tmps := map[equipmentIndex][]*pb.PBAddItemData{}
	for _, item := range items {
		index := equipmentIndex{item.Equipment.Id, item.Equipment.Quality, item.Equipment.Star}
		if _, ok := tmps[index]; !ok {
			tmps[index] = []*pb.PBAddItemData{}
		}
		tmps[index] = append(tmps[index], item)
	}
	// 倍乘
	for index, items := range tmps {
		rets = append(rets, items...)
		count := int(prob) * len(items)
		for i := 0; i < count; i++ {
			rets = append(rets, &pb.PBAddItemData{
				Kind:      items[0].Kind,
				DoingType: items[0].DoingType,
				Equipment: GetNewEquipment(index.Id, index.Quality, index.Star),
			})
		}
	}
	return
}

// 获取一个新的装备结构
func GetNewEquipment(uId uint32, uQuality uint32, uStar uint32) *pb.PBEquipment {
	cfgEquipment := cfgData.GetCfgEquipment(uId)
	//获取品质表
	cfgQuality := cfgData.GetCfgEquipmentQuality(uQuality)
	//获取星级表
	cfgStar := cfgData.GetCfgEquipmentStar(uStar)
	pEquipment := &pb.PBEquipment{
		Id:              uId,
		Quality:         uQuality,
		Star:            uStar,
		Sn:              0,
		MainProp:        &pb.PBEquipmentProp{},
		EquipProfession: base.CFG_DEFAULT_VALUE,
	}
	//获取主词条数值
	pEquipment.MainProp.PropId = cfgEquipment.MainProp
	uMainValueRate := uint32(0)
	if pEquipment.MainProp.PropId == 0 {
		arrcfgProp := cfgData.GetCfgRandEquipmentProp(cfgEquipment.Stage, uint32(cfgEnum.EHydraPropType_MainProp), cfgEquipment.PartType, uint32(1))
		if arrcfgProp != nil && len(arrcfgProp) == 1 {
			pEquipment.MainProp.Value = cfgEquipment.MainPropValue
			pEquipment.MainProp.PropId = arrcfgProp[0].PropId
			uMainValueRate = arrcfgProp[0].ValueRate
		} else {
			plog.Info("AddEquipment error")
		}
	} else {
		cfgProp := cfgData.GetCfgEquipmentProp(uint32(cfgEnum.EHydraPropType_MainProp), pEquipment.MainProp.PropId)
		pEquipment.MainProp.Value = cfgEquipment.MainPropValue
		uMainValueRate = cfgProp.ValueRate
	}

	//评分 (2000+ 2000*0.164 )*1.1 / 50
	pEquipment.MainProp.Score = base.RAND.RandUint(cfgQuality.MinScoreRate, cfgQuality.MaxScoreRate)
	pEquipment.MainProp.Value = uint32(uint64(pEquipment.MainProp.Value) * uint64(base.MIL_PERCENT+pEquipment.MainProp.Score) * uint64(cfgStar.AddRate) / uint64(uMainValueRate) / uint64(base.MIL_PERCENT*base.MIL_PERCENT))

	//获取次词条数值
	pEquipment.MinorPropList = getEquipProp(cfgEquipment, cfgQuality, cfgStar, cfgEquipment.MinorProp, cfgQuality.MinorPropCount, cfgEnum.EHydraPropType_MinorProp, cfgEquipment.MinorPropValue)
	pEquipment.VicePropList = getEquipProp(cfgEquipment, cfgQuality, cfgStar, cfgEquipment.ViceProp, cfgQuality.VicePropCount, cfgEnum.EHydraPropType_ViceProp, cfgEquipment.VicePropValue)
	return pEquipment
}

// 获取装备词条 先取固定 再随机剩余
func getEquipProp(cfgEquip *cfgData.EquipmentCfg, cfgQuality *cfgData.EquipmentQualityCfg, cfgStar *cfgData.EquipmentStarCfg, arrProp []uint32, uPropCount uint32, emPropType cfgEnum.EHydraPropType, uPropValue uint32) []*pb.PBEquipmentProp {
	arrReturn := make([]*pb.PBEquipmentProp, 0)
	if len(arrProp) > 0 {
		for i := 0; i < len(arrProp); i++ {
			if i >= int(uPropCount) {
				break
			}

			cfgProp := cfgData.GetCfgEquipmentProp(uint32(emPropType), arrProp[i])
			if cfgProp == nil {
				continue
			}

			pbProp := &pb.PBEquipmentProp{
				PropId: arrProp[i],
				Score:  base.RAND.RandUint(cfgQuality.MinScoreRate, cfgQuality.MaxScoreRate),
			}

			//无价值 取值
			if emPropType == cfgEnum.EHydraPropType_ViceProp {
				pbProp.Value = uPropValue
				pbProp.Value = pbProp.Value * pbProp.Score / (cfgProp.ValueRate * 10000)
			} else {
				var uValue uint64 = uint64(uPropValue) * (uint64(10000 + pbProp.Score)) * uint64(cfgStar.AddRate)
				uValue = uValue / uint64(cfgProp.ValueRate) / 10000 / 10000
				pbProp.Value = uint32(uValue)
			}

			arrReturn = append(arrReturn, pbProp)
		}

		uPropCount -= uint32(len(arrReturn))
	}

	if uPropCount > 0 {
		arrcfgProp := cfgData.GetCfgRandEquipmentProp(cfgEquip.Stage, uint32(emPropType), cfgEquip.PartType, uPropCount)
		if arrcfgProp != nil {
			for i := 0; i < len(arrcfgProp); i++ {
				pbProp := &pb.PBEquipmentProp{
					PropId: arrcfgProp[i].PropId,
					Score:  base.RAND.RandUint(cfgQuality.MinScoreRate, cfgQuality.MaxScoreRate),
				}

				//无价值 取值
				if emPropType == cfgEnum.EHydraPropType_ViceProp {
					pbProp.Value = uPropValue
					pbProp.Value = pbProp.Value * pbProp.Score / (arrcfgProp[i].ValueRate * 10000)
				} else {
					var uValue uint64 = uint64(uPropValue) * (uint64(10000 + pbProp.Score)) * uint64(cfgStar.AddRate)
					uValue = uValue / uint64(arrcfgProp[i].ValueRate) / 10000 / 10000
					pbProp.Value = uint32(uValue)
				}

				arrReturn = append(arrReturn, pbProp)
			}
		}
	}

	return arrReturn
}
