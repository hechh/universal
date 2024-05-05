package gate

import (
	"log"
	"universal/common/pb"
	"universal/framework/actor"
	"universal/framework/cluster"
	"universal/framework/fbasic"
	"universal/framework/network"
	"universal/framework/packet"

	"golang.org/x/net/websocket"
)

type User struct {
	uid    string
	client *network.SocketClient
}

func NewUser(conn *websocket.Conn) (*User, error) {
	log.Println("websocket connected, ", conn.RemoteAddr().String())
	client := network.NewSocketClient(conn)

	// 登陆认证(第一个请求，一定是登陆认证请求)
	pac, err := client.Read()
	if err != nil {
		return nil, err
	}
	// 判断第一个请求是否为登陆认证包
	if pac.Head.ApiCode != int32(pb.ApiCode_GATE_LOGIN_REQUEST) {
		return nil, fbasic.NewUError(1, pb.ErrorCode_GateLoginRequestExpected, pac.Head.ApiCode)
	}
	// 返回认证结果
	if err := client.Send(pac); err != nil {
		return nil, err
	}
	return &User{client: client}, nil
}

func (d *User) LoginRequest(client *network.SocketClient) error {
	// 登陆认证(第一个请求，一定是登陆认证请求)
	pac, err := client.Read()
	if err != nil {
		return err
	}
	// 判断第一个请求是否为登陆认证包
	if pac.Head.ApiCode != int32(pb.ApiCode_GATE_LOGIN_REQUEST) {
		return fbasic.NewUError(1, pb.ErrorCode_GateLoginRequestExpected, pac.Head.ApiCode)
	}
	// 登陆认证
	result, err := packet.Call(fbasic.NewContext(pac.Head, nil), pac.Buff)
	if err != nil {
		return err
	}
	// 返回认证结果
	return client.Send(result)
}

// nats消息处理
func (d *User) NatsHandle(pac *pb.Packet) {
	/*
		if err := d.client.Send(pac); err != nil {
			log.Fatal(err)
		}
	*/
}

func (d *User) LoopRead() {
	for {
		pac, err := d.client.Read()
		if err != nil {
			log.Fatal(err)
			break
		}

		// 转发
		if err := d.dispatcher(pac); err != nil {
			log.Fatal(err)
			continue
		}
	}
}

func (d *User) dispatcher(pac *pb.Packet) error {
	// 设置头信息
	head := pac.Head
	head.DstClusterType = fbasic.ApiCodeToClusterType(head.ApiCode)
	if head.DstClusterType == pb.ClusterType_NONE {
		return fbasic.NewUError(1, pb.ErrorCode_NotSupported, head.ApiCode)
	}
	local := cluster.GetLocalClusterNode()
	head.SrcClusterType = local.ClusterType
	head.SrcClusterID = local.ClusterID
	// 转发
	if head.SrcClusterType == head.DstClusterType {
		actor.Send(d.uid, pac)
		return nil
	}
	// 转发到nats
	return cluster.Publish(pac)
}
