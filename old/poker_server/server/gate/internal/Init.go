package internal

import (
	"fmt"
	"net/http"
	"poker_server/common/pb"
	"poker_server/common/token"
	"poker_server/common/yaml"
	"poker_server/framework/actor"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"poker_server/library/safe"
	"poker_server/library/util"
	"poker_server/server/gate/internal/http_api"
	"poker_server/server/gate/internal/player"
)

var (
	playerMgr = new(player.GatePlayerMgr)
)

func Init(cfg *yaml.NodeConfig, com *yaml.CommonConfig) {
	util.Must(cluster.SetBroadcastHandler(handler))
	util.Must(cluster.SetSendHandler(handler))
	util.Must(cluster.SetReplyHandler(handler))
	token.Init(com.SecretKey)

	// 初始化模块
	util.Must(playerMgr.Init(cfg.Ip, cfg.Port))

	// 初始化Actor
	safe.Go(func() { initApi(cfg) })
}

func Close() {
	playerMgr.Stop()
}

// 处理返回客户端的消息
func handler(head *pb.Head, body []byte) {
	mlog.Trace(head, "收到Nats数据包 body:%d", len(body))
	if len(head.Dst.FuncName) <= 0 || head.Dst.FuncName == "SendToClient" {
		head.Dst.ActorName = "Player"
		head.Dst.FuncName = "SendToClient"
	}
	head.ActorName = head.Dst.ActorName
	head.FuncName = head.Dst.FuncName
	head.ActorId = head.Dst.ActorId
	if err := actor.Send(head, body); err != nil {
		mlog.Errorf("Actor消息转发失败: %v", err)
	}
}

func initApi(cfg *yaml.NodeConfig) error {
	api := http.NewServeMux()
	api.HandleFunc("/api/room/token", http_api.GenToken)
	api.HandleFunc("/api/room/list", http_api.TexasRoomList)
	api.HandleFunc("/api/game/buyin", http_api.BuyInApi)
	api.HandleFunc("/api/room/sng/list", http_api.SngRoomList)
	api.HandleFunc("/api/room/rummy/list", http_api.RummyRoomList)
	api.HandleFunc("/api/game/reconnect", http_api.GameReconnect)
	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.HttpPort), api)
}
