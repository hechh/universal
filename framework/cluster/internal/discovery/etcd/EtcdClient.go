package etcd

import (
	"context"
	"log"
	"time"
	"universal/common/pb"
	"universal/framework/cluster/domain"
	"universal/framework/fbasic"

	"go.etcd.io/etcd/clientv3"
)

type EtcdClient struct {
	client   *clientv3.Client
	notifyCh chan *keyMonitor
	exitCh   chan struct{}
}

func NewEtcdClient(ends ...string) (*EtcdClient, error) {
	client, err := clientv3.New(clientv3.Config{Endpoints: ends})
	if err != nil {
		return nil, fbasic.NewUError(1, pb.ErrorCode_EtcdBuildClient, err)
	}
	return &EtcdClient{
		client:   client,
		notifyCh: make(chan *keyMonitor, 1),
		exitCh:   make(chan struct{}, 0),
	}, nil
}

func (d *EtcdClient) KeepAlive(key string, value []byte, ttl int64) {
	d.notifyCh <- &keyMonitor{
		ttl:   ttl,
		key:   key,
		value: value,
	}
}

func (d *EtcdClient) Put(key string, value string) error {
	if _, err := d.client.Put(context.Background(), key, value); err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_EtcdClientPut, err)
	}
	return nil
}

func (d *EtcdClient) Close() {
	d.exitCh <- struct{}{}
}

func (d *EtcdClient) Watch(path string, watchFunc domain.WatchFunc) error {
	resp, err := d.client.Get(context.Background(), path, clientv3.WithPrefix())
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_EtcdClientGet, err)
	}
	for _, kv := range resp.Kvs {
		watchFunc(domain.ActionTypeNone, string(kv.Key), fbasic.BytesToString(kv.Value))
	}
	// 设置监听
	watchCh := d.client.Watch(context.Background(), path, clientv3.WithPrefix())
	go d.run(watchCh, watchFunc)
	return nil
}

func (d *EtcdClient) run(watchCh clientv3.WatchChan, watchFunc domain.WatchFunc) {
	times := 0
	var key *keyMonitor
	timer := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-d.exitCh:
			return
		case <-timer.C:
			if key == nil {
				continue
			}
			if err := key.KeepAliveOnce(d.client); err != nil {
				log.Println(err)
				d.notifyCh <- key
				key = nil
			}
		case key = <-d.notifyCh:
			if err := key.Put(d.client); err != nil {
				log.Println(err)
				if times >= 5 {
					panic(err)
				}
				times++
				d.notifyCh <- key
				key = nil
			}
		case item := <-watchCh:
			for _, event := range item.Events {
				action := domain.ActionTypeAdd
				if event.Type.String() != "PUT" {
					action = domain.ActionTypeDel
				}
				watchFunc(action, string(event.Kv.Key), fbasic.BytesToString(event.Kv.Value))
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

func (d *keyMonitor) Put(client *clientv3.Client) error {
	resp, err := client.Lease.Grant(context.Background(), d.ttl)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_EtcdLeaseCreate, err)
	}
	d.leaseID = clientv3.LeaseID(resp.ID)

	// 设置key的ttl
	if _, err = client.Put(context.Background(), d.key, string(d.value), clientv3.WithLease(d.leaseID)); err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_EtcdClientPut, err)
	}
	return nil
}

func (d *keyMonitor) KeepAliveOnce(client *clientv3.Client) error {
	_, err := client.Lease.KeepAliveOnce(context.Background(), d.leaseID)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_EtcdLeaseKeepAliveOnce, err)
	}
	return nil
}
