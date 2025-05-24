package discovery

import (
	"context"
	"fmt"
	"path"
	"time"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/domain"
	"universal/library/async"
	"universal/library/mlog"
	"universal/library/uerror"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/cast"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Etcd struct {
	client    *clientv3.Client
	topic     string
	watchChan <-chan clientv3.WatchResponse
	lease     clientv3.LeaseID
	exit      chan struct{}
}

func NewEtcd(cfg *yaml.EtcdConfig) (*Etcd, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:            cfg.Endpoints,
		DialTimeout:          5 * time.Second,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
		MaxCallSendMsgSize:   10 * 1024 * 1024,
	})
	if err != nil {
		return nil, uerror.New(1, -1, "Etcd连接失败: %v", err)
	}
	return &Etcd{
		client: cli,
		topic:  cfg.Channel,
		exit:   make(chan struct{}),
	}, nil
}

// 监听服务
func (d *Etcd) Watch(cls domain.ICluster) error {
	rsp, err := d.client.Get(context.Background(), d.topic, clientv3.WithPrefix())
	if err != nil {
		return uerror.New(1, -1, "获取注册服务节点失败: %v", err)
	}
	for _, kv := range rsp.Kvs {
		node := &pb.Node{}
		if err := proto.Unmarshal(kv.Value, node); err != nil {
			return uerror.New(1, -1, "解析服务节点失败: %v", err)
		} else {
			cls.Add(node)
			mlog.Infof("添加服务节点: %s", node.String())
		}
	}
	// 监听服务
	d.watchChan = d.client.Watch(context.Background(), d.topic, clientv3.WithPrefix())
	async.SafeGo(mlog.Fatalf, func() {
		for resp := range d.watchChan {
			if resp.Err() != nil {
				mlog.Errorf("Etcd监听服务失败: %v", resp.Err())
				continue
			}

			for _, event := range resp.Events {
				mlog.Infof("Etcd监听服务事件: %s ---> %s", event.Kv.Key, string(event.Kv.Value))
				switch event.Type {
				case clientv3.EventTypePut:
					node := &pb.Node{}
					if err := proto.Unmarshal(event.Kv.Value, node); err != nil {
						mlog.Errorf("解析服务节点失败: %v", err)
					} else {
						cls.Add(node)
						mlog.Infof("添加服务节点: %v", node)
					}
				case clientv3.EventTypeDelete:
					key := string(event.Kv.Key)
					id := cast.ToInt32(path.Base(key))
					typ := pb.NodeType(cast.ToInt32(path.Base(path.Dir(key))))
					cls.Del(typ, id)
					mlog.Infof("删除服务节点: %s", key)
				}
			}
		}
	})
	return nil
}

func (d *Etcd) getTopic(nodeType pb.NodeType, nodeId int32) string {
	return fmt.Sprintf("%s/%d/%d", d.topic, nodeType, nodeId)
}

// 注册服务
func (d *Etcd) Register(nn *pb.Node, ttl int64) error {
	// 创建租约
	rsp, err := d.client.Grant(context.Background(), ttl)
	if err != nil {
		return uerror.New(1, -1, "Etcd获取租约失败, node:%v, error:%v", nn, err)
	}
	d.lease = clientv3.LeaseID(rsp.ID)

	// 注册服务
	buf, err := proto.Marshal(nn)
	if err != nil {
		return uerror.New(1, -1, "Etcd序列化服务节点失败, node:%v, error:%v", nn, err)
	}
	topic := d.getTopic(nn.Type, nn.Id)
	if _, err := d.client.Put(context.Background(), topic, string(buf), clientv3.WithLease(d.lease)); err != nil {
		return uerror.New(1, -1, "Etcd注册服务节点失败, node:%v, error:%v", nn, err)
	}

	// 定时保活
	go func() {
		tt := time.NewTicker(time.Duration(ttl/2) * time.Second)
		defer tt.Stop()

		for {
			select {
			case <-tt.C:
				i := 0
				for ; i < 3; i++ {
					if _, err := d.client.KeepAliveOnce(context.Background(), d.lease); err != nil {
						mlog.Errorf("Etcd续约失败: %v", err)
						time.Sleep(time.Second)
						continue
					}
					break
				}
				if i >= 3 {
					panic(uerror.New(1, -1, "Etcd续约失败: ndoe:%v, error:%v", nn, err))
				}
			case <-d.exit:
				return
			}
		}
	}()
	return nil
}

// 关闭服务
func (d *Etcd) Close() error {
	return d.client.Close()
}
