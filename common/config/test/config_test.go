package test

import (
	"hego/common/config/internal/manager"
	"hego/common/config/repository/route_config"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	manager.Init("../../../configure/bytes/", 1*time.Second)
	m.Run()
}

func TestApi(t *testing.T) {
	t.Log(route_config.GetByID(1))
	time.Sleep(5 * time.Second)
}
