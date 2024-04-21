package basic

import (
	"bytes"
	"encoding/gob"
	"hash/crc32"
	"reflect"
	"runtime"
	"strings"
	"universal/common/pb"
	"unsafe"

	"google.golang.org/protobuf/proto"
)

func GetCrc32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

func GetFuncName(h interface{}) string {
	v := reflect.ValueOf(h)
	name := runtime.FuncForPC(v.Pointer()).Name()
	return strings.Split(name, ".")[1]
}

func ToErrorRpcHead(rsp proto.Message, err error) proto.Message {
	if err == nil {
		return rsp
	}
	// 设置RpcHead中的错误码和错误信息
	headPtr := reflect.ValueOf(rsp).Elem().Field(0).Addr().Pointer()
	head := (*pb.RpcHead)(unsafe.Pointer(headPtr))
	switch v := err.(type) {
	case *UError:
		head.ErrMsg = v.GetErrMsg()
		head.Code = v.GetCode()
	default:
		head.ErrMsg = v.Error()
		head.Code = -1
	}
	return rsp
}

func ToErrorPacket(pac *pb.Packet, err error) *pb.Packet {
	if err == nil {
		return pac
	}
	switch v := err.(type) {
	case *UError:
		pac.ErrMsg = v.GetErrMsg()
		pac.Code = v.GetCode()
	case nil:
		break
	default:
		pac.ErrMsg = v.Error()
		pac.Code = -1
	}
	return pac
}

func UnmarhsalClusterNode(buf []byte) (*pb.ClusterNode, error) {
	node := &pb.ClusterNode{}
	if err := proto.Unmarshal(buf, node); err != nil {
		return nil, err
	}
	return node, nil
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

func ToReqPacket(head *pb.PacketHead, params ...interface{}) (*pb.Packet, error) {
	if val, ok := params[0].(proto.Message); ok && val != nil {
		buf, err := proto.Marshal(val)
		if err != nil {
			return nil, NewUError(2, pb.ErrorCode_Marhsal, err)
		}
		return &pb.Packet{Head: head, Buff: buf}, nil
	}
	return &pb.Packet{Head: head, Buff: ToGobBytes(params)}, nil
}

func ToRspPacket(head *pb.PacketHead, err error, params ...interface{}) *pb.Packet {
	if err != nil {
		switch vv := err.(type) {
		case *UError:
			return &pb.Packet{Head: head, Code: vv.GetCode(), ErrMsg: vv.GetErrMsg()}
		default:
			return &pb.Packet{Head: head, Code: -1, ErrMsg: err.Error()}
		}
	}
	if val, ok := params[0].(proto.Message); ok && val != nil {
		buf, err := proto.Marshal(val)
		if err == nil {
			return &pb.Packet{Head: head, Code: int32(pb.ErrorCode_Marhsal), ErrMsg: err.Error()}
		}
		return &pb.Packet{Head: head, Buff: buf}
	}
	return &pb.Packet{Head: head, Buff: ToGobBytes(params)}
}
