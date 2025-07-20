package framework

import (
	"poker_server/common/pb"
	"poker_server/common/yaml"
	"poker_server/framework/actor"
	"poker_server/framework/cluster"
	"poker_server/framework/domain"
	"poker_server/framework/request"
	"poker_server/library/mlog"
	"poker_server/library/pprof"
	"poker_server/library/snowflake"
	"strings"
	"sync/atomic"
)

var (
	envType pb.EnvType
)

func Init(nn *pb.Node, srvCfg *yaml.NodeConfig, cfg *yaml.Config) error {
	if err := snowflake.Init(nn); err != nil {
		return err
	}

	// 初始化集群模块
	if err := cluster.Init(nn, srvCfg, cfg); err != nil {
		return err
	}

	// 初始化全局变量
	switch strings.ToLower(cfg.Common.Env) {
	case "develop":
		envType = pb.EnvType_EnvTypeDevelop
	case "release":
		envType = pb.EnvType_EnvTypeRelease
	}

	// 依赖库初始化
	pprof.Init("", srvCfg.Port+10000)
	return nil
}

func GetEnvType() pb.EnvType {
	return envType
}

func StopAutoSendToClient(head *pb.Head) {
	atomic.AddUint32(&head.Reference, 1)
}

func DefaultHandler(head *pb.Head, buf []byte) {
	head.ActorName = head.Dst.ActorName
	head.FuncName = head.Dst.FuncName
	head.ActorId = head.Dst.ActorId
	if err := actor.Send(head, buf); err != nil {
		mlog.Errorf("跨服务调用actor错误: %v head: %v", err, head)
	}
}

func ParseActorFunc(str string) (actorName string, funcName string) {
	if pos := strings.Index(str, "."); pos > 0 {
		actorName = str[:pos]
		funcName = str[pos+1:]
	} else {
		actorName = str
	}
	return
}

func NewSrcRouter(actorId uint64, actorName string, funcs ...string) *pb.NodeRouter {
	self := cluster.GetSelf()
	rr := request.NewNodeRouter(self.Type, actorId, actorName, funcs...)
	if rr != nil {
		rr.NodeId = self.Id
	}
	return rr
}

func NewActorRouter(act domain.IActor, funcs ...string) *pb.NodeRouter {
	self := cluster.GetSelf()
	rr := request.NewNodeRouter(self.Type, act.GetId(), act.GetActorName(), funcs...)
	if rr != nil {
		rr.NodeId = self.Id
	}
	return rr
}

func NewGateRouter(actorId uint64, actorName string, funcs ...string) *pb.NodeRouter {
	return request.NewNodeRouter(pb.NodeType_NodeTypeGate, actorId, actorName, funcs...)
}

func NewGameRouter(actorId uint64, actorName string, funcs ...string) *pb.NodeRouter {
	return request.NewNodeRouter(pb.NodeType_NodeTypeGame, actorId, actorName, funcs...)
}

func NewDbRouter(actorId uint64, actorName string, funcs ...string) *pb.NodeRouter {
	return request.NewNodeRouter(pb.NodeType_NodeTypeDb, actorId, actorName, funcs...)
}

func NewBuilderRouter(actorId uint64, actorName string, funcs ...string) *pb.NodeRouter {
	return request.NewNodeRouter(pb.NodeType_NodeTypeBuilder, actorId, actorName, funcs...)
}

func NewRoomRouter(actorId uint64, actorName string, funcs ...string) *pb.NodeRouter {
	return request.NewNodeRouter(pb.NodeType_NodeTypeRoom, actorId, actorName, funcs...)
}

func NewMatchRouter(actorId uint64, actorName string, funcs ...string) *pb.NodeRouter {
	return request.NewNodeRouter(pb.NodeType_NodeTypeMatch, actorId, actorName, funcs...)
}

func NewGmRouter(actorId uint64, actorName string, funcs ...string) *pb.NodeRouter {
	return request.NewNodeRouter(pb.NodeType_NodeTypeGm, actorId, actorName, funcs...)
}

func SwapToDb(head *pb.Head, id uint64, actorName string, funcs ...string) *pb.Head {
	head.Src, head.Dst = head.Dst, head.Src
	head.Dst = NewDbRouter(id, actorName, funcs...)
	return head
}

func SwapToGate(head *pb.Head, id uint64, actorName string, funcs ...string) *pb.Head {
	head.Src, head.Dst = head.Dst, head.Src
	head.Dst = NewGateRouter(id, actorName, funcs...)
	return head
}

func SwapToGame(head *pb.Head, id uint64, actorName string, funcs ...string) *pb.Head {
	head.Src, head.Dst = head.Dst, head.Src
	head.Dst = NewGameRouter(id, actorName, funcs...)
	return head
}

func SwapToBuilder(head *pb.Head, id uint64, actorName string, funcs ...string) *pb.Head {
	head.Src, head.Dst = head.Dst, head.Src
	head.Dst = NewBuilderRouter(id, actorName, funcs...)
	return head
}

func SwapToRoom(head *pb.Head, id uint64, actorName string, funcs ...string) *pb.Head {
	head.Src, head.Dst = head.Dst, head.Src
	head.Dst = NewRoomRouter(id, actorName, funcs...)
	return head
}

func SwapToMatch(head *pb.Head, id uint64, actorName string, funcs ...string) *pb.Head {
	head.Src, head.Dst = head.Dst, head.Src
	head.Dst = NewMatchRouter(id, actorName, funcs...)
	return head
}

func SwapToGm(head *pb.Head, id uint64, actorName string, funcs ...string) *pb.Head {
	head.Src, head.Dst = head.Dst, head.Src
	head.Dst = NewGmRouter(id, actorName, funcs...)
	return head
}
