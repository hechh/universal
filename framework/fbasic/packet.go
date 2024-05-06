package fbasic

import (
	"encoding/json"
	"hash/crc32"
	"reflect"
	"runtime"
	"strings"
	"universal/common/pb"

	"google.golang.org/protobuf/proto"
)

func GetCrc32(str interface{}) uint32 {
	var buf []byte
	switch vv := str.(type) {
	case string:
		buf = []byte(vv)
	case []byte:
		buf = vv
	case proto.Message:
		buf, _ = proto.Marshal(vv)
	default:
		buf, _ = json.Marshal(str)
	}
	return crc32.ChecksumIEEE(buf)
}

func GetFuncName(h interface{}) string {
	var name string
	switch vv := h.(type) {
	case reflect.Value:
		name = runtime.FuncForPC(vv.Pointer()).Name()
	default:
		v := reflect.ValueOf(vv)
		name = runtime.FuncForPC(v.Pointer()).Name()
	}
	return strings.Split(name, ".")[1]
}

func ReqToPacket(head *pb.PacketHead, req proto.Message, params ...interface{}) (*pb.Packet, error) {
	// 设置参数
	if len(params) > 0 {
		switch vv := req.(type) {
		case *pb.ActorRequest:
			vv.Buff = AnyToEncode(params...)
		}
	}
	// 封装
	buf, err := proto.Marshal(req)
	if err != nil {
		return nil, NewUError(1, pb.ErrorCode_ProtoMarshal, err)
	}
	return &pb.Packet{Head: head, Buff: buf}, nil
}

func RspToPacket(head *pb.PacketHead, rsp proto.Message, params ...interface{}) (*pb.Packet, error) {
	// 设置参数
	if len(params) > 0 {
		switch vv := rsp.(type) {
		case *pb.ActorResponse:
			vv.Buff = AnyToEncode(params...)
		}
	}
	// 序列化
	buf, err := proto.Marshal(rsp)
	if err != nil {
		return nil, NewUError(1, pb.ErrorCode_ProtoMarshal, err)
	}
	return &pb.Packet{Head: head, Buff: buf}, nil
}

func NewActorRequest(actorName, funcName string, params ...interface{}) proto.Message {
	return &pb.ActorRequest{
		ActorName: actorName,
		FuncName:  funcName,
		Buff:      AnyToEncode(params...),
	}
}
