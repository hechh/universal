package yaml

import (
	"os"
	"poker_server/common/pb"
	"poker_server/framework/library/uerror"

	"gopkg.in/yaml.v3"
)

type MongodbConfig struct {
	DbName   string           `yaml:"dbname"`
	User     string           `yaml:"user"`
	Password string           `yaml:"password"`
	Host     string           `yaml:"host"`
	Slave    map[int32]string `yaml:"slave"`
}

type SlaveConfig struct {
	DbName   string `yaml:"dbname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type MysqlConfig struct {
	DbName   string                 `yaml:"dbname"`
	User     string                 `yaml:"user"`
	Password string                 `yaml:"password"`
	Host     string                 `yaml:"host"`
	Slave    map[int32]*SlaveConfig `yaml:"slave"`
}

type RedisConfig struct {
	DbName   string           `yaml:"dbname"`
	Db       int              `yaml:"db"`
	User     string           `yaml:"user"`
	Password string           `yaml:"password"`
	Host     string           `yaml:"host"`
	Slave    map[int32]string `yaml:"slave"`
}

type EtcdConfig struct {
	Channel   string   `yaml:"channel"`
	Endpoints []string `yaml:"endpoints"`
}

type NatsConfig struct {
	Channel   string `yaml:"channel"`
	Endpoints string `yaml:"endpoints"`
}

type NodeConfig struct {
	Env            string           `yaml:"env"`
	LogLevel       string           `yaml:"log_level"`
	LogPath        string           `yaml:"log_path"`
	RouterExpire   int64            `yaml:"router_expire"`
	DicoveryExpire int64            `yaml:"discovery_expire"`
	Nodes          map[int32]string `yaml:"nodes"`
}

type ConfigureConfig struct {
	LocalPath string `yaml:"local_path"`
}

type Config struct {
	Mysql     map[int32]*MysqlConfig   `yaml:"mysql"`
	Redis     map[int32]*RedisConfig   `yaml:"redis"`
	Mongodb   map[int32]*MongodbConfig `yaml:"mongodb"`
	Etcd      *EtcdConfig              `yaml:"etcd"`
	Nats      *NatsConfig              `yaml:"nats"`
	Cluster   map[string]*NodeConfig   `yaml:"cluster"`
	Configure *ConfigureConfig         `yaml:"configure"`
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
	nodeCfg, ok := cfg.Cluster[node.Name]
	if !ok {
		return nil, uerror.New(1, -1, "服务配置不存在: %s", node.Type.String())
	}
	addr, ok := nodeCfg.Nodes[node.Id]
	if !ok {
		return nil, uerror.New(1, -1, "服务节点不存在: %s-%d", node.Type.String(), node.Id)
	}
	node.Addr = addr
	return cfg, nil
}
