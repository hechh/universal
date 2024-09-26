package playerFun

import (
	"corps/base/cfgEnum"
	"corps/framework/plog"
	"corps/pb"
)

type (
	PlayerItemKindHeadFun struct {
		*PlayerBagFun
		emKind cfgEnum.ESystemType
	}
	PlayerItemKindHeadIconFun struct {
		*PlayerBagFun
		emKind cfgEnum.ESystemType
	}
)

func (this *PlayerItemKindHeadFun) Init(pFun *PlayerBagFun, emKind cfgEnum.ESystemType) {
	this.PlayerBagFun = pFun
	this.emKind = emKind
}

func (this *PlayerItemKindHeadFun) GetPbItem(itemId uint32, itemCount int64, emDoing pb.EmDoingType, params ...uint32) (arrPbItems []*pb.PBAddItemData) {
	arrPbItems = make([]*pb.PBAddItemData, 0)
	if itemCount <= 0 || itemId <= 0 {
		return
	}

	arrPbItems = append(arrPbItems, &pb.PBAddItemData{
		Kind:      uint32(this.emKind),
		Id:        itemId,
		DoingType: emDoing,
	})

	return arrPbItems
}
func (this *PlayerItemKindHeadFun) AddItem(head *pb.RpcHead, pbItem *pb.PBAddItemData) cfgEnum.ErrorCode {
	plog.Reward("head: %v, items: %v", head, pbItem)
	return this.getPlayerSystemCommonFun().AddHead(head, pbItem.Id, pbItem.DoingType)
}
func (this *PlayerItemKindHeadIconFun) Init(pFun *PlayerBagFun, emKind cfgEnum.ESystemType) {
	this.PlayerBagFun = pFun
	this.emKind = emKind
}

func (this *PlayerItemKindHeadIconFun) GetPbItem(itemId uint32, itemCount int64, emDoing pb.EmDoingType, params ...uint32) (arrPbItems []*pb.PBAddItemData) {
	arrPbItems = make([]*pb.PBAddItemData, 0)
	if itemCount <= 0 || itemId <= 0 {
		return
	}

	arrPbItems = append(arrPbItems, &pb.PBAddItemData{
		Kind:      uint32(this.emKind),
		Id:        itemId,
		DoingType: emDoing,
	})

	return arrPbItems
}
func (this *PlayerItemKindHeadIconFun) AddItem(head *pb.RpcHead, pbItem *pb.PBAddItemData) cfgEnum.ErrorCode {
	plog.Reward("head: %v, items: %v", head, pbItem)
	return this.getPlayerSystemCommonFun().AddHeadIcon(head, pbItem.Id, pbItem.DoingType)
}
