package main

import (
	"flag"
	"path/filepath"

	"forevernine.com/planet/server/tool/gomaker/internal/manager"

	"forevernine.com/planet/server/tool/gomaker/repository/dbsrv"
	"forevernine.com/planet/server/tool/gomaker/repository/module"
	"forevernine.com/planet/server/tool/gomaker/repository/redis"
	"forevernine.com/planet/server/tool/gomaker/repository/reward"
	"forevernine.com/planet/server/tool/gomaker/repository/secure"
	"forevernine.com/planet/server/tool/gomaker/repository/xerror"
	"forevernine.com/planet/server/tool/gomaker/repository/xlsx"

	"forevernine.com/planet/server/tool/gomaker/internal/base"
	"forevernine.com/planet/server/tool/gomaker/internal/service"
)

func init() {
	redis.Init()
	xlsx.Init()
	dbsrv.Init()
	reward.Init()
	xerror.Init()
	module.Init()
	secure.Init()
}

func main() {
	// 解析命令行
	var action, name, rule, path string
	flag.StringVar(&action, "a", "", "规则类型")
	flag.StringVar(&path, "p", "", "生成代码路径")
	flag.StringVar(&name, "n", "", "名字(例如pbname,xlsx文件名等)")
	flag.StringVar(&rule, "r", "", "规则(例如gomaker:xxx:aa|xxx|...)")
	flag.Parse()
	// 生成路径
	if len(path) > 0 {
		path = filepath.Join(base.ROOT, path)
	}
	// 添加规则
	if len(rule) > 0 {
		manager.ParseRule(name, "//@"+rule)
	}
	// 需要解析的文件
	files, err := filepath.Glob(filepath.Join(base.ROOT, "common/pbclass/*.pb.go"))
	if err != nil {
		panic(err)
	}
	files = append(files, filepath.Join(base.ROOT, "/srv/dbsrv/internal/biz/proc/internal/const.go"))

	// 解析文件
	service.ParseFiles(files...)
	// 生成代码
	service.GenCode(path, action)
}
