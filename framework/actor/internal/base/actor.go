package base

import (
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/common/uerrors"
	"universal/framework/actor/domain"
	"universal/framework/fbasic"
)

type Actor struct {
	uuid       string
	updateTime int64
	tasks      *fbasic.Async
	handle     domain.ActorHandle
	objects    map[string]fbasic.IData
}

func NewActor(uuid string, h domain.ActorHandle) *Actor {
	return &Actor{
		tasks:      fbasic.NewAsync(),
		handle:     h,
		uuid:       uuid,
		updateTime: time.Now().Unix(),
		objects:    make(map[string]fbasic.IData),
	}
}

func (d *Actor) UUID() string {
	return d.uuid
}

func (d *Actor) Stop() {
	d.tasks.Stop()
}

func (d *Actor) Start() {
	d.tasks.Start()
}

func (d *Actor) Send(pa *pb.Packet) {
	// 发送任务
	ctx := fbasic.NewContext(pa.Head, d.objects)
	d.tasks.Push(d.handle(ctx, pa.Buff))
}

func (d *Actor) SetObject(name string, data fbasic.IData) error {
	if _, ok := d.objects[name]; ok {
		return uerrors.ActorHasRegistered(name)
	}
	d.objects[name] = data
	return nil
}

func (d *Actor) GetUpdateTime() int64 {
	return atomic.LoadInt64(&d.updateTime)
}

func (d *Actor) SetUpdateTime(up int64) {
	atomic.StoreInt64(&d.updateTime, up)
}
