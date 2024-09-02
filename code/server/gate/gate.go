package main

import (
	"flag"
	"universal/common/dao"
	"universal/common/global"
	"universal/common/pb"
	"universal/framework/cluster"
	"universal/framework/plog"
)

func main() {
	var id int64
	var path string
	flag.Int64Var(&id, "id", 1, "服务节点ID")
	flag.StringVar(&path, "cfg", "./", "yaml配置文件路径")
	flag.Parse()

	// 加载配置
	if err := global.Init(pb.SERVICE_GATE, uint32(id), global.GATE, path); err != nil {
		panic(err)
	}
	// 初始化plog
	if logCfg := global.GetLogConfig(); logCfg != nil {
		plog.Init(logCfg.Level, logCfg.Path, global.GetAppName())
	}
	// 初始化redis
	if err := dao.InitRedis(global.GetRedisConfig()); err != nil {
		panic(err)
	}
	// 初始化mysql
	if err := dao.InitMysql(global.GetMysqlConfig()); err != nil {
		panic(err)
	}
	// 初始化集群
	if err := cluster.Init(global.GetConfig(), global.GetPlatform(), uint32(id), 600); err != nil {
		panic(err)
	}

}
