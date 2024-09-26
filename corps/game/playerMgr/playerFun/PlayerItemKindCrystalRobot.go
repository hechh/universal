package playerFun

import (
	"corps/base/cfgEnum"
	"corps/framework/plog"
	"corps/pb"
)

type (
	PlayerItemKindCrystalRobotFun struct {
		*PlayerBagFun
		emKind cfgEnum.ESystemType
	}
)

func (this *PlayerItemKindCrystalRobotFun) Init(pFun *PlayerBagFun, emKind cfgEnum.ESystemType) {
	this.PlayerBagFun = pFun
	this.emKind = emKind
}

func (this *PlayerItemKindCrystalRobotFun) GetPbItem(itemId uint32, itemCount int64, emDoing pb.EmDoingType, params ...uint32) (arrPbItems []*pb.PBAddItemData) {
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
func (this *PlayerItemKindCrystalRobotFun) AddItem(head *pb.RpcHead, pbItem *pb.PBAddItemData) cfgEnum.ErrorCode {
	plog.Reward("head: %v, items: %v", head, pbItem)
	return this.getPlayerCrystalFun().AddRobot(head, pbItem)
}
