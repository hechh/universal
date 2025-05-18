package test

import (
	"fmt"
	"poker_server/common/pb"
	"poker_server/common/yaml"
	"poker_server/framework/internal/core/discovery"
	"poker_server/framework/internal/service"
	"strings"
	"testing"
)

func TestEtcd(t *testing.T) {
	// 加载游戏配置
	node := &pb.Node{
		Name: strings.ToLower(pb.NodeType_Gate.String()),
		Type: pb.NodeType_Gate,
		Id:   int32(1),
	}
	yamlCfg, err := yaml.LoadConfig("./local.yaml", node)
	if err != nil {
		panic(fmt.Sprintf("游戏配置加载失败: %v", err))
	}

	// 初始化日志库
	nodeCfg := yamlCfg.Cluster[node.Name]
	client, err := discovery.NewEtcd(yamlCfg.Etcd)
	if err != nil {
		panic(err)
	}

	if err := client.Register(node, nodeCfg.DicoveryExpire); err != nil {
		panic(err)
	}

	client.Watch(service.GetCluster())

	select {}
}
