package lru

import "sync"

//防止缓存击穿：一个存在的key，在缓存过期的一刻，同时有大量的请求，这些请求都会击穿到 DB ，造成瞬时DB请求量大、压力骤增。

//代表正在进行中，或已经结束的请求
type call struct {
	sync.WaitGroup //避免重入：可能在f调用期间有n次调用，但是只执行一次
	val            interface{}
	err            error
}

//SingleFlight 管理不同key的请求
type SingleFlight struct {
	sync.Mutex //保护m并发读写
	m          map[string]*call
}

// Do 针对相同的key，无论Do被调用多少次，f只执行一次
func (g *SingleFlight) Do(key string, f func() (interface{}, error)) (interface{}, error) {
	g.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	//如果已有结果，返回
	if c, ok := g.m[key]; ok {
		g.Unlock()
		c.Wait() //等待f调用结束了，再返回
		return c.val, c.err
	}
	c := new(call)
	c.Add(1)
	g.m[key] = c //缓存执行结果
	g.Unlock()

	//调用函数
	c.val, c.err = f()
	c.Done()

	g.Lock()
	delete(g.m, key)
	defer g.Unlock()

	return c.val, c.err
}
