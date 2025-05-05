package discovery

import (
	"universal/framework/config"
	"universal/framework/define"
	"universal/library/baselib/uerror"
)

type OpOption func(*Op)

type Op struct {
	status  int32
	topic   string
	newNode func() define.INode
}

func NewOp(opts ...OpOption) *Op {
	ret := &Op{}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

func WithTopic(p string) OpOption {
	return func(o *Op) {
		o.topic = p
	}
}

func WithNode(p func() define.INode) OpOption {
	return func(o *Op) {
		o.newNode = p
	}
}

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
