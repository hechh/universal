package discovery

import (
	"context"
	"fmt"
	"path"
	"time"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/framework/library/async"
	"universal/framework/library/mlog"
	"universal/framework/library/uerror"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/cast"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Etcd struct {
	client  *clientv3.Client
	newNode func() *pb.Node
	topic   string
	exit    chan struct{}
}

func NewEtcd(endpoints []string, opts ...OpOption) (*Etcd, error) {
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		return nil, err
	}
	vals := NewOp(opts...)
	return &Etcd{
		client:  cli,
		newNode: vals.newNode,
		topic:   vals.topic,
		exit:    make(chan struct{}),
	}, nil
}

// 监听 kv 变化
func (e *Etcd) Watch(cls domain.ICluster) error {
	// 获取在线服务
	rsp, err := e.client.Get(context.Background(), e.topic, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range rsp.Kvs {
		node := &pb.Node{}
		if err := proto.Unmarshal(kv.Value, node); err != nil {
			return uerror.New(1, -1, "解析服务节点失败: %v", err)
		} else {
			cls.Add(node)
			mlog.Info(" 添加服务节点: %s", node.String())
		}
	}

	// 监听服务
	listens := e.client.Watch(context.Background(), e.topic, clientv3.WithPrefix())
	async.SafeGo(mlog.Error, func() {
		for listen := range listens {
			for _, event := range listen.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					node := &pb.Node{}
					if err := proto.Unmarshal(event.Kv.Value, node); err != nil {
						mlog.Error("解析服务节点失败: %v", err)
					} else {
						cls.Add(node)
						mlog.Info(" 添加服务节点: %v", node.String())
					}

				case clientv3.EventTypeDelete:
					id := cast.ToInt32(path.Base(string(event.Kv.Key)))
					typ := cast.ToInt32(path.Base(path.Dir(string(event.Kv.Key))))
					cls.Del(pb.NodeType(typ), id)
					mlog.Info(" 删除服务节点: %s", event.Kv.Key)
				}
			}
		}
	})
	return nil
}

// 注册服务节点
func (d *Etcd) Register(node *pb.Node, ttl int64) error {
	// 1. 创建租约
	rsp, err := d.client.Grant(context.Background(), ttl)
	if err != nil {
		return uerror.New(1, -1, "租约创建失败: %v", err)
	}
	lease := clientv3.LeaseID(rsp.ID)
	defer d.client.Revoke(context.Background(), lease)

	// 2. 序列化节点数据
	buf, err := proto.Marshal(node)
	if err != nil {
		return uerror.New(1, -1, "序列化服务节点失败: %v", err)
	}
	topic := path.Join(d.topic, cast.ToString(node.GetType()), cast.ToString(node.GetId()))

	// 3. 保存节点
	if _, err := d.client.Put(context.Background(),
		topic, string(buf),
		clientv3.WithLease(lease),
	); err != nil {
		return uerror.New(1, -1, "Etcd注册服务节点失败: %v", err)
	}

	// 定时续约保活
	go func() {
		tt := time.NewTicker(time.Duration(ttl/2) * time.Second)
		defer tt.Stop()

		for {
			select {
			case <-tt.C:
				i := 0
				for ; i < 3; i++ {
					if _, err := d.client.KeepAliveOnce(context.Background(), lease); err != nil {
						mlog.Error("Etcd租赁续约保活失败: node:%s, error:%v", node.String(), err)
						time.Sleep(1 * time.Second)
						continue
					}
					break
				}
				if i >= 3 {
					panic(fmt.Sprintf("Etcd租赁续约失败，且超过最大重试次数: %s", node.String()))
				}
			case <-d.exit:
				return
			}
		}
	}()
	return nil
}

func (d *Etcd) Close() error {
	return d.client.Close()
}
