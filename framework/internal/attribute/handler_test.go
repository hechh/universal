package attribute

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

func parseName(rr reflect.Type) string {
	name := rr.String()
	if index := strings.Index(name, "."); index > -1 {
		name = name[index+1:]
	}
	return name
}

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
	pl := &Player{}
	ll := Handler[pb.LoginReq, pb.LoginRsp](pl.Login)
	ll.Call(nil, &pb.Head{}, &pb.LoginReq{Token: "asdfasdf"}, &pb.LoginRsp{})()

	req := &pb.HeartReq{BeginTime: time.Now().Unix()}
	buf, _ := proto.Marshal(req)
	attr := Notify[pb.HeartReq](pl.Heart)
	attr.Rpc(nil, &pb.Head{}, buf)()
}
