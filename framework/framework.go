package framework

import (
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/actor"
	"universal/framework/cluster"
	"universal/framework/internal/handler"
	"universal/framework/recycle"
	"universal/library/mlog"
	"universal/library/pprof"
	"universal/library/snowflake"
)

func Init(nn *pb.Node, srvCfg *yaml.NodeConfig, cfg *yaml.Config) error {
	if err := snowflake.Init(nn); err != nil {
		return err
	}
	// 初始化集群模块
	if err := cluster.Init(nn, srvCfg, cfg); err != nil {
		return err
	}
	initOther(cfg, nn)
	return nil
}

func InitDefault(nn *pb.Node, srvCfg *yaml.NodeConfig, cfg *yaml.Config) error {
	if err := Init(nn, srvCfg, cfg); err != nil {
		return err
	}
	// 初始化集群模块
	if err := cluster.SetBroadcastHandler(defaultHandler); err != nil {
		return err
	}
	if err := cluster.SetSendHandler(defaultHandler); err != nil {
		return err
	}
	if err := cluster.SetReplyHandler(defaultHandler); err != nil {
		return err
	}
	return nil
}

func initOther(cfg *yaml.Config, nn *pb.Node) {
	pprof.Init(cfg.Common.Pprof, nn.Ip, int(nn.Port))
	recycle.Init()
}

func defaultHandler(head *pb.Head, buf []byte) {
	head.ActorName, head.FuncName = handler.GetActorFunc(head.Dst.ActorFunc)
	if head.Dst.ActorId <= 0 {
		head.ActorId = head.Dst.RouterId
	} else {
		head.ActorId = head.Dst.ActorId
	}
	if err := actor.Send(head, buf); err != nil {
		mlog.Errorf("跨服务调用actor错误: %v", err)
	}
}

func NewNodeRouterByUid(nt pb.NodeType, uid uint64, actorFunc string) *pb.NodeRouter {
	return &pb.NodeRouter{
		NodeType:  nt,
		RouterId:  handler.GenRouterId(uid, uint64(pb.RouterType_UID)),
		ActorFunc: handler.GetActorFuncId(actorFunc),
	}
}

func NewNodeRouterByRoomId(nt pb.NodeType, roomId uint64, actorFunc string) *pb.NodeRouter {
	return &pb.NodeRouter{
		NodeType:  nt,
		RouterId:  handler.GenRouterId(roomId, uint64(pb.RouterType_ROOM_ID)),
		ActorFunc: handler.GetActorFuncId(actorFunc),
	}
}

func NewNodeRouterByRandomId(nt pb.NodeType, id uint64, actorFunc string) *pb.NodeRouter {
	return &pb.NodeRouter{
		NodeType:  nt,
		RouterId:  handler.GenRouterId(id, uint64(pb.RouterType_RANDOM_ID)),
		ActorFunc: handler.GetActorFuncId(actorFunc),
	}
}

func CopyTo(head *pb.Head, dst *pb.NodeRouter) *pb.Head {
	newSrc := *head.Src
	return &pb.Head{
		SendType: head.SendType,
		Src:      &newSrc,
		Dst:      dst,
		Uid:      head.Uid,
		Seq:      head.Seq,
		Cmd:      head.Cmd,
	}
}
