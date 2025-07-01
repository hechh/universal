package config

import (
	"context"
	"path"
	"path/filepath"
	"sync"
	"time"
	"universal/common/yaml"
	"universal/library/mlog"
	"universal/library/safe"
	"universal/library/uerror"
	"universal/library/util"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Watcher struct {
	sync.WaitGroup
	topic  string
	cpath  string
	client *clientv3.Client
	exit   chan struct{}
}

func NewWatcher(cfg *yaml.DataConfig) (*Watcher, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:            cfg.Endpoints,
		DialTimeout:          5 * time.Second,
		DialKeepAliveTime:    30 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
		MaxCallSendMsgSize:   10 * 1024 * 1024,
	})
	if err != nil {
		return nil, uerror.N(1, -1, "Etcd连接失败: %v", err)
	}
	return &Watcher{
		client: cli,
		topic:  cfg.Topic,
		cpath:  cfg.Path,
		exit:   make(chan struct{}),
	}, nil
}

func (d *Watcher) Close() error {
	d.exit <- struct{}{}
	close(d.exit)
	d.Wait()
	return d.client.Close()
}

func (d *Watcher) Load(tmps map[string]struct{}) error {
	rsp, err := d.client.Get(context.Background(), d.topic, clientv3.WithPrefix())
	if err != nil {
		return uerror.N(1, -1, "获取注册服务节点失败: %v", err)
	}

	for _, kv := range rsp.Kvs {
		sheet := path.Base(string(kv.Key))
		if f, ok := fileMgr[sheet]; ok {
			if err := f(kv.Value); err != nil {
				return uerror.N(1, -1, "加载%s配置错误： %v", sheet, err)
			}
			tmps[sheet] = struct{}{}

			if err := util.Save(d.cpath, sheet+".conf", kv.Value); err == nil {
				mlog.Infof("更新配置：%s", filepath.Join(d.cpath, sheet+".conf"))
			}
		}
	}
	return nil
}

func (d *Watcher) Watch(tmps map[string]struct{}) {
	wg := &sync.WaitGroup{}
	wg.Add(len(fileMgr) - len(tmps))
	safe.Go(func() {
		for {
			watchChan := d.client.Watch(context.Background(), d.topic, clientv3.WithPrefix())
			if watchChan == nil {
				mlog.Errorf("Config监听服务失败: watchChan is nil")
				time.Sleep(1 * time.Second)
				continue
			}

			for resp := range watchChan {
				if resp.Canceled {
					mlog.Errorf("Config监听被取消，尝试重新连接")
					break
				}
				if resp.Err() != nil {
					mlog.Errorf("Config监听服务失败: %v", resp.Err())
					continue
				}
				for _, event := range resp.Events {
					if event.Type != clientv3.EventTypePut {
						continue
					}
					sheet := path.Base(string(event.Kv.Key))
					f, ok := fileMgr[sheet]
					if !ok {
						continue
					}
					if err := util.Save(d.cpath, sheet+".conf", event.Kv.Value); err == nil {
						mlog.Infof("更新配置：%s", filepath.Join(d.cpath, sheet+".conf"))
					}

					if err := f(event.Kv.Value); err != nil {
						mlog.Errorf("加载%s配置错误： %v", sheet, err)
					} else {
						if _, ok := tmps[sheet]; !ok {
							tmps[sheet] = struct{}{}
							wg.Done()
						}
					}
				}
			}
		}
	})
	wg.Wait()
}
