package playerFun

import (
	"context"
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/common/dao/player_base_display"
	"corps/common/orm/mysql"
	"corps/common/orm/redis"
	"corps/common/serverCommon"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
)

type (
	PlayerBaseFun struct {
		PlayerFun
		pbBase *pb.PBPlayerBase //基本信息
	}
)

func (this *PlayerBaseFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.pbBase = &pb.PBPlayerBase{
		Display: &pb.PBPlayerBaseDisplay{},
	}
}

// 一键领取邮件请求
func (this *PlayerBaseFun) Test(ctx context.Context, pbRequest proto.Message) {
	pb := pbRequest.(*pb.PBPlayerCrystal)
	fmt.Println(pb)
}

func (this *PlayerBaseFun) Load(pData []byte) {
	proto.Unmarshal(pData, this.pbBase)
	this.loadData()
}

func (this *PlayerBaseFun) loadData() {
	if this.pbBase.Display == nil {
		this.pbBase.Display = &pb.PBPlayerBaseDisplay{}
	}

	this.UpdateSave(true)
}

func (this *PlayerBaseFun) Save(bNow bool) {
	if !this.BSave {
		return
	}

	this.BSave = false

	this.saveData()

	//通知db保存玩家数据
	buff, _ := proto.Marshal(this.pbBase)
	cluster.SendToDb(&pb.RpcHead{Id: this.AccountId}, "DbPlayerMgr", "SavePlayerDB", this.PbType, buff, bNow)
	plog.Info("(this *PlayerBaseFun) accountid:%d  baseid:%d", this.AccountId, this.pbBase.Display.AccountId)
	if this.AccountId != this.pbBase.Display.AccountId {
		plog.Fatal("(this *PlayerBaseFun) accountid:%d  baseid:%d", this.AccountId, this.pbBase.Display.AccountId)
	}
}

func (this *PlayerBaseFun) saveData() {
}

// 保存
func (this *PlayerBaseFun) SavePb(pbData *pb.PBPlayerBase) {
	if pbData == nil {
		pbData = &pb.PBPlayerBase{}
	}
	this.saveData()
	base.DeepCopy(this.pbBase, pbData)
}

func (this *PlayerBaseFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	//设置成默认头像
	if this.pbBase.Display.AvatarID <= 0 {
		this.SetAvatar(5)
	}
	if this.pbBase.Display.AvatarFrameID <= 0 {
		this.SetAvatarFrame(1)
	}
	pbData.Base = &pb.PBPlayerBase{}
	this.SavePb(pbData.Base)
}
func (this *PlayerBaseFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerBase{}
}

// 设置玩家数据
func (this *PlayerBaseFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}

	pbSystem := pbData.(*pb.PBPlayerBase)
	if pbSystem == nil {
		return false
	}

	this.pbBase = pbSystem
	this.loadData()

	this.updateDisplay()
	return true
}

func (this *PlayerBaseFun) SetAvatar(headID uint32) {
	this.pbBase.Display.AvatarID = headID
	this.updateDisplay()
	this.UpdateSave(true)

}

func (this *PlayerBaseFun) SetAvatarFrame(id uint32) {
	this.pbBase.Display.AvatarFrameID = id
	this.updateDisplay()
	this.UpdateSave(true)
}

func (this *PlayerBaseFun) GetDisplay() *pb.PBPlayerBaseDisplay {
	return this.pbBase.Display
}

func (this *PlayerBaseFun) GetPlayerBase() *pb.PBPlayerBase {
	return this.pbBase
}

// 获取VIP等级
func (this *PlayerBaseFun) GetVipLevel() uint32 {
	return this.pbBase.Display.VipLevel
}

func (this *PlayerBaseFun) GetServerId() uint32 {
	return this.pbBase.Display.SeverId
}
func (this *PlayerBaseFun) GetOpenSeverTime() uint64 {
	return this.pbBase.SeverStartTime
}
func (this *PlayerBaseFun) GetOpenServerDays() uint32 {
	return base.DiffDays(this.pbBase.SeverStartTime, base.GetNow()) + 1
}

func (this *PlayerBaseFun) GetServerStartTime() uint64 {
	return this.pbBase.SeverStartTime
}

// 获取玩家等级
func (this *PlayerBaseFun) GetPlayerLevel() uint32 {
	return this.pbBase.Display.PlayerLevel
}
func (this *PlayerBaseFun) GetNewPlayerTypeList() []uint32 {
	return this.pbBase.NewPlayerTypeList
}
func (this *PlayerBaseFun) AddNewPlayerTypeList(uType uint32) {
	this.pbBase.NewPlayerTypeList = append(this.pbBase.NewPlayerTypeList, uType)
	this.UpdateSave(true)
}

// 获取注册天数
func (this *PlayerBaseFun) GetRegDays() uint32 {
	return base.DiffDays(this.pbBase.CreateTime, base.GetNow()) + 1
}

// 获取注册天数
func (this *PlayerBaseFun) GetRegTime() uint64 {
	return this.pbBase.CreateTime
}

// 是否新玩家
func (this *PlayerBaseFun) LoginState() pb.LoginState {
	return this.pbBase.LoginState
}

func toWeek(t time.Time) int {
	ret := int(t.Weekday())
	if ret == 0 {
		ret = 7
	}
	return ret
}

func toTime(now uint64) time.Time {
	tnow := time.Unix(int64(now), 0)
	return time.Date(tnow.Year(), tnow.Month(), tnow.Day(), 0, 0, 0, 0, tnow.Location())
}

func (this *PlayerBaseFun) CheckDay() (rets []bool) {
	now := base.GetNow()
	if this.pbBase.LastDailyTime == 0 {
		rets = append(rets, true, true, true)
	} else {
		current := toTime(now)
		last := toTime(this.pbBase.LastDailyTime)
		isMonth := last.Year() != current.Year() || last.Month() != current.Month()
		isWeek := last.Year() == current.Year() && (current.Sub(last) >= 7*24*time.Hour || toWeek(current) < toWeek(last))
		isDay := isMonth || isWeek || last.Day() != current.Day()
		rets = append(rets, isDay, isWeek, isMonth)
	}
	if rets[0] || rets[1] || rets[2] {
		// 发送通知
		this.pbBase.LastDailyTime = now
		cluster.SendToClient(&pb.RpcHead{
			Id: this.AccountId,
		}, &pb.PassTimeNotify{
			PacketHead: &pb.IPacket{},
			IsDay:      rets[0],
			IsWeek:     rets[1],
			IsMonth:    rets[2],
			CurTime:    now,
		}, cfgEnum.ErrorCode_Success)

		this.UpdateSave(true)
	}
	return
}
func (this *PlayerBaseFun) LoadComplete() {
	if this.pbBase.LoginState == pb.LoginState_None {
		this.pbBase.LoginState = pb.LoginState_Init
		this.UpdateSave(true)
	} else if this.pbBase.LoginState == pb.LoginState_Init {
		this.pbBase.LoginState = pb.LoginState_SetName
		this.UpdateSave(true)
	}
}

// 新玩家 需要初始化数据
func (this *PlayerBaseFun) NewPlayer() {
	//初始化等级
	this.pbBase.Display.PlayerLevel = 1
	this.pbBase.LoginState = pb.LoginState_Init
	this.pbBase.Display.AccountId = this.AccountId
	this.pbBase.CreateTime = base.GetNow()
	this.pbBase.Display.PlayerName = cfgData.GetCfgRandomName(this.AccountId) //fmt.Sprintf("player_%d", this.AccountId)

	redis := redis.GetCommonRedis()
	this.pbBase.Display.SeverId, this.pbBase.SeverStartTime = redis.ZGetMaxRankAndMemberByScore(base.ERK_ServerList, this.AccountId)
	this.updatePlayerName(this.pbBase.Display.PlayerName)

	//设置账号数据
	serverCommon.GetPlayerAccountData(this.pbBase)

	this.updateDisplay()
	this.UpdateSave(true)
}

// 增加玩家等级
func (this *PlayerBaseFun) AddPlayerLevel(head *pb.RpcHead, uAdd uint32, emDoing pb.EmDoingType) cfgEnum.ErrorCode {
	if uAdd <= 0 {
		return cfgEnum.ErrorCode_Success
	}

	uRealAdd := uint32(0)
	for i := uint32(1); i <= uAdd; i++ {
		cfgLevel := serverCommon.GetCfgPlayerLevel(this.pbBase.Display.PlayerLevel + i)
		if cfgLevel == nil {
			break
		}

		uRealAdd++
	}

	this.pbBase.Display.PlayerLevel += uRealAdd

	this.UpdateSave(true)
	//更新到客户端
	this.updateKvToClient(head, "PlayerLevel", int64(this.pbBase.Display.PlayerLevel))
	return cfgEnum.ErrorCode_Success
}

// 设置玩家数据
func (this *PlayerBaseFun) updateKvToClient(head *pb.RpcHead, strKey string, value int64) {
	//通知客户端
	pbResponse := &pb.PlayerUpdateKvNotify{
		PacketHead: &pb.IPacket{
			Id: this.AccountId,
		},
	}

	pbResponse.ListInfo = append(pbResponse.ListInfo, &pb.PBStringInt64{
		Key:   strKey,
		Value: value,
	})

	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)
}

// 更新玩家外显
func (this *PlayerBaseFun) updateDisplay() {
	if err := player_base_display.Set(this.AccountId, this.GetDisplay()); err != nil {
		plog.Error("redis connect error OnSaveRedis, error: %v", err)
		return
	}
}

// 修改玩家名字
func (this *PlayerBaseFun) ChangePlayerNameRequest(head *pb.RpcHead, name string) {
	if errCode := this.changePlayerName(head, name); errCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ChangePlayerNameResponse{
			PacketHead: &pb.IPacket{},
		}, errCode)
	}
}

func (this *PlayerBaseFun) updatePlayerName(strPlayerName string) cfgEnum.ErrorCode {
	db := mysql.DBMGR.GetCommonDb()
	if db == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, "updatePlayerName", strPlayerName)
	}

	// 调用存储过程
	rows, err := db.QuerySql(fmt.Sprintf("CALL p_player_name(\"%s\", %d)", strPlayerName, this.AccountId))
	if err != nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoData, "updatePlayerName", strPlayerName)
	}

	if rows.Row().Int("outResult") > 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NameRepeat, "updatePlayerName", strPlayerName)
	}

	this.pbBase.Display.PlayerName = strPlayerName

	this.updateDisplay()

	return cfgEnum.ErrorCode_Success
}

func (this *PlayerBaseFun) changePlayerName(head *pb.RpcHead, strPlayerName string) cfgEnum.ErrorCode {
	//判断长度
	if len(strPlayerName) <= 1 || len(strPlayerName) > 25 {
		return plog.Print(head.Id, cfgEnum.ErrorCode_NameNeedLength)
	}

	if strPlayerName == fmt.Sprintf("%d", this.AccountId) {
		return plog.Print(head.Id, cfgEnum.ErrorCode_NameRepeat)
	}

	if strPlayerName == this.pbBase.Display.PlayerName {
		return plog.Print(head.Id, cfgEnum.ErrorCode_NameRepeat)
	}

	//检查屏蔽字库
	uCode := serverCommon.CheckMaskWord(&serverCommon.CheckMaskWordRequest{
		Content:    strPlayerName,
		SenderID:   int(this.AccountId),
		SenderName: this.pbBase.Display.PlayerName,
		SendTime:   int(base.GetNow())})
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(head.Id, uCode, this.pbBase.Display.PlayerName, strPlayerName)
	}

	// 判断是否第一次修改
	now := base.GetNow()
	if this.pbBase.LastModifyTime > 0 {
		// 判断cd冷却时间是否超过
		cdTime := cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_CHANGE_NAME_COOLDOWN_TIME)
		if diff := now - this.pbBase.LastModifyTime; diff < uint64(cdTime) {
			return plog.Print(head.Id, cfgEnum.ErrorCode_CoolDown, uint64(cdTime), now-this.pbBase.LastModifyTime)
		}

		// 判断是否存在足够改名卡
		bagFun := this.getPlayerBagFun()
		itemID := cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_CHANGE_NAME_ITEM)
		itemCount := bagFun.GetItemCount(uint32(cfgEnum.ESystemType_Item), itemID)
		if itemCount <= 0 {
			return plog.Print(head.Id, cfgEnum.ErrorCode_ItemNotEnough, itemID, itemCount)
		}

		uCode := this.updatePlayerName(strPlayerName)
		if uCode != cfgEnum.ErrorCode_Success {
			return plog.Print(head.Id, uCode, strPlayerName)
		}
		// 消耗改名卡
		errCode := bagFun.DelItem(head, uint32(cfgEnum.ESystemType_Item), itemID, 1, pb.EmDoingType_EDT_ChangePlayerName)
		if errCode != cfgEnum.ErrorCode_Success {
			return plog.Print(head.Id, errCode, itemID, 1)
		}
	} else {
		uCode := this.updatePlayerName(strPlayerName)
		if uCode != cfgEnum.ErrorCode_Success {
			return plog.Print(head.Id, uCode, strPlayerName)
		}
	}

	//存储
	this.pbBase.LastModifyTime = now
	this.UpdateSave(true)

	// 回包
	cluster.SendToClient(head, &pb.ChangePlayerNameResponse{PacketHead: &pb.IPacket{}, PlayerName: strPlayerName}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// gm命令
func (this *PlayerBaseFun) GmFuncRequest(head *pb.RpcHead, pbRequest *pb.GmFuncRequest) {
	arrParam := make([]string, 0)
	for _, v := range pbRequest.Param {
		arrParam = append(arrParam, v)
	}

	emErrorCode := this.gmFunc(head, pbRequest.GmType, arrParam)
	cluster.SendToClient(head, &pb.GmFuncResponse{PacketHead: &pb.IPacket{}}, emErrorCode)
}

// gm命令
func (this *PlayerBaseFun) gmFunc(head *pb.RpcHead, uGmType uint32, arrParam []string) cfgEnum.ErrorCode {
	//判断gm是否开启
	if serverCommon.GetParamCfg(pb.EmGmParamType_GPT_GmOpen) <= 0 {
		return plog.Print(head.Id, cfgEnum.ErrorCode_SystemNoOpen, uGmType, arrParam)
	}

	emErrorCode := cfgEnum.ErrorCode(cfgEnum.ErrorCode_Success)
	switch pb.EmGmFuncType(uGmType) {
	case pb.EmGmFuncType_GFT_Charge:
		if len(arrParam) != 1 {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_GmParam, uGmType, arrParam)
		}

		emErrorCode = this.getPlayerSystemChargeFun().AddProductId(head, base.StringToUInt32(arrParam[0]))
	case pb.EmGmFuncType_GFT_AddItem:
		if len(arrParam) != 3 {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_GmParam, uGmType, arrParam)
		}

		emErrorCode = this.getPlayerBagFun().AddItem(head, base.StringToUInt32(arrParam[0]), base.StringToUInt32(arrParam[1]), base.StringToInt64(arrParam[2]), pb.EmDoingType_EDT_Gm, false)
	case pb.EmGmFuncType_GFT_AddEquip:
		if len(arrParam) < 3 {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_GmParam, uGmType, arrParam)
		}

		uCount := uint32(1)
		if len(arrParam) >= 4 {
			uCount = base.StringToUInt32(arrParam[3])
		}

		this.GetPlayerEquipmentFun().AddEquipments(head, base.StringToUInt32(arrParam[0]), uCount, base.StringToUInt32(arrParam[1]), base.StringToUInt32(arrParam[2]), pb.EmDoingType_EDT_Gm)

	case pb.EmGmFuncType_GFT_AddHero:
		if len(arrParam) < 2 {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_GmParam, uGmType, arrParam)
		}
		uCount := uint32(1)
		if len(arrParam) >= 3 {
			uCount = base.StringToUInt32(arrParam[2])
		}
		this.getPlayerHeroFun().AddHeros(head, base.StringToUInt32(arrParam[0]), base.StringToUInt32(arrParam[1]), uCount, pb.EmDoingType_EDT_Gm, true)
	case pb.EmGmFuncType_GFT_NB: //高级号
		this.Nb(head)
	case pb.EmGmFuncType_GFT_Relogin: //高级号
		pbResponse := &pb.AllPlayerInfoNotify{
			PacketHead: &pb.IPacket{},
		}

		//初始化玩家数据
		for key, fun := range *this.MapFun {
			if key == pb.PlayerDataType_Hero {
				continue
			}

			pbResponse.PlayerData = &pb.PBPlayerData{}
			pbResponse.Mark = 0
			pbResponse.Mark = base.SetBit32(pbResponse.Mark, uint32(key), true)
			fun.SaveDataToClient(pbResponse.PlayerData)
			cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)
		}

		//英雄放最后
		fun, ok := (*this.MapFun)[pb.PlayerDataType_Hero]
		if ok {
			pbResponse.PlayerData = &pb.PBPlayerData{}
			pbResponse.Mark = 0
			pbResponse.Mark = base.SetBit32(pbResponse.Mark, uint32(pb.PlayerDataType_Hero), true)
			fun.SaveDataToClient(pbResponse.PlayerData)
			cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)
			plog.Trace("loginSuccess %s: %v", pb.PlayerDataType_Hero.String(), pbResponse)
		}

		//完成标记
		pbResponse.Mark = 0
		pbResponse.PlayerData = &pb.PBPlayerData{}
		//plog.Debug("sendAllInfoToClient: %+v", pbResponse.PlayerData)
		cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)
	}

	return emErrorCode
}

// 牛逼
func (this *PlayerBaseFun) Nb(head *pb.RpcHead) {
	//加所有英雄
	mapHero := cfgData.GetAllCfgHero()
	for id, _ := range mapHero {
		this.getPlayerHeroFun().addHero(head, id, 15, pb.EmDoingType_EDT_Gm)
	}

	//系统解锁
	this.getPlayerSystemCommonFun().AllSystemOpen(head)

	//加所有道具 100000个
	cfgAllItem := cfgData.GetAllCfgItem()
	for id, _ := range cfgAllItem {
		uCount := 1000000

		this.getPlayerBagFun().AddItem(head, uint32(cfgEnum.ESystemType_Item), id, int64(uCount), pb.EmDoingType_EDT_Gm, false)
	}

	//加所有装备
	cfgAllEquip := cfgData.GetAllCfgEquipment()
	uMaxEquipStar := cfgData.GetCfgEquipmentMaxStar()
	for id, _ := range cfgAllEquip {
		this.GetPlayerEquipmentFun().AddEquipment(head, id, uint32(cfgEnum.EQuality_Platinum), uMaxEquipStar, true, pb.EmDoingType_EDT_Gm)
	}

	//精英到最后一关
	this.getPlayerSystemBattleNormalFun().FinishAllBattle(head)

	//挂机到最后一关
	this.getPlayerSystemBattleNormalFun().FinishAllBattle(head)

	//任务到最后一个
	this.getPlayerSystemTaskFun().FinishAllTask(head)
}

// 检查条件
func (this *PlayerBaseFun) CheckCondition(conditions []*common.ConditionInfo) (cfgEnum.ErrorCode, uint64, uint64) {
	if len(conditions) <= 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_ConditionCfgEmpty), 0, 0
	}
	uBeginTime := uint64(0)
	uEndTime := uint64(0)
	for _, condition := range conditions {
		switch cfgEnum.EConditionType(condition.Type) {
		case cfgEnum.EConditionType_OpenServerDay:
			uOpenServerDays := this.GetOpenServerDays()
			if uOpenServerDays < condition.ConditionOpenServer.BeginDay {
				return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedOpenServerDays, *condition), 0, 0
			}

			if condition.ConditionOpenServer.EndDay > 0 && uOpenServerDays > condition.ConditionOpenServer.EndDay {
				return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedOpenServerDays, *condition, uOpenServerDays), 0, 0
			}
			uRegTime := this.getPlayerBaseFun().GetOpenSeverTime()
			uBeginTime = uRegTime
			if condition.ConditionOpenServer.BeginDay > 0 {
				uBeginTime = uBeginTime + uint64(condition.ConditionOpenServer.BeginDay-1)*24*3600
			}

			if condition.ConditionOpenServer.EndDay > 0 {
				uEndTime = base.GetZeroTimestamp(uRegTime+uint64(condition.ConditionOpenServer.EndDay-1)*24*3600, 1) - 1
			}

		case cfgEnum.EConditionType_AllServerDay:
			uOpenServerDays := serverCommon.GetOpenServerDays()
			if uOpenServerDays < condition.ConditionOpenServer.BeginDay {
				return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedOpenServerDays, *condition), 0, 0
			}

			if condition.ConditionOpenServer.EndDay > 0 && uOpenServerDays > condition.ConditionOpenServer.EndDay {
				return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedOpenServerDays, *condition, uOpenServerDays), 0, 0
			}
		case cfgEnum.EConditionType_RegDay:
			uRegDays := this.getPlayerBaseFun().GetRegDays()

			if uRegDays < condition.ConditionRegDays.RegDays {
				return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedRegDays, *condition, uRegDays), 0, 0
			}
			if condition.ConditionRegDays.RegEndDays > 0 && uRegDays > condition.ConditionRegDays.RegEndDays {
				return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedRegDays, *condition, uRegDays), 0, 0
			}

			uRegTime := this.getPlayerBaseFun().GetRegTime()
			uBeginTime = uRegTime
			if condition.ConditionRegDays.RegDays > 0 {
				uBeginTime = uBeginTime + uint64(condition.ConditionRegDays.RegDays-1)*24*3600
			}

			if condition.ConditionRegDays.RegEndDays > 0 {
				uEndTime = base.GetZeroTimestamp(uRegTime+uint64(condition.ConditionRegDays.RegEndDays-1)*24*3600, 1) - 1
			}
		case cfgEnum.EConditionType_BattleMap:
			mapId, stageId := this.getPlayerSystemBattleFun().GetFinishMapIdAndStageId(pb.EmBattleType(condition.ConditionBattleMap.BattleType))
			if base.MakeU64(mapId, stageId) < base.MakeU64(condition.ConditionBattleMap.MapId, condition.ConditionBattleMap.StageId) {
				return plog.Print(this.AccountId, cfgEnum.ErrorCode_NeedBattleMap, *condition, mapId, stageId), 0, 0
			}
		case cfgEnum.EConditionType_SystemOpenType:
			if !this.getPlayerSystemCommonFun().CheckSystemTypeOpen(cfgEnum.ESystemUnlockType(condition.ConditionOpenServerType.SystemOpenType)) {
				return plog.Print(this.AccountId, cfgEnum.ErrorCode_SystemNoOpen, *condition), 0, 0
			}
		}
	}

	return cfgEnum.ErrorCode_Success, uBeginTime, uEndTime
}
