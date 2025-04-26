package discovery

import "universal/framework/define"

type OpOption func(*Op)

type Op struct {
	status int32
	root   string
	parse  define.ParseNodeFunc
}

func WithPath(p string) OpOption {
	return func(o *Op) {
		o.root = p
	}
}

func WithParse(p define.ParseNodeFunc) OpOption {
	return func(o *Op) {
		o.parse = p
	}
}
