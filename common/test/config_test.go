package test

import (
	"testing"
	"universal/common/yaml"
)

func TestConfig(t *testing.T) {
	cfg := &yaml.Config{}
	t.Log("error: ", yaml.LoadFile("../../env/config/common.yaml", cfg))
	t.Log(cfg.Server[1], cfg.Etcd)
	t.Log("error: ", yaml.LoadFile("../../env/config/game.yaml", cfg))
	t.Log(cfg.Server[1], cfg.Etcd)
}

func TestLoad(t *testing.T) {
	cfg, err := yaml.Load("../../env/config", "gate")
	t.Log(err, cfg.Server[1], cfg.Etcd)
}
