package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common/cfgData"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"

	"github.com/golang/protobuf/proto"
)

type (
	PlayerSystemHookTechFun struct {
		PlayerFun
		mapHookTech       map[uint32]*pb.PBHookTech
		mapHookTechEffect map[cfgEnum.TechEffectType]map[uint32]int32 //挂机科技效果
	}
)

func (this *PlayerSystemHookTechFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.mapHookTech = make(map[uint32]*pb.PBHookTech)
	this.mapHookTechEffect = make(map[cfgEnum.TechEffectType]map[uint32]int32)
}

// 从数据库中加载
func (this *PlayerSystemHookTechFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	this.loadData(pbSystem.HookTech)
	this.UpdateSave(false)
}

func (this *PlayerSystemHookTechFun) loadData(pbData *pb.PBPlayerSystemHookTech) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemHookTech{}
	}

	this.mapHookTech = make(map[uint32]*pb.PBHookTech)
	for _, info := range pbData.HookTechList {
		this.mapHookTech[info.Id] = info
	}

	this.UpdateSave(true)
}

// 加载完成需要计算战斗力
func (this *PlayerSystemHookTechFun) LoadComplete() {
	this.CalcHookTechEffect()
}

// 存储到数据库
func (this *PlayerSystemHookTechFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	if pbSystem.HookTech == nil {
		pbSystem.HookTech = new(pb.PBPlayerSystemHookTech)
	}

	pbSystem.HookTech.HookTechList = make([]*pb.PBHookTech, 0)
	for _, value := range this.mapHookTech {
		pbSystem.HookTech.HookTechList = append(pbSystem.HookTech.HookTechList, value)
	}

	return this.BSave
}
func (this *PlayerSystemHookTechFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemHookTech{}
}

// 设置玩家数据
func (this *PlayerSystemHookTechFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}
	pbSystem := pbData.(*pb.PBPlayerSystemHookTech)
	if pbSystem == nil {
		return false
	}
	this.loadData(pbSystem)
	this.CalcHookTechEffect()
	this.UpdateSave(true)
	return true
}

func (this *PlayerSystemHookTechFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = new(pb.PBPlayerSystem)
	}

	this.SaveSystem(pbData.System)
}

// 心跳包 检查
func (this *PlayerSystemHookTechFun) Heat() {
	uNow := base.GetNow()

	pbNotify := &pb.HookTechLevelNotify{
		PacketHead: &pb.IPacket{},
	}
	for _, info := range this.mapHookTech {
		if info.LevelTime == 0 || info.LevelTime > uNow {
			continue
		}

		info.LevelTime = 0
		info.Level++
		pbNotify.HookTechList = append(pbNotify.HookTechList, info)
	}

	if len(pbNotify.HookTechList) > 0 {
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbNotify, cfgEnum.ErrorCode_Success)
		this.CalcHookTechEffect()
	}
}

// 挂机科技升级请求
func (this *PlayerSystemHookTechFun) HookTechLevelRequest(head *pb.RpcHead, pbRequest *pb.HookTechLevelRequest) {
	uCode := this.HookTechLevel(head, pbRequest.Id)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.HookTechLevelResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 挂机科技升级请求
func (this *PlayerSystemHookTechFun) HookTechLevel(head *pb.RpcHead, uId uint32) cfgEnum.ErrorCode {
	cfgHookTech := cfgData.GetCfgHookTechConfig(uId)
	if cfgHookTech == nil {
		return plog.Print(this.AccountId, cfgData.GetHookTechConfigErrorCode(uId), uId)
	}

	pbTech, ok := this.mapHookTech[uId]
	if !ok {
		pbTech = &pb.PBHookTech{
			Id:    uId,
			Level: 0,
		}
	}

	//判断是否满级
	cfgNewLevel := cfgData.GetCfgHookTechLevelConfig(uId, pbTech.Level+1)
	if cfgNewLevel == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxLevel, uId)
	}

	if pbTech.LevelTime > 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_CoolDown, uId)
	}

	//判断前置是否满级
	for _, fatherId := range cfgHookTech.FatherId {
		if _, ok := this.mapHookTech[fatherId]; !ok {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedPreLevel, uId)
		}
		cfgFather := cfgData.GetCfgHookTechConfig(fatherId)
		if cfgFather == nil {
			return plog.Print(this.AccountId, cfgData.GetHookTechConfigErrorCode(fatherId), uId)
		}

		if cfgData.GetCfgHookTechLevelConfig(fatherId, this.mapHookTech[fatherId].Level+1) != nil {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedPreLevel, uId)
		}
	}

	//判断道具
	uCode := this.getPlayerBagFun().DelItem(head, cfgNewLevel.DelItem.Kind, cfgNewLevel.DelItem.Id, cfgNewLevel.DelItem.Count, pb.EmDoingType_EDT_HookTech)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode, uId)
	}

	pbTech.LevelTime = base.GetNow() + uint64(cfgNewLevel.LevelCoolTime)

	this.mapHookTech[uId] = pbTech
	this.UpdateSave(true)

	//成就统计
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_HookTechLevel, 1)

	//通知客户端
	cluster.SendToClient(head, &pb.HookTechLevelResponse{
		PacketHead: &pb.IPacket{},
		Id:         uId,
		LevelTime:  pbTech.LevelTime,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 挂机科技加速请求
func (this *PlayerSystemHookTechFun) HookTechSpeedRequest(head *pb.RpcHead, pbRequest *pb.HookTechSpeedRequest) {
	uCode := this.HookTechSpeed(head, pbRequest.Id, pbRequest.AdvertType)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.HookTechSpeedResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 挂机科技加速请求
func (this *PlayerSystemHookTechFun) HookTechSpeed(head *pb.RpcHead, uId uint32, uAdvertType uint32) cfgEnum.ErrorCode {
	cfgHookTech := cfgData.GetCfgHookTechConfig(uId)
	if cfgHookTech == nil {
		return plog.Print(this.AccountId, cfgData.GetHookTechConfigErrorCode(uId), uId)
	}

	pbTech, ok := this.mapHookTech[uId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, uId)
	}

	uNow := base.GetNow()
	if pbTech.LevelTime <= 0 || pbTech.LevelTime <= uNow {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HookTechSpeedTime, uId)
	}
	bLevel := false
	cfgConst := cfgData.GetCfgHookTechConst()
	uNeedMins := base.CeilU32(uint32(pbTech.LevelTime-uNow), 60)
	if uAdvertType == uint32(cfgEnum.EAdvertType_None) {
		uCode := this.getPlayerBagFun().DelItem(head, cfgConst.HookTechSpeedDelItem.Kind, cfgConst.HookTechSpeedDelItem.Id, cfgConst.HookTechSpeedDelItem.Count*int64(uNeedMins), pb.EmDoingType_EDT_HookTech)
		if uCode != cfgEnum.ErrorCode_Success {
			return plog.Print(this.AccountId, uCode, uId, *cfgConst.HookTechSpeedDelItem, uNeedMins)
		}
		pbTech.LevelTime = 0
		pbTech.Level++
		bLevel = true
	} else if uAdvertType == uint32(cfgEnum.EAdvertType_HookTechSpeed) {
		uCode := this.getPlayerSystemCommonFun().AddAdvert(head, uAdvertType)
		if uCode != cfgEnum.ErrorCode_Success {
			return plog.Print(this.AccountId, uCode, uId)
		}

		if cfgConst.HookTechLevelAdvertSpeedTime >= uNeedMins {
			pbTech.LevelTime = 0
			pbTech.Level++
			bLevel = true
		} else {
			pbTech.LevelTime -= uint64(cfgConst.HookTechLevelAdvertSpeedTime * 60)
		}
	} else {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_AdvertNotEqual, uId)
	}

	if bLevel {
		this.CalcHookTechEffect()
	}

	this.UpdateSave(true)

	//通知客户端
	cluster.SendToClient(head, &pb.HookTechSpeedResponse{
		PacketHead: &pb.IPacket{},
		Id:         uId,
		Level:      this.mapHookTech[uId].Level,
		LevelTime:  this.mapHookTech[uId].LevelTime,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

func (this *PlayerSystemHookTechFun) GetHookTechEffect(emType cfgEnum.TechEffectType) map[uint32]int32 {
	return this.mapHookTechEffect[emType]
}

// 计算挂机科技属性
func (this *PlayerSystemHookTechFun) CalcHookTechEffect() {
	this.mapHookTechEffect = make(map[cfgEnum.TechEffectType]map[uint32]int32)
	for _, info := range this.mapHookTech {
		if info.Level <= 0 {
			continue
		}
		cfgHookTech := cfgData.GetCfgHookTechConfig(info.Id)
		if cfgHookTech == nil {
			continue
		}
		cfgHookTechLevel := cfgData.GetCfgHookTechLevelConfig(info.Id, info.Level)
		if cfgHookTechLevel == nil {
			continue
		}
		if _, ok := this.mapHookTechEffect[cfgEnum.TechEffectType(cfgHookTech.TechEffectType)]; !ok {
			this.mapHookTechEffect[cfgEnum.TechEffectType(cfgHookTech.TechEffectType)] = make(map[uint32]int32)
		}

		switch cfgEnum.TechEffectType(cfgHookTech.TechEffectType) {
		case cfgEnum.TechEffectType_AddHookEquipMaxStar: //取最大值
			value, ok := this.mapHookTechEffect[cfgEnum.TechEffectType(cfgHookTech.TechEffectType)][0]
			if !ok {
				value = 0
			}

			this.mapHookTechEffect[cfgEnum.TechEffectType(cfgHookTech.TechEffectType)][0] = base.MaxInt32(value, cfgHookTechLevel.MapEffect[0])
		default:
			base.MergeMapI32I32(this.mapHookTechEffect[cfgEnum.TechEffectType(cfgHookTech.TechEffectType)], cfgHookTechLevel.MapEffect)
		}

	}

	//通知英雄系统
	if _, ok := this.mapHookTechEffect[cfgEnum.TechEffectType_HeroProp]; ok {
		this.getPlayerHeroFun().updateCalcFightpower(true)
	}
}
