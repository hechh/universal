package test

import (
	"testing"
	"time"
	"universal/common/config/internal/manager"
	"universal/common/config/repository/route"
)

func TestMain(m *testing.M) {
	manager.Init("../../../configure/bytes/", 1*time.Second)
	m.Run()
}

func TestApi(t *testing.T) {
	t.Log(route.Get(1))
}
