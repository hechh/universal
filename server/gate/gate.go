package main

import (
	"flag"
	"universal/common/dao"
	"universal/common/pb"
	"universal/common/yaml"
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
	cfg, err := yaml.Load(path, "gate")
	if err != nil {
		panic(err)
	}
	// 初始化plog
	plog.Init(cfg.Log.Level, cfg.Log.Path, "gate")

	// 初始化redis
	if err := dao.InitRedis(cfg.Redis); err != nil {
		panic(err)
	}
	// 初始化集群
	if err := cluster.Init(cfg, pb.SERVICE_GATE, uint32(id), 600); err != nil {
		panic(err)
	}

}
