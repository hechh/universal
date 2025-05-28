package yaml

import (
	"fmt"
	"os"
	"universal/common/pb"
	"universal/library/uerror"

	"gopkg.in/yaml.v3"
)

type SlaveConfig struct {
	DbName   string `yaml:"dbname"`
	Db       int    `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type DbConfig struct {
	DbName   string                 `yaml:"dbname"`
	Db       int                    `yaml:"db"`
	User     string                 `yaml:"user"`
	Password string                 `yaml:"password"`
	Host     string                 `yaml:"host"`
	Slave    map[int32]*SlaveConfig `yaml:"slave"`
}

type EtcdConfig struct {
	Topic     string   `yaml:"topic"`
	Endpoints []string `yaml:"endpoints"`
}

type NatsConfig struct {
	Topic     string `yaml:"topic"`
	Channel   string `yaml:"channel"`
	Endpoints string `yaml:"endpoints"`
}

type CommonConfig struct {
	Env            string `yaml:"env"`
	ConfigurePath  string `yaml:"configure_path"`
	RouterExpire   int64  `yaml:"router_expire"`
	DicoveryExpire int64  `yaml:"discovery_expire"`
	SecretKey      string `yaml:"secret_key"`
}

type ServerConfig struct {
	LogLevel string `yaml:"log_level"`
	LogFile  string `yaml:"log_file"`
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	HttpPort int    `yaml:"http_port"`
}

type Config struct {
	Mysql   map[int32]*DbConfig     `yaml:"mysql"`
	Redis   map[int32]*DbConfig     `yaml:"redis"`
	Mongodb map[int32]*DbConfig     `yaml:"mongodb"`
	Etcd    *EtcdConfig             `yaml:"etcd"`
	Nats    *NatsConfig             `yaml:"nats"`
	Common  *CommonConfig           `yaml:"common"`
	Gate    map[int32]*ServerConfig `yaml:"gate"`
	Db      map[int32]*ServerConfig `yaml:"db"`
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

func LoadConfig(filename string, node *pb.Node) (*Config, error) {
	cfg, err := NewConfig(filename)
	if err != nil {
		return nil, uerror.New(1, -1, "配置文件加载失败: %v", err)
	}
	switch node.Type {
	case pb.NodeType_Gate:
		srvCfg, ok := cfg.Gate[node.Id]
		if !ok {
			return nil, uerror.New(1, -1, "服务节点不存在: %s-%d", node.Type.String(), node.Id)
		}
		node.Addr = fmt.Sprintf("%s:%d", srvCfg.Ip, srvCfg.Port)
	case pb.NodeType_Db:
		srvCfg, ok := cfg.Db[node.Id]
		if !ok {
			return nil, uerror.New(1, -1, "服务节点不存在: %s-%d", node.Type.String(), node.Id)
		}
		node.Addr = fmt.Sprintf("%s:%d", srvCfg.Ip, srvCfg.Port)
	}
	return cfg, nil
}
