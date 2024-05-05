package test

import (
	"testing"
	"time"
	"universal/common/pb"
	"universal/framework/network"

	"golang.org/x/net/websocket"
)

func TestClient(t *testing.T) {
	ws, err := websocket.Dial("ws://localhost:8089/ws", "", "http://localhost")
	if err != nil {
		t.Log("Error connecting to WebSocket server:", err)
		return
	}
	defer ws.Close()

	client := network.NewSocketClient(ws)
	pac := &pb.Packet{
		Head: &pb.PacketHead{
			SendType:       pb.SendType_POINT,
			ApiCode:        int32(pb.ApiCode_GATE_BEGIN_REQUEST + 1),
			UID:            100100600,
			SrcClusterType: pb.ClusterType_GATE,
			DstClusterType: pb.ClusterType_GAME,
		},
	}

	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			t.Log("==========>", client.Send(pac))
		}
	}
}
