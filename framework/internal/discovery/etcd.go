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
	"universal/library/uerror"

	"github.com/spf13/cast"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
)

const (
	KEEP_STATUS     = 1
	REGISTER_STATUS = 2
)

type Etcd struct {
	sync.WaitGroup
	client *clientv3.Client
	lease  clientv3.LeaseID
	topic  string
	status int32
	exit   chan struct{}
}

func NewEtcd(cfg *yaml.EtcdConfig) (*Etcd, error) {
	client, err := clientv3.New(clientv3.Config{Endpoints: cfg.Endpoints})
	if err != nil {
		return nil, err
	}
	return &Etcd{client: client, topic: cfg.Topic, exit: make(chan struct{})}, nil
}

func (d *Etcd) Close() error {
	close(d.exit)
	return d.client.Close()
}

func (d *Etcd) Watch(cls domain.ICluster) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rsp, err := d.client.Get(ctx, d.topic, clientv3.WithPrefix())
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
			wchan := d.client.Watch(context.Background(), d.topic, clientv3.WithPrefix())
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
	if err := d.register(n, ttl); err != nil {
		return err
	}
	d.status = KEEP_STATUS

	safe.Go(func() {
		for {
			switch d.status {
			case KEEP_STATUS:
				if err := d.keepAlive(ttl); err != nil {
					mlog.Errorf("服务保活失败：%v", err)
					d.status = REGISTER_STATUS
				} else {
					return
				}
			case REGISTER_STATUS:
				if err := d.register(n, ttl); err != nil {
					mlog.Errorf("服务注册失败, 继续尝试注册服务: %v", err)
					time.Sleep(2 * time.Second)
				} else {
					d.status = KEEP_STATUS
				}
			}
		}
	})
	return nil
}

func (d *Etcd) register(n *pb.Node, ttl int64) error {
	rsp, err := d.client.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}
	d.lease = clientv3.LeaseID(rsp.ID)

	topic := fmt.Sprintf("%s/%d/%d", d.topic, n.Type, n.Id)
	buf, err := proto.Marshal(n)
	if err != nil {
		return err
	}
	_, err = d.client.Put(context.Background(), topic, string(buf), clientv3.WithLease(d.lease))
	return err
}

func (d *Etcd) keepAlive(ttl int64) error {
	kk, err := d.client.KeepAlive(context.Background(), d.lease)
	if err != nil {
		mlog.Errorf("Etcd 租约保活失败：%v", err)
		return err
	}

	tt := time.NewTicker(time.Duration(ttl/2) * time.Second)
	defer func() {
		tt.Stop()
	}()

	for {
		select {
		case _, ok := <-kk:
			if !ok {
				return uerror.N(1, -1, "保活异常,尝试重新注册")
			}
		case <-tt.C:
			if _, err := d.client.TimeToLive(context.Background(), d.lease); err != nil {
				return uerror.N(1, -1, "保活异常,尝试重新注册")
			}
		case <-d.exit:
			revokeCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if _, err := d.client.Revoke(revokeCtx, d.lease); err != nil {
				mlog.Errorf("租约释放失败: %v", err)
			}
			return nil
		}
	}
}
