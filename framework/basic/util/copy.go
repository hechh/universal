package util

import "universal/common/pb"

func CopyRpcHead(head *pb.RpcHead) *pb.RpcHead {
	new := *head
	newRouteInfo := *head.Route
	new.Route = &newRouteInfo
	return &new
}
