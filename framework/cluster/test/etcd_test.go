package test

import (
	"log"
	"testing"
	"time"

	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
)

func TestEtcd(t *testing.T) {
	// 创建Etcd客户端
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 创建一个新的租约
	leaseResp, err := client.Lease.Grant(context.Background(), 10) // 10秒后过期的租约
	if err != nil {
		log.Fatal(err)
	}
	leaseID := clientv3.LeaseID(leaseResp.ID)

	// 将键值对与租约相关联
	_, err = client.Put(context.Background(), "key1", "value1", clientv3.WithLease(leaseID))
	if err != nil {
		log.Fatal(err)
	}

	// 保持租约活跃
	keepAliveChan, err := client.Lease.KeepAlive(context.Background(), leaseID)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			resp := <-keepAliveChan
			if resp == nil {
				log.Println("租约已失效")
				break
			}
			log.Printf("收到租约续约响应: TTL = %d\n", resp.TTL)
		}
	}()

	// 模拟程序运行一段时间后撤销租约
	time.Sleep(20 * time.Second)
	_, err = client.Lease.Revoke(context.Background(), leaseID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("租约已撤销")
}
