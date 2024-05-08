package gate

import (
	"log"
	"universal/common/pb"
	"universal/framework"
	"universal/framework/actor"
	"universal/framework/fbasic"
	"universal/framework/network"
	"universal/framework/notify"

	"github.com/spf13/cast"
)

type User struct {
	uid    string
	client *network.SocketClient
}

func NewUser(client *network.SocketClient) *User {
	return &User{client: client}
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
