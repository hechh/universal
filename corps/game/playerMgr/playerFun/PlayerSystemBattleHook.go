package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/common/serverCommon"
	"corps/framework/cluster"
	"corps/framework/common/uerror"
	"corps/framework/plog"
	"corps/pb"

	"google.golang.org/protobuf/proto"
)

type (
	PlayerSystemBattleHookFun struct {
		*PlayerSystemBattleFun
		pbData       *pb.PBBattleHookInfo
		uBattleState base.STAT_TYPE
		uBattleTime  uint64
	}
)

func (this *PlayerSystemBattleHookFun) Init(pFun *PlayerSystemBattleFun) {
	this.pbData = &pb.PBBattleHookInfo{MapInfo: &pb.PBBattleMapInfo{}, AutoMap: true}
	this.PlayerSystemBattleFun = pFun
}

func (this *PlayerSystemBattleHookFun) GetRankValue() uint64 {
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
func (this *PlayerSystemBattleHookFun) GetFinishMapIdAndStageId() (uint32, uint32) {
	if this.pbData.MapInfo.IsSuceess > 0 {
		return this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId
	}

	cfgPre := cfgData.GetPreCfgBattleHookStage(this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId)
	if cfgPre == nil {
		return 0, 0
	}
	return cfgPre.MapId, cfgPre.StageId
}

// 获取进入关卡
func (this *PlayerSystemBattleHookFun) GetMapIdAndStageId() (uint32, uint32) {
	if this.pbData.MapInfo.IsSuceess > 0 {
		cfgNext := cfgData.GetNextCfgBattleHookStage(this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId)
		if cfgNext == nil {
			return this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId
		}

		return cfgNext.MapId, cfgNext.StageId
	}

	return this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId
}

func (this *PlayerSystemBattleHookFun) GetBattleMapInfo() *pb.PBBattleMapInfo {
	return this.pbData.MapInfo
}
func (this *PlayerSystemBattleHookFun) LoadPlayerDBFinish() {
	if this.pbData.MapInfo == nil || this.pbData.MapInfo.MapId == 0 {
		this.pbData.MapInfo = &pb.PBBattleMapInfo{
			MapId:     1,
			StageId:   1,
			Time:      base.GetNow(),
			StageRate: 0,
			UseTime:   0,
		}
		this.pbData.CurMapId = 1
		this.pbData.CurStageId = 1
		this.UpdateBattleRank(pb.EmBattleType_EBT_Hook, this.pbData.MapInfo)

		//初始化等级
		this.UpdateSave(true)
	}
}

func (this *PlayerSystemBattleHookFun) loadData(pbData *pb.PBPlayerSystemBattle) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemBattle{}
	}
	if pbData.BattleHook == nil {
		pbData.BattleHook = &pb.PBBattleHookInfo{}
	}
	this.pbData = pbData.BattleHook

	if this.pbData.MapInfo == nil {
		this.pbData.MapInfo = &pb.PBBattleMapInfo{
			MapId:      1,
			StageId:    1,
			Time:       base.GetNow(),
			StageRate:  0,
			UseTime:    0,
			FightCount: 0,
		}
		this.pbData.CurMapId = 1
		this.pbData.CurStageId = 1
	}

	this.UpdateSave(true)
}

func (this *PlayerSystemBattleHookFun) LoadComplete() {
}

// 存储到数据库
func (this *PlayerSystemBattleHookFun) SaveData(pbData *pb.PBPlayerSystemBattle) {
	pbData.BattleHook = this.pbData
}

// 挑战开始请求
func (this *PlayerSystemBattleHookFun) BattleBegin(head *pb.RpcHead, mapId uint32, stageId uint32, params []uint32) cfgEnum.ErrorCode {
	if this.uBattleState != base.STAT_TYPE_Begin && this.uBattleTime+base.BATTLE_TIME_OUT > base.GetNow() {
		//return plog.Print(head.Id, cfgEnum.ErrorCode_BattleIng, mapId, stageId, params)
	}
	uCode := this.CheckMapStage(mapId, stageId)
	if uCode != cfgEnum.ErrorCode_Success {
		return uCode
	}

	//设置当前关卡
	this.pbData.MapInfo.FightCount++
	this.pbData.CurMapId = mapId
	this.pbData.CurStageId = stageId
	this.uBattleState = base.STAT_TYPE_Ing
	this.uBattleTime = base.GetNow()
	return cfgEnum.ErrorCode_Success
}

// 检查参数 不能超过最大的关卡
func (this *PlayerSystemBattleHookFun) CheckMapStage(mapId uint32, stageId uint32) cfgEnum.ErrorCode {
	cfgBattle := cfgData.GetCfgBattleHookStage(mapId, stageId)
	if cfgBattle == nil {
		return plog.Print(this.AccountId, cfgData.GetBattleHookStageErrorCode(mapId), mapId, stageId)
	}

	if this.pbData.MapInfo.IsSuceess > 0 {
		//只能打下一关之前
		if base.MakeU64(mapId, stageId) > base.MakeU64(this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId) {
			//判断是否是下一关
			cfgNext := cfgData.GetNextCfgBattleHookStage(this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId)
			if cfgNext == nil || base.MakeU64(mapId, stageId) != base.MakeU64(cfgNext.MapId, cfgNext.StageId) {
				return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxBattleMap, mapId, stageId)
			}

			//重置数据
			this.pbData.MapInfo = &pb.PBBattleMapInfo{
				MapId:     mapId,
				StageId:   stageId,
				IsSuceess: 0,
			}
		}

	} else {
		//只能这一关关之前
		if base.MakeU64(mapId, stageId) > base.MakeU64(this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId) {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedPreBattleMap, mapId, stageId)
		}
	}
	return cfgEnum.ErrorCode_Success
}

// 挑战结束请求
func (this *PlayerSystemBattleHookFun) BattleEnd(head *pb.RpcHead, battleInfo *pb.BattleInfo) cfgEnum.ErrorCode {
	if this.uBattleState != base.STAT_TYPE_Ing {
		//return plog.Print(head.Id, cfgEnum.ErrorCode_PARAM, battleInfo)
	}
	uCode := this.CheckMapStage(battleInfo.MapId, battleInfo.StageId)
	if uCode != cfgEnum.ErrorCode_Success {
		return uCode
	}

	if base.MakeU64(battleInfo.MapId, battleInfo.StageId) != base.MakeU64(this.pbData.CurMapId, this.pbData.CurStageId) {
		return plog.Print(head.Id, cfgEnum.ErrorCode_NotEqualBattleMap, *battleInfo)
	}

	//判断新通关卡
	bNewStage := false
	if this.pbData.MapInfo.IsSuceess == 0 {
		bNewStage = base.MakeU64(battleInfo.MapId, battleInfo.StageId) == base.MakeU64(this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId)
	}

	cfgBattle := cfgData.GetCfgBattleHookStage(battleInfo.MapId, battleInfo.StageId)
	if cfgBattle == nil {
		return plog.Print(head.Id, cfgData.GetBattleHookStageErrorCode(battleInfo.MapId), battleInfo)
	}

	//如果是最大关卡 重复打最后的
	cfgNextBattle := cfgData.GetNextCfgBattleHookStage(battleInfo.MapId, battleInfo.StageId)
	if cfgNextBattle == nil {
		cfgNextBattle = cfgBattle
	}

	uCode = this.OnBattleEnd(head, battleInfo, this.pbData.MapInfo)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(head.Id, uCode, battleInfo)
	}

	arrPbItems := make([]*pb.PBAddItemData, 0)
	if battleInfo.IsSucc > 0 {
		//判断奖励 首次通关才有奖励
		if bNewStage {
			arrPbItems = this.getPlayerBagFun().GetPbItems(cfgBattle.ArrFirstAddItem, pb.EmDoingType_EDT_Battle)

		} else {
			arrPbItems = this.getPlayerBagFun().GetPbItems(cfgBattle.ArrAddItem, pb.EmDoingType_EDT_Battle)
		}

		this.getPlayerBagFun().AddPbItems(head, arrPbItems, pb.EmDoingType_EDT_Battle, false)
		//自动推关会进一关
		if this.pbData.AutoMap {
			this.pbData.CurMapId = cfgNextBattle.MapId
			this.pbData.CurStageId = cfgNextBattle.StageId
		}

		this.onBattleEndSuccess(head)
	} else {
		//自动推关会退一关，最多退10关
		cfgPreBattle := cfgData.GetPreCfgBattleHookStage(battleInfo.MapId, battleInfo.StageId)
		if cfgPreBattle != nil {
			if cfgData.GetBattleHookStageCount(battleInfo.MapId, battleInfo.StageId, this.pbData.MapInfo.MapId, this.pbData.MapInfo.StageId) < cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_HOOK_PRE_STAGE) {
				this.pbData.CurMapId = cfgPreBattle.MapId
				this.pbData.CurStageId = cfgPreBattle.StageId
			}
		}
	}

	//同步客户端
	cluster.SendToClient(head, &pb.BattleEndResponse{
		PacketHead: &pb.IPacket{},
		BattleType: pb.EmBattleType_EBT_Normal,
		MapId:      this.pbData.CurMapId,
		StageId:    this.pbData.CurStageId,
		ItemInfo:   arrPbItems,
	}, cfgEnum.ErrorCode_Success)

	this.uBattleState = base.STAT_TYPE_Begin
	this.uBattleTime = 0

	this.UpdateSave(true)

	return cfgEnum.ErrorCode_Success
}

// 挂机自动推关设置请求
func (this *PlayerSystemBattleHookFun) HookBattleAutoMapRequest(head *pb.RpcHead) {
	this.pbData.AutoMap = !this.pbData.AutoMap
	cluster.SendToClient(head, &pb.HookBattleAutoMapResponse{
		PacketHead: &pb.IPacket{},
		AutoMap:    this.pbData.AutoMap,
	}, cfgEnum.ErrorCode_Success)
}

// 设置玩家数据
func (this *PlayerSystemBattleHookFun) OnSetUserTypeInfo() {

	this.onBattleEndSuccess(&pb.RpcHead{})
}

// 挑战结束请求
func (this *PlayerSystemBattleHookFun) onBattleEndSuccess(head *pb.RpcHead) {

	mapID, StageId := this.GetFinishMapIdAndStageId()
	//成就类型
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_BattleMap, 1, uint32(pb.EmBattleType_EBT_Hook), serverCommon.MAKE_BATTLE_MAP(mapID, StageId))

	//系统开关
	this.getPlayerSystemCommonFun().CheckSystemOpen(head, cfgEnum.ESystemOpenType_BattleHookEnd)

	this.getPlayerSystemDrawFun().OnBattleEnd(head, pb.EmBattleType_EBT_Hook, mapID, StageId)
}

// 设置玩家数据
func (this *PlayerSystemBattleHookFun) FinishAllBattle(head *pb.RpcHead) {
	cfgMax := cfgData.GetMaxCfgBattleNormal()
	if cfgMax == nil {
		return
	}

	this.pbData.MapInfo.MapId = cfgMax.MapId
	this.pbData.MapInfo.StageId = cfgMax.StageId
	this.UpdateSave(true)
	this.UpdateBattleRank(pb.EmBattleType_EBT_Hook, this.pbData.MapInfo)
}

// 怪物掉落请求
func (this *PlayerSystemBattleHookFun) HookBattleLootRequest(head *pb.RpcHead, pbRequest *pb.HookBattleLootRequest) {
	uCode := this.HookBattleLoot(head, pbRequest.MonsterInfo)
	cluster.SendToClient(head, &pb.HookBattleLootResponse{
		PacketHead: &pb.IPacket{},
	}, uCode)
}

func (this *PlayerSystemBattleHookFun) AddHookEquipmentCount(uAdd uint32) {
	this.pbData.TotalLootCount += uAdd
	this.UpdateSave(true)
}

// 获取掉落概率 百分比
func (this *PlayerSystemBattleHookFun) GetHookEquipmentLootRate(cfgBattle *cfgData.BattleHookMapCfg) uint32 {
	if cfgBattle == nil {
		return 0
	}

	//判断时间
	uCurTime := base.GetNow()

	if this.pbData.BeginLootTime == 0 {
		this.pbData.BeginLootTime = uCurTime
		this.UpdateSave(true)
	}

	//小于一小时 掉落上限不掉落，
	if uCurTime < this.pbData.BeginLootTime+uint64(cfgBattle.LootTotalTime) {
		if this.pbData.TotalLootCount >= cfgBattle.LootDeclineCount {
			return cfgBattle.LootDecline
		}
	} else {
		//大于1小时 重置时间
		this.pbData.BeginLootTime = uCurTime
		this.pbData.TotalLootCount = 0
	}

	this.UpdateSave(true)
	return base.PERCENT

}

// 怪物掉落
func (this *PlayerSystemBattleHookFun) HookBattleLoot(head *pb.RpcHead, monsterList []*pb.BattleKillMonsterInfo) cfgEnum.ErrorCode {
	if len(monsterList) <= 0 {
		return plog.Print(head.Id, cfgEnum.ErrorCode_HookBattleLootParam, monsterList)
	}

	cfgBattle := cfgData.GetCfgBattleHookMap(this.pbData.CurMapId)
	if cfgBattle == nil {
		return plog.Print(head.Id, cfgData.GetBattleHookMapErrorCode(this.pbData.CurMapId), monsterList)
	}

	//衰减规则 一小时掉落装备上限 达到上限概率衰减
	uRate := this.GetHookEquipmentLootRate(cfgBattle)
	if uRate == 0 {
		return plog.Print(head.Id, cfgEnum.ErrorCode_Success, monsterList)
	}

	//给奖励 小怪 精英 boss
	arrItemInfo := make([]*common.ItemInfo, 0)
	for i := 0; i < len(monsterList); i++ {
		tmp := monsterList[i]
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
				Params: []uint32{tmpGroup[1] * uRate / 100},
			})
		}
	}

	if len(arrItemInfo) > 0 {
		this.getPlayerBagFun().AddArrItem(head, arrItemInfo, pb.EmDoingType_EDT_BattleHook, false)
	}

	return cfgEnum.ErrorCode_Success
}

// 关机通关奖励
func (this *PlayerSystemBattleHookFun) BattleHookPassRewardRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.BattleHookPassRewardRequest)
	// 判断是否通关
	mapID, stageID := this.GetFinishMapIdAndStageId()
	if req.MapID > mapID || (req.MapID == mapID && req.StageID > stageID) {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_NotPassBattleHook, "head: %v, req: %v", head, req)
	}
	// 判断是否已经领取奖励
	status := uint64(req.MapID)*10000 + uint64(req.StageID)
	for _, info := range this.pbData.PassRewards {
		if info == status {
			return uerror.NewUErrorf(1, cfgEnum.ErrorCode_BattleHookHasPassReward, "head: %v, req: %v", head, req)
		}
	}
	// 加载配置
	cfg := cfgData.GetCfgBattleHookStage(req.MapID, req.StageID)
	if cfg == nil {
		return uerror.NewUErrorf(1, cfgData.GetBattleHookStageErrorCode(req.MapID), "head: %v, req: %v", head, req)
	}
	// 判单是否有奖励
	if len(cfg.ArrPassAddItem) <= 0 {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_BattleHookNotPassReward, "head: %v, req: %v", head, req)
	}
	// 保存领取状态
	this.pbData.PassRewards = append(this.pbData.PassRewards, status)
	this.UpdateSave(true)
	// 发送奖励
	errCode := this.getPlayerBagFun().AddArrItem(head, cfg.ArrPassAddItem, pb.EmDoingType_EDT_BattleHookPassReward, true)
	if errCode != cfgEnum.ErrorCode_Success {
		return uerror.NewUErrorf(1, errCode, "head: %v, req: %v", head, req)
	}
	return nil
}
