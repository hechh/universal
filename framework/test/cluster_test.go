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
		err := cls.Put(&cluster.Node{Name: "test1", Type: uint32(i%int(define.NodeTypeMax-1)) + 1, Id: uint32(i) + 1, Addr: "192.168.1.1:22345"})
		if err != nil {
			b.Log(err)
			return
		}
	}
	for i := 0; i < b.N; i++ {
		node := cls.Get(uint32(i%int(define.NodeTypeMax-1))+1, uint32(i))
		if node != nil {
			rtr.Update(uint64(node.GetId()), &define.RouteInfo{})
		}
		if node != nil && node.GetId() == 0 {
			cls.Del(node.GetType(), node.GetId())
		}
	}
	b.Log(b.N)
}
