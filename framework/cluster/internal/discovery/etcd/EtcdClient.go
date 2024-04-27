package etcd

import (
	"context"
	"time"
	"universal/common/pb"
	"universal/framework/cluster/domain"
	"universal/framework/fbasic"

	"go.etcd.io/etcd/clientv3"
)

type EtcdClient struct {
	client   *clientv3.Client
	lease    clientv3.Lease
	key      *keyMonitor
	notifyCh chan *keyMonitor
}

func NewEtcdClient(ends ...string) (*EtcdClient, error) {
	client, err := clientv3.New(clientv3.Config{Endpoints: ends})
	if err != nil {
		return nil, fbasic.NewUError(1, pb.ErrorCode_EtcdBuildClient, err)
	}
	return &EtcdClient{
		client:   client,
		notifyCh: make(chan *keyMonitor, 2),
	}, nil
}

func (d *EtcdClient) Walk(path string, f domain.WatchFunc) error {
	resp, err := d.client.Get(context.Background(), path, clientv3.WithPrefix())
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_EtcdClientGet, err)
	}
	for _, kv := range resp.Kvs {
		f(string(kv.Key), kv.Value)
	}
	return nil
}

// 阻塞执行
func (d *EtcdClient) Watch(path string, addF, delF domain.WatchFunc) {
	go d.run(d.client.Watch(context.Background(), path, clientv3.WithPrefix()), addF, delF)
}

func (d *EtcdClient) KeepAlive(key string, value []byte, ttl int64) {
	d.notifyCh <- &keyMonitor{
		ttl:   ttl,
		key:   key,
		value: value,
	}
}

func (d *EtcdClient) run(watch clientv3.WatchChan, addF, delF domain.WatchFunc) {
	timer := time.NewTicker(4 * time.Second)
	for {
		select {
		case <-timer.C:
			if d.key == nil {
				continue
			}
			if err := d.key.KeepAliveOnce(d.client); err != nil {
				//d.notifyCh <- d.key
				//log.Println(err)
				panic(err)
			}
		case d.key = <-d.notifyCh:
			if err := d.key.Put(d.client, d.lease); err != nil {
				//d.notifyCh <- d.key
				panic(err)
			}
		case item := <-watch:
			for _, event := range item.Events {
				if event.Type.String() == "PUT" && addF != nil {
					addF(string(event.Kv.Key), event.Kv.Value)
					continue
				}
				if event.Type.String() == "DELETE" && delF != nil {
					delF(string(event.Kv.Key), event.Kv.Value)
				}
			}
		}
	}
}

type keyMonitor struct {
	ttl     int64
	key     string
	value   []byte
	leaseID clientv3.LeaseID
}

func (d *keyMonitor) Put(client *clientv3.Client, lease clientv3.Lease) error {
	resp, err := lease.Create(context.Background(), d.ttl)
	if err == nil {
		return fbasic.NewUError(1, pb.ErrorCode_EtcdLeaseCreate, err)
	}
	d.leaseID = clientv3.LeaseID(resp.ID)

	// 设置key的ttl
	if _, err = client.Put(context.TODO(), d.key, string(d.value), clientv3.WithLease(d.leaseID)); err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_EtcdClientPut, err)
	}
	return nil
}

func (d *keyMonitor) KeepAliveOnce(client *clientv3.Client) error {
	_, err := client.KeepAliveOnce(context.Background(), d.leaseID)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_EtcdLeaseKeepAliveOnce, err)
	}
	return nil
}
