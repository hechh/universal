package discovery

import (
	"universal/framework/config"
	"universal/framework/define"
	"universal/library/baselib/uerror"
)

func Init(cfg *config.Config, opts ...OpOption) (define.IDiscovery, error) {
	if cfg.Etcd != nil {
		dis, err := NewEtcd(cfg.Etcd.Endpoints, opts...)
		if err != nil {
			return nil, err
		}
		return dis, nil
	}

	if cfg.Consul != nil {
		dis, err := NewConsul(cfg.Consul.Endpoints, opts...)
		if err != nil {
			return nil, err
		}
		return dis, nil
	}
	return nil, uerror.New(1, -1, "服务注册与发现配置错误")
}
