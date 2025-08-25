package handler

import (
	"fmt"
	"testing"
	"time"
	"universal/common/pb"
	"universal/library/util"

	"google.golang.org/protobuf/proto"
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
	RegisterNotify[Player, pb.HeartReq]("Player.Heart", (*Player).Heart)
	RegisterHandler[Player, pb.LoginReq, pb.LoginRsp]("Player.Login", (*Player).Login)

	pl := &Player{name: "hhh"}

	Call(nil, pl, &pb.Head{ActorFunc: util.String2Int("Player.Login")}, &pb.LoginReq{Token: "asdfasdf"}, &pb.LoginRsp{})()

	req := &pb.HeartReq{BeginTime: time.Now().Unix()}
	buf, _ := proto.Marshal(req)
	Rpc(nil, pl, &pb.Head{ActorFunc: util.String2Int("Player.Heart")}, buf)()
}
