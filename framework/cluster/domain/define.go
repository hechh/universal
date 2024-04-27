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

type WatchFunc func(key string, value []byte)

type IDiscovery interface {
	KeepAlive(string, []byte, int64)    // 设置保活key
	Walk(string, WatchFunc) error       // 便利所有key-value
	Watch(string, WatchFunc, WatchFunc) // 监听所有变更
}

type ClusterFunc func(*pb.Packet)

// etcd
func GetNodeChannel(typ pb.ClusterType, clusterID uint32) string {
	return fmt.Sprintf("server/%s/%d", strings.ToLower(typ.String()), clusterID)
}

func GetTopicChannel(typ pb.ClusterType) string {
	return fmt.Sprintf("server/%s", strings.ToLower(typ.String()))
}
