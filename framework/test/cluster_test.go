package test

import (
	"testing"
	"universal/framework/define"
	"universal/framework/internal/cluster"
	"universal/framework/internal/router"
)

func BenchmarkCluster(b *testing.B) {
	self := &cluster.Node{Name: "test1", Type: 1, Id: 1, Addr: "192.168.1.1:22345"}
	cls := cluster.NewCluster(self)
	rtr := router.NewRouter()

	for i := 0; i < 60; i++ {
		err := cls.Put(&cluster.Node{Name: "test1", Type: int32(i%int(define.NodeTypeMax-1)) + 1, Id: int32(i) + 1, Addr: "192.168.1.1:22345"})
		if err != nil {
			b.Log(err)
			return
		}
	}
	for i := 0; i < b.N; i++ {
		node := cls.Get(int32(i%int(define.NodeTypeMax-1))+1, int32(i))
		if node != nil {
			rtr.Update(uint64(node.GetId()), node)
		}
		if node != nil && node.GetId() == 0 {
			cls.Del(node.GetType(), node.GetId())
		}
	}
	b.Log(b.N)
}
