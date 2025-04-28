package network

import "universal/framework/define"

type OpOption func(*Op)

type Op struct {
	topic  string
	parse  define.ParsePacketFunc
	newFun define.NewPacketFunc
}

func WithTopic(p string) OpOption {
	return func(o *Op) {
		o.topic = p
	}
}

func WithParse(p define.ParsePacketFunc) OpOption {
	return func(o *Op) {
		o.parse = p
	}
}

func WithNew(f define.NewPacketFunc) OpOption {
	return func(o *Op) {
		o.newFun = f
	}
}
