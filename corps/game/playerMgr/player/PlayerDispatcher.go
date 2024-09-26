package player

import (
	"context"
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common/cfgData"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
)

func (this *Player) RankRequest(ctx context.Context, req *pb.RankRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.RankResponse{}
	// 加载配置
	rankInfoCfg := cfgData.GetCfgRankInfoConfig(req.RankType)
	rankTypeCfg := cfgData.GetCfgRankTypeConfig(req.RankType)
	if rankTypeCfg == nil {
		plog.Error("rankTypeCondig(%d) not found", req.RankType)
		cluster.SendToClient(head, rsp, cfgData.GetRankTypeConfigErrorCode(req.RankType))
		return
	}
	// 设置路由规则
	head.RouteType = uint32(cfgEnum.ERouteType_ServerID)
	switch rankTypeCfg.DataType {
	case uint32(cfgEnum.EDataType_Forever): // 永久榜单
		switch rankTypeCfg.RouteType {
		case uint32(cfgEnum.ERedisType_Global): // 全局节点
			head.RegionID = 0
		case uint32(cfgEnum.ERedisType_Random): // 负载均衡
			head.RegionID = this.getPlayerBaseFun().GetServerId()
		}
		req.CreateTime = 0
	case uint32(cfgEnum.EDataType_Expire): // 限时删除榜单
		switch rankTypeCfg.RouteType {
		case uint32(cfgEnum.ERedisType_Global): // 全局节点
			head.RegionID = 0
		case uint32(cfgEnum.ERedisType_Random): // 负载均衡
			head.RegionID = this.getPlayerBaseFun().GetServerId()
		}
		if rankInfoCfg == nil {
			plog.Error("rankTypeCondig(%d) not found", req.RankType)
			cluster.SendToClient(head, &pb.RankResponse{}, cfgData.GetRankInfoConfigErrorCode(req.RankType))
			return
		}
		// 获取榜单开启时间
		createTime := this.getPlayerBaseFun().GetServerStartTime()
		if req.RankType == uint32(cfgEnum.ERankType_WorldBoss) {
			createTime = base.GetZeroTimestamp(base.GetNow(), 0)
		}
		req.CreateTime = cfgData.GetCfgRankActiveTime(rankInfoCfg, createTime)
	}
	cluster.SendToGm(head, "RankMgr", "RankRequest", req)
}

func (this *Player) RankRewardRequest(ctx context.Context, req *pb.RankRewardRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.RankResponse{}
	// 加载配置
	rankInfoCfg := cfgData.GetCfgRankInfoConfig(req.RankType)
	rankTypeCfg := cfgData.GetCfgRankTypeConfig(req.RankType)
	if rankTypeCfg == nil {
		plog.Error("rankTypeCondig(%d) not found", req.RankType)
		cluster.SendToClient(head, rsp, cfgData.GetRankTypeConfigErrorCode(req.RankType))
		return
	}
	// 设置路由规则
	head.RouteType = uint32(cfgEnum.ERouteType_ServerID)
	switch rankTypeCfg.DataType {
	case uint32(cfgEnum.EDataType_Forever): // 永久榜单
		switch rankTypeCfg.RouteType {
		case uint32(cfgEnum.ERedisType_Global): // 全局节点
			head.RegionID = 0
		case uint32(cfgEnum.ERedisType_Random): // 负载均衡
			head.RegionID = this.getPlayerBaseFun().GetServerId()
		}
		req.CreateTime = 0
	case uint32(cfgEnum.EDataType_Expire): // 限时删除榜单
		switch rankTypeCfg.RouteType {
		case uint32(cfgEnum.ERedisType_Global): // 全局节点
			head.RegionID = 0
		case uint32(cfgEnum.ERedisType_Random): // 负载均衡
			head.RegionID = this.getPlayerBaseFun().GetServerId()
		}
		if rankInfoCfg == nil {
			plog.Error("rankTypeCondig(%d) not found", req.RankType)
			cluster.SendToClient(head, &pb.RankResponse{}, cfgData.GetRankInfoConfigErrorCode(req.RankType))
			return
		}
		// 获取榜单开启时间
		createTime := this.getPlayerBaseFun().GetServerStartTime()
		if req.RankType == uint32(cfgEnum.ERankType_WorldBoss) {
			createTime = base.GetZeroTimestamp(base.GetNow(), 0)
		}
		req.CreateTime = cfgData.GetCfgRankActiveTime(rankInfoCfg, createTime)
	}
	cluster.SendToGm(head, "RankMgr", "RankRewardRequest", req)
}
