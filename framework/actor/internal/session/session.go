package session

import (
	"fmt"
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/actor/domain"
	"universal/framework/basic"
)

type Session struct {
	uuid       string
	updateTime int64
	tasks      *basic.Async
	handle     domain.PacketHandle
	actors     map[string]interface{}
}

func NewSession(uuid string, h domain.PacketHandle) *Session {
	return &Session{
		tasks:      basic.NewAsync(),
		handle:     h,
		uuid:       uuid,
		updateTime: time.Now().Unix(),
		actors:     make(map[string]interface{}),
	}
}

func (d *Session) SetActor(name string, data interface{}) error {
	if _, ok := d.actors[name]; ok {
		return basic.NewUError(3, -1, fmt.Errorf("Actor(%s) has already registered", name))
	}
	d.actors[name] = data
	return nil
}

func (d *Session) Stop() {
	d.tasks.Stop()
}

func (d *Session) Start() {
	d.tasks.Start()
}

func (d *Session) IsExpired() bool {
	return atomic.LoadInt64(&d.updateTime)+domain.SessionExpireTime <= time.Now().Unix()
}

func (d *Session) Send(pa *pb.Packet) {
	// 刷新时间
	atomic.StoreInt64(&d.updateTime, time.Now().Unix())
	// 发送任务
	ctx := basic.NewContext(pa.Head, d.actors)
	d.tasks.Push(d.handle(ctx, pa))
}
