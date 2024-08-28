package test

import (
	"testing"
	"universal/common/dao/internal/manager"
	"universal/common/global"
)

func TestMain(m *testing.M) {
	cfg := &global.Config{}
	global.LoadFile("../../../env/yaml/common.yaml", cfg)
	global.LoadFile("../../../env/yaml/game.yaml", cfg)

	if err := manager.InitRedis(cfg.Redis); err != nil {
		panic(err)
		return
	}
	m.Run()
}

func TestRedis(t *testing.T) {
	cli := manager.GetRedis(1)
	if cli == nil {
		return
	}

	result, err := cli.Get("test-hch")
	t.Log(result, err)

	if err = cli.Set("test-hch", 123); err != nil {
		t.Log(err)
		return
	}

	result, err = cli.Get("test-hch")
	t.Log(result, err)
}
