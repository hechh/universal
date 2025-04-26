package test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"universal/common/config"
	"universal/framework/internal/cluster"
	"universal/framework/internal/discovery"

	clientv3 "go.etcd.io/etcd/client/v3"
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

func TestClient(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:33601"},
		DialTimeout: 5,
	})
	if err != nil {
		log.Fatalf("Failed to create etcd client: %v", err)
	}
	defer cli.Close()

	ctx := context.Background()

	// 写入键值对
	_, err = cli.Put(ctx, "myKey", "myValue")
	if err != nil {
		log.Fatalf("Failed to put key-value pair: %v", err)
	}
	fmt.Println("Key-value pair written successfully")

	// 读取键值对
	resp, err := cli.Get(ctx, "myKey")
	if err != nil {
		log.Fatalf("Failed to get key-value pair: %v", err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("Read key-value pair: %s : %s\n", string(ev.Key), string(ev.Value))
	}

	// 删除键值对
	_, err = cli.Delete(ctx, "myKey")
	if err != nil {
		log.Fatalf("Failed to delete key-value pair: %v", err)
	}
	fmt.Println("Key-value pair deleted successfully")
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
