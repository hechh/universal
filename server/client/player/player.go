package player

import (
	"fmt"
	"time"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/actor"
	"universal/framework/domain"
	"universal/framework/network"
	"universal/framework/token"
	"universal/library/async"
	"universal/library/mlog"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

type Player struct {
	actor.Actor
	conn domain.INet
	cfg  *yaml.ServerConfig
	cmds map[uint32]func() proto.Message
	node *pb.Node
	uid  uint64
}

func NewPlayer(node *pb.Node, cfg *yaml.ServerConfig, uid uint64, cmds map[uint32]func() proto.Message) *Player {
	ret := &Player{
		node: node,
		cfg:  cfg,
		uid:  uid,
		cmds: cmds,
	}
	ret.SetId(uid)
	ret.Actor.Register(ret)
	ret.Start()
	return ret
}

func (p *Player) SendCmd(cmd uint32, buf []byte) error {
	head := &pb.Head{
		DstNodeType: p.node.Type,
		DstNodeId:   p.node.Id,
		RouteId:     p.uid,
		Id:          p.uid,
		IdType:      pb.IdType_UID,
		Cmd:         uint32(cmd),
	}
	return p.conn.Write(&pb.Packet{Head: head, Body: buf})
}

func (p *Player) Login() error {
	head := &pb.Head{ActorName: "PlayerMgr", FuncName: "Remove", IdType: pb.IdType_UID, Id: p.uid}

	// 建立连接
	ws, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%d/ws", p.cfg.Ip, p.cfg.Port), nil)
	if err != nil {
		actor.SendMsg(head, p.uid)
		return err
	}
	p.conn = network.NewSocket(ws, 1024*1024)
	p.conn.Register(&Frame{})

	// 设置 session
	tok, err := token.GenToken(&token.Token{Uid: p.uid})
	if err != nil {
		actor.SendMsg(head, p.uid)
		return err
	}

	// 发送登录请求
	buf, _ := proto.Marshal(&pb.LoginRequest{Token: tok})
	if err := p.SendCmd(uint32(pb.CMD_LOGIN_REQUEST), buf); err != nil {
		actor.SendMsg(head, p.uid)
		return err
	}

	// 接收登录返回消息
	pack := &pb.Packet{}
	if err := p.conn.Read(pack); err != nil {
		actor.SendMsg(head, p.uid)
		return err
	}
	loginRsp := &pb.LoginResponse{}
	if err := proto.Unmarshal(pack.Body, loginRsp); err != nil {
		actor.SendMsg(head, p.uid)
		return err
	}
	if loginRsp.Head != nil && loginRsp.Head.Code != 0 {
		actor.SendMsg(head, p.uid)
		return fmt.Errorf("登录失败: %v", loginRsp.Head)
	}
	async.SafeGo(mlog.Errorf, p.loop)
	async.SafeGo(mlog.Errorf, p.keepAlive)
	return nil
}

func (p *Player) loop() {
	for {
		pack := &pb.Packet{}
		if err := p.conn.Read(pack); err != nil {
			fmt.Println("读取消息失败:", err)
			break
		}

		switch pack.Head.Cmd {
		case uint32(pb.CMD_HEART_RESPONSE):
		default:
			if ff, ok := p.cmds[pack.Head.Cmd]; ok {
				msg := ff()
				if err := proto.Unmarshal(pack.Body, msg); err != nil {
					fmt.Println("反序列化失败:", err)
					break
				}
				fmt.Println("收到%s消息: %s", pb.CMD(pack.Head.Cmd).String(), msg.String())
			}
		}
	}
}

func (p *Player) keepAlive() {
	// 循环发送心跳
	tt := time.NewTicker(3 * time.Second)
	defer tt.Stop()
	buf, _ := proto.Marshal(&pb.HeartRequest{})
	for {
		<-tt.C
		if err := p.SendCmd(uint32(pb.CMD_HEART_REQUEST), buf); err != nil {
			fmt.Println("发送心跳包失败:", err)
			break
		}
	}
}
