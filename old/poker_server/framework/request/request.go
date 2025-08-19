package request

import (
	"fmt"
	"poker_server/common/pb"
	"poker_server/library/util"
	"strings"
)

var (
	cmds   = make(map[pb.CMD]*CmdInfo)
	actors = make(map[pb.NodeType]map[string]uint32)
)

type CmdInfo struct {
	cmd pb.CMD
	nt  pb.NodeType
	rt  uint32
	an  string
	fn  string
}

func RegisterActor(nt pb.NodeType, actorName string, id uint32) {
	if _, ok := actors[nt]; !ok {
		actors[nt] = make(map[string]uint32)
	}
	actors[nt][actorName] = id
}

func RegisterCmd(nt pb.NodeType, cmd pb.CMD, actorFunc string) {
	pos := strings.Index(actorFunc, ".")
	if pos <= 0 {
		panic(fmt.Sprintf("Actor接口注册错误%s", actorFunc))
	}
	actorName := actorFunc[:pos]
	funcName := actorFunc[pos+1:]
	if _, ok := actors[nt]; !ok {
		panic(fmt.Sprintf("%s未注册%s", nt.String(), actorFunc))
	}
	id, ok := actors[nt][actorName]
	if !ok {
		panic(fmt.Sprintf("%s未注册%s", nt.String(), actorName))
	}
	cmds[cmd] = &CmdInfo{cmd, nt, id, actorFunc, funcName}
}

func NewCmdRouter(cmd uint32, actorId uint64) *pb.NodeRouter {
	if info, ok := cmds[pb.CMD(cmd)]; !ok {
		return &pb.NodeRouter{
			NodeType:   info.nt,
			ActorName:  info.an,
			FuncName:   info.fn,
			ActorId:    actorId,
			RouterType: info.rt,
		}
	}
	return nil
}

func NewNodeRouter(nt pb.NodeType, actorId uint64, actorName string, funcs ...string) *pb.NodeRouter {
	infos, ok := actors[nt]
	if !ok {
		return nil
	}
	actorType, ok := infos[actorName]
	if !ok {
		return nil
	}
	return &pb.NodeRouter{
		NodeType:   nt,
		ActorName:  actorName,
		FuncName:   util.Index[string](funcs, 0, ""),
		ActorId:    actorId,
		RouterType: actorType,
	}
}
