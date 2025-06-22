package internal

import (
	"poker_server/server/builder/internal/module/rummy"
	"poker_server/server/builder/internal/module/texas"
)

var (
	genMgr      = texas.NewBuilderTexasGenerator()
	rummyGenMgr = rummy.NewBuilderRummyGenerator()
)

func Init() error {
	if err := genMgr.Load(); err != nil {
		return err
	}
	if err := rummyGenMgr.Load(); err != nil {
		return err
	}
	return nil
}

func Close() {
	genMgr.Stop()
}
