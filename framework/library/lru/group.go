package lru

import (
	"strconv"
	"sync"
)

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

type Group struct {
	name           string
	cache          []*Cache
	getter         Getter //LRU未命中时load数据的回调
	lruCap         int
	lruCacheArrLen int
	onEvited       func(Value)   //删除数据的回调
	loader         *SingleFlight //缓存击穿
}

//GetGroup 用来特定名称的 Group
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

func NewGroup(name string, cacheCap, hashLen int, getter Getter, onEvited func(Value)) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	g := &Group{
		name:           name,
		getter:         getter,
		lruCap:         cacheCap,
		lruCacheArrLen: hashLen,
		cache:          make([]*Cache, hashLen),
		onEvited:       onEvited,
		loader:         &SingleFlight{},
	}
	for i := 0; i < g.lruCacheArrLen; i++ {
		g.cache[i] = NewCache(int64(g.lruCap / g.lruCacheArrLen))
	}
	mu.Lock()
	defer mu.Unlock()
	groups[name] = g
	return g
}

func (g *Group) Put(v Value) {
	hashCode := g.getHashKey(v.GetKey())
	g.cache[hashCode].Add(v.GetKey(), v, g.onEvited)
}

func (g *Group) Get(key string) (val Value, err error, ok bool) {
	hashCode := g.getHashKey(key)
	//cache hit，直接返回
	if val, ok = g.cache[hashCode].Get(key); ok {
		return
	} else {
		//获取失败，调用回调添加数据
		val, err = g.getLocally(key)
		return
	}
}

func (g *Group) Lock(key string) {
	hashCode := g.getHashKey(key)
	m := g.cache[hashCode]
	m.Lock()
}
func (g *Group) UnLock(key string) {
	hashCode := g.getHashKey(key)
	m := g.cache[hashCode]
	m.Unlock()
}

func (g *Group) Del(key string) {
	hashCode := g.getHashKey(key)
	//cache hit，直接返回
	m := g.cache[hashCode]
	m.Del(key)
}

func (g *Group) CallRange(index int, f func(value Value, cache *LRUCache)) {
	index %= g.lruCacheArrLen
	m := g.cache[index]
	m.CallRange(f)
}

func (g *Group) Call(key string, f func(cache *LRUCache)) {
	hashCode := g.getHashKey(key)
	//cache hit，直接返回
	m := g.cache[hashCode]
	m.Call(f, g.onEvited)
}

func (g *Group) getLocally(key string) (Value, error) {
	//防止缓存击穿
	view, err := g.loader.Do(key, func() (interface{}, error) {
		val, err := g.getter.Get(key)
		if err != nil || val == nil {
			return nil, err //获取失败
		} else {
			//获取成功，添加到Cache
			hashCode := g.getHashKey(key)
			g.cache[hashCode].Add(key, val, g.onEvited)
			return val, nil
		}
	})
	if err == nil {
		return view.(Value), nil
	} else {
		return nil, err
	}

}

func (g *Group) getHashKey(uid string) int32 {
	uid64, _ := strconv.ParseInt(uid, 10, 64)
	return int32(uid64 % int64(g.lruCacheArrLen))
}
