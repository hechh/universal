package global

import (
	"hego/common/pb"
)

var (
	cfg        *Config
	serverId   uint32
	serverType pb.SERVER
)

func GetServerType() pb.SERVER {
	return serverType
}

func GetServerName() string {
	return serverType.String()
}

func GetCfg() *Config {
	return cfg
}

func GetServerCfg() *ServerConfig {
	return cfg.Server[serverId]
}

func GetMysqlCfg() map[uint32]*DbConfig {
	return cfg.Mysql
}

func GetRedisCfg() map[uint32]*DbConfig {
	return cfg.Redis
}

func GetMongodbCfg() map[uint32]*DbConfig {
	return cfg.Mongodb
}

func GetLogCfg() *LogConfig {
	return cfg.Log
}

func GetEtcdCfg() *EtcdConfig {
	return cfg.Etcd
}

func GetNatsCfg() *NatsConfig {
	return cfg.Nats
}
