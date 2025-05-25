package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/server/client/manager"

	"github.com/spf13/cast"
)

var (
	playerMgr *manager.PlayerMgr
)

func main() {
	var filename string
	var id, port int
	var begin, end int64
	flag.StringVar(&filename, "config", "local.yaml", "游戏配置")
	flag.IntVar(&id, "id", 1, " 节点id")
	flag.IntVar(&port, "port", 22345, " 节点端口")
	flag.Int64Var(&begin, "begin", 100000, "起始uid")
	flag.Int64Var(&end, "end", 100000, "终止uid")
	flag.Parse()

	// 加载游戏配置
	node := &pb.Node{
		Name: strings.ToLower(pb.NodeType_Gate.String()),
		Type: pb.NodeType_Gate,
		Id:   int32(id),
	}
	cfg, err := yaml.LoadConfig(filename, node)
	if err != nil {
		panic(fmt.Sprintf("游戏配置加载失败: %v", err))
	}
	nodeCfg := cfg.Gate[node.Id]
	playerMgr = manager.NewPlayerMgr(node, nodeCfg)

	playerMgr.Login(uint64(begin), uint64(end))

	// 请求 http 服务，接受请求
	http.HandleFunc("/api", handle)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	// 1. 只允许POST方法
	if r.Method != http.MethodGet {
		http.Error(w, "只支持post方法", http.StatusMethodNotAllowed)
		return
	}
	// 2. 解析URL参数
	query := r.URL.Query()
	cmd := query.Get("cmd")
	if len(cmd) <= 0 {
		http.Error(w, "cmd不能为空", http.StatusBadRequest)
		return
	}
	val := query.Get("value")
	playerMgr.SendCmd(cast.ToUint32(cmd), val)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
