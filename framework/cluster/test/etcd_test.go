package test

import (
	"fmt"
	"testing"
	"universal/common/pb"
	"universal/framework"
	"universal/framework/cluster/internal/discovery/etcd"
	"universal/framework/cluster/internal/nodes"
	"universal/framework/cluster/internal/service"
	"universal/framework/fbasic"
	"universal/framework/routine"

	"google.golang.org/protobuf/proto"
)

func TestDisEtcd(t *testing.T) {
	// 初始化支持类型
	typs := []pb.ClusterType{pb.ClusterType_GATE, pb.ClusterType_GAME}
	err := service.Init([]string{"localhost:2379", "172.16.126.208:33501"}, typs...)
	if err != nil {
		t.Log(err)
		return
	}
	// 添加服务节点
	client := service.GetDiscovery().(*etcd.EtcdClient)
	client.Delete("server/")
	// 注册服务节点
	if err := service.Discovery(pb.ClusterType_GATE, "127.1.0.1:10100"); err != nil {
		t.Log(err)
		return
	}
	t.Run("添加路由表", func(t *testing.T) {
		node := &pb.ClusterNode{
			ClusterType: pb.ClusterType_GATE,
			Ip:          "127.0.0.1",
			Port:        10100,
		}
		for i := 1; i < 10; i++ {
			node.Port++
			node.ClusterID = fbasic.GetCrc32(fmt.Sprintf("%s:%d", node.Ip, node.Port))
			buf, _ := proto.Marshal(node)
			client.Put(fbasic.GetNodeChannel(pb.ClusterType_GATE, node.ClusterID), string(buf))
		}
	})
	t.Run("路由表print", func(t *testing.T) {
		nodes.Print()
	})
	t.Run("路由测试", func(t *testing.T) {
		for i := uint64(1); i < 10; i++ {
			head := &pb.PacketHead{DstClusterType: pb.ClusterType_GATE, UID: 100000 + i}
			if err := framework.Dispatcher(head); err != nil {
				t.Log(err)
				return
			}
			t.Log("----->", head)
		}
		routine.Print()
	})
	t.Run("路由表删除", func(t *testing.T) {
		node := &pb.ClusterNode{
			ClusterType: pb.ClusterType_GATE,
			Ip:          "127.0.0.1",
			Port:        10100,
		}
		for i := 0; i < 5; i++ {
			node.ClusterID = fbasic.GetCrc32(fmt.Sprintf("%s:%d", node.Ip, node.Port))
			if err := client.Delete(fbasic.GetNodeChannel(pb.ClusterType_GATE, node.ClusterID)); err != nil {
				t.Log("=======>", err)
				return
			}
			node.Port++
		}
	})
	t.Run("路由表print", func(t *testing.T) {
		nodes.Print()
	})
	client.Delete("server/")
	client.Delete("/")
	service.Stop()
}
