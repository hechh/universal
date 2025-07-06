package yaml

import (
	"os"
	"strings"
	"universal/common/pb"

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
	Db       int32                  `yaml:"db"`
	Prefix   string                 `yaml:"prefix"`
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
	Endpoints string `yaml:"endpoints"`
}

type DataConfig struct {
	IsRemote  bool     `yaml:"is_remote"`
	Topic     string   `yaml:"topic"`
	Path      string   `yaml:"path"`
	Endpoints []string `yaml:"endpoints"`
}

type CommonConfig struct {
	SecretKey string `yaml:"secret_key"`
}

type NodeConfig struct {
	RouterTTL    int64  `yaml:"router_ttl"`
	DiscoveryTTL int64  `yaml:"discovery_ttl"`
	LogLevel     string `yaml:"log_level"`
	LogPath      string `yaml:"log_path"`
	Ip           string `yaml:"ip"`
	Port         int    `yaml:"port"`
	HttpPort     int    `yaml:"http_port"`
}

type Config struct {
	Env     string                `yaml:"env"`
	Mysql   map[int32]*DbConfig   `yaml:"mysql"`
	Redis   map[int32]*DbConfig   `yaml:"redis"`
	Mongodb map[int32]*DbConfig   `yaml:"mongodb"`
	Etcd    *EtcdConfig           `yaml:"etcd"`
	Nats    *NatsConfig           `yaml:"nats"`
	Data    *DataConfig           `yaml:"data"`
	Common  *CommonConfig         `yaml:"common"`
	Client  map[int32]*NodeConfig `yaml:"client"`
	Gate    map[int32]*NodeConfig `yaml:"gate"`
	Room    map[int32]*NodeConfig `yaml:"room"`
	Match   map[int32]*NodeConfig `yaml:"match"`
	Db      map[int32]*NodeConfig `yaml:"db"`
	Build   map[int32]*NodeConfig `yaml:"build"`
	Game    map[int32]*NodeConfig `yaml:"game"`
	Gm      map[int32]*NodeConfig `yaml:"gm"`
}

func ParseConfig(filename string) (*Config, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := new(Config)
	if err = yaml.Unmarshal(content, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func GetNodeConfig(cfg *Config, nodeType pb.NodeType, nodeId int32) *NodeConfig {
	switch nodeType {
	case pb.NodeType_NodeTypeGate:
		return cfg.Gate[nodeId]
	case pb.NodeType_NodeTypeRoom:
		return cfg.Room[nodeId]
	case pb.NodeType_NodeTypeMatch:
		return cfg.Match[nodeId]
	case pb.NodeType_NodeTypeDb:
		return cfg.Db[nodeId]
	case pb.NodeType_NodeTypeBuild:
		return cfg.Build[nodeId]
	case pb.NodeType_NodeTypeGame:
		return cfg.Game[nodeId]
	case pb.NodeType_NodeTypeGm:
		return cfg.Gm[nodeId]
	case pb.NodeType_NodeTypeClient:
		return cfg.Client[nodeId]
	}
	return nil
}

func GetNode(cfg *NodeConfig, nodeType pb.NodeType, nodeId int32) *pb.Node {
	return &pb.Node{
		Name: strings.TrimPrefix(nodeType.String(), "NodeType"),
		Type: nodeType,
		Id:   nodeId,
		Ip:   cfg.Ip,
		Port: int32(cfg.Port),
	}
}
