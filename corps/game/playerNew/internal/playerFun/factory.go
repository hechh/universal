package playerFun

import (
	"corps/pb"
	"corps/server/game/module/entry"
	"corps/server/game/playerNew/internal/manager"
)

type PlayerFunFactory struct {
	*manager.FunMgr // factor句柄
}

// --------------------------------------公用接口----------------------------------
// 词条技能
func (this *PlayerFunFactory) GetEntry() *entry.EntryService {
	return this.GetPlayerCrystalFun().GetEntry()
}

// 玩家系统战斗数据
func (this *PlayerFunFactory) GetPlayerSystemBattleNormalFun() *PlayerSystemBattleNormalFun {
	return this.GetPlayerSystemBattleFun().GetBattleNoramlFun()
}
func (this *PlayerFunFactory) GetPlayerSystemBattleHookFun() *PlayerSystemBattleHookFun {
	return this.GetPlayerSystemBattleFun().GetBattleHookFun()
}

func (this *PlayerFunFactory) GetPlayerBaseFun() *PlayerBaseFun {
	return this.GetIPlayerFun(pb.PlayerDataType_Base).(*PlayerBaseFun)
}
func (this *PlayerFunFactory) GetPlayerBagFun() *PlayerBagFun {
	return this.GetIPlayerFun(pb.PlayerDataType_Bag).(*PlayerBagFun)
}
func (this *PlayerFunFactory) GetPlayerEquipmentFun() *PlayerEquipmentFun {
	return this.GetIPlayerFun(pb.PlayerDataType_Equipment).(*PlayerEquipmentFun)
}
func (this *PlayerFunFactory) GetPlayerClientFun() *PlayerClientFun {
	return this.GetIPlayerFun(pb.PlayerDataType_Client).(*PlayerClientFun)
}
func (this *PlayerFunFactory) GetPlayerHeroFun() *PlayerHeroFun {
	return this.GetIPlayerFun(pb.PlayerDataType_Hero).(*PlayerHeroFun)
}
func (this *PlayerFunFactory) GetPlayerMailFun() *PlayerMailFun {
	return this.GetIPlayerFun(pb.PlayerDataType_Mail).(*PlayerMailFun)
}
func (this *PlayerFunFactory) GetPlayerCrystalFun() *PlayerCrystalFun {
	return this.GetIPlayerFun(pb.PlayerDataType_Crystal).(*PlayerCrystalFun)
}
func (this *PlayerFunFactory) GetPlayerSystemCommonFun() *PlayerSystemCommonFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemCommon).(*PlayerSystemCommonFun)
}
func (this *PlayerFunFactory) GetPlayerSystemProfessionFun() *PlayerSystemProfessionFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemProfession).(*PlayerSystemProfessionFun)
}
func (this *PlayerFunFactory) GetPlayerSystemBattleFun() *PlayerSystemBattleFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemBattle).(*PlayerSystemBattleFun)
}
func (this *PlayerFunFactory) GetPlayerSystemBoxFun() *PlayerSystemBoxFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemBox).(*PlayerSystemBoxFun)
}
func (this *PlayerFunFactory) GetPlayerSystemTaskFun() *PlayerSystemTaskFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemTask).(*PlayerSystemTaskFun)
}
func (this *PlayerFunFactory) GetPlayerSystemShopFun() *PlayerSystemShopFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemShop).(*PlayerSystemShopFun)
}
func (this *PlayerFunFactory) GetPlayerSystemDrawFun() *PlayerSystemDrawFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemDraw).(*PlayerSystemDrawFun)
}
func (this *PlayerFunFactory) GetPlayerSystemChargeFun() *PlayerSystemChargeFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemCharge).(*PlayerSystemChargeFun)
}
func (this *PlayerFunFactory) GetPlayerSystemGeneFun() *PlayerSystemGeneFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemGene).(*PlayerSystemGeneFun)
}
func (this *PlayerFunFactory) GetPlayerSystemOfflineFun() *PlayerSystemOfflineFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemOffline).(*PlayerSystemOfflineFun)
}
func (this *PlayerFunFactory) GetPlayerSystemHookTechFun() *PlayerSystemHookTechFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemHookTech).(*PlayerSystemHookTechFun)
}
func (this *PlayerFunFactory) GetPlayerSystemRepairFun() *PlayerSystemRepairFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemRepair).(*PlayerSystemRepairFun)
}
func (this *PlayerFunFactory) GetPlayerSystemSevenDayFun() *PlayerSystemSevenDayFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemSevenDay).(*PlayerSystemSevenDayFun)
}
func (this *PlayerFunFactory) GetPlayerSystemActivityFun() *PlayerSystemActivityFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemActivity).(*PlayerSystemActivityFun)
}
func (this *PlayerFunFactory) GetPlayerSystemWorldBossFun() *PlayerSystemWorldBossFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemWorldBoss).(*PlayerSystemWorldBossFun)
}
func (this *PlayerFunFactory) GetPlayerSystemChampionshipFun() *PlayerSystemChampionshipFun {
	return this.GetIPlayerFun(pb.PlayerDataType_SystemChampionship).(*PlayerSystemChampionshipFun)
}
