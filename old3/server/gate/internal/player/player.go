package player

import (
	"time"
	"universal/common/pb"
	"universal/framework"
	"universal/framework/actor"
	"universal/framework/domain"
	"universal/framework/network"
	"universal/framework/token"
	"universal/library/mlog"
	"universal/library/uerror"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

type Player struct {
	actor.Actor
	inet       domain.INet // 网络连接
	createTime int64       // 创建时间
	expireTime int64       // 过期时间
	updateTime int64       // 更新时间
}

func NewPlayer(conn *websocket.Conn) *Player {
	p := &Player{}
	p.inet = network.NewSocket(conn, 1024*1024)
	p.inet.Register(&Frame{})
	return p
}

func (p *Player) Login() error {
	// 第一个包一定是登录认证包
	p.inet.SetReadExpire(5)
	pack := &pb.Packet{}
	if err := p.inet.Read(pack); err != nil {
		return err
	}
	p.inet.SetReadExpire(0)
	req := &pb.LoginRequest{}
	if err := proto.Unmarshal(pack.Body, req); err != nil {
		return err
	}
	// 解析token
	tt, err := token.ParseToken(req.Token)
	if err != nil {
		return err
	}
	// 设置玩家ID
	p.Actor.SetId(tt.Uid)
	p.createTime = time.Now().Unix()
	p.updateTime = time.Now().Unix()
	p.expireTime = 15 * 60
	pack.Head.Cmd++
	pack.Head.Id = tt.Uid
	// 发送登录成功包
	return p.inet.Write(&pb.Packet{Head: pack.Head})
}

// 消息分发处理( 接受 websocket 传过来的消息)
func (p *Player) Dispatcher() {
	for {
		// 从客户端持续接受包消息
		pack := &pb.Packet{}
		if err := p.inet.Read(pack); err != nil {
			mlog.Errorf("读取数据包失败, websocket异常中断: %v", err)
			return
		}
		// 处理包消息
		mlog.Debugf("收到数据包: %v", pack)
		switch pack.Head.DstNodeType {
		case pb.NodeType_Gate:
			// gate直接处理
			if err := actor.Send(pack.Head, pack.Body); err != nil {
				mlog.Errorf("gate服务Actor调用: %v", err)
			}
		default:
			// 转发到其他服务
			if err := framework.Send(pack.Head, pack.Body); err != nil {
				mlog.Errorf("gate服务Nats转发: %v", err)
			}
		}
	}
}

// 心跳
func (p *Player) HeartRequest(head *pb.Head, req *pb.HeartRequest, rsp *pb.HeartResponse) error {
	now := time.Now().Unix()
	if p.updateTime+p.expireTime <= now {
		// 通知客户端
		rsp.Head = &pb.RspHead{
			Code: int32(pb.ErrorCode_TIMEOUT),
			Msg:  "连接超时断开",
		}
		p.SendToClient(head, rsp)

		// 踢掉用户
		return actor.SendMsg(framework.CopyHead(head, "PlayerMgr", "Kick"))
	}
	rsp.Utc = req.Utc
	rsp.BeginTime = req.BeginTime
	rsp.EndTime = now
	head.Cmd++
	return p.SendToClient(head, rsp)
}

func (p *Player) SendToClient(head *pb.Head, msg interface{}) error {
	switch vv := msg.(type) {
	case []byte:
		return p.inet.Write(&pb.Packet{Head: head, Body: vv})
	case proto.Message:
		buf, err := proto.Marshal(vv)
		if err != nil {
			return uerror.New(1, -1, "序列化失败: %v", err)
		}
		return p.inet.Write(&pb.Packet{Head: head, Body: buf})
	}
	return uerror.New(1, -1, "消息类型错误: %v", msg)
}
