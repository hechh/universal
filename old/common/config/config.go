package config

import (
	"io/ioutil"
	"universal/framework/common/uerror"

	"gopkg.in/yaml.v2"
)

var (
	globalCfg = &GlobalConfig{}
)

type ServerConfig struct {
	Addr  string `yaml:"addr"`
	PProf string `yaml:"pprof"`
	Gops  string `yaml:"gops"`
}

type EtcdConfig struct {
	Endpoints []string `yaml:"endpoints"`
}

type NatsConfig struct {
	Endpoints string `yaml:"endpoints"`
}

type RedisConfig struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
}

type GlobalConfig struct {
	Etcd   *EtcdConfig             `yaml:"etcd"`
	Nats   *NatsConfig             `yaml:"nats"`
	Redis  map[string]*RedisConfig `yaml:"redis"`
	Server map[int]*ServerConfig   `yaml:"server"`
}

func GetServerConfig(id int) *ServerConfig {
	return globalCfg.Server[id]
}

func GetRedisConfig() map[string]*RedisConfig {
	return globalCfg.Redis
}

func GetEtcdConfig() *EtcdConfig {
	return globalCfg.Etcd
}

func GetNatsConfig() *NatsConfig {
	return globalCfg.Nats
}

func LoadConfig(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return uerror.NewUError(1, -1, err)
	}
	if err = yaml.Unmarshal(content, globalCfg); err != nil {
		return uerror.NewUError(1, -1, err)
	}
	return nil
}
