package base

import (
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/actor/domain"
	"universal/framework/common/async"
	"universal/framework/common/fbasic"
)

type Actor struct {
	*async.Async
	uid        string
	updateTime int64
	handle     domain.ActorHandle
	objects    map[string]interface{}
}

func NewActor(uid string, h domain.ActorHandle) *Actor {
	return &Actor{
		Async:      async.NewAsync(),
		uid:        uid,
		updateTime: time.Now().Unix(),
		handle:     h,
		objects:    make(map[string]interface{}),
	}
}

func (d *Actor) GetUID() string {
	return d.uid
}

func (d *Actor) Send(head *pb.PacketHead, buf []byte) {
	ctx := fbasic.NewContext(head)
	ctx.SetReadOnly(d.objects)
	// 发送任务
	d.Push(d.handle(ctx, buf))
}

func (d *Actor) SetObject(name string, data interface{}) interface{} {
	if _, ok := d.objects[name]; !ok {
		d.objects[name] = data
		return nil
	}
	return d.objects[name]
}

func (d *Actor) GetUpdateTime() int64 {
	return atomic.LoadInt64(&d.updateTime)
}

func (d *Actor) SetUpdateTime(up int64) {
	atomic.StoreInt64(&d.updateTime, up)
}
