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
	PlayerActivityOpenServerGift struct {
		emType cfgEnum.EActivityType
		*PlayerSystemActivityFun
		mapData map[uint32]*PlayerOpenServerGift //开服特惠礼包
	}
	PlayerOpenServerGift struct {
		*pb.PBActivityOpenServerGift
		mapGift map[uint32]*PlayerOpenServerGiftInfo //开启的礼包
	}
	PlayerOpenServerGiftInfo struct {
		*pb.PBOpenServerGiftInfo
		mapStage map[uint32]*pb.PBU32U32 // 购买的档次
	}
)

func init() {
	RegisterActivity(cfgEnum.EActivityType_OpenServerGift, func() IPlayerSystemActivityFun { return new(PlayerActivityOpenServerGift) })
}

func (this *PlayerActivityOpenServerGift) Init(emType cfgEnum.EActivityType, pFun *PlayerSystemActivityFun) {
	this.emType = emType
	this.PlayerSystemActivityFun = pFun
	this.mapData = make(map[uint32]*PlayerOpenServerGift)
}

func (this *PlayerActivityOpenServerGift) LoadPlayerDBFinish() {

}

func (this *PlayerActivityOpenServerGift) LoadData(pbData *pb.PBPlayerSystemActivity) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemActivity{}
	}

	this.mapData = make(map[uint32]*PlayerOpenServerGift)
	for _, info := range pbData.OpenServerGiftList {
		this.mapData[info.Id] = &PlayerOpenServerGift{
			PBActivityOpenServerGift: info,
			mapGift:                  make(map[uint32]*PlayerOpenServerGiftInfo),
		}

		for _, giftInfo := range info.GiftList {
			this.mapData[info.Id].mapGift[giftInfo.GiftId] = &PlayerOpenServerGiftInfo{
				PBOpenServerGiftInfo: giftInfo,
				mapStage:             make(map[uint32]*pb.PBU32U32),
			}

			for _, stageInfo := range giftInfo.StageList {
				this.mapData[info.Id].mapGift[giftInfo.GiftId].mapStage[stageInfo.Key] = stageInfo
			}
		}
	}

	this.UpdateSave(true)
}

// 存储到数据库
func (this *PlayerActivityOpenServerGift) SaveData(pbData *pb.PBPlayerSystemActivity) {
	for _, info := range this.mapData {
		pbData.OpenServerGiftList = append(pbData.OpenServerGiftList, info.PBActivityOpenServerGift)
	}
}
func (this *PlayerActivityOpenServerGift) Del(uID uint32) {
	if _, ok := this.mapData[uID]; !ok {
		plog.Error("(this *PlayerActivityOpenServerGift) Del %d %d %d", uID)
		return
	}
	delete(this.mapData, uID)
	this.UpdateSave(true)
}
func (this *PlayerActivityOpenServerGift) GetRed(uID uint32) bool {
	pActivity, ok := this.mapData[uID]
	if !ok {
		return false
	}

	//有未领取的每日奖励 通知红点
	if base.GetNow() >= pActivity.NextDailyPrizeTime {
		return true
	}

	return false
}
func (this *PlayerActivityOpenServerGift) Open(uID uint32) (arrReturn []string) {
	arrReturn = make([]string, 0)
	pActivity, ok := this.mapData[uID]
	if !ok {
		return arrReturn
	}

	listStageGift := make([]*cfgData.OpenServerStageConfigCfg, 0)
	listGift := make([]*cfgData.OpenServerGiftConfigCfg, 0)
	for giftId, _ := range pActivity.mapGift {
		cfgGift := cfgData.GetCfgOpenServerGiftConfig(giftId)
		if cfgGift != nil {
			listGift = append(listGift, cfgGift)
		}

		mapTmp := cfgData.GetCfgOpenServerStageConfigByGiftId(giftId)
		for _, cfg := range mapTmp {
			listStageGift = append(listStageGift, cfg)
		}
	}

	if len(listGift) <= 0 && len(listStageGift) <= 0 {
		return arrReturn
	}

	byDatGift, err := json.Marshal(listGift)
	if err == nil {
		arrReturn = append(arrReturn, string(byDatGift))
	}

	byDatStage, err := json.Marshal(listStageGift)
	if err == nil {
		arrReturn = append(arrReturn, string(byDatStage))
	}

	return arrReturn
}
func (this *PlayerActivityOpenServerGift) Add(uID uint32, uBeginTime uint64, uEndTime uint64) {
	if _, ok := this.mapData[uID]; ok {
		plog.Error("(this *PlayerSystemChargeFun) NewChargeGift repeated %d %d %d", uID, uBeginTime, uEndTime)
		return
	}

	this.mapData[uID] = &PlayerOpenServerGift{
		PBActivityOpenServerGift: &pb.PBActivityOpenServerGift{
			Id:                 uID,
			BeginTime:          uBeginTime,
			EndTime:            uEndTime,
			NextDailyPrizeTime: 0,
		},
		mapGift: make(map[uint32]*PlayerOpenServerGiftInfo),
	}

	mapCfg := cfgData.GetCfgGiftConfigBySId(uID)
	for _, cfg := range mapCfg {
		if cfg.SystemType > 0 && !this.getPlayerSystemCommonFun().CheckSystemTypeOpen(cfgEnum.ESystemUnlockType(cfg.SystemType)) {
			continue
		}

		this.AddGift(&pb.RpcHead{Id: this.AccountId}, cfg, false)
	}

	//通知客户端
	pbNotify := &pb.ActivityDataNewNotify{
		PacketHead:   &pb.IPacket{},
		ActivityType: uint32(this.emType),
		Info:         &pb.PBPlayerSystemActivity{},
	}
	pbNotify.Info.OpenServerGiftList = append(pbNotify.Info.OpenServerGiftList, this.mapData[uID].PBActivityOpenServerGift)
	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbNotify, cfgEnum.ErrorCode_Success)
	this.UpdateSave(true)
}

// 是否能够购买
func (this *PlayerActivityOpenServerGift) canBuyProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode {
	if len(cfgCharge.Param) != 2 {
		return plog.Print(this.AccountId, cfgData.GetChargeErrorCode(cfgCharge.ProductID), cfgCharge.ProductID)
	}

	uGiftId := cfgCharge.Param[0]
	uStage := cfgCharge.Param[1]
	cfgStage := cfgData.GetCfgOpenServerStageConfig(uGiftId, uStage)
	if cfgStage == nil {
		return plog.Print(this.AccountId, cfgData.GetOpenServerStageConfigErrorCode(uGiftId), cfgCharge.ProductID, uGiftId)
	}
	//充值礼包不能买
	if cfgStage.ProductId != cfgCharge.ProductID {
		return plog.Print(this.AccountId, cfgData.GetOpenServerStageConfigErrorCode(uGiftId), uGiftId)
	}

	cfgActivity := cfgData.GetCfgOpenServerGiftConfig(cfgStage.GiftId)
	if cfgActivity == nil {
		return plog.Print(this.AccountId, cfgData.GetOpenServerGiftConfigErrorCode(cfgStage.GiftId), cfgCharge.ProductID, uGiftId)
	}
	pAcitivity, ok := this.mapData[cfgActivity.Sid]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, cfgCharge.ProductID, uGiftId)
	}

	pGift, ok := pAcitivity.mapGift[uGiftId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, cfgCharge.ProductID, uGiftId)
	}

	//判断前置是否买完
	if uStage > 1 {
		cfgPreStage := cfgData.GetCfgOpenServerStageConfig(uGiftId, uStage-1)
		if cfgPreStage == nil {
			return plog.Print(this.AccountId, cfgData.GetOpenServerStageConfigErrorCode(uGiftId), cfgCharge.ProductID, uGiftId)
		}

		pbPreBuy, ok := pGift.mapStage[cfgPreStage.Id]
		if !ok || pbPreBuy.Value < cfgPreStage.LimitCount {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedCondition, cfgCharge.ProductID, uGiftId)
		}
	}

	//购买上线
	pbBuy, okBuy := pGift.mapStage[cfgStage.Id]
	if okBuy && pbBuy.Value >= cfgStage.LimitCount {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_BuyTimesLimit, cfgCharge.ProductID, uGiftId)
	}

	return cfgEnum.ErrorCode_Success
}

// 购买直购礼包
func (this *PlayerActivityOpenServerGift) OnChargeBuyProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode {
	uCode := this.canBuyProduct(head, cfgCharge)
	if uCode != cfgEnum.ErrorCode_Success {
		return uCode
	}
	uGiftId := cfgCharge.Param[0]
	uStage := cfgCharge.Param[1]
	cfgStage := cfgData.GetCfgOpenServerStageConfig(uGiftId, uStage)
	if cfgStage == nil {
		return plog.Print(this.AccountId, cfgData.GetOpenServerStageConfigErrorCode(uGiftId), cfgCharge.ProductID, uGiftId)
	}

	cfgActivity := cfgData.GetCfgOpenServerGiftConfig(cfgStage.GiftId)
	if cfgActivity == nil {
		return plog.Print(this.AccountId, cfgData.GetOpenServerGiftConfigErrorCode(cfgStage.GiftId), cfgCharge.ProductID, uGiftId)
	}
	pAcitivity, ok := this.mapData[cfgActivity.Sid]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, cfgCharge.ProductID, uGiftId)
	}

	pGift, ok := pAcitivity.mapGift[uGiftId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, cfgCharge.ProductID, uGiftId)
	}

	//购买
	if _, okBuy := pGift.mapStage[cfgStage.Id]; !okBuy {
		pGift.mapStage[cfgStage.Id] = &pb.PBU32U32{Key: cfgStage.Id, Value: 1}
		pGift.StageList = append(pGift.StageList, pGift.mapStage[cfgStage.Id])
	} else {
		pGift.mapStage[cfgStage.Id].Value++
	}

	//加奖励
	this.getPlayerBagFun().AddArrItem(head, cfgStage.AddPrize, pb.EmDoingType_EDT_ChargeGift, true)

	this.UpdateSave(true)

	//通知客户端
	cluster.SendToClient(head, &pb.OpenServerGiftBuyNotify{
		PacketHead: &pb.IPacket{},
		Id:         cfgActivity.Sid,
		GiftId:     uGiftId,
		BuyInfo:    pGift.mapStage[cfgStage.Id],
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 加载完成需要读取有没有充值数据
func (this *PlayerActivityOpenServerGift) LoadComplete() {
	mapAllCfg := cfgData.GetAllCfgGiftConfigBySystemTypes()
	for systemType, listCfg := range mapAllCfg {
		if systemType > 0 && !this.getPlayerSystemCommonFun().CheckSystemTypeOpen(cfgEnum.ESystemUnlockType(systemType)) {
			continue
		}

		for _, cfg := range listCfg {
			this.AddGift(&pb.RpcHead{Id: this.AccountId}, cfg, true)
		}
	}

}

// 系统开启
func (this *PlayerActivityOpenServerGift) OnSystemOpenTypes(head *pb.RpcHead, arrTypes []uint32) {
	for _, systemType := range arrTypes {
		//检查pb开启
		listCfg := cfgData.GetCfgGiftConfigBySystemTypes(systemType)
		for _, cfgBp := range listCfg {
			this.AddGift(head, cfgBp, true)
		}
	}
}
func (this *PlayerActivityOpenServerGift) AddGift(head *pb.RpcHead, cfgGift *cfgData.OpenServerGiftConfigCfg, bSend bool) {
	pActivity, ok := this.mapData[cfgGift.Sid]
	if !ok {
		return
	}

	if _, ok := pActivity.mapGift[cfgGift.Id]; ok {
		return
	}

	uNow := base.GetNow()
	uEndTime := base.GetZeroTimestamp(uNow, int32(cfgGift.ContinueTime)) - 1
	if uEndTime > pActivity.EndTime {
		uEndTime = pActivity.EndTime
	}
	pActivity.mapGift[cfgGift.Id] = &PlayerOpenServerGiftInfo{
		PBOpenServerGiftInfo: &pb.PBOpenServerGiftInfo{
			BeginTime: base.GetNow(),
			EndTime:   uEndTime,
			GiftId:    cfgGift.Id,
		},
		mapStage: make(map[uint32]*pb.PBU32U32),
	}
	pActivity.GiftList = append(pActivity.GiftList, pActivity.mapGift[cfgGift.Id].PBOpenServerGiftInfo)

	//通知客户端
	if bSend {
		cluster.SendToClient(head, &pb.OpenServerGiftNewNotify{
			PacketHead: &pb.IPacket{},
			SId:        cfgGift.Sid,
			GiftInfo:   pActivity.mapGift[cfgGift.Id].PBOpenServerGiftInfo,
		}, cfgEnum.ErrorCode_Success)
	}
}

// 活动免费奖励返回
func (this *PlayerActivityOpenServerGift) FreePrize(head *pb.RpcHead, uID uint32) cfgEnum.ErrorCode {
	cfgActivity := cfgData.GetCfgOpenServerActivityConfig(uID)
	if cfgActivity == nil {
		return plog.Print(this.AccountId, cfgData.GetOpenServerActivityConfigErrorCode(uID), uID)
	}

	pActivity, ok := this.mapData[uID]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, uID)
	}

	uCurTime := base.GetNow()
	if uCurTime < pActivity.NextDailyPrizeTime {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HavePrize, uID)
	}

	this.getPlayerBagFun().AddArrItem(head, cfgActivity.AddPrize, pb.EmDoingType_EDT_Activity, true)
	pActivity.NextDailyPrizeTime = base.GetZeroTimestamp(uCurTime, 1)
	return cfgEnum.ErrorCode_Success
}
