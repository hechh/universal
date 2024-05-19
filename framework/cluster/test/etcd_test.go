package test

import (
	"fmt"
	"sync"
	"testing"
	"time"
	"universal/common/pb"
	"universal/framework/cluster/internal/discovery/etcd"
	"universal/framework/cluster/internal/nodes"
	"universal/framework/cluster/internal/router"
	"universal/framework/cluster/internal/service"
	"universal/framework/common/fbasic"

	"google.golang.org/protobuf/proto"
)

func TestDisEtcd(t *testing.T) {
	// 初始化支持类型
	err := service.Init([]string{"localhost:2379", "172.16.126.208:33501"})
	if err != nil {
		t.Log(err)
		return
	}
	// 添加服务节点
	client := service.GetDiscovery().(*etcd.EtcdClient)
	client.Delete(service.GetRootDir())
	// 注册服务节点
	if err := service.Discovery(pb.ServerType_GATE, "127.1.0.1:10100"); err != nil {
		t.Log(err)
		return
	}
	t.Run("添加路由表", func(t *testing.T) {
		node := &pb.ServerNode{
			ServerType: pb.ServerType_GATE,
			Ip:         "127.0.0.1",
			Port:       10100,
		}
		wg := sync.WaitGroup{}
		for i := int32(1); i < 10; i++ {
			newNode := *node
			wg.Add(1)
			go func(i int32, node *pb.ServerNode) {
				node.ServerID = fbasic.GetCrc32(fmt.Sprintf("%s:%d", node.Ip, node.Port+i))
				buf, _ := proto.Marshal(node)
				client.Put(service.GetNodeChannel(pb.ServerType_GATE, node.ServerID), string(buf))
				wg.Done()
			}(i, &newNode)
		}
		wg.Wait()
	})
	t.Run("路由表print", func(t *testing.T) {
		nodes.Print()
	})
	t.Run("路由测试", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := uint64(1); i < 10; i++ {
			head := &pb.PacketHead{
				DstServerType: pb.ServerType_GATE,
				UID:           100000 + i,
			}
			wg.Add(1)
			go func(head *pb.PacketHead) {
				if err := service.Dispatcher(head); err != nil {
					t.Log(err)
				}
				t.Log("----->", head)
				wg.Done()
			}(head)
		}
		wg.Wait()
		router.Print()
	})
	t.Run("路由表删除", func(t *testing.T) {
		node := &pb.ServerNode{
			ServerType: pb.ServerType_GATE,
			Ip:         "127.0.0.1",
			Port:       10100,
		}
		for i := 0; i < 5; i++ {
			node.ServerID = fbasic.GetCrc32(fmt.Sprintf("%s:%d", node.Ip, node.Port))
			if err := client.Delete(service.GetNodeChannel(pb.ServerType_GATE, node.ServerID)); err != nil {
				t.Log("=======>", err)
				return
			}
			node.Port++
		}
	})
	t.Run("路由表print", func(t *testing.T) {
		nodes.Print()
	})

	time.Sleep(5 * time.Second)

	client.Delete(service.GetRootDir())
	service.Close()
}
