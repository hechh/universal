package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/framework/cluster"
	"corps/framework/common/uerror"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/module/entry"
	"sort"

	"github.com/golang/protobuf/proto"
)

const (
	FIVE_PROPERTY_COUNT = 3
)

// ----gomaker生成的模板-------
type PlayerCrystalFun struct {
	PlayerFun
	book             *pb.PBCrystalBook             // 图鉴系统
	mapCrystal       map[uint32]*pb.PBCrystal      // 晶核系统
	mapRobot         map[uint32]*pb.PBCrystalRobot // 机器人系统
	mapCrystalRobot  map[uint32]uint32             // 晶核ID--机器人ID
	mapRobotChange   map[uint32]struct{}           // 机器人变更
	mapCrystalChange map[uint32]struct{}           // 晶核变更
	entryData        *entry.EntryService           // 词条系统
}

// --------------------通用接口实现------------------------------

// 初始化
func (this *PlayerCrystalFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
}

// 新系统
func (this *PlayerCrystalFun) NewPlayer() {
	this.book = &pb.PBCrystalBook{Stage: 1}
	this.mapCrystal = make(map[uint32]*pb.PBCrystal)
	this.mapRobot = make(map[uint32]*pb.PBCrystalRobot)
	this.mapCrystalRobot = make(map[uint32]uint32)
	this.mapRobotChange = make(map[uint32]struct{})
	this.mapCrystalChange = make(map[uint32]struct{})
	// 初始化词条系统
	this.entryData = entry.NewEntryService(this.AccountId, &pb.PBPlayerCrystal{}, this.getPlayerSystemTaskFun().AchieveBase, this.onEntryEffectChange)
	this.entryData.SetCrystal(this.mapCrystal)
	// 解锁默认关卡
	mapId, stageId := this.getPlayerSystemBattleFun().GetMapIdAndStageId(pb.EmBattleType_EBT_Normal)
	for _, robotCfg := range cfgData.GetCfgCrystalRobotByCondition(mapId, stageId) {
		this.unlockRobot(&pb.RpcHead{Id: this.AccountId}, robotCfg, false)
	}
	this.UpdateSave(true)
}

// 加载数据(非system类型数据)
func (this *PlayerCrystalFun) Load(pData []byte) {
	if len(pData) <= 0 {
		this.NewPlayer()
		return
	}
	// 序列化
	pbData := &pb.PBPlayerCrystal{}
	proto.Unmarshal(pData, pbData)
	this.initData(pbData)

	// 保存数据
	this.UpdateSave(true)
}

// 存储数据(非system类型数据)
func (this *PlayerCrystalFun) Save(bNow bool) {
	if !this.BSave {
		return
	}
	this.BSave = false
	//通知db保存玩家数据
	buff, _ := proto.Marshal(this.toProto(false))
	cluster.SendToDb(&pb.RpcHead{Id: this.AccountId}, "DbPlayerMgr", "SavePlayerDB", this.PbType, buff, bNow)
}

// 客户端数据
func (this *PlayerCrystalFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	pbData.Crystal = this.toProto(true)
}

// flag: 是否计算效果值
func (this *PlayerCrystalFun) toProto(flag bool) *pb.PBPlayerCrystal {
	pbData := &pb.PBPlayerCrystal{Book: this.book}
	// 机器人系统
	for _, robot := range this.mapRobot {
		pbData.Robots = append(pbData.Robots, robot)
	}
	sort.Slice(pbData.Robots, func(i, j int) bool {
		return pbData.Robots[i].RobotID < pbData.Robots[j].RobotID
	})
	// 晶核系统
	for _, crystal := range this.mapCrystal {
		pbData.Crystals = append(pbData.Crystals, crystal)
	}
	sort.Slice(pbData.Crystals, func(i, j int) bool {
		return pbData.Crystals[i].CrystalID < pbData.Crystals[j].CrystalID
	})
	// 词条系统
	this.entryData.ToProto(flag, pbData)
	return pbData
}
func (this *PlayerCrystalFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerCrystal{}
}

// 设置玩家数据, web管理后台
func (this *PlayerCrystalFun) SetUserTypeInfo(message proto.Message) bool {
	if message == nil {
		return false
	}
	// 设置缓存值
	pbSystem := message.(*pb.PBPlayerCrystal)
	if pbSystem == nil {
		return false
	}
	this.initData(pbSystem)
	this.UpdateSave(true)
	return true
}

// 加载数据
func (this *PlayerCrystalFun) initData(data *pb.PBPlayerCrystal) {
	this.book = &pb.PBCrystalBook{Stage: 1}
	this.mapRobot = make(map[uint32]*pb.PBCrystalRobot)
	this.mapCrystal = make(map[uint32]*pb.PBCrystal)
	this.mapCrystalRobot = make(map[uint32]uint32)
	this.mapRobotChange = make(map[uint32]struct{})
	this.mapCrystalChange = make(map[uint32]struct{})
	// 初始化图鉴系统
	if data.Book != nil {
		this.book = data.Book
	}
	// 初始化晶核系统
	for _, item := range data.Crystals {
		this.mapCrystal[item.CrystalID] = item
	}
	// 初始化机器人
	for _, item := range data.Robots {
		this.mapRobot[item.RobotID] = item
		for _, crystalID := range item.Crystals {
			this.mapCrystalRobot[crystalID] = item.RobotID
		}
	}
	// 初始化词条系统
	this.entryData = entry.NewEntryService(this.AccountId, data, this.getPlayerSystemTaskFun().AchieveBase, this.onEntryEffectChange)
	this.entryData.SetCrystal(this.mapCrystal)
}

func (this *PlayerCrystalFun) onEntryEffectChange(uEffectType uint32) {
	//重新计算属性
	if uEffectType == uint32(cfgEnum.EntryEffectType_HeroProp) {
		this.getPlayerHeroFun().updateCalcFightpower(true)
	} else if uEffectType == uint32(cfgEnum.EntryEffectType_EquipmentBagSize) {
		this.GetPlayerEquipmentFun().updateMaxPosCount()
	} else if uEffectType == uint32(cfgEnum.EntryEffectType_AddShopRefresh) {
		this.getPlayerSystemCommonFun().UpdatePrivilege(cfgEnum.PrivilegeType_BlackShopRefreshCount)
	} else if uEffectType == uint32(cfgEnum.EntryEffectType_SplitEquipmentRewardBox) {
		this.GetPlayerEquipmentFun().ClearSplitEquipCount()
	} else if uEffectType == uint32(cfgEnum.EntryEffectType_OfflineIncomeTime) {
		this.getPlayerSystemCommonFun().UpdatePrivilege(cfgEnum.PrivilegeType_OfflineTime)
	}
}

func (this *PlayerCrystalFun) GetCrystalProp(element uint32) map[uint32]float64 {
	rets := make(map[uint32]float64)
	for id, crystal := range this.mapCrystal {
		if crystal.Element != element {
			continue
		}
		cfg := cfgData.GetCfgCrystal(id)
		if cfg == nil {
			continue
		}
		tmps := map[uint32]float64{}
		for _, val := range cfg.TalentProps {
			tmps[val.Key] += float64(val.Value)
		}
		// 升级属性增加
		for _, val := range cfg.AddTalentProps {
			tmps[val.Key] += float64(val.Value * crystal.Level)
		}
		// 突破属性增加
		percent := uint32(0)
		for i := uint32(1); i <= crystal.Star; i++ {
			if qualityCfg := cfgData.GetCfgCrystalQuality(crystal.Quality, i); qualityCfg != nil {
				percent += qualityCfg.BasePercent
			}
		}
		for key, val := range tmps {
			tmps[key] += val * float64(percent) / base.MIL_PERCENT
		}
		// 晶核属性提升
		if extra, ok := entry.KeyValueToDMap(this.entryData.Get(uint32(cfgEnum.EntryEffectType_CrystalProp), element)...)[id]; ok {
			for k, val := range tmps {
				tmps[k] += val * float64(extra[k]) / base.MIL_PERCENT
			}
		}
		for k, val := range tmps {
			rets[k] += val
		}
	}
	return rets
}

// 增加晶核 注意pbItem内数据会被修改，为了方便客户端展示恭喜获得
func (this *PlayerCrystalFun) AddCrystal(head *pb.RpcHead, pbItem *pb.PBAddItemData) cfgEnum.ErrorCode {
	crystalID := pbItem.Id
	itemCount := pbItem.Count
	star := uint32(1)
	if len(pbItem.Params) > 0 {
		star = pbItem.Params[0]
	}
	// 读取配置
	cfg := cfgData.GetCfgCrystal(crystalID)
	if cfg == nil {
		plog.Error("head: %v, pbItem: %v", head, pbItem)
		return cfgData.GetCrystalErrorCode(crystalID)
	}
	// 判断晶核是否解锁
	crystal, ok := this.mapCrystal[crystalID]
	if !ok {
		// 解锁晶核
		if err := this.unlockCrystal(head, cfg, star); err != nil {
			plog.Error("%v", err)
			return uerror.GetCode(err)
		}
		itemCount--
		crystal = this.mapCrystal[crystalID]
	}
	// 转换成其他道具
	if itemCount > 0 {
		// 晶核转成晶核碎片
		qualityCfg := cfgData.GetCfgCrystalQuality(crystal.Quality, star)
		count := itemCount * int64(qualityCfg.CrystalDebris)
		errCode := this.getPlayerBagFun().AddItem(head, uint32(cfgEnum.ESystemType_Item), cfg.PieceID, count, pbItem.DoingType, false)
		if errCode != cfgEnum.ErrorCode_Success {
			plog.Error("head: %v, pbItem: %v, code: %d", head, pbItem, errCode)
		}
		pbItem.Id = cfg.PieceID
		pbItem.Kind = uint32(cfgEnum.ESystemType_Item)
		pbItem.Count = count
	}
	//成就触发
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ReceiveCrystal, uint32(pbItem.Count), crystal.Quality)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ReceiveCrystal, uint32(pbItem.Count), uint32(cfgEnum.EQuality_Any))
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ReceiveElementCrystal, uint32(pbItem.Count), crystal.Element, crystal.Quality)
	// 更新数据库
	this.UpdateSave(true)
	this.notifyCrystal(head)
	return cfgEnum.ErrorCode_Success
}

// 解锁晶核
func (this *PlayerCrystalFun) unlockCrystal(head *pb.RpcHead, cfg *cfgData.CrystalCfg, star uint32) error {
	// 解锁被动技能
	passives := []uint32{cfg.MainSkillID}
	if cfg.SpecialSkillID > 0 {
		passives = append(passives, cfg.SpecialSkillID)
	}
	cfgData.WalkCfgCrystalRedefineProp(func(pcfg *cfgData.CrystalRedefinePropCfg) bool {
		// 结束循环
		if pcfg.StageValue > star {
			return false
		}
		if index := pcfg.Id - 1; int(index) < len(cfg.PassiveSkillId) {
			passives = append(passives, cfg.PassiveSkillId[index])
		}
		return true
	})
	// 初始化
	this.mapCrystal[cfg.Id] = &pb.PBCrystal{
		CrystalID:       cfg.Id,
		Element:         cfg.Element,
		Quality:         cfg.Quality,
		Star:            star,
		RewardCoinTimes: 1,
		PassiveSkillIds: passives,
	}
	// 解锁词条
	for _, skillID := range passives {
		this.entryData.Unlock(head, skillID)
	}
	//成就统计
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_CrystalTotalLevel, this.GetTotalStar())
	this.mapCrystalChange[cfg.Id] = struct{}{}
	return nil
}

// 添加机器人
func (this *PlayerCrystalFun) AddRobot(head *pb.RpcHead, pbItem *pb.PBAddItemData) cfgEnum.ErrorCode {
	robotID := pbItem.Id
	itemCount := pbItem.Count
	// 读取配置
	cfg := cfgData.GetCfgCrystalRobot(robotID)
	if cfg == nil {
		plog.Error("head: %v, pbItem: %v", head, pbItem)
		return cfgData.GetCrystalRobotErrorCode(robotID)
	}

	//不解锁，只给碎片
	/*if _, ok := this.mapRobot[robotID]; !ok {
		this.unlockRobot(head, cfg, true)
		this.mapRobotChange[robotID] = struct{}{}
		itemCount--
	}*/
	// 转换成碎片
	count := itemCount * int64(cfg.Convert)
	errCode := this.getPlayerBagFun().AddItem(head, uint32(cfgEnum.ESystemType_Item), cfg.RobotDebris, count, pbItem.DoingType, false)
	if errCode != cfgEnum.ErrorCode_Success {
		plog.Error("head: %v, cfg: %v, count: %d, code: %d", head, cfg, count, errCode)
	}

	pbItem.Kind = uint32(cfgEnum.ESystemType_Item)
	pbItem.Id = cfg.RobotDebris
	pbItem.Count = count

	this.notifyRobot(head)
	this.UpdateSave(true)
	return cfgEnum.ErrorCode_Success
}

// 精英关卡解锁机器人
func (this *PlayerCrystalFun) UnlockRobot(head *pb.RpcHead, mapId, stageId uint32, isNotify bool) {
	for _, robotCfg := range cfgData.GetCfgCrystalRobotByCondition(mapId, stageId) {
		// 对于没有解锁的机器人，重新解锁一遍
		if _, ok := this.mapRobot[robotCfg.Id]; !ok {
			this.unlockRobot(head, robotCfg, isNotify)
		}
	}
	if isNotify {
		this.notifyRobot(head)
	}
	this.UpdateSave(true)
}

// 解锁机器人
func (this *PlayerCrystalFun) unlockRobot(head *pb.RpcHead, robotCfg *cfgData.CrystalRobotCfg, isNotify bool) {
	// 解锁共鸣技能
	linkages := []uint32{}
	cfgData.WalkCfgCrystalRobotProp(robotCfg.Id, func(cfg *cfgData.CrystalRobotPropCfg) bool {
		if cfg.Stage <= 0 {
			linkages = append(linkages, cfg.Id)
			return true
		}
		return false
	})
	// 关联晶核
	crystalIDs := base.CopyArrUint32(robotCfg.Crystal)
	// 初始化机器人
	this.mapRobot[robotCfg.Id] = &pb.PBCrystalRobot{
		RobotID:        robotCfg.Id,          // 机器人id
		Stage:          1,                    // 等级
		RoleSkillID:    robotCfg.RoleSkillId, // 技能ID
		UnlockLinkages: linkages,             // 解锁的共鸣技能
		Crystals:       crystalIDs,           // 装备的晶核
	}
	// 设置关系
	for _, crystalID := range crystalIDs {
		this.mapCrystalRobot[crystalID] = robotCfg.Id
	}
	// 解锁机器人词条
	robot := this.mapRobot[robotCfg.Id]
	filter := map[uint32]struct{}{}
	for _, ii := range robot.UnlockLinkages {
		filter[ii] = struct{}{}
	}
	cfgData.WalkCfgCrystalRobotProp(robot.RobotID, func(cfg *cfgData.CrystalRobotPropCfg) bool {
		if robot.Stage >= cfg.Stage {
			robot.UnlockLinkages = append(robot.UnlockLinkages, cfg.Id)
			return true
		}
		return false
	})

	if isNotify {
		this.mapRobotChange[robotCfg.Id] = struct{}{}
	}
	this.UpdateSave(true)
}

func (this *PlayerCrystalFun) GetTotalStar() uint32 {
	uTotal := uint32(0)
	for _, info := range this.mapCrystal {
		uTotal += info.Star
	}
	return uTotal
}

func (this *PlayerCrystalFun) notifyCrystal(head *pb.RpcHead) {
	if len(this.mapCrystalChange) <= 0 {
		return
	}
	notify := &pb.CrystalNotify{PacketHead: &pb.IPacket{}}
	for id := range this.mapCrystalChange {
		notify.CrystalInfo = append(notify.CrystalInfo, this.mapCrystal[id])
		// 删除记录
		delete(this.mapCrystalChange, id)
	}
	cluster.SendToClient(head, notify, cfgEnum.ErrorCode_Success)
	plog.Debug("Crystal: %v", notify)
}

func (this *PlayerCrystalFun) notifyRobot(head *pb.RpcHead) {
	for robotID := range this.mapRobotChange {
		cluster.SendToClient(head, &pb.CrystalRobotNotify{PacketHead: &pb.IPacket{}, RobotInfo: this.mapRobot[robotID]}, cfgEnum.ErrorCode_Success)

		plog.Debug("Robot: %v", this.mapRobotChange[robotID])

		// 删除记录
		delete(this.mapRobotChange, robotID)
	}
}

// --------------------交互接口实现------------------------------
// 词条解锁
func (this *PlayerCrystalFun) EntryUnlockRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.EntryUnlockRequest)
	// 解锁
	this.entryData.Unlock(head, req.PassiveSkillID)
	this.UpdateSave(true)
	return nil
}

// 词条触发
func (this *PlayerCrystalFun) EntryTriggerRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.EntryTriggerRequest)
	// 触发词条
	this.entryData.Trigger(head, req.EntryType, req.Times, req.Params...)
	this.UpdateSave(true)
	return nil
}

// 机器人升级请求
func (this *PlayerCrystalFun) CrystalRobotUpgradeRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.CrystalRobotUpgradeRequest)
	rsp := response.(*pb.CrystalRobotUpgradeResponse)
	// 判断机器人是否存在
	robot, ok := this.mapRobot[req.RobotID]
	if !ok {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_CrystalRobotNotExist, "head: %v, req: %v", head, req)
	}
	// 极限情况:已经到达最大等级值
	if cfgData.GetCfgCrystalRobotPropMaxStage() == robot.FinishStage {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_AllStageFinished, "head: %v, req: %v, robot: %v", head, req, robot)
	}
	robotCfg := cfgData.GetCfgCrystalRobot(req.RobotID)
	if robotCfg == nil {
		return uerror.NewUErrorf(1, cfgData.GetCrystalRobotErrorCode(req.RobotID), "head: %v, req: %v, robot: %v", head, req, robot)
	}
	// 判断下一个等级配置是否存在
	gradeCfg := cfgData.GetCfgCrystalRobotGrade(robot.Stage)
	if gradeCfg == nil {
		return uerror.NewUErrorf(1, cfgData.GetCrystalRobotGradeErrorCode(robot.Stage), "head: %v, req: %v, robot: %v", head, req, robot)
	}
	if robot.Stage == robot.FinishStage {
		if gradeCfg = cfgData.GetCfgNextCrystalRobotGrade(robot.Stage); gradeCfg == nil {
			return uerror.NewUErrorf(1, cfgEnum.ErrorCode_MaxLevel, "head: %v, req: %v, robot: %v", head, req, robot)
		}
		robot.Stage = gradeCfg.Id
	}
	// 扣除升级道具
	delItems := []*common.ItemInfo{gradeCfg.ConsumeItem, {Kind: uint32(cfgEnum.ESystemType_Item), Id: robotCfg.RobotDebris, Count: int64(gradeCfg.Debris)}}
	errCode := this.getPlayerBagFun().DelArrItem(head, delItems, pb.EmDoingType_EDT_CrystalRobotUpgrade)
	if errCode != cfgEnum.ErrorCode_Success {
		return uerror.NewUErrorf(1, errCode, "head: %v, req: %v, robot: %v, delItems: %v", head, req, robot, delItems)
	}
	// 提升科技值等级
	this.mapRobotChange[req.RobotID] = struct{}{} // 锁定变更通知
	robot.FinishStage = gradeCfg.Id               // 设置完成等级
	// 提升等级
	if nextCfg := cfgData.GetCfgNextCrystalRobotGrade(robot.Stage); nextCfg != nil {
		robot.Stage = nextCfg.Id
	}
	robot.RoleSkillPercent = gradeCfg.Add // 提升机器人共鸣技能强化百分比
	// 解锁机器人词条
	filter := map[uint32]struct{}{}
	for _, ii := range robot.UnlockLinkages {
		filter[ii] = struct{}{}
	}
	cfgData.WalkCfgCrystalRobotProp(req.RobotID, func(cfg *cfgData.CrystalRobotPropCfg) bool {
		if robot.Stage >= cfg.Stage {
			robot.UnlockLinkages = append(robot.UnlockLinkages, cfg.Id)
			return true
		}
		return false
	})

	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_CrystalRobotUpgrade, 1, req.RobotID)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_CrystalRobotUpgrade, 1, 0)

	// 组装返回数据
	rsp.CurLevel = robot.Stage
	rsp.RobotID = req.RobotID
	// 更新数据库
	this.notifyRobot(head)
	this.UpdateSave(true)
	return nil
}

// 图鉴系统领取收藏币请求
func (this *PlayerCrystalFun) BookCollectionCoinRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.BookCollectionCoinRequest)
	rsp := response.(*pb.BookCollectionCoinResponse)
	// 判断晶核是否存在
	crystal, ok := this.mapCrystal[req.CrystalID]
	if !ok {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_CrystalNotExist, "head: %v, req: %v", head, req)
	}
	// 判单是否有收藏币领取
	if crystal.RewardCoinTimes <= 0 {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_HavePrize, "head: %v, req: %v, crystal: %v", head, req, crystal)
	}
	// 一次性领取所有收藏币
	qualityCfg := cfgData.GetCfgCrystalQuality(crystal.Quality, crystal.Star)
	if qualityCfg == nil {
		return uerror.NewUErrorf(1, cfgData.GetCrystalQualityErrorCode(crystal.Quality), "head: %v, req: %v", head, req)
	}

	// 计算增加的收藏币
	total := (qualityCfg.CollectCoin * crystal.RewardCoinTimes)
	crystal.RewardCoinTimes = 0
	this.mapCrystalChange[req.CrystalID] = struct{}{}
	this.book.Coin += total

	//成就统计
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_CrystalCoinCount, total)

	// 返回客户端
	rsp.Coin = this.book.Coin
	rsp.Level = this.book.Stage

	// 更新数据库
	this.notifyCrystal(head)
	this.UpdateSave(true)
	return nil
}

// 图鉴等级升级请求
func (this *PlayerCrystalFun) BookStageRewardRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.BookStageRewardRequest)
	rsp := response.(*pb.BookStageRewardResponse)
	// 判单是否达到最大等级限制
	if maxStage := cfgData.GetCfgCrystalBookGradeMaxStage(); maxStage <= this.book.FinishedStage {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_AllStageFinished, "head: %v, req: %v", head, req)
	}
	// 判断配置是否存在
	cfg := cfgData.GetCfgCrystalBookGrade(this.book.Stage)
	if cfg == nil {
		return uerror.NewUErrorf(1, cfgData.GetCrystalBookGradeErrorCode(this.book.Stage), "head: %v, req: %v, book: %v", head, req, this.book)
	}
	if cfg.Id == this.book.FinishedStage {
		cfg = cfgData.GetCfgNextCrystalBookGrade(this.book.FinishedStage)
		if cfg == nil {
			return uerror.NewUErrorf(1, cfgData.GetCrystalBookGradeErrorCode(this.book.FinishedStage), "head: %v, req: %v, book: %v", head, req, this.book)
		}
		this.book.Stage = cfg.Id
	}
	// 判断收藏币是否足够升级
	if this.book.Coin < cfg.StageValue {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_NotExceedStageLimit, "head: %v, req: %v, cfg: %v", head, req, cfg)
	}
	// 提升等级
	this.book.Coin -= cfg.StageValue
	this.book.FinishedStage = this.book.Stage
	if cfg := cfgData.GetCfgNextCrystalBookGrade(this.book.FinishedStage); cfg != nil {
		this.book.Stage = cfg.Id
	}
	// 发放奖励
	items := []*common.ItemInfo{}
	for _, item := range cfg.RewardItem {
		newItem := *item
		items = append(items, &newItem)
	}
	this.getPlayerBagFun().AddArrItem(head, items, pb.EmDoingType_EDT_CrystalBookUpgrade, true)
	// 组转返回数据
	rsp.BookInfo = this.book
	// 更新数据库
	this.UpdateSave(true)
	return nil
}

// 晶核改造
func (this *PlayerCrystalFun) CrystalRedefineRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.CrystalRedefineRequest)
	rsp := response.(*pb.CrystalRedefineResponse)
	// 判断晶核是否存在
	crystal, ok := this.mapCrystal[req.CrystalID]
	if !ok {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_CrystalRobotNotExist, "head: %v, req: %v", head, req)
	}
	// 加载配置表
	crystalCfg := cfgData.GetCfgCrystal(req.CrystalID)
	qualityCfg := cfgData.GetCfgCrystalQuality(crystal.Quality, crystal.Star)
	if qualityCfg == nil {
		return uerror.NewUErrorf(1, cfgData.GetCrystalQualityErrorCode(crystal.Quality), "head: %v, req: %v, crystal: %v", head, req, crystal)
	}
	// 扣除晶核碎片
	delItems := []*common.ItemInfo{{Kind: uint32(cfgEnum.ESystemType_Item), Id: crystalCfg.PieceID, Count: int64(qualityCfg.RedefineConsume)}}
	for _, item := range qualityCfg.DelArrItems {
		delItems = append(delItems, &common.ItemInfo{
			Id:     item.Id,
			Kind:   item.Kind,
			Count:  item.Count,
			Params: base.CopyArrUint32(item.Params),
		})
	}
	if errCode := this.getPlayerBagFun().DelArrItem(head, delItems, pb.EmDoingType_EDT_CrystalRedefine); errCode != cfgEnum.ErrorCode_Success {
		return uerror.NewUErrorf(1, errCode, "head: %v, req: %v, crystal: %v", head, req, crystal)
	}
	// 增加星星数量 + 天赋属性百分比提升
	crystal.Star++
	crystal.RewardCoinTimes++
	this.mapCrystalChange[req.CrystalID] = struct{}{}
	// 解锁改造属性
	filter := map[uint32]struct{}{}
	for _, id := range crystal.PassiveSkillIds {
		filter[id] = struct{}{}
	}
	cfgData.WalkCfgCrystalRedefineProp(func(cfg *cfgData.CrystalRedefinePropCfg) bool {
		if cfg.StageValue > crystal.Star {
			return false
		}
		if int(cfg.Id-1) < len(crystalCfg.PassiveSkillId) {
			skillID := crystalCfg.PassiveSkillId[cfg.Id-1]
			if _, ok := filter[skillID]; !ok {
				filter[skillID] = struct{}{}
				crystal.PassiveSkillIds = append(crystal.PassiveSkillIds, skillID)

				// 词条解锁
				this.entryData.Unlock(head, skillID)
			}
		}
		return true
	})
	// 成就统计
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_CrystalTotalLevel, this.GetTotalStar())
	// 成就触发
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_CrystalQualityRedefine, 1, crystal.Quality)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_CrystalQualityRedefine, 1, uint32(cfgEnum.EQuality_Any))
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_CrystalElementQualityRedefine, 1, crystal.Element, crystal.Quality)
	// 返回
	rsp.CrystalID = req.CrystalID
	rsp.CurStar = crystal.Star
	// 更新数据库
	this.UpdateSave(true)
	this.notifyCrystal(head)
	this.entryData.NotifyMainEntry(head, crystalCfg.MainSkillID)
	return nil
}

// 晶核升级
func (this *PlayerCrystalFun) CrystalUpgradeRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.CrystalUpgradeRequest)
	rsp := response.(*pb.CrystalUpgradeResponse)
	// 判断晶核是否存在
	crystal, ok := this.mapCrystal[req.CrystalID]
	if !ok {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_CrystalRobotNotExist, "head: %v, req: %v", head, req)
	}

	cfgCrystal := cfgData.GetCfgCrystal(req.CrystalID)
	if cfgCrystal == nil {
		return uerror.NewUErrorf(1, cfgData.GetCrystalErrorCode(req.CrystalID), "head: %v, req: %v", head, req)
	}
	// 判断是否到达最大等级
	if cfgData.IsMaxLevelCfgCrystalLevel(cfgCrystal.Quality, crystal.Level) {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_AlreadyMaximumLevel, "head: %v, req: %v", head, req)
	}
	// 加载配置
	levelCfg := cfgData.GetCfgNextCrystalLevel(cfgCrystal.Quality, crystal.Level)
	if levelCfg == nil {
		return uerror.NewUErrorf(1, cfgData.GetCrystalLevelErrorCode(cfgCrystal.Quality), "head: %v, req: %v", head, req)
	}
	//判断星级
	if crystal.Star < levelCfg.NeedStar {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_NeedCrystalStar, "head: %v, req: %v", head, req)
	}

	// 扣除通用道具
	code := this.getPlayerBagFun().DelItem(
		head,
		levelCfg.Reward.Kind,
		levelCfg.Reward.Id,
		levelCfg.Reward.Count,
		pb.EmDoingType_EDT_CrystalUpgrade,
	)
	if code != cfgEnum.ErrorCode_Success {
		return uerror.NewUErrorf(1, code, "head: %v, req: %v, crystal: %v", head, req, crystal)
	}
	// 提升等级
	crystal.Level = levelCfg.Level
	this.mapCrystalChange[req.CrystalID] = struct{}{}
	// 返回
	rsp.CrystalID = req.CrystalID
	rsp.CurLevel = crystal.Level

	//成就触发
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_CrystalQualityUpgrade, 1, crystal.Quality)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_CrystalQualityUpgrade, 1, uint32(cfgEnum.EQuality_Any))

	// 更新数据库
	this.UpdateSave(true)
	this.notifyCrystal(head)
	this.entryData.NotifyMainEntry(head, cfgCrystal.MainSkillID)
	return nil
}

// 晶核改造
func (this *PlayerCrystalFun) CrystalGenerateRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.CrystalGenerateRequest)
	// 判断晶核是否存在
	if _, ok := this.mapCrystal[req.CrystalID]; ok {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_CrystalAlreadyExist, "head: %v, req: %v", head, req)
	}
	// 加载配置
	cfgCrystal := cfgData.GetCfgCrystal(req.CrystalID)
	if cfgCrystal == nil {
		return uerror.NewUErrorf(1, cfgData.GetCrystalErrorCode(req.CrystalID), "head: %v, req: %v", head, req)
	}
	// 扣除晶核碎片
	errorCode := this.getPlayerBagFun().DelItem(head, uint32(cfgEnum.ESystemType_Item), cfgCrystal.PieceID, int64(cfgCrystal.PieceCount), pb.EmDoingType_EDT_CrystalGenerate)
	if errorCode != cfgEnum.ErrorCode_Success {
		return uerror.NewUErrorf(1, errorCode, "head: %v, req: %v", head, req)
	}
	// 增加晶核
	errorCode = this.AddCrystal(head, &pb.PBAddItemData{
		Kind:      uint32(cfgEnum.ESystemType_Crystal),
		Id:        req.CrystalID,
		Count:     1,
		DoingType: pb.EmDoingType_EDT_CrystalGenerate,
	})
	if errorCode != cfgEnum.ErrorCode_Success {
		return uerror.NewUErrorf(1, errorCode, "head: %v, req: %v", head, req)
	}
	return nil
}
