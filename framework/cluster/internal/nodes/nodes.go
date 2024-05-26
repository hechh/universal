package nodes

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"universal/common/pb"
	"universal/framework/common/ulog"
)

var (
	nodeList = new(sync.Map)
)

func getKey(srvType pb.ServerType, id uint32) string {
	return fmt.Sprintf("%d_%d", srvType, id)
}

// 获取节点信息
func Get(srvType pb.ServerType, id uint32) *pb.ServerNode {
	tt, ok := nodeList.Load(getKey(srvType, id))
	if !ok {
		return nil
	}
	return tt.(*pb.ServerNode)
}

// 删除节点
func Delete(serverType pb.ServerType, srvID uint32) {
	nodeList.Delete(getKey(serverType, srvID))
	ulog.Info(1, "删除服务节点: %s-%d", serverType.String(), srvID)
}

// 添加节点
func Add(node *pb.ServerNode) {
	nodeList.Store(getKey(node.ServerType, node.ServerID), node)
	buf, _ := json.Marshal(node)
	ulog.Info(1, "新增服务节点: ", string(buf))
}

// 随机路由一个节点
func Random(head *pb.PacketHead) (ret *pb.ServerNode) {
	list := []*pb.ServerNode{}
	nodeList.Range(func(_, val interface{}) bool {
		if tt := val.(*pb.ServerNode); tt.ServerType == head.DstServerType {
			list = append(list, tt)
		}
		return true
	})
	sort.Slice(list, func(i, j int) bool {
		return list[i].ServerID < list[j].ServerID
	})
	ret = list[int(head.UID)%len(list)]
	ulog.Info(1, "随机路由节点：", ret, head)
	return
}

func Print() {
	nodeList.Range(func(_, val interface{}) bool {
		list, ok := val.(**pb.ServerNode)
		if !ok || list == nil {
			return true
		}
		buf, _ := json.Marshal(list)
		ulog.Info(1, "服务节点： %s", string(buf))
		return true
	})
}
