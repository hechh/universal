package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common/cfgData"
	"corps/common/orm/redis"
	"corps/common/serverCommon"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"

	"github.com/golang/protobuf/proto"
)

type (
	PlayerSystemBattleFun struct {
		PlayerFun
		mapBattle  map[pb.EmBattleType]IPlayerSystemBattleFun
		pbSchedule *pb.PBBattleSchedule
		pbRelive   *pb.PBBattleRelive
	}

	IPlayerSystemBattleFun interface {
		Init(pFun *PlayerSystemBattleFun)
		LoadPlayerDBFinish()
		loadData(pbData *pb.PBPlayerSystemBattle)
		LoadComplete()
		SaveData(pbData *pb.PBPlayerSystemBattle)
		OnSetUserTypeInfo()
		GetMapIdAndStageId() (uint32, uint32)       //获取通关关卡
		GetFinishMapIdAndStageId() (uint32, uint32) //获取进入关卡
		GetBattleMapInfo() *pb.PBBattleMapInfo
		BattleBegin(head *pb.RpcHead, mapId uint32, stageId uint32, params []uint32) cfgEnum.ErrorCode
		BattleEnd(head *pb.RpcHead, battleInfo *pb.BattleInfo) cfgEnum.ErrorCode
	}
)

func (this *PlayerSystemBattleFun) getBattleFun(pbType pb.EmBattleType) IPlayerSystemBattleFun {
	return this.mapBattle[pbType]
}

func (this *PlayerSystemBattleFun) GetBattleNoramlFun() *PlayerSystemBattleNormalFun {
	return this.getBattleFun(pb.EmBattleType_EBT_Normal).(*PlayerSystemBattleNormalFun)
}
func (this *PlayerSystemBattleFun) GetBattleHookFun() *PlayerSystemBattleHookFun {
	return this.getBattleFun(pb.EmBattleType_EBT_Hook).(*PlayerSystemBattleHookFun)
}

// 注册
func (this *PlayerSystemBattleFun) RegisterBattleFun() {
	this.mapBattle = make(map[pb.EmBattleType]IPlayerSystemBattleFun)
	for i := pb.EmBattleType_EBT_None; i <= pb.EmBattleType_EBT_Hook; i++ {
		switch i {
		case pb.EmBattleType_EBT_Normal:
			this.mapBattle[i] = new(PlayerSystemBattleNormalFun)
			this.mapBattle[i].Init(this)
		case pb.EmBattleType_EBT_Hook:
			this.mapBattle[i] = new(PlayerSystemBattleHookFun)
			this.mapBattle[i].Init(this)
		}
	}
}
func (this *PlayerSystemBattleFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.RegisterBattleFun()
}

func (this *PlayerSystemBattleFun) GetMapIdAndStageId(battleType pb.EmBattleType) (uint32, uint32) {
	fun := this.getBattleFun(battleType)
	if fun == nil {
		return 0, 0
	}
	return fun.GetMapIdAndStageId()
}
func (this *PlayerSystemBattleFun) GetFinishMapIdAndStageId(battleType pb.EmBattleType) (uint32, uint32) {
	fun := this.getBattleFun(battleType)
	if fun == nil {
		return 0, 0
	}

	return fun.GetFinishMapIdAndStageId()
}
func (this *PlayerSystemBattleFun) LoadPlayerDBFinish() {
	for _, info := range this.mapBattle {
		info.LoadPlayerDBFinish()
	}
	if this.pbRelive == nil {
		this.pbRelive = &pb.PBBattleRelive{AdvestReliveCount: 0, ShareReliveCount: 0}
	}
	if this.pbSchedule == nil {
		this.pbSchedule = &pb.PBBattleSchedule{}
	}
	this.UpdateSave(true)
}

// 新注册
func (this *PlayerSystemBattleFun) NewPlayer() {

}

func (this *PlayerSystemBattleFun) loadData(pbData *pb.PBPlayerSystemBattle) {
	for _, info := range this.mapBattle {
		info.loadData(pbData)
	}
	if nil == pbData.Battlechedule {
		pbData.Battlechedule = &pb.PBBattleSchedule{}
	}
	if nil == pbData.Relive {
		pbData.Relive = &pb.PBBattleRelive{}
	}
	this.pbSchedule = pbData.Battlechedule
	this.pbRelive = pbData.Relive
	this.UpdateSave(true)
}

func (this *PlayerSystemBattleFun) LoadComplete() {
	for _, info := range this.mapBattle {
		info.LoadComplete()
	}

	if this.pbRelive == nil {
		this.pbRelive = &pb.PBBattleRelive{
			AdvestReliveCount: 0,
			ShareReliveCount:  0,
		}
	}
}

// 从数据库中加载
func (this *PlayerSystemBattleFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	if pbSystem.Battle == nil {
		pbSystem.Battle = &pb.PBPlayerSystemBattle{}
	}
	this.loadData(pbSystem.Battle)
	this.UpdateSave(false)
}

// 存储到数据库
func (this *PlayerSystemBattleFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	if pbSystem.Battle == nil {
		pbSystem.Battle = new(pb.PBPlayerSystemBattle)
	}

	for _, info := range this.mapBattle {
		info.SaveData(pbSystem.Battle)
	}

	pbSystem.Battle.Battlechedule = this.pbSchedule
	pbSystem.Battle.Relive = this.pbRelive

	return this.BSave
}

func (this *PlayerSystemBattleFun) GetBattleMapInfo(battleType pb.EmBattleType) *pb.PBBattleMapInfo {
	return this.getBattleFun(battleType).GetBattleMapInfo()
}

// 挑战开始请求
func (this *PlayerSystemBattleFun) BattleBeginRequest(head *pb.RpcHead, pbRequest *pb.BattleBeginRequest) {
	uCode := this.BattleBegin(head, pbRequest.BattleType, pbRequest.MapId, pbRequest.StageId, pbRequest.Params)
	cluster.SendToClient(head, &pb.BattleBeginResponse{
		PacketHead: &pb.IPacket{},
		BattleType: pbRequest.BattleType,
		MapId:      pbRequest.MapId,
		StageId:    pbRequest.StageId,
		Params:     pbRequest.Params,
		FightCount: this.GetBattleMapInfo(pbRequest.BattleType).FightCount,
	}, uCode)
}

// 挑战开始请求
func (this *PlayerSystemBattleFun) BattleBegin(head *pb.RpcHead, battleType pb.EmBattleType, mapId uint32, stageId uint32, params []uint32) cfgEnum.ErrorCode {
	plog.Print(this.AccountId, cfgEnum.ErrorCode_Success, battleType, mapId, stageId, params)
	battleFun := this.getBattleFun(battleType)
	uCode := battleFun.BattleBegin(head, mapId, stageId, params)
	if uCode != cfgEnum.ErrorCode_Success {
		if uCode == cfgEnum.ErrorCode_BattleRepeated {
			this.BattleMapNotify(head, battleType, battleFun.GetBattleMapInfo())
		}
		return uCode
	}

	//进入成就
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_EnterBattle, 1, uint32(battleType))

	this.UpdateSave(true)
	return cfgEnum.ErrorCode_Success
}
func (this *PlayerSystemBattleFun) BattleMapNotify(head *pb.RpcHead, battleType pb.EmBattleType, battleMap *pb.PBBattleMapInfo) {
	cluster.SendToClient(head, &pb.BattleMapNotify{
		PacketHead: &pb.IPacket{},
		MapInfo:    battleMap,
		BattleType: battleType,
	}, cfgEnum.ErrorCode_Success)
}

// 挑战结束请求
func (this *PlayerSystemBattleFun) BattleEndRequest(head *pb.RpcHead, pbRequest *pb.BattleEndRequest) {
	uCode := this.BattleEnd(head, pbRequest.Battle)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.BattleEndResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 挑战结束请求
func (this *PlayerSystemBattleFun) BattleEnd(head *pb.RpcHead, battleInfo *pb.BattleInfo) cfgEnum.ErrorCode {
	plog.Print(this.AccountId, cfgEnum.ErrorCode_Success, *battleInfo)
	battleFun := this.getBattleFun(battleInfo.BattleType)
	if battleFun == nil {
		return plog.Print(head.Id, cfgEnum.ErrorCode_BattleTypeNotSupported, *battleInfo)
	}
	uCode := battleFun.BattleEnd(head, battleInfo)
	if uCode != cfgEnum.ErrorCode_Success {
		if uCode == cfgEnum.ErrorCode_BattleRepeated {
			this.BattleMapNotify(head, battleInfo.BattleType, battleFun.GetBattleMapInfo())
		}

		return uCode
	}

	//取消进度
	this.BattleScheduleEnd(head)

	//成就触发
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_PassBattle, 1, uint32(battleInfo.BattleType))
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_PassBattle, 1, 0)

	//空投成就
	if battleInfo.ClientData.DropBoxCount > 0 {
		this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_BattleNormalDropBox, 1, battleInfo.ClientData.DropBoxCount)
	}

	//使徒ID
	if len(battleInfo.ClientData.CryscalRobotId) > 0 {
		tmp := uint32(0)
		for _, robotId := range battleInfo.ClientData.CryscalRobotId {
			cfgRobot := cfgData.GetCfgCrystalRobot(robotId)
			if cfgRobot != nil {
				this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_CrystalRobotBattle, 1, uint32(battleInfo.BattleType), cfgRobot.Element)
				tmp++
			}
		}
		if tmp > 0 {
			this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_CrystalRobotBattle, uint32(len(battleInfo.ClientData.CryscalRobotId)), uint32(battleInfo.BattleType), uint32(cfgEnum.EHydraElementType_Adapt))
		}
	}

	// 触发杀怪成就触发
	uTotal := uint32(0)
	for _, info := range battleInfo.MonsterInfo {
		if info.KillCount > 0 {
			this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_KillMonster, info.KillCount, uint32(battleInfo.BattleType), info.MonsterType)
			this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_KillMonster, info.KillCount, uint32(cfgEnum.EBattleType_None), info.MonsterType)

			uTotal += info.KillCount
		}
	}

	return cfgEnum.ErrorCode_Success
}
func (this *PlayerSystemBattleFun) OnBattleEnd(head *pb.RpcHead, battleInfo *pb.BattleInfo, mapInfo *pb.PBBattleMapInfo) cfgEnum.ErrorCode {
	//不是打当前关卡，不需要更新数据
	if base.MakeU64(battleInfo.MapId, battleInfo.StageId) != base.MakeU64(mapInfo.MapId, mapInfo.StageId) {
		return cfgEnum.ErrorCode_Success
	}

	//检查作弊 todo

	//成功设置下一关
	bUpdateRank := false
	if battleInfo.IsSucc > 0 {
		//成功需要更新数据
		if battleInfo.UseTime > mapInfo.UseTime {
			mapInfo.UseTime = battleInfo.UseTime
			bUpdateRank = true
		}

		if battleInfo.StageRate > mapInfo.StageRate {
			mapInfo.StageRate = battleInfo.StageRate
			bUpdateRank = true
		}

		mapInfo.Time = base.GetNow()
		//第一次成功才发最佳整容
		if mapInfo.IsSuceess == 0 {
			mapInfo.IsSuceess = 1
			mapInfo.FightCount = 0
			bUpdateRank = true

			heroIdList := make([]uint32, 0)
			if battleInfo.ClientData != nil {
				for _, info := range battleInfo.ClientData.HeroBattleLevel {
					heroIdList = append(heroIdList, info.Key)
				}
			}

			//发送阵容到历史记录中
			cfgConfig := cfgData.GetCfgBattleConfig(uint32(battleInfo.BattleType))
			if cfgConfig != nil && cfgConfig.RecordCount > 0 {
				pbBattleData := &pb.PBPlayerBattleData{
					Time:       base.GetNow(),
					UseTime:    battleInfo.UseTime,
					Display:    this.getPlayerBaseFun().GetDisplay(),
					HeroList:   this.getPlayerHeroFun().GetBattleListInfo(heroIdList),
					FightPower: this.getPlayerHeroFun().GetFightPower(),
					ClientData: battleInfo.ClientData,
				}
				byData, err := proto.Marshal(pbBattleData)
				if err == nil {
					head.RegionID = this.getPlayerBaseFun().GetServerId()
					head.RouteType = uint32(cfgEnum.ERouteType_ServerID)
					cluster.SendToGm(head, "BattleRecordMgr", "AddRecord", battleInfo.BattleType, cfgConfig.RecordCount, battleInfo.MapId, battleInfo.StageId, byData)
				}
			}
		}
	} else {
		//成功后重复打 不处理
		if mapInfo.IsSuceess == 1 {
			return cfgEnum.ErrorCode_Success
		}

		for i := 0; i < len(battleInfo.MonsterInfo); i++ {
			tmp := battleInfo.MonsterInfo[i]
			if tmp.MaxCount <= 0 {
				continue
			}

			if tmp.KillCount > tmp.MaxCount {
				plog.Print(head.Id, cfgEnum.ErrorCode_BattleKillCountError, battleInfo)
			}
		}

		if battleInfo.UseTime > mapInfo.UseTime {
			mapInfo.UseTime = battleInfo.UseTime
			bUpdateRank = true
		}

		if battleInfo.StageRate > mapInfo.StageRate {
			mapInfo.StageRate = battleInfo.StageRate
			bUpdateRank = true
		}

	}

	if bUpdateRank {
		this.UpdateBattleRank(battleInfo.BattleType, mapInfo)
	}

	return cfgEnum.ErrorCode_Success
}

// 挑战结束请求
func (this *PlayerSystemBattleFun) BattleRecordRequest(head *pb.RpcHead, pbRequest *pb.BattleRecordRequest) {
	uCode := this.BattleRecord(head, pbRequest.BattleType, pbRequest.MapId, pbRequest.StageId)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.BattleRecordResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 挑战记录查询请求
func (this *PlayerSystemBattleFun) BattleRecord(head *pb.RpcHead, battleType pb.EmBattleType, mapId uint32, stageId uint32) cfgEnum.ErrorCode {
	redis := redis.GetDefaultRedis()
	if redis == nil {
		return plog.Print(head.Id, cfgEnum.ErrorCode_Fail, battleType, mapId, stageId)
	}

	cfgConfig := cfgData.GetCfgBattleConfig(uint32(battleType))
	if cfgConfig == nil || cfgConfig.RecordCount <= 0 {
		return plog.Print(head.Id, cfgData.GetBattleConfigErrorCode(uint32(battleType)), battleType, mapId, stageId)
	}

	pbResponse := &pb.BattleRecordResponse{
		PacketHead: &pb.IPacket{},
		BattleType: battleType,
		MapId:      mapId,
		StageId:    stageId,
	}

	//对应的redis查询
	strRedisKey := serverCommon.GetBattleRecordRedisKey(battleType, mapId, stageId)
	listData := redis.LRange(strRedisKey, -1*int64(cfgConfig.RecordCount), -1)
	if listData != nil {
		for _, strData := range listData {
			pbData := &pb.PBPlayerBattleData{}
			proto.Unmarshal([]byte(strData), pbData)
			pbResponse.RecordList = append(pbResponse.RecordList, pbData)
		}
	}

	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

func (this *PlayerSystemBattleFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = new(pb.PBPlayerSystem)
	}

	this.SaveSystem(pbData.System)
}
func (this *PlayerSystemBattleFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemBattle{}
}

// 设置玩家数据
func (this *PlayerSystemBattleFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}

	pbBattle := pbData.(*pb.PBPlayerSystemBattle)
	if pbBattle == nil {
		return false
	}

	this.loadData(pbBattle)

	for _, info := range this.mapBattle {
		info.OnSetUserTypeInfo()
	}

	this.UpdateSave(true)
	return true
}

// 挑战结束请求
func (this *PlayerSystemBattleFun) UpdateBattleRank(battleType pb.EmBattleType, mapInfo *pb.PBBattleMapInfo) {
	cfgConfig := cfgData.GetCfgBattleConfig(uint32(battleType))
	if cfgConfig == nil || cfgConfig.RankType <= uint32(cfgEnum.ERankType_None) {
		return
	}

	//为了正向排序用，客户端展示需要用9999减去这个数值
	if mapInfo.UseTime >= 9999 {
		mapInfo.UseTime = 9999
	}

	uValue := uint64(9999-mapInfo.UseTime) + uint64(mapInfo.StageRate)*10000
	uValue += uint64(mapInfo.StageId)*1000000000 + uint64(mapInfo.MapId)*1000000000000 //useTime:4位 rate：5位 stageid:3位 mapid:3位

	//永久榜单
	this.UpdateRank(0, cfgConfig.RankType, uValue)

	//锦标赛只判断精英关卡
	switch battleType {
	case pb.EmBattleType_EBT_Normal:
		if cfg := cfgData.GetCfgRankInfoConfig(uint32(cfgEnum.ERankType_ChampionshipBattle)); cfg != nil {
			createTime := this.getPlayerBaseFun().GetServerStartTime()
			this.UpdateRank(
				cfgData.GetCfgRankActiveTime(cfg, createTime),
				uint32(cfgEnum.ERankType_ChampionshipBattle),
				uint64(uValue),
			)
		}
	case pb.EmBattleType_EBT_Hook:
		if cfg := cfgData.GetCfgRankInfoConfig(uint32(cfgEnum.ERankType_ChampionshipBattleHook)); cfg != nil {
			createTime := this.getPlayerBaseFun().GetServerStartTime()
			this.UpdateRank(
				cfgData.GetCfgRankActiveTime(cfg, createTime),
				uint32(cfgEnum.ERankType_ChampionshipBattleHook),
				uint64(uValue),
			)
		}
	}
}

// 关卡功能购买请求
func (this *PlayerSystemBattleFun) BattleFunBuyRequest(head *pb.RpcHead, pbRequest *pb.BattleFunBuyRequest) {
	uCode := this.BattleFunBuy(head, pbRequest.BattleFunType)
	cluster.SendToClient(head, &pb.BattleFunBuyResponse{
		PacketHead:    &pb.IPacket{},
		BattleFunType: pbRequest.BattleFunType,
	}, uCode)
}

// 挑战记录查询请求
func (this *PlayerSystemBattleFun) BattleFunBuy(head *pb.RpcHead, uBattleFunType uint32) cfgEnum.ErrorCode {
	cfgBattleFun := cfgData.GetCfgBattleFunConfig(uBattleFunType)
	if cfgBattleFun == nil {
		return plog.Print(this.AccountId, cfgData.GetBattleFunConfigErrorCode(uBattleFunType), uBattleFunType)

	}

	uCode := this.getPlayerBagFun().DelItem(head, cfgBattleFun.DelItem.Kind, cfgBattleFun.DelItem.Id, cfgBattleFun.DelItem.Count, pb.EmDoingType_EDT_Battle)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode, uBattleFunType, *cfgBattleFun.DelItem)
	}

	return cfgEnum.ErrorCode_Success
}

// 战斗进度保存请求
func (this *PlayerSystemBattleFun) BattleScheduleSaveRequest(head *pb.RpcHead, resp *pb.BattleScheduleSaveRequest) {
	this.pbSchedule = resp.BattleSchedule

	this.UpdateSave(true)
	cluster.SendToClient(head, &pb.BattleScheduleSaveResponse{
		PacketHead:     &pb.IPacket{},
		BattleSchedule: resp.BattleSchedule,
	}, cfgEnum.ErrorCode_Success)
}

func (this *PlayerSystemBattleFun) BattleScheduleEnd(head *pb.RpcHead) {
	this.pbSchedule = &pb.PBBattleSchedule{}

	this.UpdateSave(true)
}

// 隔天刷新通知
func (this *PlayerSystemBattleFun) PassDay(isDay, isWeek, isMonth bool) {
	this.pbRelive.AdvestReliveCount = 0
	this.pbRelive.ShareReliveCount = 0

	cluster.SendToClient(&pb.RpcHead{}, &pb.BattleReliveNotify{
		PacketHead: &pb.IPacket{},
		Relive:     this.pbRelive,
	}, cfgEnum.ErrorCode_Success)

}
