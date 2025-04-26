package network

import "universal/framework/define"

type OpOption func(*Op)

type Op struct {
	root    string
	parse   define.ParsePacketFunc
	cluster define.ICluster
	table   define.ITable
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

func WithTable(r define.ITable) OpOption {
	return func(o *Op) {
		o.table = r
	}
}
