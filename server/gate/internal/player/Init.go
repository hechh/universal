package player

import (
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/cluster"
	"universal/server/gate/internal/token"
)

var (
	mgr = &PlayerMgr{}
)

func Init(cfg *yaml.Config, srvcfg *yaml.NodeConfig) error {
	token.Init(cfg.Common.SecretKey)
	if err := cluster.SetBroadcastHandler(defaultHandler); err != nil {
		return err
	}
	if err := cluster.SetSendHandler(defaultHandler); err != nil {
		return err
	}
	if err := cluster.SetReplyHandler(defaultHandler); err != nil {
		return err
	}
	return mgr.Init(srvcfg.Ip, srvcfg.Port)
}

func Close() {
	mgr.Close()
}

func defaultHandler(head *pb.Head, body []byte) {
	/*
		rpc.ParseNodeRouter(head, "Player.SendToClient")
		if err := actor.Send(head, body); err != nil {
			mlog.Errorf("Actor消息转发失败: %v", err)
		}
	*/
}
