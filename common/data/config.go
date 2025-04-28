package config

import (
	"time"
	"universal/common/config/internal/manager"
)

func Init(dir, ext string, ttl time.Duration) error {
	return manager.Init(dir, ext, ttl)
}
