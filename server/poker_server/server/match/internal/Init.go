package internal

import (
	"poker_server/server/match/internal/module/rummy"
	"poker_server/server/match/internal/module/sng_room"
	"poker_server/server/match/internal/module/texas_room"
)

var (
	texasRoomMgr = texas_room.NewMatchTexasRoomMgr()
	rummyRoomMgr = rummy.NewMatchRummyRoomMgr()
	sngRoomMgr   = sng_room.NewSngRoomMgr()
)

func Init() error {
	if err := texasRoomMgr.Load(); err != nil {
		return err
	}
	if err := rummyRoomMgr.Load(); err != nil {
		return err
	}

	return nil
}

func Close() {
	texasRoomMgr.Close()
	rummyRoomMgr.Close()
	sngRoomMgr.Close()
}
