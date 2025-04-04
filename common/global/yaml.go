package global

import (
	"fmt"
	"hego/common/pb"
	"io/ioutil"
	"path/filepath"
	"strings"

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
	Alias    string            `yaml:alias`
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

type Config struct {
	Log     *LogConfig               `yaml:log`
	Server  map[uint32]*ServerConfig `yaml:server`
	Mysql   map[uint32]*DbConfig     `yaml:mysql`
	Redis   map[uint32]*DbConfig     `yaml:redis`
	Mongodb map[uint32]*DbConfig     `yaml:mongodb`
	Etcd    *EtcdConfig              `yaml:etcd`
	Nats    *NatsConfig              `yaml:nats`
}

func LoadFile(filename string, cfg *Config) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, cfg)
}

func Init(dir string, typ pb.SERVER, srvid uint32) error {
	tmpCfg := &Config{}
	// 加载服务配置
	filename := filepath.Join(dir, fmt.Sprintf("%s.yaml", strings.ToLower(typ.String())))
	if err := LoadFile(filename, tmpCfg); err != nil {
		return err
	}
	// 加载通用配置
	if err := LoadFile(filepath.Join(dir, "common.yaml"), tmpCfg); err != nil {
		return err
	}
	cfg = tmpCfg
	serverId = srvid
	serverType = typ
	return nil
}
