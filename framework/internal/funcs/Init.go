package funcs

import (
	"reflect"
	"universal/common/pb"
	"universal/library/uerror"
	"universal/library/util"

	"github.com/golang/protobuf/proto"
)

var (
	sendRsp func(*pb.Head, proto.Message) error
	self    *pb.Node
	args    = util.ArrayPool[reflect.Value](6)
)

func Init(nn *pb.Node, f func(*pb.Head, proto.Message) error) {
	sendRsp = f
	self = nn
}

func get(size int) []reflect.Value {
	rets := args.Get().([]reflect.Value)
	return rets[:size]
}

func put(rets []reflect.Value) {
	args.Put(rets)
}

func toRspHead(err error) *pb.RspHead {
	switch vv := err.(type) {
	case nil:
		return nil
	case *uerror.UError:
		return &pb.RspHead{Code: vv.GetCode(), ErrMsg: vv.GetMsg()}
	case error:
		return &pb.RspHead{Code: -1, ErrMsg: vv.Error()}
	default:
		return nil
	}
}
