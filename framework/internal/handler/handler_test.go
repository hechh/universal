package handler

import (
	"fmt"
	"testing"
	"time"
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

type Player struct {
	name string
}

func (p *Player) Heart(h *pb.Head, req *pb.HeartReq) error {
	fmt.Println("---->", req)
	return nil
}

func (p *Player) Login(h *pb.Head, req *pb.LoginReq, rsp *pb.LoginRsp) error {
	fmt.Println(req, "=====>", rsp, p.name)
	p.name = req.Token
	fmt.Println(req, "---->", rsp, p.name)
	return nil
}

func TestHandler(t *testing.T) {
	Register1[Player, pb.HeartReq](pb.NodeType_Db, "Player.Heart", (*Player).Heart)
	Register2[Player, pb.LoginReq, pb.LoginRsp](pb.NodeType_Db, "Player.Login", (*Player).Login)

	pl := &Player{name: "hhh"}

	ff := GetHandler(pb.NodeType_Db, "Player", "Login")
	ff.Call(nil, pl, &pb.Head{}, &pb.LoginReq{Token: "asdfasdf"}, &pb.LoginRsp{})()

	req := &pb.HeartReq{BeginTime: time.Now().Unix()}
	buf, _ := proto.Marshal(req)
	rf := GetHandler(pb.NodeType_Db, "Player", "Heart")
	rf.Rpc(nil, pl, &pb.Head{}, buf)()
}
