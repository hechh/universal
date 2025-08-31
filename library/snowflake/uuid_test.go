package snowflake

import (
	"testing"
	"universal/common/pb"
)

func TestUUID(t *testing.T) {
	Init(&pb.Node{Type: 10, Id: 10})

	t.Log(GenUUID())
	t.Log(GenUUID())
	t.Log(GenUUID())
}
