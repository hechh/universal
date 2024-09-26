package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
)

type (
	PlayerSystemChargeCard struct {
		*PlayerSystemChargeFun
		mapCardInfo  map[uint32]*pb.PBChargeCard
		mapPrivilege map[uint32]uint32
	}
)

func (this *PlayerSystemChargeCard) Init(pFun *PlayerSystemChargeFun) {
	this.PlayerSystemChargeFun = pFun
	this.mapCardInfo = make(map[uint32]*pb.PBChargeCard)
	this.mapPrivilege = make(map[uint32]uint32)
}
func (this *PlayerSystemChargeCard) loadData(pbData *pb.PBPlayerSystemCharge) {
	this.mapCardInfo = make(map[uint32]*pb.PBChargeCard)
	for _, info := range pbData.CardList {
		this.mapCardInfo[info.CardType] = info
	}
}

func (this *PlayerSystemChargeCard) saveData(pbData *pb.PBPlayerSystemCharge) {
	for _, info := range this.mapCardInfo {
		pbData.CardList = append(pbData.CardList, info)
	}
}

// 加载完成需要读取有没有充值数据
func (this *PlayerSystemChargeCard) LoadComplete() {
	this.UpdatePrivilege()
}

// 是否能够购买
func (this *PlayerSystemChargeCard) canBuyProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode {
	if len(cfgCharge.Param) != 1 {
		return plog.Print(this.AccountId, cfgData.GetChargeErrorCode(cfgCharge.ProductID), cfgCharge.ProductID)
	}

	uCardType := cfgCharge.Param[0]
	cfgCard := cfgData.GetCfgChargeCardConfig(uCardType)
	if cfgCard == nil {
		return plog.Print(this.AccountId, cfgData.GetChargeCardConfigErrorCode(uCardType), cfgCharge.ProductID, uCardType)
	}

	if _, ok := this.mapCardInfo[uCardType]; ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_BuyTimesLimit, cfgCharge.ProductID, uCardType)
	}

	return cfgEnum.ErrorCode_Success
}

func (this *PlayerSystemChargeCard) OnChargeBuyProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode {
	if len(cfgCharge.Param) != 1 {
		return plog.Print(this.AccountId, cfgData.GetChargeErrorCode(cfgCharge.ProductID), cfgCharge.ProductID)
	}

	uCardType := cfgCharge.Param[0]
	cfgCard := cfgData.GetCfgChargeCardConfig(uCardType)
	if cfgCard == nil {
		return plog.Print(this.AccountId, cfgData.GetChargeCardConfigErrorCode(uCardType), cfgCharge.ProductID, uCardType)
	}

	pbCard, ok := this.mapCardInfo[uCardType]
	if ok {
		if pbCard.EndTime == 0 {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_ChargeRepeated, cfgCharge.ProductID, uCardType)
		}

		pbCard.EndTime += uint64(cfgCard.ContinueDay * 24 * 3600)
	} else {
		uCurTime := base.GetNow()

		this.mapCardInfo[uCardType] = &pb.PBChargeCard{
			CardType:  uCardType,
			BeginTime: uCurTime,
			PrizeTime: base.GetZeroTimestamp(uCurTime, 0),
		}

		if cfgCard.ContinueDay == 0 {
			this.mapCardInfo[uCardType].EndTime = 0
		} else {
			this.mapCardInfo[uCardType].EndTime = base.GetZeroTimestamp(uCurTime+uint64(cfgCard.ContinueDay*24*3600), 0) - 1
		}
	}

	//新增 加奖励
	this.getPlayerBagFun().AddArrItem(head, cfgCard.AddPrize, pb.EmDoingType_EDT_ChargeCard, true)

	this.UpdateSave(true)

	this.UpdatePrivilege()

	//通知客户端
	cluster.SendToClient(head, &pb.ChargeCardNewNotify{
		PacketHead: &pb.IPacket{},
		CardInfo:   this.mapCardInfo[uCardType],
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 获取特权
func (this *PlayerSystemChargeCard) GetPrivilege(uPrivilegeType uint32) uint32 {
	value, ok := this.mapPrivilege[uPrivilegeType]
	if !ok {
		return 0
	}
	return value
}

// 更新特权
func (this *PlayerSystemChargeCard) UpdatePrivilege() {
	mapNewPrivilege := make(map[uint32]uint32)
	for uCardType, _ := range this.mapCardInfo {
		cfgCard := cfgData.GetCfgChargeCardConfig(uCardType)
		if cfgCard == nil {
			continue
		}

		mapNewPrivilege = base.MergeMapU32U32(mapNewPrivilege, cfgCard.MapPrivilege)
	}

	this.mapPrivilege = mapNewPrivilege
	//比较map差异
	mapDiff := base.DiffMapU32U32(this.mapPrivilege, mapNewPrivilege)
	for key, _ := range mapDiff {
		this.getPlayerSystemCommonFun().UpdatePrivilege(cfgEnum.PrivilegeType(key))
	}
}

// 跨天 检查重置期数
func (this *PlayerSystemChargeCard) PassDay(isDay, isWeek, isMonth bool) {
	uNow := base.GetNow()
	//判断是否补发发奖
	pbNotify := &pb.ChargeCardUpdateNotify{
		PacketHead: &pb.IPacket{},
	}
	arrDel := make([]uint32, 0)
	for cardType, info := range this.mapCardInfo {
		if info.EndTime > 0 && info.EndTime <= uNow {
			arrDel = append(arrDel, cardType)
		}

		//判断未领取的天数
		uRealEndTime := uNow
		if info.EndTime > 0 && info.EndTime < uNow {
			uRealEndTime = info.EndTime + 1
		}

		diffDay := base.DiffDays(info.PrizeTime, uRealEndTime)
		if diffDay > 0 {
			cfgCard := cfgData.GetCfgChargeCardConfig(cardType)
			if cfgCard == nil {
				plog.Print(this.AccountId, cfgData.GetChargeCardConfigErrorCode(cardType), cardType)
				continue
			}
			info.PrizeTime = base.GetZeroTimestamp(uNow, 0)
			pbNotify.CardInfo = append(pbNotify.CardInfo, info)

			arrAddItem := make([]*common.ItemInfo, 0)
			for _, item := range cfgCard.DailyPrize {
				arrAddItem = append(arrAddItem, &common.ItemInfo{
					Id:     item.Id,
					Kind:   item.Kind,
					Count:  item.Count * int64(diffDay),
					Params: item.Params,
				})
			}

			//发邮件通知
			this.getPlayerMailFun().AddTempMail(&pb.RpcHead{Id: this.AccountId}, cfgEnum.EMailId_Card, pb.EmDoingType_EDT_ChargeCard, arrAddItem, cfgCard.Name, diffDay)
		}
	}

	//删除
	for _, cardType := range arrDel {
		delete(this.mapCardInfo, cardType)
		pbNotify.DelList = append(pbNotify.DelList, cardType)
	}

	if len(pbNotify.CardInfo) > 0 || len(pbNotify.DelList) > 0 {
		this.UpdateSave(true)
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbNotify, cfgEnum.ErrorCode_Success)
	}

	if len(arrDel) > 0 {
		this.UpdatePrivilege()
	}
}

// 充值卡领奖请求
func (this *PlayerSystemChargeCard) ChargeCardPrizeRequest(head *pb.RpcHead, pbRequest *pb.ChargeCardPrizeRequest) {
	uCode := this.ChargeCardPrize(head, pbRequest.CardType)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ChargeCardPrizeResponse{
			PacketHead: &pb.IPacket{},
			CardType:   pbRequest.CardType,
		}, uCode)
	}
}

// 充值卡领奖请求
func (this *PlayerSystemChargeCard) ChargeCardPrize(head *pb.RpcHead, uCardType uint32) cfgEnum.ErrorCode {
	pCardInfo, ok := this.mapCardInfo[uCardType]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, uCardType)
	}

	uNow := base.GetNow()
	if uNow < pCardInfo.PrizeTime {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HavePrize, uCardType)
	}

	cfgCard := cfgData.GetCfgChargeCardConfig(uCardType)
	if cfgCard == nil {
		return plog.Print(this.AccountId, cfgData.GetChargeCardConfigErrorCode(uCardType), uCardType)
	}

	//奖励
	pCardInfo.PrizeTime = base.GetZeroTimestamp(uNow, 1)
	this.getPlayerBagFun().AddArrItem(head, cfgCard.DailyPrize, pb.EmDoingType_EDT_ChargeCard, true)

	this.UpdateSave(true)

	cluster.SendToClient(head, &pb.ChargeCardPrizeResponse{
		PacketHead: &pb.IPacket{},
		CardType:   uCardType,
		PrizeTime:  pCardInfo.PrizeTime,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}
