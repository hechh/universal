package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	common2 "corps/common"
	"corps/common/cfgData"
	"corps/common/serverCommon"
	"corps/framework/common/uerror"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/module/entry"
	"sort"

	"github.com/golang/protobuf/proto"
)

// ----gomaker生成的模板-------
type PlayerSystemBoxFun struct {
	PlayerFun
	boxScore     uint32                   // 宝箱积分
	currentLevel uint32                   // 当前等级
	recycleTimes uint32                   // 宝箱积分循环次数
	mapBox       map[uint32]*pb.PBBoxInfo // 宝箱
}

// --------------------通用接口实现------------------------------
// 初始化
func (this *PlayerSystemBoxFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
}

// 新系统
func (this *PlayerSystemBoxFun) NewPlayer() {
	this.currentLevel = 1
	this.recycleTimes = 1
	this.mapBox = make(map[uint32]*pb.PBBoxInfo)

	// 保存
	this.UpdateSave(true)
}

// 加载系统数据(system类型数据)
func (this *PlayerSystemBoxFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	if pbSystem.Box == nil {
		this.NewPlayer()
		return
	}
	this.boxScore = pbSystem.Box.BoxScore
	this.currentLevel = pbSystem.Box.CurrentLevel
	this.recycleTimes = pbSystem.Box.RecycleTimes
	if this.recycleTimes <= 0 {
		this.recycleTimes = 1
	}
	this.mapBox = make(map[uint32]*pb.PBBoxInfo)
	for _, item := range pbSystem.Box.Boxs {
		this.mapBox[item.ItemID] = item
	}
}

// 存储数据 返回存储标志(system类型数据)
func (this *PlayerSystemBoxFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	data := &pb.PBPlayerSystemBox{
		BoxScore:     this.boxScore,
		CurrentLevel: this.currentLevel,
		RecycleTimes: this.recycleTimes,
	}
	for _, item := range this.mapBox {
		data.Boxs = append(data.Boxs, item)
	}
	sort.Slice(data.Boxs, func(i, j int) bool {
		return data.Boxs[i].ItemID < data.Boxs[j].ItemID
	})
	pbSystem.Box = data
	return true
}

// 客户端数据
func (this *PlayerSystemBoxFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = new(pb.PBPlayerSystem)
	}
	this.SaveSystem(pbData.System)
}
func (this *PlayerSystemBoxFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemBox{}
}

// 设置玩家数据, web管理后台
func (this *PlayerSystemBoxFun) SetUserTypeInfo(message proto.Message) bool {
	if message == nil {
		return false
	}
	box := message.(*pb.PBPlayerSystemBox)
	this.boxScore = box.BoxScore
	this.currentLevel = box.CurrentLevel
	this.recycleTimes = box.RecycleTimes
	this.mapBox = make(map[uint32]*pb.PBBoxInfo)
	for _, item := range box.Boxs {
		this.mapBox[item.ItemID] = item
	}
	return true
}

// 获取开宝箱次数
func (this *PlayerSystemBoxFun) GetBoxOpenCount(itemID uint32) uint32 {
	uCount := uint32(0)
	if itemID == 0 {
		for _, info := range this.mapBox {
			uCount += info.OpenTimes
		}
	} else {
		if info, ok := this.mapBox[itemID]; ok {
			uCount = info.OpenTimes
		}
	}
	return uCount
}

func (this *PlayerSystemBoxFun) getBox(itemID uint32) *pb.PBBoxInfo {
	if this.mapBox[itemID] == nil {
		this.mapBox[itemID] = &pb.PBBoxInfo{ItemID: itemID}
	}
	return this.mapBox[itemID]
}

// --------------------交互接口实现------------------------------
func (this *PlayerSystemBoxFun) BoxOpenRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.BoxOpenRequest)
	rsp := response.(*pb.BoxOpenResponse)
	// 读取配置
	itemCfg := cfgData.GetCfgItem(req.ItemID)
	if itemCfg == nil {
		return uerror.NewUErrorf(1, cfgData.GetItemErrorCode(req.ItemID), "head: %v, req: %v", head, req)
	}
	// 发放数量限制
	if itemCfg.MaxUse > 0 && req.ItemNum > itemCfg.MaxUse {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_UpLimit, "head: %v, req: %v, limit: %d", head, req, itemCfg.MaxUse)
	}
	if req.AdvestType == uint32(cfgEnum.EAdvertType_BoxOpen) {
		cfgAdvert := cfgData.GetCfgAdvertConfig(req.AdvestType)
		if cfgAdvert == nil {
			return uerror.NewUErrorf(1, cfgData.GetAdvertConfigErrorCode(req.AdvestType), "head: %v, req: %v", head, req)
		}
		uCode := this.getPlayerSystemCommonFun().AddAdvert(head, req.AdvestType)
		if uCode != cfgEnum.ErrorCode_Success {
			return uerror.NewUErrorf(1, uCode, "head: %v, req: %v, AdvestType: %d", head, req, req.AdvestType)
		}

		req.ItemNum = cfgAdvert.Param
	} else if req.AdvestType == uint32(cfgEnum.EAdvertType_None) {
		// 判断宝箱是否有足够数量
		curNum := this.getPlayerBagFun().GetItemCount(uint32(cfgEnum.ESystemType_Item), req.ItemID)
		if req.ItemNum > uint32(curNum) {
			return uerror.NewUErrorf(1, cfgEnum.ErrorCode_ItemNotEnough, "head: %v, req: %v, number: %d", head, req, curNum)
		}

	} else {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_AdvertNotEqual, "head: %v, req: %v", head, req)
	}

	// 词条类型 额外掉落loot
	extras := entry.KeyValueToMap(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_OpenBoxExtraLootDropProb), req.ItemID)...)
	extraLoots := []*common2.ItemInfo{}

	// 循环开宝箱
	listLootId := make([]uint32, 0)
	for i := 0; i < int(req.ItemNum); i++ {
		// 获取宝箱信息
		boxInfo := this.getBox(req.ItemID)
		// 加载配置
		boxCfg := cfgData.GetCfgBox(boxInfo.OpenTimes, req.ItemID)
		if boxCfg == nil {
			return uerror.NewUErrorf(1, cfgData.GetBoxErrorCode(req.ItemID), "head: %v, req: %v, boxInfo: %v", head, req, boxInfo)
		}
		// loot掉落次数统计
		listLootId = append(listLootId, boxCfg.DropGroupID)
		// 累加开宝箱次数
		boxInfo.OpenTimes++
		// 计算开宝箱获得的积分
		this.boxScore += uint32(boxCfg.Score)
		// 额外掉落loot
		for lootGroupId, prob := range extras {
			if base.IsRadio(prob) {
				extraLoots = append(extraLoots, &common2.ItemInfo{Kind: uint32(cfgEnum.ESystemType_LootGroup), Id: lootGroupId, Count: 1})
			}
		}
	}

	// 扣除宝箱数量 delItem
	if req.AdvestType == uint32(cfgEnum.EAdvertType_None) {
		errcode := this.getPlayerBagFun().DelItem(head, uint32(cfgEnum.ESystemType_Item), req.ItemID, int64(req.ItemNum), pb.EmDoingType_EDT_BoxConsume)
		if errcode != cfgEnum.ErrorCode_Success {
			return uerror.NewUErrorf(1, errcode, "head: %v, req: %v, boxCfg: %v, boxInfo: %v", head, req)
		}
	}

	//分别开宝箱 如果开一百个，每一次奖励合并，如果少于100, 不合并奖励
	//词条开一次
	dmapEntryIdCount := make(map[uint32]map[uint32]int64)
	dmapEntryAdd := entry.KeyValueToDMap(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_OpenBox), req.ItemID)...)

	// 走金币膨胀系数
	mapTech := this.getPlayerSystemHookTechFun().GetHookTechEffect(cfgEnum.TechEffectType_AddBoxGoldRate)
	uBoxGoldGrowthRate := uint32(0)
	if len(mapTech) > 0 {
		uBoxGoldGrowthRate = uint32(mapTech[0])
	}

	arrPbAllItems := make([]*pb.PBAddItemData, 0)
	for _, lootId := range listLootId {
		// 开宝箱
		arrItem := &common2.ItemInfo{
			Kind:  uint32(cfgEnum.ESystemType_LootGroup),
			Id:    lootId,
			Count: 1,
		}
		// 随机奖励(不入库)
		arrPbItem := this.getPlayerBagFun().GetPbItems([]*common2.ItemInfo{arrItem}, pb.EmDoingType_EDT_BoxOpen)
		if len(arrPbItem) <= 0 {
			continue
		}
		//金币膨胀系数
		if req.ItemID == 1201 && uBoxGoldGrowthRate > 0 {
			for _, item := range arrPbItem {
				if item.Kind != uint32(cfgEnum.ESystemType_Item) || item.Id != uint32(pb.EmItemExpendType_EIET_Gold) {
					continue
				}
				item.Count += int64(base.CeilU32(uint32(item.Count*int64(uBoxGoldGrowthRate)), base.MIL_PERCENT))
			}
		}
		//词条需要算合并后的加成
		arrTmpMergeItem := serverCommon.Merge_PBItem(arrPbItem)
		//10次需要合并奖励
		if req.ItemNum > 10 {
			arrPbAllItems = append(arrPbAllItems, arrTmpMergeItem...)
		} else {
			arrPbAllItems = append(arrPbAllItems, arrPbItem...)
		}
		// 词条类型 道具增加
		if len(dmapEntryAdd) == 0 {
			continue
		}
		for _, item := range arrTmpMergeItem {
			if _, ok := dmapEntryAdd[item.Kind]; !ok {
				continue
			}
			if _, ok := dmapEntryAdd[item.Kind][item.Id]; !ok {
				continue
			}
			if _, ok := dmapEntryIdCount[item.Kind]; !ok {
				dmapEntryIdCount[item.Kind] = make(map[uint32]int64)
			}
			dmapEntryIdCount[item.Kind][item.Id] += int64(base.CeilU32(uint32(item.Count*int64(dmapEntryAdd[item.Kind][item.Id])), base.MIL_PERCENT))
		}
	}

	// 判断 广告奖励词条是否生效 + 判断是否获得额外多一份
	if req.AdvestType == uint32(cfgEnum.EAdvertType_BoxOpen) {
		propValue := entry.ToValue(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_AdertiseReward), uint32(cfgEnum.EAdvertType_BoxOpen))...)
		if base.IsRadio(propValue) {
			items := serverCommon.MergePBItemWithDoingType(arrPbAllItems, pb.EmDoingType_EDT_Entry)
			plog.Trace("EmDoingType_EDT_Entry advert rewards: %v", items)
			arrPbAllItems = append(arrPbAllItems, items...)
		}
	}

	// 词条类型 道具增加
	if len(dmapEntryIdCount) > 0 {
		arrEntryItem := make([]*common2.ItemInfo, 0)
		for tkind, mapTCount := range dmapEntryIdCount {
			for tid, tcount := range mapTCount {
				arrEntryItem = append(arrEntryItem, &common2.ItemInfo{
					Kind:  tkind,
					Id:    tid,
					Count: tcount,
				})
			}
		}
		items := this.getPlayerBagFun().GetPbItems(arrEntryItem, pb.EmDoingType_EDT_Entry)
		plog.Trace("EmDoingType_EDT_Entry openbox rewards: %v", items)
		arrPbAllItems = append(arrPbAllItems, items...)
	}

	// 词条类型 额外掉落loot
	if len(extraLoots) > 0 {
		items := this.getPlayerBagFun().GetPbItems(extraLoots, pb.EmDoingType_EDT_Entry)
		if len(items) > 0 {
			arrPbAllItems = append(arrPbAllItems, items...)
		}
	}

	plog.Trace("OpenBox: %v, head: %v, req: %v", arrPbAllItems, head, req)

	// 道具入库
	this.getPlayerBagFun().AddPbItems(head, arrPbAllItems, pb.EmDoingType_EDT_BoxOpen, false)
	//成就触发
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_OpenBox, req.ItemNum, itemCfg.Id)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_OpenBox, req.ItemNum, 0)
	// 组装返回数据
	rsp.Score = this.boxScore
	rsp.ItemInfo = arrPbAllItems
	this.UpdateSave(true)
	return nil
}

func (this *PlayerSystemBoxFun) BoxProgressRewardRequest(head *pb.RpcHead, request, response proto.Message) error {
	//req := request.(*pb.BoxProgressRewardRequest)
	rsp := response.(*pb.BoxProgressRewardResponse)
	// 扣除积分，获取宝箱
	boxAdds := map[uint32]int64{}
	boxScore := this.boxScore
	currentLevel := this.currentLevel
	recycleTimes := this.recycleTimes
	isOnce := cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_BOX_SCORE_REWARD_LIMIT) > this.boxScore
	for {
		cfgData.WalkCfgBoxScoreProgress(recycleTimes, func(cfg *cfgData.BoxScoreProgressCfg) bool {
			if currentLevel == 0 || currentLevel == cfg.Stage {
				// 只领取一次宝箱
				if isOnce && currentLevel == 0 {
					currentLevel = cfg.Stage
					return false
				}
				// 领取
				if cfg.Score <= boxScore {
					boxScore -= cfg.Score
					boxAdds[cfg.BoxID]++
					currentLevel = 0
				} else {
					currentLevel = cfg.Stage
					return false
				}
			}
			return true
		})
		if isOnce {
			if currentLevel == 0 {
				recycleTimes++
				currentLevel = cfgData.GetCfgCycleMinStage(recycleTimes)
			}
			break
		} else {
			if currentLevel == 0 {
				recycleTimes++
				currentLevel = cfgData.GetCfgCycleMinStage(recycleTimes)
			} else {
				break
			}
		}
	}
	if len(boxAdds) <= 0 {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_BoxScoreNumber, "head: %v, score: %d, level: %d", head, boxScore, currentLevel)
	}
	// 加载配置
	cfg := cfgData.GetCfgBoxScoreProgress(recycleTimes, currentLevel)
	if cfg == nil {
		return uerror.NewUErrorf(1, cfgData.GetBoxScoreProgressErrorCode(recycleTimes), "head: %v, recycleTimes: %d, level: %d", head, recycleTimes, currentLevel)
	}
	// 增加宝箱数量
	arrItem := make([]*common2.ItemInfo, 0)
	for boxId, count := range boxAdds {
		arrItem = append(arrItem, &common2.ItemInfo{Id: boxId, Count: count, Kind: uint32(cfgEnum.ESystemType_Item)})
	}
	if errCode := this.getPlayerBagFun().AddArrItem(head, arrItem, pb.EmDoingType_EDT_BoxScoreExchange, true); errCode != cfgEnum.ErrorCode_Success {
		return uerror.NewUErrorf(1, errCode, "head: %v, recycleTimes: %d, level: %d", head, recycleTimes, currentLevel)
	}
	// 设置值
	this.boxScore = boxScore
	this.currentLevel = currentLevel
	this.recycleTimes = recycleTimes
	// 返回值
	rsp.Level = this.currentLevel
	rsp.Score = this.boxScore
	rsp.NeedScore = cfg.Score
	rsp.Recycle = this.recycleTimes
	// 保存数据
	this.UpdateSave(true)
	return nil
}
