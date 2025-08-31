package test

import (
	"fmt"
	"testing"
	"time"
	"universal/common/pb"
	"universal/framework/actor"
	"universal/framework/handler"
)

type Player struct {
	actor.Actor
}

func (p *Player) Init() {
	p.Actor.Register(p)
	p.Actor.Start()
	actor.Register(p)
}

func (p *Player) Kick(head *pb.Head) error {

	return nil
}

func (p *Player) Print(head *pb.Head, val uint32, str string) error {
	fmt.Println(val, "=====>", str)
	return nil
}

func (p *Player) Heart(h *pb.Head, req *pb.HeartReq) error {
	fmt.Println("---->", req)
	return nil
}

func (p *Player) Login(h *pb.Head, req *pb.LoginReq, rsp *pb.LoginRsp) error {
	fmt.Println(req, "=====>", rsp)
	return nil
}

func init() {
	actor.Init(&pb.Node{Type: pb.NodeType_Db}, nil)
	handler.RegisterGob2[Player, uint32, string](pb.NodeType_Db, "Player.Print", (*Player).Print)
	handler.RegisterEvent[Player, pb.HeartReq](pb.NodeType_Db, "Player.Heart", (*Player).Heart)
	handler.RegisterTrigger[Player](pb.NodeType_Db, "Player.Kick", (*Player).Kick)
	handler.RegisterCmd[Player, pb.LoginReq, pb.LoginRsp](pb.NodeType_Db, "Player.Login", (*Player).Login)
}

func TestHandler(t *testing.T) {
	pl := &Player{}
	pl.Init()

	pl.SendMsg(&pb.Head{FuncName: "Heart"}, &pb.HeartReq{BeginTime: time.Now().Unix()})
	pl.SendMsg(&pb.Head{FuncName: "Print"}, uint32(3421), "hchtest")
	pl.Stop()
}
