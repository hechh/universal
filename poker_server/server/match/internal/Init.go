package internal

import (
	"poker_server/framework"
	"poker_server/framework/cluster"
	"poker_server/library/util"
	"poker_server/server/match/internal/module/rummy"
	"poker_server/server/match/internal/module/sng_room"
	"poker_server/server/match/internal/module/texas_room"
)

var (
	texasRoomMgr = texas_room.NewMatchTexasRoomMgr()
	rummyRoomMgr = rummy.NewMatchRummyRoomMgr()
	sngRoomMgr   = sng_room.NewSngRoomMgr()
)

func Init() {
	util.Must(cluster.SetBroadcastHandler(framework.DefaultHandler))
	util.Must(cluster.SetSendHandler(framework.DefaultHandler))
	util.Must(cluster.SetReplyHandler(framework.DefaultHandler))
	util.Must(texasRoomMgr.Load())
	util.Must(rummyRoomMgr.Load())
}

func Close() {
	texasRoomMgr.Close()
	rummyRoomMgr.Close()
	sngRoomMgr.Close()
}
