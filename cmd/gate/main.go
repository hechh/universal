package main

import (
	"flag"
	"universal/common/pb"
	"universal/framework"
	"universal/server/gate"
)

func main() {
	var serverId int
	var lpath, ypath string
	flag.IntVar(&serverId, "id", 0, "ip:port地址")
	flag.StringVar(&lpath, "log", "", "日志输出目录")
	flag.StringVar(&ypath, "yaml", "", "日志输出目录")
	flag.Parse()
	// 设置全局变量
	framework.SetGlobal(serverId, pb.ClusterType_GATE)
	gate.Run(lpath, ypath)
}
