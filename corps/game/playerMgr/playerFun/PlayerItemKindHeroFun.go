package playerFun

import (
	"corps/base/cfgEnum"
	"corps/framework/plog"
	"corps/pb"
)

type (
	PlayerItemKindHeroFun struct {
		*PlayerBagFun
		emKind cfgEnum.ESystemType
	}
)

func (this *PlayerItemKindHeroFun) Init(pFun *PlayerBagFun, emKind cfgEnum.ESystemType) {
	this.PlayerBagFun = pFun
	this.emKind = emKind
}

func (this *PlayerItemKindHeroFun) GetPbItem(itemId uint32, itemCount int64, emDoing pb.EmDoingType, params ...uint32) (arrPbItems []*pb.PBAddItemData) {
	arrPbItems = make([]*pb.PBAddItemData, 0)
	if itemCount <= 0 || itemId <= 0 {
		return
	}

	arrPbItems = append(arrPbItems, &pb.PBAddItemData{
		Kind:      uint32(this.emKind),
		Id:        itemId,
		Count:     itemCount,
		DoingType: emDoing,
		Params:    params,
	})

	return arrPbItems
}
func (this *PlayerItemKindHeroFun) AddItem(head *pb.RpcHead, pbItem *pb.PBAddItemData) cfgEnum.ErrorCode {
	plog.Reward("head: %v, items: %v", head, pbItem)
	uStar := uint32(1)
	if len(pbItem.Params) >= 1 {
		uStar = pbItem.Params[0]
	}

	return this.getPlayerHeroFun().AddHeros(head, pbItem.Id, uStar, uint32(pbItem.Count), pbItem.DoingType, true)
}
