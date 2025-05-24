package lru

type Getter interface {
	Get(key string) (Value, error)
}

type GetterFunc func(key string) (Value, error)

func (f GetterFunc) Get(key string) (Value, error) {
	return f(key)
}
