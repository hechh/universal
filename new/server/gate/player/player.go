package player

import (
	"net"
	"poker_server/common/dao/repository/redis/login_session"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/domain"
	"poker_server/framework/library/mlog"
	"poker_server/framework/library/uerror"
	"poker_server/framework/network"
	"time"

	"github.com/golang/protobuf/proto"
)

type Player struct {
	framework.Actor
	inet       domain.INet // 网络连接
	sessionId  string      // 玩家session ID
	createTime int64       // 创建时间
	expireTime int64       // 过期时间
	updateTime int64       // 最近更新时间
	// 加密方式 todo
}

func NewPlayer(conn net.Conn, fr domain.IFrame) *Player {
	p := &Player{}
	p.inet = network.NewSocket(conn, 1024*1024)
	p.inet.Register(fr)
	return p
}

func (p *Player) login() error {
	// 第一个包一定是登录认证包
	p.inet.SetReadExpire(5)
	pack := &pb.Packet{}
	if err := p.inet.Read(pack); err != nil {
		return err
	}
	p.inet.SetReadExpire(0)
	req := &pb.GateLoginRequest{}
	if err := proto.Unmarshal(pack.Body, req); err != nil {
		return err
	}
	// 加载玩家session信息
	data, err := login_session.Get(req.SessionId)
	if err != nil {
		return err
	}
	if data == nil {
		return uerror.New(1, -1, "session不存在: %s", req.SessionId)
	}
	// 设置玩家ID
	p.Actor.SetId(data.AccountId)
	p.sessionId = req.SessionId
	p.createTime = data.CreateTime
	p.expireTime = data.ExpireTime
	p.updateTime = time.Now().Unix()
	pack.Head.Cmd++
	// 发送登录成功包
	return p.inet.Write(&pb.Packet{Head: pack.Head})
}

// 消息分发处理( 接受 websocket 传过来的消息)
func (p *Player) dispatcher() {
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
			if err := framework.Send(pack.Head, pack.Body); err != nil {
				mlog.Errorf("gate服务Actor调用: %v", err)
			}
		default:
			// 转发到其他服务
			if err := framework.SendToNode(pack.Head, pack.Body); err != nil {
				mlog.Errorf("gate服务Nats转发: %v", err)
			}
		}
	}
}

// 玩家心跳包
func (p *Player) HeartRequest(head *pb.Head, req *pb.GateHeartRequest, rsp *pb.GateHeartResponse) error {
	now := time.Now().Unix()
	if p.updateTime+p.expireTime <= now {
		// 通知客户端
		rsp.Head = &pb.RspHead{
			Code: int32(pb.ErrorCode_ERROR_TIMEOUT),
			Msg:  "连接超时断开",
		}
		head.FuncName = "SendMsgToClient"
		p.SendMsg(head, rsp)

		// 踢掉用户
		if err := framework.SendMsg(framework.CopyHead(head, "PlayerMgr", "Kick")); err != nil {
			mlog.Errorf("踢掉用户失败: %v", err)
			return uerror.New(1, -1, "心跳超时")
		}
		return uerror.New(1, -1, "连接超时断开")
	}
	rsp.Utc = req.Utc
	rsp.BeginTime = req.BeginTime
	rsp.EndTime = now
	head.Cmd++
	return p.SendMsgToClient(head, rsp)
}

// 向客户端发送数据
func (p *Player) SendMsgToClient(head *pb.Head, msg proto.Message) error {
	buf, _ := proto.Marshal(msg)
	return p.inet.Write(&pb.Packet{Head: head, Body: buf})
}

func (p *Player) SendToClient(head *pb.Head, buf []byte) error {
	return p.inet.Write(&pb.Packet{Head: head, Body: buf})
}
