package frame

import (
	"encoding/binary"
	"poker_server/common/config/repository/open_api_config"
	"poker_server/common/pb"
	"poker_server/framework/library/uerror"
)

type WsPacket struct {
	Cmd     uint32 // 消息id
	Uid     uint64 // 玩家uid
	RouteId uint64 // 路由 id
	Seq     uint32 // 序列号
	Version uint32 // 版本号
	Extra   uint32 // 扩展字段
	Body    []byte // 消息体
}

type Frame struct{}

// 最大包长限制
func (d *Frame) GetSize(pac *pb.Packet) int {
	return 32 + len(pac.Body)
}

// 解码数据包
func (d *Frame) Decode(buf []byte, msg *pb.Packet) error {
	pack := &WsPacket{}
	pos := 0
	pack.Cmd = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	pack.Uid = binary.BigEndian.Uint64(buf[pos:])
	pos += 8
	pack.RouteId = binary.BigEndian.Uint64(buf[pos:])
	pos += 8
	pack.Seq = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	pack.Version = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	pack.Extra = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	pack.Body = buf[pos:]
	// 确当路由节点
	cfg := open_api_config.MGetCmd(pack.Cmd - pack.Cmd%2)
	if cfg == nil {
		return uerror.New(1, -1, "cmd未注册: %d", pack.Cmd)
	}
	if msg.Head == nil {
		msg.Head = &pb.Head{}
	}
	msg.Head.DstNodeType = pb.NodeType(cfg.ServerType)
	msg.Head.Cmd = cfg.Cmd
	msg.Head.Id = pack.Uid
	msg.Head.RouteId = pack.RouteId
	msg.Head.Seq = pack.Seq
	msg.Head.ActorName = cfg.ActorName
	msg.Head.FuncName = cfg.FuncName
	msg.Body = pack.Body
	return nil
}

func (d *Frame) Encode(pack *pb.Packet, buf []byte) error {
	// 组包
	pos := 0
	binary.BigEndian.PutUint32(buf[pos:], pack.Head.Cmd)
	pos += 4
	binary.BigEndian.PutUint64(buf[pos:], pack.Head.Id)
	pos += 8
	binary.BigEndian.PutUint64(buf[pos:], pack.Head.RouteId)
	pos += 8
	binary.BigEndian.PutUint32(buf[pos:], pack.Head.Seq)
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], 0)
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], 0)
	pos += 4
	copy(buf[pos:], pack.Body)
	return nil
}
