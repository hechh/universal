package discovery

import (
	"fmt"
	"path"
	"sync/atomic"
	"time"
	"universal/framework/define"
	"universal/library/baselib/safe"
	"universal/library/mlog"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/spf13/cast"
)

type Consul struct {
	status  int32
	topic   string
	newNode func() define.INode
	client  *api.Client
	session *api.Session
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
		topic:   vals.topic,
		newNode: vals.newNode,
		client:  client,
		session: client.Session(),
		keys:    make(map[string]struct{}),
		exit:    make(chan struct{}),
	}, nil
}

func (c *Consul) getKey(node define.INode) string {
	return path.Join(c.topic, cast.ToString(node.GetType()), cast.ToString(node.GetId()))
}

func (c *Consul) Get() (rets []define.INode, err error) {
	kvs, _, err := c.client.KV().List(c.topic, nil)
	if err != nil {
		return nil, err
	}
	for _, kv := range kvs {
		rets = append(rets, c.newNode().Parse(kv.Value))
	}
	return
}

// 添加 key-value
func (c *Consul) Put(srv define.INode) error {
	kv := &api.KVPair{
		Key:   c.getKey(srv),
		Value: srv.ToBytes(),
	}
	_, err := c.client.KV().Put(kv, nil)
	return err
}

func (c *Consul) Del(srv define.INode) error {
	_, err := c.client.KV().Delete(c.getKey(srv), nil)
	return err
}

func (c *Consul) update(cluster define.ICluster) error {
	// 获取全部节点
	kvs, _, err := c.client.KV().List(c.topic, nil)
	if err != nil {
		return err
	}
	// 添加服务
	tmps := map[string]struct{}{}
	for _, kv := range kvs {
		tmps[kv.Key] = struct{}{}
		c.keys[kv.Key] = struct{}{}
		nn := c.newNode().Parse(kv.Value)
		cluster.Put(nn)
		mlog.Info(" 添加服务节点: %v", nn)
	}
	// 删除服务
	for k := range c.keys {
		if _, ok := tmps[k]; ok {
			continue
		}
		id := cast.ToUint32(path.Base(k))
		typ := cast.ToUint32(path.Base(path.Dir(k)))
		cluster.Del(typ, id)
		mlog.Info(" 删除服务节点: %s", k)
		delete(c.keys, k)
	}
	return nil
}

func (c *Consul) Watch(cluster define.ICluster) error {
	// 先获取在线服务
	if err := c.update(cluster); err != nil {
		return err
	}

	// 监听服务
	w, err := watch.Parse(map[string]interface{}{
		"type":   "keyprefix",
		"prefix": c.topic,
	})
	if err != nil {
		return err
	}

	// 设置监听回调
	w.Handler = func(idx uint64, data interface{}) {
		if err := c.update(cluster); err != nil {
			mlog.Error("consul 监听服务更新失败: %v", err)
		}
	}
	safe.SafeGo(mlog.Fatal, func() {
		if err := w.RunWithClientAndHclog(c.client, nil); err != nil {
			mlog.Error("consul 监听报错: %v", err)
			panic(err)
		}
	})
	return nil
}

func (c *Consul) KeepAlive(srv define.INode, ttl int64) error {
	// 设置租约
	leaseID, _, err := c.session.Create(&api.SessionEntry{
		TTL:      fmt.Sprintf("%ds", ttl),
		Behavior: "delete",
	}, nil)
	if err != nil {
		return err
	}

	// 设置 key-value
	_, err = c.client.KV().Put(&api.KVPair{
		Key:     c.getKey(srv),
		Value:   srv.ToBytes(),
		Session: leaseID,
	}, nil)
	if err != nil {
		return err
	}

	// 定时检测
	safe.SafeGo(mlog.Error, func() {
		tt := time.NewTicker(time.Duration(ttl/2) * time.Second)
		defer tt.Stop()
		atomic.AddInt32(&c.status, 1)

		for {
			select {
			case <-c.exit:
				c.session.Destroy(leaseID, nil)
				return
			case <-tt.C:
				if _, _, err := c.session.Renew(leaseID, nil); err != nil {
					mlog.Error("consul 租赁续约保活失败：%v", err)
					return
				}
			}
		}
	})
	return nil
}

func (c *Consul) Close() error {
	if atomic.LoadInt32(&c.status) > 0 {
		c.exit <- struct{}{}
	}
	time.Sleep(1 * time.Second)
	return nil
}
