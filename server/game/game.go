package main

import (
	"flag"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/cluster"
	"universal/library/mlog"
	"universal/library/util"
	"universal/message"
)

func main() {
	var conf string
	var nodeId int
	flag.StringVar(&conf, "config", "config.yaml", "游戏配置文件")
	flag.IntVar(&nodeId, "id", 1, "服务节点ID")
	flag.Parse()

	// 加载配置
	cfg, srvCfg, nn, err := yaml.LoadAndParse(conf, pb.NodeType_Game, int32(nodeId))
	if err != nil {
		panic(err)
	}

	// 初始化日志库
	mlog.Init(srvCfg.LogPath, nn.Name, srvCfg.LogLevel)

	message.Init()

	// 初始化集群
	if err := cluster.Init(nn, srvCfg, cfg); err != nil {
		panic(err)
	}

	// todo

	// 信号捕捉
	util.SignalNotify(func() {
		cluster.Close()
		mlog.Close()
	})
}
