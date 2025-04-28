package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type MongodbConfig struct {
	DbName   string            `yaml:dbname`
	User     string            `yaml:user`
	Password string            `yaml:password`
	Host     string            `yaml:host`
	Slave    map[uint32]string `yaml:slave`
}

type MysqlConfig struct {
	DbName   string            `yaml:dbname`
	User     string            `yaml:user`
	Password string            `yaml:password`
	Host     string            `yaml:host`
	Slave    map[uint32]string `yaml:slave`
}

type RedisConfig struct {
	User     string            `yaml:user`
	Password string            `yaml:password`
	Host     string            `yaml:host`
	Slave    map[uint32]string `yaml:slave`
}

type EtcdConfig struct {
	Endpoints []string `yaml:endpoints`
}

type NatsConfig struct {
	Endpoints string `yaml:endpoints`
}

type NodeConfig struct {
	Name string `yaml:name`
	Addr string `yaml:addr`
}

type GateConfig struct {
	LogLevel  string                `yaml:log_level`
	LogPath   string                `yaml:log_path`
	LogPrefix string                `yaml:log_prefix`
	RouteType int32                 `yaml:route_type`
	Nodes     map[int32]*NodeConfig `yaml:nodes`
}

type Config struct {
	Mysql   map[uint32]*MysqlConfig   `yaml:mysql`
	Redis   map[uint32]*RedisConfig   `yaml:redis`
	Mongodb map[uint32]*MongodbConfig `yaml:mongodb`
	Etcd    *EtcdConfig               `yaml:etcd`
	Nats    *NatsConfig               `yaml:nats`
	Gate    *GateConfig               `yaml:gate`
}

func (c *Config) Unmarshal(buf []byte) error {
	return yaml.Unmarshal(buf, c)
}

func NewConfig(filename string) (*Config, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := new(Config)
	if err = cfg.Unmarshal(content); err != nil {
		return nil, err
	}
	return cfg, nil
}
