package gate

import (
	"log"
	"net/http"
	"universal/common/config"
	"universal/common/pb"
	"universal/framework"
	"universal/framework/cluster"
	"universal/framework/fbasic"
	"universal/framework/network"
	"universal/framework/notify"
	"universal/framework/packet"

	"golang.org/x/net/websocket"
)

func Run(lpath, ypath string) {
	// 读取配置
	if err := config.LoadConfig(ypath); err != nil {
		panic(err)
	}
	// 核心框架初始化
	serverCfg := config.GlobalCfg.Gate[framework.GetServerID()]
	if err := framework.Init(serverCfg.Addr, config.GlobalCfg.Etcd.Endpoints, config.GlobalCfg.Nats.Endpoints); err != nil {
		panic(err)
	}
	// 注册websocket路由
	http.Handle("/ws", websocket.Handler(wsHandle))
	if err := http.ListenAndServe(":8089", nil); err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

func auth(client *network.SocketClient, pac *pb.Packet) error {
	head := pac.Head
	// 判断第一个请求是否为登陆认证包
	if head.ApiCode != int32(pb.ApiCode_GATE_LOGIN_REQUEST) {
		return fbasic.NewUError(1, pb.ErrorCode_Parameter, pb.ApiCode_GATE_LOGIN_REQUEST, head.ApiCode)
	}
	// 登陆认证
	rsp, err := packet.Call(fbasic.NewDefaultContext(head), pac.Buff)
	if err != nil {
		return err
	}
	// 返回认证结果
	if err := client.SendRsp(head, rsp); err != nil {
		return err
	}
	// 判断是否成功
	if head := rsp.(fbasic.IProto).GetHead(); head.Code > 0 {
		return fbasic.NewUError(1, pb.ErrorCode(head.Code), "gate login auth is failed")
	}
	return nil
}

func wsHandle(conn *websocket.Conn) {
	var err error
	var pac *pb.Packet
	defer func() {
		if err != nil {
			log.Fatalln("websocket connect is failed: ", err)
			conn.Close()
		} else {
			log.Println("websocket closed: ", conn.RemoteAddr().String())
			conn.Close()
		}
	}()
	client := network.NewSocketClient(conn)
	// 读取认证请求
	if pac, err = client.Read(); err != nil {
		return
	}
	// 认证
	if err := auth(client, pac); err != nil {
		return
	}
	// 订阅消息
	user := NewUser(client)
	self := cluster.GetLocalClusterNode()
	if err = notify.Subscribe(fbasic.GetPlayerChannel(self.ClusterType, self.ClusterID, pac.Head.UID), user.NatsHandle); err != nil {
		return
	}
	// 循环接受消息
	log.Println("websocket connected...", conn.RemoteAddr().String())
	user.LoopRead()
}
