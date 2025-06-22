package internal

import (
	"poker_server/common/pb"
	"poker_server/common/yaml"
	"poker_server/server/db/internal/module/generator"
	"poker_server/server/db/internal/module/player"
	"poker_server/server/db/internal/module/report"
	"poker_server/server/db/internal/module/rummy_room"
	"poker_server/server/db/internal/module/texas_room"
	"poker_server/server/db/internal/module/user_info"
)

var (
	genMgr               = generator.NewGeneratorMgr()
	playerMgr            = player.NewPlayerDataMgr()
	texasMgr             = texas_room.NewTexasRoomMgr()
	reportMgr            = report.NewReportDataMgr()
	rummyMgr             = rummy_room.NewDbRummyRoomMgr()
	rummySettlePool      = rummy_room.NewRummySettlePool()
	rummySettleMatchPool = rummy_room.NewRummySettleMatchPool()
	userInfoMgr          = user_info.NewUserInfoMgr()
)

func Init(node *pb.Node, cfg *yaml.PhpConfig) error {
	player.Init(cfg)
	if err := genMgr.Init(); err != nil {
		return err
	}
	if err := texasMgr.Init(); err != nil {
		return err
	}
	if err := playerMgr.Init(); err != nil {
		return err
	}
	if err := reportMgr.Init(); err != nil {
		return err
	}
	if err := rummyMgr.Init(); err != nil {
		return err
	}
	if err := rummySettlePool.Init(); err != nil {
		return err
	}
	if err := rummySettleMatchPool.Init(); err != nil {
		return err
	}
	if err := userInfoMgr.Init(); err != nil {
		return err
	}
	return nil
}

func Close() {
	genMgr.Close()
	texasMgr.Close()
	playerMgr.Close()
	reportMgr.Close()
	rummySettlePool.Close()
	userInfoMgr.Close()
}
