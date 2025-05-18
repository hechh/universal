package main

import (
	"flag"
	"fmt"
	"poker_server/common/config"
	"poker_server/common/dao"
	"poker_server/common/pb"
	"poker_server/common/yaml"
	"poker_server/framework"
	"poker_server/framework/library/mlog"
	"poker_server/framework/library/signal"
	"poker_server/server/gate/player"
	"strings"
)

var (
	playerMgr = new(player.PlayerMgr)
)

func main() {
	var cfg string
	var nodeId int
	flag.StringVar(&cfg, "config", "config.yaml", "游戏配置文件")
	flag.IntVar(&nodeId, "id", 1, "服务ID")
	flag.Parse()

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
	framework.SetSelf(node)

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

	// 初始化框架 + 注册内网发送client消息处理
	framework.RegisterBroadcastHandler(broadcastHandler)
	framework.RegisterSendHandler(sendHandler)
	framework.RegisterReplyHandler(sendHandler)
	if err := framework.Init(yamlCfg); err != nil {
		panic(fmt.Sprintf("框架初始化失败: %v", err))
	}

	// 初始化PlayerMgr, 建立websocket服务
	if err := playerMgr.Init(yamlCfg.Cluster[node.Name].Nodes[int32(nodeId)]); err != nil {
		panic(err)
	}

	// 服务退出
	signal.SignalNotify(func() {
		// todo
	})
}

// 处理返回客户端的消息
func sendHandler(head *pb.Head, body []byte) {
	// 发送到客户端
	if len(head.ActorName) <= 0 || len(head.FuncName) <= 0 {
		playerMgr.SendToClient(head, body)
		return
	}
	// 发送到指定Actor
	if err := playerMgr.Send(head, body); err != nil {
		mlog.Errorf("Actor消息转发失败: %v", err)
	}
}

// 处理返回客户端的消息
func broadcastHandler(head *pb.Head, body []byte) {
	// 发送到客户端
	if len(head.ActorName) <= 0 || len(head.FuncName) <= 0 {
		playerMgr.BroadcastToClient(head, body)
		return
	}

	// 广播到所有Actor
	if err := playerMgr.Send(head, body); err != nil {
		mlog.Errorf("Actor消息转发失败: %v", err)
	}
}
