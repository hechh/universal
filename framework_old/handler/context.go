package handler

import (
	"hego/common/pb"
	"sync"
)

type KValue struct {
	key   string
	value interface{}
}

type Context struct {
	sync.RWMutex             // 读写锁
	*pb.Head                 // 请求头
	data         interface{} // 关联数据
	temps        []*KValue   // 临时缓存
}

func NewContext(head *pb.Head, data interface{}) *Context {
	return &Context{data: data, Head: head}
}

// 获取关联数据
func (d *Context) GetData() interface{} {
	return d.data
}

func (d *Context) GetValue(key string) interface{} {
	d.RLock()
	defer d.RUnlock()
	for _, item := range d.temps {
		if item.key == key {
			return item.value
		}
	}
	return nil
}

func (d *Context) SetValue(key string, val interface{}) {
	d.Lock()
	defer d.Unlock()
	for _, item := range d.temps {
		if item.key == key {
			item.value = val
			return
		}
	}
	d.temps = append(d.temps, &KValue{key, val})
}
