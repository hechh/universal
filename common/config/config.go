package config

import (
	"hego/common/config/internal/manager"
	"time"
)

func Init(dir string, ttl time.Duration) error {
	return manager.Init(dir, ttl)
}
