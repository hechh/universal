package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	report2 "corps/common/report"
	"corps/common/serverCommon"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/module/entry"
	"corps/server/game/module/reward"
)

var MAX_PRIZE_STAGE uint32 = 3

type (
	PlayerSystemBattleNormalFun struct {
		*PlayerSystemBattleFun
		pbData       *pb.PBBattleNormalInfo
		uBattleState base.STAT_TYPE
		uBattleTime  uint64
	}
)

func (this *PlayerSystemBattleNormalFun) Init(pFun *PlayerSystemBattleFun) {
	this.pbData = &pb.PBBattleNormalInfo{MapInfo: &pb.PBBattleMapInfo{}}
	this.PlayerSystemBattleFun = pFun
}

func (this *PlayerSystemBattleNormalFun) GetRankValue() uint64 {
	//为了正向排序用，客户端展示需要用9999减去这个数值
	mapInfo := this.pbData.MapInfo
	if mapInfo.UseTime >= 9999 {
		mapInfo.UseTime = 9999
	}

	uValue := uint64(9999-mapInfo.UseTime) + uint64(mapInfo.StageRate)*10000
	uValue += uint64(mapInfo.StageId)*1000000000 + uint64(mapInfo.MapId)*1000000000000 //useTime:4位 rate：5位 stageid:3位 mapid:3位
	return uValue
}

// 获取通关关卡
func (this *PlayerSystemBattleNormalFun) GetFinishMapIdAndStageId() (uint32, uint32) {
	if this.pbData.MapInfo.IsSuceess > 0 {
		return this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId
	}

	cfgPre := cfgData.GetPreCfgBattleNormal(this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId)
	if cfgPre == nil {
		return 0, 0
	}
	return cfgPre.MapId, cfgPre.StageId
}

// 获取进入关卡
func (this *PlayerSystemBattleNormalFun) GetMapIdAndStageId() (uint32, uint32) {
	if this.pbData.MapInfo.IsSuceess > 0 {
		cfgNext := cfgData.GetNextCfgBattleNormal(this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId)
		if cfgNext == nil {
			return this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId
		}

		return cfgNext.MapId, cfgNext.StageId
	}

	return this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId
}

func (this *PlayerSystemBattleNormalFun) GetBattleMapInfo() *pb.PBBattleMapInfo {
	return this.pbData.MapInfo
}
func (this *PlayerSystemBattleNormalFun) LoadPlayerDBFinish() {
	if this.pbData.MapInfo == nil || this.pbData.MapInfo.MapId == 0 {
		this.pbData.MapInfo = &pb.PBBattleMapInfo{
			MapId:     1,
			StageId:   1,
			Time:      base.GetNow(),
			StageRate: 0,
			UseTime:   0,
		}
		this.UpdateBattleRank(pb.EmBattleType_EBT_Normal, this.pbData.MapInfo)

		//初始化等级
		this.UpdateSave(true)
	}

	if this.pbData.PrizeMapId == 0 {
		this.pbData.PrizeMapId = 1
		this.pbData.PrizeStageId = 1
	}
}

func (this *PlayerSystemBattleNormalFun) loadData(pbData *pb.PBPlayerSystemBattle) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemBattle{}
	}
	if pbData.BattleNormal == nil {
		pbData.BattleNormal = &pb.PBBattleNormalInfo{}
	}
	this.pbData = pbData.BattleNormal

	if this.pbData.MapInfo == nil {
		this.pbData.MapInfo = &pb.PBBattleMapInfo{
			MapId:      1,
			StageId:    1,
			Time:       base.GetNow(),
			StageRate:  0,
			UseTime:    0,
			FightCount: 0,
		}
	}

	this.UpdateSave(true)
}

func (this *PlayerSystemBattleNormalFun) LoadComplete() {
	this.getPlayerSystemCommonFun().CheckSystemOpen(&pb.RpcHead{Id: this.AccountId}, cfgEnum.ESystemOpenType_BattleNormalEnd)
	this.getPlayerSystemCommonFun().CheckSystemOpen(&pb.RpcHead{Id: this.AccountId}, cfgEnum.ESystemOpenType_FBattleNormalEnter)
}

// 存储到数据库
func (this *PlayerSystemBattleNormalFun) SaveData(pbData *pb.PBPlayerSystemBattle) {
	pbData.BattleNormal = this.pbData
}

// 检查参数 不能超过最大的关卡
func (this *PlayerSystemBattleNormalFun) CheckMapStage(mapId uint32, stageId uint32) cfgEnum.ErrorCode {
	cfgBattle := cfgData.GetCfgBattleNormal(mapId, stageId)
	if cfgBattle == nil {
		return plog.Print(this.AccountId, cfgData.GetBattleNormalErrorCode(mapId), mapId, stageId)
	}

	//必须打下一关
	if this.pbData.MapInfo.IsSuceess > 0 {
		//判断下一关
		cfgNext := cfgData.GetNextCfgBattleNormal(this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId)
		if cfgNext == nil {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxBattleMap, mapId, stageId)
		}

		//重置数据
		this.pbData.MapInfo = &pb.PBBattleMapInfo{
			MapId:     cfgNext.MapId,
			StageId:   cfgNext.StageId,
			IsSuceess: 0,
		}

	}

	if base.MakeU64(mapId, stageId) != base.MakeU64(this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_BattleRepeated, mapId, stageId)
	}

	return cfgEnum.ErrorCode_Success
}

// 挑战开始请求
func (this *PlayerSystemBattleNormalFun) BattleBegin(head *pb.RpcHead, mapId uint32, stageId uint32, params []uint32) cfgEnum.ErrorCode {
	if this.uBattleState != base.STAT_TYPE_Begin && this.uBattleTime+base.BATTLE_TIME_OUT > base.GetNow() {
		//return plog.Print(head.Id, cfgEnum.ErrorCode_BattleIng, mapId, stageId, params)
	}

	uCode := this.CheckMapStage(mapId, stageId)
	if uCode != cfgEnum.ErrorCode_Success {
		return uCode
	}

	this.uBattleState = base.STAT_TYPE_Ing
	this.uBattleTime = base.GetNow()

	this.pbData.MapInfo.FightCount++
	if this.pbData.MapInfo.FightCount == 1 {
		this.getPlayerSystemCommonFun().CheckSystemOpen(head, cfgEnum.ESystemOpenType_FBattleNormalEnter)
	}
	return cfgEnum.ErrorCode_Success
}

// 挑战结束请求
func (this *PlayerSystemBattleNormalFun) BattleEnd(head *pb.RpcHead, battleInfo *pb.BattleInfo) cfgEnum.ErrorCode {
	if this.uBattleState != base.STAT_TYPE_Ing {
		//return plog.Print(head.Id, cfgEnum.ErrorCode_PARAM, battleInfo)
	}

	uCode := this.CheckMapStage(battleInfo.MapId, battleInfo.StageId)
	if uCode != cfgEnum.ErrorCode_Success {
		return uCode
	}

	cfgBattle := cfgData.GetCfgBattleNormal(battleInfo.MapId, battleInfo.StageId)
	if cfgBattle == nil {
		return plog.Print(head.Id, cfgData.GetBattleNormalErrorCode(battleInfo.MapId), battleInfo)
	}

	bOpenHook := this.getPlayerSystemCommonFun().CheckSystemTypeOpen(cfgEnum.ESystemUnlockType_AFKSys)
	uCode = this.OnBattleEnd(head, battleInfo, this.pbData.MapInfo)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(head.Id, uCode, battleInfo)
	}

	//给奖励 小怪 精英 boss
	arrItemInfo := make([]*common.ItemInfo, 0)

	//先给首通奖励
	if battleInfo.IsSucc > 0 && len(cfgBattle.AddFirstPrize) > 0 {
		arrItemInfo = append(arrItemInfo, cfgBattle.AddFirstPrize...)
	}

	for i := 0; i < len(battleInfo.MonsterInfo); i++ {
		tmp := battleInfo.MonsterInfo[i]
		if tmp.MaxCount <= 0 {
			continue
		}
		if tmp.MonsterType < uint32(len(cfgBattle.MonsterLootGroupId)) {
			tmpGroup := cfgBattle.MonsterLootGroupId[tmp.MonsterType]
			if len(tmpGroup) != 2 || tmpGroup[0] == 0 {
				continue
			}

			arrItemInfo = append(arrItemInfo, &common.ItemInfo{
				Kind:   uint32(cfgEnum.ESystemType_LootGroup),
				Id:     tmpGroup[0],
				Count:  int64(tmp.KillCount),
				Params: []uint32{tmpGroup[1]},
			})
		}
	}
	pbResponse := &pb.BattleEndResponse{
		PacketHead: &pb.IPacket{},
		BattleType: pb.EmBattleType_EBT_Normal,
		MapId:      this.pbData.MapInfo.MapId,
		StageId:    this.pbData.MapInfo.StageId,
	}

	if len(arrItemInfo) > 0 {
		pbResponse.ItemInfo = this.getPlayerBagFun().GetPbItems(arrItemInfo, pb.EmDoingType_EDT_BattleNormal)
		//词条奖励加成
		uEntryAddRate := entry.ToValue(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_BattleEnd_IncreaseReward), uint32(cfgEnum.EBattleType_BattleNormal))...)
		if uEntryAddRate > 0 {
			arrExtra := reward.AddProbCommonItem(uEntryAddRate, arrItemInfo...)
			if len(arrExtra) > 0 {
				pbResponse.ItemInfo = append(pbResponse.ItemInfo, this.getPlayerBagFun().GetPbItems(arrExtra, pb.EmDoingType_EDT_Entry)...)
			}
		}

		// 发放奖励
		this.getPlayerBagFun().AddPbItems(head, pbResponse.ItemInfo, pb.EmDoingType_EDT_BattleNormal, false)
	}

	if battleInfo.IsSucc > 0 {
		this.onBattleEndSuccess(head)
	}
	this.uBattleState = base.STAT_TYPE_Begin
	this.uBattleTime = 0
	//同步客户端
	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)

	//增加离线时间 第一次解锁不给奖励
	if bOpenHook && this.getPlayerSystemCommonFun().CheckSystemTypeOpen(cfgEnum.ESystemUnlockType_AFKSys) {
		this.getPlayerSystemOffline().AddOfflineSeconds(head, battleInfo.UseTime, pb.EmDoingType_EDT_BattleNormal)
	}
	return cfgEnum.ErrorCode_Success
}

// 挑战结束请求
func (this *PlayerSystemBattleNormalFun) onBattleEndSuccess(head *pb.RpcHead) {

	//进入的ID
	mapId, stageId := this.GetMapIdAndStageId()
	// 请求机器人解锁
	this.getPlayerCrystalFun().UnlockRobot(head, mapId, stageId, true)

	//系统开关
	this.getPlayerSystemCommonFun().CheckSystemOpen(head, cfgEnum.ESystemOpenType_BattleNormalEnd)

	//抽奖触发
	fmapId, fstageId := this.GetFinishMapIdAndStageId()
	this.getPlayerSystemDrawFun().OnBattleEnd(head, pb.EmBattleType_EBT_Normal, fmapId, fstageId)

	//成就类型
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_BattleMap, 1, uint32(pb.EmBattleType_EBT_Normal), serverCommon.MAKE_BATTLE_MAP(this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId))
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_BattleNormalRate, 1, serverCommon.MAKE_BATTLE_MAP(this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId), this.pbData.MapInfo.StageRate)
}

// 领取精英关卡奖励请求
func (this *PlayerSystemBattleNormalFun) NormalBattlePrizeRequest(head *pb.RpcHead, pbRequest *pb.NormalBattlePrizeRequest) {
	uCode := this.NormalBattlePrize(head, pbRequest.MapId, pbRequest.StageId, pbRequest.PrizeStage)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.NormalBattlePrizeResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 领取精英关卡奖励请求
func (this *PlayerSystemBattleNormalFun) NormalBattlePrize(head *pb.RpcHead, mapId uint32, stageId uint32, prizeStage uint32) cfgEnum.ErrorCode {
	if mapId != this.pbData.PrizeMapId || stageId != this.pbData.PrizeStageId || prizeStage > MAX_PRIZE_STAGE || prizeStage <= 0 {
		return plog.Print(head.Id, cfgEnum.ErrorCode_NormalBattlePrizeParam, mapId, stageId, prizeStage)
	}

	for i := 0; i < len(this.pbData.PrizeStage); i++ {
		if this.pbData.PrizeStage[i] == prizeStage {
			return plog.Print(head.Id, cfgEnum.ErrorCode_HavePrize, mapId, stageId, prizeStage)
		}
	}

	if this.pbData.PrizeMapId == this.pbData.MapInfo.MapId {
		if this.pbData.PrizeStageId > this.pbData.MapInfo.StageId {
			return plog.Print(head.Id, cfgEnum.ErrorCode_NormalBattlePrizeStageId, mapId, stageId, prizeStage)
		} else if this.pbData.PrizeStageId == this.pbData.MapInfo.StageId {
			if prizeStage*base.MIL_PERCENT/MAX_PRIZE_STAGE > this.pbData.MapInfo.StageRate {
				return plog.Print(head.Id, cfgEnum.ErrorCode_NormalBattleMaxPrizeStageId, mapId, stageId, prizeStage)
			}
		}
	}

	cfgBattle := cfgData.GetCfgBattleNormal(this.pbData.PrizeMapId, this.pbData.PrizeStageId)
	if cfgBattle == nil {
		return plog.Print(head.Id, cfgData.GetBattleNormalErrorCode(this.pbData.PrizeMapId), mapId, stageId, prizeStage)
	}

	//加奖励
	this.getPlayerBagFun().AddArrItem(head, cfgBattle.MapAddItem[prizeStage], pb.EmDoingType_EDT_Battle, true)

	//存数据
	this.pbData.PrizeStage = append(this.pbData.PrizeStage, prizeStage)
	if len(this.pbData.PrizeStage) >= int(MAX_PRIZE_STAGE) {
		cfgNext := cfgData.GetNextCfgBattleNormal(this.pbData.PrizeMapId, this.pbData.PrizeStageId)
		if cfgNext != nil {
			this.pbData.PrizeMapId = cfgNext.MapId
			this.pbData.PrizeStageId = cfgNext.StageId
			this.pbData.PrizeStage = make([]uint32, 0)
		}
	}

	this.UpdateSave(true)

	//成就统计
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_BattleNormalPrizeCount, 1)

	cluster.SendToClient(head, &pb.NormalBattlePrizeResponse{
		PacketHead:   &pb.IPacket{},
		PrizeMapId:   this.pbData.PrizeMapId,
		PrizeStageId: this.pbData.PrizeStageId,
		PrizeStage:   this.pbData.PrizeStage,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 获取总领奖次数
func (this *PlayerSystemBattleNormalFun) GetBattleNormalTotalPrizeCount() uint32 {
	uStageCount := cfgData.GetCfgBattleNormalTotalStage(this.pbData.PrizeMapId)
	if this.pbData.PrizeStageId > 0 {
		uStageCount += this.pbData.PrizeStageId - 1
	}

	return uStageCount*MAX_PRIZE_STAGE + uint32(len(this.pbData.PrizeStage))
}

// 设置玩家数据
func (this *PlayerSystemBattleNormalFun) OnSetUserTypeInfo() {
	this.onBattleEndSuccess(&pb.RpcHead{Id: this.AccountId})
}

// 设置玩家数据
func (this *PlayerSystemBattleNormalFun) FinishAllBattle(head *pb.RpcHead) {
	cfgMax := cfgData.GetMaxCfgBattleNormal()
	if cfgMax == nil {
		return
	}

	this.pbData.MapInfo.MapId = cfgMax.MapId
	this.pbData.MapInfo.StageId = cfgMax.StageId
	this.UpdateSave(true)
	this.UpdateBattleRank(pb.EmBattleType_EBT_Normal, this.pbData.MapInfo)
}

// 战斗进度保存请求
func (this *PlayerSystemBattleNormalFun) BattleNormalCardRequest(head *pb.RpcHead, res *pb.BattleNormalCardRequest) {
	//数据上报
	report2.Send(head, &report2.ReportBattleCard{
		Stage:     res.Stage,
		CardId:    res.CardID,
		BeginTime: res.BattleBeginTime,
	})

	plog.Print(this.AccountId, cfgEnum.ErrorCode_Success, this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId, res.Stage, res.CardID, res.BattleBeginTime)
	cluster.SendToClient(head, &pb.BattleNormalCardResponse{
		PacketHead: &pb.IPacket{},
	}, cfgEnum.ErrorCode_Success)
}
