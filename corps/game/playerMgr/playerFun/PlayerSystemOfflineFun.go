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
	"sort"

	"github.com/golang/protobuf/proto"
)

// ----gomaker生成的模板-------
type PlayerSystemOfflineFun struct {
	PlayerFun
	offline *OfflineEntity
}

// --------------------通用接口实现------------------------------
// 初始化
func (this *PlayerSystemOfflineFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
}

// 新系统
func (this *PlayerSystemOfflineFun) NewPlayer() {
	this.offline = NewOfflineEntity(this, &pb.PBPlayerSystemOffline{})
	this.UpdateSave(true)
}

// 加载系统数据(system类型数据)
func (this *PlayerSystemOfflineFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	if pbSystem.Offline == nil {
		this.NewPlayer()
		return
	}
	this.offline = NewOfflineEntity(this, pbSystem.Offline)
	this.UpdateSave(false)
}

// 存储数据 返回存储标志(system类型数据)
func (this *PlayerSystemOfflineFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	pbSystem.Offline = this.offline.ToProto()
	return true
}

// 客户端数据
func (this *PlayerSystemOfflineFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = &pb.PBPlayerSystem{}
	}
	pbData.System.Offline = this.offline.ToProto()
}
func (this *PlayerSystemOfflineFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemOffline{}
}

// 设置玩家数据, web管理后台
func (this *PlayerSystemOfflineFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}

	pbSystem, ok := pbData.(*pb.PBPlayerSystemOffline)
	if !ok || pbSystem == nil {
		return false
	}
	this.offline = NewOfflineEntity(this, pbSystem)
	this.UpdateSave(true)
	return true
}

// -------------------业务层不常用接口实现---------------------------
// 对外接口（玩家心跳接口调用）
func (this *PlayerSystemOfflineFun) UpdateLogoutTime() {
	// 更新离线时间
	this.offline.UpdateLogoutTime()
	// 保存数据
	this.UpdateSave(true)
}

func (this *PlayerSystemOfflineFun) UpdateLoginTime() {
	// 更新离线时间
	this.offline.UpdateLoginTime()
	// 保存数据
	this.UpdateSave(true)
}

// 登录成功回调
func (this *PlayerSystemOfflineFun) LoadPlayerDBFinish() {
	if this.offline == nil {
		this.NewPlayer()
		// 保存数据
		this.UpdateSave(true)
	}
}

func (this *PlayerSystemOfflineFun) UpdatePrivilege(emPrivilegeType cfgEnum.PrivilegeType, uValue uint32) {
	if emPrivilegeType == cfgEnum.PrivilegeType_OfflineTime {
		this.offline.maxIncomeTime = uValue
	}
}

// 获取分钟奖励
func (this *PlayerSystemOfflineFun) AddRewardMin(head *pb.RpcHead, uMin uint32) cfgEnum.ErrorCode {
	mapID, _ := this.getPlayerSystemBattleHookFun().GetMapIdAndStageId()
	cfg := cfgData.GetCfgBattleHookMap(mapID)
	if cfg == nil {
		return cfgData.GetBattleHookMapErrorCode(mapID)
	}

	arrItem := &common.ItemInfo{
		Kind:  uint32(cfgEnum.ESystemType_LootGroup),
		Id:    cfg.OfflineLootId,
		Count: int64(uMin),
	}

	this.getPlayerBagFun().AddOneArrItem(head, arrItem, pb.EmDoingType_EDT_Offline, true)
	return cfgEnum.ErrorCode_Success
}

func (this *PlayerSystemOfflineFun) OfflineIncomeRewardRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.OfflineIncomeRewardRequest)
	//rsp := response.(*pb.OfflineIncomeRewardResponse)
	if err := this.offline.Reward(head, req.AdvertType); err != nil {
		return err
	}

	//成就触发
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_FinishedOffline, 1)

	// 保存数据
	this.UpdateSave(true)
	return nil
}

// 增加离线时间奖励
func (this *PlayerSystemOfflineFun) AddOfflineSeconds(head *pb.RpcHead, uSeconds uint32, doingType pb.EmDoingType) {
	if uSeconds < 60 {
		return
	}

	//判断挂机是否开启
	if !this.getPlayerSystemCommonFun().CheckSystemTypeOpen(cfgEnum.ESystemUnlockType_AFKSys) {
		return
	}

	this.offline.incomeTime = uSeconds
	this.offline.maxIncomeTime = this.getPlayerSystemCommonFun().GetPrivilege(cfgEnum.PrivilegeType_OfflineTime)

	mapID, _ := this.getPlayerSystemBattleHookFun().GetMapIdAndStageId()
	cfg := cfgData.GetCfgBattleHookMap(mapID)
	if cfg == nil {
		return
	}

	arrItem := &common.ItemInfo{
		Kind:  uint32(cfgEnum.ESystemType_LootGroup),
		Id:    cfg.OfflineLootId,
		Count: int64(uSeconds / 60),
	}

	// 存收益时间
	items, equips := this.offline.GetDisplayRewards(arrItem)
	this.offline.SetDisplayReward(items, equips)
	return
}

// ----------------------------------离线entity封装--------------------------------
const (
	OFFLINE_INTERVAL_TIME = 60
)

type OfflineEntity struct {
	*PlayerSystemOfflineFun
	loginTime           uint64 // 上次登录时间
	logoutTime          uint64 // 登出时间
	incomeTime          uint32 // 收益时长
	maxIncomeTime       uint32 // 最大收益时长
	totalEquipment      uint32
	addEquipmentBag     uint32
	splitEquipmentScore uint64
	items               map[uint32]*pb.PBAddItemData // 道具
	equips              []*pb.PBAddItemData          // 装备
}

func NewOfflineEntity(pFather *PlayerSystemOfflineFun, data *pb.PBPlayerSystemOffline) *OfflineEntity {
	ret := &OfflineEntity{
		PlayerSystemOfflineFun: pFather,
		loginTime:              data.LoginTime,
		logoutTime:             data.LogoutTime,
		items:                  make(map[uint32]*pb.PBAddItemData),
	}
	for _, item := range data.Rewards {
		switch cfgEnum.ESystemType(item.Kind) {
		case cfgEnum.ESystemType_Equipment:
			ret.equips = append(ret.equips, item)
		case cfgEnum.ESystemType_Item:
			ret.items[item.Id] = item
		}
	}
	return ret
}

func (d *OfflineEntity) getRewards() (rets []*pb.PBAddItemData) {
	for _, item := range d.items {
		rets = append(rets, item)
	}
	rets = append(rets, d.equips...)
	return
}

func (d *OfflineEntity) ToProto() *pb.PBPlayerSystemOffline {
	return &pb.PBPlayerSystemOffline{
		LoginTime:           d.loginTime,
		LogoutTime:          d.logoutTime,
		IncomTime:           d.incomeTime,
		MaxIncomTime:        d.maxIncomeTime,
		Rewards:             d.getRewards(),
		TotalEquipment:      d.totalEquipment,
		AddEquipmentBag:     d.addEquipmentBag,
		SplitEquipmentScore: d.splitEquipmentScore,
	}
}

// 领取奖励
func (d *OfflineEntity) Reward(head *pb.RpcHead, uAdvertType uint32) error {
	// 恭喜获得
	items := d.getRewards()
	if len(items) <= 0 {
		if uAdvertType != uint32(cfgEnum.EAdvertType_OfflineAward) {
			return uerror.NewUError(1, cfgEnum.ErrorCode_OfflineRewardEmpty, head)
		} else {
			items = make([]*pb.PBAddItemData, 0)
		}
	}

	//广告额外给奖励
	if uAdvertType == uint32(cfgEnum.EAdvertType_OfflineAward) {
		uCode := d.getPlayerSystemCommonFun().AddAdvert(head, uAdvertType)
		if uCode != cfgEnum.ErrorCode_Success {
			return uerror.NewUError(1, uCode, "Reward advert %d", uAdvertType)
		}

		cfgAdvert := cfgData.GetCfgAdvertConfig(uAdvertType)
		if cfgAdvert == nil {
			return uerror.NewUError(1, cfgData.GetAdvertConfigErrorCode(uAdvertType), "Reward advert %d", uAdvertType)
		}

		mapID, _ := d.getPlayerSystemBattleHookFun().GetMapIdAndStageId()
		cfg := cfgData.GetCfgBattleHookMap(mapID)
		if cfg == nil {
			return uerror.NewUError(1, cfgData.GetBattleHookMapErrorCode(mapID), "Reward mapID %d", mapID)
		}

		arrItem := &common.ItemInfo{
			Kind:  uint32(cfgEnum.ESystemType_LootGroup),
			Id:    cfg.OfflineLootId,
			Count: int64(cfgAdvert.Param),
		}

		mapItem, _, listEquip, _, _, _ := d.InnerGetDisplayRewards(arrItem, pb.EmDoingType_EDT_Advert)
		for _, item := range mapItem {
			items = append(items, item)
		}
		items = append(items, listEquip...)
	} else {

	}

	plog.Trace("(d *OfflineEntity) head: %v, reward: %v", head, items)
	d.getPlayerBagFun().CommonPrizeNotify(head, items, pb.EmDoingType_EDT_Offline)
	// 重置奖励
	d.items = make(map[uint32]*pb.PBAddItemData)
	d.equips = d.equips[:0]
	d.incomeTime = 0
	d.totalEquipment = 0
	d.addEquipmentBag = 0
	d.splitEquipmentScore = 0
	return nil
}

// 设置登出时间
func (d *OfflineEntity) UpdateLogoutTime() {
	d.logoutTime = base.GetNow()
}

// 登录
func (d *OfflineEntity) UpdateLoginTime() {
	if rr := d.HasOfflineReward(); rr != nil {
		items, equips := d.GetDisplayRewards(rr)

		d.SetDisplayReward(items, equips)
	}
}

func (d *OfflineEntity) HasOfflineReward() *common.ItemInfo {
	defer func() {
		d.logoutTime = d.loginTime
		plog.Trace("(d *OfflineEntity) refresh logoutTime: %d", d.logoutTime)
	}()

	// 更新登录时间
	d.loginTime = base.GetNow()

	//判断挂机是否开启
	if !d.getPlayerSystemCommonFun().CheckSystemTypeOpen(cfgEnum.ESystemUnlockType_AFKSys) {
		return nil
	}

	if d.logoutTime == 0 {
		d.logoutTime = base.GetNow()
	}

	mapID, _ := d.getPlayerSystemBattleHookFun().GetMapIdAndStageId()
	cfg := cfgData.GetCfgBattleHookMap(mapID)
	if cfg == nil {
		plog.Error("(d *OfflineEntity) BattleHook config not found, mapID: %d", mapID)
		return nil
	}

	// 计算离线收益时长
	uStepSeconds := uint32(d.loginTime - d.logoutTime)
	if d.logoutTime >= d.loginTime || uStepSeconds < 60 {
		plog.Error("(d *OfflineEntity) HasOfflineReward  uid %d logoutTime:%d loginTime:%d uStepSeconds:%d", d.AccountId, d.logoutTime, d.loginTime, uStepSeconds)
		return nil
	}

	plog.Error("(d *OfflineEntity) Offline uid: %d, loginTime: %d, logoutTime: %d, diff: %ds", d.AccountId, d.loginTime, d.logoutTime, uStepSeconds)
	d.maxIncomeTime = d.getPlayerSystemCommonFun().GetPrivilege(cfgEnum.PrivilegeType_OfflineTime)
	// 获取关卡最大离线时间 返回秒
	if d.incomeTime > d.maxIncomeTime {
		d.incomeTime = d.maxIncomeTime
		return nil
	}

	if d.incomeTime+uStepSeconds > d.maxIncomeTime {
		uStepSeconds = d.maxIncomeTime - d.incomeTime
	}
	if uStepSeconds < 60 {
		return nil
	}

	// 存收益时间
	d.incomeTime += uStepSeconds / 60 * 60
	return &common.ItemInfo{
		Kind:  uint32(cfgEnum.ESystemType_LootGroup),
		Id:    cfg.OfflineLootId,
		Count: int64(uStepSeconds / 60),
	}
}

func (d *OfflineEntity) InnerGetDisplayRewards(item *common.ItemInfo, doingType pb.EmDoingType) (mapItem map[uint32]*pb.PBAddItemData, listEquip []*pb.PBAddItemData, addEquips []*pb.PBAddItemData,
	totalEquipment uint32, addEquipmentBag uint32, splitEquipmentScore uint64) {
	mapItem = map[uint32]*pb.PBAddItemData{}                  // 合并的道具
	splitEquips := []*pb.PBAddItemData{}                      // 分解装备
	filter := d.GetPlayerEquipmentFun().GetAutoSplitQuality() // 需要自动分解的品质

	// 获取离线收益
	results := d.getPlayerBagFun().GetPbItems([]*common.ItemInfo{item}, doingType)
	plog.Trace("(d *OfflineEntity) uid: %d, rewards: %v", d.AccountId, results)
	for _, tmpitem := range results {

		switch cfgEnum.ESystemType(tmpitem.Kind) {
		case cfgEnum.ESystemType_Equipment:
			if tmpitem.Equipment == nil {
				continue
			}
			totalEquipment++
			if _, ok := filter[tmpitem.Equipment.Quality]; ok {
				// 直接分解
				splitEquips = append(splitEquips, tmpitem)
			} else {
				// 展示
				listEquip = append(listEquip, tmpitem)
			}
		case cfgEnum.ESystemType_Item:
			if val, ok := mapItem[tmpitem.Id]; ok {
				val.Count += tmpitem.Count
			} else {
				mapItem[tmpitem.Id] = tmpitem
				addEquips = append(addEquips, tmpitem)
			}
		}
	}

	// 排序品质
	sort.SliceStable(listEquip, func(i, j int) bool {
		if listEquip[i].Equipment.Quality != listEquip[j].Equipment.Quality {
			return listEquip[i].Equipment.Quality > listEquip[j].Equipment.Quality
		}
		return listEquip[i].Equipment.Star > listEquip[j].Equipment.Star
	})

	// 获取发放奖励
	spare := d.GetPlayerEquipmentFun().GetSpareBag()
	if spare <= 0 {
		// 背包没有空间，直接全部分解
		splitEquips = append(splitEquips, listEquip...)
	} else if spare >= uint32(len(listEquip)) {
		// 背包空间有剩余，全部加入背包
		addEquipmentBag = uint32(len(listEquip))
		addEquips = append(addEquips, listEquip...)
	} else {
		// 一部分加入背包，一部分被分解
		addEquipmentBag = spare
		addEquips = append(addEquips, listEquip[:spare]...)
		splitEquips = append(splitEquips, listEquip[spare:]...)
	}

	// 分解装备
	head := &pb.RpcHead{Id: d.AccountId}
	splitEquipmentScore = d.GetPlayerEquipmentFun().SplitEquipment(head, true, splitEquips...)
	// 发送奖励
	plog.Trace("离线收益：%v", addEquips)
	if errCode := d.getPlayerBagFun().AddPbItems(head, addEquips, doingType, false); errCode != cfgEnum.ErrorCode_Success {
		plog.Error("(d *OfflineEntity) Offline AddPbItems is failed, errorCode: %d, items: %v, uid: %d", errCode, addEquips, d.AccountId)
	}

	if splitEquipmentScore > 0 {
		mapItem[uint32(pb.EmItemExpendType_EIET_SplitScore)] = &pb.PBAddItemData{
			Kind:      uint32(cfgEnum.ESystemType_Item),
			Id:        uint32(pb.EmItemExpendType_EIET_SplitScore),
			Count:     int64(splitEquipmentScore),
			DoingType: doingType,
		}
	}

	plog.Trace("(d *OfflineEntity) uid: %d, history add: %d, totol: %d, splitScore: %d", d.AccountId, d.addEquipmentBag, d.totalEquipment, d.splitEquipmentScore)
	return
}

func (d *OfflineEntity) GetDisplayRewards(item *common.ItemInfo) (mapItem map[uint32]*pb.PBAddItemData, listEquip []*pb.PBAddItemData) {
	totalEquipment := uint32(0)
	addEquipmentBag := uint32(0)
	splitEquipmentScore := uint64(0)
	mapItem, listEquip, _, totalEquipment, addEquipmentBag, splitEquipmentScore = d.InnerGetDisplayRewards(item, pb.EmDoingType_EDT_Offline)
	d.totalEquipment += totalEquipment
	d.addEquipmentBag += addEquipmentBag
	d.splitEquipmentScore += splitEquipmentScore
	return
}

func (d *OfflineEntity) SetDisplayReward(items map[uint32]*pb.PBAddItemData, equips []*pb.PBAddItemData) {
	// 合并道具
	for key, vv := range items {
		if val, ok := d.items[key]; ok {
			val.Count += val.Count
		} else {
			d.items[key] = vv
		}
	}

	// 何如装备
	d.equips = append(d.equips, equips...)
	// 排序品质
	sort.SliceStable(d.equips, func(i, j int) bool {
		if d.equips[i].Equipment.Quality != d.equips[j].Equipment.Quality {
			return d.equips[i].Equipment.Quality > d.equips[j].Equipment.Quality
		}
		return d.equips[i].Equipment.Star > d.equips[j].Equipment.Star
	})
	// 截断
	if uMaxEquipCount := cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_OFFLINE_MAX_EQUIP); uint32(len(d.equips)) > uMaxEquipCount {
		d.equips = d.equips[:uMaxEquipCount]
	}

	// (再一次)主动推送一下数据
	pbResponse := &pb.AllPlayerInfoNotify{
		PacketHead: &pb.IPacket{},
		PlayerData: &pb.PBPlayerData{},
	}
	pbResponse.Mark = base.SetBit32(pbResponse.Mark, uint32(pb.PlayerDataType_SystemOffline), true)
	d.PlayerSystemOfflineFun.SaveDataToClient(pbResponse.PlayerData)
	cluster.SendToClient(&pb.RpcHead{Id: d.AccountId}, pbResponse, cfgEnum.ErrorCode_Success)
}
