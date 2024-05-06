package test

import (
	"fmt"
	"sync"
	"testing"

	"universal/common/pb"
	"universal/framework/notify/domain"

	"universal/framework/notify/internal/service"
)

var (
	clientWg = sync.WaitGroup{}
)

func TestMain(m *testing.M) {
	// 初始化支持类型
	err := service.Init(domain.NotifyTypeNats, "localhost:4222,172.16.126.208:33601,172.16.126.208:33602,172.16.126.208:33603")
	if err != nil {
		panic(err)
		return
	}
	if err := service.Subscribe("/nats", natsHandle); err != nil {
		panic(err)
		return
	}
	m.Run()
}

func natsHandle(pac *pb.Packet) {
	fmt.Println("=======>", pac)
	clientWg.Done()
}

func TestNats(t *testing.T) {
	pac := &pb.Packet{
		Head: &pb.PacketHead{
			SendType:       pb.SendType_NODE,
			ApiCode:        2,
			UID:            100100600,
			SrcClusterType: pb.ClusterType_GATE,
			DstClusterType: pb.ClusterType_GAME,
		},
	}
	clientWg.Add(1)
	if err := service.Publish("/nats", pac); err != nil {
		t.Log(err)
		return
	}
	clientWg.Wait()
}
