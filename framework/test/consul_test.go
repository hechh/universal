package test

import (
	"testing"
	"time"
	"universal/framework/internal/cluster"
	"universal/framework/internal/discovery"
)

func TestConsul(t *testing.T) {
	consul, err := discovery.NewConsul(
		cfg.Consul.Endpoints,
		discovery.WithPath("hch/consul_test"),
		discovery.WithParse(cluster.NewNode),
	)
	if err != nil {
		t.Log(err)
		return
	}

	// 监视
	self := &cluster.Node{Name: "test2", Type: 1, Id: 1, Addr: "192.168.1.1:22345"}
	cls := cluster.NewCluster(self)
	if err := consul.Watch(cls); err != nil {
		t.Fatal(err)
	}

	// 添加服务
	nnode := &cluster.Node{Name: "test2", Type: 2, Id: 2, Addr: "192.168.1.1:22345"}
	if err := consul.Put(nnode); err != nil {
		t.Log(err)
		return
	}

	// 删除节点
	if err := consul.Del(nnode); err != nil {
		t.Log(err)
		return
	}

	consul.Close()
}

func TestConsulKeepAlive(t *testing.T) {
	consul, err := discovery.NewConsul(
		cfg.Consul.Endpoints,
		discovery.WithPath("hch/consul_test"),
		discovery.WithParse(cluster.NewNode),
	)
	if err != nil {
		t.Log(err)
		return
	}

	// 监视
	self := &cluster.Node{Name: "test2", Type: 1, Id: 1, Addr: "192.168.1.1:22345"}
	cls := cluster.NewCluster(self)
	if err := consul.Watch(cls); err != nil {
		t.Fatal(err)
	}

	if err := consul.KeepAlive(self, 10); err != nil {
		t.Fatal(err)
	}

	time.Sleep(20 * time.Second)
	consul.Close()
}
