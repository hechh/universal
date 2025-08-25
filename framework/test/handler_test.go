package test

import (
	"fmt"
	"testing"
	"time"
	"universal/common/pb"
	"universal/framework/actor"
	"universal/framework/internal/handler"
)

type Player struct {
	actor.Actor
}

func (p *Player) Init() {
	p.Register(p)
	p.Start()
	actor.Register(p)
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
	handler.Register1[Player, pb.HeartReq](pb.NodeType_Db, "Player.Heart", (*Player).Heart)
	handler.Register2[Player, pb.LoginReq, pb.LoginRsp](pb.NodeType_Db, "Player.Login", (*Player).Login)
}

func TestHandler(t *testing.T) {
	pl := &Player{}
	pl.Init()

	pl.SendMsg(&pb.Head{FuncName: "Heart"}, &pb.HeartReq{BeginTime: time.Now().Unix()})
	pl.Stop()

	/*
		ff := handler.GetHandler(pb.NodeType_Db, "Player", "Login")
		ff.Call(nil, pl, &pb.Head{}, &pb.LoginReq{Token: "asdfasdf"}, &pb.LoginRsp{})()

		req := &pb.HeartReq{BeginTime: time.Now().Unix()}
		buf, _ := proto.Marshal(req)
		rf := handler.GetHandler(pb.NodeType_Db, "Player", "Heart")
		rf.Rpc(nil, pl, &pb.Head{}, buf)()
	*/
}
