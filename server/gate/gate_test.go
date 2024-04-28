package gate_test

import (
	"testing"
	"universal/common/config"
	"universal/common/pb"
	"universal/framework/cluster"
)

var (
	GateCfg GateConfig
)

type GateConfig struct {
	Servers map[int]*config.ServerConfig `yaml:"gate"`
	Etcd    *config.EtcdConfig           `yaml:"etcd"`
	Nats    *config.NatsConfig           `yaml:"nats"`
}

func TestRun(t *testing.T) {
	// 读取配置
	if err := config.LoadConfig("../../env/gate.yaml", &GateCfg); err != nil {
		panic(err)
	}
	serverId := 1
	serverCfg := GateCfg.Servers[serverId]
	node := &pb.ClusterNode{
		ClusterType: pb.ClusterType_GATE,
		Ip:          serverCfg.IP,
		Port:        int32(serverCfg.Port),
	}
	// 初始化集群
	if err := cluster.Init(node, GateCfg.Nats.Endpoints, GateCfg.Etcd.Endpoints); err != nil {
		panic(err)
	}
	// 进行服务发现
	if err := cluster.Discovery(); err != nil {
		panic(err)
	}
	// 设置消息订阅
	if err := cluster.Subscribe(natsHandle); err != nil {
		panic(err)
	}
}

func natsHandle(pac *pb.Packet) {

}
