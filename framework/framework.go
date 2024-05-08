package framework

import (
	"log"
	"universal/common/pb"
	"universal/framework/actor"
	"universal/framework/cluster"
	"universal/framework/fbasic"
	"universal/framework/notify"
	"universal/framework/packet"

	"google.golang.org/protobuf/proto"
)

// 初始化
func Init(addr string, etcds []string, natsUrl string) error {
	// 初始化集群
	typs := []pb.ClusterType{}
	for i := pb.ClusterType_NONE + 1; i < pb.ClusterType_MAX; i++ {
		typs = append(typs, i)
	}
	// 连接etcd
	if err := cluster.Init(etcds, typs...); err != nil {
		return err
	}
	// 进行服务发现
	if err := cluster.Discovery(clusterType, addr); err != nil {
		return err
	}
	// 初始化消息中间件
	if err := notify.Init(natsUrl); err != nil {
		return err
	}
	return nil
}

// 设置actor处理
func actorHandle(ctx *fbasic.Context, buf []byte) func() {
	return func() {
		// 调用接口
		rsp, err := packet.Call(ctx, buf)
		if err != nil {
			log.Fatalln(err)
			return
		}
		// 设置返回信息
		head := ctx.PacketHead
		head.SeqID++
		head.ApiCode++
		head.SrcClusterType, head.DstClusterType = head.DstClusterType, head.SrcClusterType
		head.SrcClusterID, head.DstClusterID = head.DstClusterID, head.SrcClusterID
		// 发送
		if err := PublishRspPacket(head, rsp); err != nil {
			log.Fatalln(err)
			return
		}
	}
}

func PublishRspPacket(head *pb.PacketHead, rsp proto.Message) error {
	// 封装发送包
	pac, err := fbasic.RspToPacket(head, rsp)
	if err != nil {
		return err
	}
	// 获取订阅key
	key, err := fbasic.GetHeadChannel(head)
	if err != nil {
		return err
	}
	// 发送
	return notify.Publish(key, pac)
}

func PublishReqPacket(head *pb.PacketHead, req proto.Message) error {
	// 封装发送包
	pac, err := fbasic.ReqToPacket(head, req)
	if err != nil {
		return err
	}
	// 获取订阅key
	key, err := fbasic.GetHeadChannel(head)
	if err != nil {
		return err
	}
	// 发送
	return notify.Publish(key, pac)
}

func init() {
	// 设置actor处理
	actor.SetActorHandle(actorHandle)
}
