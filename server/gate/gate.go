package gate

import (
	"flag"
	"log"
	"net/http"
	"universal/common/config"
	"universal/common/pb"
	"universal/framework/cluster"
	"universal/framework/fbasic"
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
	// 注册websocket路由
	http.Handle("/ws", websocket.Handler(wsHandle))
	if err := http.ListenAndServe(":8089", nil); err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

func wsHandle(conn *websocket.Conn) {
	log.Println("wsHandle begin, ", conn.RemoteAddr().String())
	defer func() {
		log.Println("wsHandle closed, ", conn.RemoteAddr().String())
		conn.Close()
	}()
	client := network.NewSocketClient(conn)
	// 设置消息订阅
	cluster.Subscribe(func(pac *pb.Packet) {
		log.Println("Subscribe: ", pac)
		if err := client.Send(pac); err != nil {
			log.Fatal(err)
		}
	})

	for {
		pac, err := client.Read()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("read: ", pac)
		// 转发
		if err := dispatcher(client, pac); err != nil {
			log.Println(err)
			return
		}
	}
}

func dispatcher(client *network.SocketClient, pac *pb.Packet) error {
	// 设置头信息
	head := pac.Head
	if head.ApiCode <= int32(pb.ApiCode_NONE_END_REQUEST) {
		return fbasic.NewUError(1, pb.ErrorCode_NotSupported, head.ApiCode)
	}
	head.DstClusterType = fbasic.ApiCodeToClusterType(head.ApiCode)
	local := cluster.GetLocalClusterNode()
	head.SrcClusterType = local.ClusterType
	head.SrcClusterID = local.ClusterID
	// 转发到nats
	return cluster.Publish(pac)
}
