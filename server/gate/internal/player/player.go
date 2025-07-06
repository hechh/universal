package player

import (
	"universal/framework/actor"
	"universal/framework/domain"
	"universal/framework/network"

	"github.com/gorilla/websocket"
)

type Player struct {
	actor.Actor
	inet      domain.INet
	status    int32
	loginTime int64
}

func NewPlayer(conn *websocket.Conn, fr domain.IFrame) *Player {
	p := &Player{}
	p.Actor.Register(p)
	p.inet = network.NewSocket(conn, 1024*1024)
	p.inet.Register(fr)
	return p
}
