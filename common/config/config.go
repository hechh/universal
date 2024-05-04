package config

import (
	"io/ioutil"
	"universal/common/pb"
	"universal/framework/fbasic"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	ServerID int    `yaml:"id"`
	Addr     string `yaml:"addr"`
	PProf    string `yaml:"pprof"`
}

type EtcdConfig struct {
	Endpoints []string `yaml:"endpoints"`
}

type NatsConfig struct {
	Endpoints string `yaml:"endpoints"`
}

func LoadConfig(path string, data interface{}) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_ReadYaml, err)
	}
	if err = yaml.Unmarshal(content, data); err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_YamlUnmarshal, err)
	}
	return nil
}
