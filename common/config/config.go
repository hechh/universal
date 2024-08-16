package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// 服务配置
type ServerConfig struct {
	Host  string `yaml:host`
	Pprof string `yaml:pprof`
	Gops  string `yaml:gops`
}

// 日志配置
type LogConfig struct {
	Level uint32 `yaml:level`
	Path  string `yaml:path`
}

// 数据库配置
type DbConfig struct {
	DbName   string            `yaml:dbname`
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

type StubConfig map[string]uint64

type Config struct {
	Log     *LogConfig               `yaml:log`
	Server  map[uint32]*ServerConfig `yaml:server`
	Mysql   map[uint32]*DbConfig     `yaml:mysql`
	Redis   map[uint32]*DbConfig     `yaml:redis`
	Mongodb map[uint32]*DbConfig     `yaml:mongodb`
	Etcd    *EtcdConfig              `yaml:etcd`
	Nats    *NatsConfig              `yaml:nats`
	Stub    StubConfig               `yaml:stub`
}

// 加载配置
func LoadConfig(dir, appname string) (*Config, error) {
	// 加载配置文件
	content1, err := ioutil.ReadFile(filepath.Join(dir, appname+".yaml"))
	if err != nil {
		return nil, err
	}
	content2, err := ioutil.ReadFile(filepath.Join(dir, "common.yaml"))
	if err != nil {
		return nil, err
	}
	content1 = append(content1, content2...)
	// 序列化
	result := new(Config)
	if err := yaml.Unmarshal(content1, result); err != nil {
		return nil, err
	}
	return result, nil
}
