package discovery

import (
	"context"
	"fmt"
	"path"
	"sync"
	"time"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/domain"
	"universal/library/mlog"
	"universal/library/safe"
	"universal/library/util"

	"github.com/spf13/cast"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
)

type Etcd struct {
	sync.WaitGroup
	client *clientv3.Client
	lease  clientv3.LeaseID
	exit   chan struct{}
	topic  string
}

func NewEtcd(cfg *yaml.EtcdConfig) (cli *Etcd, err error) {
	v3cfg := clientv3.Config{
		Endpoints:            cfg.Endpoints,
		DialTimeout:          5 * time.Second,
		DialKeepAliveTime:    30 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
		MaxCallSendMsgSize:   10 * 1024 * 1024,
	}
	err = util.Retry(3, time.Second, func() error {
		client, err := clientv3.New(v3cfg)
		if err == nil {
			cli = &Etcd{
				client: client,
				topic:  cfg.Topic,
				exit:   make(chan struct{}),
			}
		}
		return err
	})
	return cli, err
}

func (d *Etcd) Watch(cls domain.ICluster) (err error) {
	ctx := context.Background()
	tctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rsp, err := d.client.Get(tctx, d.topic, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range rsp.Kvs {
		node := &pb.Node{}
		if err := proto.Unmarshal(kv.Value, node); err != nil {
			return err
		}
		if cls.Add(node) {
			mlog.Infof("添加服务节点：%s", node.String())
		}
	}

	safe.Go(func() {
		for {
			wchan := d.client.Watch(ctx, d.topic, clientv3.WithPrefix())
			if wchan == nil {
				mlog.Infof("Etcd监听%s失败", d.topic)
				time.Sleep(time.Second)
				continue
			}
			for rsp := range wchan {
				if rsp.Canceled {
					mlog.Errorf("Etcd 监听被取消，尝试重新连接")
					break
				}
				if rsp.Err() != nil {
					mlog.Errorf("Etcd 监听出现错误: %v", rsp.Err().Error())
					continue
				}
				d.handleEvent(cls, rsp.Events...)
			}
		}
	})
	return nil
}

func (d *Etcd) handleEvent(cls domain.ICluster, events ...*clientv3.Event) {
	for _, event := range events {
		switch event.Type {
		case clientv3.EventTypePut:
			node := &pb.Node{}
			if err := proto.Unmarshal(event.Kv.Value, node); err != nil {
				mlog.Errorf("解析服务节点失败: %v", err)
			} else {
				if cls.Add(node) {
					mlog.Infof("添加服务节点: %v", node)
				}
			}
		case clientv3.EventTypeDelete:
			key := string(event.Kv.Key)
			id := cast.ToInt32(path.Base(key))
			typ := pb.NodeType(cast.ToInt32(path.Base(path.Dir(key))))
			if cls.Del(typ, id) {
				mlog.Infof("删除服务节点: %s", key)
			}
		}
	}
}

func (d *Etcd) Register(n *pb.Node, ttl int64) error {
	ctx := context.Background()
	tctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// 创建租约
	if err := util.Retry(3, time.Second, func() error {
		if rsp, err := d.client.Grant(tctx, ttl); err != nil {
			return err
		} else {
			d.lease = clientv3.LeaseID(rsp.ID)
			return nil
		}
	}); err != nil {
		return err
	}

	// 注册服务
	buf, err := proto.Marshal(n)
	if err != nil {
		return err
	}
	channel := fmt.Sprintf("%s/%d/%d", d.topic, n.Type, n.Id)
	if _, err = d.client.Put(ctx, channel, string(buf), clientv3.WithLease(d.lease)); err != nil {
		return err
	}

	// 保活
	keep, err := d.client.KeepAlive(ctx, d.lease)
	if err != nil {
		return err
	}

	d.Add(1)
	safe.Go(func() {
		defer d.Done()
		d.keepAlive(keep, n, ttl)
	})
	return nil
}

func (d *Etcd) keepAlive(keep <-chan *clientv3.LeaseKeepAliveResponse, n *pb.Node, ttl int64) {
	tt := time.NewTicker(time.Duration(ttl/2) * time.Second)
	defer func() {
		tt.Stop()
		d.client.Revoke(context.Background(), d.lease)
	}()
	for {
		select {
		case _, ok := <-keep:
			if !ok {
				if err := d.Register(n, ttl); err != nil {
					mlog.Errorf("尝试重新注册失败 %v", err)
				} else {
					return
				}
			}
		case <-tt.C:
			if _, err := d.client.TimeToLive(context.Background(), d.lease); err != nil {
				mlog.Errorf("保活异常, 尝试重新注册 %v", err)
				if err := d.Register(n, ttl); err != nil {
					mlog.Errorf("尝试重新注册失败 %v", err)
				} else {
					return
				}
			}
		case <-d.exit:
			return
		}
	}
}

func (d *Etcd) Close() error {
	close(d.exit)
	return d.client.Close()
}
