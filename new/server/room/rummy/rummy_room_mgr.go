package rummy

import (
	"poker_server/framework"
	"reflect"
)

type RummyRoomMgr struct {
	framework.Actor
	mgr *framework.ActorMgr
}

func NewRummyRoomMgr() *RummyRoomMgr {
	// 预先注册
	mgr := new(framework.ActorMgr)
	room := &RummyRoom{}
	mgr.Register(room)
	mgr.ParseFunc(reflect.TypeOf(room))
	framework.RegisterActor(mgr)

	// 创建房间actor 管理器
	ret := &RummyRoomMgr{mgr: mgr}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Actor.Start()
	framework.RegisterActor(ret)

	ret.loop()
	return ret
}

func (this *RummyRoomMgr) loop() {}
