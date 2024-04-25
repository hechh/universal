package service

import (
	"log"
	"universal/common/pb"
	"universal/framework/basic"
	"universal/framework/cluster/domain"
	"universal/framework/cluster/internal/balance"
	"universal/framework/cluster/internal/discovery/etcd"
	"universal/framework/cluster/internal/routine"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const (
	ROOT_DIR = "server"
)

var (
	cluster *Cluster
)

type Cluster struct {
	dis  domain.IDiscovery
	bl   *balance.Balance
	rt   *routine.Routine
	conn *nats.Conn
}

func GetCluster() *Cluster {
	return cluster
}

func InitCluster(node *pb.ClusterNode, natsUrl string, etcds []string) error {
	dis, err := etcd.NewEtcd(node, etcds...)
	if err != nil {
		return err
	}
	bl := balance.NewBalance(pb.ClusterType_GAME, pb.ClusterType_GATE, pb.ClusterType_DB)
	// 连接nats
	conn, err := nats.Connect(natsUrl)
	if err != nil {
		return basic.NewUError(1, pb.ErrorCode_NewClient, err)
	}
	// 监听
	dis.Watch(ROOT_DIR, func(action int, item *pb.ClusterNode) {
		switch action {
		case domain.ActionTypeAdd:
			bl.AddNode(item)
		case domain.ActionTypeDel:
			bl.DelNode(item)
		}
	})
	// 发现服务
	dis.Walk(ROOT_DIR, func(action int, item *pb.ClusterNode) {
		bl.AddNode(item)
	})
	// 订阅服务
	cluster = &Cluster{conn: conn, dis: dis, bl: bl, rt: routine.NewRoutine()}
	return nil
}

func (d *Cluster) GetDiscovery() domain.IDiscovery {
	return d.dis
}

func (d *Cluster) GetBalance() *balance.Balance {
	return d.bl
}

func (d *Cluster) Subscribe(h domain.ClusterFunc) {
	node := d.dis.GetSelf()
	// 订阅
	d.conn.Subscribe(domain.GetNodeChannel(node.ClusterType, node.ClusterID), func(msg *nats.Msg) {
		pac := &pb.Packet{}
		if err := proto.Unmarshal(msg.Data, pac); err != nil {
			log.Printf("Subscribe error: %v", err)
		} else {
			h(pac)
		}
	})
	// 广播
	d.conn.Subscribe(domain.GetTopicChannel(node.ClusterType), func(msg *nats.Msg) {
		pac := &pb.Packet{}
		if err := proto.Unmarshal(msg.Data, pac); err != nil {
			log.Printf("Subscribe error: %v", err)
		} else {
			h(pac)
		}
	})
}

// 跨服务转发
func (d *Cluster) Send(pac *pb.Packet) (err error) {
	buf, err := proto.Marshal(pac)
	if err != nil {
		return basic.NewUError(1, pb.ErrorCode_Marhsal, err)
	}
	// 转发
	head := pac.Head
	switch head.SendType {
	case pb.SendType_POINT:
		// 先路由
		if err = d.dispatcher(head); err != nil {
			return err
		}
		// 转发
		err = d.conn.Publish(domain.GetNodeChannel(head.DstClusterType, head.DstClusterID), buf)
	default:
		// 转发
		err = d.conn.Publish(domain.GetTopicChannel(head.DstClusterType), buf)
	}
	return basic.NewUError(1, pb.ErrorCode_NatsPublish, err)
}

// 路由
func (d *Cluster) dispatcher(head *pb.PacketHead) error {
	rlist := d.rt.GetRoutine(head)
	// 判断玩家是否已经路由
	if rinfo := rlist.Get(head.DstClusterType); rinfo == nil {
		// 路由
		if err := rlist.UpdateRoutine(head, d.bl.RandomNode(head)); err != nil {
			return err
		}
	} else {
		if head.DstClusterID <= 0 {
			head.DstClusterID = rinfo.GetClusterID()
		}
		// 节点丢失
		if head.DstClusterID != rinfo.GetClusterID() {
			// 重新路由
			if err := rlist.UpdateRoutine(head, d.bl.RandomNode(head)); err != nil {
				return err
			}
		}
		// 判断节点是否存在
		if node := d.bl.GetNode(head); node == nil {
			// 重新路由
			if err := rlist.UpdateRoutine(head, d.bl.RandomNode(head)); err != nil {
				return err
			}
		} else {
			rlist.UpdateRoutine(head, node)
		}
	}
	return nil
}
