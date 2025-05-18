package main

import (
	"fmt"
	"poker_server/common/config"
	"poker_server/common/dao"
	"poker_server/common/dao/repository/redis/login_session"
	"poker_server/common/pb"
	"poker_server/common/yaml"
	"poker_server/framework/library/async"
	"poker_server/framework/library/mlog"
	"poker_server/framework/network"
	"poker_server/server/gate/frame"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
)

func TestClient(t *testing.T) {
	cfg := "../../env/hch/local.yaml"
	nodeId := int32(1)
	/*
		flag.StringVar(&cfg, "config", "config.yaml", "游戏配置文件")
		flag.IntVar(&nodeId, "id", 1, "服务ID")
		flag.Parse()
	*/
	// 加载游戏配置
	node := &pb.Node{
		Name: strings.ToLower(pb.NodeType_Gate.String()),
		Type: pb.NodeType_Gate,
		Id:   int32(nodeId),
	}
	yamlCfg, err := yaml.LoadConfig(cfg, node)
	if err != nil {
		panic(fmt.Sprintf("游戏配置加载失败: %v", err))
	}
	// 初始化日志库
	if err := mlog.Init(yamlCfg.Cluster[node.Name]); err != nil {
		panic(fmt.Sprintf("日志库初始化失败: %v", err))
	}
	// 初始化redis
	if err := dao.InitRedis(yamlCfg.Redis); err != nil {
		panic(fmt.Sprintf("redis初始化失败: %v", err))
	}
	// 初始化游戏配置
	if err := config.InitLocal(yamlCfg.Configure); err != nil {
		panic(err)
	}
	// 写入 session
	sessionId := "123456"
	now := time.Now().Unix()
	ttl := int64(100000)
	if err := login_session.Set(sessionId, &pb.LoginSession{
		SessionId:  sessionId,
		AccountId:  1000001,
		CreateTime: now,
		ExpireTime: now + ttl,
	}, time.Duration(ttl)*time.Second); err != nil {
		panic(err)
	}
	// 建立WebSocket连接
	t.Log("======================建立链接======================")
	ws, err := websocket.Dial(fmt.Sprintf("ws://localhost%s/ws", node.Addr), "", "http://localhost/")
	if err != nil {
		panic(err)
	}
	defer ws.Close()
	sock := network.NewSocket(ws, 1024*1024)
	sock.Register(&frame.Frame{})

	// 发送登录消息
	t.Log("======================发送登录请求======================")
	head := &pb.Head{
		SendType:    pb.SendType_POINT,
		DstNodeType: pb.NodeType_Gate,
		DstNodeId:   node.Id,
		IdType:      pb.IdType_UID,
		Id:          1000001,
		RouteId:     1000001,
		Cmd:         uint32(pb.CMD_GATE_LOGIN_REQUEST),
	}
	loginReq := &pb.GateLoginRequest{SessionId: sessionId}
	if err := sock.WriteMsg(head, loginReq); err != nil {
		panic(err)
	}
	// 接收登录返回消息
	t.Log("======================等待登录应答======================")
	pack := &pb.Packet{}
	if err := sock.Read(pack); err != nil {
		panic(err)
	}
	loginRsp := &pb.GateLoginResponse{}
	if err := proto.Unmarshal(pack.Body, loginRsp); err != nil {
		panic(err)
	}
	fmt.Println("登录返回:", loginRsp)

	// 接受应答
	async.SafeGo(mlog.Fatalf, func() {
		for {
			pack := &pb.Packet{}
			if err := sock.Read(pack); err != nil {
				mlog.Errorf("读取数据包失败: %v", err)
				break
			}
			fmt.Println("--------response--------")
			head := pack.Head
			switch head.Cmd {
			case uint32(pb.CMD_GATE_HEART_RESPONSE):
				fmt.Println("心跳包返回")
			default:
				fmt.Println("未知消息:", head.Cmd)
			}
		}
	})

	// 发送心跳
	async.SafeGo(mlog.Fatalf, func() {
		head := &pb.Head{
			SendType:    pb.SendType_POINT,
			DstNodeType: pb.NodeType_Gate,
			DstNodeId:   node.Id,
			IdType:      pb.IdType_UID,
			Id:          1000001,
			RouteId:     1000001,
			Cmd:         uint32(pb.CMD_GATE_HEART_REQUEST),
		}
		// 循环发送心跳
		tt := time.NewTicker(2 * time.Second)
		for {
			<-tt.C
			fmt.Println("--------heart--------")
			if err := sock.WriteMsg(head, &pb.GateHeartRequest{}); err != nil {
				fmt.Println("发送心跳包失败:", err)
				break
			}
		}
	})
	select {}
	/*
	 */
}
