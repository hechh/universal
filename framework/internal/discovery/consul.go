package discovery

import (
	"fmt"
	"path"
	"time"
	"universal/framework/domain"
	"universal/framework/library/mlog"
	"universal/framework/library/uerror"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/spf13/cast"
)

type Consul struct {
	client  *api.Client
	session *api.Session
	newNode func() domain.INode
	topic   string
	keys    map[string]struct{}
	exit    chan struct{}
}

func NewConsul(endpoints string, opts ...OpOption) (*Consul, error) {
	client, err := api.NewClient(&api.Config{Address: endpoints})
	if err != nil {
		return nil, err
	}
	vals := new(Op)
	for _, opt := range opts {
		opt(vals)
	}
	// 返回
	return &Consul{
		client:  client,
		session: client.Session(),
		newNode: vals.newNode,
		topic:   vals.topic,
		keys:    make(map[string]struct{}),
		exit:    make(chan struct{}),
	}, nil
}

func (d *Consul) update(cls domain.ICluster) error {
	// 添加服务
	kvs, _, err := d.client.KV().List(d.topic, nil)
	if err != nil {
		return err
	}
	tmps := map[string]struct{}{}
	for _, kv := range kvs {
		node := d.newNode()
		if err := node.ReadFrom(kv.Value); err != nil {
			return err
		}
		tmps[kv.Key] = struct{}{}
		d.keys[kv.Key] = struct{}{}
		cls.Add(node)
		mlog.Info("添加服务节点: %s", node.String())
	}

	// 删除服务
	for k := range d.keys {
		if _, ok := tmps[k]; ok {
			continue
		}
		id := cast.ToInt32(path.Base(k))
		typ := cast.ToInt32(path.Base(path.Dir(k)))
		cls.Del(typ, id)
		mlog.Info("删除服务节点: %s", k)
		delete(d.keys, k)
	}
	return nil
}

func (d *Consul) Watch(cls domain.ICluster) error {
	if err := d.update(cls); err != nil {
		return uerror.New(1, -1, " 获取服务节点失败: %v", err)
	}
	// 监听服务
	w, err := watch.Parse(map[string]interface{}{
		"type":   "keyprefix",
		"prefix": d.topic,
	})
	if err != nil {
		return err
	}

	// 设置监听回调
	w.Handler = func(idx uint64, data interface{}) {
		if err := d.update(cls); err != nil {
			mlog.Error("consul监听服务更新失败: %v", err)
		}
	}
	go func() {
		if err := w.RunWithClientAndHclog(d.client, nil); err != nil {
			mlog.Error("consul 监听报错: %v", err)
			panic(err)
		}
	}()
	return nil
}

func (d *Consul) Register(node domain.INode, ttl int64) error {
	// 设置租约
	leaseId, _, err := d.session.Create(&api.SessionEntry{
		TTL:      fmt.Sprintf("%ds", ttl),
		Behavior: "delete",
	}, nil)
	if err != nil {
		return uerror.New(1, -1, " 设置租约失败: %v", err)
	}

	// 2. 序列化节点数据
	topic := path.Join(d.topic, cast.ToString(node.GetType()), cast.ToString(node.GetId()))
	buf := make([]byte, node.GetSize())
	if err := node.WriteTo(buf); err != nil {
		return err
	}

	// 注册服务
	if _, err := d.client.KV().Put(&api.KVPair{
		Key:     topic,
		Value:   buf,
		Session: leaseId,
	}, nil); err != nil {
		return uerror.New(1, -1, "服务注册失败: %v", err)
	}

	// 定时检测
	go func() {
		tt := time.NewTicker(time.Duration(ttl/2) * time.Second)
		defer tt.Stop()

		for {
			select {
			case <-tt.C:
				i := 0
				for ; i < 3; i++ {
					if _, _, err := d.session.Renew(leaseId, nil); err != nil {
						mlog.Error("Consul租赁续约保活失败: node:%s, error:%v", node.String(), err)
						time.Sleep(1 * time.Second)
						continue
					}
					break
				}
				if i >= 3 {
					panic(fmt.Sprintf("Etcd租赁续约失败，且超过最大重试次数: %s", node.String()))
				}
			case <-d.exit:
				d.session.Destroy(leaseId, nil)
				return
			}
		}
	}()
	return nil
}

func (d *Consul) Close() error {
	return nil
}
