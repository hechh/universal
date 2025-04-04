package util

import "hego/common/pb"

func CopyHead(head *pb.Head) *pb.Head {
	new := *head
	newRouteInfo := *head.Table
	new.Table = &newRouteInfo
	return &new
}
