package network

import (
	"universal/framework/config"
	"universal/framework/define"
	"universal/library/baselib/uerror"
)

func Init(cfg *config.Config, opts ...OpOption) (define.INetwork, error) {
	if cfg.Nats != nil {
		dd, err := NewNats(cfg.Nats.Endpoints, opts...)
		if err != nil {
			return nil, err
		}
		return dd, nil
	}

	return nil, uerror.New(1, -1, " 消息中间件配置错误")
}
