package test

import (
	"testing"
	"time"
	"universal/framework/config"
	"universal/framework/internal/cluster"
	"universal/framework/internal/discovery"
)

var (
	cfg *config.Config
)

func TestMain(m *testing.M) {
	var err error
	cfg, err = config.NewConfig("../../configure/yaml/local.yaml")
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestEtcd(t *testing.T) {
	etcd, err := discovery.NewEtcd(
		cfg.Etcd.Endpoints,
		discovery.WithTopic("/hch/etcd_test/"),
		discovery.WithNode(cluster.NewNode),
	)
	if err != nil {
		t.Log(err)
		return
	}

	// 监视
	self := &cluster.Node{Name: "test1", Type: 1, Id: 1, Addr: "192.168.1.1:22345"}
	cls := cluster.NewCluster(self)
	if err := etcd.Watch(cls); err != nil {
		t.Fatal(err)
	}

	// 添加服务
	nnode := &cluster.Node{Name: "test1", Type: 2, Id: 2, Addr: "192.168.1.1:22345"}
	if err := etcd.Put(nnode); err != nil {
		t.Fatal(err)
	}

	// 查询
	list, err := etcd.Get()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range list {
		t.Log("--->", item)
	}

	// 删除节点
	etcd.Del(nnode)
	t.Log(" 删除链接")
	if err := etcd.Close(); err != nil {
		t.Log(err)
	}
}

func TestEtcdKeepAlive(t *testing.T) {
	etcd, err := discovery.NewEtcd(
		cfg.Etcd.Endpoints,
		discovery.WithTopic("/hch/etcd_test/"),
		discovery.WithNode(cluster.NewNode),
	)
	if err != nil {
		t.Log(err)
		return
	}

	// 监视
	self := &cluster.Node{Name: "test1", Type: 1, Id: 1, Addr: "192.168.1.1:22345"}
	cls := cluster.NewCluster(self)
	if err := etcd.Watch(cls); err != nil {
		t.Fatal(err)
	}

	if err := etcd.KeepAlive(self, 6); err != nil {
		t.Fatal(err)
	}

	time.Sleep(20 * time.Second)
}
