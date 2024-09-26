package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	common2 "corps/common"
	"corps/common/cfgData"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/module/entry"
	"corps/server/game/module/reward"
	"sort"

	"github.com/golang/protobuf/proto"
)

type (
	ProfLink struct {
		uLinkStar   uint32
		uLinkLevel  uint32
		uLinkRefine uint32
	}

	PlayerSystemProfessionFun struct {
		PlayerFun
		arrProf     []*PlayerProfInfo    //内存数据
		mapProfLink map[uint32]*ProfLink //共鸣数据
		maxProfLink *ProfLink            //最大共鸣数据
	}

	PlayerProfInfo struct {
		*pb.PBPlayerSystemProfInfo
		mapProp map[uint32]float64
	}
)

func (this *PlayerSystemProfessionFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.mapProfLink = make(map[uint32]*ProfLink)
	this.maxProfLink = &ProfLink{}
}

// 加载完成
func (this *PlayerSystemProfessionFun) LoadPlayerDBFinish() {
	if len(this.arrProf) > 0 {
		return
	}

	//初始化五个职业
	for i := cfgEnum.EHydraProf_Tank; i <= cfgEnum.EHydraProf_Machinist; i++ {
		pbProp := &PlayerProfInfo{
			PBPlayerSystemProfInfo: &pb.PBPlayerSystemProfInfo{
				ProfType:  uint32(i),
				Level:     1,
				Grade:     0,
				PeakLevel: 0,
			},
			mapProp: make(map[uint32]float64),
		}

		for j := cfgEnum.EHydraEquipment_Weapon; j <= cfgEnum.EHydraEquipment_Ring; j++ {
			pbProp.PartList = append(pbProp.PartList, &pb.PBPlayerSystemProfPartInfo{
				Part:       uint32(j),
				Level:      1,
				EquipSn:    0,
				Refine:     1,
				RefineTupo: 0,
			})
		}

		this.arrProf = append(this.arrProf, pbProp)

	}

	//初始化等级
	this.UpdateSave(true)
}

func (this *PlayerSystemProfessionFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = new(pb.PBPlayerSystem)
	}

	this.SaveSystem(pbData.System)
}

// 从数据库中加载
func (this *PlayerSystemProfessionFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	this.BSave = false
	if pbSystem.Prof == nil || pbSystem.Prof.ProfList == nil {
		return
	}

	this.loadData(pbSystem.Prof)
	this.UpdateSave(false)
}
func (this *PlayerSystemProfessionFun) loadData(pbData *pb.PBPlayerSystemProfession) {
	this.arrProf = make([]*PlayerProfInfo, 0)
	for i := 0; i < len(pbData.ProfList); i++ {
		this.arrProf = append(this.arrProf, &PlayerProfInfo{
			PBPlayerSystemProfInfo: pbData.ProfList[i],
			mapProp:                make(map[uint32]float64),
		})
	}
	this.UpdateSave(true)
}
func (this *PlayerSystemProfessionFun) LoadComplete() {
	for i := 0; i < len(this.arrProf); i++ {
		this.CalcProp(uint32(i))
	}
}
func (this *PlayerSystemProfessionFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemProfession{}
}

// 设置玩家数据
func (this *PlayerSystemProfessionFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}

	pbSystem := pbData.(*pb.PBPlayerSystemProfession)
	if pbSystem == nil {
		return false
	}

	this.loadData(pbSystem)
	return true
}

// 存储到数据库
func (this *PlayerSystemProfessionFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	if pbSystem.Prof == nil {
		pbSystem.Prof = new(pb.PBPlayerSystemProfession)
	}
	for i := 0; i < len(this.arrProf); i++ {
		pbSystem.Prof.ProfList = append(pbSystem.Prof.ProfList, this.arrProf[i].PBPlayerSystemProfInfo)
	}

	return this.BSave
}

// 获取职业等级
func (this *PlayerSystemProfessionFun) GetProfLevel(uType uint32, isUp bool) uint32 {
	uLevel := uint32(0)
	if cfgEnum.EHydraProf(uType) == cfgEnum.EHydraProf_Any {
		for i := 0; i < len(this.arrProf); i++ {
			uLevel += this.arrProf[i].Level
			if isUp && this.arrProf[i].Level > 0 {
				uLevel--
			}

		}
	} else {
		if uType < uint32(len(this.arrProf)) {
			uLevel = this.arrProf[uType].Level
			if isUp && this.arrProf[uType].Level > 0 {
				uLevel--
			}
		}
	}

	return uLevel
}

// 获取职业突破等级
func (this *PlayerSystemProfessionFun) GetProfGrade(uType uint32) uint32 {
	uLevel := uint32(0)
	if cfgEnum.EHydraProf(uType) == cfgEnum.EHydraProf_Any {
		for i := 0; i < len(this.arrProf); i++ {
			uLevel += this.arrProf[i].Grade
		}
	} else {
		uLevel = this.arrProf[uType].Grade
	}

	return uLevel
}

// 获取职业部位等级 升级 等级要减去1
func (this *PlayerSystemProfessionFun) GetProfPartLevel(uType uint32, bUp bool) uint32 {
	uLevel := uint32(0)
	if cfgEnum.EHydraProf(uType) == cfgEnum.EHydraProf_Any {
		for i := 0; i < len(this.arrProf); i++ {
			for j := 0; j < len(this.arrProf[i].PartList); j++ {
				uLevel += this.arrProf[i].PartList[j].Level
				if bUp && this.arrProf[i].PartList[j].Level > 0 {
					uLevel--
				}
			}
		}
	} else {
		if uType < uint32(len(this.arrProf)) {
			for j := 0; j < len(this.arrProf[uType].PartList); j++ {
				uLevel += this.arrProf[uType].PartList[j].Level
				if bUp && this.arrProf[uType].PartList[j].Level > 0 {
					uLevel--
				}
			}
		}
	}

	return uLevel
}

// 获取职业最大等级
func (this *PlayerSystemProfessionFun) GetProfMaxLevel() uint32 {
	uMax := uint32(0)
	for i := 0; i < len(this.arrProf); i++ {
		if uMax < this.arrProf[i].Level {
			uMax = this.arrProf[i].Level
		}
	}

	return uMax
}

// 获取最大职业穿戴装备数量
func (this *PlayerSystemProfessionFun) GetMaxProfEquipCount() uint32 {
	mapCount := make(map[uint32]uint32)
	for i := 0; i < len(this.arrProf); i++ {
		for j := 0; j < len(this.arrProf[i].PartList); j++ {
			if this.arrProf[i].PartList[j].EquipSn > 0 {
				mapCount[uint32(i)] += 1
			}
		}
	}

	uMaxCount := uint32(0)
	for _, count := range mapCount {
		if count > uMaxCount {
			uMaxCount = count
		}
	}

	return uMaxCount
}

// 获取职业穿戴装备数量
func (this *PlayerSystemProfessionFun) GetProfEquipCount(uType uint32) uint32 {
	uCount := uint32(0)
	if cfgEnum.EHydraProf(uType) == cfgEnum.EHydraProf_Any {
		for i := 0; i < len(this.arrProf); i++ {
			for j := 0; j < len(this.arrProf[i].PartList); j++ {
				if this.arrProf[i].PartList[j].EquipSn > 0 {
					uCount++
				}
			}
		}
	} else {
		if uType < uint32(len(this.arrProf)) {
			for j := 0; j < len(this.arrProf[uType].PartList); j++ {
				if this.arrProf[uType].PartList[j].EquipSn > 0 {
					uCount++
				}
			}
		}
	}

	return uCount
}
func (this *PlayerSystemProfessionFun) Test(head *pb.RpcHead, pbRequest proto.Message) {

}

// 职业升级请求
func (this *PlayerSystemProfessionFun) ProfessionLevelRequest(head *pb.RpcHead, pbRequest *pb.ProfessionLevelRequest) {
	uCode := this.ProfessionLevel(head, pbRequest.ProfType, pbRequest.CurLevel, pbRequest.AddLevel)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ProfessionLevelResponse{
			PacketHead: &pb.IPacket{},
			ProfType:   pbRequest.ProfType,
			CurLevel:   pbRequest.CurLevel,
		}, uCode)
	}
}

// 职业升级
func (this *PlayerSystemProfessionFun) ProfessionLevel(head *pb.RpcHead, uProfType uint32, uCurLevel uint32, uAddLevel uint32) cfgEnum.ErrorCode {
	if uAddLevel <= 0 || uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotSupported, uProfType, uAddLevel)
	}

	stProf := this.arrProf[uProfType]
	if uCurLevel != stProf.Level {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionLevelParam, uProfType, uCurLevel, stProf.Level)
	}

	//巅峰等级激活不能升级
	if stProf.PeakLevel > 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxLevel, uProfType, uCurLevel, uAddLevel, stProf.PeakLevel)
	}

	if uCurLevel+uAddLevel > cfgData.GetCfgProfMaxLevel() {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxLevel, uProfType, uCurLevel, uAddLevel, cfgData.GetCfgProfMaxLevel())
	}

	//限制突破等级
	cfgGrade := cfgData.GetCfgProfGrade(stProf.Grade)
	if cfgGrade == nil {
		return plog.Print(this.AccountId, cfgData.GetProfGradeErrorCode(stProf.Grade), uProfType, uCurLevel, uAddLevel, stProf.Grade)
	}
	if uCurLevel+uAddLevel > cfgGrade.MaxLevel {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxLevel, uProfType, uCurLevel, uAddLevel, cfgGrade.MaxLevel)
	}

	uRealAdd := uint32(0)
	mapNeedItem := make(map[uint32]map[uint32]int64)
	for i := uint32(0); i < uAddLevel; i++ {
		cfgLevel := cfgData.GetCfgProfLevel(uCurLevel + i)
		if cfgLevel == nil {
			return plog.Print(this.AccountId, cfgData.GetProfLevelErrorCode(uCurLevel+i), uProfType, uCurLevel, uAddLevel, uCurLevel+i)
		}
		for _, v := range cfgLevel.ListNeedItem {
			if _, ok := mapNeedItem[v.Kind]; !ok {
				mapNeedItem[v.Kind] = make(map[uint32]int64)
			}

			mapNeedItem[v.Kind][v.Id] += v.Count
		}

		uRealAdd++
	}

	listNeedItem := make([]*common2.ItemInfo, 0)
	for k, v := range mapNeedItem {
		for sk, sv := range v {
			// 词条效果
			listNeedItem = append(listNeedItem, &common2.ItemInfo{
				Kind:  k,
				Id:    sk,
				Count: sv,
			})
		}
	}

	// 消耗减免
	items := entry.ToItem(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_ConsumeReduce), uProfType)...)
	listNeedItem = reward.SubProbReward(items, listNeedItem...)

	//扣道具
	uCode := this.getPlayerBagFun().DelArrItem(head, listNeedItem, pb.EmDoingType_EDT_ProfessionLevel)
	if uCode != cfgEnum.ErrorCode_Success {
		return uCode
	}

	//等级
	stProf.Level += uRealAdd
	this.arrProf[uProfType] = stProf

	this.onProfLevel(head, uProfType, uRealAdd)
	this.UpdateSave(true)

	cluster.SendToClient(head, &pb.ProfessionLevelResponse{
		PacketHead: &pb.IPacket{},
		ProfType:   uProfType,
		CurLevel:   stProf.Level,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 职业突破请求
func (this *PlayerSystemProfessionFun) ProfessionGradeRequest(head *pb.RpcHead, pbRequest *pb.ProfessionGradeRequest) {
	uCode := this.ProfessionGrade(head, pbRequest.ProfType)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ProfessionGradeResponse{
			PacketHead: &pb.IPacket{},
			ProfType:   pbRequest.ProfType,
		}, uCode)
	}
}

// 获取总突破等级
func (this *PlayerSystemProfessionFun) GetAllGrade() uint32 {
	uTotal := uint32(0)
	for _, v := range this.arrProf {
		uTotal += v.Grade
	}

	return uTotal
}

// 获取属性
func (this *PlayerSystemProfessionFun) GetProf(uProfType uint32) *PlayerProfInfo {
	if uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return nil
	}

	if this.arrProf == nil || uProfType >= uint32(len(this.arrProf)) {
		return nil
	}

	return this.arrProf[uProfType]
}

// 职业升级
func (this *PlayerSystemProfessionFun) ProfessionGrade(head *pb.RpcHead, uProfType uint32) cfgEnum.ErrorCode {
	if uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotSupported, uProfType)
	}

	stProf := this.arrProf[uProfType]

	cfgGrade := cfgData.GetCfgProfGrade(stProf.Grade)
	if cfgGrade == nil {
		return plog.Print(this.AccountId, cfgData.GetProfGradeErrorCode(stProf.Grade), uProfType, stProf.Grade)
	}
	//判断等级是否达到突破等级
	if cfgGrade.MaxLevel != stProf.Level {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotMaxLevel, uProfType, cfgGrade.MaxLevel, stProf.Level)
	}

	//限制总等级
	if cfgGrade.NeedTotalLevel > this.GetProfLevel(uint32(cfgEnum.EHydraProf_Any), false) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxLevel, uProfType, cfgGrade.NeedTotalLevel, this.GetProfLevel(uint32(cfgEnum.EHydraProf_Any), false))
	}

	if cfgData.GetCfgProfGrade(stProf.Grade+1) == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxLevel, uProfType, stProf.Grade)
	}

	//扣道具
	uCode := this.getPlayerBagFun().DelArrItem(head, cfgGrade.ListNeedItem, pb.EmDoingType_EDT_ProfessionLevel)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode, uProfType, cfgGrade.ListNeedItem)
	}

	//存库
	stProf.Grade++
	this.arrProf[uProfType] = stProf
	this.UpdateSave(true)

	//成就统计
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_PofessionTupoCount, 1, uint32(cfgEnum.EHydraProf_Any))
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_PofessionTupoCount, 1, uProfType)

	//计算属性
	this.calcFightPower(uProfType)

	cluster.SendToClient(head, &pb.ProfessionGradeResponse{
		PacketHead: &pb.IPacket{},
		ProfType:   uProfType,
		Grade:      stProf.Grade,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 巅峰等级突破请求
func (this *PlayerSystemProfessionFun) ProfessionPeakRequest(head *pb.RpcHead, pbRequest *pb.ProfessionPeakRequest) {
	uCode := this.ProfessionPeak(head, pbRequest.ProfType)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ProfessionPeakResponse{
			PacketHead: &pb.IPacket{},
			ProfType:   pbRequest.ProfType,
			PeakLevel:  0,
		}, uCode)
	}
}

// 巅峰等级突破
func (this *PlayerSystemProfessionFun) ProfessionPeak(head *pb.RpcHead, uProfType uint32) cfgEnum.ErrorCode {
	if uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotSupported, uProfType)
	}

	stProf := this.arrProf[uProfType]

	if stProf.PeakLevel > 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionPeakLevelHasFinished, stProf.PeakLevel)
	}

	//是否最大等级
	if stProf.Level < cfgData.GetCfgProfMaxLevel() {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotMaxLevel, stProf.Level, cfgData.GetCfgProfMaxLevel())
	}

	stProf.PeakLevel = 1
	this.arrProf[uProfType] = stProf

	this.UpdateSave(true)

	cluster.SendToClient(head, &pb.ProfessionPeakResponse{
		PacketHead: &pb.IPacket{},
		ProfType:   uProfType,
		PeakLevel:  stProf.PeakLevel,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 职业巅峰升级请求
func (this *PlayerSystemProfessionFun) ProfessionPeakLevelRequest(head *pb.RpcHead, pbRequest *pb.ProfessionPeakLevelRequest) {
	uCode := this.ProfessionPeakLevel(head, pbRequest.ProfType, pbRequest.CurLevel)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ProfessionPeakLevelResponse{
			PacketHead: &pb.IPacket{},
			ProfType:   pbRequest.ProfType,
			CurLevel:   pbRequest.CurLevel,
		}, uCode)
	}
}

// 职业巅峰升级
func (this *PlayerSystemProfessionFun) ProfessionPeakLevel(head *pb.RpcHead, uProfType uint32, uCurLevel uint32) cfgEnum.ErrorCode {
	if uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotSupported, uProfType, uCurLevel)
	}

	stProf := this.arrProf[uProfType]

	if stProf.PeakLevel == 0 || stProf.PeakLevel != uCurLevel {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionPeakLevel, uProfType, uCurLevel, stProf.PeakLevel)
	}

	//是否最大等级
	if cfgData.GetCfgProfPeakLevel(uCurLevel+1) == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxLevel, uProfType, uCurLevel)
	}

	cfgLevel := cfgData.GetCfgProfPeakLevel(uCurLevel)
	if cfgLevel == nil {
		return plog.Print(this.AccountId, cfgData.GetProfPeakLevelErrorCode(uCurLevel), uProfType, uCurLevel)
	}

	//扣道具
	uCode := this.getPlayerBagFun().DelArrItem(head, cfgLevel.ListNeedItem, pb.EmDoingType_EDT_ProfessionPeakLevel)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode, uProfType, uCurLevel, cfgLevel.ListNeedItem)
	}

	//存数据
	uCurLevel++
	stProf.PeakLevel = uCurLevel
	this.arrProf[uProfType] = stProf

	this.UpdateSave(true)

	//计算属性
	this.calcFightPower(uProfType)

	cluster.SendToClient(head, &pb.ProfessionPeakLevelResponse{
		PacketHead: &pb.IPacket{},
		ProfType:   uProfType,
		CurLevel:   uCurLevel,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 职业部位升级请求
func (this *PlayerSystemProfessionFun) ProfessionPartLevelRequest(head *pb.RpcHead, pbRequest *pb.ProfessionPartLevelRequest) {
	uCode := this.ProfessionPartLevel(head, pbRequest.ProfType, pbRequest.PartType, pbRequest.CurLevel)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ProfessionPartLevelResponse{
			PacketHead: &pb.IPacket{},
			ProfType:   pbRequest.ProfType,
			PartType:   pbRequest.PartType,
			CurLevel:   pbRequest.CurLevel,
		}, uCode)
	}
}

// 获取部位能强化的最大等级
func (this *PlayerSystemProfessionFun) getMaxProfPartLevel(pProf *PlayerProfInfo) uint32 {
	if pProf == nil {
		return 0
	}

	uMaxLevel := uint32(0)
	if pProf.PeakLevel > 0 {
		cfgLevel := cfgData.GetCfgProfPeakLevel(pProf.PeakLevel)
		if cfgLevel != nil {
			uMaxLevel = cfgLevel.MaxPartLevel
		}
	} else {
		cfgLevel := cfgData.GetCfgProfLevel(pProf.Level)
		if cfgLevel != nil {
			uMaxLevel = cfgLevel.MaxPartLevel
		}
	}

	return uMaxLevel
}

// 获取部位能强化的最大精炼
func (this *PlayerSystemProfessionFun) getMaxProfPartRefine(pProf *PlayerProfInfo) uint32 {
	if pProf == nil {
		return 0
	}

	cfgLevel := cfgData.GetCfgProfLevel(pProf.Level)
	if cfgLevel == nil {
		return 0
	}

	return cfgLevel.MaxPartRefine
}

// 职业部位升级
func (this *PlayerSystemProfessionFun) ProfessionPartLevel(head *pb.RpcHead, uProfType uint32, uPartType uint32, uCurLevel uint32) cfgEnum.ErrorCode {
	if uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotSupported, uProfType, uPartType, uCurLevel)
	}
	if uPartType > uint32(cfgEnum.EHydraEquipment_Ring) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_EquipmentPartTypeNotSupported, uProfType, uPartType, uCurLevel)
	}

	stProf := this.arrProf[uProfType]

	pbPartInfo := stProf.PartList[uPartType]

	if pbPartInfo.Level != uCurLevel {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionPartLevel, uProfType, uPartType, uCurLevel, pbPartInfo.Level)
	}

	//真正升级
	uCode := this.innerPartLevel(head, uProfType, uPartType)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode, uProfType, uPartType, uCurLevel, pbPartInfo.Level)
	}

	cluster.SendToClient(head, &pb.ProfessionPartLevelResponse{
		PacketHead: &pb.IPacket{},
		ProfType:   uProfType,
		PartType:   uPartType,
		CurLevel:   pbPartInfo.Level,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 职业部位升级
func (this *PlayerSystemProfessionFun) innerPartLevel(head *pb.RpcHead, uProfType uint32, uPartType uint32) cfgEnum.ErrorCode {
	if uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotSupported, uProfType, uPartType, uPartType)
	}
	if uPartType > uint32(cfgEnum.EHydraEquipment_Ring) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_EquipmentPartTypeNotSupported, uProfType, uPartType, uPartType)
	}

	stProf := this.arrProf[uProfType]
	pbPartInfo := stProf.PartList[uPartType]
	if pbPartInfo.Level >= this.getMaxProfPartLevel(stProf) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxLevel, uProfType, uPartType, uPartType, pbPartInfo.Level, this.getMaxProfPartLevel(stProf))
	}

	//职业等级限制 需要算上巅峰等级
	cfgLevel := cfgData.GetCfgProfPartLevel(pbPartInfo.Level)

	//是否最大等级
	if cfgData.GetCfgProfPartLevel(pbPartInfo.Level+1) == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxLevel, uProfType, uPartType, uPartType, pbPartInfo.Level)
	}

	//扣道具
	uCode := this.getPlayerBagFun().DelItem(head, cfgLevel.ArrNeedItem.Kind, cfgLevel.ArrNeedItem.Id, cfgLevel.ArrNeedItem.Count, pb.EmDoingType_EDT_ProfessionPartLevel)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode, uProfType, uPartType, uPartType, pbPartInfo.Level, *cfgLevel.ArrNeedItem)
	}

	//成功率
	pbPartInfo.Level++

	// 词条加成效果
	prob := entry.ToValue(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_FinishedEquipmentRedefine), uint32(cfgEnum.EntryWorkTag_None))...)
	if prob > 0 && base.RandRange(0, base.MIL_PERCENT) < prob {
		if pbPartInfo.Level < this.getMaxProfPartLevel(stProf) && cfgData.GetCfgProfPartLevel(pbPartInfo.Level+1) != nil {
			pbPartInfo.Level++
		}
	}

	//如果穿戴装备 需要更新属性
	if pbPartInfo.EquipSn > 0 {
		this.CalcProp(uProfType)
	}

	this.onProfPartLevel(head, uProfType, 1)
	this.UpdateSave(true)
	return cfgEnum.ErrorCode_Success
}

func (this *PlayerSystemProfessionFun) GetMaxProfPartLevelLink() *ProfLink {
	stReturn := &ProfLink{}
	for profType, _ := range this.arrProf {
		if _, ok := this.mapProfLink[uint32(profType)]; ok {
			stReturn.uLinkLevel = base.MaxUint32(this.mapProfLink[uint32(profType)].uLinkLevel, stReturn.uLinkLevel)
			stReturn.uLinkStar = base.MaxUint32(this.mapProfLink[uint32(profType)].uLinkStar, stReturn.uLinkStar)
			stReturn.uLinkRefine = base.MaxUint32(this.mapProfLink[uint32(profType)].uLinkRefine, stReturn.uLinkRefine)
		}
	}

	return stReturn
}

// 职业部位一键升级请求
func (this *PlayerSystemProfessionFun) ProfessionPartOnekeyLevelRequest(head *pb.RpcHead, pbRequest *pb.ProfessionPartOnekeyLevelResponse) {
	uCode := this.ProfessionPartOnekeyLevel(head, pbRequest.ProfType)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ProfessionPartOnekeyLevelResponse{
			PacketHead: &pb.IPacket{},
			ProfType:   pbRequest.ProfType,
		}, uCode)
	}
}

// 职业部位一键升级请求
func (this *PlayerSystemProfessionFun) ProfessionPartOnekeyLevel(head *pb.RpcHead, uProfType uint32) cfgEnum.ErrorCode {
	if uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotSupported, uProfType, uProfType)
	}
	stProf := this.arrProf[uProfType]

	cfgPartLevel := cfgData.GetCfgProfPartLevel(stProf.PartList[0].Level)
	if cfgPartLevel == nil {
		return plog.Print(head.Id, cfgData.GetProfPartLevelErrorCode(stProf.PartList[0].Level), uProfType)
	}

	cfgProfLevel := cfgData.GetCfgProfLevel(stProf.Level)
	if cfgProfLevel == nil {
		return plog.Print(head.Id, cfgData.GetProfLevelErrorCode(stProf.Level), uProfType)
	}

	//道具不足
	uCurItemCount := this.getPlayerBagFun().GetItemCount(cfgPartLevel.ArrNeedItem.Kind, cfgPartLevel.ArrNeedItem.Id)
	if uCurItemCount <= 0 {
		return plog.Print(head.Id, cfgEnum.ErrorCode_ItemNotEnough, uProfType)
	}

	//先找未到突破等级上限 未穿戴不升级
	bBuqi := false
	sortProf := make([]*pb.PBPlayerSystemProfPartInfo, 0)
	for i := 0; i < len(stProf.PartList); i++ {
		if stProf.PartList[i].EquipSn == 0 {
			continue
		}

		if stProf.PartList[i].Level >= cfgProfLevel.MaxPartLevel {
			bBuqi = true
			continue
		}

		sortProf = append(sortProf, stProf.PartList[i])
	}

	if len(sortProf) <= 0 {
		return plog.Print(head.Id, cfgEnum.ErrorCode_NotFoundProfessionLevel, uProfType)
	}

	//按照等级从低到高排序
	if len(sortProf) > 1 {
		sort.SliceStable(sortProf, func(i, j int) bool {
			if sortProf[i].Level != sortProf[j].Level {
				return sortProf[i].Level < sortProf[j].Level
			} else {
				return sortProf[i].Part < sortProf[j].Part
			}

			return false
		})
	}

	//过滤需要升级的部位，优先升级最低一档的等级
	MAX_CRICLE := uint32(200)
	uCircle := uint32(0)
	if len(sortProf) == 1 {
		uCircle = MAX_CRICLE
	} else {
		uMaxLevel := sortProf[0].Level
		for i := 0; i < len(sortProf); i++ {
			if sortProf[i].Level > uMaxLevel {
				uMaxLevel = sortProf[i].Level
				sortProf = sortProf[:i]
				break
			}
		}

		//等级一样
		if uMaxLevel == sortProf[0].Level {
			uCircle = MAX_CRICLE
		} else {
			uCircle = uMaxLevel - sortProf[0].Level
			bBuqi = true
		}
	}

	//一键升级将会对当前职业全部装备进行级别+1的操作
	//若当前材料较为大量，依次检测全部装备升200级、100级、50级、10级、1级的情况，材料满足升级需要时，执行对应升级选项，直至材料不够，停止升级
	uPreLen := uint32(len(sortProf))
	uPreCircle := uint32(0)
	uPreNeedItemCount := uint32(0)
	for i := uint32(0); i < uCircle; i++ {
		cfgPartLevel = cfgData.GetCfgProfPartLevel(sortProf[0].Level + i)
		if cfgPartLevel == nil {
			break
		}

		if sortProf[0].Level+i >= cfgProfLevel.MaxPartLevel {
			break
		}

		if cfgData.GetCfgProfPartLevel(sortProf[0].Level+i+1) == nil {
			break
		}

		if uPreNeedItemCount+uint32(cfgPartLevel.ArrNeedItem.Count)*uPreLen > uint32(uCurItemCount) {
			break
		}

		uPreNeedItemCount += uint32(cfgPartLevel.ArrNeedItem.Count) * uPreLen

		uPreCircle++
	}

	//补齐不做限制
	if !bBuqi {
		if uPreCircle < 5 {
			uCircle = 1
		} else if uPreCircle < 10 {
			uCircle = 5
		} else if uPreCircle < 50 {
			uCircle = 10
		} else if uPreCircle < 100 {
			uCircle = 50
		} else {
			uCircle = 100
		}
	}

	uCostItemCount := int64(0)
	uCircle = base.MinUint32(MAX_CRICLE, uCircle)
	mapUpdate := make(map[uint32]uint32)
	uAddLevel := uint32(0)
	for uCircle > 0 {
		uCircle--

		bBreak := false
		for i := 0; i < len(sortProf); i++ {
			cfgPartLevel = cfgData.GetCfgProfPartLevel(sortProf[i].Level)
			if cfgPartLevel == nil {
				bBreak = true
				break
			}

			if sortProf[i].Level >= cfgProfLevel.MaxPartLevel {
				break
			}

			//是否满级
			if cfgData.GetCfgProfPartLevel(sortProf[i].Level+1) == nil {
				bBreak = true
				break
			}

			if uCostItemCount+cfgPartLevel.ArrNeedItem.Count > uCurItemCount {
				bBreak = true
				break
			}

			uCostItemCount += cfgPartLevel.ArrNeedItem.Count
			sortProf[i].Level++
			mapUpdate[sortProf[i].Part] = sortProf[i].Level
			uAddLevel++
		}

		if bBreak {
			break
		}
	}

	if uCostItemCount <= 0 {
		return plog.Print(head.Id, cfgEnum.ErrorCode_ItemNotEnough, uProfType)
	}

	//扣道具
	this.getPlayerBagFun().DelItem(head, cfgPartLevel.ArrNeedItem.Kind, cfgPartLevel.ArrNeedItem.Id, uCostItemCount, pb.EmDoingType_EDT_ProfessionPartLevel)

	//每个部位升一级，道具不够才停止
	pbResponse := &pb.ProfessionPartOnekeyLevelResponse{
		PacketHead: &pb.IPacket{},
		ProfType:   uProfType,
	}

	for partType, _ := range mapUpdate {
		pbResponse.PartList = append(pbResponse.PartList, this.arrProf[uProfType].PartList[partType])
	}

	//更新属性
	this.CalcProp(uProfType)

	this.onProfPartLevel(head, uProfType, uAddLevel)
	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 职业升级
func (this *PlayerSystemProfessionFun) onProfLevel(head *pb.RpcHead, uProfType uint32, uAdd uint32) {
	//成就统计
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfessionUpgrade, uAdd, uProfType)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfessionUpgrade, uAdd, uint32(cfgEnum.EHydraProf_Any))

	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfessionLevel, this.GetProfLevel(uProfType, false), uProfType)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfessionLevel, this.GetProfLevel(uint32(cfgEnum.EHydraProf_Any), false), uint32(cfgEnum.EHydraProf_Any))

	//上阵成就
	this.getPlayerHeroFun().OnHeroGameHeroList()

	this.calcFightPower(uProfType)
}

// 职业部位升级
func (this *PlayerSystemProfessionFun) onProfPartLevel(head *pb.RpcHead, uProfType uint32, uAdd uint32) {
	//成就统计
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfEquipmentUpgrade, uAdd, uProfType)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfEquipmentUpgrade, uAdd, uint32(cfgEnum.EHydraProf_Any))

	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfEquipmentLevel, this.GetProfPartLevel(uProfType, false), uProfType)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfEquipmentLevel, this.GetProfPartLevel(uint32(cfgEnum.EHydraProf_Any), false), uint32(cfgEnum.EHydraProf_Any))

	uMaxLinkLevel := this.GetMaxProfPartLevelLink().uLinkLevel
	if uMaxLinkLevel > this.maxProfLink.uLinkLevel {
		this.maxProfLink.uLinkLevel = uMaxLinkLevel
		this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfPartLevelLinkLevel, uMaxLinkLevel, uint32(cfgEnum.EHydraProf_Any))
	}
}

// 职业部位精炼
func (this *PlayerSystemProfessionFun) onProfPartRefine(head *pb.RpcHead, uProfType uint32, uAdd uint32) {
	//成就触发
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfessionRefine, uAdd, uProfType)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfessionRefine, uAdd, uint32(cfgEnum.EHydraProf_Any))

	uMaxLinkRefine := this.GetMaxProfPartLevelLink().uLinkRefine
	if uMaxLinkRefine > this.maxProfLink.uLinkRefine {
		this.maxProfLink.uLinkRefine = uMaxLinkRefine
		this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfPartRefineLinkLevel, uMaxLinkRefine, uint32(cfgEnum.EHydraProf_Any))
	}

}

// 职业穿戴装备请求
func (this *PlayerSystemProfessionFun) ProfessionEquipRequest(head *pb.RpcHead, pbRequest *pb.ProfessionEquipRequest) {
	uCode := this.ProfessionEquip(head, pbRequest.ProfType, pbRequest.PartType, pbRequest.Sn)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ProfessionEquipResponse{
			PacketHead: &pb.IPacket{},
			ProfType:   pbRequest.ProfType,
			PartType:   pbRequest.PartType,
			Sn:         pbRequest.Sn,
		}, uCode)
	}
}

// 职业穿戴装备
func (this *PlayerSystemProfessionFun) ProfessionEquip(head *pb.RpcHead, uProfType uint32, uPartType uint32, uSn uint32) cfgEnum.ErrorCode {
	if uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotSupported, uProfType, uPartType, uPartType)
	}
	if uPartType > uint32(cfgEnum.EHydraEquipment_Ring) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_EquipmentPartTypeNotSupported, uProfType, uPartType, uPartType)
	}

	stProf := this.arrProf[uProfType]
	pbPartInfo := stProf.PartList[uPartType]

	if pbPartInfo.EquipSn == uSn {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionEquipIsUsing, uProfType, uPartType, uSn)
	}

	//穿戴
	uCode := this.GetPlayerEquipmentFun().Equip(head, uProfType, uPartType, uSn, pbPartInfo.EquipSn)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode, uProfType, uPartType, uSn, pbPartInfo.EquipSn)
	}

	pbPartInfo.EquipSn = uSn
	this.arrProf[uProfType].PartList[uPartType] = pbPartInfo

	this.UpdateSave(true)
	this.CalcProp(uProfType)

	//成就统计
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfEquipmentCount, this.GetProfEquipCount(uProfType), uProfType)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfEquipmentCount, this.GetProfEquipCount(uint32(cfgEnum.EHydraProf_Any)), uint32(cfgEnum.EHydraProf_Any))
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfEquipmentMaxCount, this.GetMaxProfEquipCount())

	uMaxLinkStar := this.GetMaxProfPartLevelLink().uLinkStar
	if uMaxLinkStar > this.maxProfLink.uLinkStar {
		this.maxProfLink.uLinkStar = uMaxLinkStar
		this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfPartStarLinkLevel, uMaxLinkStar, uint32(cfgEnum.EHydraProf_Any))
	}

	cluster.SendToClient(head, &pb.ProfessionEquipResponse{
		PacketHead: &pb.IPacket{},
		ProfType:   uProfType,
		PartType:   uPartType,
		Sn:         uSn,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 职业一键脱装备请求
func (this *PlayerSystemProfessionFun) ProfessionOnekeyUnEquipRequest(head *pb.RpcHead, pbRequest *pb.ProfessionOnekeyUnEquipRequest) {
	uCode := this.ProfessionOnekeyUnEquip(head, pbRequest.ProfType)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ProfessionOnekeyUnEquipResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 职业一键脱装备
func (this *PlayerSystemProfessionFun) ProfessionOnekeyUnEquip(head *pb.RpcHead, uProfType uint32) cfgEnum.ErrorCode {
	if uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotSupported, uProfType)
	}

	stProf := this.arrProf[uProfType]

	pbResponse := &pb.ProfessionOnekeyUnEquipResponse{
		PacketHead: &pb.IPacket{},
		ProfType:   uProfType,
	}
	for upartId, pPartInfo := range stProf.PartList {
		if pPartInfo.EquipSn > 0 {
			this.GetPlayerEquipmentFun().Equip(head, uProfType, uint32(upartId), 0, pPartInfo.EquipSn)
			pbResponse.EquipSnList = append(pbResponse.EquipSnList, pPartInfo.EquipSn)
			pPartInfo.EquipSn = 0
			stProf.PartList[upartId] = pPartInfo
		}
	}

	this.arrProf[uProfType] = stProf

	//成就统计
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfEquipmentCount, this.GetProfEquipCount(uProfType), uProfType)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfEquipmentCount, this.GetProfEquipCount(uint32(cfgEnum.EHydraProf_Any)), uint32(cfgEnum.EHydraProf_Any))
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_ProfEquipmentMaxCount, this.GetMaxProfEquipCount())

	this.UpdateSave(true)
	this.CalcProp(uProfType)

	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 计算属性
func (this *PlayerSystemProfessionFun) CalcProp(uProfType uint32) {
	if uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return
	}

	stProf := this.arrProf[uProfType]
	stProf.mapProp = make(map[uint32]float64)
	cfgProfConfig := cfgData.GetCfgProfConfig(uProfType)
	if cfgProfConfig != nil {
		for k, v := range cfgProfConfig.MapInitProp {
			stProf.mapProp[k] += float64(v)
		}
	}

	//算套装
	bSuitActive := true
	uMinLinkLevel := uint32(0)
	uMinLinkRefine := uint32(0)
	for _, v := range stProf.PartList {
		if v.EquipSn <= 0 {
			bSuitActive = false
			continue
		}
		//查找装备
		pbEquipment := this.GetPlayerEquipmentFun().getEquipment(v.EquipSn)
		if pbEquipment == nil {
			continue
		}
		//精炼共鸣
		if uMinLinkRefine == 0 || v.Refine < uMinLinkRefine {
			uMinLinkRefine = v.Refine
		}
		//改造共鸣
		if uMinLinkLevel == 0 || v.Level < uMinLinkLevel {
			uMinLinkLevel = v.Level
		}
		cfgQuality := cfgData.GetCfgEquipmentQuality(pbEquipment.Quality)

		//主词条 升级级影响主词条
		stProf.mapProp[pbEquipment.MainProp.PropId] += float64(pbEquipment.MainProp.Value)
		stProf.mapProp[pbEquipment.MainProp.PropId] += float64(cfgQuality.LevelRate*pbEquipment.MainProp.Value*(v.Level-1)) / base.MIL_PERCENT

		//次词条 精炼等级和突破等级 突破等级0开始
		for i := 0; i < len(pbEquipment.MinorPropList); i++ {
			stProf.mapProp[pbEquipment.MinorPropList[i].PropId] += float64(pbEquipment.MinorPropList[i].Value)
			stProf.mapProp[pbEquipment.MinorPropList[i].PropId] += float64(cfgQuality.RefineRate*pbEquipment.MinorPropList[i].Value*(v.Refine-1)) / base.MIL_PERCENT
			stProf.mapProp[pbEquipment.MinorPropList[i].PropId] += float64(cfgQuality.RefineTupoRate*pbEquipment.MinorPropList[i].Value*v.RefineTupo) / base.MIL_PERCENT
		}

		//副词条
		for i := 0; i < len(pbEquipment.VicePropList); i++ {
			stProf.mapProp[pbEquipment.VicePropList[i].PropId] += float64(pbEquipment.VicePropList[i].Value)
		}

		//共鸣词条
		for i := 0; i < len(pbEquipment.LinkPropList); i++ {
			stProf.mapProp[pbEquipment.LinkPropList[i].PropId] += float64(pbEquipment.LinkPropList[i].Value)
		}
	}

	//算共鸣 必须穿戴装备
	if bSuitActive {
		//改造共鸣
		cfgLevelLink := cfgData.GetCfgProfPartLevelLink(uMinLinkLevel)
		if cfgLevelLink != nil {
			for k, v := range cfgLevelLink.MapAddProp {
				stProf.mapProp[k] += float64(v)
			}

			if _, ok := this.mapProfLink[uProfType]; !ok {
				this.mapProfLink[uProfType] = &ProfLink{}
			}

			this.mapProfLink[uProfType].uLinkLevel = cfgLevelLink.Id
		}

		//精炼共鸣
		cfgRefineLink := cfgData.GetCfgProfPartRefineLink(uMinLinkRefine)
		if cfgRefineLink != nil {
			for k, v := range cfgRefineLink.MapAddProp {
				stProf.mapProp[k] += float64(v)
			}

			if _, ok := this.mapProfLink[uProfType]; !ok {
				this.mapProfLink[uProfType] = &ProfLink{}
			}
			this.mapProfLink[uProfType].uLinkRefine = cfgRefineLink.Id
		}
	}

	this.arrProf[uProfType] = stProf

	//通知英雄系统算战斗力
	this.calcFightPower(uProfType)

}

func (this *PlayerSystemProfessionFun) calcFightPower(uProfType uint32) {

	//通知英雄系统
	this.getPlayerHeroFun().UpdateProfFightPower(uProfType)

}
func (this *PlayerSystemProfessionFun) DipUpdateProf(head *pb.RpcHead, pbRequest *pb.PBPlayerSystemProfInfo) {
	this.arrProf[pbRequest.ProfType].PBPlayerSystemProfInfo = pbRequest
	this.UpdateSave(true)
	this.CalcProp(pbRequest.ProfType)
}

// 职业部位精炼请求
func (this *PlayerSystemProfessionFun) ProfessionPartRefineRequest(head *pb.RpcHead, pbRequest *pb.ProfessionPartRefineRequest) {
	uCode := this.ProfessionPartRefine(head, pbRequest.ProfType, pbRequest.PartType, pbRequest.CurLevel)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ProfessionPartRefineResponse{
			PacketHead: &pb.IPacket{},
			ProfType:   pbRequest.ProfType,
			PartType:   pbRequest.PartType,
			CurLevel:   pbRequest.CurLevel,
		}, uCode)
	}
}

// 职业部位精炼请求
func (this *PlayerSystemProfessionFun) ProfessionPartRefine(head *pb.RpcHead, uProfType uint32, uPartType uint32, uCurLevel uint32) cfgEnum.ErrorCode {
	if uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotSupported, uProfType, uPartType, uPartType)
	}
	if uPartType > uint32(cfgEnum.EHydraEquipment_Ring) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_EquipmentPartTypeNotSupported, uProfType, uPartType, uPartType)
	}

	stProf := this.arrProf[uProfType]

	pbPartInfo := stProf.PartList[uPartType]

	if pbPartInfo.Refine != uCurLevel {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionPartRefine, uProfType, uPartType, uCurLevel, pbPartInfo.RefineTupo)
	}

	//职业等级限制
	if pbPartInfo.Refine >= this.getMaxProfPartRefine(stProf) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxLevel, uProfType, uPartType, pbPartInfo.Refine, this.getMaxProfPartRefine(stProf))
	}

	cfgRefine := cfgData.GetCfgProfPartRefine(pbPartInfo.Refine)
	if cfgRefine == nil {
		return plog.Print(this.AccountId, cfgData.GetProfPartRefineErrorCode(pbPartInfo.Refine), uProfType, uPartType, uCurLevel, pbPartInfo.RefineTupo)
	}

	cfgTupo := cfgData.GetCfgProfPartRefineTupo(pbPartInfo.RefineTupo)
	if cfgTupo == nil {
		return plog.Print(this.AccountId, cfgData.GetProfPartRefineTupoErrorCode(pbPartInfo.RefineTupo), uProfType, uPartType, uCurLevel, pbPartInfo.RefineTupo)
	}

	if pbPartInfo.Refine >= cfgTupo.MaxLevel {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxLevel, uProfType, uPartType, uCurLevel, pbPartInfo.Level, cfgTupo.MaxLevel)
	}

	//是否最大等级
	if cfgData.GetCfgProfPartRefine(pbPartInfo.Refine+1) == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxLevel, uProfType, uPartType, uCurLevel, pbPartInfo.Level, cfgTupo.MaxLevel)
	}

	//扣道具
	uCode := this.getPlayerBagFun().DelItem(head, cfgRefine.ArrNeedItem.Kind, cfgRefine.ArrNeedItem.Id, cfgRefine.ArrNeedItem.Count, pb.EmDoingType_EDT_ProfessionPartRefine)
	if uCode != cfgEnum.ErrorCode_Success {
		return uCode
	}

	//成功率
	pbPartInfo.Refine++

	// 词条加成效果
	prob := entry.ToValue(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_FinishedEquipmentProfessionRefine), uint32(cfgEnum.EntryWorkTag_None))...)
	if prob > 0 && base.RandRange(0, base.MIL_PERCENT) < prob {
		if pbPartInfo.Refine < cfgTupo.MaxLevel && cfgData.GetCfgProfPartRefine(pbPartInfo.Refine+1) != nil {
			pbPartInfo.Refine++
		}
	}

	//如果穿戴装备 需要更新属性
	if pbPartInfo.EquipSn > 0 {
		this.CalcProp(uProfType)
	}

	this.onProfPartRefine(head, uProfType, 1)

	this.UpdateSave(true)

	cluster.SendToClient(head, &pb.ProfessionPartRefineResponse{
		PacketHead: &pb.IPacket{},
		ProfType:   uProfType,
		PartType:   uPartType,
		CurLevel:   pbPartInfo.Refine,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 职业部位一键精炼请求
func (this *PlayerSystemProfessionFun) ProfessionPartOnekeyRefineRequest(head *pb.RpcHead, pbRequest *pb.ProfessionPartOnekeyRefineRequest) {
	uCode := this.ProfessionPartOnekeyRefine(head, pbRequest.ProfType)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ProfessionPartOnekeyRefineResponse{
			PacketHead: &pb.IPacket{},
			ProfType:   pbRequest.ProfType,
		}, uCode)
	}
}

// 职业部位一键精炼请求
func (this *PlayerSystemProfessionFun) ProfessionPartOnekeyRefine(head *pb.RpcHead, uProfType uint32) cfgEnum.ErrorCode {
	if uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotSupported, uProfType)
	}
	stProf := this.arrProf[uProfType]

	cfgPartRefine := cfgData.GetCfgProfPartRefine(stProf.PartList[0].Refine)
	if cfgPartRefine == nil {
		return plog.Print(head.Id, cfgData.GetProfPartRefineErrorCode(stProf.PartList[0].Refine), uProfType)
	}

	//职业等级限制
	cfgProfLevel := cfgData.GetCfgProfLevel(stProf.Level)
	if cfgProfLevel == nil {
		return plog.Print(head.Id, cfgData.GetProfLevelErrorCode(stProf.Level), stProf.Level, "cfgProfLevel")
	}

	//道具不足
	uCurItemCount := this.getPlayerBagFun().GetItemCount(cfgPartRefine.ArrNeedItem.Kind, cfgPartRefine.ArrNeedItem.Id)

	//按照精炼等级排序 未穿戴不升级 到达突破等级的不参与
	sortProf := make([]*pb.PBPlayerSystemProfPartInfo, 0)
	for i := 0; i < len(stProf.PartList); i++ {
		if stProf.PartList[i].EquipSn == 0 {
			continue
		}

		//判断是否需要突破
		cfgPartRefineTupo := cfgData.GetCfgProfPartRefineTupo(stProf.PartList[i].RefineTupo)
		if cfgPartRefineTupo == nil || stProf.PartList[i].Refine >= cfgPartRefineTupo.MaxLevel {
			continue
		}

		//职业限制
		if stProf.PartList[i].Refine >= cfgProfLevel.MaxPartRefine {
			continue
		}

		sortProf = append(sortProf, stProf.PartList[i])
	}

	if len(sortProf) <= 0 {
		return plog.Print(head.Id, cfgEnum.ErrorCode_NotFoundProfessionPartRefine, uProfType)
	}

	//按照等级从低到高排序
	if len(sortProf) > 1 {
		sort.SliceStable(sortProf, func(i, j int) bool {
			if sortProf[i].Refine != sortProf[j].Refine {
				return sortProf[i].Refine < sortProf[j].Refine
			} else {
				return sortProf[i].Part < sortProf[j].Part
			}

			return false
		})
	}

	//过滤需要升级的部位，优先升级最低一档的等级
	MAX_CRICLE := uint32(200)
	uCircle := uint32(0)
	if len(sortProf) == 1 {
		uCircle = MAX_CRICLE
	} else {
		uMaxLevel := sortProf[0].Refine
		for i := 0; i < len(sortProf); i++ {
			if sortProf[i].Refine > uMaxLevel {
				uMaxLevel = sortProf[i].Refine
				sortProf = sortProf[:i]
				break
			}
		}

		//等级一样
		if uMaxLevel == sortProf[0].Refine {
			uCircle = MAX_CRICLE
		} else {
			uCircle = uMaxLevel - sortProf[0].Refine
		}
	}

	uCode := cfgEnum.ErrorCode_Fail
	uCostItemCount := int64(0)
	uCircle = base.MinUint32(MAX_CRICLE, uCircle)
	mapUpdate := make(map[uint32]uint32)
	uAddRefine := uint32(0)
	for uCircle > 0 {
		uCircle--

		bBreak := false
		for i := 0; i < len(sortProf); i++ {
			//判断是否需要突破
			cfgPartRefineTupo := cfgData.GetCfgProfPartRefineTupo(sortProf[i].RefineTupo)
			if cfgPartRefineTupo == nil || sortProf[i].Refine >= cfgPartRefineTupo.MaxLevel {
				uCode = cfgEnum.ErrorCode_MaxLevel
				bBreak = true
				break
			}

			//职业限制
			if sortProf[i].Refine >= cfgProfLevel.MaxPartRefine {
				uCode = cfgEnum.ErrorCode_MaxLevel
				bBreak = true
				break
			}

			//判断是否需要突破
			cfgPartRefine = cfgData.GetCfgProfPartRefine(sortProf[i].Refine)
			if cfgPartRefine == nil {
				uCode = cfgEnum.ErrorCode_MaxLevel
				bBreak = true
				break
			}

			//是否满级
			if cfgData.GetCfgProfPartRefine(sortProf[i].Refine+1) == nil {
				uCode = cfgEnum.ErrorCode_MaxLevel
				bBreak = true
				break
			}

			if uCostItemCount+cfgPartRefine.ArrNeedItem.Count > uCurItemCount {
				uCode = cfgEnum.ErrorCode_ItemNotEnough
				bBreak = true
				break
			}

			uCostItemCount += cfgPartRefine.ArrNeedItem.Count
			sortProf[i].Refine++
			mapUpdate[sortProf[i].Part] = 1
			uAddRefine++
		}

		if bBreak {
			break
		}
	}

	//无可以升级的
	if len(mapUpdate) <= 0 || uCostItemCount <= 0 {
		return plog.Print(head.Id, uCode, uProfType)
	}

	//扣道具
	this.getPlayerBagFun().DelItem(head, cfgPartRefine.ArrNeedItem.Kind, cfgPartRefine.ArrNeedItem.Id, uCostItemCount, pb.EmDoingType_EDT_ProfessionPartRefine)

	//每个部位升一级，道具不够才停止
	pbResponse := &pb.ProfessionPartOnekeyRefineResponse{
		PacketHead: &pb.IPacket{},
		ProfType:   uProfType,
	}

	for partType, _ := range mapUpdate {
		pbResponse.PartList = append(pbResponse.PartList, this.arrProf[uProfType].PartList[partType])
	}

	//更新属性
	this.CalcProp(uProfType)

	this.onProfPartRefine(head, uProfType, uAddRefine)
	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 职业部位精炼突破请求
func (this *PlayerSystemProfessionFun) ProfessionPartRefineTupoRequest(head *pb.RpcHead, pbRequest *pb.ProfessionPartRefineTupoRequest) {
	uCode := this.ProfessionPartRefineTupo(head, pbRequest.ProfType, pbRequest.PartType, pbRequest.CurLevel)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ProfessionPartRefineTupoResponse{
			PacketHead: &pb.IPacket{},
			ProfType:   pbRequest.ProfType,
			PartType:   pbRequest.PartType,
			CurLevel:   pbRequest.CurLevel,
		}, uCode)
	}
}

// 职业部位精炼突破请求
func (this *PlayerSystemProfessionFun) ProfessionPartRefineTupo(head *pb.RpcHead, uProfType uint32, uPartType uint32, uCurLevel uint32) cfgEnum.ErrorCode {
	if uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotSupported, uProfType, uPartType, uPartType)
	}
	if uPartType > uint32(cfgEnum.EHydraEquipment_Ring) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_EquipmentPartTypeNotSupported, uProfType, uPartType, uPartType)
	}

	stProf := this.arrProf[uProfType]
	pbPartInfo := stProf.PartList[uPartType]
	if pbPartInfo.RefineTupo != uCurLevel {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionPartRefineTupo, uProfType, uPartType, uCurLevel, pbPartInfo.RefineTupo)
	}

	cfgRefineTupo := cfgData.GetCfgProfPartRefineTupo(pbPartInfo.RefineTupo)
	if cfgRefineTupo == nil {
		return plog.Print(this.AccountId, cfgData.GetProfPartRefineTupoErrorCode(pbPartInfo.RefineTupo), uProfType, uPartType, uCurLevel, pbPartInfo.RefineTupo)
	}

	if cfgData.GetCfgProfPartRefineTupo(pbPartInfo.RefineTupo+1) == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_MaxLevel, uProfType, uPartType, uCurLevel, pbPartInfo.RefineTupo)
	}

	//判断是否到达突破等级
	if pbPartInfo.Refine != cfgRefineTupo.MaxLevel {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedRefineLevel, uProfType, uPartType, uCurLevel, pbPartInfo.RefineTupo)
	}

	//扣道具
	uCode := this.getPlayerBagFun().DelItem(head, cfgRefineTupo.ArrNeedItem.Kind, cfgRefineTupo.ArrNeedItem.Id, cfgRefineTupo.ArrNeedItem.Count, pb.EmDoingType_EDT_ProfessionPartRefineTupo)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode, uProfType, uPartType, uCurLevel, pbPartInfo.RefineTupo)
	}

	//成功率
	pbPartInfo.RefineTupo++

	//如果穿戴装备 需要更新属性
	if pbPartInfo.EquipSn > 0 {
		this.CalcProp(uProfType)
	}

	this.UpdateSave(true)

	cluster.SendToClient(head, &pb.ProfessionPartRefineTupoResponse{
		PacketHead: &pb.IPacket{},
		ProfType:   uProfType,
		PartType:   uPartType,
		CurLevel:   pbPartInfo.RefineTupo,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 职业部位重置请求
func (this *PlayerSystemProfessionFun) ProfessionPartResetRequest(head *pb.RpcHead, pbRequest *pb.ProfessionPartResetRequest) {
	uCode := this.ProfessionPartReset(head, pbRequest.ProfType, pbRequest.PartType)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ProfessionPartResetResponse{
			PacketHead: &pb.IPacket{},
			ProfType:   pbRequest.ProfType,
		}, uCode)
	}
}

// 职业部位精炼突破请求
func (this *PlayerSystemProfessionFun) ProfessionPartReset(head *pb.RpcHead, uProfType uint32, uPartType uint32) cfgEnum.ErrorCode {
	if uProfType > uint32(cfgEnum.EHydraProf_Machinist) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfessionNotSupported, uProfType, uPartType)
	}
	profInfo := this.arrProf[uProfType]
	mapAddItem := make(map[uint32]int64)
	pbResponse := &pb.ProfessionPartResetResponse{
		PacketHead: &pb.IPacket{},
		ProfType:   uProfType,
	}
	for key, partInfo := range profInfo.PartList {
		if uPartType != uint32(cfgEnum.EQuality_Any) && uPartType != uint32(key) {
			continue
		}

		if partInfo.Level <= 1 && partInfo.Refine <= 1 {
			continue
		}

		//给升级奖励
		if partInfo.Level > 1 {
			cfgLevel := cfgData.GetCfgProfPartLevel(partInfo.Level)
			if cfgLevel == nil || cfgLevel.ArrTotalItem.Id == 0 {
				plog.Info("(this *PlayerSystemProfessionFun) ProfessionPartReset level not find")
				continue
			}

			mapAddItem[cfgLevel.ArrTotalItem.Id] += cfgLevel.ArrTotalItem.Count
			partInfo.Level = 1
		}

		//给精炼奖励
		if partInfo.Refine > 1 {
			cfgRefine := cfgData.GetCfgProfPartRefine(partInfo.Refine)
			if cfgRefine == nil || cfgRefine.ArrTotalItem.Id == 0 {
				plog.Info("(this *PlayerSystemProfessionFun) ProfessionPartReset level not find")
				continue
			}

			mapAddItem[cfgRefine.ArrTotalItem.Id] += cfgRefine.ArrTotalItem.Count
			partInfo.Refine = 1
		}

		//给精炼突破奖励
		if partInfo.RefineTupo > 0 {
			cfgRefineTupo := cfgData.GetCfgProfPartRefineTupo(partInfo.RefineTupo)
			if cfgRefineTupo == nil || cfgRefineTupo.ArrTotalItem.Id == 0 {
				plog.Info("(this *PlayerSystemProfessionFun) ProfessionPartReset level not find")
				continue
			}

			mapAddItem[cfgRefineTupo.ArrTotalItem.Id] += cfgRefineTupo.ArrTotalItem.Count
			partInfo.RefineTupo = 0
		}

		pbResponse.PartList = append(pbResponse.PartList, partInfo)

	}

	//判断是否是所有重置
	if len(pbResponse.PartList) <= 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ProfNoPartReset, uProfType, uPartType)
	}

	//给奖励
	arrItem := make([]*common2.ItemInfo, 0)
	for id, value := range mapAddItem {
		arrItem = append(arrItem, &common2.ItemInfo{Kind: uint32(cfgEnum.ESystemType_Item), Id: id, Count: value})
	}
	this.getPlayerBagFun().AddArrItem(head, arrItem, pb.EmDoingType_EDT_Reset, true)

	//更新 属性
	this.calcFightPower(uProfType)

	this.UpdateSave(true)

	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}
