package fbasic

import (
	"bytes"
	"encoding/gob"
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

func ToGobBytes(params interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	switch vals := params.(type) {
	case []interface{}:
		for _, param := range vals {
			enc.Encode(param)
		}
	case []reflect.Value:
		for _, param := range vals {
			enc.EncodeValue(param)
		}
	}
	return buf.Bytes()
}

func toErrorRsp(err error, rsp proto.Message) {
	code, errMsg := GetCodeMsg(err)
	vv := reflect.ValueOf(rsp).Elem().Field(3)
	if vv.IsNil() {
		vv.Set(reflect.ValueOf(&pb.RpcHead{Code: code, ErrMsg: errMsg}))
	} else if head, ok := vv.Interface().(*pb.RpcHead); ok {
		head.Code = code
		head.ErrMsg = errMsg
	}
}

func RspToPacket(head *pb.PacketHead, err error, params ...interface{}) (*pb.Packet, error) {
	var rsp proto.Message
	if len(params) <= 0 {
		rsp = &pb.ActorResponse{Head: &pb.RpcHead{}}
	} else {
		if val, ok := params[0].(proto.Message); !ok {
			rsp = &pb.ActorResponse{Head: &pb.RpcHead{}, Buff: ToGobBytes(params)}
		} else {
			rsp = val
		}
	}
	if err != nil {
		toErrorRsp(err, rsp)
	}
	buf, err := proto.Marshal(rsp)
	if err != nil {
		return nil, NewUError(1, pb.ErrorCode_ProtoMarshal, err)
	}
	return &pb.Packet{Head: head, Buff: buf}, nil
}

func ReqToPacket(head *pb.PacketHead, params ...interface{}) (*pb.Packet, error) {
	var req proto.Message
	if len(params) <= 0 {
		req = &pb.ActorRequest{Head: &pb.RpcHead{}}
	} else {
		if val, ok := params[0].(proto.Message); !ok {
			req = &pb.ActorRequest{Head: &pb.RpcHead{}, Buff: ToGobBytes(params)}
		} else {
			req = val
		}
	}
	// 封装
	buf, err := proto.Marshal(req)
	if err != nil {
		return nil, NewUError(1, pb.ErrorCode_ProtoMarshal, err)
	}
	return &pb.Packet{Head: head, Buff: buf}, nil
}
