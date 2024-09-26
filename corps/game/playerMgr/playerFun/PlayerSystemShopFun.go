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
	PlayerSystemShopFun struct {
		PlayerFun
		pbBlackShop   *pb.PBBlackShop
		mapBlackGoods map[uint32]*pb.PBShopGoodInfo
		mapShopInfo   map[uint32]*PlayerShopInfo
	}

	PlayerShopInfo struct {
		*pb.PBShopInfo
		mapBuyData map[uint32]*pb.PBU32U32
	}
)

// 注册时，获取player的列表
func (this *PlayerSystemShopFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.pbBlackShop = &pb.PBBlackShop{RefreshInfo: &pb.PBShopRefreshInfo{}}
	this.mapBlackGoods = make(map[uint32]*pb.PBShopGoodInfo)
	this.mapShopInfo = make(map[uint32]*PlayerShopInfo)
}

func (this *PlayerSystemShopFun) LoadPlayerDBFinish() {
	if this.pbBlackShop == nil {
		this.pbBlackShop = &pb.PBBlackShop{RefreshInfo: &pb.PBShopRefreshInfo{}}
	}

	if this.pbBlackShop.RefreshInfo == nil {
		this.pbBlackShop.RefreshInfo = &pb.PBShopRefreshInfo{}
	}
}
func (this *PlayerSystemShopFun) LoadComplete() {
	this.checkShopTypeOpen()
}

// 新成员初始化
func (this *PlayerSystemShopFun) NewPlayer() {
}

func (this *PlayerSystemShopFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	if pbSystem.Shop == nil {
		pbSystem.Shop = &pb.PBPlayerSystemShop{}
	}

	if pbSystem.Shop.BlackShop == nil {
		pbSystem.Shop.BlackShop = &pb.PBBlackShop{}
	}

	this.loadData(pbSystem.Shop)

	this.pbBlackShop = pbSystem.Shop.BlackShop

	this.UpdateSave(false)
}

func (this *PlayerSystemShopFun) loadData(pbData *pb.PBPlayerSystemShop) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemShop{}
	}

	this.pbBlackShop = pbData.BlackShop

	this.mapBlackGoods = make(map[uint32]*pb.PBShopGoodInfo)
	for i := 0; i < len(this.pbBlackShop.Items); i++ {
		this.mapBlackGoods[this.pbBlackShop.Items[i].GoodsID] = this.pbBlackShop.Items[i]
	}

	this.mapShopInfo = make(map[uint32]*PlayerShopInfo)
	for _, shopInfo := range pbData.ShopList {
		this.mapShopInfo[shopInfo.ShopType] = &PlayerShopInfo{
			PBShopInfo: shopInfo,
			mapBuyData: make(map[uint32]*pb.PBU32U32),
		}

		for _, item := range shopInfo.Items {
			this.mapShopInfo[shopInfo.ShopType].mapBuyData[item.Key] = item
		}
	}
	this.UpdateSave(true)
}

func (this *PlayerSystemShopFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	if pbSystem.Shop == nil {
		pbSystem.Shop = &pb.PBPlayerSystemShop{}
	}

	this.SavePb(pbSystem.Shop)

	return this.BSave
}

// 保存
func (this *PlayerSystemShopFun) SavePb(pbShopData *pb.PBPlayerSystemShop) {

	this.pbBlackShop.Items = make([]*pb.PBShopGoodInfo, 0)
	for _, v := range this.mapBlackGoods {
		this.pbBlackShop.Items = append(this.pbBlackShop.Items, v)
	}

	pbShopData.BlackShop = this.pbBlackShop

	pbShopData.ShopList = make([]*pb.PBShopInfo, 0)
	for _, shopInfo := range this.mapShopInfo {
		pbShopData.ShopList = append(pbShopData.ShopList, shopInfo.PBShopInfo)
	}

}

// 更新客户端数据
func (this *PlayerSystemShopFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = new(pb.PBPlayerSystem)
	}

	this.SaveSystem(pbData.System)
}
func (this *PlayerSystemShopFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemShop{}
}

// 设置玩家数据
func (this *PlayerSystemShopFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}

	pbSystem := pbData.(*pb.PBPlayerSystemShop)
	if pbSystem == nil {
		return false
	}

	this.loadData(pbSystem)

	this.UpdateSave(true)
	return true
}

func (this *PlayerSystemShopFun) ShopRefreshTimeNotify() {
	pbNotify := &pb.ShopRefreshTimeNotify{
		PacketHead:  &pb.IPacket{},
		ShopType:    uint32(cfgEnum.EShopType_BlackShop),
		RefreshInfo: this.pbBlackShop.RefreshInfo,
	}

	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbNotify, cfgEnum.ErrorCode_Success)
}
func (this *PlayerSystemShopFun) Heat() {
	if this.pbBlackShop == nil {
		return
	}
	uCurTime := base.GetNow()
	if this.pbBlackShop.RefreshInfo.NextFreeRefreshTime > 0 && this.pbBlackShop.RefreshInfo.NextFreeRefreshTime <= uCurTime {
		if this.pbBlackShop.RefreshInfo.DailyFreeUseCount > 0 {
			this.pbBlackShop.RefreshInfo.DailyFreeUseCount--

			cfgShop := cfgData.GetCfgShopConfig(uint32(cfgEnum.EShopType_BlackShop))
			if cfgShop != nil {
				this.pbBlackShop.RefreshInfo.NextFreeRefreshTime = uCurTime + uint64(cfgShop.RefreshRecoverTime)

			}

		}

		if this.pbBlackShop.RefreshInfo.DailyFreeUseCount == 0 {
			this.pbBlackShop.RefreshInfo.NextFreeRefreshTime = 0
		}

		this.ShopRefreshTimeNotify()
	}
}

// 更新最大刷新次数
func (this *PlayerSystemShopFun) UpdatePrivilege(emPrivilegeType cfgEnum.PrivilegeType, uValue uint32) {
	if emPrivilegeType == cfgEnum.PrivilegeType_BlackShopRefreshCount {
		this.pbBlackShop.RefreshInfo.DailyFreeMaxCount = this.GetDailyFreeMaxCount(uint32(cfgEnum.EShopType_BlackShop))
		this.ShopRefreshTimeNotify()
	}
}

// 获取最大次数
func (this *PlayerSystemShopFun) GetDailyFreeMaxCount(shopType uint32) uint32 {
	cfgShop := cfgData.GetCfgShopConfig(shopType)
	if cfgShop == nil {
		return 0
	}

	uMax := cfgShop.FreeRefreshCount
	if shopType == uint32(cfgEnum.EShopType_BlackShop) {
		uMax += this.getPlayerSystemCommonFun().GetPrivilege(cfgEnum.PrivilegeType_BlackShopRefreshCount)
	}
	return uMax
}

// 添加一个新商店
func (this *PlayerSystemShopFun) getNextRefreshTime(cfgShop *cfgData.ShopConfigCfg) uint64 {

	uNextFreshTime := uint64(0)

	uCurTime := base.GetNow()
	switch cfgEnum.EActivityRefreshType(cfgShop.RefreshType) {
	case cfgEnum.EActivityRefreshType_Daily:
		uNextFreshTime = base.GetZeroTimestamp(uCurTime, 1)
	case cfgEnum.EActivityRefreshType_Week:
		uNextFreshTime = base.GetWeekZeroTimestamp(uCurTime, 1)
	case cfgEnum.EActivityRefreshType_Month:
		uNextFreshTime = base.GetMonthZeroTimestamp(uCurTime, 1)
	case cfgEnum.EActivityRefreshType_FixDay:
		if cfgShop.RefreshParam <= 0 {
			plog.Error("(this *PlayerSystemShopFun) PassDay %d RefreshParam is 0", cfgShop.ShopType)
			return 0
		}

		//循环判断
		uBeginTime := base.GetNow()
		if cfgShop.ConditionInfo.Type == uint32(cfgEnum.EConditionType_OpenServerDay) {
			uBeginTime = this.getPlayerBaseFun().GetOpenSeverTime()
		} else if cfgShop.ConditionInfo.Type == uint32(cfgEnum.EConditionType_RegDay) {
			uBeginTime = this.getPlayerBaseFun().GetRegTime()
		} else if cfgShop.ConditionInfo.Type == uint32(cfgEnum.EConditionType_AllServerDay) {
			uBeginTime = serverCommon.GetOpenServerTime()
		} else if cfgShop.ConditionInfo.Type == uint32(cfgEnum.EConditionType_SystemOpenType) {
			uBeginTime = serverCommon.GetOpenServerTime()
		}
		if cfgShop.ConditionInfo.Type != uint32(cfgEnum.EConditionType_SystemOpenType) {
			uBeginTime = base.GetZeroTimestamp(uBeginTime, int32(cfgShop.ConditionInfo.BeginDay))
		} else {
			uBeginTime = base.GetZeroTimestamp(uBeginTime, 0)
		}

		uPassDay := base.DiffDays(uBeginTime, base.GetNow()) + 1 // 1,2,3 / 3  3
		uNextFreshTime = uBeginTime + uint64(base.CeilU32(uPassDay, cfgShop.RefreshParam)*cfgShop.RefreshParam*3600*24)
	}

	return uNextFreshTime
}

// 隔天刷新通知
func (this *PlayerSystemShopFun) PassDay(isDay, isWeek, isMonth bool) {
	//更新刷新时间
	this.pbBlackShop.RefreshInfo.DailyBuyCount = 0
	this.pbBlackShop.RefreshInfo.DailyFreeUseCount = 0
	this.pbBlackShop.RefreshInfo.NextFreeRefreshTime = 0
	this.pbBlackShop.RefreshInfo.DailyFreeMaxCount = this.GetDailyFreeMaxCount(uint32(cfgEnum.EShopType_BlackShop))
	this.pbBlackShop.NextRefreshTime = base.GetZeroTimestamp(base.GetNow(), 1)

	//判断刷新时间
	this.checkShopTypeOpen()

	this.refreshBlackShop(true)
}

// 系统开启
func (this *PlayerSystemShopFun) OnSystemOpenTypes(head *pb.RpcHead, arrTypes []uint32) {
	bOpen := false
	for _, systemType := range arrTypes {
		//检查pb开启
		mapCfg := cfgData.GetSystemOpenShopConfig(systemType)
		if mapCfg != nil || len(mapCfg) > 0 {
			bOpen = true
			break
		}
	}
	if bOpen {
		this.checkShopTypeOpen()
	}
}

// 商店是否有红点
func (this *PlayerSystemShopFun) getShopTypeHaveRed(ShopType uint32, mapBuy map[uint32]*pb.PBU32U32) uint32 {
	//是否有红点 如果有免费的
	listCfg := cfgData.GetCfgExchangeShopListConfig(ShopType)
	for _, cfgGood := range listCfg {
		if cfgGood.Currency == nil && cfgGood.ProductId == 0 {
			if mapBuy != nil {
				if pbBuy, ok := mapBuy[cfgGood.Id]; ok {
					if pbBuy.Value >= cfgGood.LimitBuyNum {
						continue
					}
				}
			}

			return 1
		}
	}

	return 0
}
func (this *PlayerSystemShopFun) repairShopRed() {
	pbShopListRedNotify := &pb.ShopRedNotify{
		PacketHead: &pb.IPacket{Id: this.AccountId},
	}
	for _, shopData := range this.mapShopInfo {
		if shopData.HaveRed > 0 {
			continue
		}

		shopData.HaveRed = this.getShopTypeHaveRed(shopData.ShopType, shopData.mapBuyData)
		if shopData.HaveRed > 0 {
			pbShopListRedNotify.ShopRedList = append(pbShopListRedNotify.ShopRedList, &pb.PBU32U32{
				Key:   shopData.ShopType,
				Value: 1,
			})
		}
	}

	//同步红点
	if len(pbShopListRedNotify.ShopRedList) > 0 {
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbShopListRedNotify, cfgEnum.ErrorCode_Success)
		this.UpdateSave(true)
	}
}

// 检查商店刷新
func (this *PlayerSystemShopFun) checkShopTypeOpen() {
	//判断刷新时间
	uCurTime := base.GetNow()
	mapCfg := cfgData.GetAllCfgShopConfig()
	pbShopListNotify := &pb.ShopListNotify{
		PacketHead: &pb.IPacket{Id: this.AccountId},
	}
	pbShopListRedNotify := &pb.ShopRedNotify{
		PacketHead: &pb.IPacket{Id: this.AccountId},
	}
	for _, cfgShop := range mapCfg {
		//未到刷新时间 不处理
		if _, ok := this.mapShopInfo[cfgShop.ShopType]; ok {
			if this.mapShopInfo[cfgShop.ShopType].NextRefreshTime == 0 || uCurTime < this.mapShopInfo[cfgShop.ShopType].NextRefreshTime {
				continue
			}
		}

		//判断条件
		if uCode, _, _ := this.getPlayerBaseFun().CheckCondition([]*common.ConditionInfo{cfgShop.ConditionInfo}); uCode != cfgEnum.ErrorCode_Success {
			continue
		}

		this.mapShopInfo[cfgShop.ShopType] = &PlayerShopInfo{
			PBShopInfo: &pb.PBShopInfo{
				ShopType:        cfgShop.ShopType,
				NextRefreshTime: this.getNextRefreshTime(cfgShop),
				HaveRed:         this.getShopTypeHaveRed(cfgShop.ShopType, nil),
			},
			mapBuyData: make(map[uint32]*pb.PBU32U32),
		}

		//通知客户端
		if this.mapShopInfo[cfgShop.ShopType].HaveRed > 0 {
			pbShopListRedNotify.ShopRedList = append(pbShopListRedNotify.ShopRedList, &pb.PBU32U32{
				Key:   cfgShop.ShopType,
				Value: 1,
			})
		}

		pbShopListNotify.ShopList = append(pbShopListNotify.ShopList, &pb.PBU32U64{Key: cfgShop.ShopType, Value: this.mapShopInfo[cfgShop.ShopType].NextRefreshTime})
	}

	//同步时间
	if len(pbShopListNotify.ShopList) > 0 {
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbShopListNotify, cfgEnum.ErrorCode_Success)
		this.UpdateSave(true)
	}
	//同步红点
	if len(pbShopListRedNotify.ShopRedList) > 0 {
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbShopListRedNotify, cfgEnum.ErrorCode_Success)
		this.UpdateSave(true)
	}
}

// 手动刷新黑市商店
func (this *PlayerSystemShopFun) refreshBlackShop(isDaily bool) {
	// 获取掉落折扣率
	vals := this.getEntry().Get(uint32(cfgEnum.EntryEffectType_DiscountProbability), uint32(cfgEnum.EShopType_BlackShop))
	// 获取广告观看次数
	adTimes := entry.ToValue(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_WatchAdvertise), uint32(cfgEnum.EAdvertType_BlackShop))...)
	//不是自动刷新的 需要继承上一次的免费购买次数 如果没按费次数用完 直接不刷
	uPreFreeCount := uint32(0)
	uPreFreeBuyCount := uint32(0)
	bHaveFree := false
	if isDaily {
		bHaveFree = true
	} else {
		for _, info := range this.mapBlackGoods {
			cfg := cfgData.GetCfgBlackShop(info.GoodsID)
			if cfg != nil && cfg.FreeTimes > 0 && info.BuyTimes < cfg.LimitBuyNum+adTimes {
				uPreFreeCount = info.FreeTimes
				uPreFreeBuyCount = info.BuyTimes
				bHaveFree = true
				break
			}
		}
	}

	mapBlackGoods := make(map[uint32]*pb.PBShopGoodInfo)
	mapId, stageId := this.getPlayerSystemBattleFun().GetMapIdAndStageId(pb.EmBattleType_EBT_Hook)
	cfgs := cfgData.RandBlackShopCfg(mapId, stageId, !bHaveFree)
	if cfgs != nil {
		now := base.GetNow()
		for _, cfg := range cfgs {
			if cfgEnum.EConditionType(cfg.ConditionInfo.Type) == cfgEnum.EConditionType_OpenServerDay {
				// 检查条件
				uOpenServerDays := serverCommon.GetOpenServerDays()
				if (uOpenServerDays < cfg.ConditionInfo.ConditionOpenServer.BeginDay) ||
					(cfg.ConditionInfo.ConditionOpenServer.EndDay > 0 && uOpenServerDays > cfg.ConditionInfo.ConditionOpenServer.EndDay) ||
					(cfg.OpenTime > 0 && (now < cfg.OpenTime || now > cfg.EndTime)) {
					continue
				}
			}
			if !bHaveFree && cfg.FreeTimes > 0 {
				continue
			}
			pbData := &pb.PBShopGoodInfo{
				GoodsID:  cfg.Id,
				Discount: uint32(base.RandInt64Probability(int64(cfg.Discount.Min), int64(cfg.Discount.Max), vals...)),
			}
			if cfg.FreeTimes > 0 {
				pbData.FreeTimes = uPreFreeCount
				pbData.BuyTimes = uPreFreeBuyCount
			}

			if cfg.Goods.Kind == uint32(cfgEnum.ESystemType_Equipment) {
				pbAddItem := this.getPlayerBagFun().GetPbItems([]*common.ItemInfo{
					{
						Id:     cfg.Goods.Id,
						Kind:   cfg.Goods.Kind,
						Count:  1,
						Params: cfg.Goods.Params,
					}}, pb.EmDoingType_EDT_BlackShop)

				if len(pbAddItem) > 0 {
					pbData.Equipment = pbAddItem[0].Equipment
				}
			}

			mapBlackGoods[cfg.Id] = pbData
		}
	}

	this.mapBlackGoods = make(map[uint32]*pb.PBShopGoodInfo)
	this.pbBlackShop.Items = make([]*pb.PBShopGoodInfo, 0)
	for tmpid, tmpinfo := range mapBlackGoods {
		this.mapBlackGoods[tmpid] = tmpinfo
		this.pbBlackShop.Items = append(this.pbBlackShop.Items, tmpinfo)
	}

	//存储数据
	this.UpdateSave(true)

	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.ShopUpdateNotify{
		PacketHead: &pb.IPacket{},
		ShopType:   uint32(pb.EmShopType_EST_BlackShop),
		Shop: &pb.PBPlayerSystemShop{
			BlackShop: this.pbBlackShop,
		},
	}, cfgEnum.ErrorCode_Success)
}

func (this *PlayerSystemShopFun) ShopBuyRequest(head *pb.RpcHead, req *pb.ShopBuyRequest) {
	uCode := this.shopBuyRequest(head, req.ShopType, req.GoodsID, req.AdvertType)

	//成就触发
	if uCode == cfgEnum.ErrorCode_Success {
		this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ShopBuy, 1, req.ShopType)
		this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ShopBuy, 1, 0)

		//成就触发
		if req.AdvertType > uint32(cfgEnum.EAdvertType_None) {
			this.getPlayerSystemCommonFun().OnAdvert(head, req.AdvertType)
		}
	}
}

// 通用购买接口
func (this *PlayerSystemShopFun) shopBuyRequest(head *pb.RpcHead, shopType, goodsID uint32, uAdvertType uint32) cfgEnum.ErrorCode {
	uCode := cfgEnum.ErrorCode_Success
	switch pb.EmShopType(shopType) {
	case pb.EmShopType_EST_BlackShop:
		uCode = this.shopBuyBlackShop(head, goodsID, uAdvertType)
		cluster.SendToClient(head, &pb.ShopBuyResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	default:
		plog.Print(this.AccountId, cfgEnum.ErrorCode_ShopTypeNotSupported, shopType, goodsID, uAdvertType)

	}
	return uCode

}

func (this *PlayerSystemShopFun) shopBuyBlackShop(head *pb.RpcHead, uGoodId uint32, uAdvertType uint32) cfgEnum.ErrorCode {
	pbGoodInfo, ok := this.mapBlackGoods[uGoodId]
	if !ok {
		return plog.Print(head.Id, cfgEnum.ErrorCode_NotSupported, "GoodsID", uGoodId)
	}

	// 加载配置
	cfgBlackShop := cfgData.GetCfgBlackShop(uGoodId)
	if cfgBlackShop == nil {
		return plog.Print(head.Id, cfgData.GetBlackShopErrorCode(uGoodId), "BlackShopConfig", uGoodId)
	}

	//判断背包已满
	if cfgBlackShop.Goods.Kind == uint32(cfgEnum.ESystemType_Equipment) {
		if this.GetPlayerEquipmentFun().GetSpareBag() < uint32(cfgBlackShop.Goods.Count) {
			cluster.SendCommonToClient(head, cfgEnum.ErrorCode_BagFull)
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_BagFull, uGoodId, uAdvertType)
		}
	} else if cfgBlackShop.Goods.Kind == uint32(cfgEnum.ESystemType_Hero) {
		if this.getPlayerHeroFun().GetSpareBag() < uint32(cfgBlackShop.Goods.Count) {
			cluster.SendCommonToClient(head, cfgEnum.ErrorCode_HeroBagFull)
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_HeroBagFull, uGoodId, uAdvertType)
		}
	}

	// 获取广告观看次数
	adTimes := entry.ToValue(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_WatchAdvertise), uint32(cfgEnum.EAdvertType_BlackShop))...)

	// 判断是否为免费
	if pbGoodInfo.FreeTimes < cfgBlackShop.FreeTimes {
		// 扣除免费次数
		pbGoodInfo.FreeTimes++
	} else {
		// 判断是否限购
		if pbGoodInfo.BuyTimes >= cfgBlackShop.LimitBuyNum+adTimes {
			return plog.Print(head.Id, cfgEnum.ErrorCode_BuyTimesLimit, "BlackShop", pbGoodInfo.BuyTimes, cfgBlackShop.LimitBuyNum)
		}

		// 判断玩家货币是否足够
		if cfgEnum.ESystemType(cfgBlackShop.Currency.Kind) == cfgEnum.ESystemType_Adverting {
			if uAdvertType == uint32(cfgEnum.EAdvertType_None) {
				return plog.Print(head.Id, cfgEnum.ErrorCode_AdvertTypeNotSupported, cfgBlackShop.Currency)
			}
		} else {
			errCode := this.getPlayerBagFun().DelItem(head, cfgBlackShop.Currency.Kind, cfgBlackShop.Currency.Id, cfgBlackShop.Currency.Count*int64(pbGoodInfo.Discount)/10, pb.EmDoingType_EDT_BlackShop)
			if errCode != cfgEnum.ErrorCode_Success {
				return plog.Print(head.Id, errCode, cfgBlackShop.Currency)
			}
		}

		// 增加购买次数
		pbGoodInfo.BuyTimes++
	}

	//奖励
	pbAddItem := &pb.PBAddItemData{
		Id:        cfgBlackShop.Goods.Id,
		Kind:      cfgBlackShop.Goods.Kind,
		Count:     cfgBlackShop.Goods.Count,
		Equipment: pbGoodInfo.Equipment,
		DoingType: pb.EmDoingType_EDT_BlackShop,
		Params:    cfgBlackShop.Goods.Params,
	}

	//判断词条 双倍奖励
	arrAddItem := make([]*pb.PBAddItemData, 0)
	arrAddItem = append(arrAddItem, pbAddItem)
	if uAdvertType == uint32(cfgEnum.EAdvertType_BlackShop) {
		propValue := entry.ToValue(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_AdertiseReward), uint32(cfgEnum.EAdvertType_BlackShop))...)
		if base.IsRadio(propValue) {
			pbEntryAddItem := *pbAddItem
			pbEntryAddItem.DoingType = pb.EmDoingType_EDT_Entry
			plog.Trace("EntryEffectType_AdertiseReward rewards: %v", pbEntryAddItem)
			arrAddItem = append(arrAddItem, &pbEntryAddItem)
		}
	}

	this.getPlayerBagFun().AddPbItems(head, arrAddItem, pb.EmDoingType_EDT_BlackShop, true)

	this.UpdateSave(true)

	//同步数据给客户端
	cluster.SendToClient(head, &pb.ShopUpdateOneGoodsNotify{
		PacketHead: &pb.IPacket{},
		ShopType:   uint32(pb.EmShopType_EST_BlackShop),
		ShopGood:   pbGoodInfo,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

func (this *PlayerSystemShopFun) ShopExchangeRequest(head *pb.RpcHead, req *pb.ShopExchangeRequest) {
	uCode := this.shopExchange(head, req.ShopType, req.GoodsID)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ShopExchangeResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 兑换商店购买
func (this *PlayerSystemShopFun) shopExchange(head *pb.RpcHead, uShopType uint32, uGoodId uint32) cfgEnum.ErrorCode {
	cfgGood := cfgData.GetCfgExchangeShopConfig(uGoodId)
	if cfgGood == nil || cfgGood.ShopType != uShopType {
		return plog.Print(head.Id, cfgData.GetExchangeShopConfigErrorCode(uGoodId), cfgGood.ShopType, uShopType, uGoodId)
	}

	pShopInfo, ok := this.mapShopInfo[uShopType]
	if !ok {
		return plog.Print(head.Id, cfgEnum.ErrorCode_NoData, uShopType, uGoodId)
	}

	if _, ok := pShopInfo.mapBuyData[uGoodId]; ok {
		if pShopInfo.mapBuyData[uGoodId].Value >= cfgGood.LimitBuyNum {
			return plog.Print(head.Id, cfgEnum.ErrorCode_BuyTimesLimit, "shopExchangeShop", pShopInfo.mapBuyData[uGoodId], cfgGood.LimitBuyNum)
		}
	}

	//扣道具 有免费的
	if cfgGood.Currency != nil {
		errCode := this.getPlayerBagFun().DelItem(head, cfgGood.Currency.Kind, cfgGood.Currency.Id, cfgGood.Currency.Count, pb.EmDoingType_EDT_Shop)
		if errCode != cfgEnum.ErrorCode_Success {
			return plog.Print(head.Id, errCode, cfgGood.Currency)
		}
	}

	//存数据
	if _, ok := pShopInfo.mapBuyData[uGoodId]; ok {
		pShopInfo.mapBuyData[uGoodId].Value++
	} else {
		pShopInfo.mapBuyData[uGoodId] = &pb.PBU32U32{Key: uGoodId, Value: 1}
		pShopInfo.Items = append(pShopInfo.Items, pShopInfo.mapBuyData[uGoodId])
	}

	//如果免费的 去掉红点
	if cfgGood.Currency == nil && pShopInfo.mapBuyData[uGoodId].Value >= cfgGood.LimitBuyNum {
		pShopInfo.HaveRed = 0

		//通知客户端
		cluster.SendToClient(head, &pb.ShopRedNotify{
			PacketHead:  &pb.IPacket{},
			ShopRedList: []*pb.PBU32U32{&pb.PBU32U32{Key: uShopType, Value: 0}},
		}, cfgEnum.ErrorCode_Success)
	}

	this.UpdateSave(true)

	//加道具
	this.getPlayerBagFun().AddArrItem(head, cfgGood.Goods, pb.EmDoingType_EDT_Shop, true)

	cluster.SendToClient(head, &pb.ShopExchangeResponse{
		PacketHead: &pb.IPacket{},
		ShopType:   uShopType,
		GoodsID:    uGoodId,
		BuyTimes:   pShopInfo.mapBuyData[uGoodId].Value,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 商店刷新请求
func (this *PlayerSystemShopFun) ShopRefreshRequest(head *pb.RpcHead, req *pb.ShopRefreshRequest) {
	uCode := this.ShopRefresh(head, req.ShopType, req.IsFree)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ShopRefreshResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

func (this *PlayerSystemShopFun) UpdateShopFreeRefreshTime() {

}

// 商店刷新请求
func (this *PlayerSystemShopFun) ShopRefresh(head *pb.RpcHead, ShopType uint32, IsFree bool) cfgEnum.ErrorCode {
	if ShopType != uint32(pb.EmShopType_EST_BlackShop) {
		return plog.Print(head.Id, cfgEnum.ErrorCode_ShopTypeNotSupported, ShopType, IsFree)
	}

	cfgShop := cfgData.GetCfgShopConfig(ShopType)
	if cfgShop == nil {
		return plog.Print(head.Id, cfgData.GetShopConfigErrorCode(ShopType), ShopType, IsFree)
	}

	if IsFree {
		if this.pbBlackShop.RefreshInfo.DailyFreeUseCount >= this.pbBlackShop.RefreshInfo.DailyFreeMaxCount {
			return plog.Print(head.Id, cfgEnum.ErrorCode_NeedRefreshCount, ShopType, IsFree)
		}

		this.pbBlackShop.RefreshInfo.DailyFreeUseCount++

		if this.pbBlackShop.RefreshInfo.DailyFreeUseCount > 0 && this.pbBlackShop.RefreshInfo.NextFreeRefreshTime == 0 {
			this.pbBlackShop.RefreshInfo.NextFreeRefreshTime = base.GetNow() + uint64(cfgShop.RefreshRecoverTime)
		}

	} else {
		cfgRefresh := cfgData.GetCfgShopRefresh(ShopType, this.pbBlackShop.RefreshInfo.DailyBuyCount+1)
		if cfgRefresh == nil {
			return plog.Print(head.Id, cfgData.GetShopRefreshErrorCode(ShopType), ShopType, this.pbBlackShop.RefreshInfo.DailyBuyCount)
		}

		uCode := this.getPlayerBagFun().DelItem(head, cfgRefresh.DelItem.Kind, cfgRefresh.DelItem.Id, cfgRefresh.DelItem.Count, pb.EmDoingType_EDT_Shop)
		if uCode != cfgEnum.ErrorCode_Success {
			return plog.Print(head.Id, uCode, ShopType, cfgRefresh.DelItem)
		}

		this.pbBlackShop.RefreshInfo.DailyBuyCount++
	}

	this.refreshBlackShop(false)

	this.UpdateSave(true)
	cluster.SendToClient(head, &pb.ShopRefreshResponse{
		PacketHead:  &pb.IPacket{},
		ShopType:    ShopType,
		RefreshInfo: this.pbBlackShop.RefreshInfo,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 打开商店请求
func (this *PlayerSystemShopFun) ShopOpenRequest(head *pb.RpcHead, req *pb.ShopOpenRequest) {
	uCode := this.ShopOpen(head, req.ShopType)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ShopOpenResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 打开商店请求
func (this *PlayerSystemShopFun) ShopOpen(head *pb.RpcHead, ShopType uint32) cfgEnum.ErrorCode {
	cfgShop := cfgData.GetCfgShopConfig(ShopType)
	if cfgShop == nil {
		return plog.Print(this.AccountId, cfgData.GetShopConfigErrorCode(ShopType), ShopType)
	}

	pShopInfo, ok := this.mapShopInfo[ShopType]
	if !ok {
		pShopInfo = &PlayerShopInfo{
			PBShopInfo: &pb.PBShopInfo{
				ShopType: ShopType,
			},
			mapBuyData: make(map[uint32]*pb.PBU32U32),
		}

		//return plog.Print(head.Id, cfgEnum.ErrorCode_NoData, ShopType)
	}

	pbResponse := &pb.ShopOpenResponse{
		ShopType:   ShopType,
		PacketHead: &pb.IPacket{},
	}
	listCfgGoods := cfgData.GetCfgExchangeShopListConfig(ShopType)
	for _, cfgGood := range listCfgGoods {
		uBuyCount := uint32(0)
		pbBuy, ok := pShopInfo.mapBuyData[cfgGood.Id]
		if ok {
			uBuyCount = pbBuy.Value
		}

		pbGoodCfg := &pb.PBShopGoodCfg{
			GoodsID:     cfgGood.Id,
			MaxTimes:    cfgGood.LimitBuyNum,
			BuyTimes:    uBuyCount,
			Discount:    cfgGood.Discount,
			AddItem:     serverCommon.ChangeCommonItem(cfgGood.Goods, pb.EmDoingType_EDT_Shop),
			ProductId:   cfgGood.ProductId,
			ProductName: cfgGood.ProductName,
			ValueTips:   cfgGood.ValueTips,
			SortTag:     cfgGood.SortTag,
		}

		//如果有装备 todo

		if cfgGood.Currency != nil {
			pbGoodCfg.NeedItem = &pb.PBAddItem{
				Kind:   cfgGood.Currency.Kind,
				Id:     cfgGood.Currency.Id,
				Count:  cfgGood.Currency.Count,
				Params: cfgGood.Currency.Params,
			}
		}
		if cfgGood.ProductId > 0 {
			cfgCharge := cfgData.GetCfgCharge(cfgGood.ProductId, this.getPlayerBaseFun().GetPlayerBase().PlatSystemType)
			if cfgCharge == nil {
				plog.Error("(this *PlayerSystemShopFun) ShopOpen no productid:%d", cfgGood.ProductId)
				continue
			}
			pbGoodCfg.ProductName = cfgCharge.ProductName
			pbGoodCfg.Price = cfgCharge.Price
		}

		pbResponse.GoodList = append(pbResponse.GoodList, pbGoodCfg)
	}

	this.UpdateSave(true)
	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 充值回调
func (this *PlayerSystemShopFun) OnChargeBuyProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode {
	if len(cfgCharge.Param) != 2 {
		return plog.Print(this.AccountId, cfgData.GetChargeErrorCode(cfgCharge.ProductID), cfgCharge.ProductID)
	}

	shopType := cfgCharge.Param[0]
	cfgShopGood := cfgData.GetCfgExchangeShopConfig(cfgCharge.Param[1])
	if cfgShopGood == nil || cfgShopGood.ShopType != shopType || cfgShopGood.ProductId != cfgCharge.ProductID {
		return plog.Print(this.AccountId, cfgData.GetExchangeShopConfigErrorCode(cfgCharge.Param[1]), cfgShopGood.ShopType, shopType, cfgShopGood.ProductId, cfgCharge.ProductID)
	}

	pShopInfo, ok := this.mapShopInfo[shopType]
	if !ok {
		pShopInfo = &PlayerShopInfo{
			PBShopInfo: &pb.PBShopInfo{
				ShopType: shopType,
			},
			mapBuyData: make(map[uint32]*pb.PBU32U32),
		}
	}

	_, okBuy := pShopInfo.mapBuyData[cfgShopGood.Id]
	//如果是钻石商店，首次购买给双倍
	if shopType == uint32(cfgEnum.EShopType_DiamonShop) && !okBuy {
		arrDiamondAdd := make([]*common.ItemInfo, 0)
		for _, v := range cfgShopGood.Goods {
			arrDiamondAdd = append(arrDiamondAdd, &common.ItemInfo{
				Kind:   v.Kind,
				Count:  v.Count * 2,
				Id:     v.Id,
				Params: v.Params,
			})
		}
		this.getPlayerBagFun().AddArrItem(head, arrDiamondAdd, pb.EmDoingType_EDT_Charge, true)
	} else {
		this.getPlayerBagFun().AddArrItem(head, cfgShopGood.Goods, pb.EmDoingType_EDT_Charge, true)
	}

	//存数据
	if !okBuy {
		pShopInfo.mapBuyData[cfgShopGood.Id] = &pb.PBU32U32{Key: cfgShopGood.Id, Value: 1}
		pShopInfo.PBShopInfo.Items = append(pShopInfo.PBShopInfo.Items, pShopInfo.mapBuyData[cfgShopGood.Id])
	} else {
		pShopInfo.mapBuyData[cfgShopGood.Id].Value++
	}

	this.mapShopInfo[shopType] = pShopInfo
	this.UpdateSave(true)

	//通知客户端
	cluster.SendToClient(head, &pb.ShopExchangeGoodNotify{
		PacketHead: &pb.IPacket{},
		ShopType:   shopType,
		GoodInfo:   pShopInfo.mapBuyData[cfgShopGood.Id],
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 是否能够购买
func (this *PlayerSystemShopFun) canBuyProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode {
	if len(cfgCharge.Param) != 2 {
		return plog.Print(this.AccountId, cfgData.GetChargeErrorCode(cfgCharge.ProductID), cfgCharge.ProductID)
	}

	shopType := cfgCharge.Param[0]
	cfgShopGood := cfgData.GetCfgExchangeShopConfig(cfgCharge.Param[1])
	if cfgShopGood == nil || cfgShopGood.ShopType != shopType || cfgShopGood.ProductId != cfgCharge.ProductID {
		return plog.Print(this.AccountId, cfgData.GetExchangeShopConfigErrorCode(cfgCharge.Param[1]), cfgShopGood.ShopType, shopType, cfgShopGood.ProductId, cfgCharge.ProductID)
	}

	//不限制
	if cfgShopGood.LimitBuyNum <= 0 {
		return cfgEnum.ErrorCode_Success
	}

	pShopInfo, ok := this.mapShopInfo[shopType]
	if !ok {
		pShopInfo = &PlayerShopInfo{
			PBShopInfo: &pb.PBShopInfo{
				ShopType: shopType,
			},
			mapBuyData: make(map[uint32]*pb.PBU32U32),
		}
	}

	_, okBuy := pShopInfo.mapBuyData[cfgShopGood.Id]
	if okBuy && pShopInfo.mapBuyData[cfgShopGood.Id].Value >= cfgShopGood.LimitBuyNum {
		return cfgEnum.ErrorCode_BuyTimesLimit
	}

	return cfgEnum.ErrorCode_Success
}
