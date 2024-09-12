package discovery

import (
	"context"
	"time"
	"universal/framework/basic"
	"universal/framework/plog"

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
	return &EtcdClient{
		client: cli,
		exit:   make(chan struct{}, 1),
	}, nil
}

func (d *EtcdClient) Close() {
	d.exit <- struct{}{}
}

func (d *EtcdClient) Put(key string, value string) error {
	_, err := d.client.Put(context.Background(), key, value)
	return err
}

func (d *EtcdClient) Delete(key string) error {
	_, err := d.client.Delete(context.Background(), key, clientv3.WithPrefix())
	return err
}

// 遍历指定path目录下面的所有key-value
func (d *EtcdClient) Walk(path string, f func(string, string)) error {
	rsp, err := d.client.Get(context.Background(), path, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range rsp.Kvs {
		f(string(kv.Key), basic.BytesToString(kv.Value))
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
	rsp := d.client.Watch(context.Background(), path, clientv3.WithPrefix())
	basic.SafeGo(plog.Catch, func() {
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
					del(string(event.Kv.Key), basic.BytesToString(event.Kv.Value))
				} else {
					add(string(event.Kv.Key), basic.BytesToString(event.Kv.Value))
				}
			}
		}
	})
	return nil
}

func (d *EtcdClient) KeepAlive(key, value string, ttl int64) error {
	// 设置租赁
	rsp, err := d.client.Lease.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}
	// 设置key-value
	lease := clientv3.LeaseID(rsp.ID)
	if _, err := d.client.Put(context.Background(), key, value, clientv3.WithLease(lease)); err != nil {
		return err
	}
	// 定时器执行
	tt := time.NewTicker(time.Duration(ttl/2) * time.Second)
	basic.SafeGo(plog.Catch, func() {
		for {
			select {
			case <-d.exit:
				select {
				case <-d.exit:
				default:
				}
				return
			case <-tt.C:
				_, err := d.client.Lease.KeepAliveOnce(context.Background(), lease)
				if err != nil {
					plog.Error("Etcd租赁续约保活失败: %v", err)
					return
				}
			}
		}
	})
	return nil
}
