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
	PlayerSystemChargeBP struct {
		*PlayerSystemChargeFun
		mapBPInfo map[uint32]*PlayerPBInfo
	}

	PlayerPBInfo struct {
		*pb.PBBPInfo
		mapStageInfo map[uint32]*pb.PBBPStageInfo
	}
)

func (this *PlayerSystemChargeBP) Init(pFun *PlayerSystemChargeFun) {
	this.PlayerSystemChargeFun = pFun
	this.mapBPInfo = make(map[uint32]*PlayerPBInfo)
}
func (this *PlayerSystemChargeBP) loadData(pbData *pb.PBPlayerSystemCharge) {
	this.mapBPInfo = make(map[uint32]*PlayerPBInfo)
	for _, info := range pbData.BPList {
		this.mapBPInfo[info.BPType] = &PlayerPBInfo{
			PBBPInfo:     info,
			mapStageInfo: make(map[uint32]*pb.PBBPStageInfo),
		}

		for _, stageInfo := range info.StageList {
			this.mapBPInfo[info.BPType].mapStageInfo[stageInfo.StageId] = stageInfo
		}
	}

}

func (this *PlayerSystemChargeBP) saveData(pbData *pb.PBPlayerSystemCharge) {
	for _, info := range this.mapBPInfo {
		pbData.BPList = append(pbData.BPList, info.PBBPInfo)
	}
}

// 加载完成需要读取有没有充值数据
func (this *PlayerSystemChargeBP) LoadComplete() {
	//默认解锁BP
	arrPb := cfgData.GetAllCfgBpActConfigBySystemTypes()
	for systemId, listBp := range arrPb {
		if systemId > 0 && !this.getPlayerSystemCommonFun().CheckSystemTypeOpen(cfgEnum.ESystemUnlockType(systemId)) {
			continue
		}

		//检查pb开启
		for _, cfgBp := range listBp {
			this.AddNewPb(&pb.RpcHead{Id: this.AccountId}, cfgBp)
		}
	}
}

// 系统开启
func (this *PlayerSystemChargeBP) OnSystemOpenTypes(head *pb.RpcHead, arrIds []uint32) {
	for _, systemId := range arrIds {
		//检查pb开启
		listBp := cfgData.GetCfgBpActConfigBySystemTypes(systemId)
		for _, cfgBp := range listBp {
			this.AddNewPb(head, cfgBp)
		}
	}
}
func (this *PlayerSystemChargeBP) DelBP(head *pb.RpcHead, cfgBp *cfgData.BpActConfigCfg) {
	pBPInfo, ok := this.mapBPInfo[cfgBp.BpType]
	if !ok {
		return
	}

	for index, info := range pBPInfo.StageList {
		if info.StageId == cfgBp.Stage {
			pBPInfo.StageList = pBPInfo.StageList
			pBPInfo.StageList = append(pBPInfo.StageList[:index], pBPInfo.StageList[index+1:]...)
		}
	}
	delete(pBPInfo.mapStageInfo, cfgBp.Stage)
}

// 新增BP
func (this *PlayerSystemChargeBP) AddNewPb(head *pb.RpcHead, cfgBp *cfgData.BpActConfigCfg) {
	if cfgBp == nil {
		return
	}

	//只能加最高期数的
	_, ok := this.mapBPInfo[cfgBp.BpType]
	if ok {
		if cfgBp.Stage <= this.mapBPInfo[cfgBp.BpType].MaxStage {
			return
		}
	} else {
		this.mapBPInfo[cfgBp.BpType] = &PlayerPBInfo{
			PBBPInfo: &pb.PBBPInfo{
				BPType: cfgBp.BpType,
				Value:  this.GetPbValue(cfgBp, 0, 0),
			},
			mapStageInfo: make(map[uint32]*pb.PBBPStageInfo),
		}
	}

	uCurTime := base.GetNow()
	pbStage := &pb.PBBPStageInfo{
		StageId:   cfgBp.Stage,
		BeginTime: uCurTime,
	}

	if cfgBp.AliveDay > 0 {
		pbStage.EndTime = base.GetZeroTimestamp(uCurTime, 0) + uint64(cfgBp.AliveDay*24*3600) - 1
	}

	this.mapBPInfo[cfgBp.BpType].StageList = append(this.mapBPInfo[cfgBp.BpType].StageList, pbStage)
	this.mapBPInfo[cfgBp.BpType].mapStageInfo[pbStage.StageId] = pbStage

	this.mapBPInfo[cfgBp.BpType].MaxStage = base.MaxUint32(this.mapBPInfo[cfgBp.BpType].MaxStage, cfgBp.Stage)
	this.UpdateSave(true)

	//同步客户端
	cluster.SendToClient(head, &pb.BPNewNotify{
		PacketHead: &pb.IPacket{},
		BPInfo:     this.mapBPInfo[cfgBp.BpType].PBBPInfo,
	}, cfgEnum.ErrorCode_Success)
}

// 新增BPstage
func (this *PlayerSystemChargeBP) AddNewPbStage(head *pb.RpcHead, cfgBp *cfgData.BpActConfigCfg, bSend bool) *pb.PBBPStageInfo {
	if cfgBp == nil {
		return nil
	}

	pBPInfo, ok := this.mapBPInfo[cfgBp.BpType]
	if !ok {
		return nil
	}

	if _, ok := pBPInfo.mapStageInfo[cfgBp.Stage]; ok {
		return nil
	}

	uCurTime := base.GetNow()
	pbStage := &pb.PBBPStageInfo{
		StageId:   cfgBp.Stage,
		BeginTime: uCurTime,
	}

	if cfgBp.AliveDay > 0 {
		pbStage.EndTime = base.GetZeroTimestamp(uCurTime, 0) + uint64(cfgBp.AliveDay*24*3600) - 1
	}

	pBPInfo.StageList = append(this.mapBPInfo[cfgBp.BpType].StageList, pbStage)
	pBPInfo.mapStageInfo[pbStage.StageId] = pbStage

	this.UpdateSave(true)

	if bSend {
		cluster.SendToClient(head, &pb.BPNewStageNotify{
			PacketHead: &pb.IPacket{},
			BPType:     pBPInfo.BPType,
			StageList:  []*pb.PBBPStageInfo{pbStage},
		}, cfgEnum.ErrorCode_Success)
	}

	return pbStage
}

// 获取pb值
func (this *PlayerSystemChargeBP) GetPbValue(cfgBp *cfgData.BpActConfigCfg, uBPValue uint32, uAdd uint32) uint32 {
	if cfgBp == nil {
		return 0
	}
	uValue := uint32(0)
	switch cfgEnum.AchieveType(cfgBp.AchieveType) {
	case cfgEnum.AchieveType_BattleMap:
		mapId, stageId := this.getPlayerSystemBattleFun().GetFinishMapIdAndStageId(pb.EmBattleType(cfgBp.AchieveParams[0]))
		uValue = mapId*1000 + stageId
	default:
		if cfgBp.IsTotal > 0 {
			uValue = this.getPlayerSystemTaskFun().GetAchieveValue(cfgBp.AchieveType, cfgBp.AchieveParams...)
		} else {
			uValue = uBPValue + uAdd
		}
	}

	return uValue
}

// BP领奖请求
func (this *PlayerSystemChargeBP) BPPrizeRequest(head *pb.RpcHead, pbRequest *pb.BPPrizeRequest) {
	uCode := this.BPPrize(head, pbRequest.BPType, pbRequest.StageId)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.BPPrizeResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// BP领奖请求
func (this *PlayerSystemChargeBP) BPPrize(head *pb.RpcHead, uBPType uint32, uStageId uint32) cfgEnum.ErrorCode {
	pBPInfo, ok := this.mapBPInfo[uBPType]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, uBPType, uStageId)
	}

	pbStageInfo, ok := pBPInfo.mapStageInfo[uStageId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, uBPType, uStageId)
	}

	//奖励
	arrAddItem := this.getStagePrize(uBPType, pBPInfo.Value, pbStageInfo)
	if len(arrAddItem) <= 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HavePrize, uBPType, uStageId)
	}

	this.getPlayerBagFun().AddArrItem(head, arrAddItem, pb.EmDoingType_EDT_BP, true)
	this.UpdateSave(true)

	cluster.SendToClient(head, &pb.BPPrizeResponse{
		PacketHead:    &pb.IPacket{},
		BPType:        uBPType,
		StageId:       uStageId,
		NoramlPrizeId: pbStageInfo.NoramlPrizeId,
		ExtralPrizeId: pbStageInfo.ExtralPrizeId,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}
func (this *PlayerSystemChargeBP) canBuyProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode {
	if len(cfgCharge.Param) != 2 {
		return plog.Print(this.AccountId, cfgData.GetChargeErrorCode(cfgCharge.ProductID), cfgCharge.ProductID)
	}

	uBpType := cfgCharge.Param[0]
	uStageId := cfgCharge.Param[1]
	cfgBP := cfgData.GetCfgBpActConfig(uBpType, uStageId)
	if cfgBP == nil {
		return plog.Print(this.AccountId, cfgData.GetBpActConfigErrorCode(uBpType), cfgCharge.ProductID, uBpType, uStageId)
	}

	pPlayerPBInfo, ok := this.mapBPInfo[uBpType]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, cfgCharge.ProductID, uBpType, uStageId)
	}

	pbStageInfo, ok := pPlayerPBInfo.mapStageInfo[uStageId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, cfgCharge.ProductID, uBpType, uStageId)
	}

	//已经激活了
	if pbStageInfo.ChargeTime > 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HavePrize, cfgCharge.ProductID, uBpType, uStageId)
	}
	return cfgEnum.ErrorCode_Success
}
func (this *PlayerSystemChargeBP) OnChargeBuyProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode {
	if len(cfgCharge.Param) != 2 {
		return plog.Print(this.AccountId, cfgData.GetChargeErrorCode(cfgCharge.ProductID), cfgCharge.ProductID)
	}

	uBpType := cfgCharge.Param[0]
	uStageId := cfgCharge.Param[1]
	cfgBP := cfgData.GetCfgBpActConfig(uBpType, uStageId)
	if cfgBP == nil {
		return plog.Print(this.AccountId, cfgData.GetBpActConfigErrorCode(uBpType), cfgCharge.ProductID, uBpType, uStageId)
	}

	pPlayerPBInfo, ok := this.mapBPInfo[uBpType]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, cfgCharge.ProductID, uBpType, uStageId)
	}

	pbStageInfo, ok := pPlayerPBInfo.mapStageInfo[uStageId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, cfgCharge.ProductID, uBpType, uStageId)
	}

	//已经激活了
	if pbStageInfo.ChargeTime > 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HavePrize, cfgCharge.ProductID, uBpType, uStageId)
	}

	pbStageInfo.ChargeTime = base.GetNow()

	//加奖励
	this.getPlayerBagFun().AddArrItem(head, cfgBP.AddPrize, pb.EmDoingType_EDT_BP, true)

	this.UpdateSave(true)

	//通知客户端
	cluster.SendToClient(head, &pb.BPAcitiveNotify{
		PacketHead: &pb.IPacket{},
		BPType:     uBpType,
		StageId:    uStageId,
		ChargeTime: pbStageInfo.ChargeTime,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 成就类型
func (this *PlayerSystemChargeBP) TriggerAchieve(head *pb.RpcHead, emAchieveType cfgEnum.AchieveType, uAdd uint32, params ...uint32) {
	mapData := cfgData.GetCfgBpActConfigByAchieveType(uint32(emAchieveType))
	if mapData == nil {
		return
	}

	for _, cfgBp := range mapData {
		this.UpdateBPStage(head, cfgBp, uAdd)
	}
}

// 更新进度 自动开启下一期
func (this *PlayerSystemChargeBP) UpdateBPStage(head *pb.RpcHead, cfgBp *cfgData.BpActConfigCfg, uAdd uint32) {
	if cfgBp == nil {
		return
	}

	pPlayerPBInfo, ok := this.mapBPInfo[cfgBp.BpType]
	if !ok {
		return
	}

	//有改变通知客户端
	uOldValue := pPlayerPBInfo.Value
	pPlayerPBInfo.Value = this.GetPbValue(cfgBp, pPlayerPBInfo.Value, uAdd)
	if uOldValue != pPlayerPBInfo.Value {
		this.UpdateSave(true)
		//通知客户端
		cluster.SendToClient(head, &pb.BPValueNotify{
			PacketHead: &pb.IPacket{},
			BPType:     cfgBp.BpType,
			Value:      pPlayerPBInfo.Value,
		}, cfgEnum.ErrorCode_Success)
	}

	uCurTime := base.GetNow()

	//根据时间来开启
	pbMaxStage, ok := pPlayerPBInfo.mapStageInfo[uint32(len(pPlayerPBInfo.mapStageInfo))]
	if !ok {
		return
	}

	cfgMax := cfgData.GetCfgBpActConfig(cfgBp.BpType, pbMaxStage.StageId)
	if cfgMax == nil {
		return
	}

	cfgNextMax := cfgData.GetCfgBpActConfig(cfgBp.BpType, pbMaxStage.StageId+1)
	if cfgNextMax == nil {
		return
	}

	if cfgMax.PrePrizeIdOpenNext <= 0 {
		return
	}

	//时间未满足
	uDiffDay := base.DiffDays(pbMaxStage.BeginTime, uCurTime)
	if uDiffDay < cfgMax.OpenNextDay {
		return
	}

	if pPlayerPBInfo.Value < cfgData.GetCfgBpRewardPreIdNeedValue(cfgMax.BpType, cfgMax.Stage, cfgMax.PrePrizeIdOpenNext) {
		return
	}

	this.AddNewPbStage(head, cfgNextMax, true)
}

// 获取未领取的奖励
func (this *PlayerSystemChargeBP) getStagePrize(uBPType uint32, uValue uint32, pbStageInfo *pb.PBBPStageInfo) (arrAddItem []*common.ItemInfo) {
	arrAddItem = make([]*common.ItemInfo, 0)

	//奖励
	listPrize := cfgData.GetCfgBpRewardList(uBPType, pbStageInfo.StageId)
	if len(listPrize) <= 0 {
		return
	}

	for _, cfgPrize := range listPrize {
		if uValue < cfgPrize.NeedValue {
			break
		}

		//普通奖励
		if pbStageInfo.NoramlPrizeId < cfgPrize.PrizeId {
			pbStageInfo.NoramlPrizeId = cfgPrize.PrizeId
			arrAddItem = append(arrAddItem, cfgPrize.NormalPrize...)
		}

		//领取奖励
		if pbStageInfo.ChargeTime > 0 && pbStageInfo.ExtralPrizeId < cfgPrize.PrizeId {
			pbStageInfo.ExtralPrizeId = cfgPrize.PrizeId
			arrAddItem = append(arrAddItem, cfgPrize.ExtralPrize...)
		}
	}

	return
}

// 跨天 检查重置期数
func (this *PlayerSystemChargeBP) PassDay(isDay, isWeek, isMonth bool) {
	uCurTime := base.GetNow()
	for _, pBpInfo := range this.mapBPInfo {
		if len(pBpInfo.StageList) <= 0 {
			continue
		}

		pbNotify := &pb.BPNewStageNotify{
			PacketHead: &pb.IPacket{},
			BPType:     pBpInfo.BPType,
		}
		arrNewStage := make([]*cfgData.BpActConfigCfg, 0)
		arrDel := make([]*cfgData.BpActConfigCfg, 0)
		for _, pStageInfo := range pBpInfo.StageList {
			cfgBp := cfgData.GetCfgBpActConfig(pBpInfo.BPType, pStageInfo.StageId)
			if cfgBp == nil {
				continue
			}

			//判断一下是否开启  满足开启天数 和数值
			cfgNextBp := cfgData.GetCfgBpActConfig(pBpInfo.BPType, pStageInfo.StageId+1)
			if cfgNextBp != nil {
				passDay := base.DiffDays(pStageInfo.BeginTime, uCurTime)
				if _, ok := pBpInfo.mapStageInfo[pStageInfo.StageId+1]; !ok {
					if passDay >= cfgBp.OpenNextDay && pBpInfo.Value >= cfgData.GetCfgBpRewardPreIdNeedValue(pBpInfo.BPType, pStageInfo.StageId, cfgBp.PrePrizeIdOpenNext) {
						arrNewStage = append(arrNewStage, cfgNextBp)
					}
				}
			}

			if pStageInfo.EndTime > 0 && pStageInfo.EndTime <= uCurTime {
				//补发奖励邮件
				arrItem := this.getStagePrize(pBpInfo.BPType, pBpInfo.Value, pStageInfo)
				if len(arrItem) > 0 {
					this.getPlayerMailFun().AddTempMail(&pb.RpcHead{Id: this.AccountId}, cfgEnum.EMailId_BP, pb.EmDoingType_EDT_BP, arrItem, cfgBp.Title)
				}

				//结束 如果是循环活动，需要开启下一期
				if cfgBp.IsCircle > 0 {
					pStageInfo.BeginTime = uCurTime
					pStageInfo.EndTime = base.GetZeroTimestamp(uCurTime, 0) + uint64(cfgBp.AliveDay*3600*24) - 1
					pStageInfo.NoramlPrizeId = 0
					pStageInfo.ExtralPrizeId = 0
					pStageInfo.ChargeTime = 0
					//需要清理数据
					pBpInfo.Value = 0

					pbNotify.StageList = append(pbNotify.StageList, pStageInfo)

					//通知客户端
					cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.BPValueNotify{
						PacketHead: &pb.IPacket{},
						BPType:     pBpInfo.BPType,
						Value:      pBpInfo.Value,
					}, cfgEnum.ErrorCode_Success)
				} else {
					//判断结束 直接删除
					arrDel = append(arrDel, cfgBp)
				}
			}
		}

		//新增
		for _, cfgNew := range arrNewStage {
			pInfo := this.AddNewPbStage(&pb.RpcHead{Id: this.AccountId}, cfgNew, false)
			if pInfo != nil {
				pbNotify.StageList = append(pbNotify.StageList, pInfo)
			}
		}

		for _, cfgDel := range arrDel {
			this.DelBP(&pb.RpcHead{Id: this.AccountId}, cfgDel)
		}

		//通知客户端
		if len(pbNotify.StageList) > 0 || len(pbNotify.DelList) > 0 {
			cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbNotify, cfgEnum.ErrorCode_Success)
		}
	}
}
