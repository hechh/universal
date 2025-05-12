package network

import (
	"universal/framework/config"
	"universal/framework/domain"
	"universal/framework/library/uerror"
)

type OpOption func(*Op)

type Op struct {
	topic     string
	routeMgr  domain.IRouterMgr
	newRoute  func() domain.IRouter
	newHeader func() domain.IHead
	newPacket func() domain.IPacket
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

func WithPacket(p func() domain.IPacket) OpOption {
	return func(o *Op) {
		o.newPacket = p
	}
}

func WithHead(f func() domain.IHead) OpOption {
	return func(o *Op) {
		o.newHeader = f
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

func Init(cfg *config.Config, opts ...OpOption) (domain.INetwork, error) {
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
