package player

import (
	"fmt"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/actor"
	"universal/framework/domain"
	"universal/framework/network"
	"universal/framework/token"
	"universal/library/async"
	"universal/library/mlog"

	"github.com/gorilla/websocket"
)

type Player struct {
	actor.Actor
	conn domain.INet
	cfg  *yaml.ServerConfig
	node *pb.Node
	uid  uint64
}

func NewPlayer(node *pb.Node, cfg *yaml.ServerConfig, uid uint64) *Player {
	return &Player{
		node: node,
		cfg:  cfg,
		uid:  uid,
	}
}

func (p *Player) Login() error {
	// 设置 session
	tok, err := token.GenToken(&token.Token{Uid: p.uid})
	if err != nil {
		return err
	}
	// 建立连接
	cfg := p.cfg.Nodes[p.node.Id]
	ws, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws", d.node.Addr), nil)
	if err != nil {
		return err
	}
	d.conn = network.NewSocket(ws, 1024*1024)
	d.conn.Register(&frame.Frame{})
	// 发送登录请求
	if err := d.login(tok); err != nil {
		return err
	}
	async.SafeGo(mlog.Fatalf, d.loop)
	// 启动心跳
	d.keepAlive()
	return nil
}
