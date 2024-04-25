package test

import (
	"fmt"
	"testing"
	"universal/framework/basic"
	"universal/framework/packet"

	"universal/common/pb"

	"google.golang.org/protobuf/proto"
)

func Print(name string, val int32) {
	fmt.Println(name, "----->", val)
}

func Login(ctx *basic.Context, req, rsp proto.Message) error {
	fmt.Println(ctx, "--Login-->", req, rsp)
	return nil
}

type Person struct {
	name string
	age  int32
}

func (d *Person) SetName(name string) {
	d.name = name
}

func (d *Person) GetName() string {
	return d.name
}

func (d *Person) SetAge(age int32) {
	d.age = age
}

func (d *Person) GetAge() int32 {
	return d.age
}

func (d *Person) Print() {
	fmt.Println(d.name, "-------->", d.age)
}

func TestMain(m *testing.M) {
	packet.RegisterApi(1, Login, &pb.LoginRequest{}, &pb.LoginResponse{})
	packet.RegisterFunc(2, Print)
	packet.RegisterStruct(2, &Person{})
	m.Run()
}

func TestApi(t *testing.T) {
	data := map[string]interface{}{
		"Person": &Person{"hch", 120},
	}
	t.Run("Login", func(t *testing.T) {
		head := &pb.PacketHead{UID: 1234000, ApiCode: 1}
		ctx := basic.NewContext(head, data)
		req := &pb.LoginRequest{}
		buf, _ := proto.Marshal(req)
		ret, err := packet.Call(ctx, buf)
		fmt.Println(ret, "-----------", err)
	})
	t.Run("Func", func(t *testing.T) {
		head := &pb.PacketHead{UID: 1234000, ApiCode: 2}
		ctx := basic.NewContext(head, data)
		req := &pb.ActorRequest{ActorName: "Person", FuncName: "GetName"}
		buf, _ := proto.Marshal(req)
		ret, err := packet.Call(ctx, buf)
		fmt.Println(ret, "-----------", err)
	})
}
