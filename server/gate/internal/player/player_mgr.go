package player

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"universal/common/pb"
	"universal/framework/actor"
	"universal/framework/handler"
	"universal/library/mlog"
	"universal/library/safe"
	"universal/server/gate/internal/frame"

	"github.com/gorilla/websocket"
)

type PlayerMgr struct {
	actor.Actor
	mgr    *actor.ActorMgr // 玩家管理器
	status int32           // 运行状态
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10 * 1024,
	WriteBufferSize: 10 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (d *PlayerMgr) Close() {
	atomic.StoreInt32(&d.status, 0)
	d.mgr.Stop()
	d.Actor.Stop()
	mlog.Infof("PlayerMgr服务关闭")
}

func (d *PlayerMgr) Init(ip string, port int) error {
	// 初始化ActorMgr
	d.mgr = new(actor.ActorMgr)
	player := &Player{}
	d.mgr.Register(player)
	actor.Register(d.mgr)

	// 初始化Actor
	d.Actor.Register(d)
	d.Actor.Start()
	actor.Register(d)

	// 启动ws服务
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&d.status) <= 0 {
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil || conn == nil {
			mlog.Errorf("WebSocket连接失败: %v", err)
			return
		}
		d.accept(conn)
	})

	safe.Go(func() {
		mlog.Infof("启动WebSocket服务, 地址: %s:%d", ip, port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			mlog.Errorf("WebSocket服务启动失败:%v", err)
		}
		atomic.StoreInt32(&d.status, 0)
	})
	atomic.AddInt32(&d.status, 1)
	mlog.Infof("WebSocket服务(%s:%d)启动成功，等待连接中!!!", ip, port)
	return nil
}

func (d *PlayerMgr) accept(conn *websocket.Conn) {
	usr := NewPlayer(conn, &frame.Frame{})
	if err := usr.Login(); err != nil {
		mlog.Errorf("玩家登录失败: %v", err)
		conn.Close()
		return
	}

	// 检查玩家是否已存在
	if act := d.mgr.GetActor(usr.GetId()); act != nil {
		d.Add(1)
		safe.Go(func() {
			act.Stop()
			d.Done()
		})
	}

	// 登录成功，添加玩家
	usr.Start()
	d.mgr.AddActor(usr)
	mlog.Infof("客户端: %d(%s)连接成功!!!", usr.GetId(), conn.RemoteAddr().String())

	// 循环接受消息
	usr.Dispatcher()
}

func (d *PlayerMgr) Kick(head *pb.Head) error {
	if act := d.mgr.GetActor(head.Uid); act != nil {
		d.mgr.DelActor(head.Uid)

		d.Add(1)
		safe.Go(func() {
			act.Stop()
			d.Done()
		})
	}
	return nil
}

func init() {
	handler.RegisterTrigger[PlayerMgr](pb.NodeType_Gate, "PlayerMgr.Kick", (*PlayerMgr).Kick)
}
