package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/framework/cluster"
	"corps/framework/common/uerror"
	"corps/framework/plog"
	"corps/pb"
	"encoding/json"

	"github.com/golang/protobuf/proto"
)

// gomaker生成模板

type PlayerActivityAdventure struct {
	*PlayerSystemActivityFun
	emType  cfgEnum.EActivityType
	mapData map[uint32]*pb.PBActivityAdventure
}

func init() {
	RegisterActivity(cfgEnum.EActivityType_Adventure, func() IPlayerSystemActivityFun { return new(PlayerActivityAdventure) })
}

func (this *PlayerActivityAdventure) Init(emType cfgEnum.EActivityType, pFun *PlayerSystemActivityFun) {
	this.mapData = make(map[uint32]*pb.PBActivityAdventure)
	this.PlayerSystemActivityFun = pFun
	this.emType = emType
}

func (this *PlayerActivityAdventure) LoadData(pbData *pb.PBPlayerSystemActivity) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemActivity{}
	}

	this.mapData = make(map[uint32]*pb.PBActivityAdventure)
	for _, adventure := range pbData.AdventureList {
		this.mapData[adventure.Id] = adventure
	}
}
func (this *PlayerActivityAdventure) LoadComplete() {
}

// 存储到数据库
func (this *PlayerActivityAdventure) SaveData(pbData *pb.PBPlayerSystemActivity) {
	for _, data := range this.mapData {
		pbData.AdventureList = append(pbData.AdventureList, data)
	}
}

// 活动过期删除
func (this *PlayerActivityAdventure) Del(uID uint32) {
	pbActivity, ok := this.mapData[uID]
	if !ok {
		plog.Error("(this *PlayerActivityGrowRoad) sid: %d", uID)
		return
	}

	// 获取计算时间
	now := base.GetNow()
	registerTime := this.getPlayerBaseFun().GetRegTime()
	if registerTime >= now {
		return
	}

	//经过的天数
	regDays := base.DiffDays(registerTime, now) + 1
	//最后一天的时间
	lastDayPassMin := this.getLastDayPassMin()

	arrItem := make([]*common.ItemInfo, 0)
	listAllCfg := cfgData.GetAllCfgAdventureConfig(uID)
	for _, cfg := range listAllCfg {
		if base.ArrayContainsValue(pbActivity.PrizeIdList, cfg.Id) {
			continue
		}

		if regDays > cfg.Day {
			arrItem = append(arrItem, cfg.Reward)
		} else if regDays == cfg.Day {
			if lastDayPassMin >= cfg.Expire {
				arrItem = append(arrItem, cfg.Reward)
			}
		}
	}

	//补发奖励邮件
	if len(arrItem) > 0 {
		cfgAcitivty := cfgData.GetCfgActivityConfig(uID)
		if cfgAcitivty != nil {
			this.getPlayerMailFun().AddTempMail(&pb.RpcHead{Id: this.AccountId}, cfgEnum.EMailId_BP, pb.EmDoingType_EDT_BP, arrItem, cfgAcitivty.Name)

		}
	}

	delete(this.mapData, uID)
	this.UpdateSave(true)
}

// 活动上线，动态添加
func (this *PlayerActivityAdventure) Add(uID uint32, uBeginTime uint64, uEndTime uint64) {
	if _, ok := this.mapData[uID]; ok {
		plog.Error("(this *PlayerActivityAdventure) Add repeated %d %d %d", uID, uBeginTime, uEndTime)
		return
	}
	item := &pb.PBActivityAdventure{
		Id:           uID,
		BeginTime:    uBeginTime,
		EndTime:      uEndTime,
		RegisterTime: this.getPlayerBaseFun().GetRegTime(),
	}
	this.mapData[uID] = item
	//通知客户端
	pbNotify := &pb.ActivityDataNewNotify{
		PacketHead:   &pb.IPacket{},
		ActivityType: uint32(this.emType),
		Info:         &pb.PBPlayerSystemActivity{},
	}
	pbNotify.Info.AdventureList = append(pbNotify.Info.AdventureList, this.mapData[uID])

	// 发送数据
	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbNotify, cfgEnum.ErrorCode_Success)
	this.UpdateSave(true)
}
func (this *PlayerActivityAdventure) FreePrize(head *pb.RpcHead, uID uint32) cfgEnum.ErrorCode {
	return cfgEnum.ErrorCode_NoData
}

// 判断是否有红点
func (this *PlayerActivityAdventure) GetRed(sid uint32) bool {
	// 判单活动数据是否存在
	pbActivity, ok := this.mapData[sid]
	if !ok {
		return false
	}
	// 获取计算时间
	now := base.GetNow()
	registerTime := this.getPlayerBaseFun().GetRegTime()
	if registerTime >= now {
		return false
	}

	//经过的天数
	regDays := base.DiffDays(registerTime, now) + 1
	//最后一天的时间
	lastDayPassMin := this.getLastDayPassMin()

	listAllCfg := cfgData.GetAllCfgAdventureConfig(sid)
	for _, cfg := range listAllCfg {
		if base.ArrayContainsValue(pbActivity.PrizeIdList, cfg.Id) {
			continue
		}

		if regDays > cfg.Day {
			return true
		} else if regDays == cfg.Day {
			if lastDayPassMin >= cfg.Expire {
				return true
			}
		}

		return false

	}

	return false
}
func (this *PlayerActivityAdventure) getLastDayPassMin() uint32 {
	// 获取计算时间
	now := base.GetNow()
	registerTime := this.getPlayerBaseFun().GetRegTime()
	if registerTime >= now {
		return 0
	}

	//经过的天数
	regDays := base.DiffDays(registerTime, now) + 1
	//最后一天的时间
	lastDayPassMin := uint32(0)
	if regDays > 1 {
		lastDayPassMin = uint32(now-base.GetZeroTimestamp(registerTime, int32(regDays-1))) / 60
	} else {
		lastDayPassMin = uint32(now-registerTime) / 60
	}

	return lastDayPassMin
}
func (this *PlayerActivityAdventure) Open(sid uint32) (rets []string) {
	if cfgs := cfgData.GetCfgAdventureConfigList(sid); len(cfgs) > 0 {
		buf, err := json.Marshal(cfgs)
		if err != nil {
			return
		}
		rets = append(rets, string(buf))
		return
	}
	return
}

// --------------------交互接口实现------------------------------

func (this *PlayerActivityAdventure) AdventureRewardRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.AdventureRewardRequest)
	// 判单sid是否存在
	pbActivity, ok := this.mapData[req.ID]
	if !ok {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_NoActivity, "head: %v, req: %v", head, req)
	}

	cfgAdventure := cfgData.GetCfgAdventureConfig(req.CfgID)
	if cfgAdventure == nil {
		return uerror.NewUErrorf(1, cfgData.GetAdventureConfigErrorCode(req.CfgID), "head: %v, req: %v", head, req)
	}

	//已经领取过
	if base.ArrayContainsValue(pbActivity.PrizeIdList, req.CfgID) {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_HavePrize, "head: %v, req: %v", head, req)
	}

	// 获取计算时间
	now := base.GetNow()
	registerTime := this.getPlayerBaseFun().GetRegTime()
	if registerTime >= now {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_RegisterTimeExceedNow, "head: %v, req: %v", head, req)
	}

	//经过的天数
	regDays := base.DiffDays(registerTime, now) + 1

	//判断条件
	if regDays < cfgAdventure.Day {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_AdventureRewardParam, "head: %v, req: %v", head, req)
	} else if regDays == cfgAdventure.Day {
		//最后一天的时间
		lastDayPassMin := uint32(0)
		if regDays > 1 {
			lastDayPassMin = uint32(now-base.GetZeroTimestamp(registerTime, int32(regDays-1))) / 60
		} else {
			lastDayPassMin = uint32(now-registerTime) / 60
		}

		if lastDayPassMin < cfgAdventure.Expire {
			return uerror.NewUErrorf(1, cfgEnum.ErrorCode_AdventureRewardLocked, "head: %v, req: %v", head, req)
		}
	}

	pbActivity.PrizeIdList = append(pbActivity.PrizeIdList, req.CfgID)

	// 发送道具
	errCode := this.getPlayerBagFun().AddOneArrItem(head, cfgAdventure.Reward, pb.EmDoingType_EDT_Adventure, true)
	if errCode != cfgEnum.ErrorCode_Success {
		return uerror.NewUErrorf(1, errCode, "head: %v, req: %v", head, req)
	}
	this.UpdateSave(true)
	return nil
}
