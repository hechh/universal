package test

import (
	"testing"
	"universal/common/config"
)

func TestConfig(t *testing.T) {
	cfg := &config.Config{}
	t.Log("error: ", config.LoadFile("../../env/config/common.yaml", cfg))
	t.Log(cfg.Server[1], cfg.Etcd)
	t.Log("error: ", config.LoadFile("../../env/config/game.yaml", cfg))
	t.Log(cfg.Server[1], cfg.Etcd)
}

func TestLoad(t *testing.T) {
	cfg, err := config.LoadConfig("../../env/config", "gate")
	t.Log(err, cfg.Server[1], cfg.Etcd)
}
