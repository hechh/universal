package playerFun

import (
	"corps/base/cfgEnum"
	"corps/framework/plog"
	"corps/pb"
)

type (
	PlayerItemKindItemFun struct {
		*PlayerBagFun
		emKind cfgEnum.ESystemType
	}
)

func (this *PlayerItemKindItemFun) Init(pFun *PlayerBagFun, emKind cfgEnum.ESystemType) {
	this.PlayerBagFun = pFun
	this.emKind = emKind
}

func (this *PlayerItemKindItemFun) GetPbItem(itemId uint32, itemCount int64, emDoing pb.EmDoingType, params ...uint32) (arrPbItems []*pb.PBAddItemData) {
	arrPbItems = make([]*pb.PBAddItemData, 0)
	if itemId <= 0 {
		return
	}

	arrPbItems = append(arrPbItems, &pb.PBAddItemData{
		Kind:      uint32(this.emKind),
		Id:        itemId,
		Count:     itemCount,
		DoingType: emDoing,
	})

	return arrPbItems
}
func (this *PlayerItemKindItemFun) AddItem(head *pb.RpcHead, pbItem *pb.PBAddItemData) (uCode cfgEnum.ErrorCode) {
	plog.Reward("head: %v, items: %v", head, pbItem)
	uCode = this.addBagItem(head, pbItem) //加道具
	// 触发道具消耗效果
	if uCode == cfgEnum.ErrorCode_Success {
		if pbItem.Count >= 0 {
			//成就触发
			this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ItemAdd, uint32(pbItem.Count), pbItem.Id)
		} else if pbItem.Count < 0 {
			//成就触发
			this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ItemConsume, uint32(pbItem.Count*-1), pbItem.Id)
		}

	}
	return
}
