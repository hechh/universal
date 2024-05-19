package domain

import (
	"universal/framework/common/fbasic"

	"google.golang.org/protobuf/proto"
)

// 对外接口定义
type ApiFunc func(*fbasic.Context, proto.Message, proto.Message) error

type IHandler interface {
	Call(*fbasic.Context, []byte) proto.Message
}
