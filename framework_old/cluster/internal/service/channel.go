package service

import (
	"fmt"
	"path"
	"strings"
	"universal/common/pb"
	"universal/framework/common/uerror"

	"github.com/spf13/cast"
)

const (
	ROOT_DIR = "server"
)

// 玩家消息
func GetPlayerChannel(typ pb.ServerType, serverID uint32, uid uint64) string {
	return path.Join(ROOT_DIR, strings.ToLower(typ.String()), fmt.Sprintf("%d/%d", serverID, uid))
}

// 服务节点消息
func GetNodeChannel(typ pb.ServerType, serverID uint32) string {
	return path.Join(ROOT_DIR, strings.ToLower(typ.String()), cast.ToString(serverID))
}

// 所有节点消息
func GetClusterChannel(typ pb.ServerType) string {
	return path.Join(ROOT_DIR, strings.ToLower(typ.String()))
}

// 获取channel的key
func GetHeadChannel(head *pb.PacketHead) (str string, err error) {
	switch head.SendType {
	case pb.SendType_PLAYER:
		str = GetPlayerChannel(head.DstServerType, head.DstServerID, head.UID)
	case pb.SendType_NODE:
		str = GetNodeChannel(head.DstServerType, head.DstServerID)
	case pb.SendType_CLUSTER:
		str = GetClusterChannel(head.DstServerType)
	default:
		err = uerror.NewUErrorf(1, -1, "%v", head)
	}
	return
}

// 解析channel
func ParseChannel(str string) (serverType pb.ServerType, serverID uint32, uid uint64) {
	strs := strings.Split(str, "/")
	// 解析serverType
	if len(strs) > 2 {
		serverType = pb.ServerType(pb.ServerType_value[strings.ToUpper(strs[1])])
	}
	// 解析serverID
	if len(strs) > 3 {
		serverID = cast.ToUint32(strs[2])
	}
	// 解析玩家uid
	if len(strs) > 4 {
		uid = cast.ToUint64(strs[3])
	}
	return
}
