package test

import (
	"sync"
	"testing"
	"time"
	"universal/common/pb"
	"universal/server/gate/internal/player"

	"golang.org/x/net/websocket"
)

func TestClient(t *testing.T) {
	ws, err := websocket.Dial("ws://127.0.0.1:10100/ws", "", "http://localhost")
	if err != nil {
		t.Log("Error connecting to WebSocket server:", err)
		return
	}
	defer ws.Close()

	client := player.NewPlayer(ws)
	pac := &pb.Packet{
		Head: &pb.PacketHead{
			SendType:      pb.SendType_PLAYER,
			ApiCode:       int32(pb.ApiCode_GATE_LOGIN_REQUEST),
			UID:           100100600,
			SrcServerType: pb.ServerType_GATE,
			DstServerType: pb.ServerType_GATE,
		},
	}

	times := 0
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			times++
			t.Log(times, "==========>", client.Send(pac))
			if times > 10 {
				return
			}
		}
	}
}

func BenchmarkClient(b *testing.B) {
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(uid uint64) {
			defer wg.Done()
			ws, err := websocket.Dial("ws://127.0.0.1:8089/ws", "", "http://localhost")
			if err != nil {
				b.Log("Error connecting to WebSocket server:", err)
				return
			}
			defer ws.Close()
			client := player.NewPlayer(ws)
			pac := &pb.Packet{
				Head: &pb.PacketHead{
					SendType:      pb.SendType_PLAYER,
					ApiCode:       int32(pb.ApiCode_GATE_LOGIN_REQUEST),
					UID:           uid,
					SrcServerType: pb.ServerType_GATE,
					DstServerType: pb.ServerType_GATE,
				},
			}
			for i := 0; i < 10; i++ {
				client.Send(pac)
			}
		}(100100100 + uint64(i))
	}
	wg.Wait()
}
