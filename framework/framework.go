package framework

import (
	"log"
	"universal/common/pb"
	"universal/framework/actor"
	"universal/framework/cluster"
	"universal/framework/fbasic"
	"universal/framework/notify"
	"universal/framework/packet"
)

var (
	serverId    int
	clusterType pb.ClusterType
)

func SetGlobal(id int, t pb.ClusterType) {
	serverId = id
	clusterType = t
}

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
		rsp, err := packet.Call(ctx, buf)
		if err != nil {
			log.Fatalln(err)
			return
		}
		// 设置返回信息
		head := ctx.PacketHead
		head.SeqID++
		head.Status = pb.StatusType_RESPONSE
		head.SrcClusterType, head.DstClusterType = head.DstClusterType, head.SrcClusterType
		head.SrcClusterID, head.DstClusterID = head.DstClusterID, head.SrcClusterID
		// 发送
		pac, err := fbasic.RspToPacket(head, rsp)
		if err != nil {
			log.Fatalln(err)
			return
		}
		if err := notify.Publish(pac); err != nil {
			log.Fatalln(err)
		}
	}
}

func init() {
	// 设置actor处理
	actor.SetActorHandle(actorHandle)
}
