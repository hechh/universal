package test

import (
	"fmt"
	"testing"
	"time"
	"universal/common/pb"
	"universal/framework/actor"
	"universal/framework/handler"
	"universal/library/encode"
)

type Player struct {
	actor.Actor
}

func (p *Player) Init() {
	p.Register(p)
	p.Start()
	actor.Register(p)
}

func (p *Player) Print(val uint32, str string) error {
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
	handler.RegisterCmd[Player, pb.LoginReq, pb.LoginRsp](pb.NodeType_Db, "Player.Login", (*Player).Login)
}

func TestHandler(t *testing.T) {
	pl := &Player{}
	pl.Init()

	pl.SendMsg(&pb.Head{FuncName: "Heart"}, &pb.HeartReq{BeginTime: time.Now().Unix()})
	pl.SendMsg(&pb.Head{FuncName: "Print"}, uint32(3421), "hchtest")
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

func TestGob(t *testing.T) {
	arg1 := uint32(100)
	buf, err := encode.Encode(arg1)
	if err != nil {
		t.Log("=====1=====>", err)
		return
	}

	param1 := new(uint32)
	if err := encode.Decode(buf, param1); err != nil {
		t.Log("=====2=====>", err)
		return
	}
	t.Log(*param1)
}
