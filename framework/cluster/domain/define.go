package domain

import (
	"fmt"
	"strings"
	"universal/common/pb"
)

const (
	ActionTypeNone = 0
	ActionTypeAdd  = 1
	ActionTypeDel  = 2
)

type DiscoveryFunc func(int, *pb.ClusterNode)

type IDiscovery interface {
	GetSelf() *pb.ClusterNode         // 获取自身节点
	Watch(string, DiscoveryFunc)      // 监听所有变更
	Walk(string, DiscoveryFunc) error // 便利所有key-value
}

type ClusterFunc func(*pb.Packet)

// etcd
func GetNodeChannel(typ pb.ClusterType, clusterID uint32) string {
	return fmt.Sprintf("server/%s/%d", strings.ToLower(typ.String()), clusterID)
}

func GetTopicChannel(typ pb.ClusterType) string {
	return fmt.Sprintf("server/%s", strings.ToLower(typ.String()))
}
