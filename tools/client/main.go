package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"universal/common/dao"
	"universal/common/global"
	"universal/library/mlog"
	"universal/tools/client/internal/httpkit"
	"universal/tools/client/internal/player"
)

var (
	cfg *global.Config
)

func main() {
	// 解析命令行
	var level, plat uint64
	var sid, path string
	flag.StringVar(&path, "cfg", "配置文件目录", "yaml配置文件")
	flag.StringVar(&sid, "gate", "1", "websocket连接的gate服务节点列表")
	flag.Uint64Var(&plat, "plat", 0, "0: 本地服务， 1: 内网服务, 2:改时间专服")
	flag.Uint64Var(&level, "log", mlog.LOG_DEFAULT, "plog输出日志等级")
	flag.Parse()

	// 加载配置
	tmpcfg, err := global.Load(path, "gate")
	if err != nil {
		panic(fmt.Errorf("配置文件加载错误: %v", err))
	}
	cfg = tmpcfg

	// 初始化日志
	mlog.Init(uint32(level), "./log", "client")

	// 初始化redis
	if err := dao.InitRedis(cfg.Redis); err != nil {
		panic(fmt.Errorf("redis连接池初始化失败: %v", err))
	}

	player.Init(cfg.Server, plat, sid)
	httpkit.Init()

	// 阻塞
	signalHandleBlock(func(sig os.Signal) {
	}, os.Interrupt, os.Kill)
}

func signalHandleBlock(f func(os.Signal), sig ...os.Signal) {
	ch := make(chan os.Signal, 0)
	args := append([]os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL}, sig...)
	signal.Notify(ch, args...)
	for item := range ch {
		f(item)
		os.Exit(0)
	}
}
