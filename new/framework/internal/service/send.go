package service

import (
	"poker_server/common/pb"
	"poker_server/framework/library/mlog"
	"poker_server/framework/library/uerror"

	"github.com/golang/protobuf/proto"
)

func SendMsgToNode(head *pb.Head, msg proto.Message) error {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return uerror.New(1, -1, "序列化消息失败: %v", err)
	}
	return SendToNode(head, buf)
}

func SendToNode(head *pb.Head, msg []byte) error {
	// 判断参数是否正确
	if head.Id <= 0 {
		return uerror.New(1, -1, "唯一ID为空: %v", head)
	}
	if head.DstNodeType >= pb.NodeType_End || head.DstNodeType <= pb.NodeType_Begin {
		return uerror.New(1, -1, "服务类型不支持: %v", head)
	}
	if clusterObj.GetCount(head.DstNodeType) <= 0 {
		return uerror.New(1, -1, "未找到服务节点: %v", head)
	}
	// 分发
	if err := dispatch(head); err != nil {
		return err
	}
	return busObj.Send(head, msg)
}

func RequestMsgToNode(head *pb.Head, msg proto.Message, reply proto.Message) error {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return uerror.New(1, -1, "序列化消息失败: %v", err)
	}
	return RequestToNode(head, buf, reply)
}

func RequestToNode(head *pb.Head, buf []byte, reply proto.Message) error {
	// 判断参数是否正确
	if head.Id <= 0 {
		return uerror.New(1, -1, "唯一ID为空: %v", head)
	}
	if head.DstNodeType >= pb.NodeType_End || head.DstNodeType <= pb.NodeType_Begin {
		return uerror.New(1, -1, "服务类型不支持: %v", head)
	}
	if clusterObj.GetCount(head.DstNodeType) <= 0 {
		return uerror.New(1, -1, "未找到服务节点: %v", head)
	}
	// 分发
	if err := dispatch(head); err != nil {
		return err
	}
	return busObj.Request(head, buf, reply)
}

// 路由分发
func dispatch(head *pb.Head) error {
	// 设置值
	head.SendType = pb.SendType_POINT
	head.SrcNodeId = self.Id
	head.SrcNodeType = self.Type

	// 更新路由
	router, err := tableObj.Get(head.IdType, head.Id)
	if err != nil {
		return err
	}
	router.Set(self.Type, self.Id)

	// 业务层直接指定具体节点
	if head.DstNodeId > 0 {
		if clusterObj.Get(head.DstNodeType, head.DstNodeId) != nil {
			router.Set(head.DstNodeType, head.DstNodeId)
			return nil
		}
		return uerror.New(1, -1, "未找到服务节点: head:%v", head)
	}

	// 优先从路由中选择
	if nodeId := router.Get(head.DstNodeType); nodeId > 0 {
		if clusterObj.Get(head.DstNodeType, nodeId) != nil {
			head.DstNodeId = nodeId
			return nil
		}
	}

	//从集群中随机获取一个节点
	if node := clusterObj.Random(head.DstNodeType, head.RouteId); node != nil {
		head.DstNodeId = node.Id
		router.Set(head.DstNodeType, node.Id)
		return nil
	}
	return uerror.New(1, -1, "未找到服务节点: head:%v", head)
}

// 将数据发送到客户端
func SendMsgToClient(head *pb.Head, msg proto.Message) error {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return uerror.New(1, -1, "序列化消息失败: %v", err)
	}
	return SendToClient(head, buf)
}

func SendToClient(head *pb.Head, msg []byte) error {
	if head.Id <= 0 {
		return uerror.New(1, -1, "唯一ID为空: %v", head)
	}
	// 设置值
	head.SendType = pb.SendType_POINT
	head.SrcNodeType = self.Type
	head.SrcNodeId = self.Id
	head.DstNodeType = pb.NodeType_Gate
	head.ActorName = ""
	head.FuncName = ""
	// 获取路由信息
	router, err := tableObj.Get(head.IdType, head.Id)
	if err != nil {
		return err
	}
	head.DstNodeId = router.Get(head.DstNodeType)
	// 判断节点是否存在
	if clusterObj.Get(head.DstNodeType, head.DstNodeId) == nil {
		return uerror.New(1, -1, "网关节点不存在: %v", head)
	}
	return busObj.Send(head, msg)
}

func NotifyMsgToClient(uids []uint64, head *pb.Head, msg proto.Message) error {
	if len(uids) <= 0 {
		return uerror.New(1, -1, "通知人员为空: %v", head)
	}
	buf, err := proto.Marshal(msg)
	if err != nil {
		return uerror.New(1, -1, "序列化消息失败: %v", err)
	}
	return NotifyToClient(uids, head, buf)
}

func NotifyToClient(uids []uint64, head *pb.Head, msg []byte) error {
	// 设置值
	head.SendType = pb.SendType_POINT
	head.SrcNodeType = self.Type
	head.SrcNodeId = self.Id
	head.DstNodeType = pb.NodeType_Gate
	head.ActorName = ""
	head.FuncName = ""
	for _, uid := range uids {
		head.Id = uid
		// 获取路由信息
		router, err := tableObj.Get(head.IdType, head.Id)
		if err != nil {
			mlog.Errorf("玩家路由不存在:%v", head)
			continue
		}
		head.DstNodeId = router.Get(head.DstNodeType)
		// 判断节点是否存在
		if clusterObj.Get(head.DstNodeType, head.DstNodeId) == nil {
			mlog.Errorf("网关节点不存在: %v", head)
			continue
		}
		if err := busObj.Send(head, msg); err != nil {
			mlog.Errorf("通知玩家失败：%v, error:%v", head, err)
		}
	}
	return nil
}
