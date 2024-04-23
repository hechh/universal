package etcd

import (
	"context"
	"fmt"
	"log"
	"time"
	"universal/common/pb"
	"universal/framework/basic"
	"universal/framework/cluster/domain"

	"go.etcd.io/etcd/clientv3"
	"google.golang.org/protobuf/proto"
)

type Etcd struct {
	client *clientv3.Client
	self   *ClusterNode
}

type ClusterNode struct {
	node     *pb.ClusterNode
	lease    clientv3.Lease
	leaseID  clientv3.LeaseID
	notifyCh chan struct{}
}

func (d *ClusterNode) KeepAliveOnce(client *clientv3.Client) {
	if _, err := client.KeepAliveOnce(context.Background(), d.leaseID); err != nil {
		d.notifyCh <- struct{}{}
		log.Print(err)
	}
}

func (d *ClusterNode) Put(client *clientv3.Client) {
	leaseResp, err := d.lease.Create(context.Background(), 10)
	if err != nil {
		log.Print(err)
		d.notifyCh <- struct{}{}
		return
	}
	d.leaseID = clientv3.LeaseID(leaseResp.ID)
	// 重新注册服务
	key := domain.GetNodeChannel(d.node.ClusterType, d.node.ClusterID)
	buf, _ := proto.Marshal(d.node)
	if _, err = client.Put(context.TODO(), key, string(buf), clientv3.WithLease(d.leaseID)); err != nil {
		log.Print(err)
		d.notifyCh <- struct{}{}
	}
}
func NewEtcd(node *pb.ClusterNode, endpoints ...string) (*Etcd, error) {
	client, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		return nil, err
	}
	node.ClusterID = basic.GetCrc32(fmt.Sprintf("%s:%d", node.Ip, node.Port))
	return &Etcd{
		client: client,
		self: &ClusterNode{
			node:     node,
			lease:    clientv3.NewLease(client),
			notifyCh: make(chan struct{}, 1),
		},
	}, nil
}

func (d *Etcd) GetSelf() *pb.ClusterNode {
	return d.self.node
}

func (d *Etcd) Walk(path string, f domain.DiscoveryFunc) error {
	resp, err := d.client.Get(context.Background(), path, clientv3.WithPrefix())
	if err != nil {
		return basic.NewUError(2, pb.ErrorCode_BuildEtcdClient, err)
	}
	for _, kv := range resp.Kvs {
		if node, err := basic.UnmarhsalClusterNode(kv.Value); err == nil {
			f(domain.ActionTypeNone, node)
		} else {
			return basic.NewUError(2, pb.ErrorCode_Unmarshal, err)
		}
	}
	// 自身节点
	f(domain.ActionTypeNone, d.self.node)
	return nil
}

func (d *Etcd) Watch(path string, f domain.DiscoveryFunc) {
	d.self.notifyCh <- struct{}{}
	timer := time.NewTicker(3 * time.Second)
	result := d.client.Watch(context.Background(), path, clientv3.WithPrefix())

	go func() {
		select {
		case <-timer.C:
			d.self.KeepAliveOnce(d.client)
		case <-d.self.notifyCh:
			d.self.Put(d.client)
		case item := <-result:
			for _, event := range item.Events {
				action := domain.ActionTypeDel
				if event.Type.String() == "PUT" {
					action = domain.ActionTypeAdd
				}

				if node, err := basic.UnmarhsalClusterNode(event.Kv.Value); err == nil {
					f(action, node)
				} else {
					log.Print(err)
				}
			}
		}
	}()
}
