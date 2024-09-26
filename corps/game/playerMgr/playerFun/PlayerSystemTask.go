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
	PlayerSystemTaskFun struct {
		PlayerFun
		pbMainTask           *pb.PBTaskStageInfo     //主线任务
		pbDailyTask          *pb.PBDailyTask         //每日任务
		*achieve.AchieveBase                         //成就系统
		pMainTaskAchieve     *achieve.AchieveService //主线任务成就系统
		pDailyTaskAchieve    *achieve.AchieveService //每日任务成就系统
	}
)

func (this *PlayerSystemTaskFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.AchieveBase = achieve.NewAchieveBase()
	this.pMainTaskAchieve = achieve.NewAchieveService(this.AccountId, cfgEnum.EAchieveSystemType_MainTask, this.AchieveBase)
	this.pDailyTaskAchieve = achieve.NewAchieveService(this.AccountId, cfgEnum.EAchieveSystemType_DailyTask, this.AchieveBase)
}

// 新玩家 需要初始化数据
func (this *PlayerSystemTaskFun) NewPlayer() {

	//初始化等级
	this.UpdateSave(true)
}

// 从数据库中加载
func (this *PlayerSystemTaskFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	if pbSystem.Task == nil {
		return
	}

	this.loadData(pbSystem.Task)

	this.UpdateSave(false)
}
func (this *PlayerSystemTaskFun) loadData(pbData *pb.PBPlayerSystemTask) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemTask{}
	}

	this.pbMainTask = pbData.MainTask
	this.pbDailyTask = pbData.DailyTask

	for _, info := range pbData.AchieveList {
		this.LoadAchieveBase(info)
	}

	this.UpdateSave(true)
}

// 设置数据
func (this *PlayerSystemTaskFun) LoadPlayerDBFinish() {
	if this.pbMainTask == nil {
		cfgTask := cfgData.GetCfgMainTask(1)
		if cfgTask == nil {
			plog.Info("(this *PlayerSystemTaskFun) LoadPlayerDBFinish task is not find 1")
			return
		}

		this.pbMainTask = &pb.PBTaskStageInfo{
			Id:       cfgTask.Id,
			MaxValue: cfgTask.Value,
			Value:    0,
			State:    pb.EmTaskState_ETS_Ing,
		}
		if cfgTask.IsTotal > 0 {
			this.pbMainTask.Value = this.GetAchieveValue(cfgTask.Type, cfgTask.Param...)
			if this.pbMainTask.Value >= this.pbMainTask.MaxValue {
				this.pbMainTask.State = pb.EmTaskState_ETS_Finish
			}
		}

	}

	cfgTask := cfgData.GetCfgMainTask(this.pbMainTask.Id)
	this.pMainTaskAchieve.Clear()
	this.pMainTaskAchieve.AddAchieve(cfgTask.Type, cfgTask.Param, this.pbMainTask)

	if this.pbDailyTask == nil {
		this.pbDailyTask = &pb.PBDailyTask{}
	}

	//注册到每日任务成就中
	for _, info := range this.pbDailyTask.TaskList {
		cfgDailyTask := cfgData.GetCfgDailyTaskConfig(info.Id)
		if cfgDailyTask == nil {
			continue
		}

		this.pDailyTaskAchieve.AddAchieve(cfgDailyTask.AchieveType, cfgDailyTask.AchieveParams, info)
	}
}

// 存储到数据库
func (this *PlayerSystemTaskFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	if pbSystem.Task == nil {
		pbSystem.Task = new(pb.PBPlayerSystemTask)
	}

	pbSystem.Task.MainTask = this.pbMainTask
	pbSystem.Task.DailyTask = this.pbDailyTask
	if pbSystem.Task.AchieveList == nil {
		pbSystem.Task.AchieveList = make([]*pb.PBAchieveInfo, 0)
	}

	this.SaveAchieveBase(&pbSystem.Task.AchieveList)

	return this.BSave
}
func (this *PlayerSystemTaskFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemTask{}
}

// 设置玩家数据
func (this *PlayerSystemTaskFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}

	pbSystem := pbData.(*pb.PBPlayerSystemTask)
	if pbSystem == nil {
		return false
	}

	this.loadData(pbSystem)

	this.LoadPlayerDBFinish()
	this.UpdateSave(true)

	//通知
	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.AchieveTaskInfoNotify{
		PacketHead: &pb.IPacket{},
		SystemType: uint32(cfgEnum.EAchieveSystemType_MainTask),
		TaskList:   []*pb.PBTaskStageInfo{this.pbMainTask},
	}, cfgEnum.ErrorCode_Success)

	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.DailyTaskNotify{
		PacketHead: &pb.IPacket{},
		DailyTask:  this.pbDailyTask,
	}, cfgEnum.ErrorCode_Success)
	return true
}

// 任务完成请求
func (this *PlayerSystemTaskFun) MainTaskFinishRequest(head *pb.RpcHead, pbRequest proto.Message) {
	uCode := this.MainTaskFinish(head)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.MainTaskFinishResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

func (this *PlayerSystemTaskFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = new(pb.PBPlayerSystem)
	}

	this.SaveSystem(pbData.System)
}

// 跨天
func (this *PlayerSystemTaskFun) PassDay(isDay, isWeek, isMonth bool) {
	//生成每日任务
	this.pDailyTaskAchieve.Clear()
	this.pbDailyTask.TaskList = make([]*pb.PBTaskStageInfo, 0)
	this.pbDailyTask.PrizeScore = 0
	this.pbDailyTask.Score = 0

	mapCfg := cfgData.GetAllCfgDailyTaskConfig()
	for _, cfg := range mapCfg {
		if cfg.SystemOpenId > 0 && !this.getPlayerSystemCommonFun().IsSystemOpen(cfg.SystemOpenId) {
			continue
		}

		pbTask := &pb.PBTaskStageInfo{
			Id:       cfg.Id,
			MaxValue: cfg.AchieveValue,
			Value:    0,
			State:    pb.EmTaskState_ETS_Ing,
		}
		this.pbDailyTask.TaskList = append(this.pbDailyTask.TaskList, pbTask)
		this.pDailyTaskAchieve.AddAchieve(cfg.AchieveType, cfg.AchieveParams, pbTask)
	}

	this.UpdateSave(true)
	//更新每日任务
	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.DailyTaskNotify{
		PacketHead: &pb.IPacket{},
		DailyTask:  this.pbDailyTask,
	}, cfgEnum.ErrorCode_Success)
}

// 任务完成请求
func (this *PlayerSystemTaskFun) MainTaskFinish(head *pb.RpcHead) cfgEnum.ErrorCode {
	//判断是否完成
	if this.pbMainTask.State == pb.EmTaskState_ETS_Ing {
		return plog.Print(head.Id, cfgEnum.ErrorCode_TaskNoFinish, this.pbMainTask)
	}

	cfgTask := cfgData.GetCfgMainTask(this.pbMainTask.Id)
	if cfgTask == nil {
		return plog.Print(head.Id, cfgData.GetMainTaskErrorCode(this.pbMainTask.Id), this.pbMainTask)
	}

	//给奖励
	if this.pbMainTask.State == pb.EmTaskState_ETS_Finish {
		this.pbMainTask.State = pb.EmTaskState_ETS_Award
		this.getPlayerBagFun().AddArrItem(head, cfgTask.ListAddItem, pb.EmDoingType_EDT_Task, true)
	}

	if cfgTask.NextId > 0 {
		cfgNext := cfgData.GetCfgMainTask(cfgTask.NextId)
		if cfgNext == nil {
			return plog.Print(head.Id, cfgData.GetMainTaskErrorCode(cfgTask.NextId), this.pbMainTask)
		}

		this.pMainTaskAchieve.Clear()
		this.pbMainTask = &pb.PBTaskStageInfo{
			Id:       cfgNext.Id,
			MaxValue: cfgNext.Value,
			State:    pb.EmTaskState_ETS_Ing,
			Value:    0,
		}

		if cfgNext.IsTotal > 0 {
			this.pbMainTask.Value = this.GetAchieveValue(cfgNext.Type, cfgNext.Param...)
		}

		if this.pbMainTask.Value >= this.pbMainTask.MaxValue {
			this.pbMainTask.State = pb.EmTaskState_ETS_Finish
		}

		this.pMainTaskAchieve.AddAchieve(cfgNext.Type, cfgNext.Param, this.pbMainTask)
	}

	this.UpdateSave(true)

	//通知客户端
	cluster.SendToClient(head, &pb.MainTaskFinishResponse{
		PacketHead: &pb.IPacket{},
		MainTask:   this.pbMainTask,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

func (this *PlayerSystemTaskFun) FinishAllTask(head *pb.RpcHead) {
	uMaxTask := cfgData.GetAllCfgMainTask()
	this.pbMainTask.Id = uMaxTask
	this.pbMainTask.State = pb.EmTaskState_ETS_Award
	this.pbMainTask.Value = 0

	this.UpdateSave(true)
}

// 触发成就
func (this *PlayerSystemTaskFun) TriggerAchieve(head *pb.RpcHead, emAchieveType cfgEnum.AchieveType, uAdd uint32, params ...uint32) {
	//存储成就
	if this.AddAchieveBase(emAchieveType, uAdd, params...) {
		this.UpdateSave(true)
	}

	//词条触发
	if this.getEntry() != nil {
		this.getEntry().Trigger(head, uint32(emAchieveType), uAdd, params...)
	}

	//BP
	this.getPlayerSystemChargeBPFun().TriggerAchieve(head, emAchieveType, uAdd, params...)

}

// 每日任务完成请求
func (this *PlayerSystemTaskFun) DailyTaskFinishRequest(head *pb.RpcHead, pbRequest *pb.DailyTaskFinishRequest) {
	uCode := this.DailyTaskFinish(head, pbRequest.TaskId)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.MainTaskFinishResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 每日任务完成请求
func (this *PlayerSystemTaskFun) DailyTaskFinish(head *pb.RpcHead, uTaskId uint32) cfgEnum.ErrorCode {
	cfgTask := cfgData.GetCfgDailyTaskConfig(uTaskId)
	if cfgTask == nil {
		return plog.Print(head.Id, cfgData.GetDailyTaskConfigErrorCode(uTaskId), uTaskId)
	}

	index := 0
	for index = 0; index < len(this.pbDailyTask.TaskList); index++ {
		if this.pbDailyTask.TaskList[index].Id == uTaskId {
			break
		}
	}

	pbTask := this.pbDailyTask.TaskList[index]
	if pbTask.State == pb.EmTaskState_ETS_Ing {
		if pbTask.Value < pbTask.MaxValue {
			return plog.Print(head.Id, cfgEnum.ErrorCode_TaskNoFinish, uTaskId)
		}

	} else if pbTask.State == pb.EmTaskState_ETS_Award {
		return plog.Print(head.Id, cfgEnum.ErrorCode_HavePrize, uTaskId)
	}

	pbTask.State = pb.EmTaskState_ETS_Award

	//给奖励
	this.pbDailyTask.Score += cfgTask.AddScore

	this.UpdateSave(true)
	cluster.SendToClient(head, &pb.DailyTasFinishResponse{
		PacketHead: &pb.IPacket{},
		Task:       pbTask,
		Score:      this.pbDailyTask.Score,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 星源秘宝积分奖励请求
func (this *PlayerSystemTaskFun) DailyTaskScorePrizeRequest(head *pb.RpcHead) {
	uCode := this.DailyTaskScorePrize(head)
	cluster.SendToClient(head, &pb.DailyTaskScorePrizeResponse{
		PacketHead: &pb.IPacket{},
		PrizeScore: this.pbDailyTask.PrizeScore,
	}, uCode)
}

// 星源秘宝积分奖励请求
func (this *PlayerSystemTaskFun) DailyTaskScorePrize(head *pb.RpcHead) cfgEnum.ErrorCode {
	listCfg := cfgData.GetCfgAllDailyTaskScoreConfig()

	arrPrize := make([]*common.ItemInfo, 0)
	uNewValue := uint32(0)
	for _, cfg := range listCfg {
		if this.pbDailyTask.PrizeScore >= cfg.Value {
			continue
		}

		if this.pbDailyTask.Score < cfg.Value {
			break
		}

		arrPrize = append(arrPrize, cfg.AddPrize...)
		uNewValue = cfg.Value
	}

	if uNewValue <= this.pbDailyTask.PrizeScore || len(arrPrize) <= 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HavePrize)
	}

	//存数据 加奖励
	this.pbDailyTask.PrizeScore = uNewValue
	this.getPlayerBagFun().AddArrItem(head, arrPrize, pb.EmDoingType_EDT_DailyTask, true)
	this.UpdateSave(true)

	return cfgEnum.ErrorCode_Success
}

// 星源秘宝积分奖励请求
func (this *PlayerSystemTaskFun) repairMainTask0828() {
	arrRepairTask := []uint32{24, 25, 27, 28, 29, 39, 40}
	if this.pbMainTask.State == pb.EmTaskState_ETS_Finish {
		return
	}

	if !base.ArrayContainsValue(arrRepairTask, this.pbMainTask.Id) {
		return
	}

	cfgMainTask := cfgData.GetCfgMainTask(this.pbMainTask.Id)
	if cfgMainTask == nil {
		return
	}
	if cfgMainTask.IsTotal > 0 {
		this.pbMainTask.Value = this.GetAchieveValue(cfgMainTask.Type, cfgMainTask.Param...)
	} else {
		this.pbMainTask.Value = 0
	}

	this.pbMainTask.MaxValue = cfgMainTask.Value

	if this.pbMainTask.Value >= this.pbMainTask.MaxValue {
		this.pbMainTask.State = pb.EmTaskState_ETS_Finish
	}

	this.pMainTaskAchieve.Clear()
	this.pMainTaskAchieve.AddAchieve(cfgMainTask.Type, cfgMainTask.Param, this.pbMainTask)

	//通知客户端 成就更新
	//判断比较条件
	pbNotify := &pb.AchieveTaskInfoNotify{
		PacketHead: &pb.IPacket{Id: this.AccountId},
		SystemType: uint32(cfgEnum.EAchieveSystemType_MainTask),
	}
	pbNotify.TaskList = append(pbNotify.TaskList, this.pbMainTask)
	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbNotify, cfgEnum.ErrorCode_Success)
}
