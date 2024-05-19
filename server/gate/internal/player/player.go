package player

import (
	"log"
	"net"
	"universal/common/pb"
	"universal/framework"
	"universal/framework/actor"
	"universal/framework/cluster"
	"universal/framework/common/fbasic"
	"universal/framework/common/socket"
	"universal/framework/common/uerror"
	"universal/framework/network"
	"universal/framework/packet"

	"github.com/spf13/cast"
	"google.golang.org/protobuf/proto"
)

type Player struct {
	uid    uint64
	client *socket.Socket
}

func NewPlayer(ws net.Conn) *Player {
	return &Player{
		client: socket.NewSocket(&socket.Frame{}, ws),
	}
}

func (d *Player) GetUID() uint64 {
	return d.uid
}

func (d *Player) Auth() error {
	pac, err := d.Read()
	if err != nil {
		return err
	}
	// 判断第一个请求是否为登陆认证包
	d.uid = pac.Head.UID
	head := pac.Head
	if head.ApiCode != int32(pb.ApiCode_GATE_LOGIN_REQUEST) {
		return uerror.NewUError(1, -1, pb.ApiCode_GATE_LOGIN_REQUEST, head.ApiCode)
	}
	// 登陆认证
	rsp, err := packet.Call(fbasic.NewContext(head), pac.Buff)
	if err != nil {
		return err
	}
	// 返回认证结果
	if err := d.SendRsp(head, rsp); err != nil {
		return err
	}
	// 判断是否成功
	if head := rsp.(fbasic.IRpcHead).GetHead(); head.Code > 0 {
		return uerror.NewUError(1, pb.ErrorCode(head.Code), "gate login auth is failed")
	}
	return nil
}

func (d *Player) Read() (*pb.Packet, error) {
	buf, err := d.client.ReadBytes()
	if err != nil {
		return nil, err
	}
	ret := &pb.Packet{}
	if err := proto.Unmarshal(buf, ret); err != nil {
		return nil, uerror.NewUError(1, -1, err)
	}
	return ret, nil
}

func (d *Player) Send(pac *pb.Packet) error {
	buf, err := proto.Marshal(pac)
	if err != nil {
		return uerror.NewUError(1, -1, err)
	}
	return d.client.SendBytes(buf)
}

func (d *Player) SendRsp(head *pb.PacketHead, item proto.Message, params ...interface{}) error {
	pac, err := network.RspToPacket(head, item, params...)
	if err != nil {
		return err
	}
	return d.Send(pac)
}

// nats消息处理
func (d *Player) NatsHandle(pac *pb.Packet) {
	head := pac.Head
	switch head.ApiCode & 0x01 {
	case 1:
		if err := d.Send(pac); err != nil {
			log.Fatalln(err)
		}
	case 0:
		actor.Send(cast.ToString(d.uid), framework.ActorHandle, pac)
	}
	log.Println("player nats handler finished: ", head)
}

// 循环读取客户端请求
func (d *Player) LoopRead() {
	for {
		// 接受数据包
		pac, err := d.Read()
		if err != nil {
			log.Println("Read: ", err)
			return
		}
		log.Println(pac, "---->", err)
		// 更新head路由信息
		head := pac.Head
		if err := cluster.Dispatcher(head); err != nil {
			log.Println("Dispatcher: ", head, string(pac.Buff))
			continue
		}
		// 转发
		if head.SrcServerType == head.DstServerType {
			actor.Send(cast.ToString(head.UID), framework.ActorHandle, pac)
			continue
		}
		// 转发到nats
		if key, err := cluster.GetHeadChannel(head); err != nil {
			log.Println("nats: ", err)
		} else if err = network.Publish(key, pac); err != nil {
			log.Println("Puslish: ", err)
		}
	}
}
