package test

import (
	"testing"
	"time"
	"universal/common/config"
	"universal/framework/define"
	"universal/framework/internal/cluster"
	"universal/framework/internal/discovery"
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
	types := map[int32]int32{1: 1, 2: 2, 3: 3, 4: 4}
	cls := cluster.NewCluster(self, types)
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
