package gate

import (
	"log"
	"universal/common/pb"
	"universal/framework"
	"universal/framework/actor"
	"universal/framework/cluster"
	"universal/framework/fbasic"
	"universal/framework/network"
	"universal/framework/notify"
	"universal/framework/packet"

	"github.com/spf13/cast"
)

type User struct {
	uid    string
	client *network.SocketClient
}

func NewUser(client *network.SocketClient) *User {
	return &User{client: client}
}

func (d *User) Init() error {
	self := cluster.GetLocalClusterNode()
	uid := cast.ToUint64(d.uid)
	return notify.Subscribe(fbasic.GetPlayerChannel(self.ClusterType, self.ClusterID, uid), d.NatsHandle)
}

// nats消息处理
func (d *User) NatsHandle(pac *pb.Packet) {
	head := pac.Head
	switch head.ApiCode & 0x01 {
	case 1:
		if err := d.client.Send(pac); err != nil {
			log.Fatalln(err)
		}
	case 0:
		actor.Send(d.uid, pac)
	}
	log.Println("user nats handler finished: ", head)
}

func (d *User) Auth() (flag bool) {
	pac, err := d.client.Read()
	if err != nil {
		log.Println("auth failed: ", err)
		return
	}
	head := pac.Head
	d.uid = cast.ToString(head.UID)
	// 判断第一个请求是否为登陆认证包
	if head.ApiCode != int32(pb.ApiCode_GATE_LOGIN_REQUEST) {
		log.Println("GateLoginRequest is expected", head)
		return
	}
	// 登陆认证
	rsp, err := packet.Call(fbasic.NewDefaultContext(head), pac.Buff)
	if err != nil {
		log.Println("ApiCode not supported: ", err)
		return
	}
	// 返回认证结果
	if pp, err := fbasic.RspToPacket(head, rsp); err != nil {
		log.Println("RspToPacket is failed: ", err)
		return
	} else if err := d.client.Send(pp); err != nil {
		log.Println("auth reply is failed: ", err)
		return
	}
	// 判断是否成功
	flag = rsp.(fbasic.IProto).GetHead().Code == 0
	return
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
		// 更新head路由信息
		head := pac.Head
		if err := framework.Dispatcher(head); err != nil {
			log.Fatalln(head, string(pac.Buff))
			continue
		}
		// 转发
		if head.SrcClusterType == head.DstClusterType {
			actor.Send(cast.ToString(head.UID), pac)
			continue
		}
		// 转发到nats
		if key, err := fbasic.GetHeadChannel(head); err != nil {
			log.Fatalln(err)
		} else if err = notify.Publish(key, pac); err != nil {
			log.Fatalln(err)
		}
	}
}
