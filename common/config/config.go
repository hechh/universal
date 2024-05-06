package config

import (
	"io/ioutil"
	"universal/common/pb"
	"universal/framework/fbasic"

	"gopkg.in/yaml.v2"
)

var (
	GlobalCfg = &GlobalConfig{}
)

type ServerConfig struct {
	Addr  string `yaml:"addr"`
	PProf string `yaml:"pprof"`
}

type EtcdConfig struct {
	Endpoints []string `yaml:"endpoints"`
}

type NatsConfig struct {
	Endpoints string `yaml:"endpoints"`
}

type GlobalConfig struct {
	Gate map[int]*ServerConfig `yaml:"gate"`
	Etcd *EtcdConfig           `yaml:"etcd"`
	Nats *NatsConfig           `yaml:"nats"`
}

func LoadConfig(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_ReadYaml, err)
	}
	if err = yaml.Unmarshal(content, GlobalCfg); err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_YamlUnmarshal, err)
	}
	return nil
}
