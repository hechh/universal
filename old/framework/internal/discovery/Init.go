package discovery

import (
	"universal/common/pb"
	"universal/framework/domain"
	"universal/framework/global"
	"universal/framework/library/uerror"
)

type OpOption func(*Op)

type Op struct {
	topic   string
	newNode func() *pb.Node
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

func WithNode(p func() *pb.Node) OpOption {
	return func(o *Op) {
		o.newNode = p
	}
}

func Init(cfg *global.Config, opts ...OpOption) (domain.IDiscovery, error) {
	if cfg.Etcd != nil {
		dis, err := NewEtcd(cfg.Etcd.Endpoints, opts...)
		if err != nil {
			return nil, err
		}
		return dis, nil
	}

	/*
		if cfg.Consul != nil {
			dis, err := NewConsul(cfg.Consul.Endpoints, opts...)
			if err != nil {
				return nil, err
			}
			return dis, nil
		}
	*/
	return nil, uerror.New(1, -1, "服务注册与发现配置错误")
}
