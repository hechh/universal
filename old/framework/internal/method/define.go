package method

import (
	"reflect"
	"universal/common/pb"
	"universal/framework/define"

	"github.com/golang/protobuf/proto"
)

const (
	HEAD_FLAG           = 1 << 0
	REQ_FLAG            = 1 << 1
	RSP_FLAG            = 1 << 2
	BYTES_FLAG          = 1 << 3
	INTERFACE_FLAG      = 1 << 4
	GOB_FLAG            = 1 << 5
	HEAD_REQ_RSP_MASK   = HEAD_FLAG | REQ_FLAG | RSP_FLAG
	HEAD_REQ_MASK       = HEAD_FLAG | REQ_FLAG
	HEAD_RSP_MASK       = HEAD_FLAG | RSP_FLAG
	REQ_RSP_MASK        = REQ_FLAG | RSP_FLAG
	HEAD_BYTES_MASK     = HEAD_FLAG | BYTES_FLAG
	HEAD_INTERFACE_MASK = HEAD_FLAG | INTERFACE_FLAG
)

var (
	headType      = reflect.TypeOf((*pb.Head)(nil))
	reqType       = reflect.TypeOf((*proto.Message)(nil)).Elem()
	rspType       = reflect.TypeOf((*define.IRspProto)(nil)).Elem()
	bytesType     = reflect.TypeOf((*[]byte)(nil)).Elem()
	errorType     = reflect.TypeOf((*error)(nil)).Elem()
	interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
	nilError      = reflect.ValueOf((*error)(nil)).Elem()
)
