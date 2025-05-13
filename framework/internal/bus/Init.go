package network

import (
	"universal/framework/domain"
	"universal/framework/global"
	"universal/framework/library/uerror"
)

type OpOption func(*Op)

type Op struct {
	topic    string
	routeMgr domain.IRouterMgr
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

func WithRoute(f func() domain.IRouter) OpOption {
	return func(o *Op) {
		o.newRoute = f
	}
}

func WithRouteMgr(rr domain.IRouterMgr) OpOption {
	return func(o *Op) {
		o.routeMgr = rr
	}
}

func Init(cfg *global.Config, opts ...OpOption) (domain.INetwork, error) {
	if cfg.Nats != nil {
		dd, err := NewNats(cfg.Nats.Endpoints, opts...)
		if err != nil {
			return nil, err
		}
		return dd, nil
	}

	if cfg.Nsq != nil {
		dd, err := NewNsq(cfg.Nsq.Nsqd, opts...)
		if err != nil {
			return nil, err
		}
		return dd, nil
	}

	return nil, uerror.New(1, -1, " 消息中间件配置错误")
}
