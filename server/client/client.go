package main

import (
	"flag"
	"fmt"
	"net/http"
	"universal/common/config"
	"universal/common/pb"
	"universal/common/token"
	"universal/common/yaml"
	"universal/framework/cluster"
	"universal/library/mlog"
	"universal/library/safe"
	"universal/library/util"
	"universal/message"
	"universal/server/client/internal/player"
	"universal/server/client/internal/stat"

	"github.com/spf13/cast"
)

var (
	playerMgr *player.ClientPlayerMgr
)

func main() {
	var filename string
	var id, port int
	var begin, end int64
	flag.StringVar(&filename, "config", "config.yaml", "游戏配置")
	flag.IntVar(&id, "id", 1, " 节点id")
	flag.IntVar(&port, "port", 22345, " 节点端口")
	flag.Int64Var(&begin, "begin", 1000001, "起始uid")
	flag.Int64Var(&end, "end", 1000222, "终止uid")
	flag.Parse()

	// 加载游戏配置
	yamlcfg, node, err := yaml.LoadConfig(filename, pb.NodeType_NodeTypeClient, int32(id))
	if err != nil {
		panic(fmt.Sprintf("游戏配置加载失败: %v", err))
	}
	nodeCfg := yamlcfg.Client[node.Id]

	// 初始化日志库
	mlog.Init(node.Name, node.Id, nodeCfg.LogLevel, nodeCfg.LogPath)
	token.Init(yamlcfg.Common.SecretKey)

	mlog.Infof("初始化游戏配置")
	util.Must(config.Init(yamlcfg.Etcd, yamlcfg.Data))
	message.Init()

	// 初始化PlayerMgr
	stat.NewStatistics()
	playerMgr = player.NewClientPlayerMgr(node, nodeCfg)
	playerMgr.Login(uint64(begin), uint64(end))

	// 请求 http 服务，接受请求
	http.HandleFunc("/api", handle)
	safe.Go(func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	})

	util.SignalNotify(func() {
		cluster.Close()
		mlog.Close()
	})
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
	var routeId uint64
	rr := query.Get("route_id")
	if len(rr) <= 0 {
		routeId = 0
	} else {
		routeId = cast.ToUint64(rr)
	}
	jsonVal := query.Get("json")
	uid := uint64(0)
	if uidStr := query.Get("uid"); len(uidStr) > 0 {
		uid = cast.ToUint64(uidStr)
	}

	playerMgr.SendCmd(cast.ToUint32(cmd), uid, routeId, jsonVal)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
