package test

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"universal/common/pb"
	"universal/framework/network"

	"golang.org/x/net/websocket"
)

const (
	LimitTimes = 50000
)

func TestMain(m *testing.M) {
	go func() {
		http.Handle("/ws", websocket.Handler(wsHandle))
		http.ListenAndServe(":8089", nil)
	}()
	m.Run()
}

func wsHandle(conn *websocket.Conn) {
	client := network.NewSocketClient(conn)
	defer conn.Close()
	sendCh := make(chan *pb.Packet, 1)
	// 接受数据
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			pac, err := client.Read()
			if err != nil {
				return
			}
			sendCh <- pac
			if pac.Head.SeqID >= LimitTimes {
				fmt.Println("server read finished")
				return
			}
		}
	}()
	// 循环处理
	defer wg.Wait()
	for {
		select {
		case item := <-sendCh:
			fmt.Println("---server--->", item)
			if item.Head.SeqID >= LimitTimes {
				fmt.Println("select server finished")
				return
			}
			if err := client.Send(item); err != nil {
				return
			}
		}
	}
}
func TestClient(t *testing.T) {
	ws, err := websocket.Dial("ws://localhost:8089/ws", "", "http://localhost")
	if err != nil {
		t.Log("Error connecting to WebSocket server:", err)
		return
	}
	defer ws.Close()
	client := network.NewSocketClient(ws)
	sendCh := make(chan *pb.Packet, 1)
	// 发送数据
	sendCh <- &pb.Packet{
		Head: &pb.PacketHead{
			SendType:       pb.SendType_NODE,
			ApiCode:        2,
			UID:            100100600,
			SrcClusterType: pb.ClusterType_GATE,
			DstClusterType: pb.ClusterType_GAME,
		},
	}
	// 接受协程
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			pac, err := client.Read()
			if err != nil {
				return
			}
			pac.Head.SeqID++
			sendCh <- pac
			if pac.Head.SeqID >= LimitTimes {
				fmt.Println("client read finished")
				return
			}
		}
	}()
	// 循环处理
	defer wg.Wait()
	for {
		select {
		case item := <-sendCh:
			fmt.Println("---client--->", item)
			if item.Head.SeqID >= LimitTimes {
				fmt.Println("select client finished")
				return
			}
			if err := client.Send(item); err != nil {
				return
			}
		}
	}
}
