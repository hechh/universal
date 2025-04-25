package discovery

import (
	"fmt"
	"path"
	"time"
	"universal/framework/define"
	"universal/library/baselib/safe"
	"universal/library/mlog"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/spf13/cast"
)

type Consul struct {
	options *options
	client  *api.Client
	session *api.Session
	keys    map[string]struct{}
	exit    chan struct{}
}

func NewConsul(endpoints string, opts ...Option) (*Consul, error) {
	client, err := api.NewClient(&api.Config{Address: endpoints})
	if err != nil {
		return nil, err
	}
	vals := new(options)
	for _, opt := range opts {
		opt(vals)
	}
	return &Consul{
		options: vals,
		client:  client,
		session: client.Session(),
		keys:    make(map[string]struct{}),
		exit:    make(chan struct{}),
	}, nil
}

func (c *Consul) getKey(node define.INode) string {
	return fmt.Sprintf("%s/%d/%d", c.options.root, node.GetType(), node.GetId())
}

func (c *Consul) Close() error {
	return nil
}

// 添加 key-value
func (c *Consul) Put(srv define.INode) error {
	_, err := c.client.KV().Put(&api.KVPair{Key: c.getKey(srv), Value: srv.ToBytes()}, nil)
	return err
}

func (c *Consul) Del(srv define.INode) error {
	_, err := c.client.KV().Delete(c.getKey(srv), nil)
	return err
}

func (c *Consul) Watch(cluster define.ICluster) error {
	// 先获取在线服务
	kv := c.client.KV()
	kvs, _, err := kv.List(c.options.root, nil)
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		c.keys[kv.Key] = struct{}{}
		cluster.Put(c.options.parse(kv.Value))
	}

	// 监听服务
	w, err := watch.Parse(map[string]interface{}{
		"type":   "keyprefix",
		"prefix": c.options.root,
	})
	if err != nil {
		return err
	}
	// 设置监听回调
	w.Handler = func(idx uint64, data interface{}) {
		kvs := data.(api.KVPairs)
		tmps := map[string]struct{}{}
		// 添加服务
		for _, kv := range kvs {
			tmps[kv.Key] = struct{}{}
			if err := cluster.Put(c.options.parse(kv.Value)); err != nil {
				mlog.Error("consul发现新服务，新服务添加失败: %v", err)
			}
		}
		// 删除服务
		for k := range c.keys {
			if _, ok := tmps[k]; ok {
				continue
			}
			id := cast.ToInt32(path.Base(k))
			typ := cast.ToInt32(path.Base(path.Dir(k)))
			if err := cluster.Del(typ, id); err != nil {
				mlog.Error("consul发现服务下线，删除服务失败: %v", err)
			} else {
				delete(c.keys, k)
			}
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
	entry := &api.SessionEntry{
		TTL:      fmt.Sprintf("%ds", ttl),
		Behavior: "delete",
	}
	leaseID, _, err := c.session.Create(entry, nil)
	if err != nil {
		return err
	}

	// 设置 key-value
	kv := &api.KVPair{
		Key:     c.getKey(srv),
		Value:   srv.ToBytes(),
		Session: leaseID,
	}
	if _, err := c.client.KV().Put(kv, nil); err != nil {
		return err
	}

	// 定时检测
	tt := time.NewTicker(time.Duration(ttl/2) * time.Second)
	defer tt.Stop()
	safe.SafeGo(mlog.Error, func() {
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
