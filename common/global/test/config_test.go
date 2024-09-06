package test

import (
	"testing"
	"universal/common/global"
)

func TestConfig(t *testing.T) {
	cfg := &global.Config{}
	t.Log("error: ", global.LoadFile("../../env/config/common.yaml", cfg))
	t.Log(cfg.Server[1], cfg.Etcd)
	t.Log("error: ", global.LoadFile("../../env/config/game.yaml", cfg))
	t.Log(cfg.Server[1], cfg.Etcd)
}

func TestLoad(t *testing.T) {
	cfg, err := global.Load("../../env/config", "gate")
	t.Log(err, cfg.Server[1], cfg.Etcd)
}
