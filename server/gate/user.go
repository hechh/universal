package gate

import (
	"log"
	"universal/common/pb"
	"universal/framework/actor"
	"universal/framework/cluster"
	"universal/framework/fbasic"
	"universal/framework/network"
	"universal/framework/packet"

	"github.com/spf13/cast"
	"google.golang.org/protobuf/proto"
)

type User struct {
	uid    string
	client *network.SocketClient
}

func NewUser(client *network.SocketClient) *User {
	return &User{client: client}
}

// nats消息处理,(point_send)
func (d *User) NatsHandle(pac *pb.Packet) {
	head := pac.Head
	switch head.Status {
	case pb.StatusType_REQUEST:
		local := cluster.GetLocalClusterNode()
		if local.ClusterType == head.DstClusterType {
			actor.Send(d.uid, pac)
		} else {
			// 转发到nats
			if err := cluster.Publish(pac); err != nil {
				log.Println(err)
			}
		}
	case pb.StatusType_RESPONSE:
		if err := d.client.Send(pac); err != nil {
			log.Fatal(err)
		}
	}
}

// 认证
func (d *User) Auth() (flag bool) {
	var pac *pb.Packet
	var err error
	if pac, err = d.client.Read(); err != nil {
		log.Println("auth failed: ", err)
		return
	}
	d.uid = cast.ToString(pac.Head.UID)
	// 判断第一个请求是否为登陆认证包
	if pac.Head.ApiCode != int32(pb.ApiCode_GATE_LOGIN_REQUEST) {
		log.Println("GateLoginRequest is expected", pac.Head)
		return
	}
	// 登陆认证
	rsp, err := packet.Call(fbasic.NewDefaultContext(pac.Head), pac.Buff)
	if err != nil {
		log.Println("ApiCode not supported: ", err)
		return
	}
	// 返回认证结果
	if err := d.Reply(pac.Head, rsp); err != nil {
		log.Println("auth reply is failed: ", err)
		return
	}
	// 判断是否成功
	return rsp.(fbasic.IProto).GetHead().Code == 0
}

// 回复客户端
func (d *User) Reply(head *pb.PacketHead, rsp proto.Message) error {
	pp, err := fbasic.RspToPacket(head, rsp)
	if err != nil {
		return err
	}
	return d.client.Send(pp)
}

// 循环读取客户端请求
func (d *User) LoopRead() {
	for {
		// 接受数据包
		pac, err := d.client.Read()
		if err != nil {
			log.Fatal(err)
			return
		}
		// 本地信息
		head := pac.Head
		local := cluster.GetLocalClusterNode()
		head.SrcClusterType = local.ClusterType
		head.SrcClusterID = local.ClusterID
		// 设置头信息
		head.DstClusterType = fbasic.ApiCodeToClusterType(head.ApiCode)
		if head.DstClusterType == pb.ClusterType_NONE {
			log.Println(head.ApiCode, head)
			continue
		}
		// 转发
		if head.SrcClusterType == head.DstClusterType {
			actor.Send(cast.ToString(pac.Head.UID), pac)
		} else {
			// 转发到nats
			if err := cluster.Publish(pac); err != nil {
				log.Println(err)
			}
		}
	}
}
