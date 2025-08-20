package attribute

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"universal/common/pb"
	"universal/framework/define"

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

func (p *Player) Login(h *pb.Head, req proto.Message, rsp proto.Message) error {
	fmt.Println(req, "=====>", rsp, p.name)
	p.name = req.(*pb.LoginReq).Token
	fmt.Println(req, "---->", rsp, p.name)
	return nil
}

func TestHandler(t *testing.T) {
	pl := &Player{}
	attr := NewAttribute(define.HeadTwoFunc(pl.Login), "LoginReq", "LoginRsp")
	f := attr.Call(nil, &pb.Head{}, &pb.LoginReq{Token: "asdfasdf"}, &pb.LoginRsp{})
	f()

	attr2 := NewAttribute(define.HeadTwoFunc(pl.Login), "LoginReq", "LoginRsp")
	f2 := attr2.Call(nil, &pb.Head{}, &pb.LoginReq{Token: "fasdf"}, &pb.LoginRsp{})
	f2()
}
