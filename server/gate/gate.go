package gate

import (
	"flag"
	"log"
	"net/http"
	"universal/common/config"
	"universal/common/pb"
	"universal/framework"
	"universal/framework/network"

	"golang.org/x/net/websocket"
)

func Run() {
	var serverId int
	var lpath, ypath string
	flag.IntVar(&serverId, "id", 0, "ip:port地址")
	flag.StringVar(&lpath, "log", "", "日志输出目录")
	flag.StringVar(&ypath, "yaml", "", "日志输出目录")
	flag.Parse()
	// 读取配置
	if err := config.LoadConfig(ypath); err != nil {
		panic(err)
	}
	// 设置全局变量
	framework.SetGlobal(serverId, pb.ClusterType_GATE)
	// 核心框架初始化
	serverCfg := config.GlobalCfg.Gate[serverId]
	if err := framework.Init(serverCfg.Addr, config.GlobalCfg.Etcd.Endpoints, config.GlobalCfg.Nats.Endpoints); err != nil {
		panic(err)
	}
	// 注册websocket路由
	http.Handle("/ws", websocket.Handler(wsHandle))
	if err := http.ListenAndServe(":8089", nil); err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

func wsHandle(conn *websocket.Conn) {
	user := NewUser(network.NewSocketClient(conn))
	// 认证
	if !user.Auth() {
		conn.Close()
		return
	}
	// 初始化
	if err := user.Init(); err != nil {
		log.Fatalln("user init nats handler: ", err)
		conn.Close()
		return
	}
	// 循环接受消息
	log.Println("websocket connected...", conn.RemoteAddr().String())
	defer func() {
		log.Println("wsHandle closed: ", conn.RemoteAddr().String())
		conn.Close()
	}()
	user.LoopRead()
}
