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

	ROOT_DIR = "server/cluster/"
)

type WatchFunc func(action int, key string, value string)

type IDiscovery interface {
	KeepAlive(string, string, int64) // 设置保活key
	Watch(string, WatchFunc) error   // 开启协程watch+keepalive
	Close()                          // 停止协程watch+keepalive
}

type ClusterFunc func(*pb.Packet)

// etcd
func GetNodeChannel(typ pb.ClusterType, clusterID uint32) string {
	return fmt.Sprintf(ROOT_DIR+"%s/%d", strings.ToLower(typ.String()), clusterID)
}

func GetTopicChannel(typ pb.ClusterType) string {
	return fmt.Sprintf(ROOT_DIR+"%s", strings.ToLower(typ.String()))
}
