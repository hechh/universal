package network

import (
	"universal/framework/config"
	"universal/framework/define"
	"universal/library/baselib/uerror"
)

type OpOption func(*Op)

type Op struct {
	topic     string
	newPacket func() define.IPacket
	newHeader func() define.IHeader
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

func WithPacket(p func() define.IPacket) OpOption {
	return func(o *Op) {
		o.newPacket = p
	}
}

func WithHeader(f func() define.IHeader) OpOption {
	return func(o *Op) {
		o.newHeader = f
	}
}

func Init(cfg *config.Config, opts ...OpOption) (define.INetwork, error) {
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
