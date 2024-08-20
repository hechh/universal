package discovery

import (
	"context"
	"universal/framework/basic/util"

	"go.etcd.io/etcd/clientv3"
)

type EtcdClient struct {
	client *clientv3.Client
	exit   chan struct{}
}

func NewEtcdClient(ends ...string) (*EtcdClient, error) {
	cli, err := clientv3.New(clientv3.Config{Endpoints: ends})
	if err != nil {
		return nil, err
	}
	return &EtcdClient{client: cli, exit: make(chan struct{}, 1)}, nil
}

// 遍历指定path目录下面的所有key-value
func (d *EtcdClient) Walk(path string, f func(string, string)) error {
	rsp, err := d.client.Get(context.Background(), path, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range rsp.Kvs {
		f(string(kv.Key), util.BytesToString(kv.Value))
	}
	return nil
}

// 监控指定目录下面的所有key-value
func (d *EtcdClient) Watch(path string, add, del func(string, string)) error {
	// 先遍历所有key-value
	if err := d.Walk(path, add); err != nil {
		return err
	}
	// 监听最新变更
	monitorF := func(rsp <-chan clientv3.WatchResponse) {
		select {
		case <-d.exit:
			select {
			case <-d.exit:
			default:
			}
			return
		case item := <-rsp:
			for _, event := range item.Events {
				if event.Type == clientv3.EventTypeDelete {
					del(string(event.Kv.Key), util.BytesToString(event.Kv.Value))
				} else {
					add(string(event.Kv.Key), util.BytesToString(event.Kv.Value))
				}
			}
		}
	}
	util.SafeGo(nil, func() {
		monitorF(d.client.Watch(context.Background(), path, clientv3.WithPrefix()))
	})
	return nil
}

/*
func (d *EtcdClient) Put(key, val string, ttl int64) error {
	if ttl <= 0 {
		_, err := d.client.Put(context.Background(), key, val)
		return err
	}

	// 设置过期时间
	if rsp, err := d.client.Lease.Grant(context.Background(), ttl); err != nil {
		return err
	} else {
		d.lease = clientv3.LeaseID(rsp.ID)
	}
	// 设置key-value
	_, err := d.client.Put(context.Background(), key, val, clientv3.WithLease(d.lease))
	return err
}

type Monitor struct {
	client *clientv3.Client
	lease  clientv3.LeaseID
	key    string
	value  string
	ttl    int64
}

func NewMonitor(cli *clientv3.Client, key, value string, ttl int64) *Monitor {

	return nil
}

func (d *Monitor) Put(cli *clientv3.Client) error {
	// 设置过期时间
	if rsp, err := cli.Lease.Grant(context.Background(), d.ttl); err != nil {
		return err
	} else {
		d.lease = clientv3.LeaseID(rsp.ID)
	}
	// 设置key-value
	_, err := cli.Put(context.Background(), d.key, d.value, clientv3.WithLease(d.lease))
	return err
}

func (d *Monitor) KeepAliveOnce(cli *clientv3.Client) error {
	_, err := cli.Lease.KeepAliveOnce(context.Background(), d.lease)
	return err
}
*/
