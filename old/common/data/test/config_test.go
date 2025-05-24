package test

import (
	"testing"
	"time"
	"universal/common/config/internal/manager"
	"universal/common/config/repository/RouteData"
)

func TestMain(m *testing.M) {
	manager.Init("../../../configure/json/", "json", 1*time.Second)
	m.Run()
}

func TestApi(t *testing.T) {
	t.Log(RouteData.SGet())
}
