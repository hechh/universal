package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/common/orm/redis"
	"corps/common/serverCommon"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"fmt"
	"math"

	"github.com/golang/protobuf/proto"
)

type (
	PlayerSystemChargeFun struct {
		PlayerFun
		pbData         *pb.PBCharge
		mapFirstCharge map[uint32]*pb.PBFirstCharge

		mapFun  map[cfgEnum.EChargeType]IPlayerSystemChargeFun
		bCharge bool
	}

	IPlayerSystemChargeFun interface {
		Init(pFun *PlayerSystemChargeFun)
		loadData(pbData *pb.PBPlayerSystemCharge)
		saveData(pbData *pb.PBPlayerSystemCharge)
		LoadComplete()
		PassDay(isDay, isWeek, isMonth bool)
		OnChargeBuyProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode
		canBuyProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode
	}
)

func (this *PlayerSystemChargeFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.pbData = &pb.PBCharge{}
	this.mapFirstCharge = make(map[uint32]*pb.PBFirstCharge)
	this.bCharge = false

	this.RegisterFun()
}

// 注册
func (this *PlayerSystemChargeFun) RegisterFun() {
	this.mapFun = make(map[cfgEnum.EChargeType]IPlayerSystemChargeFun)
	for i := cfgEnum.EChargeType_ShopCharge; i <= cfgEnum.EChargeType_ChargeGift; i++ {
		switch i {
		case cfgEnum.EChargeType_BP:
			this.mapFun[i] = new(PlayerSystemChargeBP)
			this.mapFun[i].Init(this)
		case cfgEnum.EChargeType_ChargeCard:
			this.mapFun[i] = new(PlayerSystemChargeCard)
			this.mapFun[i].Init(this)
		}
	}
}
func (this *PlayerSystemChargeFun) getChargeFun(emType cfgEnum.EChargeType) IPlayerSystemChargeFun {
	return this.mapFun[emType]
}
func (this *PlayerSystemChargeFun) GetChargeBPFun() *PlayerSystemChargeBP {
	return this.getChargeFun(cfgEnum.EChargeType_BP).(*PlayerSystemChargeBP)
}
func (this *PlayerSystemChargeFun) GetChargeCardFun() *PlayerSystemChargeCard {
	return this.getChargeFun(cfgEnum.EChargeType_ChargeCard).(*PlayerSystemChargeCard)
}

// 从数据库中加载
func (this *PlayerSystemChargeFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	if pbSystem.Charge == nil {
		pbSystem.Charge = &pb.PBPlayerSystemCharge{}
		return
	}

	this.loadData(pbSystem.Charge)

	this.UpdateSave(false)
}
func (this *PlayerSystemChargeFun) loadData(pbData *pb.PBPlayerSystemCharge) {
	if pbData.Charge == nil {
		pbData.Charge = &pb.PBCharge{}
	}

	this.pbData = pbData.Charge
	this.mapFirstCharge = make(map[uint32]*pb.PBFirstCharge)
	for _, info := range pbData.FirstChargeList {
		this.mapFirstCharge[info.FirstChargeId] = info
	}

	for _, info := range this.mapFun {
		info.loadData(pbData)
	}

	this.UpdateSave(true)
}

func (this *PlayerSystemChargeFun) saveData(pbData *pb.PBPlayerSystemCharge) {
	if pbData.Charge == nil {
		pbData.Charge = &pb.PBCharge{}
	}

	pbData.Charge = this.pbData
	for _, info := range this.mapFirstCharge {
		pbData.FirstChargeList = append(pbData.FirstChargeList, info)
	}

	for _, info := range this.mapFun {
		info.saveData(pbData)
	}
}

// 存储到数据库
func (this *PlayerSystemChargeFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	if pbSystem.Charge == nil {
		pbSystem.Charge = new(pb.PBPlayerSystemCharge)
	}
	this.saveData(pbSystem.Charge)

	return this.BSave
}
func (this *PlayerSystemChargeFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemCharge{}
}

// 设置玩家数据
func (this *PlayerSystemChargeFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}
	pbSystem := pbData.(*pb.PBPlayerSystemCharge)
	if pbSystem == nil {
		return false
	}
	this.loadData(pbSystem)
	this.UpdateSave(true)
	return true
}

func (this *PlayerSystemChargeFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = new(pb.PBPlayerSystem)
	}

	this.SaveSystem(pbData.System)
}

// 加载完成需要读取有没有充值数据
func (this *PlayerSystemChargeFun) LoadComplete() {
	for _, info := range this.mapFirstCharge {
		//数据兼容 改成注册时间
		if info.ActiveTime > 0 && info.OpenTime == 0 {
			info.OpenTime = this.getPlayerBaseFun().GetRegTime()
			//通知客户端
			cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.FirstChargeNotify{
				PacketHead:      &pb.IPacket{},
				FirstChargeInfo: info,
			}, cfgEnum.ErrorCode_Success)
		}
	}

	this.UpdatePlayerCharge(&pb.RpcHead{Id: this.AccountId})

	for _, info := range this.mapFun {
		info.LoadComplete()
	}
}

// 心跳处理
func (this *PlayerSystemChargeFun) HeartbeatRequest(head *pb.RpcHead) {
	//需要同步充值数据
	if this.bCharge {
		this.UpdatePlayerCharge(head)
	}

}

// 充值订单请求
func (this *PlayerSystemChargeFun) ChargeOrderRequest(head *pb.RpcHead, pbRequest *pb.ChargeOrderRequest) {
	uCode := this.ChargeOrder(head, pbRequest.ProductId, pbRequest.IsNeigou)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ChargeOrderResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 充值订单请求
func (this *PlayerSystemChargeFun) ChargeOrder(head *pb.RpcHead, uProductId uint32, bNeigou bool) cfgEnum.ErrorCode {
	cfgCharge := cfgData.GetCfgCharge(uProductId, this.getPlayerBaseFun().GetPlayerBase().PlatSystemType)
	if cfgCharge == nil {
		return plog.Print(this.AccountId, cfgData.GetChargeErrorCode(uProductId), uProductId)
	}

	//判断是否能购买
	uCode := cfgEnum.ErrorCode_Success
	switch cfgEnum.EChargeType(cfgCharge.ChargeType) {
	case cfgEnum.EChargeType_SevenDay:
		uCode = this.getPlayerSystemSevenDayFun().canBuyProduct(head, cfgCharge)
	case cfgEnum.EChargeType_BP, cfgEnum.EChargeType_ChargeCard:
		this.getChargeFun(cfgEnum.EChargeType(cfgCharge.ChargeType)).canBuyProduct(head, cfgCharge)
	case cfgEnum.EChargeType_ChargeGift:
		uCode = this.getPlayerActivityChargeGiftFun().canBuyProduct(head, cfgCharge)
	case cfgEnum.EChargeType_ShopCharge:
		uCode = this.getPlayerSystemShopFun().canBuyProduct(head, cfgCharge)
	case cfgEnum.EChargeType_OpenServerGift:
		uCode = this.getPlayerActivityOpenServerGiftFun().canBuyProduct(head, cfgCharge)
	case cfgEnum.EChargeType_FirstCharge:
		uCode = this.canBuyFirstChargeProduct(head, cfgCharge)
	}

	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode, uProductId)
	}

	this.pbData.OrderId++
	OrderNo := fmt.Sprintf("%d_%d", this.pbData.OrderId, base.GetNow())

	//向数据库插入
	if !serverCommon.InsertChargeOrder(this.AccountId, uProductId, cfgCharge.PlatSystemType, cfgCharge.Price, OrderNo, base.GetNow()) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_UpdateTOrderWrong, uProductId)
	}

	// 成就
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ChargeCount, cfgCharge.Price)

	this.UpdateSave(true)

	//通知客户端
	cluster.SendToClient(head, &pb.ChargeOrderResponse{
		PacketHead:     &pb.IPacket{},
		BingchuanOrder: this.getBingchuanOrder(head, cfgCharge, OrderNo, bNeigou),
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 更新玩家充值
func (this *PlayerSystemChargeFun) getBingchuanOrder(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg, orderNo string, isNeigou bool) *pb.PBChargeBingchuanOrder {
	pbBingchuan := &pb.PBChargeBingchuanOrder{}

	//商品ID * 单价（分）* 数量*商品名称 示例 8*600*1*6元充值 注意如果是IOS或GooglePlay内购时要多加上内购定义商品ID
	pbBingchuan.OrderItem = fmt.Sprintf("%d*%d*%d*%s", cfgCharge.ProductID, cfgCharge.Price, 1, cfgCharge.ProductName)
	if isNeigou {
		pbBingchuan.OrderItem = fmt.Sprintf("%s*%s", pbBingchuan.OrderItem, cfgCharge.PlatProductID)
	}

	pbBingchuan.OrderNo = orderNo
	if _, remainder := math.Modf(float64(cfgCharge.PayNum)); remainder == 0 {
		pbBingchuan.PayNum = fmt.Sprintf("%d", uint32(cfgCharge.PayNum))
	} else {
		pbBingchuan.PayNum = fmt.Sprintf("%.2f", cfgCharge.PayNum)
	}

	pbBingchuan.UserId = fmt.Sprintf("%s", this.getPlayerBaseFun().GetPlayerBase().AccountName)
	pbBingchuan.ActorId = fmt.Sprintf("%d", this.AccountId)

	//MD5(orderItem * orderNo * payNum * userID * buyKey) 国内版不用加入CurrencyType，海外版必加
	pbBingchuan.OrderSign = fmt.Sprintf("%s%s%s%s", pbBingchuan.OrderItem, pbBingchuan.OrderNo, pbBingchuan.PayNum, pbBingchuan.UserId)
	if cfgCharge.PlatSystemType == uint32(cfgEnum.EPlatSystemType_O_Andriod) || cfgCharge.PlatSystemType == uint32(cfgEnum.EPlatSystemType_O_Ios) {
		pbBingchuan.CurrencyType = cfgCharge.CurrencyType
		pbBingchuan.OrderSign = fmt.Sprintf("%s%s", pbBingchuan.OrderSign, pbBingchuan.CurrencyType)
	}
	strBuyKey := serverCommon.PlatformConfig.PayBuyKey
	pbBingchuan.OrderSign = fmt.Sprintf("%s%s", pbBingchuan.OrderSign, strBuyKey)
	pbBingchuan.OrderSign = base.MD5(pbBingchuan.OrderSign)

	//透传参数 商品ID
	pbBingchuan.DeveloperPayload = fmt.Sprintf("%d_%d", cfgCharge.ProductID, cfgCharge.PlatSystemType)
	return pbBingchuan
}

// 更新玩家充值
func (this *PlayerSystemChargeFun) UpdatePlayerCharge(head *pb.RpcHead) []uint32 {
	this.bCharge = false
	//读取所有的充值数据
	redisGame := redis.GetRedisByAccountID(this.AccountId)
	if redisGame == nil {
		return nil
	}

	//遍历邮件
	arrList := make([]uint32, 0)
	strKey := fmt.Sprintf("%s%d", base.ERK_GamePlayerCharge, this.AccountId)
	for {
		strData := redisGame.RPop(strKey)
		if strData == "" {
			break
		}

		uProductId := base.StringToUInt32(strData)
		if cfgEnum.ErrorCode_Success == this.AddProductId(head, uProductId) {
			arrList = append(arrList, uProductId)
		}
	}
	return arrList
}

// 隔天刷新通知
func (this *PlayerSystemChargeFun) PassDay(isDay, isWeek, isMonth bool) {
	this.pbData.DailyCharge = 0
	if isWeek {
		this.pbData.WeekCharge = 0
	}

	if isMonth {
		this.pbData.MonthCharge = 0
	}

	//同步充值所有的数据
	pbSystemCharge := &pb.PBPlayerSystemCharge{}
	this.saveData(pbSystemCharge)

	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.ChargeNotify{
		PacketHead: &pb.IPacket{},
		Charge:     pbSystemCharge.Charge,
	}, cfgEnum.ErrorCode_Success)

	for _, info := range this.mapFun {
		info.PassDay(isDay, isWeek, isMonth)
	}
}

func (this *PlayerSystemChargeFun) AddProductId(head *pb.RpcHead, uProductId uint32) cfgEnum.ErrorCode {
	cfgCharge := cfgData.GetCfgCharge(uProductId, this.getPlayerBaseFun().GetPlayerBase().PlatSystemType)
	if cfgCharge == nil {
		return plog.Print(this.AccountId, cfgData.GetChargeErrorCode(uProductId), uProductId, this.getPlayerBaseFun().GetPlayerBase().PlatSystemType)
	}

	this.pbData.DailyCharge += cfgCharge.Price
	this.pbData.WeekCharge += cfgCharge.Price
	this.pbData.MonthCharge += cfgCharge.Price
	this.pbData.TotalCharge += cfgCharge.Price

	this.UpdateSave(true)

	//同步充值所有的数据
	pbSystemCharge := &pb.PBPlayerSystemCharge{}
	this.saveData(pbSystemCharge)

	cluster.SendToClient(head, &pb.ChargeNotify{
		PacketHead: &pb.IPacket{},
		Charge:     pbSystemCharge.Charge,
	}, cfgEnum.ErrorCode_Success)

	//根据充值类型 对应各个系统中
	switch cfgEnum.EChargeType(cfgCharge.ChargeType) {
	case cfgEnum.EChargeType_SevenDay:
		this.getPlayerSystemSevenDayFun().OnChargeBuyProduct(head, cfgCharge)
	case cfgEnum.EChargeType_BP, cfgEnum.EChargeType_ChargeCard:
		this.getChargeFun(cfgEnum.EChargeType(cfgCharge.ChargeType)).OnChargeBuyProduct(head, cfgCharge)
	case cfgEnum.EChargeType_ChargeGift:
		this.getPlayerActivityChargeGiftFun().OnChargeBuyGiftProduct(head, cfgCharge)
	case cfgEnum.EChargeType_ShopCharge:
		this.getPlayerSystemShopFun().OnChargeBuyProduct(head, cfgCharge)
	case cfgEnum.EChargeType_OpenServerGift:
		this.getPlayerActivityOpenServerGiftFun().OnChargeBuyProduct(head, cfgCharge)
	case cfgEnum.EChargeType_FirstCharge:
		this.OnChargeBuyFirstChargeProduct(head, cfgCharge)
	}

	return cfgEnum.ErrorCode_Success
}

// 系统解锁开启首冲
func (this *PlayerSystemChargeFun) OnSystemOpenFirstCharge(head *pb.RpcHead, bNormal bool) {
	cfgAll := cfgData.GetAllCfgFirstChargeConfig()
	for _, cfg := range cfgAll {
		if _, ok := this.mapFirstCharge[cfg.Id]; !ok {
			uOpenTime := base.GetNow()
			if !bNormal {
				uOpenTime = this.getPlayerBaseFun().GetRegTime()
			}
			this.mapFirstCharge[cfg.Id] = &pb.PBFirstCharge{
				FirstChargeId: cfg.Id,
				OpenTime:      uOpenTime,
			}

			//通知客户端
			cluster.SendToClient(head, &pb.FirstChargeNotify{
				PacketHead:      &pb.IPacket{},
				FirstChargeInfo: this.mapFirstCharge[cfg.Id],
			}, cfgEnum.ErrorCode_Success)
		}
	}
}

// 是否能够购买
func (this *PlayerSystemChargeFun) canBuyFirstChargeProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) cfgEnum.ErrorCode {
	if len(cfgCharge.Param) != 1 {
		return plog.Print(this.AccountId, cfgData.GetChargeErrorCode(cfgCharge.ProductID), cfgCharge.ProductID)
	}

	uId := cfgCharge.Param[0]
	cfgFirstCharge := cfgData.GetCfgFirstChargePrizeConfig(uId)
	if cfgFirstCharge == nil {
		return plog.Print(this.AccountId, cfgData.GetFirstChargePrizeConfigErrorCode(cfgCharge.ProductID), cfgCharge.ProductID)
	}

	if _, ok := this.mapFirstCharge[uId]; !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_SystemNoOpen, cfgCharge.ProductID, uId)
	}

	//充值过不能重复
	if this.mapFirstCharge[uId].ActiveTime > 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ChargeRepeated, cfgCharge.ProductID, uId)
	}

	return cfgEnum.ErrorCode_Success
}

// 首冲领奖请求
func (this *PlayerSystemChargeFun) OnChargeBuyFirstChargeProduct(head *pb.RpcHead, cfgCharge *cfgData.ChargeCfg) {
	if cfgCharge == nil || cfgCharge.ChargeType != uint32(cfgEnum.EChargeType_FirstCharge) || len(cfgCharge.Param) != 1 {
		plog.Error("OnChargeBuyFirstChargeProduct accountid:%d", this.AccountId)
		return
	}
	if this.canBuyFirstChargeProduct(head, cfgCharge) != cfgEnum.ErrorCode_Success {
		return
	}

	uId := cfgCharge.Param[0]

	this.mapFirstCharge[uId].ActiveTime = base.GetNow()

	//通知客户端
	cluster.SendToClient(head, &pb.FirstChargeNotify{
		PacketHead:      &pb.IPacket{},
		FirstChargeInfo: this.mapFirstCharge[uId],
	}, cfgEnum.ErrorCode_Success)
}

// 首冲领奖请求
func (this *PlayerSystemChargeFun) FirstChargePrizeRequest(head *pb.RpcHead, pbRequest *pb.FirstChargePrizeRequest) {
	uCode := this.FirstChargePrize(head, pbRequest.FirstChargeId)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.FirstChargePrizeResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 首冲领奖请求
func (this *PlayerSystemChargeFun) FirstChargePrize(head *pb.RpcHead, uFirstChargeId uint32) cfgEnum.ErrorCode {
	pbFirstCharge, ok := this.mapFirstCharge[uFirstChargeId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, uFirstChargeId)
	}

	uCurTime := base.GetNow()
	if pbFirstCharge.ActiveTime <= 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, uFirstChargeId)
	}

	uPassDay := base.DiffDays(pbFirstCharge.OpenTime, uCurTime) + 1
	listPrizeCfg := cfgData.GetCfgFirstChargePrizeConfig(uFirstChargeId)
	if pbFirstCharge.PrizeDay >= uint32(len(listPrizeCfg)) || pbFirstCharge.PrizeDay >= uPassDay {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HavePrize, uFirstChargeId)
	}

	arrAddPrize := make([]*common.ItemInfo, 0)
	for _, cfgPrize := range listPrizeCfg {
		if uPassDay < cfgPrize.NeedDay {
			break
		}

		if pbFirstCharge.PrizeDay >= cfgPrize.NeedDay {
			continue
		}

		pbFirstCharge.PrizeDay = cfgPrize.NeedDay
		arrAddPrize = append(arrAddPrize, cfgPrize.AddPrize...)
	}

	if len(arrAddPrize) <= 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HavePrize, uFirstChargeId)
	}

	//加奖励
	this.getPlayerBagFun().AddArrItem(head, arrAddPrize, pb.EmDoingType_EDT_FirstCharge, true)

	this.UpdateSave(true)
	cluster.SendToClient(head, &pb.FirstChargePrizeResponse{
		PacketHead:    &pb.IPacket{},
		FirstChargeId: uFirstChargeId,
		PrizeDay:      pbFirstCharge.PrizeDay,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}
func (this *PlayerSystemChargeFun) ChargeQueryRequest(head *pb.RpcHead) {
	cluster.SendToClient(head, &pb.ChargeQueryResponse{
		PacketHead: &pb.IPacket{},
		ProductIds: this.UpdatePlayerCharge(head),
	}, cfgEnum.ErrorCode_Success)
}
func (this *PlayerSystemChargeFun) DipUpdatePlayerCharge(head *pb.RpcHead, ProductId uint32) {
	this.bCharge = true

	cluster.SendToClient(head, &pb.ChargeQueryNotify{
		PacketHead: &pb.IPacket{},
		ProductId:  ProductId,
	}, cfgEnum.ErrorCode_Success)
}
