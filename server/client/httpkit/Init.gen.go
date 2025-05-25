/*
* 本代码由cfgtool工具生成，请勿手动修改
 */

package httpkit

import (
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

var (
	cmds = make(map[uint32]func() proto.Message)
)

func init() {
	cmds[16777216] = func() proto.Message { return &pb.LoginRequest{} }
	cmds[16777217] = func() proto.Message { return &pb.LoginResponse{} }
	cmds[16777218] = func() proto.Message { return &pb.HeartRequest{} }
	cmds[16777219] = func() proto.Message { return &pb.HeartResponse{} }
}
