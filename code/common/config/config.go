package config

import (
	"universal/common/config/domain"
	"universal/common/config/internal/manager"
)

func Init(dir string) error {
	return manager.Init(dir)
}

func Register(name string, cfgs ...domain.IConfig) {
	manager.Register(name, cfgs...)
}
