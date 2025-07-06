package message

import (
	"universal/common/pb"
	"universal/framework/rpc"
)

func Init() {
	rpc.Register(pb.NodeType_NodeTypeGate, 0, "Player", "SendToClient", "LoginSuccess")
	rpc.Register(pb.NodeType_NodeTypeGate, 1, "PlayerMgr", "Kick")
}
