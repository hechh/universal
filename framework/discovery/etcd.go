package discovery

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"time"
	"universal/framework/define"
	"universal/library/baselib/safe"
	"universal/library/mlog"

	"go.etcd.io/etcd/clientv3"
)

type ParseFunc func([]byte, []byte) define.IServer

type Etcd struct {
	root   string
	parse  ParseFunc
	client *clientv3.Client
}

func NewEtcd(root string, parse ParseFunc, endpoints ...string) (*Etcd, error) {
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		return nil, err
	}
	return &Etcd{
		root:   root,
		parse:  parse,
		client: cli,
	}, nil
}

func (e *Etcd) getKey(node define.IServer) string {
	return path.Join(
		e.root,
		strconv.Itoa(int(node.GetServerType())),
		strconv.Itoa(int(node.GetServerId())),
	)
}

func (e *Etcd) Put(ctx context.Context, info define.IServer) error {
	_, err := e.client.Put(ctx, e.getKey(info), info.GetAddress())
	return err
}

func (e *Etcd) Delete(ctx context.Context, info define.IServer) error {
	_, err := e.client.Delete(ctx, e.getKey(info), clientv3.WithPrefix())
	return err
}

func (e *Etcd) Watch(ctx context.Context, add, del func(define.IServer)) error {
	// 先获取在线服务
	rsp, err := e.client.Get(ctx, e.root, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range rsp.Kvs {
		add(e.parse(kv.Key, kv.Value))
	}

	// 先监听服务
	safe.SafeGo(mlog.Error, func() {
		for listen := range e.client.Watch(ctx, e.root, clientv3.WithPrefix()) {
			for _, event := range listen.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					add(e.parse(event.Kv.Key, event.Kv.Value))
				case clientv3.EventTypeDelete:
					del(e.parse(event.Kv.Key, nil))
				}
			}
		}
	})
	return nil
}

func (e *Etcd) KeepAlive(ctx context.Context, srv define.IServer, ttl int64) error {
	rsp, err := e.client.Grant(ctx, ttl)
	if err != nil {
		return err
	}

	// 设置租赁时间
	lease := clientv3.LeaseID(rsp.ID)
	_, err = e.client.Put(ctx, e.getKey(srv), srv.GetAddress(), clientv3.WithLease(lease))
	if err != nil {
		return err
	}

	// 多次重试
	keepAlive := func(ctx context.Context, lease clientv3.LeaseID) error {
		for i := 0; i < 3; i++ {
			if _, err := e.client.Lease.KeepAliveOnce(ctx, lease); err == nil {
				return nil
			}
			mlog.Error("Etcd租赁续约保活失败: %v", err)
			time.Sleep(1 * time.Second)
		}
		return fmt.Errorf("Etcd 租赁续约失败，且超过最大重试次数")
	}

	// 定时检测
	tt := time.NewTicker(time.Duration(ttl/2) * time.Second)
	defer tt.Stop()
	safe.SafeGo(mlog.Error, func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-tt.C:
				if err := keepAlive(ctx, lease); err != nil {
					mlog.Error("Etcd租赁续约保活失败: %v", err)
					return
				}
			}
		}
	})
	return nil
}
