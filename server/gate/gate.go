package main

import (
	"flag"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/cluster"
	"universal/library/mlog"
	"universal/library/signal"
	"universal/message"
	"universal/server/gate/internal/player"
)

func main() {
	var conf string
	var nodeId int
	flag.StringVar(&conf, "config", "config.yaml", "游戏配置文件")
	flag.IntVar(&nodeId, "id", 1, "服务节点ID")
	flag.Parse()

	// 加载配置
	cfg, srvCfg, nn, err := yaml.LoadAndParse(conf, pb.NodeType_NodeTypeGate, int32(nodeId))
	if err != nil {
		panic(err)
	}

	// 初始化日志库
	message.Init()
	mlog.Init(srvCfg.LogPath, nn.Name, mlog.StringToLevel(srvCfg.LogLevel))

	// 初始化集群
	if err := cluster.Init(cfg, srvCfg, nn); err != nil {
		panic(err)
	}

	if err := player.Init(cfg, srvCfg); err != nil {
		panic(err)
	}

	// 信号捕捉
	signal.SignalNotify(func() {
		player.Close()
		cluster.Close()
		mlog.Close()
	})
}
