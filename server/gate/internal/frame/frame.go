package frame

import (
	"encoding/binary"
	"universal/common/pb"
	"universal/framework/rpc"
)

type Frame struct{}

func (d *Frame) GetSize(pac *pb.Packet) int {
	return 32 + len(pac.Body)
}

// 解码数据包
func (d *Frame) Decode(buf []byte, msg *pb.Packet) error {
	pos := 0
	cmd := binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	uid := binary.BigEndian.Uint64(buf[pos:])
	pos += 8
	routeId := binary.BigEndian.Uint64(buf[pos:])
	pos += 8
	seq := binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	//	version := binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	//	extra := binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	msg.Body = buf[pos:]
	// 确当路由节点
	msg.Head = &pb.Head{
		Uid: uid,
		Src: rpc.NewNodeRouter(pb.NodeType_NodeTypeGate, "Player.SendToClient", uid),
		Dst: rpc.NewNodeRouterByCmd(pb.CMD(cmd), routeId),
		Cmd: cmd,
		Seq: seq,
	}
	return nil
}

// 组包
func (d *Frame) Encode(pack *pb.Packet, buf []byte) error {
	pos := 0
	binary.BigEndian.PutUint32(buf[pos:], pack.Head.Cmd) // cmd
	pos += 4
	binary.BigEndian.PutUint64(buf[pos:], pack.Head.Uid) // uid
	pos += 8
	binary.BigEndian.PutUint64(buf[pos:], pack.Head.Src.ActorId) // router id
	pos += 8
	binary.BigEndian.PutUint32(buf[pos:], pack.Head.Seq) // seq
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], 0) // version
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], 0) // extra
	pos += 4
	copy(buf[pos:], pack.Body)
	return nil
}

/*
type WsPacket struct {
	Cmd     uint32 // 消息id
	Uid     uint64 // 玩家uid
	RouteId uint64 // 路由 id
	Seq     uint32 // 序列号
	Version uint32 // 版本号
	Extra   uint32 // 扩展字段
	Body    []byte // 消息体
}
*/
