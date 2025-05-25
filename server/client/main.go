package main

import (
	"flag"
	"fmt"
	"strings"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/library/signal"
	"universal/server/client/manager"
)

var (
	playerMgr *manager.PlayerMgr
)

func main() {
	var filename string
	var id int
	var begin, end int64
	flag.StringVar(&filename, "config", "config.yaml", "游戏配置")
	flag.IntVar(&id, "id", 1, " 节点id")
	flag.Int64Var(&begin, "begin", 100000, "起始uid")
	flag.Int64Var(&end, "end", 100000, "终止uid")
	flag.Parse()

	// 加载游戏配置
	node := &pb.Node{
		Name: strings.ToLower(pb.NodeType_Gate.String()),
		Type: pb.NodeType_Gate,
		Id:   int32(id),
	}
	cfg, err := yaml.LoadConfig(filename, node)
	if err != nil {
		panic(fmt.Sprintf("游戏配置加载失败: %v", err))
	}
	nodeCfg := cfg.Gate[node.Id]
	playerMgr = manager.NewPlayerMgr(node, nodeCfg)

	playerMgr.Login(uint64(begin), uint64(end))

	// 服务退出
	signal.SignalNotify(func() {
		// todo
	})
}
