package internal

import (
	"poker_server/framework"
	"poker_server/framework/cluster"
	"poker_server/library/util"
	"poker_server/server/builder/internal/rummy"
	"poker_server/server/builder/internal/texas"
)

var (
	genMgr      = texas.NewBuilderTexasGenerator()
	rummyGenMgr = rummy.NewBuilderRummyGenerator()
)

func Init() {
	util.Must(cluster.SetBroadcastHandler(framework.DefaultHandler))
	util.Must(cluster.SetSendHandler(framework.DefaultHandler))
	util.Must(cluster.SetReplyHandler(framework.DefaultHandler))
	util.Must(genMgr.Load())
	util.Must(rummyGenMgr.Load())
}

func Close() {
	genMgr.Stop()
}
