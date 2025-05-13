package lru

//cache.go 的实现非常简单，实例化 lru，封装 get 和 add 方法，并添加互斥锁 mu。
import (
	"sync"
)

type Cache struct {
	sync.RWMutex
	lru      *LRUCache
	cacheLen int64
}

func NewCache(cacheLen int64) *Cache {
	return &Cache{
		cacheLen: cacheLen,
	}
}

// Add 向LRU添加数据
func (c *Cache) Add(key string, value Value, onEvited func(Value)) {
	c.Lock()
	defer c.Unlock()
	if c.lru == nil {
		c.lru = New(c.cacheLen, onEvited)
	}
	c.lru.Add(key, value)
}

// Get 从lru获取数据
func (c *Cache) Get(key string) (value Value, ok bool) {
	c.RLock()
	defer c.RUnlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(Value), ok
	}
	return
}

func (c *Cache) Del(key string) {
	c.Lock()
	defer c.Unlock()
	if c.lru == nil {
		return
	}
	c.lru.Remove(key)
}

func (c *Cache) IsFull() bool {
	c.RLock()
	defer c.RUnlock()
	if c.lru == nil {
		return false
	}
	return c.lru.IsFull()
}

func (c *Cache) CallRange(f func(value Value, cache *LRUCache)) {
	c.Lock()
	defer c.Unlock()
	if c.lru == nil {
		return
	}
	for _, ele := range c.lru.cache { //如果键对应的链表节点存在，返回查找到的值。
		kv := ele.Value.(*entry)
		f(kv.value, c.lru)
	}
}

func (c *Cache) Call(f func(cache *LRUCache), onEvited func(Value)) {
	c.Lock()
	defer c.Unlock()
	if c.lru == nil {
		c.lru = New(c.cacheLen, onEvited)
	}
	f(c.lru)
}
