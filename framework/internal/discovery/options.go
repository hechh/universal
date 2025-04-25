package discovery

import "universal/framework/define"

type Option func(*options)

type options struct {
	root  string
	parse define.ParseNodeFunc
}

func WithPath(p string) Option {
	return func(o *options) {
		o.root = p
	}
}

func WithParse(p define.ParseNodeFunc) Option {
	return func(o *options) {
		o.parse = p
	}
}
