package gate

import (
	"flag"
	"fmt"
	"strings"
	"universal/common/dao"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework"
	"universal/library/mlog"
	"universal/library/signal"
)

func main() {
	var filename string
	var id int
	flag.StringVar(&filename, "config", "config.yaml", "游戏配置")
	flag.IntVar(&id, "id", 1, " 节点id")
	flag.Parse()

	// 加载游戏配置
	node := &pb.Node{
		Name: strings.ToLower(pb.NodeType_Gate.String()),
		Type: pb.NodeType_Gate,
		Id:   int32(id),
	}
	cfg, err := yaml.LoadConfig(filename, node)
	if err != nil {
		panic(fmt.Sprintf("游戏配置加载失败: %v", err))
	}

	// 初始化日志库
	if err := mlog.Init(cfg.Cluster[node.Name]); err != nil {
		panic(fmt.Sprintf("日志库初始化失败: %v", err))
	}

	// 初始化redis
	mlog.Infof("初始化redis配置")
	if err := dao.InitRedis(cfg.Redis); err != nil {
		panic(fmt.Sprintf("redis初始化失败: %v", err))
	}

	// 初始化框架服务
	mlog.Infof("初始化redis配置")
	if err := framework.Init(node, cfg); err != nil {
		panic(fmt.Sprintf("框架初始化失败：%v", err))
	}

	// 服务退出
	signal.SignalNotify(func() {
		// todo
	})
}

/*
// 处理返回客户端的消息
func sendHandler(head *pb.Head, body []byte) {
	// 发送到客户端
	mlog.Debugf("receive send msg head:%v", head)
	if len(head.ActorName) <= 0 || len(head.FuncName) <= 0 {
		playerMgr.SendToClient(head, body)
		return
	}
	// 发送到指定Actor
	if err := playerMgr.Send(head, body); err != nil {
		mlog.Errorf("Actor消息转发失败: %v", err)
	}
}

// 处理返回客户端的消息
func broadcastHandler(head *pb.Head, body []byte) {
	// 发送到客户端
	mlog.Debugf("receive broadcast msg head:%v", head)
	if len(head.ActorName) <= 0 || len(head.FuncName) <= 0 {
		playerMgr.BroadcastToClient(head, body)
		return
	}

	// 广播到所有Actor
	if err := playerMgr.Send(head, body); err != nil {
		mlog.Errorf("Actor消息转发失败: %v", err)
	}
}
*/
