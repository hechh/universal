package test

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"universal/common/pb"
	"universal/framework/network"

	"golang.org/x/net/websocket"
)

const (
	LimitTimes = 10000
)

var (
	times int32
	cli   = sync.WaitGroup{}
	sr    = sync.WaitGroup{}
)

func wsHandle(conn *websocket.Conn) {
	fmt.Println("======> websocket connected")
	client := network.NewSocketClient(conn, 2*time.Second, 2*time.Second)
	defer conn.Close()

	for {
		pac, err := client.Read()
		fmt.Println(err, "----->", pac)
		if err != nil {
			return
		}
		cli.Done()
		sr.Add(1)
		pac.Head.SeqID++
		if err := client.Send(pac); err != nil {
			fmt.Println("---->", err)
			return
		}
		sr.Wait()
		atomic.AddInt32(&times, 1)
		if atomic.LoadInt32(&times) > LimitTimes {
			return
		}
	}
}
func TestMain(m *testing.M) {
	go func() {
		http.Handle("/ws", websocket.Handler(wsHandle))
		http.ListenAndServe(":8089", nil)
	}()
	m.Run()
}

func TestClient(t *testing.T) {
	ws, err := websocket.Dial("ws://localhost:8089/ws", "", "http://localhost")
	if err != nil {
		t.Log("Error connecting to WebSocket server:", err)
		return
	}
	defer ws.Close()

	pac := &pb.Packet{
		Head: &pb.PacketHead{
			SendType:       pb.SendType_POINT,
			ApiCode:        2,
			UID:            100100600,
			SrcClusterType: pb.ClusterType_GATE,
			DstClusterType: pb.ClusterType_GAME,
		},
	}
	client := network.NewSocketClient(ws, 2*time.Second, 2*time.Second)
	for {
		cli.Add(1)
		if err := client.Send(pac); err != nil {
			t.Log("=====>", err)
			return
		}
		cli.Wait()

		pac, err = client.Read()
		t.Log(err, "=====>", pac)
		if err != nil {
			return
		}
		sr.Done()
		pac.Head.SeqID++
		if atomic.LoadInt32(&times) > LimitTimes {
			return
		}
	}
}
