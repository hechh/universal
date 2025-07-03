package funcs

import (
	"reflect"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/library/uerror"
	"universal/library/util"

	"github.com/golang/protobuf/proto"
)

const (
	HEAD_FLAG          = 1 << 0
	REQ_FLAG           = 1 << 1
	RSP_FLAG           = 1 << 2
	BYTES_FLAG         = 1 << 3
	CMD_HANDLER        = HEAD_FLAG | REQ_FLAG | RSP_FLAG // *pb.head, proto.Message, domain.IRspProto
	NOTIFY_HANDLER     = HEAD_FLAG | REQ_FLAG            // *pb.Head, proto.Message
	HEAD_BYTES_HANDLER = HEAD_FLAG | BYTES_FLAG          // *pb.Head, []byte
	BYTES_HANDLER      = BYTES_FLAG                      // []byte
)

var (
	headType  = reflect.TypeOf((*pb.Head)(nil))
	reqType   = reflect.TypeOf((*proto.Message)(nil)).Elem()
	rspType   = reflect.TypeOf((*domain.IRspProto)(nil)).Elem()
	bytesType = reflect.TypeOf((*[]byte)(nil)).Elem()
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	nilValue  = reflect.ValueOf((*error)(nil))
	args      = util.ArrayPool[reflect.Value](6)
	sendRsp   func(*pb.Head, proto.Message) error
)

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
