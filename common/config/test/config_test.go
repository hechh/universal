package test

import (
	"testing"
	"universal/common/config"
)

func TestConfig(t *testing.T) {
	cfg, err := config.LoadConfig("../../../env/config", "game")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(cfg.Server[1])
}
