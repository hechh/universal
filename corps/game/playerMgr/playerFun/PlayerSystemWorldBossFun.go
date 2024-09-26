package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/common/dao/rank_info"
	"corps/common/orm/redis"
	"corps/common/serverCommon"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
)

type (
	PlayerSystemWorldBossFun struct {
		PlayerFun
		pbData *pb.PBPlayerSystemWorldBoss
	}
)

func (this *PlayerSystemWorldBossFun) Init(pbType pb.PlayerDataType, pCommon *FunCommon) {
	this.PlayerFun.Init(pbType, pCommon)
}

// 从数据库中加载
func (this *PlayerSystemWorldBossFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	this.loadData(pbSystem.WorldBoss)
	this.UpdateSave(false)
}

func (this *PlayerSystemWorldBossFun) loadData(pbData *pb.PBPlayerSystemWorldBoss) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemWorldBoss{}
	}

	this.pbData = pbData

	this.UpdateSave(true)
}

// 加载完成
func (this *PlayerSystemWorldBossFun) LoadPlayerDBFinish() {
	if this.pbData == nil {
		this.pbData = &pb.PBPlayerSystemWorldBoss{}
	}

	if this.pbData.BossId == 0 {
		uOpenServerWeeks := base.DiffWeeks(this.getPlayerBaseFun().GetOpenSeverTime(), base.GetNow())
		this.pbData.BossId = cfgData.GetCfgWorldBosssMapByWeek(uOpenServerWeeks)
	}
}

// 加载完成
func (this *PlayerSystemWorldBossFun) LoadComplete() {
	if this.pbData.BossId == 0 {
		uOpenServerWeeks := base.DiffWeeks(this.getPlayerBaseFun().GetOpenSeverTime(), base.GetNow())
		this.pbData.BossId = cfgData.GetCfgWorldBosssMapByWeek(uOpenServerWeeks)

		//同步数据
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.WorldBossNotify{
			PacketHead: &pb.IPacket{},
			WorldBoss:  this.pbData,
		}, cfgEnum.ErrorCode_Success)
	}
}

// 存储到数据库
func (this *PlayerSystemWorldBossFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	if pbSystem.WorldBoss == nil {
		pbSystem.WorldBoss = new(pb.PBPlayerSystemWorldBoss)
	}
	pbSystem.WorldBoss = this.pbData

	return this.BSave
}
func (this *PlayerSystemWorldBossFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemWorldBoss{}
}

// 设置玩家数据
func (this *PlayerSystemWorldBossFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}
	pbSystem := pbData.(*pb.PBPlayerSystemWorldBoss)
	if pbSystem == nil {
		return false
	}
	this.loadData(pbSystem)
	this.UpdateSave(true)
	return true
}

func (this *PlayerSystemWorldBossFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = new(pb.PBPlayerSystem)
	}

	this.SaveSystem(pbData.System)
}

// 跨天
func (this *PlayerSystemWorldBossFun) PassDay(isDay, isWeek, isMonth bool) {
	//进度奖励补发
	if this.pbData.DailyMaxDamage > 0 {
		arrAddItem := this.GetStagePrize(&pb.RpcHead{Id: this.AccountId})
		if len(arrAddItem) > 0 {
			this.getPlayerMailFun().AddTempMail(&pb.RpcHead{Id: this.AccountId}, cfgEnum.EMailId_BossDaily, pb.EmDoingType_EDT_WorldBoss, arrAddItem)
		}
	}

	//清理数据
	this.pbData = &pb.PBPlayerSystemWorldBoss{
		BossId:    this.pbData.BossId,
		MaxDamage: this.pbData.MaxDamage,
	}

	//重新随机boss
	if isWeek {
		uOpenServerWeeks := base.DiffWeeks(this.getPlayerBaseFun().GetOpenSeverTime(), base.GetNow())
		this.pbData.BossId = cfgData.GetCfgWorldBosssMapByWeek(uOpenServerWeeks)
		this.getPlayerHeroFun().ClearTeamList(uint32(cfgEnum.EHeroTeam_WorldBoss))
	}

	this.UpdateSave(true)

	//同步数据
	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.WorldBossNotify{
		PacketHead: &pb.IPacket{},
		WorldBoss:  this.pbData,
	}, cfgEnum.ErrorCode_Success)
}

// 七天活跃奖励领取请求
func (this *PlayerSystemWorldBossFun) WorldBossStagePrizeRequest(head *pb.RpcHead, pbRequest *pb.WorldBossStagePrizeRequest) {
	uCode := this.WorldBossStagePrize(head)
	cluster.SendToClient(head, &pb.WorldBossStagePrizeResponse{
		PacketHead:        &pb.IPacket{},
		DailyPrizeStageId: this.pbData.DailyPrizeStageId,
	}, uCode)
}
func (this *PlayerSystemWorldBossFun) GetStagePrize(head *pb.RpcHead) (arrAddItem []*common.ItemInfo) {
	arrAddItem = make([]*common.ItemInfo, 0)
	listCfg := cfgData.GetCfgWorldBossHp(this.pbData.DailyMaxDamage)
	if len(listCfg) <= 0 {
		return nil
	}

	for _, cfg := range listCfg {
		if cfg.Id <= this.pbData.DailyPrizeStageId {
			continue
		}

		if cfg.Hp > this.pbData.DailyMaxDamage {
			break
		}

		this.pbData.DailyPrizeStageId = cfg.Id
		arrAddItem = append(arrAddItem, cfg.AddPrize...)
	}

	arrAddItem = serverCommon.Merge_CommonItem(arrAddItem)
	return
}

// 世界boss阶段奖励请求
func (this *PlayerSystemWorldBossFun) WorldBossStagePrize(head *pb.RpcHead) cfgEnum.ErrorCode {
	arrAddItem := this.GetStagePrize(head)
	if len(arrAddItem) <= 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoPrize)
	}

	this.getPlayerBagFun().AddArrItem(head, arrAddItem, pb.EmDoingType_EDT_WorldBoss, true)
	this.UpdateSave(true)
	return cfgEnum.ErrorCode_Success
}

// 世界boss购买次数请求
func (this *PlayerSystemWorldBossFun) WorldBossBuyCountRequest(head *pb.RpcHead, pbRequest *pb.WorldBossBuyCountRequest) {
	uCode := this.WorldBossBuyCount(head)
	cluster.SendToClient(head, &pb.WorldBossBuyCountResponse{
		PacketHead:    &pb.IPacket{},
		DailyBuyCount: this.pbData.DailyBuyCount,
	}, uCode)
}

// 世界boss购买次数请求
func (this *PlayerSystemWorldBossFun) WorldBossBuyCount(head *pb.RpcHead) cfgEnum.ErrorCode {
	//最大次数
	cfgConst := cfgData.GetCfgWorldBosssConst()
	if cfgConst == nil {
		return plog.Print(this.AccountId, cfgData.GetWorldBosssConstErrorCode(0))
	}

	if this.pbData.DailyBuyCount >= cfgConst.DailyBuyCount+this.getPlayerSystemCommonFun().GetPrivilege(cfgEnum.PrivilegeType_WorldBossFightCount) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_DailyMaxCount)
	}

	//判断进入时间
	uPassSeconds := base.GetZeroSeconds(0)
	if uPassSeconds < cfgConst.BeginSecond || uPassSeconds >= cfgConst.PrizeSecond {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity)
	}

	uCode := this.getPlayerBagFun().DelItem(head, cfgConst.BuyDelItem.Kind, cfgConst.BuyDelItem.Id, cfgConst.BuyDelItem.Count, pb.EmDoingType_EDT_WorldBoss)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode)
	}

	this.pbData.DailyBuyCount++
	this.UpdateSave(true)
	return cfgEnum.ErrorCode_Success
}

// 获取最大奖励次数
func (this *PlayerSystemWorldBossFun) getMaxPrizeCount() uint32 {
	return cfgData.GetCfgWorldBosssConst().DailyEndCount + this.pbData.DailyBuyCount
}

// 世界boss扫荡请求
func (this *PlayerSystemWorldBossFun) WorldBossSweepRequest(head *pb.RpcHead, pbRequest *pb.WorldBossSweepRequest) {
	uCode := this.WorldBossSweep(head)
	cluster.SendToClient(head, &pb.WorldBossSweepResponse{
		PacketHead:      &pb.IPacket{},
		DailyPrizeCount: this.pbData.DailyPrizeCount,
		DailyEnterCount: this.pbData.DailyEnterCount,
	}, uCode)

}

// 世界boss购买次数请求
func (this *PlayerSystemWorldBossFun) WorldBossSweep(head *pb.RpcHead) cfgEnum.ErrorCode {
	//最大次数 todo战令
	cfgConst := cfgData.GetCfgWorldBosssConst()
	if cfgConst == nil {
		return plog.Print(this.AccountId, cfgData.GetWorldBosssConstErrorCode(0))
	}

	//每日挑战结算次数
	if this.pbData.DailyPrizeCount >= this.getMaxPrizeCount() {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_DailyMaxCount)
	}

	if this.pbData.DailyPrizeCount <= 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_DailyMaxCount)
	}

	//判断进入时间
	uPassSeconds := base.GetZeroSeconds(0)
	if uPassSeconds < cfgConst.BeginSecond || uPassSeconds >= cfgConst.PrizeSecond {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity)
	}

	this.pbData.DailyPrizeCount++

	//给奖励
	this.OnBattleEnd(head, this.pbData.DailyMaxDamage)

	this.UpdateSave(true)
	return cfgEnum.ErrorCode_Success
}

func (this *PlayerSystemWorldBossFun) GetMaxDamage() uint64 {
	return this.pbData.MaxDamage
}

func (this *PlayerSystemWorldBossFun) OnBattleEnd(head *pb.RpcHead, uDamage uint64) {
	if uDamage > this.pbData.DailyMaxDamage {
		this.pbData.DailyMaxDamage = uDamage
	}

	// 更新历史最大战力排行榜
	if this.pbData.MaxDamage < uDamage {
		this.pbData.MaxDamage = uDamage
	}
	if cfg := cfgData.GetCfgRankInfoConfig(uint32(cfgEnum.ERankType_ChampionshipDamage)); cfg != nil {
		createTime := this.getPlayerBaseFun().GetServerStartTime()
		this.UpdateRank(
			cfgData.GetCfgRankActiveTime(cfg, createTime),
			uint32(cfgEnum.ERankType_ChampionshipDamage),
			this.pbData.MaxDamage,
		)
	}

	this.pbData.DailyTotalDamage += uDamage

	//添加到排行榜中
	if cfg := cfgData.GetCfgRankInfoConfig(uint32(cfgEnum.ERankType_WorldBoss)); cfg != nil {
		createTime := base.GetZeroTimestamp(base.GetNow(), 0)
		this.UpdateRank(
			cfgData.GetCfgRankActiveTime(cfg, createTime),
			uint32(cfgEnum.ERankType_WorldBoss),
			this.pbData.DailyTotalDamage,
		)
	}

	//相关成就回调
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_BattleWorldBoss, 1)

	//给奖励
	cfgBoss := cfgData.GetCfgWorldBosssMap(this.pbData.BossId)
	if cfgBoss == nil {
		plog.Info("(this *PlayerSystemWorldBossFun) OnBattleEnd no boss %d", this.pbData.BossId)
		return
	}

	return
}

// 挑战开始请求
func (this *PlayerSystemWorldBossFun) WorldBossBattleBeginRequest(head *pb.RpcHead, pbRequest *pb.WorldBossBattleBeginRequest) {
	uCode := this.WorldBossBattleBegin(head)
	cluster.SendToClient(head, &pb.WorldBossBattleBeginResponse{
		PacketHead:      &pb.IPacket{},
		DailyEnterCount: this.pbData.DailyEnterCount,
	}, uCode)
}

// 挑战开始请求
func (this *PlayerSystemWorldBossFun) WorldBossBattleBegin(head *pb.RpcHead) cfgEnum.ErrorCode {
	cfgConst := cfgData.GetCfgWorldBosssConst()
	if cfgConst == nil {
		return plog.Print(this.AccountId, cfgData.GetWorldBosssConstErrorCode(0))
	}

	//判断进入次数
	if this.pbData.DailyPrizeCount >= this.getMaxPrizeCount() || this.pbData.DailyPrizeCount+this.pbData.DailyEnterCount > cfgConst.DailyEnterCount+this.getMaxPrizeCount() {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_DailyUpperLimit)
	}

	//判断进入时间
	uPassSeconds := base.GetZeroSeconds(0)
	if uPassSeconds < cfgConst.BeginSecond || uPassSeconds >= cfgConst.EndSecond {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity)
	}

	this.pbData.DailyEnterCount++

	this.UpdateSave(true)
	return cfgEnum.ErrorCode_Success
}

// 挑战开始请求
func (this *PlayerSystemWorldBossFun) WorldBossBattleEndRequest(head *pb.RpcHead, pbRequest *pb.WorldBossBattleEndRequest) {
	uCode := this.WorldBossBattleEnd(head, pbRequest.Battle)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.WorldBossBattleEndResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 挑战开始请求
func (this *PlayerSystemWorldBossFun) WorldBossBattleEnd(head *pb.RpcHead, battleInfo *pb.BattleInfo) cfgEnum.ErrorCode {
	cfgConst := cfgData.GetCfgWorldBosssConst()
	if cfgConst == nil {
		return plog.Print(this.AccountId, cfgData.GetWorldBosssConstErrorCode(0))
	}

	//判断进入次数
	if this.pbData.DailyPrizeCount >= this.getMaxPrizeCount() {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_DailyUpperLimit)
	}

	//判断进入时间
	uPassSeconds := base.GetZeroSeconds(0)
	if uPassSeconds < cfgConst.BeginSecond || uPassSeconds >= cfgConst.PrizeSecond {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity)
	}

	this.pbData.DailyPrizeCount++

	//添加通关记录
	if battleInfo.TotalDamage > this.pbData.DailyMaxDamage {
		heroIdList := make([]uint32, 0)
		if battleInfo.ClientData != nil {
			for _, info := range battleInfo.ClientData.HeroBattleLevel {
				heroIdList = append(heroIdList, info.Key)
			}
		}
		pbBattleData := &pb.PBPlayerBattleData{
			Time:       base.GetNow(),
			UseTime:    battleInfo.UseTime,
			Display:    this.getPlayerBaseFun().GetDisplay(),
			HeroList:   this.getPlayerHeroFun().GetBattleListInfo(heroIdList),
			FightPower: this.getPlayerHeroFun().GetFightPower(),
			ClientData: battleInfo.ClientData,
		}
		byData, err := json.Marshal(pbBattleData)
		if err == nil {
			redis := redis.GetRedisBySeverID(this.getPlayerBaseFun().GetServerId())
			if redis != nil {
				redis.HSet(base.ERK_WorldBoss_Record, this.AccountId, byData)
			}
		}
	}

	//加奖励
	this.OnBattleEnd(head, battleInfo.TotalDamage)

	//增加离线时间
	this.getPlayerSystemOffline().AddOfflineSeconds(head, battleInfo.UseTime, pb.EmDoingType_EDT_BattleNormal)

	this.UpdateSave(true)

	cluster.SendToClient(head, &pb.WorldBossBattleEndResponse{
		PacketHead:      &pb.IPacket{},
		DailyMaxDamage:  this.pbData.DailyMaxDamage,
		DailyPrizeCount: this.pbData.DailyPrizeCount,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 挑战开始请求
func (this *PlayerSystemWorldBossFun) WorldBossRecordRequest(head *pb.RpcHead, pbRequest *pb.WorldBossRecordRequest) {
	uCode := this.WorldBossRecord(head, pbRequest.AccountId, pbRequest.ServerId)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.WorldBossRecordResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 挑战开始请求
func (this *PlayerSystemWorldBossFun) WorldBossRecord(head *pb.RpcHead, AccountId uint64, ServerId uint32) cfgEnum.ErrorCode {
	redisClient := redis.GetRedisBySeverID(ServerId)
	if redisClient == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, AccountId, ServerId)
	}

	byData := redisClient.HGet(base.ERK_WorldBoss_Record, fmt.Sprintf("%d", this.AccountId))
	if byData == "" {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, AccountId, ServerId)
	}

	pbData := &pb.PBPlayerBattleData{}
	err := json.Unmarshal([]byte(byData), pbData)
	if err != nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, AccountId, ServerId)
	}

	cluster.SendToClient(head, &pb.WorldBossRecordResponse{
		PacketHead: &pb.IPacket{},
		RecordInfo: pbData,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 打开boss界面请求
func (this *PlayerSystemWorldBossFun) OpenBossRequest(head *pb.RpcHead) {
	pbResponse := &pb.OpenBossResponse{
		PacketHead: &pb.IPacket{},
	}

	cfgWorldBoss := cfgData.GetCfgRankTypeConfig(uint32(cfgEnum.ERankType_WorldBoss))
	if cfgWorldBoss != nil {
		pbRank, err := rank_info.GetMemberRankInfo(cfgWorldBoss, this.getPlayerBaseFun().GetServerId(), base.GetZeroTimestamp(base.GetNow(), 0), fmt.Sprintf("%d", this.AccountId))
		if err == nil && pbRank != nil {
			pbResponse.WorldBossRank = &pb.PBU32U64{
				Key:   pbRank.Rank,
				Value: pbRank.Value,
			}
		}

	}

	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)
}
