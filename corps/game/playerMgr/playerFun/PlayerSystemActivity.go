package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"

	"github.com/golang/protobuf/proto"
)

type NewActivityFun func() IPlayerSystemActivityFun

var (
	activitys = make(map[cfgEnum.EActivityType]NewActivityFun)
)

func RegisterActivity(activityType cfgEnum.EActivityType, f NewActivityFun) {
	activitys[activityType] = f
}

type (
	PlayerSystemActivityFun struct {
		PlayerFun
		mapActivity map[uint32]*PlayerActivityInfo
		mapFun      map[cfgEnum.EActivityType]IPlayerSystemActivityFun
	}

	PlayerActivityInfo struct {
		*pb.PBPlayerActivityInfo
		Send    bool //是否已经发送过
		HaveRed bool //是否有红点
	}

	ActivityRed struct {
		Send    bool //是否已经发送过
		HaveRed bool //是否有红点
	}

	IPlayerSystemActivityFun interface {
		Init(emType cfgEnum.EActivityType, pFun *PlayerSystemActivityFun)
		LoadData(pbData *pb.PBPlayerSystemActivity)
		SaveData(pbData *pb.PBPlayerSystemActivity)
		LoadComplete()
		Add(uID uint32, uBeginTime uint64, uEndTime uint64)
		Del(uID uint32)
		GetRed(uID uint32) bool
		Open(uID uint32) []string
		FreePrize(head *pb.RpcHead, uID uint32) cfgEnum.ErrorCode
	}
)

func (this *PlayerSystemActivityFun) Init(pbType pb.PlayerDataType, pcommon *FunCommon) {
	this.PlayerFun.Init(pbType, pcommon)
	this.mapActivity = make(map[uint32]*PlayerActivityInfo)
	this.RegisterFun()
}

// 注册
func (this *PlayerSystemActivityFun) RegisterFun() {
	this.mapFun = make(map[cfgEnum.EActivityType]IPlayerSystemActivityFun)
	for i, f := range activitys {
		data := f()
		data.Init(i, this)
		this.mapFun[i] = data
	}
}

// 从数据库中加载
func (this *PlayerSystemActivityFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	this.loadData(pbSystem.Activity)
	this.UpdateSave(false)
}

func (this *PlayerSystemActivityFun) loadData(pbData *pb.PBPlayerSystemActivity) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemActivity{}
	}

	this.mapActivity = make(map[uint32]*PlayerActivityInfo)
	for _, v := range pbData.ActivityList {
		this.mapActivity[v.ActivityId] = &PlayerActivityInfo{
			PBPlayerActivityInfo: v,
			Send:                 false,
			HaveRed:              false,
		}
	}

	for _, info := range this.mapFun {
		info.LoadData(pbData)
	}

	this.UpdateSave(true)
}

// 加载完成
func (this *PlayerSystemActivityFun) LoadComplete() {
	//重置红点
	for _, pActivity := range this.mapActivity {
		pActivity.Send = false
		pActivity.HaveRed = false
	}

	this.CheckActivity()

	for _, info := range this.mapFun {
		info.LoadComplete()
	}
}

// 存储到数据库
func (this *PlayerSystemActivityFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	if pbSystem.Activity == nil {
		pbSystem.Activity = new(pb.PBPlayerSystemActivity)
	}

	for _, v := range this.mapActivity {
		pbSystem.Activity.ActivityList = append(pbSystem.Activity.ActivityList, v.PBPlayerActivityInfo)
	}

	for _, info := range this.mapFun {
		info.SaveData(pbSystem.Activity)
	}

	return this.BSave
}
func (this *PlayerSystemActivityFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemActivity{}
}

// 设置玩家数据
func (this *PlayerSystemActivityFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}
	pbSystem := pbData.(*pb.PBPlayerSystemActivity)
	if pbSystem == nil {
		return false
	}
	this.loadData(pbSystem)
	this.UpdateSave(true)
	return true
}

func (this *PlayerSystemActivityFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = new(pb.PBPlayerSystem)
	}

	this.SaveSystem(pbData.System)
}

func (this *PlayerSystemActivityFun) getActivityFun(emType cfgEnum.EActivityType) IPlayerSystemActivityFun {
	return this.mapFun[emType]
}

func (this *PlayerSystemActivityFun) GetActivityGrowRoadFun() *PlayerActivityGrowRoad {
	return this.getActivityFun(cfgEnum.EActivityType_GrowRoad).(*PlayerActivityGrowRoad)
}

func (this *PlayerSystemActivityFun) GetActivityChargeGiftFun() *PlayerActivityChargeGift {
	return this.getActivityFun(cfgEnum.EActivityType_ChargeBuy).(*PlayerActivityChargeGift)
}

func (this *PlayerSystemActivityFun) GetActivityAdventureFun() *PlayerActivityAdventure {
	return this.getActivityFun(cfgEnum.EActivityType_Adventure).(*PlayerActivityAdventure)
}
func (this *PlayerSystemActivityFun) GetActivityOpenServerGiftFun() *PlayerActivityOpenServerGift {
	return this.getActivityFun(cfgEnum.EActivityType_OpenServerGift).(*PlayerActivityOpenServerGift)
}

// 跨天
func (this *PlayerSystemActivityFun) PassDay(isDay, isWeek, isMonth bool) {
	this.CheckActivity()
}

// 心跳包
func (this *PlayerSystemActivityFun) Heat() {
	pbRedNotify := &pb.ActivityRedNotify{PacketHead: &pb.IPacket{Id: this.AccountId}}
	for _, pActivity := range this.mapActivity {
		cfgActiviy := cfgData.GetCfgActivityConfig(pActivity.ActivityId)
		if cfgActiviy == nil {
			continue
		}

		fun := this.getActivityFun(cfgEnum.EActivityType(cfgActiviy.ActivityType))
		if fun == nil {
			continue
		}

		bSend := false
		if !pActivity.HaveRed {
			pActivity.HaveRed = fun.GetRed(cfgActiviy.Sid)
			if pActivity.HaveRed {
				bSend = true
			}
		} else if !pActivity.Send {
			bSend = true
		}

		if bSend {
			pActivity.Send = true
			pbRedNotify.IdList = append(pbRedNotify.IdList, cfgActiviy.Sid)
		}
	}

	if len(pbRedNotify.IdList) > 0 {
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbRedNotify, cfgEnum.ErrorCode_Success)
	}
}

// 检查活动开启
func (this *PlayerSystemActivityFun) CheckActivity() {
	//判断活动
	pbListNotify := &pb.ActivityListNotify{PacketHead: &pb.IPacket{Id: this.AccountId}}
	mapCfg := cfgData.GetAllCfgActivityConfig()
	for _, cfg := range mapCfg {
		uCode, uBeginTime, uEndTime := this.getPlayerBaseFun().CheckCondition([]*common.ConditionInfo{cfg.ConditionInfo})

		bOpen := uCode == cfgEnum.ErrorCode_Success
		if !bOpen {
			if _, ok := this.mapActivity[cfg.Sid]; ok {
				pbListNotify.DelIdList = append(pbListNotify.DelIdList, cfg.Sid)
				delete(this.mapActivity, cfg.Sid)
			}
			continue
		}

		//新开
		if _, ok := this.mapActivity[cfg.Sid]; !ok {
			if bOpen {
				this.mapActivity[cfg.Sid] = &PlayerActivityInfo{
					PBPlayerActivityInfo: &pb.PBPlayerActivityInfo{
						ActivityId: cfg.Sid,
						BeginTime:  uBeginTime,
						EndTime:    uEndTime,
					},
				}

				pbListNotify.ActivityList = append(pbListNotify.ActivityList, this.mapActivity[cfg.Sid].PBPlayerActivityInfo)
			}
		}
	}

	//同步列表变化
	if len(pbListNotify.ActivityList) > 0 || len(pbListNotify.DelIdList) > 0 {
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbListNotify, cfgEnum.ErrorCode_Success)
	}

	//判断活动
	for _, info := range pbListNotify.ActivityList {
		cfgActivity := cfgData.GetCfgActivityConfig(info.ActivityId)
		if cfgActivity == nil {
			continue
		}

		this.OnAddActivity(cfgData.GetCfgActivityConfig(info.ActivityId))
	}

	for _, id := range pbListNotify.DelIdList {
		this.OnDelActivity(cfgData.GetCfgActivityConfig(id))
	}
}
func (this *PlayerSystemActivityFun) OnAddActivity(cfgActivity *cfgData.ActivityConfigCfg) {
	info, ok := this.mapActivity[cfgActivity.Sid]
	if !ok {
		return
	}

	fun := this.getActivityFun(cfgEnum.EActivityType(cfgActivity.ActivityType))
	if fun == nil {
		return
	}

	fun.Add(cfgActivity.Sid, info.BeginTime, info.EndTime)
}

func (this *PlayerSystemActivityFun) OnDelActivity(cfgActivity *cfgData.ActivityConfigCfg) {
	if cfgActivity == nil {
		return
	}

	fun := this.getActivityFun(cfgEnum.EActivityType(cfgActivity.ActivityType))
	if fun == nil {
		return
	}

	fun.Del(cfgActivity.Sid)
}

// 活动打开请求
func (this *PlayerSystemActivityFun) ActivityOpenRequest(head *pb.RpcHead, pbRequest *pb.ActivityOpenRequest) {
	uCode := this.ActivityOpen(head, pbRequest.Id)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ActivityOpenResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 活动打开请求
func (this *PlayerSystemActivityFun) ActivityOpen(head *pb.RpcHead, uId uint32) cfgEnum.ErrorCode {
	cfgActivity := cfgData.GetCfgActivityConfig(uId)
	if cfgActivity == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, uId)
	}

	fun := this.getActivityFun(cfgEnum.EActivityType(cfgActivity.ActivityType))
	if fun == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, uId)
	}

	if _, ok := this.mapActivity[cfgActivity.Sid]; !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, uId)
	}

	//通知客户端
	cluster.SendToClient(head, &pb.ActivityOpenResponse{
		PacketHead: &pb.IPacket{},
		Id:         uId,
		JsonData:   fun.Open(uId),
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 活动免费奖励请求
func (this *PlayerSystemActivityFun) ActivityFreePrizeRequest(head *pb.RpcHead, pbRequest *pb.ActivityFreePrizeRequest) {
	uCode := this.ActivityFreePrize(head, pbRequest.Id)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.ActivityOpenResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 活动免费奖励请求
func (this *PlayerSystemActivityFun) ActivityFreePrize(head *pb.RpcHead, uId uint32) cfgEnum.ErrorCode {
	cfgActivity := cfgData.GetCfgActivityConfig(uId)
	if cfgActivity == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, uId)
	}

	fun := this.getActivityFun(cfgEnum.EActivityType(cfgActivity.ActivityType))
	if fun == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, uId)
	}

	if _, ok := this.mapActivity[cfgActivity.Sid]; !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoActivity, uId)
	}

	uCode := fun.FreePrize(head, uId)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(this.AccountId, uCode, uId)
	}

	//通知客户端
	cluster.SendToClient(head, &pb.ActivityFreePrizeResponse{
		PacketHead:         &pb.IPacket{},
		Id:                 uId,
		NextDailyPrizeTime: base.GetZeroTimestamp(base.GetNow(), 1),
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}
