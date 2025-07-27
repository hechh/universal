package player

import (
	"bytes"
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
	"poker_server/server/client/internal/stat"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var (
	pool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(nil)
		},
	}
)

func get() *bytes.Buffer {
	return pool.Get().(*bytes.Buffer)
}

func put(val *bytes.Buffer) {
	pool.Put(val)
}

type ClientPlayer struct {
	actor.Actor
	cmds map[uint32]*atomic.Pointer[stat.CmdStat]
	conn domain.INet
	cfg  *yaml.NodeConfig
	node *pb.Node
	uid  uint64
}

func NewClientPlayer(node *pb.Node, cfg *yaml.NodeConfig, uid uint64) *ClientPlayer {
	ret := &ClientPlayer{
		cmds: make(map[uint32]*atomic.Pointer[stat.CmdStat]),
		node: node,
		cfg:  cfg,
		uid:  uid,
	}
	for cmd := range request.Cmds {
		ret.cmds[cmd] = new(atomic.Pointer[stat.CmdStat])
	}
	ret.Actor.Register(ret)
	ret.SetId(uid)
	ret.Start()
	return ret
}

func (p *ClientPlayer) Login(st *stat.CmdStat) error {
	head := &pb.Head{ActorName: "PlayerMgr", FuncName: "Remove", Uid: p.uid}
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
	if err := p.SendCmd(st, uint32(pb.CMD_GATE_LOGIN_REQUEST), p.uid, buf); err != nil {
		actor.SendMsg(head, p.uid)
		return err
	}

	// 接收登录返回消息
	pack := &pb.Packet{}
	if err := p.conn.Read(pack); err != nil {
		actor.SendMsg(head, p.uid)
		return err
	}
	p.finish(pack.Head.Cmd)

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
	mlog.Infof("登录成功: %d", p.uid)
	return nil
}

func (p *ClientPlayer) finish(cmd uint32) {
	if rr := p.cmds[cmd-cmd%2].Load(); rr != nil {
		if ret := rr.Get(p.uid); ret != nil {
			ret.Finish()
			rr.Done()
		}
	}
}

func (p *ClientPlayer) SendCmd(st *stat.CmdStat, cmd uint32, routeId uint64, buf []byte) error {
	if rr := p.cmds[cmd].Load(); rr != nil {
		if ret := rr.Get(p.uid); ret != nil && !ret.IsFinish() {
			return uerror.New(pb.ErrorCode_REQUEST_FAIELD, "尚未接收到应答")
		}
	}
	if st != nil {
		if rr := st.Get(p.uid); rr != nil {
			rr.Start()
			st.Add(1)
			p.cmds[cmd].Store(st)
		}
	}

	if routeId <= 0 {
		routeId = p.uid
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
	}
	return p.conn.Write(&pb.Packet{Head: head, Body: buf})
}

// 循环发送心跳
func (p *ClientPlayer) keepAlive() {
	tt := time.NewTicker(3 * time.Second)
	defer tt.Stop()
	buf, _ := proto.Marshal(&pb.GateHeartRequest{})
	for {
		<-tt.C
		if err := p.SendCmd(nil, uint32(pb.CMD_GATE_HEART_REQUEST), p.uid, buf); err != nil {
			mlog.Errorf("发送心跳包失败: %v", err)
			break
		}
	}
}

func (p *ClientPlayer) loop() {
	for {
		pack := &pb.Packet{}
		if err := p.conn.Read(pack); err != nil {
			mlog.Errorf("读取消息失败: %v", err)
			break
		}
		p.finish(pack.Head.Cmd)

		/*
			switch pack.Head.Cmd {
			case uint32(pb.CMD_GATE_HEART_RESPONSE):
			default:
				ttl := p.stat.Finish(pack.Head.Cmd, pack.Head.Seq)
				if ff, ok := request.Cmds[pack.Head.Cmd]; ok && pack.Head.Cmd%2 == 1 {
					msg := ff()
					if err := proto.Unmarshal(pack.Body, msg); err != nil {
						mlog.Errorf("反序列化失败: %v", err)
						break
					}
					mlog.Infof("[%d] [%s] 接口调用时长: %dms, %v, rsp:%s", p.uid, pb.CMD(pack.Head.Cmd).String(), ttl, pack.Head, msg.String())
				} else {
					mlog.Infof("[%d] [%s] 接口调用时长:%dms, %v, body:%s", p.uid, pb.CMD(pack.Head.Cmd).String(), ttl, pack.Head, string(pack.Body))
				}
			}
		*/
	}
}
