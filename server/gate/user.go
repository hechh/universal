package gate

import (
	"log"
	"universal/common/pb"
	"universal/framework/actor"
	"universal/framework/cluster"
	"universal/framework/fbasic"
	"universal/framework/network"
	"universal/framework/notify"

	"github.com/spf13/cast"
	"google.golang.org/protobuf/proto"
)

type User struct {
	uid    string
	client *network.SocketClient
}

func NewUser(uid uint64, client *network.SocketClient) (*User, error) {
	ret := &User{cast.ToString(uid), client}
	// 设置nats消息处理
	self := cluster.GetLocalClusterNode()
	err := notify.Subscribe(fbasic.GetPlayerChannel(self.ClusterType, self.ClusterID, uid), ret.NatsHandle)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// nats消息处理,(point_send)
func (d *User) NatsHandle(pac *pb.Packet) {

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
		// 更新head路由信息
		head := pac.Head
		if err := cluster.Dispatcher(head); err != nil {
			log.Fatalln(head, string(pac.Buff))
			continue
		}
		// 转发
		if head.SrcClusterType == head.DstClusterType {
			actor.Send(cast.ToString(head.UID), pac)
			continue
		}
		// 转发到nats
		if err := notify.Publish(pac); err != nil {
			log.Println(err)
		}
	}
}
