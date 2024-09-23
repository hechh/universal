package global

import (
	"fmt"
	"path/filepath"
	"strings"
	"universal/common/pb"
)

var (
	cfg        *Config
	serverId   uint32
	serverType pb.SERVER
)

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

func GetServer() pb.SERVER {
	return serverType
}

func GetServerName() string {
	return serverType.String()
}

func GetConfig() *Config {
	return cfg
}

func GetServerCfg() *ServerConfig {
	return cfg.Server[serverId]
}

func GetMysqlConfig() map[uint32]*DbConfig {
	return cfg.Mysql
}

func GetRedisConfig() map[uint32]*DbConfig {
	return cfg.Redis
}

func GetMongodbConfig() map[uint32]*DbConfig {
	return cfg.Mongodb
}

func GetLogConfig() *LogConfig {
	return cfg.Log
}

func GetEtcdConfig() *EtcdConfig {
	return cfg.Etcd
}

func GetNatsConfig() *NatsConfig {
	return cfg.Nats
}
