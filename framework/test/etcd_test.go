package test

import (
	"testing"
	"time"
	"universal/framework/config"
	"universal/framework/define"
	"universal/framework/internal/cluster"
	"universal/framework/internal/discovery"
	"universal/framework/internal/router"
)

var (
	cfg  *config.Config
	etcd define.IDiscovery
)

func TestMain(m *testing.M) {
	var err error
	cfg, err = config.NewConfig("../../configure/yaml/local.yaml")
	if err != nil {
		panic(err)
	}
	dis, err := discovery.NewEtcd(
		cfg.Etcd.Endpoints,
		discovery.WithPath("/hch/etcd_test/"),
		discovery.WithParse(cluster.NewNode),
	)
	if err != nil {
		panic(err)
	}
	etcd = dis
	m.Run()
}

func TestEtcdPut(t *testing.T) {
	// 监视
	self := &cluster.Node{Name: "test1", Type: 1, Id: 1, Addr: "192.168.1.1:22345"}
	cls := cluster.NewCluster(self)
	if err := etcd.Watch(cls); err != nil {
		t.Fatal(err)
	}

	// 添加服务
	nnode := &cluster.Node{Name: "test1", Type: 2, Id: 2, Addr: "192.168.1.1:22345"}
	if err := etcd.Put(nnode); err != nil {
		t.Log(err)
		return
	}

	// 查询
	list, err := etcd.Get()
	if err != nil {
		t.Log(err)
		return
	}
	for _, item := range list {
		t.Log("--->", item)
	}

	// 删除节点
	etcd.Del(nnode)
	time.Sleep(3 * time.Second)
}

func BenchmarkCluster(b *testing.B) {
	self := &cluster.Node{Name: "test1", Type: 1, Id: 1, Addr: "192.168.1.1:22345"}
	cls := cluster.NewCluster(self)
	rtr := router.NewRouter()

	for i := 0; i < 60; i++ {
		cls.Put(&cluster.Node{Name: "test1", Type: int32(i % 5), Id: int32(i), Addr: "192.168.1.1:22345"})
	}
	for i := 0; i < b.N; i++ {
		node := cls.Get(int32(i%5), int32(i))
		if node != nil {
			rtr.Update(uint64(node.GetId()), node)
		}
		if node != nil && node.GetId() == 0 {
			cls.Del(node.GetType(), node.GetId())
		}
	}
	b.Log(b.N)
}
