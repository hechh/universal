package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common/cfgData"
	"corps/common/dao/player_base_display"
	"corps/common/dao/rank_info"
	"corps/framework/cluster"
	"corps/framework/common/uerror"
	"corps/pb"
	"corps/server/game/module/achieve"

	"github.com/golang/protobuf/proto"
)

const (
	MAX_DWORD_VALUE = 0xffffffff
)

var (
	rankTypes = []uint32{
		uint32(cfgEnum.ERankType_ChampionshipBattleHook),
		uint32(cfgEnum.ERankType_ChampionshipBattle),
		uint32(cfgEnum.ERankType_ChampionshipDamage),
		uint32(cfgEnum.ERankType_ChampionshipPower),
	}
)

// ----gomaker生成的模板-------
type PlayerSystemChampionshipFun struct {
	PlayerFun
	pbData *pb.PBPlayerSystemChampionship
	battle *achieve.AchieveService
	power  *achieve.AchieveService
	damage *achieve.AchieveService
	hook   *achieve.AchieveService
	flag   uint32
}

func (this *PlayerSystemChampionshipFun) getAchieveService(rankType uint32) *achieve.AchieveService {
	switch rankType {
	case uint32(cfgEnum.ERankType_ChampionshipBattleHook):
		return this.hook
	case uint32(cfgEnum.ERankType_ChampionshipBattle):
		return this.battle
	case uint32(cfgEnum.ERankType_ChampionshipDamage):
		return this.damage
	}
	return this.power
}

func (this *PlayerSystemChampionshipFun) registerTask(tt *pb.PBTaskStageInfo, rankType uint32) (flag bool) {
	if tt == nil {
		flag = true
		tt = &pb.PBTaskStageInfo{Id: 0, Value: 0, MaxValue: 0, State: pb.EmTaskState_ETS_Award}
	}
	// 判断任务是否已经完成
	isNext, taskID := false, tt.Id
	if tt.State == pb.EmTaskState_ETS_Award {
		isNext = true
		taskID = tt.Id
	}
	if isNext {
		if cfg := cfgData.GetCfgNextChampionshipTaskConfig(rankType, taskID); cfg != nil {
			tt = &pb.PBTaskStageInfo{Id: cfg.Id, MaxValue: cfg.AchieveValue, State: pb.EmTaskState_ETS_Ing}
			if cfg.IsTotal > 0 {
				tt.Value = this.getPlayerSystemTaskFun().GetAchieveValue(cfg.AchieveType, cfg.AchieveParams...)
				if tt.MaxValue <= tt.Value && tt.State != pb.EmTaskState_ETS_Award {
					tt.State = pb.EmTaskState_ETS_Finish
				}
			}
			flag = true
			this.getAchieveService(rankType).AddAchieve(cfg.AchieveType, cfg.AchieveParams, tt)
		}
	} else {
		if cfg := cfgData.GetCfgChampionshipTaskConfig(rankType, taskID); cfg != nil {
			if tt.MaxValue <= tt.Value && tt.State != pb.EmTaskState_ETS_Award {
				tt.State = pb.EmTaskState_ETS_Finish
			}
			this.getAchieveService(rankType).AddAchieve(cfg.AchieveType, cfg.AchieveParams, tt)
		} else {
			// 任务已经丢弃，直接进入下一个任务
			if cfg := cfgData.GetCfgNextChampionshipTaskConfig(rankType, taskID); cfg != nil {
				tt = &pb.PBTaskStageInfo{Id: cfg.Id, MaxValue: cfg.AchieveValue, State: pb.EmTaskState_ETS_Ing}
				flag = true
				this.getAchieveService(rankType).AddAchieve(cfg.AchieveType, cfg.AchieveParams, tt)
			}
		}
	}
	switch rankType {
	case uint32(cfgEnum.ERankType_ChampionshipBattleHook):
		this.pbData.Hook = tt
	case uint32(cfgEnum.ERankType_ChampionshipBattle):
		this.pbData.Battle = tt
	case uint32(cfgEnum.ERankType_ChampionshipDamage):
		this.pbData.Damage = tt
	case uint32(cfgEnum.ERankType_ChampionshipPower):
		this.pbData.Power = tt
	}
	return
}

// --------------------通用接口实现------------------------------
// 初始化
func (this *PlayerSystemChampionshipFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.hook = achieve.NewAchieveService(this.AccountId, cfgEnum.EAchieveSystemType_ChampionshipBattleHook, this.getPlayerSystemTaskFun().AchieveBase)
	this.battle = achieve.NewAchieveService(this.AccountId, cfgEnum.EAchieveSystemType_ChampionshipBattle, this.getPlayerSystemTaskFun().AchieveBase)
	this.damage = achieve.NewAchieveService(this.AccountId, cfgEnum.EAchieveSystemType_ChampionshipDamage, this.getPlayerSystemTaskFun().AchieveBase)
	this.power = achieve.NewAchieveService(this.AccountId, cfgEnum.EAchieveSystemType_ChampionshipPower, this.getPlayerSystemTaskFun().AchieveBase)
}

// 新系统
func (this *PlayerSystemChampionshipFun) NewPlayer() {
	this.pbData = &pb.PBPlayerSystemChampionship{}
	flag01 := this.registerTask(this.pbData.Hook, uint32(cfgEnum.ERankType_ChampionshipBattleHook))
	flag02 := this.registerTask(this.pbData.Battle, uint32(cfgEnum.ERankType_ChampionshipBattle))
	flag03 := this.registerTask(this.pbData.Damage, uint32(cfgEnum.ERankType_ChampionshipDamage))
	flag04 := this.registerTask(this.pbData.Power, uint32(cfgEnum.ERankType_ChampionshipPower))
	this.UpdateSave(flag01 || flag02 || flag03 || flag04)
}

// 加载系统数据(system类型数据)
func (this *PlayerSystemChampionshipFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	if pbSystem.Championship == nil {
		this.NewPlayer()
		return
	}
	this.pbData = pbSystem.Championship

	flag01 := this.registerTask(this.pbData.Hook, uint32(cfgEnum.ERankType_ChampionshipBattleHook))
	flag02 := this.registerTask(this.pbData.Battle, uint32(cfgEnum.ERankType_ChampionshipBattle))
	flag03 := this.registerTask(this.pbData.Damage, uint32(cfgEnum.ERankType_ChampionshipDamage))
	flag04 := this.registerTask(this.pbData.Power, uint32(cfgEnum.ERankType_ChampionshipPower))
	this.UpdateSave(flag01 || flag02 || flag03 || flag04)
}

// 存储数据 返回存储标志(system类型数据)
func (this *PlayerSystemChampionshipFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	pbSystem.Championship = this.pbData
	return true
}

// 客户端数据
func (this *PlayerSystemChampionshipFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = &pb.PBPlayerSystem{}
	}
	pbData.System.Championship = this.pbData
}

func (this *PlayerSystemChampionshipFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemChampionship{}
}

// 设置玩家数据, web管理后台
func (this *PlayerSystemChampionshipFun) SetUserTypeInfo(message proto.Message) bool {
	if message == nil {
		return false
	}
	pbSystem, ok := message.(*pb.PBPlayerSystemChampionship)
	if !ok || pbSystem == nil {
		return false
	}
	this.flag = 0
	this.pbData = pbSystem
	this.registerTask(this.pbData.Hook, uint32(cfgEnum.ERankType_ChampionshipBattleHook))
	this.registerTask(this.pbData.Battle, uint32(cfgEnum.ERankType_ChampionshipBattle))
	this.registerTask(this.pbData.Damage, uint32(cfgEnum.ERankType_ChampionshipDamage))
	this.registerTask(this.pbData.Power, uint32(cfgEnum.ERankType_ChampionshipPower))
	this.UpdateSave(true)
	return true
}

func (this *PlayerSystemChampionshipFun) InitRank() {
	now := base.GetNow()
	createTime := this.getPlayerBaseFun().GetServerStartTime()
	for _, rankType := range rankTypes {
		// 加载排行榜配置，获取开启时间
		cfg := cfgData.GetCfgRankInfoConfig(rankType)
		if cfg == nil {
			continue
		}
		// 判断是否解锁
		startTime := cfgData.GetCfgRankActiveTime(cfg, createTime)
		if now < startTime {
			continue
		}
		// 设置排行榜初始值
		switch rankType {
		case uint32(cfgEnum.ERankType_ChampionshipBattleHook):
			if (this.flag&1) <= 0 && this.pbData.HookHasReward <= 0 {
				this.UpdateRank(startTime, rankType, this.getPlayerSystemBattleHookFun().GetRankValue())
				this.flag |= 1
			}
		case uint32(cfgEnum.ERankType_ChampionshipBattle):
			if (this.flag&(1<<1)) <= 0 && this.pbData.BattleHasReward <= 0 {
				this.UpdateRank(startTime, rankType, this.getPlayerSystemBattleNormalFun().GetRankValue())
				this.flag |= (1 << 1)
			}
		case uint32(cfgEnum.ERankType_ChampionshipPower):
			if (this.flag&(1<<2)) <= 0 && this.pbData.PowerHasReward <= 0 {
				this.UpdateRank(startTime, rankType, this.getPlayerHeroFun().GetMaxHistoryFightPower())
				this.flag |= (1 << 2)
			}
		case uint32(cfgEnum.ERankType_ChampionshipDamage):
			if (this.flag&(1<<3)) <= 0 && this.pbData.DamageHasReward <= 0 {
				this.UpdateRank(startTime, rankType, this.getPlayerSystemWorldBossFun().GetMaxDamage())
				this.flag |= (1 << 3)
			}
		}
	}
}

func (this *PlayerSystemChampionshipFun) LoadComplete() {
	createTime := this.getPlayerBaseFun().GetServerStartTime()
	notify := &pb.ChampionshipNotify{
		PacketHead: &pb.IPacket{},
		CreateTime: createTime,
	}
	for _, rankType := range rankTypes {
		// 读取配置
		cfg := cfgData.GetCfgRankInfoConfig(rankType)
		if cfg == nil {
			continue
		}

		active := cfgData.GetCfgRankActiveTime(cfg, createTime)
		reward := cfgData.GetCfgRankRewardTime(cfg, createTime)
		show := cfgData.GetCfgRankShowTime(cfg, createTime)
		close := cfgData.GetCfgRankCloseTime(cfg, createTime)

		// 组装数据
		notify.List = append(notify.List, &pb.ChampionshipTimeInfo{
			RankType: rankType,
			Interval: active - createTime,
			Active:   reward - active,
			Reward:   show - reward,
			Show:     close - show,
		})
		if expire := close - createTime; expire > notify.Expire {
			notify.Expire = expire
		}
	}
	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, notify, cfgEnum.ErrorCode_Success)
}

func (this *PlayerSystemChampionshipFun) SetChampionshipFlag(rankType uint32, val uint32) {
	switch rankType {
	case uint32(cfgEnum.ERankType_ChampionshipDamage):
		this.pbData.DamageHasReward = val
		this.UpdateSave(true)
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.RankRewardResponse{PacketHead: &pb.IPacket{}}, cfgEnum.ErrorCode_Success)
	case uint32(cfgEnum.ERankType_ChampionshipBattleHook):
		this.pbData.HookHasReward = val
		this.UpdateSave(true)
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.RankRewardResponse{PacketHead: &pb.IPacket{}}, cfgEnum.ErrorCode_Success)
	case uint32(cfgEnum.ERankType_ChampionshipPower):
		this.pbData.PowerHasReward = val
		this.UpdateSave(true)
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.RankRewardResponse{PacketHead: &pb.IPacket{}}, cfgEnum.ErrorCode_Success)
	case uint32(cfgEnum.ERankType_ChampionshipBattle):
		this.pbData.BattleHasReward = val
		this.UpdateSave(true)
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.RankRewardResponse{PacketHead: &pb.IPacket{}}, cfgEnum.ErrorCode_Success)
	}
}

// --------------------交互接口实现------------------------------

func (this *PlayerSystemChampionshipFun) getTaskStageInfo(rankType uint32) *pb.PBTaskStageInfo {
	switch rankType {
	case uint32(cfgEnum.ERankType_ChampionshipDamage):
		return this.pbData.Damage
	case uint32(cfgEnum.ERankType_ChampionshipBattleHook):
		return this.pbData.Hook
	case uint32(cfgEnum.ERankType_ChampionshipPower):
		return this.pbData.Power
	case uint32(cfgEnum.ERankType_ChampionshipBattle):
		return this.pbData.Battle
	}
	return nil
}

func (this *PlayerSystemChampionshipFun) ChampionshipTaskRewardRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.ChampionshipTaskRewardRequest)
	rsp := response.(*pb.ChampionshipTaskRewardResponse)
	task := this.getTaskStageInfo(req.RankType)

	// 判断参数是否正确
	if task == nil || task.Id != req.TaskID {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_ChampionshipTaskRewardParam, "head: %v, req: %v", head, req)
	}

	// 判断是否完成任务
	if task.State != pb.EmTaskState_ETS_Finish {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_TaskNoFinish, "head: %v, req: %v", head, req)
	}

	// 判断配置是否存在
	taskCfg := cfgData.GetCfgChampionshipTaskConfig(req.RankType, req.TaskID)
	if taskCfg == nil {
		return uerror.NewUErrorf(1, cfgData.GetChampionshipTaskConfigErrorCode(req.RankType), "head: %v, req: %v", head, req)
	}

	// 领取奖励
	task.State = pb.EmTaskState_ETS_Award
	this.getPlayerBagFun().AddArrItem(head, taskCfg.ListAddItem, pb.EmDoingType_EDT_Championship, true)

	// 设置下一个任务
	this.UpdateSave(this.registerTask(task, req.RankType))
	rsp.Task = this.getTaskStageInfo(req.RankType)

	// 返回客户端
	cluster.SendToClient(head, rsp, cfgEnum.ErrorCode_Success)
	return nil
}

func (this *PlayerSystemChampionshipFun) ChampionshipInfoRequest(head *pb.RpcHead, request, response proto.Message) error {
	rsp := response.(*pb.ChampionshipInfoResponse)
	regionID := this.getPlayerBaseFun().GetServerId()

	tmps := map[uint64]struct{}{}
	uids := []uint64{}
	for _, rankType := range rankTypes {
		// 读取配置
		cfg := cfgData.GetCfgRankTypeConfig(rankType)
		if cfg == nil {
			return uerror.NewUErrorf(1, cfgData.GetRankTypeConfigErrorCode(rankType), "RankTypeConfig(%d) not found", rankType)
		}
		infoCfg := cfgData.GetCfgRankInfoConfig(rankType)
		if infoCfg == nil {
			return uerror.NewUErrorf(1, cfgData.GetRankInfoConfigErrorCode(rankType), "RankInfoConfig(%d) not found", rankType)
		}

		// 获取排行榜
		createTime := cfgData.GetCfgRankActiveTime(infoCfg, this.getPlayerBaseFun().GetServerStartTime())
		members, err := rank_info.ZRevRangeWithScores(cfg, regionID, createTime, 0, 0)
		if err != nil {
			return uerror.NewUErrorf(1, cfgEnum.ErrorCode_DATABASE, "RankType: %d, regionID: %d", rankType, regionID)
		}

		if len(members) > 0 {
			// 过滤
			uid := members[0].Display.AccountId
			if _, ok := tmps[uid]; !ok {
				uids = append(uids, uid)
				tmps[uid] = struct{}{}
			}

			// 组装数据
			rsp.List = append(rsp.List, &pb.ChampionshipRankInfo{
				RankType: rankType,
				First:    members[0],
			})
		}
	}

	// 加载玩家信息
	if dis, err := player_base_display.MGet(uids...); err != nil {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_DATABASE, "head: %v, uids: %v", head, uids)
	} else {
		for _, top := range rsp.List {
			if vv, ok := dis[top.First.Display.AccountId]; ok {
				top.First.Display = vv
			}
		}
	}

	return nil
}
