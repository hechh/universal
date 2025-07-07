package message

import (
	"universal/common/pb"
	"universal/framework/rpc"
)

func Init() {
	rpc.Register(pb.NodeType_Gate, 0, "Player", "SendToClient", "LoginSuccess")
	rpc.Register(pb.NodeType_Gate, 1, "PlayerMgr", "Kick")
}
