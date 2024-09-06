package global

import (
	"path/filepath"
	"universal/common/pb"
)

const (
	GAME = "game"
	GATE = "gate"
	DB   = "db"
	GM   = "gm"
)

var (
	platform pb.SERVER
	appId    uint32
	appName  string
	cfg      *Config
)

func Init(typ pb.SERVER, appid uint32, appname string, dir string) error {
	tmpCfg := &Config{}
	// 加载服务配置
	if err := LoadFile(filepath.Join(dir, appname+".yaml"), tmpCfg); err != nil {
		return err
	}
	// 加载通用配置
	if err := LoadFile(filepath.Join(dir, "common.yaml"), tmpCfg); err != nil {
		return err
	}
	cfg = tmpCfg
	platform = typ
	appId = appid
	appName = appname
	return nil
}

func GetPlatform() pb.SERVER {
	return platform
}

func GetAppName() string {
	return appName
}

func GetConfig() *Config {
	return cfg
}

func GetServerCfg() *ServerConfig {
	return cfg.Server[appId]
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

func GetStub() StubConfig {
	return cfg.Stub
}
