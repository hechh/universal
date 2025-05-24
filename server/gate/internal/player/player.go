package player

import (
	"time"
	"universal/common/pb"
	"universal/framework/actor"
	"universal/framework/domain"
	"universal/framework/network"

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
