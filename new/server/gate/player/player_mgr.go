package player

import (
	"net/http"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/library/async"
	"poker_server/framework/library/mlog"
	"poker_server/server/gate/frame"
	"reflect"

	"golang.org/x/net/websocket"
)

type PlayerMgr struct {
	framework.Actor
	mgr *framework.ActorMgr
}

func (mgr *PlayerMgr) GetActorMgr() *framework.ActorMgr {
	return mgr.mgr
}

// 初始化ws
func (mgr *PlayerMgr) Init(addr string) error {
	// 初始化ActorMgr
	mgr.mgr = new(framework.ActorMgr)
	player := &Player{}
	mgr.mgr.Register(player)
	mgr.mgr.ParseFunc(reflect.TypeOf(player))
	framework.RegisterActor(mgr.mgr)

	// 初始化Actor
	mgr.Actor.Register(mgr)
	mgr.Actor.ParseFunc(reflect.TypeOf(mgr))
	framework.RegisterActor(mgr)

	// 启动ws服务
	http.Handle("/ws", websocket.Handler(mgr.accept))
	async.SafeGo(mlog.Fatalf, func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			mlog.Errorf("WebSocket服务启动失败:%v", err)
		}
	})
	mlog.Infof("WebSocket服务启动成功，等待连接中!!!")
	return nil
}

func (mgr *PlayerMgr) accept(conn *websocket.Conn) {
	// 新建玩家
	usr := NewPlayer(conn, &frame.Frame{})

	// 登录
	if err := usr.login(); err != nil {
		mlog.Errorf("玩家登录失败: %v", err)
		conn.Close()
		return
	}

	// 登录成功，添加玩家
	usr.Actor.Register(usr)
	usr.Start()
	mgr.mgr.AddActor(usr.GetId(), usr)
	mlog.Infof("客户端: %s连接成功!!!", conn.RemoteAddr().String())

	// 循环接受消息
	usr.dispatcher()
}

// 剔除玩家
func (mgr *PlayerMgr) Kick(head *pb.Head) {
	if act := mgr.mgr.GetActor(head.Id); act != nil {
		// 删除玩家
		mgr.mgr.DelActor(head.Id)

		// 等待消息处理完成，然后关闭连接
		async.SafeGo(mlog.Fatalf, func() {
			act.Stop()
		})
	}
}

// 发送消息到客户端
func (mgr *PlayerMgr) SendToClient(head *pb.Head, msg []byte) {
	if act := mgr.mgr.GetActor(head.Id); act != nil {
		mlog.Errorf("调用接口不存在, head:%v", head)
	} else {
		head.FuncName = "SendToClient"
		if err := act.Send(head, msg); err != nil {
			mlog.Errorf("发送消息到客户端失败: head:%v, error:%v", head, err)
		}
	}
}

// 广播消息到客户端
func (mgr *PlayerMgr) BroadcastToClient(head *pb.Head, msg []byte) {
	head.FuncName = "SendToClient"
	if err := mgr.mgr.Send(head, msg); err != nil {
		mlog.Errorf("广播消息到客户端失败: head:%v, error:%v", head, err)
	}
}
