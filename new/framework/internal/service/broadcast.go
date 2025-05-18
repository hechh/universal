package service

import (
	"poker_server/common/pb"
	"poker_server/framework/library/uerror"

	"github.com/golang/protobuf/proto"
)

func BroadcastMsgToNode(head *pb.Head, msg proto.Message) error {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return uerror.New(1, -1, "序列化消息失败: %v", err)
	}
	return BroadcastToNode(head, buf)
}

func BroadcastToNode(head *pb.Head, msg []byte) error {
	// 判断参数是否正确
	if head.DstNodeType >= pb.NodeType_End || head.DstNodeType <= pb.NodeType_Begin {
		return uerror.New(1, -1, "服务类型不支持: %v", head)
	}
	if clusterObj.GetCount(head.DstNodeType) <= 0 {
		return uerror.New(1, -1, "未找到服务节点: %v", head)
	}
	// 设置值
	head.SendType = pb.SendType_BROADCAST
	head.SrcNodeId = self.Id
	head.SrcNodeType = self.Type
	// 更新路由
	if head.Id > 0 {
		router, err := tableObj.Get(head.IdType, head.Id)
		if err != nil {
			return err
		}
		router.Set(self.Type, self.Id)
	}
	// 发送
	return busObj.Broadcast(head, msg)
}

func BroadcastMsgToClient(head *pb.Head, msg proto.Message) error {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return uerror.New(1, -1, "序列化消息失败: %v", err)
	}
	return BroadcastToClient(head, buf)
}

func BroadcastToClient(head *pb.Head, msg []byte) error {
	head.SendType = pb.SendType_BROADCAST
	head.SrcNodeType = self.Type
	head.SrcNodeId = self.Id
	head.DstNodeType = pb.NodeType_Gate
	head.ActorName = ""
	head.FuncName = ""
	return busObj.Broadcast(head, msg)
}
