package test

import (
	"testing"
	"universal/common/config"
)

func TestConfig(t *testing.T) {
	t.Log(config.LoadConfig("gate.yaml"))
}
