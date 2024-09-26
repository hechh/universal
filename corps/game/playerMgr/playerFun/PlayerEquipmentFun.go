package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	report2 "corps/common/report"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/module/entry"

	"github.com/golang/protobuf/proto"
)

type (
	PlayerEquipmentFun struct {
		PlayerFun
		pbData           pb.PBPlayerEquipment                   //基本信息
		maxPosCount      uint32                                 //最大格子数量
		usePosCount      uint32                                 //当前使用的格子数量
		mapEquipment     map[uint32]*PlayerEquipment            //装备数据
		mapEquipmentById map[uint32]map[uint32]*PlayerEquipment //装备数据
		mapHookEquipment map[uint32]*PlayerEquipment            //挂机装备数据
	}

	PlayerEquipment struct {
		*pb.PBEquipment //装备数据
	}
)

func (this *PlayerEquipmentFun) getMaxPosCount() uint32 {

	return this.maxPosCount
}
func (this *PlayerEquipmentFun) GetSpareBag() uint32 {
	if this.usePosCount > this.maxPosCount {
		plog.Info("GetSpareBag is full %d,%d", this.usePosCount, this.maxPosCount)
		return 0
	}
	return this.maxPosCount - this.usePosCount
}
func (this *PlayerEquipmentFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.mapEquipment = make(map[uint32]*PlayerEquipment)
	this.mapHookEquipment = make(map[uint32]*PlayerEquipment)
	this.mapEquipmentById = make(map[uint32]map[uint32]*PlayerEquipment)
	this.maxPosCount = 0
	this.usePosCount = 0
}

// 加载背包数据
func (this *PlayerEquipmentFun) Load(pData []byte) {
	proto.Unmarshal(pData, &this.pbData)
	this.loadData()
	this.UpdateSave(false)
}

func (this *PlayerEquipmentFun) loadData() {
	this.mapEquipment = make(map[uint32]*PlayerEquipment)
	this.mapEquipmentById = make(map[uint32]map[uint32]*PlayerEquipment)
	this.usePosCount = 0
	for i := 0; i < len(this.pbData.EquipmentList); i++ {
		pbPlayerEquipment := &PlayerEquipment{}
		pbPlayerEquipment.PBEquipment = this.pbData.EquipmentList[i]

		this.mapEquipment[pbPlayerEquipment.Sn] = pbPlayerEquipment

		if _, ok := this.mapEquipmentById[pbPlayerEquipment.Id]; !ok {
			this.mapEquipmentById[pbPlayerEquipment.Id] = make(map[uint32]*PlayerEquipment)
		}

		this.mapEquipmentById[pbPlayerEquipment.Id][pbPlayerEquipment.Sn] = pbPlayerEquipment

		//记录使用的格子
		if pbPlayerEquipment.EquipProfession == base.CFG_DEFAULT_VALUE {
			this.usePosCount++
		}
	}

	for i := 0; i < len(this.pbData.HookEquipmentList); i++ {
		pbPlayerEquipment := &PlayerEquipment{}
		pbPlayerEquipment.PBEquipment = this.pbData.HookEquipmentList[i]
		this.mapHookEquipment[pbPlayerEquipment.Sn] = pbPlayerEquipment
	}

	this.UpdateSave(true)
}

func (this *PlayerEquipmentFun) LoadComplete() {
	this.updateMaxPosCount()
}

// 保存
func (this *PlayerEquipmentFun) Save(bNow bool) {
	if !this.BSave {
		return
	}

	this.BSave = false

	pbData := &pb.PBPlayerEquipment{}
	this.SavePb(pbData)

	//通知db保存玩家数据
	buff, _ := proto.Marshal(pbData)
	cluster.SendToDb(&pb.RpcHead{Id: this.AccountId}, "DbPlayerMgr", "SavePlayerDB", this.PbType, buff, bNow)
}

func (this *PlayerEquipmentFun) saveData() {
	this.pbData.EquipmentList = make([]*pb.PBEquipment, 0)
	for _, v := range this.mapEquipment {
		this.pbData.EquipmentList = append(this.pbData.EquipmentList, v.PBEquipment)
	}

	this.pbData.HookEquipmentList = make([]*pb.PBEquipment, 0)
	for _, v := range this.mapHookEquipment {
		this.pbData.HookEquipmentList = append(this.pbData.HookEquipmentList, v.PBEquipment)
	}
}

// 保存
func (this *PlayerEquipmentFun) SavePb(pbData *pb.PBPlayerEquipment) {
	if pbData == nil {
		pbData = &pb.PBPlayerEquipment{}
	}
	this.saveData()
	base.DeepCopy(&this.pbData, pbData)
}

func (this *PlayerEquipmentFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.Equipment == nil {
		pbData.Equipment = new(pb.PBPlayerEquipment)
	}

	this.SavePb(pbData.Equipment)
}
func (this *PlayerEquipmentFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerEquipment{}
}

// 设置玩家数据
func (this *PlayerEquipmentFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}

	pbSystem := pbData.(*pb.PBPlayerEquipment)
	if pbSystem == nil {
		return false
	}

	this.pbData = *pbSystem
	this.loadData()
	return true
}

// 新玩家 需要初始化数据
func (this *PlayerEquipmentFun) NewPlayer() {
	//送初始装备

}

func (this *PlayerEquipmentFun) getOrderId() uint32 {
	this.pbData.OrderId++
	return this.pbData.OrderId
}
func (this *PlayerEquipmentFun) GetSplitAddBoxCount() uint32 {
	return this.pbData.SplitAddBoxCount
}

// 加装备
func (this *PlayerEquipmentFun) AddEquipments(head *pb.RpcHead, uId uint32, uCount uint32, uQuality uint32, uStar uint32, doingType pb.EmDoingType) cfgEnum.ErrorCode {
	uCode := cfgEnum.ErrorCode_Success
	arrList := make([]*pb.PBEquipment, 0)
	for i := uint32(0); i < uCount; i++ {
		uCode, pbEquip := this.AddEquipment(head, uId, uQuality, uStar, false, doingType)
		if uCode != cfgEnum.ErrorCode_Success {
			return uCode
		}

		if pbEquip == nil {
			continue
		}

		arrList = append(arrList, pbEquip)
	}

	this.EquipmentNotify(head, arrList, doingType == pb.EmDoingType_EDT_BattleHook, doingType)
	return uCode
}

// 获取一个新的装备结构
func (this *PlayerEquipmentFun) GetNewEquipment(uId uint32, uQuality uint32, uStar uint32, doingType pb.EmDoingType) *pb.PBEquipment {
	//从battle表中随机
	if uStar == 0 {
		uStar = 1
	}

	cfgEquipment := cfgData.GetCfgEquipment(uId)
	if nil == cfgEquipment {
		return nil
	}

	//获取品质表
	cfgQuality := cfgData.GetCfgEquipmentQuality(uQuality)
	if nil == cfgQuality {
		return nil
	}

	//获取星级表
	cfgStar := cfgData.GetCfgEquipmentStar(uStar)
	if nil == cfgStar {
		return nil
	}

	pEquipment := &pb.PBEquipment{
		Id:              uId,
		Quality:         uQuality,
		Star:            uStar,
		Sn:              0,
		MainProp:        &pb.PBEquipmentProp{},
		EquipProfession: base.CFG_DEFAULT_VALUE,
	}

	//获取主词条数值
	pEquipment.MainProp.PropId = cfgEquipment.MainProp
	uMainValueRate := uint32(0)
	if pEquipment.MainProp.PropId == 0 {
		arrcfgProp := cfgData.GetCfgRandEquipmentProp(cfgEquipment.Stage, uint32(cfgEnum.EHydraPropType_MainProp), cfgEquipment.PartType, uint32(1))
		if arrcfgProp != nil && len(arrcfgProp) == 1 {
			pEquipment.MainProp.Value = cfgEquipment.MainPropValue
			pEquipment.MainProp.PropId = arrcfgProp[0].PropId
			uMainValueRate = arrcfgProp[0].ValueRate
		} else {
			plog.Info("AddEquipment error")
		}
	} else {
		cfgProp := cfgData.GetCfgEquipmentProp(uint32(cfgEnum.EHydraPropType_MainProp), pEquipment.MainProp.PropId)
		pEquipment.MainProp.Value = cfgEquipment.MainPropValue
		uMainValueRate = cfgProp.ValueRate
	}

	if uMainValueRate <= 0 {
		uMainValueRate = base.MIL_PERCENT
	}
	//挂机装备
	uAddStarRate := cfgStar.AddRate
	if doingType == pb.EmDoingType_EDT_BattleHook || doingType == pb.EmDoingType_EDT_Offline {
		mapid, stageid := this.getPlayerSystemBattleHookFun().GetMapIdAndStageId()
		cfgHook := cfgData.GetCfgBattleHookStage(mapid, stageid)
		if cfgHook != nil && cfgHook.StarRate > 0 {
			uAddStarRate = uAddStarRate * cfgHook.StarRate / base.MIL_PERCENT
		}
	}

	//评分 (2000+ 2000*0.164 )*1.1 / 50
	pEquipment.MainProp.Score = base.RAND.RandUint(cfgQuality.MinScoreRate, cfgQuality.MaxScoreRate)
	pEquipment.MainProp.Value = uint32(uint64(pEquipment.MainProp.Value) * uint64(pEquipment.MainProp.Score) * uint64(uAddStarRate) / uint64(uMainValueRate) / uint64(base.MIL_PERCENT*base.MIL_PERCENT))

	//获取次词条数值
	pEquipment.MinorPropList = this.getEquipProp(cfgEquipment, cfgQuality, uAddStarRate, cfgEquipment.MinorProp, cfgQuality.MinorPropCount, cfgEnum.EHydraPropType_MinorProp, cfgEquipment.MinorPropValue)
	pEquipment.VicePropList = this.getEquipProp(cfgEquipment, cfgQuality, uAddStarRate, cfgEquipment.ViceProp, cfgQuality.VicePropCount, cfgEnum.EHydraPropType_ViceProp, cfgEquipment.VicePropValue)

	return pEquipment
}

// 加装备
func (this *PlayerEquipmentFun) AddEquipment(head *pb.RpcHead, uId uint32, uQuality uint32, uStar uint32, bSend bool, doingType pb.EmDoingType) (cfgEnum.ErrorCode, *pb.PBEquipment) {
	pEquipment := this.GetNewEquipment(uId, uQuality, uStar, doingType)
	if pEquipment == nil {
		return plog.Print(this.AccountId, cfgData.GetEquipmentErrorCode(uId), uId, uQuality, uStar, bSend, doingType), nil
	}

	uCode, pEquipment := this.AddPBEquipment(head, pEquipment, bSend, doingType)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode, uId, uQuality, uStar, bSend, doingType), nil
	}
	return uCode, pEquipment
}

func (this *PlayerEquipmentFun) addInnerEquipment(head *pb.RpcHead, pbEquipment *pb.PBEquipment, bSend bool, doingType pb.EmDoingType, isHook bool) {
	pbEquipment.Sn = this.getOrderId()

	if isHook {
		//更新
		this.updateHookEquipment(head, &PlayerEquipment{PBEquipment: pbEquipment}, bSend)

	} else {
		this.usePosCount++

		//更新
		this.updateEquipment(head, &PlayerEquipment{PBEquipment: pbEquipment}, bSend)
	}

	// 数据上报
	report2.Send(head, &report2.ReportAddItem{
		Kind:   uint32(cfgEnum.ESystemType_Equipment),
		Doing:  uint32(doingType),
		ItemID: pbEquipment.Id,
		Params: []uint32{pbEquipment.Quality, pbEquipment.Star},
		Sn:     pbEquipment.Sn,
	})
}

// 加装备
func (this *PlayerEquipmentFun) AddPBEquipment(head *pb.RpcHead, pbEquipment *pb.PBEquipment, bSend bool, doingType pb.EmDoingType) (cfgEnum.ErrorCode, *pb.PBEquipment) {
	if pbEquipment == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_AddEquipmentParam, bSend, doingType), nil
	}

	//成就触发
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_EquipmentReward, 1, pbEquipment.Quality)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_EquipmentReward, 1, base.CFG_DEFAULT_VALUE)

	//是否自动分解
	if base.ArrayContainsValue(this.pbData.AutoSplitQuality, pbEquipment.Quality) {
		this.InnerSplit(head, map[uint32]uint32{pbEquipment.Quality: 1}, true)
		return cfgEnum.ErrorCode_Success, nil
	}

	//挂机掉落 优先掉落到挂机背包中 只要20件红色品质以上的
	if doingType == pb.EmDoingType_EDT_BattleHook {
		this.getPlayerSystemBattleHookFun().AddHookEquipmentCount(1)
		uMinQuality := cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_MIN_HOOK_EQUIPMENT_QUALITY)
		mapId, stageId := this.getPlayerSystemBattleHookFun().GetFinishMapIdAndStageId()
		if mapId > 0 {
			cfgHook := cfgData.GetCfgBattleHookStage(mapId, stageId)
			if cfgHook != nil && cfgHook.EquipQuality > 0 {
				uMinQuality = base.MaxUint32(uMinQuality, cfgHook.EquipQuality-1)
			}
		}

		if uMinQuality > 0 && pbEquipment.Quality >= uMinQuality &&
			uint32(len(this.mapHookEquipment)) < cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_MAX_HOOK_EQUIPMENT) {
			this.addInnerEquipment(head, pbEquipment, bSend, doingType, true)
			return cfgEnum.ErrorCode_Success, pbEquipment
		}
	}

	//背包已满，自动分解 购买的不分解
	if this.usePosCount >= this.getMaxPosCount() && doingType != pb.EmDoingType_EDT_Charge && doingType != pb.EmDoingType_EDT_FirstCharge {
		this.InnerSplit(head, map[uint32]uint32{pbEquipment.Quality: 1}, true)
		return cfgEnum.ErrorCode_Success, nil
	}

	this.addInnerEquipment(head, pbEquipment, bSend, doingType, false)
	return cfgEnum.ErrorCode_Success, pbEquipment
}

// 更新装备数据
func (this *PlayerEquipmentFun) EquipmentNotify(head *pb.RpcHead, equipment []*pb.PBEquipment, isHook bool, doingType pb.EmDoingType) {
	cluster.SendToClient(head, &pb.EquipmentNotify{
		PacketHead: &pb.IPacket{},
		Equipment:  equipment,
		IsHook:     isHook,
		DoingType:  doingType,
	}, cfgEnum.ErrorCode_Success)
}

// 获取装备词条 先取固定 再随机剩余
func (this *PlayerEquipmentFun) getEquipProp(cfgEquip *cfgData.EquipmentCfg, cfgQuality *cfgData.EquipmentQualityCfg, uAddStarRate uint32, arrProp []uint32, uPropCount uint32, emPropType cfgEnum.EHydraPropType, uPropValue uint32) []*pb.PBEquipmentProp {
	arrReturn := make([]*pb.PBEquipmentProp, 0)
	if len(arrProp) > 0 {
		for i := 0; i < len(arrProp); i++ {
			if i >= int(uPropCount) {
				break
			}

			cfgProp := cfgData.GetCfgEquipmentProp(uint32(emPropType), arrProp[i])
			if cfgProp == nil {
				continue
			}

			pbProp := &pb.PBEquipmentProp{
				PropId: arrProp[i],
				Score:  base.RAND.RandUint(cfgQuality.MinScoreRate, cfgQuality.MaxScoreRate),
			}

			//无价值 取值
			if emPropType == cfgEnum.EHydraPropType_ViceProp {
				pbProp.Value = uPropValue
				pbProp.Value = pbProp.Value * pbProp.Score / (cfgProp.ValueRate * 10000)
			} else {
				var uValue uint64 = uint64(uPropValue) * uint64(pbProp.Score) * uint64(uAddStarRate)
				uValue = uValue / uint64(cfgProp.ValueRate) / 10000 / 10000
				pbProp.Value = uint32(uValue)
			}

			if pbProp.Value <= 0 {
				pbProp.Value = 1
			}

			arrReturn = append(arrReturn, pbProp)
		}

		uPropCount -= uint32(len(arrReturn))
	}

	if uPropCount > 0 {
		arrcfgProp := cfgData.GetCfgRandEquipmentProp(cfgEquip.Stage, uint32(emPropType), cfgEquip.PartType, uPropCount)
		if arrcfgProp != nil {
			for i := 0; i < len(arrcfgProp); i++ {
				pbProp := &pb.PBEquipmentProp{
					PropId: arrcfgProp[i].PropId,
					Score:  base.RAND.RandUint(cfgQuality.MinScoreRate, cfgQuality.MaxScoreRate),
				}

				//无价值 取值
				if emPropType == cfgEnum.EHydraPropType_ViceProp {
					pbProp.Value = uPropValue
					pbProp.Value = pbProp.Value * pbProp.Score / (arrcfgProp[i].ValueRate * 10000)
				} else {
					var uValue uint64 = uint64(uPropValue) * uint64(pbProp.Score) * uint64(uAddStarRate)
					uValue = uValue / uint64(arrcfgProp[i].ValueRate) / 10000 / 10000
					pbProp.Value = uint32(uValue)
				}

				if pbProp.Value <= 0 {
					pbProp.Value = 1
				}

				arrReturn = append(arrReturn, pbProp)
			}
		}
	}

	return arrReturn
}

// 装备分解请求
func (this *PlayerEquipmentFun) EquipmentSplitRequest(head *pb.RpcHead, pbRequest *pb.EquipmentSplitRequest) {
	uCode := this.EquipmentSplit(head, pbRequest.SnList)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.EquipmentSplitResponse{
			PacketHead: &pb.IPacket{},
			SnList:     pbRequest.SnList,
		}, uCode)
	}
}

func (this *PlayerEquipmentFun) getEquipment(uSn uint32) *PlayerEquipment {
	if uSn <= 0 {
		return nil
	}
	pEquip, ok := this.mapEquipment[uSn]
	if !ok {
		return nil
	}

	return pEquip
}

func (this *PlayerEquipmentFun) updateEquipment(head *pb.RpcHead, pEquipment *PlayerEquipment, bSend bool) {
	if pEquipment == nil {
		return
	}

	this.mapEquipment[pEquipment.Sn] = pEquipment

	this.UpdateSave(true)

	if bSend {
		this.EquipmentNotify(head, []*pb.PBEquipment{pEquipment.PBEquipment}, false, pb.EmDoingType_EDT_Other)
	}
}
func (this *PlayerEquipmentFun) updateHookEquipment(head *pb.RpcHead, pEquipment *PlayerEquipment, bSend bool) {
	if pEquipment == nil {
		return
	}

	this.mapHookEquipment[pEquipment.Sn] = pEquipment

	if bSend {
		this.EquipmentNotify(head, []*pb.PBEquipment{pEquipment.PBEquipment}, true, pb.EmDoingType_EDT_Other)
	}
}

// 穿戴
func (this *PlayerEquipmentFun) Equip(head *pb.RpcHead, uProf uint32, uPart uint32, uEquipSn uint32, uUnEquipSn uint32) cfgEnum.ErrorCode {
	//判断装备是否已满
	if uUnEquipSn > 0 {
		uNeedSpare := uint32(1)
		if uEquipSn > 0 {
			uNeedSpare = 0
		}

		if this.GetSpareBag() < uNeedSpare {
			cluster.SendCommonToClient(head, cfgEnum.ErrorCode_BagFull)
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_BagFull, uProf, uPart, uEquipSn, uUnEquipSn)
		}
	}

	uCode := cfgEnum.ErrorCode_Success
	for {
		if uEquipSn > 0 {
			pEquip := this.getEquipment(uEquipSn)
			if pEquip == nil {
				uCode = cfgEnum.ErrorCode_EquipmentNotFound
				break
			}

			if !base.IsDefaultCfg(pEquip.EquipProfession) {
				uCode = cfgEnum.ErrorCode_EquipProfessionDefaultValue
				break
			}

			//职业限制
			cfgEquip := cfgData.GetCfgEquipment(pEquip.Id)
			if cfgEquip.PartType != uPart {
				uCode = cfgEnum.ErrorCode_EquipmentPartInconsistent
				break
			}

			//部位限制
			if !base.IsDefaultCfg(cfgEquip.ProfType) && cfgEquip.ProfType != uProf {
				uCode = cfgEnum.ErrorCode_EquipmentProfInconsistent
				break
			}

			//职业等级限制
			cfgStar := cfgData.GetCfgEquipmentStar(pEquip.Star)
			if cfgStar == nil {
				uCode = plog.Print(head.Id, cfgData.GetEquipmentStarErrorCode(pEquip.Star), pEquip.Star)
				break
			}

			if cfgStar.NeedProfLevel > 0 && cfgStar.NeedProfLevel > this.getPlayerSystemProfessionFun().GetProf(uProf).Level {
				uCode = plog.Print(head.Id, cfgEnum.ErrorCode_NeedProfLevel)
				break
			}

			//脱下
			pEquip.EquipProfession = uProf
			this.usePosCount--
			this.updateEquipment(head, pEquip, false)
		}

		if uUnEquipSn > 0 {
			pUnEquip := this.getEquipment(uUnEquipSn)
			if pUnEquip == nil {
				break
			}

			if base.IsDefaultCfg(pUnEquip.EquipProfession) {
				break
			}

			//脱下
			pUnEquip.EquipProfession = base.CFG_DEFAULT_VALUE
			this.usePosCount++
			this.updateEquipment(head, pUnEquip, false)
		}

		break
	}

	return uCode
}

func (this *PlayerEquipmentFun) GetTotalSplitScore() uint64 {
	return this.pbData.TotalSplitScore
}

// 返回分解获取的积分
func (this *PlayerEquipmentFun) SplitEquipment(head *pb.RpcHead, bAuto bool, args ...*pb.PBAddItemData) uint64 {
	if len(args) <= 0 {
		return 0
	}
	oldSplitScore := this.pbData.TotalSplitScore
	// 统计品质
	mapQualityCount := make(map[uint32]uint32)
	for _, item := range args {
		mapQualityCount[item.Equipment.Quality] += 1
	}
	// 装备分解
	this.InnerSplit(head, mapQualityCount, bAuto)
	return this.pbData.TotalSplitScore - oldSplitScore
}

// 内部分解
func (this *PlayerEquipmentFun) InnerSplit(head *pb.RpcHead, mapQuality map[uint32]uint32, bAuto bool) {
	uAddScore := uint32(0)
	uCount := uint32(0)
	for k, v := range mapQuality {
		cfgQuality := cfgData.GetCfgEquipmentQuality(k)
		if cfgQuality == nil {
			continue
		}

		uCount += v
		uAddScore += v * cfgQuality.SplitAddScore

		//成就触发
		this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_SpliteQualityEquip, v, k)
	}

	this.pbData.SplitEquipCount += uCount

	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_SpliteQualityEquip, uCount, uint32(cfgEnum.EQuality_Any))

	//加道具
	this.addSplitScore(head, uAddScore, bAuto)

	this.UpdateSave(true)

}
func (this *PlayerEquipmentFun) ClearSplitEquipCount() {
	this.pbData.SplitEquipCount = 0
	this.UpdateSave(true)
}

// 增加分解积分
func (this *PlayerEquipmentFun) addSplitScore(head *pb.RpcHead, uAdd uint32, bAuto bool) {
	this.pbData.SplitScore += uAdd
	this.pbData.TotalSplitScore += uint64(uAdd)

	//手动分解 如果溢出给宝箱
	if !bAuto {
		arrAddItem := make([]*pb.PBAddItemData, 0)
		// 词条效果加成 每回收多少次加宝箱
		val := entry.ToValue(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_SplitEquipmentRewardBox), 1201)...)
		if val > 0 {
			uAddBoxCount := this.pbData.SplitEquipCount / val
			if uAddBoxCount > 0 {
				arrAddItem = this.getPlayerBagFun().GetPbItems([]*common.ItemInfo{
					{Kind: uint32(cfgEnum.ESystemType_Item), Id: cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_CFG_SPLITSCORE_BOXID), Count: int64(uAddBoxCount)}},
					pb.EmDoingType_EDT_Entry)
				this.pbData.SplitEquipCount = this.pbData.SplitEquipCount % val
			}
		}
		uMaxScore := cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_CFG_MAX_SPLITSCORE)
		if uMaxScore > 0 && this.pbData.SplitScore >= uMaxScore {
			uAddCount := this.pbData.SplitScore / uMaxScore
			this.pbData.SplitScore -= uAddCount * uMaxScore

			this.pbData.SplitAddBoxCount += uAddCount
			uAddBigBox := this.pbData.SplitAddBoxCount / 10
			if uAddBigBox > 0 {
				arrAddItem = append(arrAddItem, this.getPlayerBagFun().GetPbItems([]*common.ItemInfo{
					{Kind: uint32(cfgEnum.ESystemType_Item), Id: cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_CFG_SPLITSCORE_BIGBOXID), Count: int64(uAddBigBox)}},
					pb.EmDoingType_EDT_EquipSplit)...)
			}
			this.pbData.SplitAddBoxCount = this.pbData.SplitAddBoxCount % 10

			if uAddCount > uAddBigBox {
				uAddSmallCount := uAddCount - uAddBigBox
				arrAddItem = append(arrAddItem, this.getPlayerBagFun().GetPbItems([]*common.ItemInfo{
					{Kind: uint32(cfgEnum.ESystemType_Item), Id: cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_CFG_SPLITSCORE_BOXID), Count: int64(uAddSmallCount)}},
					pb.EmDoingType_EDT_EquipSplit)...)
			}

			//恭喜获得
			if len(arrAddItem) > 0 {
				this.getPlayerBagFun().AddPbItems(head, arrAddItem, pb.EmDoingType_EDT_EquipSplit, true)
			}

			//成就统计
			this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_AddSpliteBoxCount, uAddCount)
		}
	}

	cluster.SendToClient(head, &pb.EquipmentSplitScoreNotify{
		PacketHead: &pb.IPacket{},
		SplitScore: this.pbData.SplitScore,
	}, cfgEnum.ErrorCode_Success)

}

// 装备分解请求 可以为空
func (this *PlayerEquipmentFun) EquipmentSplit(head *pb.RpcHead, arrSn []uint32) cfgEnum.ErrorCode {
	arrDelSn := make([]uint32, 0)

	mapQuality := make(map[uint32]uint32)
	for i := 0; i < len(arrSn); i++ {
		pEquipment := this.getEquipment(arrSn[i])
		if pEquipment == nil {
			continue
		}

		//穿戴中不能被分解
		if !base.IsDefaultCfg(pEquipment.EquipProfession) {
			continue
		}

		//加锁不能分解
		if pEquipment.IsLock {
			continue
		}

		mapQuality[pEquipment.Quality] += 1

		arrDelSn = append(arrDelSn, arrSn[i])
	}

	for i := 0; i < len(arrDelSn); i++ {
		pEquipment := this.getEquipment(arrDelSn[i])
		this.delEquipment(head, pEquipment)
	}

	//给分解道具
	this.InnerSplit(head, mapQuality, false)

	cluster.SendToClient(head, &pb.EquipmentSplitResponse{
		PacketHead: &pb.IPacket{},
		SnList:     arrDelSn,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 装备格子购买请求
func (this *PlayerEquipmentFun) EquipmentBuyPosRequest(head *pb.RpcHead, pbRequest *pb.EquipmentBuyPosRequest) {
	uCode := this.EquipmentBuyPos(head, pbRequest.CurPosBuyCount)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.EquipmentBuyPosResponse{
			PacketHead:  &pb.IPacket{},
			PosBuyCount: pbRequest.CurPosBuyCount,
		}, uCode)
	}
}

// 装备格子购买请求
func (this *PlayerEquipmentFun) EquipmentBuyPos(head *pb.RpcHead, uCurCount uint32) cfgEnum.ErrorCode {
	if uCurCount != this.pbData.PosBuyCount {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_EquipmentPosBuyInconsistent, uCurCount, this.pbData.PosBuyCount)
	}
	//扣道具
	cfgPosBuy := cfgData.GetCfgEquipmentBuyPos(this.pbData.PosBuyCount)
	if cfgPosBuy == nil {
		return plog.Print(this.AccountId, cfgData.GetEquipmentBuyPosErrorCode(this.pbData.PosBuyCount), uCurCount, this.pbData.PosBuyCount)
	}

	uCode := this.getPlayerBagFun().DelItem(head, cfgPosBuy.ArrNeedItem.Kind, cfgPosBuy.ArrNeedItem.Id, cfgPosBuy.ArrNeedItem.Count, pb.EmDoingType_EDT_EquipSplit)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode, uCurCount, this.pbData.PosBuyCount, *cfgPosBuy)
	}

	this.pbData.PosBuyCount++
	this.updateMaxPosCount()
	this.UpdateSave(true)

	cluster.SendToClient(head, &pb.EquipmentBuyPosResponse{
		PacketHead:  &pb.IPacket{},
		PosBuyCount: this.pbData.PosBuyCount,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 计算格子最大数量
func (this *PlayerEquipmentFun) updateMaxPosCount() {
	entryCount := entry.ToValue(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_EquipmentBagSize), uint32(cfgEnum.EntryWorkTag_None))...)

	this.maxPosCount = cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_MAX_EQUIPMENT)
	this.maxPosCount += this.pbData.PosBuyCount * 5
	this.maxPosCount += entryCount
}

// 装备自动分解请求
func (this *PlayerEquipmentFun) EquipmentAutoSplitRequest(head *pb.RpcHead, pbRequest *pb.EquipmentAutoSplitRequest) {
	uCode := this.EquipmentAutoSplit(head, pbRequest.QualityList)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.EquipmentAutoSplitResponse{
			PacketHead:  &pb.IPacket{},
			QualityList: pbRequest.QualityList,
		}, uCode)
	}
}

func (this *PlayerEquipmentFun) GetAutoSplitQuality() map[uint32]struct{} {
	tmp := make(map[uint32]struct{})
	for _, qq := range this.pbData.AutoSplitQuality {
		tmp[qq] = struct{}{}
	}
	return tmp
}

// 装备格子购买请求
func (this *PlayerEquipmentFun) EquipmentAutoSplit(head *pb.RpcHead, arrQuality []uint32) cfgEnum.ErrorCode {
	this.pbData.AutoSplitQuality = arrQuality
	this.UpdateSave(true)

	cluster.SendToClient(head, &pb.EquipmentAutoSplitResponse{
		PacketHead:  &pb.IPacket{},
		QualityList: arrQuality,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 装备锁定请求
func (this *PlayerEquipmentFun) EquipmentLockRequest(head *pb.RpcHead, pbRequest *pb.EquipmentLockRequest) {
	uCode := this.EquipmentLock(head, pbRequest.Sn, pbRequest.IsLock)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.EquipmentLockResponse{
			PacketHead: &pb.IPacket{},
			Sn:         pbRequest.Sn,
			IsLock:     pbRequest.IsLock,
		}, uCode)
	}
}

// 装备锁定请求
func (this *PlayerEquipmentFun) EquipmentLock(head *pb.RpcHead, Sn uint32, IsLock bool) cfgEnum.ErrorCode {
	pEquip := this.getEquipment(Sn)
	if pEquip == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_EquipmentSnNotFound, Sn, IsLock)
	}

	if pEquip.IsLock != IsLock {
		pEquip.IsLock = IsLock
		this.updateEquipment(head, pEquip, false)
		this.UpdateSave(true)
	}

	cluster.SendToClient(head, &pb.EquipmentLockResponse{
		PacketHead: &pb.IPacket{},
		Sn:         Sn,
		IsLock:     IsLock,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 装备删除
func (this *PlayerEquipmentFun) delEquipment(head *pb.RpcHead, pEquip *PlayerEquipment) {
	if pEquip == nil {
		return
	}

	//需要脱下
	if !base.IsDefaultCfg(pEquip.EquipProfession) {
		cfgEquip := cfgData.GetCfgEquipment(pEquip.Id)
		if cfgEquip == nil {
			plog.Print(this.AccountId, cfgData.GetEquipmentErrorCode(pEquip.Id), *pEquip)
			return
		}

		this.getPlayerSystemProfessionFun().ProfessionEquip(head, pEquip.EquipProfession, cfgEquip.PartType, 0)
	}

	delete(this.mapEquipmentById[pEquip.Id], pEquip.Sn)
	delete(this.mapEquipment, pEquip.Sn)

	this.usePosCount--
	this.UpdateSave(true)
}

// 删除所有装备
func (this *PlayerEquipmentFun) DipDelAllEquipment(head *pb.RpcHead) cfgEnum.ErrorCode {
	if len(this.mapEquipment) <= 0 {
		return cfgEnum.ErrorCode_Success
	}

	//先脱去装备
	for _, pEquip := range this.mapEquipment {
		if !base.IsDefaultCfg(pEquip.EquipProfession) {
			cfgEquip := cfgData.GetCfgEquipment(pEquip.Id)
			if cfgEquip == nil {
				return plog.Print(this.AccountId, cfgData.GetEquipmentErrorCode(pEquip.Id), *pEquip)
			}

			this.getPlayerSystemProfessionFun().ProfessionEquip(head, pEquip.EquipProfession, cfgEquip.PartType, 0)
		}
	}

	this.mapEquipmentById = make(map[uint32]map[uint32]*PlayerEquipment)
	this.mapEquipment = make(map[uint32]*PlayerEquipment)
	this.usePosCount = 0
	this.UpdateSave(true)
	return cfgEnum.ErrorCode_Success
}

// 挂机装备领取请求
func (this *PlayerEquipmentFun) HookEquipmentAwardRequest(head *pb.RpcHead, pbRequest *pb.HookEquipmentAwardRequest) {
	uCode := this.HookEquipmentAward(head, pbRequest.Sn)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.HookEquipmentAwardResponse{
			PacketHead: &pb.IPacket{},
			Sn:         pbRequest.Sn,
		}, uCode)
	}
}

// 装备锁定请求
func (this *PlayerEquipmentFun) HookEquipmentAward(head *pb.RpcHead, Sn uint32) cfgEnum.ErrorCode {
	pEquip, ok := this.mapHookEquipment[Sn]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_EquipmentSnNotFound, Sn)
	}

	//判断背包已满
	if this.usePosCount >= this.getMaxPosCount() {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_BagFull, Sn)
	}

	this.addInnerEquipment(head, pEquip.PBEquipment, false, pb.EmDoingType_EDT_BattleHook, false)
	delete(this.mapHookEquipment, Sn)

	cluster.SendToClient(head, &pb.HookEquipmentAwardResponse{
		PacketHead: &pb.IPacket{},
		Sn:         pEquip.PBEquipment.Sn,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}
