package test

import (
	"hego/common/config/internal/manager"
	"hego/common/config/repository/RouteData"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	manager.Init("../../../configure/json/", "json", 1*time.Second)
	m.Run()
}

func TestApi(t *testing.T) {
	t.Log(RouteData.SGet())
}
