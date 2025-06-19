package test

import (
	"testing"
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

func TestRouter(t *testing.T) {
	rr := &pb.Router{Gate: 1<<16 | 1, Room: 123543, Db: int32(13123<<16) | 100}
	buf, _ := proto.Marshal(rr)
	t.Log(len(buf))
}
