package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/module/achieve"

	"github.com/golang/protobuf/proto"
)

type (
	PlayerSystemSevenDayFun struct {
		PlayerFun
		mapData    map[uint32]*pb.PBSevenDayInfo
		mapAchieve map[*common.AchieveKey]map[uint32][]*pb.PBTaskStageInfo //key：成就类型+参数 key2:活动ID key3:任务ID 列表
		mapTask    map[uint32]*pb.PBTaskStageInfo                          // 任务进度列表
		mapGift    map[uint32]*pb.PBU32U32                                 // 礼包列表

		*achieve.AchieveService
	}
)

func (this *PlayerSystemSevenDayFun) Init(pbType pb.PlayerDataType, pcommon *FunCommon) {
	this.PlayerFun.Init(pbType, pcommon)
	this.mapData = make(map[uint32]*pb.PBSevenDayInfo)
	this.mapAchieve = make(map[*common.AchieveKey]map[uint32][]*pb.PBTaskStageInfo)
	this.mapTask = make(map[uint32]*pb.PBTaskStageInfo)
	this.mapGift = make(map[uint32]*pb.PBU32U32)
	this.AchieveService = achieve.NewAchieveService(this.AccountId, cfgEnum.EAchieveSystemType_SevenDay, this.getPlayerSystemTaskFun().AchieveBase)
}

// 从数据库中加载
func (this *PlayerSystemSevenDayFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	this.loadData(pbSystem.SevenDay)
	this.UpdateSave(false)
}

func (this *PlayerSystemSevenDayFun) loadData(pbData *pb.PBPlayerSystemSevenDay) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemSevenDay{}
	}

	this.mapData = make(map[uint32]*pb.PBSevenDayInfo)
	for _, info := range pbData.SevenDayList {
		//任务进度
		for i := 0; i < len(info.TaskList); i++ {
			this.mapTask[info.TaskList[i].Id] = info.TaskList[i]

			if info.TaskList[i].Value >= info.TaskList[i].MaxValue && info.TaskList[i].State == pb.EmTaskState_ETS_Ing {
				info.TaskList[i].State = pb.EmTaskState_ETS_Finish
			}

			//成就数据
			if cfgTask := cfgData.GetCfgSevenDayTaskConfig(info.Id, info.TaskList[i].Id); cfgTask != nil {
				this.AddAchieve(cfgTask.AchieveType, cfgTask.AchieveParams, info.TaskList[i])
			}
		}

		//礼包进度
		for i := 0; i < len(info.GiftList); i++ {
			this.mapGift[info.GiftList[i].Key] = info.GiftList[i]
		}

		this.mapData[info.Id] = info
	}

	this.UpdateSave(true)
}

// 加载完成
func (this *PlayerSystemSevenDayFun) LoadComplete() {
}

// 存储到数据库
func (this *PlayerSystemSevenDayFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	if pbSystem.SevenDay == nil {
		pbSystem.SevenDay = new(pb.PBPlayerSystemSevenDay)
	}

	pbSystem.SevenDay.SevenDayList = make([]*pb.PBSevenDayInfo, 0)
	for _, info := range this.mapData {
		pbSystem.SevenDay.SevenDayList = append(pbSystem.SevenDay.SevenDayList, info)
	}

	return this.BSave
}
func (this *PlayerSystemSevenDayFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemSevenDay{}
}

// 设置玩家数据
func (this *PlayerSystemSevenDayFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}
	pbSystem := pbData.(*pb.PBPlayerSystemSevenDay)
	if pbSystem == nil {
		return false
	}
	this.loadData(pbSystem)
	this.UpdateSave(true)
	return true
}

func (this *PlayerSystemSevenDayFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = new(pb.PBPlayerSystem)
	}

	this.SaveSystem(pbData.System)
}

// 七天活跃奖励领取请求
func (this *PlayerSystemSevenDayFun) SevenDayActivePrizeRequest(head *pb.RpcHead, pbRequest *pb.SevenDayActivePrizeRequest) {
	uCode := this.SevenDayActivePrize(head, pbRequest.Id)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.SevenDayActivePrizeResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 跨天
func (this *PlayerSystemSevenDayFun) PassDay(isDay, isWeek, isMonth bool) {
	//判断活动
	pbListNotify := &pb.SevenDayListNotify{PacketHead: &pb.IPacket{Id: this.AccountId}}
	pbTaskNotify := &pb.AchieveTaskInfoNotify{
		PacketHead: &pb.IPacket{Id: this.AccountId},
		SystemType: uint32(cfgEnum.EAchieveSystemType_SevenDay),
	}
	mapCfg := cfgData.GetAllCfgSevenDayConfig()
	uRegisterDay := this.getPlayerBaseFun().GetRegDays()
	for _, cfg := range mapCfg {
		//结束 需要删除
		if uRegisterDay < cfg.BeginRegisterDays || uRegisterDay > cfg.EndRegisterDays {
			if sevenDayInfo, ok := this.mapData[cfg.Id]; ok {
				//需要删除成就列表 未领取进度值的给进度值
				for _, pbTmpTask := range this.mapData[cfg.Id].TaskList {
					cfgTask := cfgData.GetCfgSevenDayTaskConfig(cfg.Id, pbTmpTask.Id)
					if cfgTask != nil {
						if pbTmpTask.State == pb.EmTaskState_ETS_Finish {
							sevenDayInfo.Value += cfgTask.Value
						}

						this.DeleteAchieve(cfgTask.AchieveType, pbTmpTask)

					}
				}

				//进度奖励
				listCfg := cfgData.GetCfgSevenDayActivePrize(cfg.Id)
				arrPrize := make([]*common.ItemInfo, 0)
				for _, cfgActivie := range listCfg {
					if sevenDayInfo.PrizeValue >= cfgActivie.Value {
						continue
					}

					if sevenDayInfo.Value < cfgActivie.Value {
						break
					}

					arrPrize = append(arrPrize, cfgActivie.AddPrize...)
				}

				if len(arrPrize) > 0 {
					this.getPlayerMailFun().AddTempMail(&pb.RpcHead{Id: this.AccountId}, cfgEnum.EMailId_BP, pb.EmDoingType_EDT_BP, arrPrize, cfg.Name)
				}

				delete(this.mapData, cfg.Id)
				pbListNotify.Delist = append(pbListNotify.Delist, cfg.Id)
			}
			continue
		}

		//需要解锁新活动
		bNew := false
		if _, ok := this.mapData[cfg.Id]; !ok {
			uBeginTime := base.GetZeroTimestamp(base.GetNow(), int32(cfg.BeginRegisterDays)-int32(uRegisterDay))
			bNew = true
			this.mapData[cfg.Id] = &pb.PBSevenDayInfo{
				Id:        cfg.Id,
				BeginTime: uBeginTime,
				EndTime:   uBeginTime + uint64(cfg.EndRegisterDays*24*3600),
			}

			pbListNotify.AddList = append(pbListNotify.AddList, this.mapData[cfg.Id])
		}

		//触发新天数 解锁新任务
		uSevenDayPass := uRegisterDay - cfg.BeginRegisterDays + 1
		for i := uint32(1); i <= uSevenDayPass; i++ {
			listTask := cfgData.GetCfgSevenDayTaskByDays(cfg.Id, i)
			for _, cfgTask := range listTask {
				if _, ok := this.mapTask[cfgTask.Id]; ok {
					continue
				}

				//创建新任务
				pbTask := this.NewTaskInfo(cfgTask.Id, cfgTask.IsTotal, cfgTask.AchieveValue, cfgTask.AchieveType, cfgTask.AchieveParams...)

				this.mapData[cfg.Id].TaskList = append(this.mapData[cfg.Id].TaskList, pbTask)
				this.mapTask[cfgTask.Id] = pbTask
				//不是新的 需要同步变化
				if !bNew {
					pbTaskNotify.TaskList = append(pbTaskNotify.TaskList, pbTask)
				}
				pbTaskNotify.TaskList = append(pbTaskNotify.TaskList, pbTask)
			}

			//解锁新礼包
			listGift := cfgData.GetCfgSevenDayGiftByDays(cfg.Id, i)
			for _, cfgGift := range listGift {
				if _, ok := this.mapGift[cfgGift.Id]; ok {
					continue
				}

				//创建新任务
				pbGift := &pb.PBU32U32{
					Key:   cfgGift.Id,
					Value: 0,
				}

				this.mapData[cfg.Id].GiftList = append(this.mapData[cfg.Id].GiftList, pbGift)
				this.mapGift[cfgGift.Id] = pbGift
			}
		}

	}

	//同步列表变化
	if len(pbListNotify.AddList) > 0 || len(pbListNotify.Delist) > 0 {
		//需要清除登录数据
		this.getPlayerSystemTaskFun().DeleteAchieveBase(cfgEnum.AchieveType_LoginDay, uint32(cfgEnum.EAchieveSystemType_SevenDay))

		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbListNotify, cfgEnum.ErrorCode_Success)
	}

	//同步任务变化
	if len(pbTaskNotify.TaskList) > 0 {
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbTaskNotify, cfgEnum.ErrorCode_Success)
	}

	this.getPlayerSystemTaskFun().TriggerAchieve(&pb.RpcHead{Id: this.AccountId}, cfgEnum.AchieveType_LoginDay, 1, uint32(cfgEnum.EAchieveSystemType_SevenDay))
}

// 七天活跃奖励领取请求
func (this *PlayerSystemSevenDayFun) SevenDayActivePrize(head *pb.RpcHead, uId uint32) cfgEnum.ErrorCode {
	listCfg := cfgData.GetCfgSevenDayActivePrize(uId)
	if listCfg == nil {
		return plog.Print(this.AccountId, cfgData.GetSevenDayActivePrizeErrorCode(uId), uId)
	}

	pbInfo, ok := this.mapData[uId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, uId)
	}

	arrPrize := make([]*common.ItemInfo, 0)
	uNewValue := uint32(0)
	for _, cfg := range listCfg {
		if pbInfo.PrizeValue >= cfg.Value {
			continue
		}

		if pbInfo.Value < cfg.Value {
			break
		}

		arrPrize = append(arrPrize, cfg.AddPrize...)
		uNewValue = cfg.Value
	}

	if uNewValue <= pbInfo.PrizeValue || len(arrPrize) <= 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HavePrize, uId)
	}

	//存数据 加奖励
	pbInfo.PrizeValue = uNewValue
	this.getPlayerBagFun().AddArrItem(head, arrPrize, pb.EmDoingType_EDT_SevenDay, true)
	this.UpdateSave(true)

	//通知客户端
	cluster.SendToClient(head, &pb.SevenDayActivePrizeResponse{
		PacketHead: &pb.IPacket{},
		Id:         uId,
		PrizeValue: pbInfo.PrizeValue,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 七天任务奖励领取请求
func (this *PlayerSystemSevenDayFun) SevenDayTaskPrizeRequest(head *pb.RpcHead, pbRequest *pb.SevenDayTaskPrizeRequest) {
	uCode := this.SevenDayTaskPrize(head, pbRequest.Id, pbRequest.TaskId)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.SevenDayTaskPrizeResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 星源秘宝请求
func (this *PlayerSystemSevenDayFun) SevenDayTaskPrize(head *pb.RpcHead, uId uint32, uTaskId uint32) cfgEnum.ErrorCode {
	pbInfo, ok := this.mapData[uId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, uId)
	}

	pbTask, ok := this.mapTask[uTaskId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedCondition, uId)
	}

	if pbTask.State == pb.EmTaskState_ETS_Award {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HavePrize, uId)
	}

	if pbTask.State != pb.EmTaskState_ETS_Finish {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedCondition, uId)
	}

	cfgTask := cfgData.GetCfgSevenDayTaskConfig(uId, uTaskId)
	if !ok {
		return plog.Print(this.AccountId, cfgData.GetSevenDayTaskConfigErrorCode(uId), uId, uTaskId)
	}

	pbTask.State = pb.EmTaskState_ETS_Award
	pbInfo.Value += cfgTask.Value

	this.UpdateSave(true)
	//通知客户端
	cluster.SendToClient(head, &pb.SevenDayTaskPrizeResponse{
		PacketHead:    &pb.IPacket{},
		Id:            uId,
		Value:         pbInfo.Value,
		TaskStageInfo: pbTask,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 七天任务奖励领取请求
func (this *PlayerSystemSevenDayFun) SevenDayBuyGiftRequest(head *pb.RpcHead, pbRequest *pb.SevenDayBuyGiftRequest) {
	uCode := this.SevenDayBuyGift(head, pbRequest.Id, pbRequest.GiftId)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.SevenDayBuyGiftResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 七天任务奖励领取请求
func (this *PlayerSystemSevenDayFun) SevenDayBuyGift(head *pb.RpcHead, uId uint32, uGiftId uint32) cfgEnum.ErrorCode {
	pbInfo, ok := this.mapData[uId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, uId, uGiftId)
	}

	pbGift, ok := this.mapGift[uGiftId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, uId, uGiftId)
	}

	cfg := cfgData.GetCfgSevenDayGiftConfig(uGiftId)
	if cfg == nil {
		return plog.Print(this.AccountId, cfgData.GetSevenDayGiftConfigErrorCode(uGiftId), uId, uGiftId)
	}

	if uId != cfg.ActivityID {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_SeveDayActivityNotOnline, uId, uGiftId)
	}

	//充值礼包不能买
	if cfg.ProductId > 0 {
		return plog.Print(this.AccountId, cfgData.GetSevenDayGiftConfigErrorCode(uId), uId, uGiftId)
	}

	//最大次数限制
	if pbGift.Value >= cfg.LimitCount {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_BuyTimesLimit, uId, uGiftId)
	}

	//扣道具
	uCode := this.getPlayerBagFun().DelItem(head, cfg.DelItem.Kind, cfg.DelItem.Id, cfg.DelItem.Count, pb.EmDoingType_EDT_SevenDay)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode, uId, uGiftId)
	}

	pbGift.Value++
	pbInfo.Value += cfg.AddValue

	//加奖励
	this.getPlayerBagFun().AddArrItem(head, cfg.AddPrize, pb.EmDoingType_EDT_SevenDay, true)

	this.UpdateSave(true)

	//通知客户端
	cluster.SendToClient(head, &pb.SevenDayBuyGiftResponse{
		PacketHead: &pb.IPacket{},
		Id:         uId,
		Value:      pbInfo.Value,
		GiftInfo:   pbGift,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}
func (this *PlayerSystemSevenDayFun) canBuyProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode {
	if len(cfgCharge.Param) != 1 {
		return plog.Print(this.AccountId, cfgData.GetChargeErrorCode(cfgCharge.ProductID), cfgCharge.ProductID)
	}

	uGiftId := cfgCharge.Param[0]
	cfg := cfgData.GetCfgSevenDayGiftConfig(uGiftId)
	if cfg == nil {
		return plog.Print(this.AccountId, cfgData.GetSevenDayGiftConfigErrorCode(uGiftId), uGiftId)
	}

	//充值礼包不能买
	if cfg.ProductId != cfgCharge.ProductID {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_SeveDayGiftPID, uGiftId)
	}

	pbGift, ok := this.mapGift[uGiftId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, uGiftId)
	}

	if pbGift.Value >= cfg.LimitCount {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_BuyTimesLimit, uGiftId)
	}
	return cfgEnum.ErrorCode_Success
}

func (this *PlayerSystemSevenDayFun) OnChargeBuyProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode {
	if len(cfgCharge.Param) != 1 {
		return plog.Print(this.AccountId, cfgData.GetChargeErrorCode(cfgCharge.ProductID), cfgCharge.ProductID)
	}

	uGiftId := cfgCharge.Param[0]
	cfg := cfgData.GetCfgSevenDayGiftConfig(uGiftId)
	if cfg == nil {
		return plog.Print(this.AccountId, cfgData.GetSevenDayGiftConfigErrorCode(uGiftId), uGiftId)
	}

	//充值礼包不能买
	if cfg.ProductId != cfgCharge.ProductID {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_SeveDayGiftPID, uGiftId)
	}

	pbGift, ok := this.mapGift[uGiftId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, uGiftId)
	}

	if pbGift.Value >= cfg.LimitCount {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_BuyTimesLimit, uGiftId)
	}

	//加奖励
	this.getPlayerBagFun().AddArrItem(head, cfg.AddPrize, pb.EmDoingType_EDT_Charge, true)

	this.mapData[cfg.ActivityID].Value += cfg.AddValue
	pbGift.Value++

	//通知客户端
	cluster.SendToClient(head, &pb.SevenDayGiftNotify{
		PacketHead: &pb.IPacket{},
		GiftInfo:   pbGift,
		Value:      this.mapData[cfg.ActivityID].Value,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}
