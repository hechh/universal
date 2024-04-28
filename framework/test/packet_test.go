package test

import (
	"fmt"
	"reflect"
	"testing"
	"universal/framework/fbasic"
	"universal/framework/packet"

	"universal/common/pb"

	"google.golang.org/protobuf/proto"
)

func Print(name string, val int32) {
	fmt.Println(name, "----->", val)
}

func Login(ctx *fbasic.Context, req, rsp proto.Message) error {
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
	fmt.Println("=====>", d.name, d.age)
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
	data := map[string]fbasic.IData{
		"Person": &Person{"hch", 120},
	}
	t.Run("Person", func(t *testing.T) {
		head := &pb.PacketHead{UID: 1234000, ApiCode: 2}
		ctx := fbasic.NewContext(head, data)
		req := &pb.ActorRequest{ActorName: "Person", FuncName: "GetAge"}
		buf, _ := proto.Marshal(req)
		ret, err := packet.Call(ctx, buf)
		rsp := &pb.ActorResponse{}
		proto.Unmarshal(ret.Buff, rsp)
		t.Log(string(rsp.Buff), "-----Actor Result------", err)
	})
	t.Run("Login", func(t *testing.T) {
		head := &pb.PacketHead{UID: 1234000, ApiCode: 1}
		ctx := fbasic.NewContext(head, data)
		req := &pb.LoginRequest{}
		buf, _ := proto.Marshal(req)
		ret, err := packet.Call(ctx, buf)
		t.Log(ret, "-----Login Result------", err)
	})
	t.Run("Func", func(t *testing.T) {
		head := &pb.PacketHead{UID: 1234000, ApiCode: 2}
		ctx := fbasic.NewContext(head, data)
		req := &pb.ActorRequest{ActorName: "Person", FuncName: "GetName"}
		buf, _ := proto.Marshal(req)
		ret, err := packet.Call(ctx, buf)
		t.Log(ret, "-----Func Result------", err)
	})
}

func ToTurn(a []byte, b []reflect.Value) {
	for i := 0; i < len(b); i++ {
		b[i] = reflect.ValueOf(string(a))
	}
}

func TestPrint(t *testing.T) {
	b := make([]reflect.Value, 10)
	ToTurn([]byte("hch"), b)
	t.Log(b, cap(b))

	vv := []reflect.Value{
		reflect.ValueOf(12),
		reflect.ValueOf(12),
		reflect.ValueOf("test"),
		reflect.ValueOf(12),
	}

	t.Log(fbasic.EncodeValues(vv).Encode())

	/*
		bb := bytes.NewBuffer(nil)
		enc := gob.NewEncoder(bb)
		for _, param := range []interface{}{123, 1, 2, 4, "aserfa"} {
			enc.Encode(param)
		}
		buf := bb.Bytes()
		bb.Reset()
		fmt.Println("q----->", len(buf), cap(buf))
	*/
}
