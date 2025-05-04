package discovery

import (
	"context"
	"fmt"
	"path"
	"sync/atomic"
	"time"
	"universal/framework/define"
	"universal/library/baselib/safe"
	"universal/library/mlog"

	"github.com/spf13/cast"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Etcd struct {
	status   int32
	root     string
	parseFun define.ParseNodeFunc
	client   *clientv3.Client
	exit     chan struct{}
}

func NewEtcd(endpoints []string, opts ...OpOption) (*Etcd, error) {
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		return nil, err
	}
	vals := new(Op)
	for _, opt := range opts {
		opt(vals)
	}
	return &Etcd{
		status:   vals.status,
		root:     vals.root,
		parseFun: vals.parse,
		client:   cli,
		exit:     make(chan struct{}),
	}, nil
}

func (e *Etcd) getKey(node define.INode) string {
	return path.Join(e.root, cast.ToString(node.GetType()), cast.ToString(node.GetId()))
}

func (e *Etcd) Get() (rets []define.INode, err error) {
	rsp, err := e.client.Get(context.Background(), e.root, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, kv := range rsp.Kvs {
		rets = append(rets, e.parseFun(kv.Value))
	}
	return
}

// 添加 key-value
func (e *Etcd) Put(node define.INode) error {
	_, err := e.client.Put(context.Background(), e.getKey(node), string(node.ToBytes()))
	return err
}

// 删除 key-value
func (e *Etcd) Del(info define.INode) error {
	_, err := e.client.Delete(context.Background(), e.getKey(info), clientv3.WithPrefix())
	return err
}

// 监听 kv 变化
func (e *Etcd) Watch(cluster define.ICluster) error {
	// 获取在线服务
	rsp, err := e.client.Get(context.Background(), e.root, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range rsp.Kvs {
		node := e.parseFun(kv.Value)
		if err := cluster.Put(node); err != nil {
			mlog.Error("Etcd获取在线服务失败: %v, node: %v", err, node)
			return err
		}
	}
	// 监听服务
	listens := e.client.Watch(context.Background(), e.root, clientv3.WithPrefix())
	safe.SafeGo(mlog.Error, func() {
		for listen := range listens {
			for _, event := range listen.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					if err := cluster.Put(e.parseFun(event.Kv.Value)); err != nil {
						mlog.Error("Etcd发现新服务，新服务添加失败: %v", err)
					}
				case clientv3.EventTypeDelete:
					id := cast.ToUint32(path.Base(string(event.Kv.Key)))
					typ := cast.ToUint32(path.Base(path.Dir(string(event.Kv.Key))))
					if err := cluster.Del(typ, id); err != nil {
						mlog.Error("Etcd发现服务下线，删除服务失败: %v", err)
					}
				}
			}
		}
	})
	return nil
}

func (e *Etcd) KeepAlive(srv define.INode, ttl int64) error {
	rsp, err := e.client.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}

	// 设置租赁时间
	lease := clientv3.LeaseID(rsp.ID)
	_, err = e.client.Put(context.Background(), e.getKey(srv), string(srv.ToBytes()), clientv3.WithLease(lease))
	if err != nil {
		return err
	}

	// 多次重试
	keepAlive := func(lease clientv3.LeaseID) error {
		for i := 0; i < 3; i++ {
			if _, err := e.client.Lease.KeepAliveOnce(context.Background(), lease); err == nil {
				return nil
			}
			mlog.Error("Etcd租赁续约保活失败: %v", err)
			time.Sleep(1 * time.Second)
		}
		return fmt.Errorf("Etcd 租赁续约失败，且超过最大重试次数")
	}

	// 定时续约保活
	safe.SafeGo(mlog.Error, func() {
		atomic.AddInt32(&e.status, 1)
		tt := time.NewTicker(time.Duration(ttl/2) * time.Second)
		defer tt.Stop()

		for {
			select {
			case <-e.exit:
				return
			case <-tt.C:
				if err := keepAlive(lease); err != nil {
					mlog.Error("Etcd租赁续约保活失败: %v", err)
				}
			}
		}
	})
	return nil
}

func (e *Etcd) Close() error {
	if atomic.LoadInt32(&e.status) > 0 {
		e.exit <- struct{}{}
	}
	return e.client.Close()
}
