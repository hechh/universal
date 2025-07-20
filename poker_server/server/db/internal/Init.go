package internal

import (
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/cluster"
	"poker_server/library/util"
	"poker_server/server/db/internal/generator"
	"poker_server/server/db/internal/player"
	"poker_server/server/db/internal/report"
	"poker_server/server/db/internal/room_info"
	"poker_server/server/db/internal/rummy_room"
	"poker_server/server/db/internal/texas_room"
	"poker_server/server/db/internal/user_info"
)

var (
	genMgr                  = generator.NewGeneratorMgr()
	playerMgr               = player.NewPlayerDataMgr()
	texasMgr                = texas_room.NewTexasRoomMgr()
	reportMgr               = report.NewReportDataMgr()
	rummyMgr                = rummy_room.NewDbRummyRoomMgr()
	rummySettlePool         = rummy_room.NewRummySettlePool()
	rummySettleMatchPool    = rummy_room.NewRummySettleMatchPool()
	rummyExtSettleMatchPool = rummy_room.NewRummyExtSettleMatchPool()
	userInfoMgr             = user_info.NewUserInfoMgr()
	roomInfoMgr             = room_info.NewRoomInfoMgr()
)

func Close() {
	genMgr.Close()
	texasMgr.Close()
	playerMgr.Close()
	reportMgr.Close()
	rummySettlePool.Close()
	userInfoMgr.Close()
	rummyExtSettleMatchPool.Close()
	roomInfoMgr.Close()
}

func Init(node *pb.Node) {
	util.Must(cluster.SetBroadcastHandler(framework.DefaultHandler))
	util.Must(cluster.SetSendHandler(framework.DefaultHandler))
	util.Must(cluster.SetReplyHandler(framework.DefaultHandler))
	util.Must(genMgr.Init())
	util.Must(texasMgr.Init())
	util.Must(playerMgr.Init())
	util.Must(reportMgr.Init())
	util.Must(rummyMgr.Init())
	util.Must(rummySettlePool.Init())
	util.Must(rummySettleMatchPool.Init())
	util.Must(rummyExtSettleMatchPool.Init())
	util.Must(userInfoMgr.Init())
	util.Must(roomInfoMgr.Init())
}
