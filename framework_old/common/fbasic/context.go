package fbasic

import (
	"sync"
	"universal/common/pb"
)

type IRpcHead interface {
	GetHead() *pb.RpcHead
}

type IContext interface {
	GetTemp(string) interface{}
	SetTemp(string, interface{})
	GetValue(string) interface{}
	SetReadOnly(map[string]interface{})
}

type KValue struct {
	key   string
	value interface{}
}

type Context struct {
	sync.RWMutex                          // 读写锁
	*pb.PacketHead                        // rpc请求头
	readOnly       map[string]interface{} // 只读数据
	temps          []*KValue              // 临时缓存
}

func NewContext(head *pb.PacketHead) *Context {
	return &Context{PacketHead: head}
}

func (d *Context) SetReadOnly(data map[string]interface{}) {
	d.readOnly = data
}

func (d *Context) GetValue(key string) interface{} {
	if d.readOnly == nil {
		return nil
	}
	return d.readOnly[key]
}

func (d *Context) GetTemp(key string) interface{} {
	d.RLock()
	defer d.RUnlock()
	for _, item := range d.temps {
		if item.key == key {
			return item
		}
	}
	return nil
}

func (d *Context) SetTemp(key string, val interface{}) {
	d.Lock()
	defer d.Unlock()
	for _, item := range d.temps {
		if item.key == key {
			item.value = val
			return
		}
	}
	d.temps = append(d.temps, &KValue{key: key, value: val})
}
