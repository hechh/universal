package player

import (
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/actor"
	"universal/framework/cluster"
	"universal/framework/define"
	"universal/framework/network"
	"universal/library/mlog"
	"universal/server/gate/internal/token"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

type Player struct {
	actor.Actor
	inet      define.INet
	status    int32
	loginTime int64
}

func NewPlayer(conn *websocket.Conn, fr define.IFrame) *Player {
	p := &Player{}
	p.Actor.Register(p)
	p.inet = network.NewSocket(conn, 1024*1024)
	p.inet.SetFrame(fr)
	return p
}

func (p *Player) Stop() {
	p.inet.Close()
	p.Actor.Stop()
}

func (p *Player) Login() error {
	p.inet.SetReadExpire(5)
	pack := &pb.Packet{}
	if err := p.inet.Read(pack); err != nil {
		return err
	}
	p.inet.SetReadExpire(0)
	req := &pb.LoginReq{}
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
	p.loginTime = now

	/*
		head := framework.SwapToDb(pack.Head, tt.Uid, "PlayerDataMgr", "Login")
		head.Src = framework.NewSrcRouter(pb.RouterType_RouterTypeUid, tt.Uid)
		head.Dst.RouterType = pb.RouterType_RouterTypeUid
		return framework.Send(head, req)
	*/
	return nil
}

func (p *Player) LoginSuccess(head *pb.Head, rsp *pb.LoginRsp) error {
	p.status = 1
	return p.SendToClient(head, rsp)
}

func (p *Player) SendToClient(head *pb.Head, msg interface{}) error {
	var buf []byte
	switch vv := msg.(type) {
	case []byte:
		buf = vv
	case proto.Message:
		buf, _ = proto.Marshal(vv)
	}
	atomic.AddUint32(&head.Reference, 1)
	return p.inet.Write(&pb.Packet{Head: head, Body: buf})
}

func (p *Player) Dispatcher() {
	for {
		pack := &pb.Packet{}
		if err := p.inet.Read(pack); err != nil {
			mlog.Errorf("读取数据包失败, websocket异常中断: %v", err)
			return
		}

		if p.status <= 0 {
			continue
		}

		switch pack.Head.Dst.NodeType {
		case pb.NodeType_Gate:
			/*
				rpc.ParseNodeRouter(pack.Head)
				mlog.Debugf("收到websocket数据包 pack:%v", pack)
				if err := actor.Send(pack.Head, pack.Body); err != nil {
					mlog.Errorf("gate服务Actor调用: %v", err)
				}
			*/
		default:
			if err := cluster.Send(pack.Head, pack.Body); err != nil {
				mlog.Errorf("转发websocket数据包失败: pack:%v, error:%v", pack, err)
			} else {
				mlog.Debugf("转发websocket数据包 pack:%v", pack)
			}
		}
	}
}
