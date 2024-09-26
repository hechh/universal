package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common/cfgData"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"encoding/json"
)

type (
	PlayerActivityChargeGift struct {
		emType cfgEnum.EActivityType
		*PlayerSystemActivityFun
		mapChargeGift map[uint32]*PlayerChargeGift //直购礼包
	}
	PlayerChargeGift struct {
		*pb.PBActivityChargeGift
		mapBuy map[uint32]*pb.PBU32U32 // 购买的
	}
)

func init() {
	RegisterActivity(cfgEnum.EActivityType_ChargeBuy, func() IPlayerSystemActivityFun { return new(PlayerActivityChargeGift) })
}

func (this *PlayerActivityChargeGift) Init(emType cfgEnum.EActivityType, pFun *PlayerSystemActivityFun) {
	this.emType = emType
	this.PlayerSystemActivityFun = pFun
	this.mapChargeGift = make(map[uint32]*PlayerChargeGift)
}

func (this *PlayerActivityChargeGift) LoadPlayerDBFinish() {

}
func (this *PlayerActivityChargeGift) FreePrize(head *pb.RpcHead, uID uint32) cfgEnum.ErrorCode {
	return cfgEnum.ErrorCode_NoData
}
func (this *PlayerActivityChargeGift) LoadData(pbData *pb.PBPlayerSystemActivity) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemActivity{}
	}

	this.mapChargeGift = make(map[uint32]*PlayerChargeGift)
	for _, info := range pbData.GiftList {
		this.mapChargeGift[info.Id] = &PlayerChargeGift{
			PBActivityChargeGift: info,
			mapBuy:               make(map[uint32]*pb.PBU32U32),
		}

		for _, buyInfo := range info.BuyList {
			this.mapChargeGift[info.Id].mapBuy[buyInfo.Key] = buyInfo
		}
	}

	this.UpdateSave(true)
}

func (this *PlayerActivityChargeGift) LoadComplete() {
}

// 存储到数据库
func (this *PlayerActivityChargeGift) SaveData(pbData *pb.PBPlayerSystemActivity) {
	for _, info := range this.mapChargeGift {
		pbData.GiftList = append(pbData.GiftList, info.PBActivityChargeGift)
	}
}
func (this *PlayerActivityChargeGift) Del(uID uint32) {
	if _, ok := this.mapChargeGift[uID]; !ok {
		plog.Error("(this *PlayerActivityChargeGift) Del %d %d %d", uID)
		return
	}
	delete(this.mapChargeGift, uID)
	this.UpdateSave(true)
}
func (this *PlayerActivityChargeGift) GetRed(uID uint32) bool {

	return false
}
func (this *PlayerActivityChargeGift) Open(uID uint32) []string {
	listCfg := cfgData.GetAllCfgChargeGiftConfig(uID)
	if len(listCfg) <= 0 {
		return nil
	}
	byDat, err := json.Marshal(listCfg)
	if err != nil {
		return nil
	}

	return []string{string(byDat)}
}
func (this *PlayerActivityChargeGift) Add(uID uint32, uBeginTime uint64, uEndTime uint64) {
	if _, ok := this.mapChargeGift[uID]; ok {
		plog.Error("(this *PlayerSystemChargeFun) NewChargeGift repeated %d %d %d", uID, uBeginTime, uEndTime)
		return
	}

	this.mapChargeGift[uID] = &PlayerChargeGift{
		PBActivityChargeGift: &pb.PBActivityChargeGift{
			Id:        uID,
			BeginTime: uBeginTime,
			EndTime:   uEndTime,
		},
		mapBuy: make(map[uint32]*pb.PBU32U32),
	}

	//通知客户端
	pbNotify := &pb.ActivityDataNewNotify{
		PacketHead:   &pb.IPacket{},
		ActivityType: uint32(this.emType),
		Info:         &pb.PBPlayerSystemActivity{},
	}
	pbNotify.Info.GiftList = append(pbNotify.Info.GiftList, this.mapChargeGift[uID].PBActivityChargeGift)
	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbNotify, cfgEnum.ErrorCode_Success)
}

// 购买直购礼包
func (this *PlayerActivityChargeGift) OnChargeBuyGiftProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode {
	if len(cfgCharge.Param) != 1 {
		return plog.Print(this.AccountId, cfgData.GetChargeErrorCode(cfgCharge.ProductID), cfgCharge.ProductID)
	}

	uGiftId := cfgCharge.Param[0]
	cfg := cfgData.GetCfgChargeGiftConfig(uGiftId)
	if cfg == nil {
		return plog.Print(this.AccountId, cfgData.GetChargeGiftConfigErrorCode(uGiftId), uGiftId)
	}

	//充值礼包不能买
	if cfg.ProductId != cfgCharge.ProductID {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ChargeGiftPID, uGiftId)
	}

	pGift, ok := this.mapChargeGift[cfg.Sid]
	if !ok {
		return plog.Print(this.AccountId, cfgData.GetChargeGiftConfigErrorCode(uGiftId), uGiftId)
	}

	pbBuy, okBuy := pGift.mapBuy[uGiftId]
	if okBuy && pbBuy.Value >= cfg.LimitCount {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_BuyTimesLimit, uGiftId)
	}

	//加奖励
	this.getPlayerBagFun().AddArrItem(head, cfg.AddPrize, pb.EmDoingType_EDT_ChargeGift, true)

	if !okBuy {
		pGift.mapBuy[uGiftId] = &pb.PBU32U32{Key: uGiftId, Value: 1}
		pGift.PBActivityChargeGift.BuyList = append(pGift.PBActivityChargeGift.BuyList, pGift.mapBuy[uGiftId])
	} else {
		pGift.mapBuy[uGiftId].Value++
	}

	this.UpdateSave(true)

	//通知客户端
	cluster.SendToClient(head, &pb.ChargeGiftBuyNotify{
		PacketHead: &pb.IPacket{},
		Id:         cfg.Sid,
		BuyInfo:    pGift.mapBuy[uGiftId],
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 是否能够购买
func (this *PlayerActivityChargeGift) canBuyProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode {
	if len(cfgCharge.Param) != 1 {
		return plog.Print(this.AccountId, cfgData.GetChargeErrorCode(cfgCharge.ProductID), cfgCharge.ProductID)
	}

	uGiftId := cfgCharge.Param[0]
	cfg := cfgData.GetCfgChargeGiftConfig(uGiftId)
	if cfg == nil {
		return plog.Print(this.AccountId, cfgData.GetChargeGiftConfigErrorCode(uGiftId), uGiftId)
	}

	pGift, ok := this.mapChargeGift[cfg.Sid]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, uGiftId)
	}

	pbBuy, okBuy := pGift.mapBuy[uGiftId]
	if okBuy && pbBuy.Value >= cfg.LimitCount {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_BuyTimesLimit, uGiftId)
	}

	//判断天数
	if base.DiffDays(pGift.BeginTime, base.GetNow())+1 < cfg.Day {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, uGiftId)
	}

	return cfgEnum.ErrorCode_Success
}
