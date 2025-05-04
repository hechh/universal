package test

import (
	"reflect"
	"testing"
	"universal/framework"
	"universal/framework/internal/packet"
)

type Player struct {
	framework.Actor
	count int64
}

type PlayerMgr struct {
	framework.ActorGroup
}

func (d *Player) Print() {
	//fmt.Println("-----------------")
	d.count++
}

func TestPlayerMgr(t *testing.T) {
	usr := new(Player)
	usr.Register(usr, nil)
	usr.Start()

	mgr := new(PlayerMgr)
	mgr.Register(&Player{}, nil)
	mgr.AddActor(1, usr)

	head := &packet.Header{
		RouteId:   1,
		Uid:       1,
		ActorName: "Player",
		FuncName:  "Print",
	}
	if err := mgr.Send(head, nil); err != nil {
		t.Log("=====>", err)
	}
	usr.Stop()
}

func TestActor(t *testing.T) {
	usr := new(Player)
	usr.Register(usr, reflect.TypeOf(usr))
	usr.Start()
	head := &packet.Header{
		ActorName: "Player",
		FuncName:  "Print",
	}
	if err := usr.Send(head, nil); err != nil {
		t.Log("=====>", err)
	}
	usr.Stop()
}

func BenchmarkActor(b *testing.B) {
	usr := new(Player)
	usr.Register(usr, reflect.TypeOf(usr))
	usr.Start()
	for i := 0; i < b.N; i++ {
		head := &packet.Header{
			ActorName: "Player",
			FuncName:  "Print",
		}
		if err := usr.Send(head, nil); err != nil {
			b.Log("=====>", err)
		}
	}
	usr.Stop()
	b.Log(usr.count)
}
