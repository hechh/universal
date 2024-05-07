package framework

import (
	"universal/common/pb"
)

var (
	serverId    int
	clusterType pb.ClusterType
)

func SetGlobal(id int, t pb.ClusterType) {
	serverId = id
	clusterType = t
}

func GetServerID() int {
	return serverId
}

func GetClusterType() pb.ClusterType {
	return clusterType
}
