package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	common2 "corps/common"
	"corps/common/cfgData"
	report2 "corps/common/report"
	"corps/common/serverCommon"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"

	"github.com/golang/protobuf/proto"
)

type (
	PlayerBagFun struct {
		PlayerFun
		mapPlayerItemFun map[uint32]IPlayerItemKindFun

		mapItem         map[uint32]*pb.PBItem //背包数据
		mapDailyItemBuy map[uint32]uint32     //道具购买
		mapChangeItem   map[uint32]int64      //道具改变通知
	}

	IPlayerItemKindFun interface {
		Init(pFun *PlayerBagFun, systemType cfgEnum.ESystemType)
		GetPbItem(itemId uint32, itemCount int64, emDoing pb.EmDoingType, params ...uint32) (arrPbItems []*pb.PBAddItemData)
		AddItem(head *pb.RpcHead, pbItem *pb.PBAddItemData) cfgEnum.ErrorCode
	}
)

func (this *PlayerBagFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.mapItem = make(map[uint32]*pb.PBItem)
	this.mapDailyItemBuy = make(map[uint32]uint32)
	this.mapChangeItem = make(map[uint32]int64)
	this.registerPlayerItemFun()

}

// 注册道具
func (this *PlayerBagFun) registerPlayerItemFun() {
	this.mapPlayerItemFun = make(map[uint32]IPlayerItemKindFun)
	for i := cfgEnum.ESystemType_None; i <= cfgEnum.ESystemType_CrystalRobot; i++ {
		switch i {
		case cfgEnum.ESystemType_Item:
			this.mapPlayerItemFun[uint32(i)] = new(PlayerItemKindItemFun)
		case cfgEnum.ESystemType_Equipment:
			this.mapPlayerItemFun[uint32(i)] = new(PlayerItemKindEquipmentFun)
		case cfgEnum.ESystemType_Hero:
			this.mapPlayerItemFun[uint32(i)] = new(PlayerItemKindHeroFun)
		case cfgEnum.ESystemType_LootGroup:
			this.mapPlayerItemFun[uint32(i)] = new(PlayerItemKindLootGroupFun)
		case cfgEnum.ESystemType_LootEquipment:
			this.mapPlayerItemFun[uint32(i)] = new(PlayerItemKindLootEquipmentFun)
		case cfgEnum.ESystemType_Head:
			this.mapPlayerItemFun[uint32(i)] = new(PlayerItemKindHeadFun)
		case cfgEnum.ESystemType_HeadIcon:
			this.mapPlayerItemFun[uint32(i)] = new(PlayerItemKindHeadIconFun)
		case cfgEnum.ESystemType_Crystal:
			this.mapPlayerItemFun[uint32(i)] = new(PlayerItemKindCrystalFun)
		case cfgEnum.ESystemType_CrystalRobot:
			this.mapPlayerItemFun[uint32(i)] = new(PlayerItemKindCrystalRobotFun)
		}

		if _, ok := this.mapPlayerItemFun[uint32(i)]; ok {
			this.mapPlayerItemFun[uint32(i)].Init(this, i)
		}
	}
}

func (this *PlayerBagFun) getPlayerItemFun(kind uint32) IPlayerItemKindFun {
	return this.mapPlayerItemFun[kind]
}

// 加载背包数据
func (this *PlayerBagFun) Load(pData []byte) {
	pbData := &pb.PBPlayerBag{}
	proto.Unmarshal(pData, pbData)

	this.loadData(pbData)
	this.UpdateSave(false)
}

func (this *PlayerBagFun) loadData(pbData *pb.PBPlayerBag) {
	this.mapItem = make(map[uint32]*pb.PBItem)
	for i := 0; i < len(pbData.ItemList); i++ {
		this.mapItem[pbData.ItemList[i].Id] = pbData.ItemList[i]
	}

	this.mapDailyItemBuy = make(map[uint32]uint32)
	for i := 0; i < len(pbData.DailyBuyItem); i++ {
		this.mapDailyItemBuy[pbData.DailyBuyItem[i].Key] = pbData.DailyBuyItem[i].Value
	}

	this.UpdateSave(true)
}

// 保存
func (this *PlayerBagFun) Save(bNow bool) {
	if !this.BSave {
		return
	}

	this.BSave = false

	pbData := &pb.PBPlayerBag{}
	this.SavePb(pbData)

	//通知db保存玩家数据
	buff, _ := proto.Marshal(pbData)
	cluster.SendToDb(&pb.RpcHead{Id: this.AccountId}, "DbPlayerMgr", "SavePlayerDB", this.PbType, buff, bNow)
}

// 保存
func (this *PlayerBagFun) SavePb(pbData *pb.PBPlayerBag) {
	if pbData == nil {
		pbData = &pb.PBPlayerBag{}
	}

	for _, v := range this.mapItem {
		pbData.ItemList = append(pbData.ItemList, v)
	}

	for k, v := range this.mapDailyItemBuy {
		pbData.DailyBuyItem = append(pbData.DailyBuyItem, &pb.PBU32U32{Key: k, Value: v})
	}
}

func (this *PlayerBagFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	pbData.Bag = new(pb.PBPlayerBag)
	this.SavePb(pbData.Bag)
}
func (this *PlayerBagFun) LoadComplete() {
	this.mapChangeItem = make(map[uint32]int64)
}
func (this *PlayerBagFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerBag{}
}

// 设置玩家数据
func (this *PlayerBagFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}

	pbSystem := pbData.(*pb.PBPlayerBag)
	if pbSystem == nil {
		return false
	}

	this.loadData(pbSystem)
	return true
}

// 新玩家 需要初始化数据
func (this *PlayerBagFun) NewPlayer() {
	//送初始道具
	cfgGlobal := cfgData.GetCfgGlobalArrConfig(cfgEnum.GlobalArrConfig_GLOBAL_CFG_NewAddItem)
	if cfgGlobal != nil && len(cfgGlobal.Value) > 0 {
		arrItem := make([]*common2.ItemInfo, 0)
		for i := 0; i < len(cfgGlobal.Value)/2; i++ {
			arrItem = append(arrItem, &common2.ItemInfo{Kind: uint32(cfgEnum.ESystemType_Item), Id: cfgGlobal.Value[i*2], Count: int64(cfgGlobal.Value[i*2+1])})
		}

		this.AddArrItem(&pb.RpcHead{Id: this.AccountId}, arrItem, pb.EmDoingType_EDT_Other, false)
	}

	this.UpdateSave(true)
}

// 跨天
func (this *PlayerBagFun) PassDay(isDay, isWeek, isMonth bool) {
	if len(this.mapDailyItemBuy) > 0 {
		this.mapDailyItemBuy = make(map[uint32]uint32)
		this.UpdateSave(true)
	}
}

// 加道具
func (this *PlayerBagFun) GetItemCount(itemKind uint32, itemId uint32) int64 {
	uCount := int64(0)

	//判断类型
	switch cfgEnum.ESystemType(itemKind) {
	case cfgEnum.ESystemType_Item:
		uCount = this.getBagItemCount(itemId)
		break
	}

	return uCount
}

// 加道具 itemCount>0加道具 itemCount<0扣道具
func (this *PlayerBagFun) AddItem(head *pb.RpcHead, itemKind uint32, itemId uint32, itemCount int64, emDoing pb.EmDoingType, bNotice bool, params ...uint32) (uErrorCode cfgEnum.ErrorCode) {
	itemInfo := &common2.ItemInfo{
		Kind:   itemKind,
		Id:     itemId,
		Count:  itemCount,
		Params: params,
	}
	return this.AddArrItem(head, []*common2.ItemInfo{itemInfo}, emDoing, bNotice)
}

func (this *PlayerBagFun) GetPbItems(arrItem []*common2.ItemInfo, emDoing pb.EmDoingType) (arrPbItems []*pb.PBAddItemData) {
	arrPbItems = make([]*pb.PBAddItemData, 0)
	for _, v := range arrItem {
		itemFun := this.getPlayerItemFun(v.Kind)
		if itemFun == nil {
			plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, v.Kind, v.Id, v.Count, emDoing, v.Params)
			continue
		}
		pbItems := itemFun.GetPbItem(v.Id, v.Count, emDoing, v.Params...)
		if len(pbItems) <= 0 {
			//plog.Print(this.AccountId, cfgEnum.ErrorCode_Cfg, v.Id, v.Count, emDoing, v.Params)
			//serverCommon.CHECK_CODE(this.AccountId, cfgEnum.ErrorCode_Cfg, v.Id, v.Count, emDoing, v.Params)
			continue
		}

		arrPbItems = append(arrPbItems, pbItems...)
	}

	return this.MergePbItems(arrPbItems, emDoing)
}
func (this *PlayerBagFun) AddOneArrItem(head *pb.RpcHead, arrItem *common2.ItemInfo, emDoing pb.EmDoingType, bNotice bool) cfgEnum.ErrorCode {
	if arrItem == nil || arrItem.Id == 0 {
		return plog.Print(this.AccountId, cfgData.GetItemErrorCode(0), emDoing)
	}
	return this.AddArrItem(head, []*common2.ItemInfo{arrItem}, emDoing, bNotice)
}

// 增加道具数组
func (this *PlayerBagFun) AddArrItem(head *pb.RpcHead, arrItem []*common2.ItemInfo, emDoing pb.EmDoingType, bNotice bool) cfgEnum.ErrorCode {
	uErrorCode := cfgEnum.ErrorCode(cfgEnum.ErrorCode_Success)
	if len(arrItem) <= 0 {
		//plog.Info("AddArrItem arrItem is 0, emDoing:%d", emDoing)
		return uErrorCode
	}

	arrPbItems := this.GetPbItems(arrItem, emDoing)
	if len(arrPbItems) <= 0 {
		return plog.Print(this.AccountId, cfgData.GetItemErrorCode(0), *arrItem[0])
	}

	uErrorCode = this.AddPbItems(head, arrPbItems, emDoing, bNotice)
	if uErrorCode != cfgEnum.ErrorCode_Success {
		return uErrorCode
	}

	return uErrorCode
}
func (this *PlayerBagFun) MergePbItems(arrPbItems []*pb.PBAddItemData, emDoing pb.EmDoingType) []*pb.PBAddItemData {
	arrReturn := make([]*pb.PBAddItemData, 0)

	//英雄强制展开
	uLen := len(arrPbItems)
	for i := 0; i < uLen; i++ {
		if arrPbItems[i].Count <= 1 || (arrPbItems[i].Kind != uint32(cfgEnum.ESystemType_Hero) &&
			arrPbItems[i].Kind != uint32(cfgEnum.ESystemType_Crystal) &&
			arrPbItems[i].Kind != uint32(cfgEnum.ESystemType_CrystalRobot)) {
			arrReturn = append(arrReturn, &pb.PBAddItemData{
				Id:        arrPbItems[i].Id,
				Kind:      arrPbItems[i].Kind,
				Count:     arrPbItems[i].Count,
				DoingType: arrPbItems[i].DoingType,
				Params:    arrPbItems[i].Params,
				Equipment: arrPbItems[i].Equipment,
			})
			continue
		}

		for j := int64(1); j < arrPbItems[i].Count; j++ {
			arrReturn = append(arrReturn, &pb.PBAddItemData{
				Id:        arrPbItems[i].Id,
				Kind:      arrPbItems[i].Kind,
				Count:     1,
				DoingType: arrPbItems[i].DoingType,
				Params:    arrPbItems[i].Params,
				Equipment: arrPbItems[i].Equipment,
			})
		}

		arrReturn = append(arrReturn, &pb.PBAddItemData{
			Id:        arrPbItems[i].Id,
			Kind:      arrPbItems[i].Kind,
			DoingType: arrPbItems[i].DoingType,
			Count:     1,
			Params:    arrPbItems[i].Params,
			Equipment: arrPbItems[i].Equipment,
		})
	}

	//多余10个合并
	if emDoing == pb.EmDoingType_EDT_BoxOpen || emDoing == pb.EmDoingType_EDT_Draw || emDoing == pb.EmDoingType_EDT_StarSourceDraw || emDoing == pb.EmDoingType_EDT_EquipSplit {
		if len(arrReturn) <= 20 {
			return arrReturn
		}
	}

	//计算叠加
	arrRealReturn := make([]*pb.PBAddItemData, 0)
	mapKindItem := make(map[uint32]*pb.PBAddItemData)
	for _, info := range arrReturn {
		if info.Kind != uint32(cfgEnum.ESystemType_Item) {
			arrRealReturn = append(arrRealReturn, info)
			continue
		}
		if info.DoingType == pb.EmDoingType_EDT_Entry {
			arrRealReturn = append(arrRealReturn, info)
			continue
		}

		if _, ok := mapKindItem[info.Id]; !ok {
			mapKindItem[info.Id] = info
		} else {
			mapKindItem[info.Id].Count += info.Count
		}
	}

	for _, info := range mapKindItem {
		arrRealReturn = append(arrRealReturn, info)
	}
	return arrRealReturn
}

// 恭喜获得通知 必须StarCommonPrizeNotify 宝箱不合并 英雄强制展开
func (this *PlayerBagFun) CommonPrizeNotify(head *pb.RpcHead, arrPbItems []*pb.PBAddItemData, emDoing pb.EmDoingType) {
	if arrPbItems == nil {
		return
	}

	pbResponse := &pb.CommonPrizeNotify{
		PacketHead: &pb.IPacket{},
		DoingType:  emDoing,
		ItemInfo:   this.MergePbItems(arrPbItems, emDoing),
	}

	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)
}

func (this *PlayerBagFun) AddPbItems(head *pb.RpcHead, arrItem []*pb.PBAddItemData, emDoing pb.EmDoingType, bNotice bool) (uErrorCode cfgEnum.ErrorCode) {
	if len(arrItem) <= 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_Success, emDoing)
	}

	uCode := cfgEnum.ErrorCode_Success
	for _, pbItem := range arrItem {
		uCode = this.getPlayerItemFun(pbItem.Kind).AddItem(head, pbItem)
		if uCode != cfgEnum.ErrorCode_Success {
			return plog.Print(this.AccountId, uCode, pbItem)
		}
	}

	// 恭喜获得
	if bNotice {
		this.CommonPrizeNotify(head, arrItem, emDoing)
	}

	return uCode
}

// 加道具 bForce为true表示 道具不够扣成负数
func (this *PlayerBagFun) DelItem(head *pb.RpcHead, itemKind uint32, itemId uint32, itemCount int64, emDoing pb.EmDoingType) cfgEnum.ErrorCode {
	return this.AddItem(head, itemKind, itemId, itemCount*-1, emDoing, false)
}
func (this *PlayerBagFun) DelCommonItem(head *pb.RpcHead, itemInfo *common2.ItemInfo, emDoing pb.EmDoingType) cfgEnum.ErrorCode {
	return this.AddItem(head, itemInfo.Kind, itemInfo.Id, itemInfo.Count*-1, emDoing, false)
}

// 扣除道具
func (this *PlayerBagFun) DelArrItem(head *pb.RpcHead, arrItem []*common2.ItemInfo, emDoing pb.EmDoingType) cfgEnum.ErrorCode {
	//判断道具是否足够
	for _, v := range arrItem {
		if this.GetItemCount(v.Kind, v.Id) < v.Count {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_ItemNotEnough, v.Kind, v.Id, v.Count)
		}
	}

	for _, v := range arrItem {
		uErrorCode := this.AddItem(head, v.Kind, v.Id, v.Count*-1, emDoing, false)
		if uErrorCode != cfgEnum.ErrorCode_Success {
			return plog.Print(this.AccountId, uErrorCode, v.Kind, v.Id, v.Count, emDoing)
		}
	}

	return cfgEnum.ErrorCode_Success

}

// 加道具
func (this *PlayerBagFun) getBagItemCount(itemId uint32) int64 {
	uCount := int64(0)
	itemData, ok := this.mapItem[itemId]
	if !ok {
		return uCount
	}

	return itemData.Count
}

// 扣除道具
func (this *PlayerBagFun) delBagItem(head *pb.RpcHead, itemId uint32, itemCount int64, emDoing pb.EmDoingType) cfgEnum.ErrorCode {
	if itemCount < 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ItemCountIsNegative, itemId, itemCount, emDoing)
	}

	if this.getBagItemCount(itemId) < itemCount {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ItemNotEnough, itemId, itemCount, this.getBagItemCount(itemId), emDoing)
	}

	//查找id ，先从最少得开始扣除
	pItem := this.getItem(itemId)
	if pItem == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ItemNotEnough, itemId, itemCount, emDoing)
	}

	pItem.Count -= itemCount
	this.updateItem(head, pItem, emDoing)

	//发送道具更新
	this.ItemsUpdateNotify(head, []*pb.PBItem{pItem})

	//数据上报
	report2.Send(head, &report2.ReportAddItem{
		Kind:   uint32(cfgEnum.ESystemType_Item),
		ItemID: itemId,
		Add:    itemCount,
		Total:  pItem.Count,
		Doing:  uint32(emDoing),
	})
	return cfgEnum.ErrorCode_Success
}

// 加道具
func (this *PlayerBagFun) addBagItem(head *pb.RpcHead, pbItem *pb.PBAddItemData) cfgEnum.ErrorCode {
	if pbItem.Count < 0 {
		return this.delBagItem(head, pbItem.Id, pbItem.Count*-1, pbItem.DoingType)
	}

	uErrorCode := cfgEnum.ErrorCode_Success

	cfgItem := cfgData.GetCfgItem(pbItem.Id)
	if cfgItem == nil {
		uErrorCode = cfgEnum.ErrorCode_ItemNotExist
		return uErrorCode
	}

	//自动使用
	if cfgItem.AutoUse > 0 {
		uCode := this.use(head, cfgItem, pbItem)
		return uCode
	}

	//是否新增
	pItem := this.getItem(pbItem.Id)
	if pItem == nil {
		pItem = this.newItem(pbItem.Id, pbItem.Count, pbItem.DoingType)
	} else {
		pItem.Count += pbItem.Count
	}

	this.updateItem(head, pItem, pbItem.DoingType)

	//发送消息
	this.ItemsUpdateNotify(head, []*pb.PBItem{pItem})

	//数据上报
	report2.Send(head, &report2.ReportAddItem{
		Kind:   uint32(cfgEnum.ESystemType_Item),
		ItemID: pbItem.Id,
		Add:    pbItem.Count,
		Total:  pItem.Count,
		Doing:  uint32(pbItem.DoingType),
	})
	return uErrorCode
}

// 补满可以叠加的道具
func (this *PlayerBagFun) addItemLog(itemId uint32, oldCount int64, newCount int64, emDoing pb.EmDoingType) {
	plog.Info("addItemLog itemId:%d oldCount:%d newCount:%d emDoing:%d", itemId, oldCount, newCount, emDoing)

	//存流水日志
	/*
		packet.SendRecordMsg(this.AccountId, serverCommon.EmRecordTypeAddItem, &serverCommon.RecordAddItem{
			Time:      base.GetNow(),
			AccountId: this.AccountId,
			Id:        itemId,
			Add:       oldCount - newCount,
			Total:     newCount,
			Doing:     uint32(emDoing),
		})
	*/

	this.BSave = true
}

// 补满可以叠加的道具
func (this *PlayerBagFun) newItem(itemId uint32, itemCount int64, emDoing pb.EmDoingType) *pb.PBItem {
	itemInfo := &pb.PBItem{
		Id:    itemId,
		Count: itemCount,
	}

	//写日志
	this.addItemLog(itemId, 0, int64(itemCount), emDoing)

	//存数据
	this.mapItem[itemId] = itemInfo
	return itemInfo
}

// 补满可以叠加的道具
func (this *PlayerBagFun) realDelItem(uId uint32) {
	_, ok := this.mapItem[uId]
	if ok {
		delete(this.mapItem, uId)
	}
}

// 设置玩家道具 itemCount剩余大局数量
func (this *PlayerBagFun) UpdatePlayerItem(head *pb.RpcHead, itemId uint32, itemCount int64, emType pb.EmDoingType) {
	pPlayerItem := this.getItem(itemId)
	if pPlayerItem == nil {
		return
	}

	pPlayerItem.Count = itemCount
	this.updateItem(head, pPlayerItem, emType)
}

// 使用道具请求
func (this *PlayerBagFun) ItemUseRequest(head *pb.RpcHead, pbRequest *pb.ItemUseRequest) {
	uCode := this.ItemUse(head, pbRequest.Id, pbRequest.Count)
	cluster.SendToClient(head, &pb.ItemUseResponse{
		PacketHead: &pb.IPacket{},
		Id:         pbRequest.Id,
		Count:      pbRequest.Count}, uCode)
}
func (this *PlayerBagFun) use(head *pb.RpcHead, cfgItem *cfgData.ItemCfg, pbItem *pb.PBAddItemData) cfgEnum.ErrorCode {
	bNotice := pbItem.DoingType == pb.EmDoingType_EDT_ItemUse
	emErrorCode := cfgEnum.ErrorCode_Success
	switch cfgEnum.EHydraItemUseType(cfgItem.UseType) {
	case cfgEnum.EHydraItemUseType_AddItem:
		{
			if len(cfgItem.UseParam) != 3 {
				return plog.Print(this.AccountId, cfgEnum.ErrorCode_CfgItemParam, cfgItem.UseParam, cfgItem.Id, pbItem.Count)
			}

			//增加等级
			emErrorCode = this.AddItem(head, cfgItem.UseParam[0], cfgItem.UseParam[1], int64(cfgItem.UseParam[2])*pbItem.Count, pb.EmDoingType_EDT_ItemUse, bNotice)
			if emErrorCode != cfgEnum.ErrorCode_Success {
				return plog.Print(this.AccountId, emErrorCode, cfgItem.Id, cfgItem.UseParam, pbItem.Count)
			}
		}
	case cfgEnum.EHydraItemUseType_GoldBox:
		{
			if len(cfgItem.UseParam) != 2 {
				return plog.Print(this.AccountId, cfgEnum.ErrorCode_CfgItemParam, cfgItem.Id, cfgItem.UseParam, pbItem.Count)
			}

			// 走金币膨胀系数
			mapTech := this.getPlayerSystemHookTechFun().GetHookTechEffect(cfgEnum.TechEffectType_AddBoxGoldRate)
			uBoxGoldGrowthRate := uint32(0)
			if len(mapTech) > 0 {
				uBoxGoldGrowthRate = uint32(mapTech[0])
			}

			// 宝箱金币膨胀
			arrItem := []*common2.ItemInfo{&common2.ItemInfo{
				Kind:  uint32(cfgEnum.ESystemType_LootGroup),
				Id:    cfgItem.UseParam[0],
				Count: int64(cfgItem.UseParam[1]) * pbItem.Count,
			}}

			arrPbItem := this.getPlayerBagFun().GetPbItems(arrItem, pb.EmDoingType_EDT_ItemUse)
			for _, item := range arrPbItem {
				if item.Kind != uint32(cfgEnum.ESystemType_Item) || item.Id != uint32(pb.EmItemExpendType_EIET_Gold) {
					continue
				}
				item.Count += item.Count * int64(uBoxGoldGrowthRate) / base.MIL_PERCENT
			}

			emErrorCode = this.AddPbItems(head, arrPbItem, pb.EmDoingType_EDT_ItemUse, bNotice)
			if emErrorCode != cfgEnum.ErrorCode_Success {
				return plog.Print(this.AccountId, emErrorCode, cfgItem.Id, cfgItem.UseParam, pbItem.Count)
			}
		}
	case cfgEnum.EHydraItemUseType_Loot:
		{
			if len(cfgItem.UseParam) != 1 {
				return plog.Print(this.AccountId, cfgEnum.ErrorCode_CfgItemParam, cfgItem.Id, cfgItem.UseParam, pbItem.Count)
			}

			arrItem := []*common2.ItemInfo{&common2.ItemInfo{
				Kind:  uint32(cfgEnum.ESystemType_LootGroup),
				Id:    cfgItem.UseParam[0],
				Count: pbItem.Count,
			}}

			arrPbItem := this.getPlayerBagFun().GetPbItems(arrItem, pb.EmDoingType_EDT_ItemUse)
			if len(arrPbItem) <= 0 {
				return plog.Print(this.AccountId, emErrorCode, cfgItem.Id, cfgItem.UseParam, pbItem.Count)
			}
			emErrorCode = this.AddPbItems(head, arrPbItem, pb.EmDoingType_EDT_ItemUse, bNotice)
			if emErrorCode != cfgEnum.ErrorCode_Success {
				return plog.Print(this.AccountId, emErrorCode, cfgItem.Id, cfgItem.UseParam, pbItem.Count)
			}

			pbItem.Id = arrPbItem[0].Id
			pbItem.Kind = arrPbItem[0].Kind
			pbItem.Count = arrPbItem[0].Count
			pbItem.Params = arrPbItem[0].Params
		}
	default:
		return plog.Print(this.AccountId, cfgData.GetItemErrorCode(cfgItem.Id), cfgItem.UseType, cfgItem.Id, pbItem.Count)
	}
	return emErrorCode
}

// 使用道具
func (this *PlayerBagFun) ItemUse(head *pb.RpcHead, uId uint32, uCount uint32) cfgEnum.ErrorCode {
	if uCount < 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ItemCountIsNegative, uId, uCount)
	}

	pbItem := this.getItem(uId)
	if pbItem == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ItemNotEnough, uId, uCount)
	}

	//修正道具个数
	if pbItem.Count < int64(uCount) {
		uCount = uint32(pbItem.Count)
	}

	cfgItem := cfgData.GetCfgItem(pbItem.Id)
	if cfgItem == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ItemNotExist, *pbItem)
	}

	//使用个数限制
	if cfgItem.MaxUse > 0 && uCount >= cfgItem.MaxUse {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ItemMaxUse, uId, cfgItem.MaxUse, uCount)
	}

	uCode := this.use(head, cfgItem, &pb.PBAddItemData{Kind: uint32(cfgEnum.ESystemType_Item), Id: pbItem.Id, Count: int64(uCount), DoingType: pb.EmDoingType_EDT_ItemUse})
	if uCode != cfgEnum.ErrorCode_Success {
		return uCode
	}

	//扣除道具
	pbItem = this.getItem(uId)
	if pbItem == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ItemNotEnough, cfgItem.UseType, uId, uCount)
	}

	if int64(uCount) <= pbItem.Count {
		pbItem.Count -= int64(uCount)
	} else {
		pbItem.Count = 0
	}

	this.updateItem(head, pbItem, pb.EmDoingType_EDT_ItemUse)

	//更新给客户端
	this.ItemUpdateNotify(head, pbItem)

	return cfgEnum.ErrorCode_Success
}

func (this *PlayerBagFun) getItem(uId uint32) *pb.PBItem {
	info, ok := this.mapItem[uId]
	if !ok {
		return nil
	}

	return info
}

// 道具使用预览请求
func (this *PlayerBagFun) ItemUseShowRequest(head *pb.RpcHead, pbRequest *pb.ItemUseShowRequest) {
	uCode := this.ItemUseShow(head, pbRequest.Id)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ItemUseShowResponse{
			PacketHead: &pb.IPacket{}}, uCode)
	}
}

// 道具使用预览请求
func (this *PlayerBagFun) ItemUseShow(head *pb.RpcHead, uId uint32) cfgEnum.ErrorCode {
	pbItem := this.getItem(uId)
	if pbItem == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_BagNotFoundItem, uId)
	}

	cfgItem := cfgData.GetCfgItem(pbItem.Id)
	if cfgItem == nil {
		return plog.Print(this.AccountId, cfgData.GetItemErrorCode(pbItem.Id), pbItem.Id)
	}

	pbResponse := &pb.ItemUseShowResponse{PacketHead: &pb.IPacket{}}
	if cfgItem.UseType == uint32(cfgEnum.EHydraItemUseType_Loot) {
		if len(cfgItem.UseParam) != 1 {
			return plog.Print(this.AccountId, cfgData.GetItemErrorCode(pbItem.Id), pbItem.Id, cfgItem.UseParam)
		}

	} else {
		if len(cfgItem.UseParam) != 1 {
			return plog.Print(this.AccountId, cfgData.GetItemErrorCode(pbItem.Id), pbItem.Id, cfgItem.UseParam)
		}

		uGroupId := cfgItem.UseParam[0]
		listCfg := cfgData.GetCfgItemSelectGroup(uGroupId)
		if len(listCfg) < 0 {
			return plog.Print(this.AccountId, cfgData.GetItemSelectErrorCode(uGroupId), pbItem.Id, cfgItem.UseParam)
		}

		for _, cfg := range listCfg {
			pbResponse.ItemList = append(pbResponse.ItemList, &pb.ItemUseShowInfo{
				Id: cfg.Id,
				Item: &pb.PBAddItem{
					Id:     cfg.AddPrize.Id,
					Kind:   cfg.AddPrize.Kind,
					Count:  cfg.AddPrize.Count,
					Params: cfg.AddPrize.Params,
				}})
		}
	}

	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 道具选择请求
func (this *PlayerBagFun) ItemSelectRequest(head *pb.RpcHead, pbRequest *pb.ItemSelectRequest) {
	uCode := this.ItemSelect(head, pbRequest.Id, pbRequest.SelectList)
	cluster.SendToClient(head, &pb.ItemSelectResponse{
		PacketHead: &pb.IPacket{},
		Id:         pbRequest.Id,
		SelectList: pbRequest.SelectList}, uCode)
}

// 道具选择请求
func (this *PlayerBagFun) ItemSelect(head *pb.RpcHead, uId uint32, selectList []*pb.PBU32U32) cfgEnum.ErrorCode {
	cfgItem := cfgData.GetCfgItem(uId)
	if cfgItem == nil {
		return plog.Print(this.AccountId, cfgData.GetItemErrorCode(uId), uId)
	}

	if len(cfgItem.UseParam) != 1 {
		return plog.Print(this.AccountId, cfgData.GetItemErrorCode(uId), uId, cfgItem.UseParam)
	}

	uGroupId := cfgItem.UseParam[0]

	uItemCount := uint32(0)
	arrAddItem := make([]*common2.ItemInfo, 0)

	if cfgItem.UseType == uint32(cfgEnum.EHydraItemUseType_SelectFixItem) {
		for _, info := range selectList {
			uItemCount += info.Value

			cfgSelect := cfgData.GetCfgItemSelect(info.Key)
			if cfgSelect == nil {
				return plog.Print(this.AccountId, cfgData.GetItemSelectErrorCode(info.Key), uId, cfgItem.UseParam)
			}

			if cfgSelect.GroupId != uGroupId {
				return plog.Print(this.AccountId, cfgData.GetItemSelectErrorCode(cfgSelect.GroupId), uId, cfgItem.UseParam)
			}

			arrAddItem = append(arrAddItem, cfgSelect.AddPrize)
		}

	} else if cfgItem.UseType == uint32(cfgEnum.EHydraItemUseType_SelectAllItem) {
		listCfg := cfgData.GetCfgItemSelectGroup(uGroupId)
		if len(listCfg) < 0 {
			return plog.Print(this.AccountId, cfgData.GetItemSelectErrorCode(uGroupId), uId, cfgItem.UseParam)
		}

		uItemCount = 1
		for _, cfgSelect := range listCfg {
			arrAddItem = append(arrAddItem, cfgSelect.AddPrize)
		}
	} else {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_CfgItemUseTypeNotSupported, uId, cfgItem.UseParam)
	}

	if this.GetItemCount(uint32(cfgEnum.ESystemType_Item), uId) < int64(uItemCount) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ItemNotEnough, uId, cfgItem.UseParam)
	}

	//装备需要判断格子
	arrPbItems := this.GetPbItems(arrAddItem, pb.EmDoingType_EDT_ItemUse)
	uEquipCount, uHeroCount := serverCommon.GetPBItemEquipmentAndHeroCount(arrPbItems)
	if uEquipCount > this.GetPlayerEquipmentFun().GetSpareBag() {
		cluster.SendCommonToClient(head, cfgEnum.ErrorCode_BagFull)
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_BagFull, uId, cfgItem.UseParam)
	}

	if uHeroCount > 0 && this.getPlayerHeroFun().GetSpareBag() < uHeroCount {
		cluster.SendCommonToClient(head, cfgEnum.ErrorCode_HeroBagFull)
		return plog.Print(head.Id, cfgEnum.ErrorCode_HeroBagFull, uId, cfgItem.UseParam)
	}

	//口道具
	this.DelItem(head, uint32(cfgEnum.ESystemType_Item), uId, int64(uItemCount), pb.EmDoingType_EDT_ItemUse)
	this.AddPbItems(head, arrPbItems, pb.EmDoingType_EDT_ItemUse, true)

	return cfgEnum.ErrorCode_Success
}

// 同步给客户端
func (this *PlayerBagFun) ItemUpdateNotify(head *pb.RpcHead, pbItem *pb.PBItem) {
	this.ItemsUpdateNotify(head, []*pb.PBItem{pbItem})
}

// 同步给客户端
func (this *PlayerBagFun) ItemsUpdateNotify(head *pb.RpcHead, arrPlayerItem []*pb.PBItem) {
	for _, item := range arrPlayerItem {
		this.mapChangeItem[item.Id] += item.Count
	}

	//判断叠加个数
	cluster.SendToClient(head, &pb.ItemUpdateNotify{
		PacketHead: &pb.IPacket{},
		ItemList:   arrPlayerItem,
	}, cfgEnum.ErrorCode_Success)
}

// 同步给客户端
func (this *PlayerBagFun) SendChangeItem(head *pb.RpcHead) {
	return
	if len(this.mapChangeItem) <= 0 {
		return
	}

	arrPlayerItem := make([]*pb.PBItem, 0)
	for id, count := range this.mapChangeItem {
		if count == 0 {
			continue
		}

		arrPlayerItem = append(arrPlayerItem, &pb.PBItem{
			Id:    id,
			Count: this.getBagItemCount(id),
		})
	}

	this.mapChangeItem = make(map[uint32]int64)
	//判断叠加个数
	if len(arrPlayerItem) > 0 {
		cluster.SendToClient(head, &pb.ItemUpdateNotify{
			PacketHead: &pb.IPacket{},
			ItemList:   arrPlayerItem,
		}, cfgEnum.ErrorCode_Success)
	}
}

func (this *PlayerBagFun) Heat() {
	//this.SendChangeItem(&pb.RpcHead{Id: this.AccountId})
}

// 设置玩家道具
func (this *PlayerBagFun) updateItem(head *pb.RpcHead, pbItem *pb.PBItem, emType pb.EmDoingType) {
	oldCount := pbItem.Count
	this.mapItem[pbItem.Id] = pbItem
	if pbItem.Count == 0 {
		this.realDelItem(pbItem.Id)
	}

	//写日志
	this.addItemLog(pbItem.Id, int64(oldCount), int64(pbItem.Count), emType)
}

// 道具购买请求
func (this *PlayerBagFun) ItemBuyRequest(head *pb.RpcHead, pbRequest *pb.ItemBuyRequest) {
	uCode := this.ItemBuy(head, pbRequest.Id, pbRequest.Count)
	cluster.SendToClient(head, &pb.ItemBuyResponse{
		PacketHead: &pb.IPacket{},
		Id:         pbRequest.Id,
		Count:      pbRequest.Count,
	}, uCode)
}

// 道具购买请求
func (this *PlayerBagFun) ItemBuy(head *pb.RpcHead, uId uint32, uCount uint32) cfgEnum.ErrorCode {
	uCurCount, ok := this.mapDailyItemBuy[uId]
	if !ok {
		uCurCount = 0
	}

	if uCurCount+uCount > cfgData.GetCfgItemBuyMaxCount(uId) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_DailyMaxCount, uId)
	}

	arrNeedItem := cfgData.GetCfgItemBuyNeedItem(uId, uCurCount, uCount)
	if arrNeedItem.Count <= 0 {
		return plog.Print(this.AccountId, cfgData.GetItemBuyErrorCode(uId), uId)
	}

	//扣道具
	uCode := this.DelItem(head, arrNeedItem.Kind, arrNeedItem.Id, arrNeedItem.Count, pb.EmDoingType_EDT_ItemBuy)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode, uId, uCount)
	}

	this.mapDailyItemBuy[uId] = uCurCount + uCount
	this.AddItem(head, uint32(cfgEnum.ESystemType_Item), uId, int64(uCount), pb.EmDoingType_EDT_ItemBuy, true)

	this.UpdateSave(true)
	return cfgEnum.ErrorCode_Success
}
