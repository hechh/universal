package test

import (
	"fmt"
	"testing"
	"universal/framework/fbasic"
	"universal/framework/packet"
	"universal/framework/packet/internal/manager"

	"universal/common/pb"

	"google.golang.org/protobuf/proto"
)

func LoginRequest(ctx *fbasic.Context, req, rsp proto.Message) error {
	resp := rsp.(*pb.GateLoginResponse)
	resp.Head.Code = 10
	resp.Head.ErrMsg = "this is a test"
	return nil
}

func Print(name string, val int32) {
	fmt.Println("---Print--->", name, val)
}

type Person struct {
	name string
	age  int32
}

func (d *Person) SetName(name string) {
	d.name = name
	fmt.Println("------SetName------->", d.name, d.age)
}

func (d *Person) GetName() string {
	fmt.Println("------GetName------->", d.name, d.age)
	return d.name
}

func (d *Person) SetAge(age int32) {
	d.age = age
	fmt.Println("------SetAge------->", d.name, d.age)
}

func (d *Person) GetAge() int32 {
	fmt.Println("------GetAge------->", d.name, d.age)
	return d.age
}

func TestMain(m *testing.M) {
	packet.RegisterApi(1, LoginRequest, &pb.GateLoginRequest{}, &pb.GateLoginResponse{})
	packet.RegisterFunc(2, Print)
	packet.RegisterStruct(2, &Person{})
	m.Run()
}

func TestApi(t *testing.T) {
	data := map[string]fbasic.IData{
		"Person": &Person{"hch", 120},
	}
	t.Run("LoginRequest调用测试", func(t *testing.T) {
		head := &pb.PacketHead{UID: 1234000, ApiCode: 1}
		ctx := fbasic.NewContext(head, data)
		req := &pb.GateLoginRequest{}
		buf, _ := proto.Marshal(req)
		ret, err := packet.Call(ctx, buf)
		rsp := &pb.GateLoginResponse{}
		proto.Unmarshal(ret.Buff, rsp)
		t.Log("-----LoginRequest Result------", rsp, err)
	})
	t.Run("Print调用测试", func(t *testing.T) {
		head := &pb.PacketHead{UID: 1234000, ApiCode: 2}
		ctx := fbasic.NewContext(head, data)
		// 设置参数
		req := &pb.ActorRequest{FuncName: "Print", Buff: fbasic.AnyToEncode("print", 423)}
		buf, _ := proto.Marshal(req)
		ret, err := packet.Call(ctx, buf)
		rsp := &pb.ActorResponse{}
		proto.Unmarshal(ret.Buff, rsp)
		t.Log("------Print Result------", rsp, err)
	})
	t.Run("Person.GetAge调用测试", func(t *testing.T) {
		head := &pb.PacketHead{UID: 1234000, ApiCode: 2}
		ctx := fbasic.NewContext(head, map[string]fbasic.IData{"Person": &Person{"hch10", 10}})
		req := &pb.ActorRequest{ActorName: "Person", FuncName: "GetAge"}
		buf, _ := proto.Marshal(req)
		// 返回值
		ret, err := packet.Call(ctx, buf)
		rsp := &pb.ActorResponse{}
		proto.Unmarshal(ret.Buff, rsp)
		rets, err01 := manager.ParseReturns(2, "Person", "GetAge", rsp.Buff)
		t.Log(ret.Head, "-----GetAge Result------", err, err01, rets, rsp)
	})
	t.Run("Person.SetAge调用测试", func(t *testing.T) {
		head := &pb.PacketHead{UID: 1234000, ApiCode: 2}
		pp := &Person{"hch10", 10}
		ctx := fbasic.NewContext(head, map[string]fbasic.IData{"Person": pp})
		req := &pb.ActorRequest{ActorName: "Person", FuncName: "SetAge", Buff: fbasic.AnyToEncode(120)}
		buf, _ := proto.Marshal(req)
		// 返回值
		ret, err := packet.Call(ctx, buf)
		rsp := &pb.ActorResponse{}
		proto.Unmarshal(ret.Buff, rsp)
		rets, err01 := manager.ParseReturns(2, "Person", "SetAge", rsp.Buff)
		t.Log(ret.Head, pp, "-----SetAge Result------", err, err01, rets, rsp)
	})
}
