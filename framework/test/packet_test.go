package test

import (
	"testing"
	"universal/framework/internal/packet"
)

func TestPacket(t *testing.T) {
	head := &packet.Header{
		SrcNodeType: 1,
		SrcNodeId:   2,
		DstNodeType: 3,
		DstNodeId:   4,
		Uid:         5,
		RouteId:     6,
		Cmd:         7,
		ActorName:   "test_actor",
		FuncName:    "test_func",
	}
	body := []byte("test_body")

	pack := packet.NewPacket(head, body)
	buf := pack.ToBytes()

	pack2 := packet.ParsePacket(buf)
	t.Log(pack2.GetHeader())
}

func BenchmarkPacket(b *testing.B) {
	for i := 0; i < b.N; i++ {
		head := &packet.Header{
			SrcNodeType: uint32(i),
			SrcNodeId:   uint32(i),
			DstNodeType: uint32(i),
			DstNodeId:   uint32(i),
			Uid:         uint64(i),
			RouteId:     uint64(i),
			Cmd:         uint32(i),
			ActorName:   "test_actor",
			FuncName:    "test_func",
		}
		body := []byte("test_body")
		pack := packet.NewPacket(head, body)
		buf := pack.ToBytes()
		packet.ParsePacket(buf)
	}
	b.Log(b.N)
}
