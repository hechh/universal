package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common/cfgData"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/playerNew/domain"
	"corps/server/game/playerNew/internal/manager"

	"github.com/spf13/cast"
)

type PlayerFun struct {
	*PlayerFunFactory
	accountId uint64            // 玩家uid
	dataType  pb.PlayerDataType // 数据类型
	isSave    bool              // 是否需要保存
}

func NewPlayerFun(uid uint64, dataType pb.PlayerDataType, mgr *manager.FunMgr) *PlayerFun {
	return &PlayerFun{
		accountId:        uid,
		dataType:         dataType,
		PlayerFunFactory: &PlayerFunFactory{FunMgr: mgr},
	}
}

func (d *PlayerFun) GetPlayerDataType() pb.PlayerDataType {
	return d.dataType
}

// 注册定时器
func (d *PlayerFun) RegisterTimer() {}

// 加载数据(非system数据)
func (d *PlayerFun) Load([]byte) error {
	return nil
}

// db数据加载完成回调，在LoadComplete之前调用
func (d *PlayerFun) LoadPlayerDBFinish() {}

// 初始化新开启的模块数据，在LoadComplete之后调用
/*
func (d *PlayerFun) NewPlayer() error {
	return nil
}
*/

// 加载完成，在NewPlayer之后调用
func (d *PlayerFun) LoadComplete() {}

// 心跳包
func (d *PlayerFun) Heat() {}

// 是否跨天
func (d *PlayerFun) PassDay(isDay, isWeek, isMonth bool) {}

// 设置数据
func (d *PlayerFun) SetUserTypeInfo([]byte) error {
	return nil
}

// 深度拷贝数据
func (d *PlayerFun) CopyTo(pbData *pb.PBPlayerData) error {
	return nil
}

// 设置保存状态
func (d *PlayerFun) UpdateSave(bSave bool) {
	d.isSave = bSave
}

// 判断是否存储数据
func (d *PlayerFun) IsSave() bool {
	return d.isSave
}

// ---------------------------------公用接口----------------------------------------
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
		val = (val << 32) | uint64(domain.MAX_DWORD_VALUE-uint32(now))
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
		Id:        this.accountId,
		RegionID:  this.GetPlayerBaseFun().GetServerId(),
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
		Member:     cast.ToString(this.accountId),
		Score:      float64(val),
	}
	// 发送gm消息，更新排行榜
	this.GetPlayerSystemChampionshipFun().SetChampionshipFlag(rankType, 1)
	cluster.SendToGm(head, "RankMgr", "UpdateRequest", notify)
}
