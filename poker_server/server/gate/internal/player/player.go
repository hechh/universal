package player

import (
	"poker_server/common/pb"
	"poker_server/common/token"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/framework/cluster"
	"poker_server/framework/domain"
	"poker_server/framework/network"
	"poker_server/library/mlog"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

type Player struct {
	actor.Actor
	inet       domain.INet // 网络连接
	status     int32       // 玩家登录状态
	createTime int64       // 创建时间
	extra      uint32      // 设备唯一 id
	version    uint32
	// 加密方式 todo
}

func NewPlayer(conn *websocket.Conn, fr domain.IFrame) *Player {
	p := &Player{}
	p.Actor.Register(p)
	p.Actor.Start()
	p.inet = network.NewSocket(conn, fr)
	return p
}

func (p *Player) Start() {
	p.Actor.Start()
}

func (p *Player) Stop() {
	uid := p.GetId()
	p.Actor.Stop()
	p.inet.Close()
	mlog.Infof("关闭玩家actor(%d)", uid)
}

func (p *Player) GetExtra() uint32 {
	return atomic.LoadUint32(&p.extra)
}

func (p *Player) CheckToken() error {
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

	// 解析token
	tt, err := token.ParseToken(req.Token)
	if err != nil {
		return err
	}

	// 设置玩家ID
	pack.Head.Uid = tt.Uid
	now := time.Now().Unix()
	p.Actor.SetId(tt.Uid)
	p.createTime = now
	p.version = pack.Head.Version
	p.extra = pack.Head.Extra

	return cluster.Send(&pb.Head{
		Src: framework.NewSrcRouter(pack.Head.Uid, "Player"),
		Dst: framework.NewDbRouter(0, "PlayerDataMgr", "Login"),
		Uid: pack.Head.Uid,
		Seq: pack.Head.Seq,
		Cmd: pack.Head.Cmd,
	})
}

func (p *Player) Kick(extra uint32) {
	// 不是同一设备，发送剔除消息
	if extra != p.extra {
		uid := p.GetId()
		p.SendToClient(&pb.Head{
			Src: framework.NewSrcRouter(uid, "Player", "Kick"),
			Cmd: uint32(pb.CMD_KICK_PLAYER_NOTIFY),
			Uid: uid,
		}, &pb.KickPlayerNotify{})
	}
}

// 登录成功请求
func (p *Player) LoginSuccess(head *pb.Head, req *pb.GateLoginRequest, rsp *pb.GateLoginResponse) error {
	p.status = 1
	return p.SendToClient(head, rsp)
}

// 向客户端发送数据
func (p *Player) SendToClient(head *pb.Head, msg interface{}) error {
	var buf []byte
	switch vv := msg.(type) {
	case []byte:
		buf = vv
	case proto.Message:
		buf, _ = proto.Marshal(vv)
	}
	if head.Cmd%2 == 0 {
		if _, ok := pb.CMD_name[int32(head.Cmd)+1]; ok {
			head.Cmd++
			head.Seq++
		}
	}
	atomic.AddUint32(&head.Reference, 1)
	return p.inet.Write(&pb.Packet{Head: head, Body: buf})
}

// 消息分发处理( 接受 websocket 传过来的消息)
func (p *Player) Dispatcher() {
	for {
		// 从客户端持续接受包消息
		pack := &pb.Packet{}
		if err := p.inet.Read(pack); err != nil {
			mlog.Errorf("读取数据包失败, websocket异常中断: %v", err)
			actor.SendMsg(&pb.Head{ActorName: "GatePlayerMgr", FuncName: "Kick", Uid: p.GetId()})
			return
		}

		// 为登录成功，任何请求直接丢弃
		if p.status <= 0 {
			continue
		}

		// 处理包消息
		switch pack.Head.Dst.NodeType {
		case pb.NodeType_NodeTypeGate:
			pack.Head.ActorName = pack.Head.Dst.ActorName
			pack.Head.FuncName = pack.Head.Dst.FuncName
			pack.Head.ActorId = pack.Head.Dst.ActorId
			if err := actor.Send(pack.Head, pack.Body); err != nil {
				mlog.Error(pack.Head, "gate服务Actor调用 error:%v", err)
			}
		case pb.NodeType_NodeTypeGame:
			//pack.Head.Dst.ActorId = pack.Head.Uid
			if err := cluster.Send(pack.Head, pack.Body); err != nil {
				mlog.Error(pack.Head, "转发websocket数据包失败 error:%v", err)
			}
		default:
			if err := cluster.Send(pack.Head, pack.Body); err != nil {
				mlog.Error(pack.Head, "转发websocket数据包失败 error:%v", err)
			}
		}
	}
}
