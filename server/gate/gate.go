package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"universal/common/config"
	"universal/common/pb"
	"universal/framework"
	"universal/framework/actor"
	"universal/framework/cluster"
	"universal/framework/network"
	"universal/framework/profiler"
	"universal/server/gate/internal/broadcast"
	"universal/server/gate/internal/player"

	"golang.org/x/net/websocket"
)

var (
	broads = broadcast.NewBroadcast()
)

func main() {
	var srvID int
	var yaml string
	flag.IntVar(&srvID, "id", 0, "ip:port地址")
	flag.StringVar(&yaml, "yaml", "", "日志输出目录")
	flag.Parse()
	// 加载配置
	if err := config.LoadConfig(yaml); err != nil {
		panic(err)
	}
	// 获取配置
	srvCfg := config.GetServerConfig(srvID)
	if srvCfg == nil {
		panic(fmt.Sprintf("Server id(%d) not found", srvID))
	}
	// 设置过期清理
	actor.SetActorClearExpire(int64(10 * 60 * time.Second))
	cluster.SetRouterClearExpire(int64(10 * 60 * time.Second))
	// 设置性能分析工具
	if err := profiler.InitGops(srvCfg.Gops); err != nil {
		panic(err)
	}
	profiler.InitPprof(srvCfg.PProf)
	// 初始化框架
	etcdCfg := config.GetEtcdConfig()
	natsCfg := config.GetNatsConfig()
	if err := framework.Init(pb.ServerType_GATE, srvCfg.Addr, etcdCfg.Endpoints, natsCfg.Endpoints); err != nil {
		panic(err)
	}
	// 注册节点广播
	self := cluster.GetSelfServerNode()
	if err := network.Subscribe(cluster.GetNodeChannel(self.ServerType, self.ServerID), broads.Send); err != nil {
		panic(err)
	}
	// 注册全局广播
	if err := network.Subscribe(cluster.GetClusterChannel(self.ServerType), broads.Send); err != nil {
		panic(err)
	}
	// 设置信号处理
	framework.SignalHandle(func(sig os.Signal) {
		actor.StopAll()
		cluster.Close()
	})
	// 开启websocket服务
	log.Println("websocket start: ", srvCfg.Addr)
	http.Handle("/ws", websocket.Handler(wsHandle))
	if err := http.ListenAndServe(srvCfg.Addr, nil); err != nil {
		panic(err)
	}
}

func wsHandle(conn *websocket.Conn) {
	var err error
	var uid uint64
	defer func() {
		if err != nil {
			log.Fatalln("websocket connect is failed: ", err)
		} else {
			broads.Delete(uid)
			log.Println("websocket closed: ", conn.RemoteAddr().String())
		}
		conn.Close()
	}()
	// 创建玩家
	pl := player.NewPlayer(conn)
	if err = pl.Auth(); err != nil {
		return
	}
	// 订阅消息
	uid = pl.GetUID()
	self := cluster.GetSelfServerNode()
	key := cluster.GetPlayerChannel(self.ServerType, self.ServerID, uid)
	if err = network.Subscribe(key, pl.NatsHandle); err != nil {
		return
	}
	broads.Add(uid, pl.NatsHandle)
	// 循环接受消息
	log.Println("websocket connected...", conn.RemoteAddr().String())
	pl.LoopRead()
}
