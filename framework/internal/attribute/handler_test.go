package attribute

import (
	"fmt"
	"testing"
	"time"
	"universal/common/pb"

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

func TestMethod(t *testing.T) {
	mm := NewMethod()
	mm.Register("Heart", Notify[Player, pb.HeartReq]((*Player).Heart))
	mm.Register("Login", Handler[Player, pb.LoginReq, pb.LoginRsp]((*Player).Login))

	pl := &Player{name: "hhh"}
	mm.Call(nil, pl, &pb.Head{FuncName: "Heart"}, &pb.HeartReq{BeginTime: time.Now().Unix()})()
}

func TestHandler(t *testing.T) {
	pl := &Player{name: "hhh"}

	ll := Handler[Player, pb.LoginReq, pb.LoginRsp]((*Player).Login)
	ll.Call(nil, pl, &pb.Head{}, &pb.LoginReq{Token: "asdfasdf"}, &pb.LoginRsp{})()

	req := &pb.HeartReq{BeginTime: time.Now().Unix()}
	buf, _ := proto.Marshal(req)
	attr := Notify[Player, pb.HeartReq]((*Player).Heart)
	attr.Rpc(nil, pl, &pb.Head{}, buf)()
}
