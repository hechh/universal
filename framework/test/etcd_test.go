package test

import (
	"testing"
	"universal/common/config"
	"universal/framework/internal/cluster"
	"universal/framework/internal/discovery"
)

var cfg *config.Config

func TestMain(m *testing.M) {
	var err error
	cfg, err = config.NewConfig("../../configure/yaml/local.yaml")
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestEtcdPut(t *testing.T) {
	etcd, err := discovery.NewEtcd(
		cfg.Etcd.Endpoints,
		discovery.WithPath("/hch/etcd_test/"),
		discovery.WithParse(cluster.NewNode),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer etcd.Close()

	if err := etcd.Put(&cluster.Node{
		Name: "test",
		Type: 1,
		Id:   12,
		Addr: "192.168.1.1:22345",
	}); err != nil {
		t.Log(err)
		return
	}

	list, err := etcd.Get()
	if err != nil {
		t.Log(err)
		return
	}
	for _, item := range list {
		t.Log("--->", item)
	}
}
