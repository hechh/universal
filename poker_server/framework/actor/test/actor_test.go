package test

import (
	"fmt"
	"poker_server/common/pb"
	"poker_server/framework/actor"
	"reflect"
	"testing"
)

type Player struct {
	actor.Actor
}

type PlayerMgr struct {
	actor.Actor
	mgr *actor.ActorMgr
}

func NewPlayerMgr() *PlayerMgr {
	mgr := new(actor.ActorMgr)
	pp := &Player{}
	mgr.Register(pp)
	mgr.ParseFunc(reflect.TypeOf(pp))
	actor.Register(mgr)

	ret := &PlayerMgr{mgr: mgr}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Start()
	actor.Register(ret)
	return ret
}

func (d *PlayerMgr) Stop() {
	d.mgr.Stop()
	d.Actor.Stop()
}

func (d *PlayerMgr) Heart(head *pb.Head, req *pb.GateHeartRequest, rsp *pb.GateHeartResponse) {
	fmt.Println("---------->", head)
}

func (d *PlayerMgr) HeartRequest(head *pb.Head, req *pb.GateHeartRequest, rsp *pb.GateHeartResponse) error {
	fmt.Println("---------->", head)
	return nil
}

func TestActor(t *testing.T) {
	mgr := NewPlayerMgr()

	actor.SendMsg(&pb.Head{
		ActorName: "PlayerMgr",
		FuncName:  "HeartRequest",
	}, &pb.GateHeartRequest{})

	mgr.Stop()
}

func TestActor2(t *testing.T) {
	mgr := NewPlayerMgr()

	actor.SendMsg(&pb.Head{
		ActorName: "PlayerMgr",
		FuncName:  "Heart",
	}, &pb.GateHeartRequest{})

	mgr.Stop()
}
