package playerFun

import (
	"corps/base"
	cfgEnum "corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/common/orm/redis"
	orm "corps/common/orm/redis"
	report2 "corps/common/report"
	"corps/framework/cluster"
	"corps/pb"
	"corps/server/game/module/entry"
	"encoding/json"
	"fmt"

	"corps/framework/plog"

	"github.com/golang/protobuf/proto"
)

type (
	PlayerSystemCommonFun struct {
		PlayerFun
		pbData             pb.PBPlayerSystemCommon         //系统
		mapAvatars         map[uint32]*pb.PBAvatar         // 头像
		mapAvatarFrames    map[uint32]*pb.PBAvatarFrame    // 头像框
		mapAcode           map[string]*pb.PBPlayerGiftCode //激活码
		mapAdvest          map[uint32]*pb.PBAdvertInfo     //广告数据
		listSystemOpenType []uint32                        //系统开关
	}
)

// 登录成功回调
func (this *PlayerSystemCommonFun) LoadPlayerDBFinish() {
	// 默认解锁头像框
	if _, ok := this.mapAvatarFrames[1]; !ok {
		this.mapAvatarFrames[1] = &pb.PBAvatarFrame{FrameID: 1}
	}
}

func (this *PlayerSystemCommonFun) getAvatars() (rets []*pb.PBAvatar) {
	for _, item := range this.mapAvatars {
		rets = append(rets, item)
	}
	return
}

func (this *PlayerSystemCommonFun) getAvatarFrames() (rets []*pb.PBAvatarFrame) {
	for _, item := range this.mapAvatarFrames {
		rets = append(rets, item)
	}
	return
}

func (this *PlayerSystemCommonFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.mapAcode = make(map[string]*pb.PBPlayerGiftCode)
	this.mapAvatars = make(map[uint32]*pb.PBAvatar)
	this.mapAvatarFrames = make(map[uint32]*pb.PBAvatarFrame)
	this.mapAdvest = make(map[uint32]*pb.PBAdvertInfo)
}

func (this *PlayerSystemCommonFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = new(pb.PBPlayerSystem)
	}

	this.saveData()
	pbData.System.Common = &this.pbData
}
func (this *PlayerSystemCommonFun) loadData(pbData *pb.PBPlayerSystemCommon) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemCommon{}
	}

	this.pbData = *pbData
	this.mapAcode = make(map[string]*pb.PBPlayerGiftCode)
	for i := 0; i < len(this.pbData.GiftCode); i++ {
		this.mapAcode[this.pbData.GiftCode[i].Acode] = this.pbData.GiftCode[i]
	}
	this.mapAvatars = make(map[uint32]*pb.PBAvatar)
	for _, item := range this.pbData.Avatars {
		this.mapAvatars[item.AvatarID] = item
	}
	this.mapAvatarFrames = make(map[uint32]*pb.PBAvatarFrame)
	for _, item := range this.pbData.AvatarFrames {
		this.mapAvatarFrames[item.FrameID] = item
	}
	// 默认解锁头像框
	if _, ok := this.mapAvatarFrames[1]; !ok {
		this.mapAvatarFrames[1] = &pb.PBAvatarFrame{FrameID: 1}
	}

	this.mapAdvest = make(map[uint32]*pb.PBAdvertInfo)
	for _, item := range this.pbData.AdvertList {
		this.mapAdvest[item.Type] = item
	}

	this.listSystemOpenType = make([]uint32, 0)
	for _, id := range this.pbData.SystemOpenIds {
		cfg := cfgData.GetCfgSystemOpenById(id)
		if cfg != nil {
			this.listSystemOpenType = append(this.listSystemOpenType, cfg.SystemTypes...)
		}
	}

	if this.pbData.EjectAdvertInfo == nil {
		this.pbData.EjectAdvertInfo = &pb.PBEjectAdvertInfo{
			Id:              0,
			NextRefreshTime: 0,
		}
	}
	this.UpdateSave(true)
}

func (this *PlayerSystemCommonFun) saveData() {
	this.pbData.GiftCode = make([]*pb.PBPlayerGiftCode, 0)
	for _, v := range this.mapAcode {
		this.pbData.GiftCode = append(this.pbData.GiftCode, v)
	}
	this.pbData.Avatars = this.pbData.Avatars[:0]
	for _, item := range this.mapAvatars {
		this.pbData.Avatars = append(this.pbData.Avatars, item)
	}
	this.pbData.AvatarFrames = this.pbData.AvatarFrames[:0]
	for _, item := range this.mapAvatarFrames {
		this.pbData.AvatarFrames = append(this.pbData.AvatarFrames, item)
	}
	this.pbData.AdvertList = this.pbData.AdvertList[:0]
	for _, item := range this.mapAdvest {
		this.pbData.AdvertList = append(this.pbData.AdvertList, item)
	}
}
func (this *PlayerSystemCommonFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	if pbSystem.Common == nil {
		pbSystem.Common = &pb.PBPlayerSystemCommon{}
	}
	this.loadData(pbSystem.Common)
	this.UpdateSave(false)
}

func (this *PlayerSystemCommonFun) LoadComplete() {
	this.CheckSystemOpen(&pb.RpcHead{Id: this.AccountId}, cfgEnum.ESystemOpenType_Default)

	if this.CheckSystemTypeOpen(cfgEnum.ESystemUnlockType_EjectAdvert) {
		this.OnActiveEjectAdvert()
	}

	if this.CheckSystemTypeOpen(cfgEnum.ESystemUnlockType_FirstCharge) {
		this.getPlayerSystemChargeFun().OnSystemOpenFirstCharge(&pb.RpcHead{Id: this.AccountId}, false)
	}
	//发送公告 是否有新的直接弹窗
	//this.NoticeNotify()
}

func (this *PlayerSystemCommonFun) OnActiveEjectAdvert() {
	if this.pbData.EjectAdvertInfo == nil || (this.pbData.EjectAdvertInfo.Id == 0 && this.pbData.EjectAdvertInfo.NextRefreshTime == 0) {
		this.pbData.EjectAdvertInfo = &pb.PBEjectAdvertInfo{
			Id:              1,
			NextRefreshTime: 0,
		}

		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.EjectAdvertNotify{
			PacketHead:      &pb.IPacket{},
			EjectAdvertInfo: this.pbData.EjectAdvertInfo,
		}, cfgEnum.ErrorCode_Success)
	}
}

func (this *PlayerSystemCommonFun) CheckSystemTypeOpen(openType cfgEnum.ESystemUnlockType) bool {
	return base.ArrayContainsValue(this.listSystemOpenType, uint32(openType))
}

func (this *PlayerSystemCommonFun) AllSystemOpen(head *pb.RpcHead) {
	mapCfg := cfgData.GetAllCfgSystemOpen()
	arrAdd := make([]uint32, 0)
	for _, mapInfo := range mapCfg {
		for _, cfg := range mapInfo {
			if base.ArrayContainsValue(this.pbData.SystemOpenIds, cfg.Id) {
				continue
			}

			arrAdd = append(arrAdd, cfg.Id)
			this.pbData.SystemOpenIds = append(this.pbData.SystemOpenIds, cfg.Id)
			this.listSystemOpenType = append(this.listSystemOpenType, cfg.SystemTypes...)
		}
	}

	if len(arrAdd) > 0 {
		this.UpdateSave(true)

		cluster.SendToClient(head, &pb.SystemOpenNotify{
			PacketHead:    &pb.IPacket{},
			SystemOpenIds: arrAdd,
		}, cfgEnum.ErrorCode_Success)
	}
}

func (this *PlayerSystemCommonFun) CheckSystemOpen(head *pb.RpcHead, eType cfgEnum.ESystemOpenType) {
	mapCfg := cfgData.GetCfgSystemOpen(eType)
	if len(mapCfg) <= 0 {
		return
	}

	mapId, stageId := this.getPlayerSystemBattleNormalFun().GetFinishMapIdAndStageId()
	hmapId, hstageId := this.getPlayerSystemBattleHookFun().GetFinishMapIdAndStageId()
	regDays := this.getPlayerBaseFun().GetRegDays()
	arrAdd := make([]uint32, 0)
	arrAddTypes := make([]uint32, 0)
	for id, v := range mapCfg {
		if base.ArrayContainsValue(this.pbData.SystemOpenIds, id) {
			continue
		}

		bSuc := true
		for _, Param := range v.Condition {
			if len(Param) < 1 {
				bSuc = false
				break
			}

			switch cfgEnum.ESystemOpenType(Param[0]) {
			case cfgEnum.ESystemOpenType_BattleNormalEnd:
				if len(Param) != 3 {
					bSuc = false
					break
				}

				if mapId*1000+stageId < Param[1]*1000+Param[2] {
					bSuc = false
					break
				}
			case cfgEnum.ESystemOpenType_BattleHookEnd:
				if len(Param) != 3 {
					bSuc = false
					break
				}

				if hmapId*1000+hstageId < Param[1]*1000+Param[2] {
					bSuc = false
					break
				}
			case cfgEnum.ESystemOpenType_FBattleNormalEnter:
				if len(Param) != 3 {
					bSuc = false
					break
				}
				fmapId, fstageId := this.getPlayerSystemBattleNormalFun().GetMapIdAndStageId()
				if fmapId*1000+fstageId > Param[1]*1000+Param[2] {
					continue
				} else if fmapId*1000+fstageId == Param[1]*1000+Param[2] {
					fightCount := this.getPlayerSystemBattleNormalFun().GetBattleMapInfo().FightCount
					if fightCount == 0 {
						bSuc = false
						break
					}
				} else {
					bSuc = false
					break
				}

			case cfgEnum.ESystemOpenType_RegisterDay:
				if len(Param) != 2 {
					bSuc = false
					break
				}

				if regDays < Param[1] {
					bSuc = false
					break
				}
			}
		}
		if !bSuc {
			continue
		}

		arrAdd = append(arrAdd, v.Id)
		this.pbData.SystemOpenIds = append(this.pbData.SystemOpenIds, v.Id)
		this.listSystemOpenType = append(arrAddTypes, v.SystemTypes...)
		arrAddTypes = append(arrAddTypes, v.SystemTypes...)
		//plog.Info("CheckSystemOpen add ", v.Id)
	}

	if len(arrAdd) > 0 {
		this.UpdateSave(true)

		cluster.SendToClient(head, &pb.SystemOpenNotify{
			PacketHead:    &pb.IPacket{},
			SystemOpenIds: arrAdd,
		}, cfgEnum.ErrorCode_Success)

		this.getPlayerSystemChargeBPFun().OnSystemOpenTypes(head, arrAddTypes)
		this.getPlayerActivityOpenServerGiftFun().OnSystemOpenTypes(head, arrAddTypes)
		this.getPlayerSystemShopFun().OnSystemOpenTypes(head, arrAddTypes)
		for _, sysType := range arrAddTypes {
			if sysType == uint32(cfgEnum.ESystemUnlockType_EjectAdvert) {
				this.OnActiveEjectAdvert()
			} else if sysType == uint32(cfgEnum.ESystemUnlockType_FirstCharge) {
				this.getPlayerSystemChargeFun().OnSystemOpenFirstCharge(&pb.RpcHead{Id: this.AccountId}, true)
			}

			this.getPlayerSystemDrawFun().OnSystemTypeOpen(head, sysType)

		}
	}
}
func (this *PlayerSystemCommonFun) IsSystemOpen(uSystemId uint32) bool {
	return base.ArrayContainsValue(this.pbData.SystemOpenIds, uSystemId)
}
func (this *PlayerSystemCommonFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	this.saveData()

	pbSystem.Common = &this.pbData
	return this.BSave
}

// 使用道具请求
func (this *PlayerSystemCommonFun) GiftCodeRequest(head *pb.RpcHead, pbRequest *pb.GiftCodeRequest) {
	uCode := this.GiftCode(head, pbRequest.Acode)
	cluster.SendToClient(head, &pb.GiftCodeResponse{
		PacketHead: &pb.IPacket{},
		Acode:      pbRequest.Acode}, uCode)
}

// 隔天刷新通知
func (this *PlayerSystemCommonFun) PassDay(isDay, isWeek, isMonth bool) {
	//更新刷新时间
	bNotice := false
	pbNotice := &pb.AdvertNotify{
		PacketHead: &pb.IPacket{Id: this.AccountId},
	}
	for _, info := range this.mapAdvest {
		if info.DailyCount > 0 {
			info.DailyCount = 0
			bNotice = true
			pbNotice.AdvertList = append(pbNotice.AdvertList, info)
		}
	}

	if bNotice {
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbNotice, cfgEnum.ErrorCode_Success)
		this.UpdateSave(true)
	}

	this.CheckSystemOpen(&pb.RpcHead{Id: this.AccountId}, cfgEnum.ESystemOpenType_RegisterDay)

	//刷新广告弹出 必须激活
	if this.pbData.EjectAdvertInfo != nil && this.pbData.EjectAdvertInfo.Id > 0 {
		this.pbData.EjectAdvertInfo = &pb.PBEjectAdvertInfo{
			Id:              1,
			NextRefreshTime: 0,
		}

		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.EjectAdvertNotify{
			PacketHead:      &pb.IPacket{},
			EjectAdvertInfo: this.pbData.EjectAdvertInfo,
		}, cfgEnum.ErrorCode_Success)
	}
}

// 兑换
func (this *PlayerSystemCommonFun) GiftCode(head *pb.RpcHead, strAcode string) cfgEnum.ErrorCode {
	emErrorCode := cfgEnum.ErrorCode_Success

	if len(strAcode) < 5 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_GiftCodeInvalid, strAcode)
	}

	redisCommon := redis.GetCommonRedis()

	//需要加锁
	redisLock := orm.AddRedisLock(redisCommon, fmt.Sprintf("%s_%s_lock", base.ERK_CommonGiftCodeList, strAcode))
	if redisLock == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ServerBusy, strAcode)
	}

	//释放锁
	defer redisLock.UnLock()

	//通过激活码查询礼包数据
	uNowTime := base.GetNow()
	strUseKey := strAcode
	strListKey := strAcode
	jsonGiftCodeData := redisCommon.HGet(base.ERK_CommonGiftCodeCommon, strAcode)
	if jsonGiftCodeData == "" {

		//取前六位
		if len(strAcode) <= base.COMMON_GIFTCODE_KEYLENGTH {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_GiftCodeInvalid, strAcode)
		}

		strUseKey = strAcode[0:base.COMMON_GIFTCODE_KEYLENGTH]
		strListKey = strAcode[0:(base.COMMON_GIFTCODE_KEYLENGTH - 1)]
		jsonGiftCodeData = redisCommon.HGet(fmt.Sprintf("%s%s", base.ERK_CommonGiftCode, strUseKey), strAcode)
	}

	//不存在
	if jsonGiftCodeData == "" {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_GiftCodeError, strAcode)
	}

	giftCodeInfo := &common.GiftCodeInfo{}

	//先判断是否是通用兑换码
	if json.Unmarshal([]byte(jsonGiftCodeData), giftCodeInfo) != nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_JsonUnmarshal, strAcode)
	}

	//兑换码
	jsonGiftCodeListData := redisCommon.HGet(base.ERK_CommonGiftCodeList, strListKey)
	if jsonGiftCodeListData == "" {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_GiftCodeError, strAcode)
	}

	giftCodeListData := &common.GiftCodeListInfo{}
	if json.Unmarshal([]byte(jsonGiftCodeListData), giftCodeListData) != nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_JsonUnmarshal, jsonGiftCodeListData)
	}

	//查询兑换码
	jsonGiftListData := redisCommon.HGet(base.ERK_CommonGift, giftCodeInfo.ActKey)
	if jsonGiftListData == "" {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_GiftCodeError, strAcode)
	}

	giftListData := &common.GiftInfo{}
	if json.Unmarshal([]byte(jsonGiftListData), giftListData) != nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_JsonUnmarshal, strAcode)
	}

	//判断时效
	if uNowTime > giftListData.ETime || uNowTime < giftListData.STime {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_GiftCodeExpired, strAcode)
	}

	playerBaseFun := this.getPlayerBaseFun()
	if giftListData.RegTime > 0 && giftListData.RegTime < playerBaseFun.GetRegTime() {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedRegTime, strAcode)
	}

	if playerBaseFun.GetVipLevel() < giftListData.VipLevel {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedVipLevel, strAcode)
	}

	//判断是否是限制次数
	pbAcode, ok := this.mapAcode[strUseKey]
	uCurCount := uint32(0)
	if ok {
		uCurCount = pbAcode.Count
		if uCurCount > giftListData.UseCount && pbAcode.Time < giftListData.ETime {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_GiftCodeUsed, strAcode)
		}
	}

	//判断总兑换次数
	if giftListData.TotalCount > 0 && giftCodeListData.TotalUseCount >= giftListData.TotalCount {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_GiftCodepNoCount, strAcode)
	}

	if giftCodeInfo.Uid > 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_GiftCodeUsed, strAcode)
	}

	if uCurCount >= giftListData.UseCount {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_GiftCodepNoCount, strAcode)
	}

	//更新累计使用数量
	giftCodeListData.TotalUseCount++
	byteGiftCodeListData, err := json.Marshal(giftCodeListData)
	if err != nil {
		plog.Info("GiftCode json.Marshal error:%s", err.Error())
		byteGiftCodeListData = []byte(jsonGiftCodeListData)
	}
	redisCommon.HSet(base.ERK_CommonGiftCodeList, strUseKey, byteGiftCodeListData)

	//更新个人兑换码
	if giftListData.Type != uint32(pb.EmGiftCodeType_GAT_Common) {
		giftCodeInfo.Uid = this.AccountId
		giftCodeInfo.UTime = uNowTime
		byteGiftCodeData, err := json.Marshal(giftCodeInfo)
		if err != nil {
			plog.Info("GiftCode json.Marshal error:%s", err.Error())
			byteGiftCodeData = []byte(jsonGiftCodeData)
		}

		redisCommon.HSet(fmt.Sprintf("%s%s", base.ERK_CommonGiftCode, strUseKey), strAcode, byteGiftCodeData)
	}

	//存储个人兑换
	this.mapAcode[strUseKey] = &pb.PBPlayerGiftCode{
		Acode: strAcode,
		Time:  uNowTime,
		Count: uCurCount + 1,
	}

	//发道具
	this.getPlayerBagFun().AddArrItem(head, giftListData.Items, pb.EmDoingType_EDT_GiftCode, true)

	this.UpdateSave(true)
	return emErrorCode
}
func (this *PlayerSystemCommonFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemCommon{}
}

// 设置玩家数据
func (this *PlayerSystemCommonFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}

	pbSystem := pbData.(*pb.PBPlayerSystemCommon)
	if pbSystem == nil {
		return false
	}

	this.loadData(pbSystem)

	return true
}

func (this *PlayerSystemCommonFun) ChangeAvatarRequest(head *pb.RpcHead, headID uint32) {
	if errCode := this.changeHead(head, headID); errCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ChangeAvatarResponse{
			PacketHead: &pb.IPacket{},
		}, errCode)
	}
}

func (this *PlayerSystemCommonFun) changeHead(head *pb.RpcHead, headID uint32) cfgEnum.ErrorCode {
	// 判断这个头像是否合法
	cfg := cfgData.GetCfgAvatar(headID)
	if cfg == nil {
		return plog.Print(head.Id, cfgData.GetAvatarErrorCode(headID), "HeadConfig", headID)
	}
	// 判单是否已经解锁
	if _, ok := this.mapAvatars[headID]; !ok {
		return plog.Print(head.Id, cfgEnum.ErrorCode_ItemNotExist, "headId", headID)
	}
	// 判断当前头像是否同一个
	baseFun := this.getPlayerBaseFun()
	if baseFun.GetDisplay().AvatarID == headID {
		return plog.Print(head.Id, cfgEnum.ErrorCode_InUse, "headId", headID)
	}
	// 切换头像
	baseFun.SetAvatar(headID)
	// 回包
	cluster.SendToClient(head, &pb.ChangeAvatarResponse{
		PacketHead: &pb.IPacket{},
		AvatarID:   headID,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

func (this *PlayerSystemCommonFun) ChangeAvatarFrameRequest(head *pb.RpcHead, headID uint32) {
	if errCode := this.changeHeadIcon(head, headID); errCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ChangeAvatarFrameResponse{
			PacketHead: &pb.IPacket{},
		}, errCode)
	}
}

func (this *PlayerSystemCommonFun) changeHeadIcon(head *pb.RpcHead, headIconID uint32) cfgEnum.ErrorCode {
	// 判断这个头像框是否合法
	cfg := cfgData.GetCfgAvatarFrame(headIconID)
	if cfg == nil {
		return plog.Print(head.Id, cfgData.GetAvatarFrameErrorCode(headIconID), "HeadIconConfig", headIconID)
	}
	// 判断头像框是否已经解锁
	if _, ok := this.mapAvatarFrames[headIconID]; !ok {
		return plog.Print(head.Id, cfgEnum.ErrorCode_ItemNotExist, "headIconId", headIconID)
	}
	// 判单当前头像框是否同一个
	baseFun := this.getPlayerBaseFun()
	if baseFun.GetDisplay().AvatarFrameID == headIconID {
		return plog.Print(head.Id, cfgEnum.ErrorCode_InUse, "headIconId", headIconID)
	}
	// 切换头像框
	baseFun.SetAvatarFrame(headIconID)
	// 回报
	cluster.SendToClient(head, &pb.ChangeAvatarFrameResponse{
		PacketHead: &pb.IPacket{},
		FrameID:    headIconID,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 新增头像
func (this *PlayerSystemCommonFun) AddHead(head *pb.RpcHead, id uint32, doingType pb.EmDoingType) cfgEnum.ErrorCode {
	// 判断这个头像框是否合法
	cfg := cfgData.GetCfgAvatar(id)
	if cfg == nil {
		return plog.Print(head.Id, cfgData.GetAvatarErrorCode(id), "HeadConfig", id)
	}
	// 判断是否已经存在这个头像
	if _, ok := this.mapAvatars[id]; ok {
		plog.Debug("Head(%d) already have", id)
		return cfgEnum.ErrorCode_Success
	}
	// 增加头像
	this.mapAvatars[id] = &pb.PBAvatar{AvatarID: id, Type: uint32(doingType)}

	// 通知客户端
	cluster.SendToClient(head, &pb.AvatarNotify{PacketHead: &pb.IPacket{}, Avatars: this.getAvatars()}, cfgEnum.ErrorCode_Success)
	// 数据上报
	report2.Send(head, &report2.ReportAddItem{
		Kind:   uint32(cfgEnum.ESystemType_Head),
		ItemID: id,
		Add:    1,
		Total:  1,
		Doing:  uint32(doingType),
	})
	return cfgEnum.ErrorCode_Success
}

// 新增头像框
func (this *PlayerSystemCommonFun) AddHeadIcon(head *pb.RpcHead, id uint32, doing pb.EmDoingType) cfgEnum.ErrorCode {
	// 判断这个头像框是否合法
	cfg := cfgData.GetCfgAvatarFrame(id)
	if cfg == nil {
		return plog.Print(head.Id, cfgData.GetAvatarFrameErrorCode(id), "HeadIconConfig", id)
	}
	// 判断是否已经存在这个头像
	if _, ok := this.mapAvatarFrames[id]; ok {
		plog.Debug("Head(%d) already have", id)
		return cfgEnum.ErrorCode_Success
	}
	// 增加头像
	this.mapAvatarFrames[id] = &pb.PBAvatarFrame{FrameID: id, Type: uint32(doing)}

	// 通知客户端
	cluster.SendToClient(head, &pb.AvatarFrameNotify{Frames: this.getAvatarFrames()}, cfgEnum.ErrorCode_Success)
	// 数据上报
	report2.Send(head, &report2.ReportAddItem{
		Kind:   uint32(cfgEnum.ESystemType_HeadIcon),
		ItemID: id,
		Add:    1,
		Total:  1,
		Doing:  uint32(doing),
	})
	return cfgEnum.ErrorCode_Success
}

// 公告
func (this *PlayerSystemCommonFun) NoticeRequest(head *pb.RpcHead) {
	pbResponse := &pb.NoticeResponse{
		PacketHead: &pb.IPacket{},
	}

	//查询redis公告内容
	redisCommon := redis.GetCommonRedis()
	if redisCommon == nil {
		cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_DATABASE)
		return
	}

	uCurTime := base.GetNow()

	mapKey := redisCommon.HGetAll(base.ERK_CommonNotice)
	arrKey := base.SortStringMapByKey(mapKey, true)
	for _, key := range arrKey {
		stNotice := &common.NoticeInfo{}
		err := json.Unmarshal([]byte(mapKey[key]), stNotice)
		if err != nil {
			continue
		}
		if uCurTime < stNotice.StartTime || uCurTime > stNotice.EndTime {
			continue
		}

		pbNotice := &pb.PBNotice{
			Id:        stNotice.Id,
			Title:     stNotice.Title,
			BeginTime: stNotice.StartTime,
			EndTime:   stNotice.EndTime,
			Content:   stNotice.Content,
			Sender:    stNotice.Sender,
		}

		pbResponse.NoticeList = append(pbResponse.NoticeList, pbNotice)
	}

	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)
}

// 公告弹窗
func (this *PlayerSystemCommonFun) NoticeNotify() {
	//查询redis公告内容
	redisCommon := redis.GetCommonRedis()
	if redisCommon == nil {
		return
	}

	//遍历所有的获取最大
	mapKey := redisCommon.HGetAll(base.ERK_CommonNotice)
	if len(mapKey) <= 0 {
		return
	}

	arrKey := base.SortStringMapByKey(mapKey, true)
	if len(arrKey) <= 0 {
		return
	}

	srtData := mapKey[arrKey[0]]
	stNotice := &common.NoticeInfo{}
	err := json.Unmarshal([]byte(srtData), stNotice)
	if err != nil {
		return
	}

	uMaxId := base.StringToUInt32(arrKey[0])

	uCurTime := base.GetNow()
	if uCurTime < stNotice.StartTime || uCurTime > stNotice.EndTime {
		return
	}

	pbResponse := &pb.NoticeNotify{
		PacketHead: &pb.IPacket{},
		Notice: &pb.PBNotice{
			Id:        stNotice.Id,
			Title:     stNotice.Title,
			BeginTime: stNotice.StartTime,
			EndTime:   stNotice.EndTime,
			Content:   stNotice.Content,
		},
		IsNew: false,
	}

	//判断是否最新的
	if this.pbData.MaxNoticeId != uMaxId {
		this.pbData.MaxNoticeId = uMaxId
		pbResponse.IsNew = true
		this.UpdateSave(true)
	}

	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbResponse, cfgEnum.ErrorCode_Success)
}

// 观看广告完成请求
func (this *PlayerSystemCommonFun) AdvertRequest(head *pb.RpcHead, pbRequest *pb.AdvertRequest) {
	uCode := this.AddAdvert(head, pbRequest.AdvestType)
	cluster.SendToClient(head, &pb.AdvertResponse{
		PacketHead: &pb.IPacket{},
		AdvestInfo: this.mapAdvest[pbRequest.AdvestType],
	}, uCode)
}

func (this *PlayerSystemCommonFun) OnAdvert(head *pb.RpcHead, uAdvertType uint32) {
	if uAdvertType <= uint32(cfgEnum.EAdvertType_None) || uAdvertType >= uint32(cfgEnum.EAdvertType_ShareBegin) {
		return
	}

	//成就触发
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_WatchAdvertise, 1, uAdvertType)
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_WatchAdvertise, 1, uint32(cfgEnum.EAdvertType_None))
}

// 广告系统
func (this *PlayerSystemCommonFun) AddAdvert(head *pb.RpcHead, uAdvertType uint32) cfgEnum.ErrorCode {
	cfgAdvert := cfgData.GetCfgAdvertConfig(uAdvertType)
	if cfgAdvert == nil {
		return plog.Print(this.AccountId, cfgData.GetAdvertConfigErrorCode(uAdvertType), uAdvertType)
	}

	pbData, ok := this.mapAdvest[uAdvertType]
	if !ok {
		pbData = &pb.PBAdvertInfo{Type: uAdvertType, DailyCount: 0}
		this.mapAdvest[uAdvertType] = pbData
	}

	if pbData.DailyCount >= cfgAdvert.DailyCount {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_DailyMaxCount, uAdvertType)
	}

	pbData.DailyCount++
	this.UpdateSave(true)

	pbNotify := &pb.AdvertNotify{
		PacketHead: &pb.IPacket{},
	}
	pbNotify.AdvertList = append(pbNotify.AdvertList, pbData)
	cluster.SendToClient(head, pbNotify, cfgEnum.ErrorCode_Success)

	this.OnAdvert(head, uAdvertType)
	return cfgEnum.ErrorCode_Success
}

// 查看页面请求
func (this *PlayerSystemCommonFun) PageOpenRequest(head *pb.RpcHead, pbRequest *pb.PageOpenRequest) {
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_PageOpen, 1, pbRequest.PageType)
	cluster.SendToClient(head, &pb.PageOpenResponse{
		PacketHead: &pb.IPacket{}}, cfgEnum.ErrorCode_Success)
}

// 系统开启奖励领取
func (this *PlayerSystemCommonFun) SystemOpenPrizeRequest(head *pb.RpcHead, pbRequest *pb.SystemOpenPrizeRequest) {
	uCode := this.SystemOpenPrize(head, pbRequest.Id)
	cluster.SendToClient(head, &pb.SystemOpenPrizeResponse{
		PacketHead: &pb.IPacket{},
		Id:         pbRequest.Id},
		uCode)
}

// 广告系统
func (this *PlayerSystemCommonFun) SystemOpenPrize(head *pb.RpcHead, uId uint32) cfgEnum.ErrorCode {
	cfgSytemOpen := cfgData.GetCfgSystemOpenById(uId)
	if cfgSytemOpen == nil {
		return plog.Print(this.AccountId, cfgData.GetSystemOpenErrorCode(uId), uId)
	}

	if len(cfgSytemOpen.AddPrize) <= 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoPrize, uId)
	}

	//是否领取过
	if !base.ArrayContainsValue(this.pbData.SystemOpenIds, uId) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NotYetUnlocked, uId)
	}

	//是否领取过
	if base.ArrayContainsValue(this.pbData.SystemOpenPrizeList, uId) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HavePrize, uId)
	}

	this.pbData.SystemOpenPrizeList = append(this.pbData.SystemOpenPrizeList, uId)

	this.getPlayerBagFun().AddArrItem(head, cfgSytemOpen.AddPrize, pb.EmDoingType_EDT_System, true)

	this.UpdateSave(true)
	return cfgEnum.ErrorCode_Success
}

// 获取特权
func (this *PlayerSystemCommonFun) GetPrivilege(emPrivilegeType cfgEnum.PrivilegeType) uint32 {
	uValue := this.getPlayerSystemChargeCardFun().GetPrivilege(uint32(emPrivilegeType))
	switch emPrivilegeType {
	case cfgEnum.PrivilegeType_OfflineTime:
		uValue = uValue * 3600
		mapID, _ := this.getPlayerSystemBattleHookFun().GetMapIdAndStageId()
		cfg := cfgData.GetCfgBattleHookMap(mapID)
		if cfg != nil {
			uValue = uValue + cfg.MaxOfflineTime*60
		}

		// 词条加成
		vv := entry.ToValue(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_OfflineIncomeTime), uint32(cfgEnum.EntryWorkTag_None))...)
		uValue += vv
	case cfgEnum.PrivilegeType_BlackShopRefreshCount:
		vals := this.getEntry().Get(uint32(cfgEnum.EntryEffectType_AddShopRefresh), uint32(cfgEnum.EShopType_BlackShop))
		if len(vals) > 0 && len(vals[0].List) > 0 {
			uValue += vals[0].List[0]
		}
	}

	return uValue
}

func (this *PlayerSystemCommonFun) UpdatePrivilege(emPrivilegeType cfgEnum.PrivilegeType) {
	switch emPrivilegeType {
	case cfgEnum.PrivilegeType_OfflineTime:
		this.getPlayerSystemOffline().UpdatePrivilege(emPrivilegeType, this.GetPrivilege(emPrivilegeType))
	case cfgEnum.PrivilegeType_BlackShopRefreshCount:
		this.getPlayerSystemShopFun().UpdatePrivilege(emPrivilegeType, this.GetPrivilege(emPrivilegeType))
	}
}

// 广告弹出请求
func (this *PlayerSystemCommonFun) EjectAdvertRequest(head *pb.RpcHead, pbRequest *pb.EjectAdvertRequest) {
	uCode := this.EjectAdvert(head, pbRequest.Id)
	cluster.SendToClient(head, &pb.EjectAdvertResponse{
		PacketHead:      &pb.IPacket{},
		EjectAdvertInfo: this.pbData.EjectAdvertInfo},
		uCode)
}

// 广告弹出请求
func (this *PlayerSystemCommonFun) EjectAdvert(head *pb.RpcHead, uId uint32) cfgEnum.ErrorCode {
	if this.pbData.EjectAdvertInfo == nil || this.pbData.EjectAdvertInfo.Id == 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoPrize, uId)
	}
	uCurTime := base.GetNow()
	if this.pbData.EjectAdvertInfo.NextRefreshTime > 0 && this.pbData.EjectAdvertInfo.NextRefreshTime > uCurTime {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_CoolDown, uId)
	}

	cfgEject := cfgData.GetCfgEjectAdvertConfig(uId)
	if cfgEject == nil {
		return plog.Print(this.AccountId, cfgData.GetEjectAdvertConfigErrorCode(uId), uId)
	}

	//跨天id变了，不再随机
	if uId == this.pbData.EjectAdvertInfo.Id {
		this.pbData.EjectAdvertInfo.Id = cfgData.GetRandCfgEjectAdvertId()
	}

	// 广告加成
	bag := this.getPlayerBagFun()
	rewards := []*pb.PBAddItemData{}
	propValue := entry.ToValue(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_AdertiseReward), uint32(cfgEnum.EAdvertType_EjectAward))...)
	if base.IsRadio(propValue) {
		items := bag.GetPbItems([]*common.ItemInfo{cfgEject.AddPrize}, pb.EmDoingType_EDT_Entry)
		plog.Trace("EntryEffectType_AdertiseReward rewards: %v", items)
		rewards = append(rewards, items...)
	}

	// 奖励
	addItems := bag.GetPbItems([]*common.ItemInfo{cfgEject.AddPrize}, pb.EmDoingType_EDT_AdvertEject)
	rewards = append(rewards, addItems...)

	// 发送奖励
	this.pbData.EjectAdvertInfo.NextRefreshTime = base.GetNow() + uint64(cfgEject.CoolTime)
	bag.AddPbItems(head, rewards, pb.EmDoingType_EDT_AdvertEject, true)

	// 回调
	this.OnAdvert(head, uint32(cfgEnum.EAdvertType_EjectAward))
	this.UpdateSave(true)
	return cfgEnum.ErrorCode_Success
}
