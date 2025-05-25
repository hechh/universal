package internal

import (
	"fmt"
	"net/http"
	"reflect"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/actor"
	"universal/library/async"
	"universal/library/mlog"
	"universal/server/gate/internal/player"

	"github.com/gorilla/websocket"
)

type PlayerMgr struct {
	actor.Actor
	mgr *actor.ActorMgr
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10 * 1024,
	WriteBufferSize: 10 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 初始化ws
func (d *PlayerMgr) Init(cfg *yaml.ServerConfig) error {
	// 初始化ActorMgr
	d.mgr = new(actor.ActorMgr)
	player := &player.Player{}
	d.mgr.Register(player)
	d.mgr.ParseFunc(reflect.TypeOf(player))
	actor.Register(d.mgr)

	// 初始化Actor
	d.Actor.Register(d)
	d.Actor.ParseFunc(reflect.TypeOf(d))
	d.Start()
	actor.Register(d)

	// 启动ws服务
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil || conn == nil {
			mlog.Errorf("WebSocket连接失败: %v", err)
			return
		}
		d.accept(conn)
	})
	async.SafeGo(mlog.Fatalf, func() {
		mlog.Infof("启动WebSocket服务, 地址: %s:%d", cfg.Ip, cfg.Port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil); err != nil {
			mlog.Errorf("WebSocket服务启动失败:%v", err)
		}
	})
	mlog.Infof("WebSocket服务启动成功，等待连接中!!!")
	return nil
}

func (p *PlayerMgr) accept(conn *websocket.Conn) {
	usr := player.NewPlayer(conn)

	// 登录
	if err := usr.Login(); err != nil {
		mlog.Errorf("玩家登录失败: %v", err)
		conn.Close()
		return
	}

	// 登录成功，添加玩家
	usr.Actor.Register(usr)
	usr.Start()
	p.mgr.AddActor(usr)
	mlog.Infof("客户端: %s连接成功!!!", conn.RemoteAddr().String())

	// 循环接受消息
	usr.Dispatcher()
}

// 剔除玩家
func (p *PlayerMgr) Kick(head *pb.Head) {
	if act := p.mgr.GetActor(head.Id); act != nil {
		// 删除玩家
		p.mgr.DelActor(head.Id)

		// 等待消息处理完成，然后关闭连接
		async.SafeGo(mlog.Fatalf, func() {
			act.Stop()
		})
	}
}

// 发送消息到客户端
func (p *PlayerMgr) SendToClient(head *pb.Head, msg []byte) {
	if act := p.mgr.GetActor(head.Id); act == nil {
		mlog.Errorf("调用接口不存在, head:%v", head)
	} else {
		head.FuncName = "SendToClient"
		if err := act.SendMsg(head, msg); err != nil {
			mlog.Errorf("发送消息到客户端失败: head:%v, error:%v", head, err)
		}
	}
}

// 广播消息到客户端
func (p *PlayerMgr) BroadcastToClient(head *pb.Head, msg []byte) {
	head.FuncName = "SendToClient"
	if err := p.mgr.SendMsg(head, msg); err != nil {
		mlog.Errorf("广播消息到客户端失败: head:%v, error:%v", head, err)
	}
}
