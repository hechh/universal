package test

import (
	"testing"
	"time"
	"universal/framework/define"
	"universal/framework/internal/cluster"
	"universal/framework/internal/network"
	"universal/framework/internal/packet"
)

func TestNsq(t *testing.T) {
	cli, err := network.NewNsq(cfg.Nsq.Nsqd, network.WithTopic("hch_test"), network.WithPacket(packet.NewPacket), network.WithHeader(packet.NewHeader))
	if err != nil {
		t.Fatalf("nats connect err: %v", err)
		return
	}

	// 接受消息
	self := &cluster.Node{Name: "test1", Type: 1, Id: 1, Addr: "192.168.1.1:22345"}
	err = cli.Read(self, func(header define.IHeader, body []byte) {
		t.Logf("header: %v, recv: %v", header, string(body))
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	// 发送消息
	head := &packet.Header{
		SrcNodeType: 1,
		SrcNodeId:   1,
		DstNodeType: 1,
		DstNodeId:   1,
		Cmd:         1,
		Uid:         1,
	}
	for i := 0; i < 5; i++ {
		if err := cli.Send(head.SetDstNode(self), []byte("hello world")); err != nil {
			t.Fatal(err)
			return
		}
	}
	time.Sleep(1 * time.Second)
}
