package framework

import (
	"universal/framework/config"
	"universal/framework/domain"
	"universal/framework/internal/actor"
	"universal/library/baselib/uerror"
)

type Actor struct{ actor.Actor }
type ActorMgr struct{ actor.ActorMgr }

type Framework struct {
	nodeType int32
	nodeId   int32
	rte      domain.IRouteMgr
	cls      domain.ICluster
	dis      domain.IDiscovery
	net      domain.INetwork
	actors   map[string]domain.IActor
	newNode  func() domain.INode
	newHead  func() domain.IHead
	newRoute func() domain.IRoute
	newPack  func() domain.IPacket
}

func (f *Framework) Init(cfg *config.Config, nodeType domain.NodeType, nodeId int32) error {
	f.nodeType = int32(nodeType)
	f.nodeId = nodeId
	nodecfg, ok := cfg.Cluster[domain.NodeType_name[int32(nodeType)]]
	if !ok || nodecfg == nil {
		return uerror.New(1, -1, "服务节点配置不存在：%s", domain.NodeType_name[int32(nodeType)])
	}

	return nil
}
