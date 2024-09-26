package achieve

//成就公共数据
import (
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/common/serverCommon"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
)

type IAchieve interface {
	AddAchieve(uAchieveType uint32, params []uint32, pbTask *pb.PBTaskStageInfo)
	NewTaskInfo(id uint32, isTotal uint32, maxValue uint32, achieveType uint32, params ...uint32) *pb.PBTaskStageInfo
	Trigger(pAchieveKey *common.AchieveKey, uAdd uint32)
}

type AchieveService struct {
	IAchieve
	accountId    uint64
	emType       cfgEnum.EAchieveSystemType
	pAchieveBase *AchieveBase
	mapAchieve   map[uint32]map[common.AchieveKey]map[uint32]*pb.PBTaskStageInfo //key：成就类型+参数 key2:任务列表
}

func NewAchieveService(AccountId uint64, emType cfgEnum.EAchieveSystemType, pBase *AchieveBase) *AchieveService {
	pService := &AchieveService{
		accountId:    AccountId,
		emType:       emType,
		pAchieveBase: pBase,
		mapAchieve:   make(map[uint32]map[common.AchieveKey]map[uint32]*pb.PBTaskStageInfo),
	}
	pBase.RegisterAchieve(emType, pService)
	return pService
}

// 新解锁任务
func (this *AchieveService) NewTaskInfo(id uint32, isTotal uint32, maxValue uint32, achieveType uint32, params ...uint32) *pb.PBTaskStageInfo {
	pbTask := &pb.PBTaskStageInfo{
		Id:       id,
		Value:    0,
		MaxValue: maxValue,
		State:    pb.EmTaskState_ETS_Ing,
	}

	if isTotal > 0 {
		pbTask.Value = this.pAchieveBase.GetAchieveValue(achieveType, params...)
	}

	if pbTask.Value >= maxValue {
		pbTask.State = pb.EmTaskState_ETS_Finish
	}

	//增加成就触发
	this.AddAchieve(achieveType, params, pbTask)

	return pbTask
}

func (this *AchieveService) Clear() {
	this.mapAchieve = make(map[uint32]map[common.AchieveKey]map[uint32]*pb.PBTaskStageInfo)
}

// 触发成就
func (this *AchieveService) DeleteAchieve(uAchieveType uint32, pbTask *pb.PBTaskStageInfo) {
	if _, ok := this.mapAchieve[uAchieveType]; !ok {
		return
	}

	//删除
	for _, mapInfo := range this.mapAchieve[uAchieveType] {
		if _, ok := mapInfo[pbTask.Id]; ok {
			delete(mapInfo, pbTask.Id)
		}
	}

}

// 增加成就触发
func (this *AchieveService) AddAchieve(uAchieveType uint32, params []uint32, pbTask *pb.PBTaskStageInfo) {
	if _, ok := this.mapAchieve[uAchieveType]; !ok {
		this.mapAchieve[uAchieveType] = make(map[common.AchieveKey]map[uint32]*pb.PBTaskStageInfo, 0)
	}

	//增加成就触发
	pAchieveKey := common.AchieveKey{
		AchieveType: uAchieveType,
	}
	if len(params) > 0 {
		pAchieveKey.Param1 = params[0]
	}
	if len(params) > 1 {
		pAchieveKey.Param2 = params[1]
	}

	//修正三个参数
	if len(params) > 2 {
		if uAchieveType == uint32(cfgEnum.AchieveType_BattleMap) {
			pAchieveKey.Param2 = serverCommon.MAKE_BATTLE_MAP(params[1], params[2])
		}
	}

	if _, ok := this.mapAchieve[uAchieveType][pAchieveKey]; !ok {
		this.mapAchieve[uAchieveType][pAchieveKey] = make(map[uint32]*pb.PBTaskStageInfo, 0)
	}

	this.mapAchieve[uAchieveType][pAchieveKey][pbTask.Id] = pbTask
}

// 增加成就触发
func (this *AchieveService) Trigger(pAchieveKey *common.AchieveKey, uAdd uint32) {
	//是否取累计
	cfgAchieveType := cfgData.GetCfgAchieveTypeConfig(pAchieveKey.AchieveType)
	if cfgAchieveType == nil {
		plog.Info("(this *PlayerSystemTaskFun) TriggerAchieve cfg error", pAchieveKey.AchieveType, uAdd)
		return
	}

	//判断是否有影响的
	mapTasks, ok := this.mapAchieve[pAchieveKey.AchieveType]
	if !ok {
		return
	}

	//判断比较条件
	pbNotify := &pb.AchieveTaskInfoNotify{
		PacketHead: &pb.IPacket{Id: this.accountId},
		SystemType: uint32(this.emType),
	}

	//强等于
	if cfgAchieveType.CompareType == uint32(cfgEnum.ECompareType_Equal) {
		listTask, ok := mapTasks[*pAchieveKey]
		if !ok {
			return
		}

		//更新数据
		for _, pbTask := range listTask {
			if pbTask.State != pb.EmTaskState_ETS_Ing {
				continue
			}

			if cfgAchieveType.IsSet > 0 {
				pbTask.Value = this.pAchieveBase.GetAchieveValue(pAchieveKey.AchieveType, pAchieveKey.Param1, pAchieveKey.Param2)
			} else {
				pbTask.Value += uAdd
			}

			if pbTask.Value >= pbTask.MaxValue {
				pbTask.Value = pbTask.MaxValue
				pbTask.State = pb.EmTaskState_ETS_Finish
			}
			pbNotify.TaskList = append(pbNotify.TaskList, pbTask)

			if cfgAchieveType.AchieveType == uint32(cfgEnum.AchieveType_IncreasePower) {
				plog.Info("MaxHistoryFightPower updatestage id:%d cur:%d add:%d", this.accountId, pbTask.Value, uAdd)
			}
		}
	} else if cfgAchieveType.CompareType == uint32(cfgEnum.ECompareType_NotSmall) {
		for pTmpKey, pbTasks := range mapTasks {
			for _, pbTask := range pbTasks {
				if pbTask.State != pb.EmTaskState_ETS_Ing {
					continue
				}

				pbTask.Value = this.pAchieveBase.GetAchieveValue(pAchieveKey.AchieveType, pTmpKey.Param1, pTmpKey.Param2)
				if pbTask.Value >= pbTask.MaxValue {
					pbTask.Value = pbTask.MaxValue
					pbTask.State = pb.EmTaskState_ETS_Finish
				}
				pbNotify.TaskList = append(pbNotify.TaskList, pbTask)
			}

		}
	}

	if len(pbNotify.TaskList) > 0 {
		//通知客户端 成就更新
		cluster.SendToClient(&pb.RpcHead{Id: this.accountId}, pbNotify, cfgEnum.ErrorCode_Success)
	}
}
