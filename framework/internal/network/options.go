package network

type OpOption func(*Op)

type Op struct {
	root string
}

func WithPath(p string) OpOption {
	return func(o *Op) {
		o.root = p
	}
}
