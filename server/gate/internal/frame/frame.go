package frame

import (
	"encoding/binary"
	"universal/common/pb"
	"universal/framework"
	"universal/framework/handler"
)

type Frame struct{}

func (d *Frame) GetSize(pac *pb.Packet) int {
	return 32 + len(pac.Body)
}

// 解码数据包
func (d *Frame) Decode(buf []byte, msg *pb.Packet) error {
	pos := 0
	cmdVal := binary.BigEndian.Uint32(buf[pos:])
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
		Src: framework.NewNodeRouterByUid(pb.NodeType_Gate, uid, 0, "Player.SendToClient"),
		Dst: handler.NewNodeRouterByCmd(cmdVal, routeId, 0),
		Cmd: cmdVal,
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
