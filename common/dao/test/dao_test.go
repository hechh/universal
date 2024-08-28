package test

import (
	"testing"
	"universal/common/global"
)

func TestMain(m *testing.M) {
	cfg := &global.Config{}
	global.LoadFile("../../../env/yaml/common.global", cfg)
	global.LoadFile("../../../env/yaml/game.global", cfg)
	m.Run()
}
