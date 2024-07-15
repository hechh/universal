package router

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"
	"universal/common/pb"
)

type RouteInfo struct {
	ServerType int32
	ServerID   uint32
}

type RouteList struct {
	uid        uint64       // 玩家uid
	updateTime int64        // 更新时间
	list       []*RouteInfo // 该玩家的路由信息
}

func (d *RouteList) String() string {
	buf, _ := json.Marshal(d.list)
	return fmt.Sprintf("%d: %s", d.updateTime, string(buf))
}

// 查询路由节点
func (d *RouteList) GetRouteInfo(srvType int32) (dst *RouteInfo) {
	for _, routeritem := range d.list {
		if routeritem.ServerType == srvType {
			dst = routeritem
			break
		}
	}
	return
}

// 创建路由节点
func (d *RouteList) NewRouteInfo(srvType int32) (dst *RouteInfo) {
	for _, routeritem := range d.list {
		if routeritem.ServerType == srvType {
			dst = routeritem
			break
		}
	}
	if dst == nil {
		dst = &RouteInfo{ServerType: srvType}
		d.list = append(d.list, dst)
	}
	return
}

// 更新玩家路由信息
func (d *RouteList) UpdateRouteInfo(head *pb.PacketHead, node *pb.ServerNode) {
	// 更新路由信息
	dst := d.NewRouteInfo(int32(head.DstServerType))
	dst.ServerID = node.ServerID
	// 更新时间
	atomic.StoreInt64(&d.updateTime, time.Now().Unix())
}
