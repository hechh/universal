package player

import (
	"fmt"
	"poker_server/common/pb"
	"poker_server/common/token"
	"poker_server/common/yaml"
	"poker_server/framework/actor"
	"poker_server/framework/domain"
	"poker_server/framework/network"
	"poker_server/library/mlog"
	"poker_server/library/safe"
	"poker_server/library/uerror"
	"poker_server/server/client/internal/frame"
	"poker_server/server/client/internal/request"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

type TimeInfo struct {
	Seq       uint32
	StartTime int64
	EndTime   int64
}

type ApiInfo struct {
	seq   uint32
	infos map[uint32]*TimeInfo
}

type ClientPlayer struct {
	actor.Actor
	conn domain.INet
	cfg  *yaml.NodeConfig
	apis map[uint32]*ApiInfo
	node *pb.Node
	uid  uint64
}

func NewClientPlayer(node *pb.Node, cfg *yaml.NodeConfig, uid uint64) *ClientPlayer {
	ret := &ClientPlayer{
		node: node,
		cfg:  cfg,
		uid:  uid,
		apis: make(map[uint32]*ApiInfo),
	}
	for cmd := range request.Cmds {
		ret.apis[cmd] = &ApiInfo{infos: make(map[uint32]*TimeInfo)}
	}
	ret.Actor.Register(ret)
	ret.SetId(uid)
	ret.Start()
	return ret
}

func (p *ClientPlayer) SendCmd(cmd uint32, routeId uint64, buf []byte) error {
	if routeId <= 0 {
		routeId = p.uid
	}
	api, ok := p.apis[cmd]
	if !ok {
		return uerror.New(1, pb.ErrorCode_CMD_NOT_FOUND, "cmd不存在%s", pb.CMD(cmd).String())
	}
	head := &pb.Head{
		Src: &pb.NodeRouter{ActorId: routeId},
		Dst: &pb.NodeRouter{
			NodeType: p.node.Type,
			NodeId:   p.node.Id,
			ActorId:  routeId,
		},
		Uid: p.uid,
		Cmd: uint32(cmd),
		Seq: api.seq,
	}
	defer func() {
		api.infos[api.seq] = &TimeInfo{Seq: api.seq, StartTime: time.Now().UnixMilli()}
		api.seq++
	}()
	return p.conn.Write(&pb.Packet{Head: head, Body: buf})
}

func (p *ClientPlayer) update(cmd uint32, seq uint32) int64 {
	if api, ok := p.apis[cmd-(cmd%2)]; ok && cmd%2 == 1 {
		if tt, ok := api.infos[seq-1]; ok {
			tt.EndTime = time.Now().UnixMilli()
			return tt.EndTime - tt.StartTime
		}
	}
	return 0
}

func (p *ClientPlayer) Login() error {
	head := &pb.Head{
		ActorName: "PlayerMgr",
		FuncName:  "Remove",
		Uid:       p.uid,
	}

	// 建立连接
	wsUrl := fmt.Sprintf("ws://%s:%d/ws", p.cfg.Ip, p.cfg.Port)
	ws, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		actor.SendMsg(head, p.uid)
		return err
	}
	p.conn = network.NewSocket(ws, frame.New(p.node))

	// 设置 session
	tok, err := token.GenToken(&token.Token{Uid: p.uid})
	if err != nil {
		actor.SendMsg(head, p.uid)
		return err
	}

	// 发送登录请求
	buf, _ := proto.Marshal(&pb.GateLoginRequest{Token: tok})
	if err := p.SendCmd(uint32(pb.CMD_GATE_LOGIN_REQUEST), p.uid, buf); err != nil {
		actor.SendMsg(head, p.uid)
		return err
	}

	// 接收登录返回消息
	pack := &pb.Packet{}
	if err := p.conn.Read(pack); err != nil {
		actor.SendMsg(head, p.uid)
		return err
	}
	ttl := p.update(pack.Head.Cmd, pack.Head.Seq)
	mlog.Infof("uid:%d, 登录时间：%dms", p.uid, ttl)
	loginRsp := &pb.GateLoginResponse{}
	if err := proto.Unmarshal(pack.Body, loginRsp); err != nil {
		actor.SendMsg(head, p.uid)
		return err
	}
	if loginRsp.Head != nil {
		return uerror.ToError(loginRsp.Head)
	}
	safe.Go(p.loop)
	safe.Go(p.keepAlive)
	return nil
}

func (p *ClientPlayer) loop() {
	for {
		pack := &pb.Packet{}
		if err := p.conn.Read(pack); err != nil {
			mlog.Errorf("读取消息失败: %v", err)
			break
		}

		switch pack.Head.Cmd {
		case uint32(pb.CMD_GATE_HEART_RESPONSE):
		default:
			ttl := p.update(pack.Head.Cmd, pack.Head.Seq)
			if ff, ok := request.Cmds[pack.Head.Cmd]; ok && pack.Head.Cmd%2 == 1 {
				msg := ff()
				if err := proto.Unmarshal(pack.Body, msg); err != nil {
					mlog.Errorf("反序列化失败: %v", err)
					break
				}
				mlog.Infof("[%d] [%s] 接口调用时长: %dms, %v, rsp:%s", p.uid, pb.CMD(pack.Head.Cmd).String(), ttl, pack.Head, msg.String())
			} else {
				mlog.Infof("[%d] 接口调用时长:%dms, %v, body:%s", p.uid, pack.Head.Cmd, ttl, pack.Head, string(pack.Body))
			}
		}
	}
}

func (p *ClientPlayer) keepAlive() {
	// 循环发送心跳
	tt := time.NewTicker(3 * time.Second)
	defer tt.Stop()
	buf, _ := proto.Marshal(&pb.GateHeartRequest{})
	for {
		<-tt.C
		if err := p.SendCmd(uint32(pb.CMD_GATE_HEART_REQUEST), p.uid, buf); err != nil {
			mlog.Errorf("发送心跳包失败: %v", err)
			break
		}
	}
}
