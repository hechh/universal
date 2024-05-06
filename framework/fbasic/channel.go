package fbasic

import (
	"fmt"
	"strings"
	"universal/common/pb"

	"github.com/spf13/cast"
)

const (
	rootDir = "/server"
)

func GetRootDir() string {
	return rootDir
}

// 玩家消息
func GetPlayerChannel(typ pb.ClusterType, clusterID uint32, uid uint64) string {
	return fmt.Sprintf("%s/%s/%d/%d", rootDir, strings.ToLower(typ.String()), clusterID, uid)
}

// 服务节点消息
func GetNodeChannel(typ pb.ClusterType, clusterID uint32) string {
	return fmt.Sprintf("%s/%s/%d", rootDir, strings.ToLower(typ.String()), clusterID)
}

// 所有节点消息
func GetClusterChannel(typ pb.ClusterType) string {
	return fmt.Sprintf("%s/%s", rootDir, strings.ToLower(typ.String()))
}

// 解析channel
func ParseChannel(str string) (clusterType pb.ClusterType, clusterID uint32, uid uint64) {
	strs := strings.Split(strings.TrimPrefix(str, rootDir), "/")
	// 解析clusterType
	if len(strs) > 0 {
		clusterType = pb.ClusterType(pb.ClusterType_value[strings.ToUpper(strs[0])])
	}
	// 解析clusterID
	if len(strs) > 1 {
		clusterID = cast.ToUint32(strs[1])
	}
	// 解析玩家uid
	if len(strs) > 2 {
		uid = cast.ToUint64(strs[2])
	}
	return
}
