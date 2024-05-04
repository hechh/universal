package gate

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
	"universal/common/config"
	"universal/common/pb"
	"universal/framework/cluster"
	"universal/framework/network"

	"golang.org/x/net/websocket"
)

var (
	GateCfg GateConfig
)

type GateConfig struct {
	Servers map[int]*config.ServerConfig `yaml:"gate"`
	Etcd    *config.EtcdConfig           `yaml:"etcd"`
	Nats    *config.NatsConfig           `yaml:"nats"`
}

func Run() {
	var serverId int
	var lpath, ypath string
	flag.IntVar(&serverId, "id", 0, "ip:port地址")
	flag.StringVar(&lpath, "log", "", "日志输出目录")
	flag.StringVar(&ypath, "yaml", "", "日志输出目录")
	flag.Parse()
	// 读取配置
	if err := config.LoadConfig(ypath, &GateCfg); err != nil {
		panic(err)
	}
	// 初始化集群
	if err := cluster.Init(GateCfg.Nats.Endpoints, GateCfg.Etcd.Endpoints); err != nil {
		panic(err)
	}
	// 进行服务发现
	serverCfg := GateCfg.Servers[serverId]
	if err := cluster.Discovery(pb.ClusterType_GATE, serverCfg.Addr); err != nil {
		panic(err)
	}
	// 设置消息订阅
	if err := cluster.Subscribe(natsHandle); err != nil {
		panic(err)
	}
	// 注册websocket路由
	http.Handle("/ws", websocket.Handler(wsHandle))
	if err := http.ListenAndServe(":8089", nil); err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

func wsHandle(conn *websocket.Conn) {
	client := network.NewSocketClient(conn, 2*time.Second, 2*time.Second)
	for {
		pac, err := client.Read()
		if err != nil {
			fmt.Sprintln(err)
			return
		}
		// 发送到nats
		cluster.Publish(pac)
	}
}

func natsHandle(pac *pb.Packet) {

}
