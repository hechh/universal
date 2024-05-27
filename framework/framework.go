package framework

import (
	"os"
	"os/signal"
	"syscall"
	"universal/common/pb"
	"universal/framework/cluster"
	"universal/framework/common/fbasic"
	"universal/framework/common/ulog"
	"universal/framework/network"
	"universal/framework/packet"

	"google.golang.org/protobuf/proto"
)

const (
	API_CODE = 1000000
)

func ApiCodeToServerType(val int32) (typ pb.ServerType) {
	switch val / API_CODE {
	case 1:
		typ = pb.ServerType_GATE
	case 2:
		typ = pb.ServerType_GAME
	default:
		typ = pb.ServerType_NONE
	}
	return
}

// 发送到其他服务
func SendTo(sendType pb.SendType, apiCode int32, uid uint64, req proto.Message, params ...interface{}) error {
	self := cluster.GetSelfServerNode()
	head := &pb.PacketHead{
		SendType:      sendType,
		SrcServerType: self.ServerType,
		SrcServerID:   self.ServerID,
		DstServerType: ApiCodeToServerType(apiCode),
		ApiCode:       apiCode,
		Time:          fbasic.GetNow(),
		UID:           uid,
	}
	// 路由
	if err := cluster.Dispatcher(head); err != nil {
		return err
	}
	// 获取订阅key
	key, err := cluster.GetHeadChannel(head)
	if err != nil {
		return err
	}
	return network.PublishReq(key, head, req, params...) //(head, req)
}

// 发送客户端
func SendToClient(sendType pb.SendType, apiCode int32, uid uint64, rsp proto.Message, params ...interface{}) error {
	self := cluster.GetSelfServerNode()
	head := &pb.PacketHead{
		SendType:      sendType,
		SrcServerType: self.ServerType,
		SrcServerID:   self.ServerID,
		DstServerType: pb.ServerType_GATE,
		ApiCode:       apiCode + 1,
		Time:          fbasic.GetNow(),
		UID:           uid,
	}
	// 路由
	if err := cluster.Dispatcher(head); err != nil {
		return err
	}
	// 获取订阅key
	key, err := cluster.GetHeadChannel(head)
	if err != nil {
		return err
	}
	return network.PublishRsp(key, head, rsp, params...)
}

// 设置actor处理
func ActorHandle(ctx *fbasic.Context, buf []byte) func() {
	return func() {
		// 调用接口
		rsp, err := packet.Call(ctx, buf)
		if err != nil {
			ulog.Error(1, "head: %v, error: %v", ctx.PacketHead, err)
			return
		}
		// 设置返回信息
		head := ctx.PacketHead
		head.SeqID++
		head.ApiCode++
		head.SrcServerType, head.DstServerType = head.DstServerType, head.SrcServerType
		head.SrcServerID, head.DstServerID = head.DstServerID, head.SrcServerID
		// 获取订阅key
		key, err := cluster.GetHeadChannel(head)
		if err != nil {
			ulog.Error(1, "head: %v, error: %v", head, err)
			return
		}
		// 发送
		if err := network.PublishRsp(key, head, rsp); err != nil {
			ulog.Error(1, "head: %v, key: %s, rsp: %v, error: %v", head, key, rsp, err)
			return
		}
	}
}

func SignalHandle(f func(os.Signal)) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		for sig := range ch {
			f(sig)
		}
	}()
}

// 初始化
func Init(serverType pb.ServerType, addr string, etcds []string, natsUrl string) error {
	// 连接etcd
	if err := cluster.Init(etcds); err != nil {
		return err
	}
	// 进行服务发现
	if err := cluster.Discovery(serverType, addr); err != nil {
		return err
	}
	// 初始化消息中间件
	if err := network.InitNats(natsUrl); err != nil {
		return err
	}
	return nil
}
