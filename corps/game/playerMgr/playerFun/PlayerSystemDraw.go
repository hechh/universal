package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/common/serverCommon"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/module/entry"

	"github.com/golang/protobuf/proto"
)

type (
	PlayerSystemDrawFun struct {
		PlayerFun
		mapDrawData map[uint32]*pb.PBDrawInfo
		mapUpLoot   map[uint32][]uint32 //up池
	}
)

func (this *PlayerSystemDrawFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.mapDrawData = make(map[uint32]*pb.PBDrawInfo)
	this.mapUpLoot = make(map[uint32][]uint32)
}

// 从数据库中加载
func (this *PlayerSystemDrawFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	if pbSystem.Draw == nil {
		this.mapDrawData = make(map[uint32]*pb.PBDrawInfo)
		return
	}

	this.loadData(pbSystem.Draw)

	this.UpdateSave(false)
}
func (this *PlayerSystemDrawFun) loadData(pbData *pb.PBPlayerSystemDraw) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemDraw{}
	}

	this.mapDrawData = make(map[uint32]*pb.PBDrawInfo)
	for _, info := range pbData.DrawList {
		this.mapDrawData[info.DrawId] = info
	}

	this.UpdateSave(true)
}

// 存储到数据库
func (this *PlayerSystemDrawFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	if pbSystem.Draw == nil {
		pbSystem.Draw = new(pb.PBPlayerSystemDraw)
	}
	pbSystem.Draw.DrawList = make([]*pb.PBDrawInfo, 0)
	for _, info := range this.mapDrawData {
		pbSystem.Draw.DrawList = append(pbSystem.Draw.DrawList, info)
	}

	return this.BSave
}
func (this *PlayerSystemDrawFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemDraw{}
}

// 设置玩家数据
func (this *PlayerSystemDrawFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}
	pbSystem := pbData.(*pb.PBPlayerSystemDraw)
	if pbSystem == nil {
		return false
	}
	this.loadData(pbSystem)
	this.UpdateSave(true)
	return true
}

func (this *PlayerSystemDrawFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = new(pb.PBPlayerSystem)
	}

	this.SaveSystem(pbData.System)
}
func (this *PlayerSystemDrawFun) LoadComplete() {
	this.CheckDrawOpen(&pb.RpcHead{Id: this.AccountId})
	this.checkUpLoot()

}
func (this *PlayerSystemDrawFun) PassDay(isDay, isWeek, isMonth bool) {
	this.CheckDrawOpen(&pb.RpcHead{Id: this.AccountId})
}
func (this *PlayerSystemDrawFun) GetUpLoot(uGroupId uint32) []uint32 {
	return this.mapUpLoot[uGroupId]
}

// 检查up池子
func (this *PlayerSystemDrawFun) checkUpLoot() {
	mapAllCfg := cfgData.GetAllCfgDraw()

	//英雄up池子 有日期
	this.mapUpLoot = make(map[uint32][]uint32)
	uNow := base.GetNow()
	for _, cfg := range mapAllCfg {
		if cfg.DrawType != uint32(cfgEnum.EDrawType_HeroActivity) {
			continue
		}

		//判断是否结束
		if cfg.EndTime > uNow {
			continue
		}

		for k, v := range cfg.MapUpLoot {
			if _, ok := this.mapUpLoot[k]; !ok {
				this.mapUpLoot[k] = make([]uint32, 0)
			}

			this.mapUpLoot[k] = append(this.mapUpLoot[k], v)
		}
	}
}

func (this *PlayerSystemDrawFun) OnSystemTypeOpen(head *pb.RpcHead, uSystemType uint32) {
	mapAllCfg := cfgData.GetAllCfgDraw()
	bCheck := false
	for id, cfg := range mapAllCfg {
		if cfg.ConditionInfo[0].Type != uint32(cfgEnum.EConditionType_SystemOpenType) {
			continue
		}

		if cfg.ConditionInfo[0].SystemOpenType != uSystemType {
			continue
		}

		if _, ok := this.mapDrawData[id]; ok {
			continue
		}

		bCheck = true
		break
	}

	if bCheck {
		this.CheckDrawOpen(head)
	}
}

func (this *PlayerSystemDrawFun) OnBattleEnd(head *pb.RpcHead, battleType pb.EmBattleType, mapId uint32, stageId uint32) {
	mapAllCfg := cfgData.GetAllCfgDraw()
	bCheck := false
	for id, cfg := range mapAllCfg {
		if cfg.ConditionInfo[0].Type != uint32(cfgEnum.EConditionType_BattleMap) {
			continue
		}
		//检查条件
		if uCode, _, _ := this.getPlayerBaseFun().CheckCondition(cfg.ConditionInfo); uCode != cfgEnum.ErrorCode_Success {
			continue
		}

		if _, ok := this.mapDrawData[id]; ok {
			continue
		}

		bCheck = true
		break
	}

	if bCheck {
		this.CheckDrawOpen(head)
	}
}

// 检查抽奖开启关闭
func (this *PlayerSystemDrawFun) CheckDrawOpen(head *pb.RpcHead) {
	mapAllCfg := cfgData.GetAllCfgDraw()
	uCurTime := base.GetNow()

	uOpenServerTime := serverCommon.GetOpenServerTime()
	uOpenServerDays := serverCommon.GetOpenServerDays()
	pbResponse := &pb.DrawNotify{}

	for id, cfg := range mapAllCfg {
		uBeginTime := base.GetNow()
		uEndTime := uint64(0)
		if uCurTime < cfg.OpenTime {
			continue
		}

		if cfg.EndTime > 0 && uCurTime > cfg.EndTime {
			continue
		}

		//检查条件
		if uCode, _, _ := this.getPlayerBaseFun().CheckCondition(cfg.ConditionInfo); uCode != cfgEnum.ErrorCode_Success {
			continue
		}

		//是否结束
		if cfg.EndConditionInfo != nil {
			if uCode, _, _ := this.getPlayerBaseFun().CheckCondition([]*common.ConditionInfo{cfg.EndConditionInfo}); uCode == cfgEnum.ErrorCode_Success {
				//结束强制修改endtime 最多存在一天
				pDraw := this.mapDrawData[id]
				if pDraw != nil {
					pDraw.EndTime = pDraw.BeginTime + 24*3600
				}

				continue
			}
		}

		//开服时间 需要检查循环
		if cfg.ConditionInfo[0].Type == uint32(cfgEnum.EConditionType_OpenServerDay) {
			uBeginDay := uint32(base.Min(int(cfg.OpenTime), 1))

			//算循环
			if cfg.CircleDay > 0 && cfg.ContinueDay > 0 {
				uCircle := (uOpenServerDays - uBeginDay) / cfg.CircleDay
				uBeginTime = uOpenServerTime + uint64((uBeginDay+uCircle*cfg.CircleDay-1)*24*3600)
				uEndTime = uBeginTime + uint64(cfg.ContinueDay*24*3600) - 1
			} else {
				uBeginTime = uOpenServerTime + uint64((uBeginDay-1)*24*3600)
				if cfg.EndTime > 0 {
					uEndTime = uOpenServerTime + cfg.EndTime*24*3600
				}
			}
		}

		if uBeginTime < cfg.OpenTime {
			uBeginTime = cfg.OpenTime
		}

		if uEndTime == 0 || uEndTime > cfg.EndTime {
			uEndTime = cfg.EndTime
		}

		//检查是否是新增
		pDraw := this.mapDrawData[id]
		if pDraw != nil {
			if pDraw.EndTime != uEndTime {
				pDraw = &pb.PBDrawInfo{
					DrawId:    id,
					BeginTime: uBeginTime,
					EndTime:   uEndTime,
				}
				pbResponse.DrawList = append(pbResponse.DrawList, this.mapDrawData[id])
			}

		} else {
			this.mapDrawData[id] = &pb.PBDrawInfo{
				DrawId:    id,
				BeginTime: uBeginTime,
				EndTime:   uEndTime,
			}

			pDraw = this.mapDrawData[id]
			pbResponse.DrawList = append(pbResponse.DrawList, this.mapDrawData[id])
		}
	}

	//判断删除
	for did, info := range this.mapDrawData {
		if info.EndTime > 0 && info.EndTime < uCurTime {
			delete(this.mapDrawData, did)
			pbResponse.DelDrawList = append(pbResponse.DelDrawList, did)
		}
	}

	if len(pbResponse.DrawList) > 0 || len(pbResponse.DelDrawList) > 0 {
		cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)
	}
}

// 抽奖请求
func (this *PlayerSystemDrawFun) DrawRequest(head *pb.RpcHead, pbRequest *pb.DrawRequest) {
	uCode := this.Draw(head, pbRequest.DrawId, pbRequest.DrawCount, pbRequest.AdvertType)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.DrawResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}

}

// 抽奖请求
func (this *PlayerSystemDrawFun) Draw(head *pb.RpcHead, uDrawId uint32, uDrawCount uint32, uAdvertType uint32) cfgEnum.ErrorCode {
	cfgDraw := cfgData.GetCfgDraw(uDrawId)
	if cfgDraw == nil {
		return plog.Print(this.AccountId, cfgData.GetDrawErrorCode(uDrawId), uDrawId, uDrawCount)
	}

	arrNeedItem, ok := cfgDraw.MapDelItem[uDrawCount]
	if !ok {
		return plog.Print(this.AccountId, cfgData.GetDrawErrorCode(uDrawId), uDrawId, uDrawCount)
	}

	if arrNeedItem == nil || arrNeedItem.Id <= 0 || arrNeedItem.Count <= 0 {
		return plog.Print(this.AccountId, cfgData.GetDrawErrorCode(uDrawId), uDrawId, uDrawCount)
	}

	//英雄需要判断背包
	if cfgDraw.DrawType == uint32(cfgEnum.EDrawType_HeroActivity) || cfgDraw.DrawType == uint32(cfgEnum.EDrawType_Hero) {
		if this.getPlayerHeroFun().GetSpareBag() < uDrawCount {
			cluster.SendCommonToClient(head, cfgEnum.ErrorCode_HeroBagFull)
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_HeroBagFull, uDrawId, uDrawCount)
		}
	}

	drawInfo, ok := this.mapDrawData[uDrawId]
	if !ok {
		drawInfo = &pb.PBDrawInfo{
			DrawId: uDrawId,
		}
		this.mapDrawData[uDrawId] = drawInfo
	}

	//如果是广告 需要检查冷却时间
	if uAdvertType > uint32(cfgEnum.EAdvertType_None) {
		uCurTime := base.GetNow()
		if cfgDraw.AdvertingCoolTime <= 0 {
			return plog.Print(head.Id, cfgEnum.ErrorCode_NotSupportAdvert, uDrawId)
		}

		if drawInfo.AdvertNextTime > uCurTime {
			return plog.Print(head.Id, cfgEnum.ErrorCode_CooltimeIng, uDrawId)
		}

		uDrawCount = 1
		drawInfo.AdvertNextTime = uCurTime + uint64(cfgDraw.AdvertingCoolTime)
	} else {

		//扣道具
		uCode := this.getPlayerBagFun().DelItem(head, arrNeedItem.Kind, arrNeedItem.Id, arrNeedItem.Count, pb.EmDoingType_EDT_Draw)
		if uCode != cfgEnum.ErrorCode_Success {
			return plog.Print(this.AccountId, uCode, uDrawId, uDrawCount, *arrNeedItem)
		}
	}

	//加奖励
	mapAddLootGroupId := make(map[uint32]uint32)
	// 获取词条效果
	dropProb := entry.KeyValueToMap(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_DrawDropProbability), uDrawId)...)
	for i := uint32(0); i < uDrawCount; i++ {
		drawInfo.DrawCount++
		drawInfo.GuarCount++
		drawInfo.Guar2Count++
		drawInfo.Guar3Count++

		//获取奖池
		mapLootGroup := cfgDraw.MapRandLootGroup

		//前N次
		if cfgDraw.MaxFirstCount > 0 && drawInfo.DrawCount <= cfgDraw.MaxFirstCount {
			for _, tmp := range cfgDraw.ListFirstCountLootGroup {
				if drawInfo.DrawCount <= tmp.Key {
					mapLootGroup = make(map[uint32]uint32)
					mapLootGroup[tmp.Value] = base.MIL_PERCENT
					break
				}
			}
		}

		if cfgDraw.Guar3Count > 0 && drawInfo.Guar3Count >= cfgDraw.Guar3Count {
			drawInfo.Guar3Count = 0
			mapLootGroup = cfgDraw.MapRandGuar3LootGroup
		} else if cfgDraw.Guar2Count > 0 && drawInfo.Guar2Count >= cfgDraw.Guar2Count {
			drawInfo.Guar2Count = 0
			mapLootGroup = cfgDraw.MapRandGuar2LootGroup
		} else if cfgDraw.GuarCount > 0 && drawInfo.GuarCount >= cfgDraw.GuarCount {
			drawInfo.GuarCount = 0
			mapLootGroup = cfgDraw.MapRandGuarLootGroup
		}
		if len(mapLootGroup) <= 0 {
			plog.Print(this.AccountId, 0, "lootGroup is nill", uDrawId, uDrawCount)
			continue
		}

		//随机一个lootgroupid
		wrand := base.NewWeightedRandom()
		for key, value := range mapLootGroup {
			if _, ok := dropProb[key]; ok {
				value += dropProb[key]
			}
			wrand.Add(key, value)
		}

		uRandLootId := wrand.GetRandomKey()
		plog.Trace("draw uid:%d drawid:%d count:%d randlootid:%d", this.AccountId, uDrawId, i, uRandLootId)

		//判断是否重置guar,在池子里内就重置
		bHaveGua := false
		if _, ok := cfgDraw.MapRandGuar3LootGroup[uRandLootId]; !bHaveGua && ok {
			drawInfo.Guar3Count = 0
			bHaveGua = true
		}

		if _, ok := cfgDraw.MapRandGuar2LootGroup[uRandLootId]; !bHaveGua && ok {
			drawInfo.Guar2Count = 0
			bHaveGua = true
		}

		if _, ok := cfgDraw.MapRandGuarLootGroup[uRandLootId]; !bHaveGua && ok {
			drawInfo.GuarCount = 0
			bHaveGua = true
		}

		mapAddLootGroupId[uRandLootId] += 1
	}

	//添加奖励
	arrItemInfo := make([]*common.ItemInfo, 0)
	for lootGroupId, lootGroupCount := range mapAddLootGroupId {
		arrItemInfo = append(arrItemInfo, &common.ItemInfo{
			Kind:  uint32(cfgEnum.ESystemType_LootGroup),
			Id:    lootGroupId,
			Count: int64(lootGroupCount),
		})
	}

	arrPbItems := this.getPlayerBagFun().GetPbItems(arrItemInfo, pb.EmDoingType_EDT_Draw)

	// 获取词条效果
	if uDrawCount > 1 {

		vals := entry.KeyValueToMap(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_MultipleDraw), uDrawId)...)
		for key, val := range vals {
			//计算概率
			uAddEntryCount := 0
			for i := uint32(0); i < uDrawCount/10; i++ {
				if base.RandRange(0, base.MIL_PERCENT) <= val {
					uAddEntryCount++
				}
			}
			if uAddEntryCount <= 0 {
				continue
			}

			arrItemInfo = make([]*common.ItemInfo, 0)
			arrItemInfo = append(arrItemInfo, &common.ItemInfo{
				Kind:  uint32(cfgEnum.ESystemType_LootGroup),
				Id:    key,
				Count: int64(uAddEntryCount),
			})

			arrPbItems = append(arrPbItems, this.getPlayerBagFun().GetPbItems(arrItemInfo, pb.EmDoingType_EDT_Entry)...)
		}
	}

	uCode := this.getPlayerBagFun().AddPbItems(head, arrPbItems, pb.EmDoingType_EDT_Draw, true)
	if uCode != cfgEnum.ErrorCode_Success {
		plog.Print(this.AccountId, uCode, uDrawId, uDrawCount, arrItemInfo)
	}

	this.UpdateSave(true)

	if uAdvertType > uint32(cfgEnum.EAdvertType_None) {
		this.getPlayerSystemCommonFun().OnAdvert(head, uAdvertType)
	}

	//成就触发
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_FinishDraw, uDrawCount, uDrawId)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_FinishDraw, uDrawCount, 0)

	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_FinishDrawType, uDrawCount, cfgDraw.DrawType)

	//通知客户端
	cluster.SendToClient(head, &pb.DrawResponse{
		PacketHead: &pb.IPacket{},
		DrawInfo:   drawInfo,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 抽奖请求
func (this *PlayerSystemDrawFun) DrawPrizeInfoRequest(head *pb.RpcHead, pbRequest *pb.DrawPrizeInfoRequest) {
	uCode := this.DrawPrizeInfo(head, pbRequest.DrawId)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.DrawPrizeInfoResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 抽奖请求
func (this *PlayerSystemDrawFun) DrawPrizeInfo(head *pb.RpcHead, uDrawId uint32) cfgEnum.ErrorCode {
	cfgDraw := cfgData.GetCfgDraw(uDrawId)
	if cfgDraw == nil {
		return plog.Print(this.AccountId, cfgData.GetDrawErrorCode(uDrawId), uDrawId)
	}

	pbResponse := &pb.DrawPrizeInfoResponse{
		PacketHead: &pb.IPacket{},
		DrawId:     uDrawId,
	}

	for lootGroupID, rate := range cfgDraw.MapShowRandLootGroup {
		cfgLootGroup := cfgData.GetCfgLootGroup(lootGroupID)
		if cfgLootGroup == nil {
			return plog.Print(this.AccountId, cfgData.GetLootGroupErrorCode(lootGroupID), lootGroupID)
		}

		pbPrize := &pb.PBDrawPrizeInfo{
			Rate:  rate.Key,
			Rate2: rate.Value,
			Name:  cfgLootGroup.Name,
		}

		//需要增加抽奖up池子
		lootPool := make([]uint32, 0)
		lootPool = append(lootPool, cfgLootGroup.LootPool...)
		arrUpLoot := this.getPlayerSystemDrawFun().GetUpLoot(lootGroupID)
		if len(arrUpLoot) > 0 {
			for _, lootid := range arrUpLoot {
				if !base.ArrayContainsValue(lootPool, lootid) {
					lootPool = append(lootPool, lootid)
				}
			}
		}

		for _, lootid := range lootPool {
			cfgLoot := cfgData.GetCfgLoot(lootid)
			if cfgLoot == nil {
				return plog.Print(this.AccountId, cfgData.GetLootErrorCode(lootid), lootid)
			}

			pbPrize.ItemList = append(pbPrize.ItemList, &pb.PBAddItem{
				Id:     cfgLoot.LootId,
				Kind:   cfgLoot.LootKind,
				Count:  int64(cfgLoot.MaxCount),
				Params: cfgLoot.LootParam,
			})
		}

		pbResponse.PrizeList = append(pbResponse.PrizeList, pbPrize)
	}

	//通知客户端
	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 抽奖积分奖励请求
func (this *PlayerSystemDrawFun) DrawScorePrizeRequest(head *pb.RpcHead, pbRequest *pb.DrawScorePrizeRequest) {
	uCode := this.DrawScorePrize(head, pbRequest.DrawId, pbRequest.Id)
	cluster.SendToClient(head, &pb.DrawScorePrizeResponse{
		PacketHead: &pb.IPacket{},
		DrawId:     pbRequest.DrawId,
		Id:         pbRequest.Id,
	}, uCode)
}

// 抽奖积分奖励请求
func (this *PlayerSystemDrawFun) DrawScorePrize(head *pb.RpcHead, uDrawId, uId uint32) cfgEnum.ErrorCode {
	cfgScore := cfgData.GetCfgDrawScoreConfig(uId)
	if cfgScore == nil {
		return plog.Print(this.AccountId, cfgData.GetDrawScoreConfigErrorCode(uId), uId)
	}

	if cfgScore.DrawId != uDrawId {
		return plog.Print(this.AccountId, cfgData.GetDrawScoreConfigErrorCode(uId), uId)
	}

	drawInfo, ok := this.mapDrawData[uDrawId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, uId)
	}

	//已经领取
	if base.ArrayContainsValue(drawInfo.ScorePrize, uId) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HavePrize, uId)
	}

	//积分判断
	if drawInfo.DrawCount < cfgScore.Value {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedCondition, uId)
	}

	//给奖励 恭喜获得
	this.getPlayerBagFun().AddOneArrItem(head, cfgScore.AddPrize, pb.EmDoingType_EDT_StarSource, true)

	drawInfo.ScorePrize = append(drawInfo.ScorePrize, uId)
	this.UpdateSave(true)
	return cfgEnum.ErrorCode_Success
}
