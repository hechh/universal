package playerFun

import (
	"corps/base/cfgEnum"
	"corps/common/cfgData"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/module/achieve"
	"encoding/json"
)

type (
	PlayerActivityGrowRoad struct {
		*PlayerSystemActivityFun
		emType  cfgEnum.EActivityType
		mapData map[uint32]*pb.PBActivityGrowRoadInfo
		*achieve.AchieveService
		mapTask map[uint32]*pb.PBTaskStageInfo // 任务进度列表
	}
)

func init() {
	RegisterActivity(cfgEnum.EActivityType_GrowRoad, func() IPlayerSystemActivityFun { return new(PlayerActivityGrowRoad) })
}

func (this *PlayerActivityGrowRoad) Init(emType cfgEnum.EActivityType, pFun *PlayerSystemActivityFun) {
	this.emType = emType
	this.PlayerSystemActivityFun = pFun
	this.mapData = make(map[uint32]*pb.PBActivityGrowRoadInfo)
	this.AchieveService = achieve.NewAchieveService(this.AccountId, cfgEnum.EAchieveSystemType_GrowRoad, this.getPlayerSystemTaskFun().AchieveBase)
	this.mapTask = make(map[uint32]*pb.PBTaskStageInfo)

}

func (this *PlayerActivityGrowRoad) LoadPlayerDBFinish() {

}

func (this *PlayerActivityGrowRoad) LoadData(pbData *pb.PBPlayerSystemActivity) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemActivity{}
	}

	this.mapData = make(map[uint32]*pb.PBActivityGrowRoadInfo)
	for _, info := range pbData.GrowRoadList {
		this.mapData[info.Id] = info

		//任务进度
		for i := 0; i < len(info.TaskList); i++ {
			if info.TaskList[i].Value >= info.TaskList[i].MaxValue && info.TaskList[i].State == pb.EmTaskState_ETS_Ing {
				info.TaskList[i].State = pb.EmTaskState_ETS_Finish
			}

			//成就数据 需要注册
			if info.TaskList[i].State != pb.EmTaskState_ETS_Finish {
				if cfgTask := cfgData.GetCfgGrowRoadTaskConfig(info.Id, info.TaskList[i].Id); cfgTask != nil {
					this.AddAchieve(cfgTask.AchieveType, cfgTask.AchieveParams, info.TaskList[i])
				}
			}

			this.mapTask[info.TaskList[i].Id] = info.TaskList[i]
		}
	}

	this.UpdateSave(true)
}

func (this *PlayerActivityGrowRoad) LoadComplete() {
}

// 存储到数据库
func (this *PlayerActivityGrowRoad) SaveData(pbData *pb.PBPlayerSystemActivity) {
	for _, info := range this.mapData {
		pbData.GrowRoadList = append(pbData.GrowRoadList, info)
	}
}
func (this *PlayerActivityGrowRoad) Del(uID uint32) {
	if _, ok := this.mapData[uID]; !ok {
		plog.Error("(this *PlayerActivityGrowRoad) Del %d", uID)
		return
	}
	delete(this.mapData, uID)
	this.UpdateSave(true)
}
func (this *PlayerActivityGrowRoad) GetRed(uID uint32) bool {
	for _, info := range this.mapData {
		for _, taskInfo := range info.TaskList {
			if taskInfo.State == pb.EmTaskState_ETS_Finish {
				return true
			}
		}
	}

	return false
}
func (this *PlayerActivityGrowRoad) FreePrize(head *pb.RpcHead, uID uint32) cfgEnum.ErrorCode {
	return cfgEnum.ErrorCode_NoData
}
func (this *PlayerActivityGrowRoad) Open(uID uint32) []string {
	mapCfg := cfgData.GetAllCfgGrowRoadTaskConfig(uID)
	if len(mapCfg) <= 0 {
		return nil
	}
	byDat, err := json.Marshal(mapCfg)
	if err != nil {
		return nil
	}

	return []string{string(byDat)}
}

func (this *PlayerActivityGrowRoad) Add(uID uint32, uBeginTime uint64, uEndTime uint64) {
	if _, ok := this.mapData[uID]; ok {
		plog.Error("(this *PlayerActivityGrowRoad) Add repeated %d %d %d", uID, uBeginTime, uEndTime)
		return
	}

	this.mapData[uID] = &pb.PBActivityGrowRoadInfo{
		Id:        uID,
		BeginTime: uBeginTime,
		EndTime:   uEndTime,
	}

	//注册任务
	mapCfg := cfgData.GetAllCfgGrowRoadTaskConfig(uID)
	for _, cfgTask := range mapCfg {
		//创建新任务
		pbTask := this.NewTaskInfo(cfgTask.Id, 0, cfgTask.NeedValue, cfgTask.AchieveType, cfgTask.AchieveParams...)
		this.mapData[uID].TaskList = append(this.mapData[uID].TaskList, pbTask)
		this.mapTask[cfgTask.Id] = pbTask
	}

	this.UpdateSave(true)

	//通知客户端
	pbNotify := &pb.ActivityDataNewNotify{
		PacketHead:   &pb.IPacket{},
		ActivityType: uint32(this.emType),
		Info:         &pb.PBPlayerSystemActivity{},
	}
	pbNotify.Info.GrowRoadList = append(pbNotify.Info.GrowRoadList, this.mapData[uID])

	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbNotify, cfgEnum.ErrorCode_Success)
}

// 成长之路领奖请求
func (this *PlayerActivityGrowRoad) GrowRoadTaskPrizeRequest(head *pb.RpcHead, pbRequest *pb.GrowRoadTaskPrizeRequest) {
	uCode := this.GrowRoadTaskPrize(head, pbRequest.Id, pbRequest.TaskId)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.GrowRoadTaskPrizeResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 成长之路领奖请求
func (this *PlayerActivityGrowRoad) GrowRoadTaskPrize(head *pb.RpcHead, uId uint32, uTaskId uint32) cfgEnum.ErrorCode {
	_, ok := this.mapData[uId]
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

	cfgTask := cfgData.GetCfgGrowRoadTaskConfig(uId, uTaskId)
	if cfgTask == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_GrowRoadTaskIDNotFound, uId)
	}

	pbTask.State = pb.EmTaskState_ETS_Award

	this.getPlayerBagFun().AddArrItem(head, cfgTask.AddPrize, pb.EmDoingType_EDT_GrowRoad, true)
	this.UpdateSave(true)
	//通知客户端
	cluster.SendToClient(head, &pb.GrowRoadTaskPrizeResponse{
		PacketHead:    &pb.IPacket{},
		Id:            uId,
		TaskStageInfo: pbTask,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}
