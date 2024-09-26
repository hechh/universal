package playerFun

import (
	"corps/pb"
	"corps/server/game/player/internal/manager"
)

type PlayerFun struct {
	*manager.FunMgr
	accountId uint64
}

func NewPlayerFun() *PlayerFun {
	return &PlayerFun{FunMgr: manager.NewFunMgr()}
}

func (d *PlayerFun) GetPlayerDataType() pb.PlayerDataType {
	return pb.PlayerDataType_Max
}

// 注册定时器
func (d *PlayerFun) RegisterTimer() {}

// 初始化
func (d *PlayerFun) Init(uid uint64, data interface{}) {
	d.accountId = uid
}

// 加载数据(非system数据)
func (d *PlayerFun) Load([]byte) {}

// 判断是否存储数据
func (d *PlayerFun) IsSave() bool { return false }

// 存储数据(非system数据)
func (d *PlayerFun) Save(bNow bool) {}

// db数据加载完成回调，在LoadComplete之前调用
func (d *PlayerFun) LoadPlayerDBFinish() {}

// 初始化新开启的模块数据，在LoadComplete之后调用
func (d *PlayerFun) NewPlayer() {}

// 加载完成，在NewPlayer之后调用
func (d *PlayerFun) LoadComplete() {}

// 设置保存状态
func (d *PlayerFun) UpdateSave(bSave bool) {}

// 心跳包
func (d *PlayerFun) Heat() {}

// 是否跨天
func (d *PlayerFun) PassDay(isDay, isWeek, isMonth bool) {}

// 设置缓存数据
func (d *PlayerFun) SetUserTypeInfo([]byte) error {
	return nil
}

// 深度拷贝数据
func (d *PlayerFun) CopyTo(pbData *pb.PBPlayerData) {}

func (this *PlayerFun) GetPlayerBaseFun() *PlayerBaseFun {
	return this.GetIPlayerFun(pb.PlayerDataType_Base).(*PlayerBaseFun)
}
func (this *PlayerFun) GetPlayerBagFun() *PlayerBagFun {
	return this.GetIPlayerFun(pb.PlayerDataType_Bag).(*PlayerBagFun)
}
func (this *PlayerFun) GetPlayerEquipmentFun() *PlayerEquipmentFun {
	return this.GetIPlayerFun(pb.PlayerDataType_Equipment).(*PlayerEquipmentFun)
}
func (this *PlayerFun) GetPlayerClientFun() *PlayerClientFun {
	return this.GetIPlayerFun(pb.PlayerDataType_Client).(*PlayerClientFun)
}
func (this *PlayerFun) GetPlayerHeroFun() *PlayerHeroFun {
	return this.GetIPlayerFun(pb.PlayerDataType_Hero).(*PlayerHeroFun)
}
func (this *PlayerFun) GetPlayerMailFun() *PlayerMailFun {
	return this.GetIPlayerFun(pb.PlayerDataType_Mail).(*PlayerMailFun)
}
func (this *PlayerFun) GetPlayerCrystalFun() *PlayerCrystalFun {
	return this.GetIPlayerFun(pb.PlayerDataType_Crystal).(*PlayerCrystalFun)
}
func (this *PlayerFun) GetPlayerSystemCommonFun() *PlayerSystemCommonFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemCommon).(*PlayerSystemCommonFun)
}
func (this *PlayerFun) GetPlayerSystemProfessionFun() *PlayerSystemProfessionFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemProfession).(*PlayerSystemProfessionFun)
}
func (this *PlayerFun) GetPlayerSystemBattleFun() *PlayerSystemBattleFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemBattle).(*PlayerSystemBattleFun)
}
func (this *PlayerFun) GetPlayerSystemBoxFun() *PlayerSystemBoxFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemBox).(*PlayerSystemBoxFun)
}
func (this *PlayerFun) GetPlayerSystemTaskFun() *PlayerSystemTaskFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemTask).(*PlayerSystemTaskFun)
}
func (this *PlayerFun) GetPlayerSystemShopFun() *PlayerSystemShopFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemShop).(*PlayerSystemShopFun)
}
func (this *PlayerFun) GetPlayerSystemDrawFun() *PlayerSystemDrawFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemDraw).(*PlayerSystemDrawFun)
}
func (this *PlayerFun) GetPlayerSystemChargeFun() *PlayerSystemChargeFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemCharge).(*PlayerSystemChargeFun)
}
func (this *PlayerFun) GetPlayerSystemGeneFun() *PlayerSystemGeneFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemGene).(*PlayerSystemGeneFun)
}
func (this *PlayerFun) GetPlayerSystemOfflineFun() *PlayerSystemOfflineFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemOffline).(*PlayerSystemOfflineFun)
}
func (this *PlayerFun) GetPlayerSystemHookTechFun() *PlayerSystemHookTechFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemHookTech).(*PlayerSystemHookTechFun)
}
func (this *PlayerFun) GetPlayerSystemRepairFun() *PlayerSystemRepairFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemRepair).(*PlayerSystemRepairFun)
}
func (this *PlayerFun) GetPlayerSystemSevenDayFun() *PlayerSystemSevenDayFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemSevenDay).(*PlayerSystemSevenDayFun)
}
func (this *PlayerFun) GetPlayerSystemActivityFun() *PlayerSystemActivityFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemActivity).(*PlayerSystemActivityFun)
}
func (this *PlayerFun) GetPlayerSystemWorldBossFun() *PlayerSystemWorldBossFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemWorldBoss).(*PlayerSystemWorldBossFun)
}
func (this *PlayerFun) GetPlayerSystemChampionshipFun() *PlayerSystemChampionshipFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemChampionship).(*PlayerSystemChampionshipFun)
}
