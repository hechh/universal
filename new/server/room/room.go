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
	"poker_server/server/room/texas"
	"strings"
)

func main() {
	var cfg string
	var nodeId int
	flag.StringVar(&cfg, "config", "config.yaml", "游戏配置文件")
	flag.IntVar(&nodeId, "id", 1, "服务ID")
	flag.Parse()

	// 加载游戏配置
	node := &pb.Node{
		Name: strings.ToLower(pb.NodeType_Room.String()),
		Type: pb.NodeType_Gate,
		Id:   int32(nodeId),
	}
	yamlcfg, err := yaml.LoadConfig(cfg, node)
	if err != nil {
		panic(fmt.Sprintf("游戏配置加载失败: %v", err))
	}
	framework.SetSelf(node)

	// 初始化日志库
	if err := mlog.Init(yamlcfg.Cluster[node.Name]); err != nil {
		panic(fmt.Sprintf("日志库初始化失败: %v", err))
	}

	// 初始化redis
	if err := dao.InitRedis(yamlcfg.Redis); err != nil {
		panic(fmt.Sprintf("redis初始化失败: %v", err))
	}

	// 初始化mysql

	// 初始化游戏配置
	if err := config.InitLocal(yamlcfg.Configure); err != nil {
		panic(err)
	}

	// 初始化框架
	if err := framework.Init(yamlcfg); err != nil {
		panic(fmt.Sprintf("框架初始化失败: %v", err))
	}

	// 功能模块初始化
	Init()

	// 服务退出
	signal.SignalNotify(func() {

	})
}

func Init() {
	texas.Init()
}
