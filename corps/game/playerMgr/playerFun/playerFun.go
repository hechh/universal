package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common/cfgData"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/module/entry"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/cast"
)

type (
	PlayerFun struct {
		FunCommon
		PbType pb.PlayerDataType
		Self   IPlayerFun
		BSave  bool //是否存库
	}

	FunCommon struct {
		AccountId uint64
		MapFun    *map[pb.PlayerDataType]IPlayerFun
	}

	IPlayerFun interface {
		GetPbType() pb.PlayerDataType                     //初始化
		Init(pbType pb.PlayerDataType, common *FunCommon) //初始化
		UpdateCommon(common *FunCommon)                   //初始化
		Load(pData []byte)                                //加载数据
		Save(bNow bool)                                   //存储数据
		IsSave() bool                                     //是否存储数据
		LoadComplete()                                    //加载完成
		LoadSystem(pbSystem *pb.PBPlayerSystem)           //加载系统数据
		SaveSystem(pbSystem *pb.PBPlayerSystem) bool      //存储数据 返回存储标志
		LoadPlayerDBFinish()                              //加载系统数据DB数据 数据初始化用
		NewPlayer()                                       //新系统
		SaveDataToClient(pbData *pb.PBPlayerData)         //拷贝数据
		Heat()                                            //心跳包
		PassDay(isDay, isWeek, isMonth bool)              //是否跨天
		UpdateSave(bSave bool)                            //保存
		GetProtoPtr() proto.Message                       //获取pb指针
		SetUserTypeInfo(message proto.Message) bool       //设置玩家数据
	}
)

func (this *PlayerFun) GetPbType() pb.PlayerDataType {
	return this.PbType
}
func (this *PlayerFun) IsSave() bool {
	return this.BSave
}
func (this *PlayerFun) SetUserTypeInfo(message proto.Message) bool {
	return true
}
func (this *PlayerFun) Load(pData []byte) {
}
func (this *PlayerFun) LoadComplete() {
}
func (this *PlayerFun) Save(bNow bool) {
}

func (this *PlayerFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
}

func (this *PlayerFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	return true
}

func (this *PlayerFun) UpdateSave(bSave bool) {
	this.BSave = bSave
}

func (this *PlayerFun) UpdateCommon(common *FunCommon) {
	this.FunCommon = *common
}

func (this *PlayerFun) Heat() {
}

func (this *PlayerFun) PassDay(isDay, isWeek, isMonth bool) {
}

func (this *PlayerFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PbType = pbType
	this.FunCommon = *common
}

func (this *PlayerFun) LoadPlayerDBFinish() {
}
func (this *PlayerFun) NewPlayer() {
}

// 客户端通信
func (this *PlayerFun) getPlayerFun(emType pb.PlayerDataType) IPlayerFun {
	fun, ok := (*this.MapFun)[emType]
	if !ok {
		return nil
	}

	return fun
}

// 晶核系统
func (this *PlayerFun) getPlayerCrystalFun() *PlayerCrystalFun {
	fun := this.getPlayerFun(pb.PlayerDataType_Crystal)
	if fun == nil {
		return nil
	}
	return this.getPlayerFun(pb.PlayerDataType_Crystal).(*PlayerCrystalFun)
}

// 玩家基本数据
func (this *PlayerFun) getPlayerBaseFun() *PlayerBaseFun {
	return this.getPlayerFun(pb.PlayerDataType_Base).(*PlayerBaseFun)
}

// 玩家背包数据
func (this *PlayerFun) getPlayerBagFun() *PlayerBagFun {
	return this.getPlayerFun(pb.PlayerDataType_Bag).(*PlayerBagFun)
}

// 邮件数据
func (this *PlayerFun) getPlayerMailFun() *PlayerMailFun {
	return this.getPlayerFun(pb.PlayerDataType_Mail).(*PlayerMailFun)
}

// 玩家系统通用数据
func (this *PlayerFun) getPlayerSystemCommonFun() *PlayerSystemCommonFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemCommon).(*PlayerSystemCommonFun)
}

// 商店系统
func (this *PlayerFun) getPlayerSystemShopFun() *PlayerSystemShopFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemShop).(*PlayerSystemShopFun)
}
func (this *PlayerFun) getPlayerSystemDrawFun() *PlayerSystemDrawFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemDraw).(*PlayerSystemDrawFun)
}
func (this *PlayerFun) getPlayerSystemHookTechFun() *PlayerSystemHookTechFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemHookTech).(*PlayerSystemHookTechFun)
}
func (this *PlayerFun) getPlayerSystemSevenDayFun() *PlayerSystemSevenDayFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemSevenDay).(*PlayerSystemSevenDayFun)
}
func (this *PlayerFun) getPlayerSystemActivityFun() *PlayerSystemActivityFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemActivity).(*PlayerSystemActivityFun)
}
func (this *PlayerFun) getPlayerActivityChargeGiftFun() *PlayerActivityChargeGift {
	return this.getPlayerSystemActivityFun().GetActivityChargeGiftFun()
}
func (this *PlayerFun) getPlayerActivityOpenServerGiftFun() *PlayerActivityOpenServerGift {
	return this.getPlayerSystemActivityFun().GetActivityOpenServerGiftFun()
}
func (this *PlayerFun) getPlayerSystemWorldBossFun() *PlayerSystemWorldBossFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemWorldBoss).(*PlayerSystemWorldBossFun)
}

func (this *PlayerFun) getPlayerSystemChampionshipFun() *PlayerSystemChampionshipFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemChampionship).(*PlayerSystemChampionshipFun)
}

// 玩家系统职业数据
func (this *PlayerFun) getPlayerSystemProfessionFun() *PlayerSystemProfessionFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemProfession).(*PlayerSystemProfessionFun)
}

func (this *PlayerFun) getPlayerSystemBoxFun() *PlayerSystemBoxFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemBox).(*PlayerSystemBoxFun)
}

// 玩家系统战斗数据
func (this *PlayerFun) getPlayerSystemBattleFun() *PlayerSystemBattleFun {
	fun := this.getPlayerFun(pb.PlayerDataType_SystemBattle)
	if fun == nil {
		return nil
	}
	return this.getPlayerFun(pb.PlayerDataType_SystemBattle).(*PlayerSystemBattleFun)
}

// 玩家系统战斗数据
func (this *PlayerFun) getPlayerSystemBattleNormalFun() *PlayerSystemBattleNormalFun {
	return this.getPlayerSystemBattleFun().GetBattleNoramlFun()
}
func (this *PlayerFun) getPlayerSystemBattleHookFun() *PlayerSystemBattleHookFun {
	return this.getPlayerSystemBattleFun().GetBattleHookFun()
}

// 玩家系统任务数据
func (this *PlayerFun) getPlayerSystemTaskFun() *PlayerSystemTaskFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemTask).(*PlayerSystemTaskFun)
}
func (this *PlayerFun) getPlayerSystemChargeFun() *PlayerSystemChargeFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemCharge).(*PlayerSystemChargeFun)
}
func (this *PlayerFun) getPlayerSystemChargeBPFun() *PlayerSystemChargeBP {
	return this.getPlayerSystemChargeFun().GetChargeBPFun()
}
func (this *PlayerFun) getPlayerSystemChargeCardFun() *PlayerSystemChargeCard {
	return this.getPlayerSystemChargeFun().GetChargeCardFun()
}
func (this *PlayerFun) getPlayerSystemClientFun() *PlayerClientFun {
	return this.getPlayerFun(pb.PlayerDataType_Client).(*PlayerClientFun)
}

// 玩家系统装备数据
func (this *PlayerFun) GetPlayerEquipmentFun() *PlayerEquipmentFun {
	return this.getPlayerFun(pb.PlayerDataType_Equipment).(*PlayerEquipmentFun)
}

// 玩家伙伴数据
func (this *PlayerFun) getPlayerHeroFun() *PlayerHeroFun {
	return this.getPlayerFun(pb.PlayerDataType_Hero).(*PlayerHeroFun)
}

// 玩家系统宝箱数据
func (this *PlayerFun) getPlayerSystemBox() *PlayerSystemBoxFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemBox).(*PlayerSystemBoxFun)
}

// 玩家基因系统
func (this *PlayerFun) getPlayerSystemGene() *PlayerSystemGeneFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemGene).(*PlayerSystemGeneFun)
}

// 离线收益系统
func (this *PlayerFun) getPlayerSystemOffline() *PlayerSystemOfflineFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemOffline).(*PlayerSystemOfflineFun)
}

// 词条技能
func (this *PlayerFun) getEntry() *entry.EntryService {
	fun := this.getPlayerCrystalFun()
	if fun == nil {
		return nil
	}
	return fun.entryData
}

// 更新排行榜
func (this *PlayerFun) UpdateRank(active uint64, rankType uint32, val uint64) {
	if val <= 0 {
		return
	}
	// 加载配置
	rankTypeCfg := cfgData.GetCfgRankTypeConfig(rankType)
	if rankTypeCfg == nil {
		plog.Error("rankTypeCondig(%d) not found", rankType)
		return
	}
	// 积分按时间排序
	now := base.GetNow()
	if rankTypeCfg.SortType > 0 {
		val = (val << 32) | uint64(MAX_DWORD_VALUE-uint32(now))
	}
	// 判断是否在线
	switch rankTypeCfg.DataType {
	case uint32(cfgEnum.EDataType_Expire):
		infoCfg := cfgData.GetCfgRankInfoConfig(rankType)
		if infoCfg == nil {
			plog.Error("rankInfoCondig(%d) not found", rankType)
			return
		}
		if now < active || cfgData.GetCfgRankRewardTimeByActiveTime(infoCfg, active) <= now {
			return
		}
	case uint32(cfgEnum.EDataType_Forever):
	}
	// 设置转发头
	head := &pb.RpcHead{
		Id:        this.AccountId,
		RegionID:  this.getPlayerBaseFun().GetServerId(),
		RouteType: uint32(cfgEnum.ERouteType_ServerID),
	}
	switch rankTypeCfg.RouteType {
	case uint32(cfgEnum.ERedisType_Global): // 全局节点
		head.RegionID = 0
	case uint32(cfgEnum.ERedisType_Random): // 负载均衡
	}
	// 发送更新排行榜通知
	notify := &pb.RankUpdateNotify{
		PacketHead: &pb.IPacket{},
		RankType:   rankType,
		CreateTime: active,
		Member:     cast.ToString(this.AccountId),
		Score:      float64(val),
	}
	// 发送gm消息，更新排行榜
	this.getPlayerSystemChampionshipFun().SetChampionshipFlag(rankType, 1)
	cluster.SendToGm(head, "RankMgr", "UpdateRequest", notify)
}
