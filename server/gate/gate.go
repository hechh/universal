package main

import (
	"flag"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/cluster"
	"universal/library/mlog"
	"universal/library/signal"
)

func main() {
	var conf string
	var nodeId int
	flag.StringVar(&conf, "config", "config.yaml", "游戏配置文件")
	flag.IntVar(&nodeId, "id", 1, "服务节点ID")
	flag.Parse()

	// 加载配置
	cfg, err := yaml.ParseConfig(conf)
	if err != nil {
		panic(err)
	}
	srvCfg := yaml.GetNodeConfig(cfg, pb.NodeType_NodeTypeGate, int32(nodeId))
	if srvCfg == nil {
		panic("节点配置不存在")
	}
	nn := yaml.GetNode(srvCfg, pb.NodeType_NodeTypeGate, int32(nodeId))

	// 初始化日志库
	mlog.Init(srvCfg.LogPath, nn.Name, mlog.StringToLevel(srvCfg.LogLevel))

	// 初始化集群
	if err := cluster.Init(cfg, srvCfg, nn); err != nil {
		panic(err)
	}

	// 信号捕捉
	signal.SignalNotify(func() {

	})
}
