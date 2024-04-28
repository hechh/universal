package gate

/*
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
*/
