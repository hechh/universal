package config

import (
	"hego/common/config/internal/manager"
	"time"
)

func Init(dir, ext string, ttl time.Duration) error {
	return manager.Init(dir, ext, ttl)
}
