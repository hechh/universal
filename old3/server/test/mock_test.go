package test

import (
	"testing"
	"universal/common/pb"
	"universal/framework/mock"
)

func TestMock(t *testing.T) {
	cmd := uint32(pb.CMD_HEART_REQUEST)
	err := mock.Request(cmd, &pb.HeartRequest{})
	t.Log(err)
}
