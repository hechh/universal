package main

import (
	"flag"
	"fmt"
	"poker_server/common/config"
	"poker_server/common/dao"
	"poker_server/common/pb"
	"poker_server/common/yaml"
	"poker_server/framework"
	"poker_server/library/async"
	"poker_server/library/mlog"
	"poker_server/library/signal"
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
	if err := mlog.Init(yamlcfg.Common.Env, nodeCfg.LogLevel, nodeCfg.LogFile); err != nil {
		panic(fmt.Sprintf("日志库初始化失败: %v", err))
	}
	async.Init(mlog.Errorf)

	// 初始化游戏配置
	mlog.Infof("初始化游戏配置")
	if err := config.Init(yamlcfg.Etcd, yamlcfg.Common); err != nil {
		panic(err)
	}

	// 初始化redis
	mlog.Infof("初始化redis配置")
	if err := dao.InitRedis(yamlcfg.Redis); err != nil {
		panic(fmt.Sprintf("redis初始化失败: %v", err))
	}

	// 初始化框架
	mlog.Infof("启动框架服务: %v", node)
	if err := framework.InitDefault(node, nodeCfg, yamlcfg); err != nil {
		panic(fmt.Sprintf("框架初始化失败: %v", err))
	}

	// 功能模块初始化
	if err := internal.Init(); err != nil {
		panic(fmt.Sprintf("功能模块初始化失败: %v", err))
	}

	// 服务退出
	signal.SignalNotify(func() {
		internal.Close()
		framework.Close()
		mlog.Close()
	})
}
