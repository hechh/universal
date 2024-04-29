package test

import (
	"context"
	"log"
	"testing"
	"time"
	"universal/common/pb"
	"universal/framework/cluster/domain"
	"universal/framework/cluster/internal/discovery/etcd"
	"universal/framework/cluster/internal/nodes"
	"universal/framework/cluster/internal/service"
	"universal/framework/fbasic"

	"go.etcd.io/etcd/clientv3"
	"google.golang.org/protobuf/proto"
)

func TestEtcd(t *testing.T) {
	// 创建Etcd客户端
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"172.16.126.208:33501"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 创建一个新的租约
	leaseResp, err := client.Lease.Grant(context.Background(), 10) // 10秒后过期的租约
	if err != nil {
		log.Fatal(err)
	}
	leaseID := clientv3.LeaseID(leaseResp.ID)

	// 将键值对与租约相关联
	_, err = client.Put(context.Background(), "key1", "value1", clientv3.WithLease(leaseID))
	if err != nil {
		log.Fatal(err)
	}

	// 保持租约活跃
	keepAliveChan, err := client.Lease.KeepAlive(context.Background(), leaseID)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			resp := <-keepAliveChan
			if resp == nil {
				log.Println("租约已失效")
				break
			}
			log.Printf("收到租约续约响应: TTL = %d\n", resp.TTL)
		}
	}()

	// 模拟程序运行一段时间后撤销租约
	time.Sleep(20 * time.Second)
	_, err = client.Lease.Revoke(context.Background(), leaseID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("租约已撤销")
}

func TestDisEtcd(t *testing.T) {
	client, err := etcd.NewEtcdClient("172.16.126.208:33501")
	if err != nil {
		t.Log(err)
		return
	}
	node := &pb.ClusterNode{
		ClusterType: pb.ClusterType_GATE,
		Ip:          "127.0.0.1",
		Port:        10100,
	}
	buf, _ := proto.Marshal(node)
	client.KeepAlive(domain.GetNodeChannel(pb.ClusterType_GATE, 123), buf, 10)
	client.Watch(domain.ROOT_DIR, func(action int, key string, value string) {
		vv := &pb.ClusterNode{}
		proto.Unmarshal(fbasic.StringToBytes(value), vv)
		switch action {
		case domain.ActionTypeDel:
			nodes.DeleteNode(vv)
		default:
			nodes.AddNode(vv)
		}
		t.Log(action, "kv==add>", key, value)
	})
	for i := 0; i < 10; i++ {
		node.Port++
		buf, _ := proto.Marshal(node)
		client.Put(domain.GetNodeChannel(pb.ClusterType_GATE, uint32(i)+1), string(buf))
	}

	nodes.InitNodes(pb.ClusterType_GATE, pb.ClusterType_GAME, pb.ClusterType_DB)
	t.Run("路由测试", func(t *testing.T) {
		for i := uint64(1); i < 10; i++ {
			head := &pb.PacketHead{DstClusterType: pb.ClusterType_GATE, UID: 100000 + i}
			if err := service.Dispatcher(head); err != nil {
				t.Log("=========>", err)
				return
			}
			t.Log("----->", head)
		}
	})

	client.Close()
}
