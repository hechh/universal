package playerFun

import (
	"corps/base/cfgEnum"
	"corps/common/cfgData"
	"corps/framework/plog"
	"corps/pb"
)

type (
	PlayerItemKindEquipmentFun struct {
		*PlayerBagFun
		emKind cfgEnum.ESystemType
	}
)

func (this *PlayerItemKindEquipmentFun) Init(pFun *PlayerBagFun, emKind cfgEnum.ESystemType) {
	this.PlayerBagFun = pFun
	this.emKind = emKind
}

func (this *PlayerItemKindEquipmentFun) GetPbItem(itemId uint32, itemCount int64, emDoing pb.EmDoingType, params ...uint32) (arrPbItems []*pb.PBAddItemData) {
	arrPbItems = make([]*pb.PBAddItemData, 0)
	uQuality := uint32(1)
	uStar := uint32(1)
	if len(params) >= 2 {
		uQuality = params[0]
		uStar = params[1]
	}
	if itemCount <= 0 || itemId <= 0 {
		return
	}

	for i := uint32(0); i < uint32(itemCount); i++ {
		pEquipment := this.GetPlayerEquipmentFun().GetNewEquipment(itemId, uQuality, uStar, emDoing)
		if pEquipment == nil {
			plog.Print(this.AccountId, cfgData.GetEquipmentErrorCode(itemId), itemId, uQuality, uStar)
			continue
		}

		arrPbItems = append(arrPbItems, &pb.PBAddItemData{
			Kind:      uint32(this.emKind),
			DoingType: emDoing,
			Equipment: pEquipment,
		})
	}

	return arrPbItems
}
func (this *PlayerItemKindEquipmentFun) AddItem(head *pb.RpcHead, pbItem *pb.PBAddItemData) cfgEnum.ErrorCode {
	uCode, _ := this.GetPlayerEquipmentFun().AddPBEquipment(head, pbItem.Equipment, true, pbItem.DoingType)
	plog.Reward("head: %v, code: %d, items: %v", head, uCode, pbItem)
	return uCode
}
