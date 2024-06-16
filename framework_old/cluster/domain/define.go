package domain

const (
	NONE   = 0
	ADD    = 1
	DELETE = 2
)

type WatchFunc func(action int, key string, value string)

type IDiscovery interface {
	KeepAlive(string, string, int64) // 设置保活key
	Watch(string, WatchFunc) error   // 开启协程watch+keepalive
	Close()                          // 停止协程watch+keepalive
}
