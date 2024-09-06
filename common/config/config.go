package config

import (
	"time"
	"universal/common/config/domain"
	"universal/common/config/internal/manager"
)

func Init(dir string, ttl time.Duration) error {
	return manager.Init(dir, ttl)
}

func Register(name string, cfgs ...domain.IConfig) {
	manager.Register(name, cfgs...)
}
