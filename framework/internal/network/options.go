package network

import "universal/framework/define"

type OpOption func(*Op)

type Op struct {
	root      string
	parse     define.ParsePacketFunc
	cluster   define.ICluster
	routerMgr define.IRouterMgr
}

func WithPath(p string) OpOption {
	return func(o *Op) {
		o.root = p
	}
}

func WithParse(p define.ParsePacketFunc) OpOption {
	return func(o *Op) {
		o.parse = p
	}
}

func WithCluster(cls define.ICluster) OpOption {
	return func(o *Op) {
		o.cluster = cls
	}
}

func WithRouterMgr(r define.IRouterMgr) OpOption {
	return func(o *Op) {
		o.routerMgr = r
	}
}
