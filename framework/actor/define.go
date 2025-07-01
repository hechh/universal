package actor

import (
	"reflect"
	"strings"
	"universal/common/pb"
	"universal/framework/domain"

	"github.com/golang/protobuf/proto"
)

const (
	RSP_FLAG   = 1 << 0
	REQ_FLAG   = 1 << 1
	HEAD_FLAG  = 1 << 2
	BYTES_FLAG = 1 << 3

	CMD_HANDLER    = HEAD_FLAG | REQ_FLAG | RSP_FLAG // *pb.head, proto.Message, proto.Message
	NOTIFY_HANDLER = HEAD_FLAG | REQ_FLAG            // *pb.Head, proto.Message
	BYTES_HANDLER  = HEAD_FLAG | BYTES_FLAG          // *pb.Head, []byte
)

var (
	rspType   = reflect.TypeOf((*domain.IRspProto)(nil)).Elem()
	reqType   = reflect.TypeOf((*proto.Message)(nil)).Elem()
	headType  = reflect.TypeOf((*pb.Head)(nil))
	bytesType = reflect.TypeOf((*[]byte)(nil)).Elem()
)

func parseName(rr reflect.Type) string {
	name := rr.String()
	if index := strings.Index(name, "."); index > -1 {
		name = name[index+1:]
	}
	return name
}

type Method struct {
	reflect.Method
	ins  int
	flag uint32
}

func parseMethod(m reflect.Method) *Method {
	ins, outs := m.Type.NumIn(), m.Type.NumOut()
	if outs > 1 {
		return nil
	}

	return nil
}
