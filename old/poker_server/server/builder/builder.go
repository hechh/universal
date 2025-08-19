package main

import (
	"flag"
	"fmt"
	"poker_server/common/config"
	"poker_server/common/pb"
	"poker_server/common/redis"
	"poker_server/common/yaml"
	"poker_server/framework"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"poker_server/library/signal"
	"poker_server/library/util"
	"poker_server/message"
	"poker_server/server/builder/internal"
)

func main() {
	var cfg string
	var nodeId int
	flag.StringVar(&cfg, "config", "config.yaml", "游戏配置文件")
	flag.IntVar(&nodeId, "id", 1, "服务ID")
	flag.Parse()

	// 加载游戏配置
	yamlcfg, node, err := yaml.LoadConfig(cfg, pb.NodeType_NodeTypeBuilder, int32(nodeId))
	if err != nil {
		panic(fmt.Sprintf("游戏配置加载失败: %v", err))
	}
	nodeCfg := yamlcfg.Builder[node.Id]

	// 初始化日志库
	mlog.Init(node.Name, node.Id, nodeCfg.LogLevel, nodeCfg.LogPath)

	// 初始化游戏配置
	mlog.Infof("初始化游戏配置")
	util.Must(config.Init(yamlcfg.Etcd, yamlcfg.Data))

	// 初始化redis
	mlog.Infof("初始化redis配置")
	util.Must(redis.Init(yamlcfg.Redis))

	// 初始化框架
	mlog.Infof("启动框架服务")
	util.Must(framework.Init(node, nodeCfg, yamlcfg))

	// 功能模块初始化
	mlog.Infof("初始化功能模块")
	message.Init()
	internal.Init()

	// 服务退出
	mlog.Infof("服务关闭: %v", node)
	signal.SignalNotify(func() {
		internal.Close()
		cluster.Close()
		mlog.Close()
	})
}
