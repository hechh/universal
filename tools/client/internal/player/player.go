package player

import (
	"encoding/binary"
	"fmt"
	"sync/atomic"
	"time"
	"universal/common/dao/repository/token"
	"universal/common/pb"
	"universal/framework/async"
	"universal/framework/handler"
	"universal/framework/plog"
	"universal/framework/socket"
	"universal/framework/util"
	"universal/tools/client/domain"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
)

type ApiInfo struct {
	status    int32        // 状态
	beginTime int64        // 开始时间戳
	endTime   int64        // 结束时间戳
	callback  atomic.Value // 回调
}

type Player struct {
	*async.Async
	uid     uint64              // 玩家uid
	conn    *socket.Socket      // websocket连接
	apis    map[uint32]*ApiInfo // API调用统计
	closeCh chan<- uint64       // 关闭通知
	exitCh  chan struct{}       // 退出通知
}

func NewPlayer(ws *websocket.Conn, uid uint64, closeCh chan<- uint64) (*Player, error) {
	apis := make(map[uint32]*ApiInfo)
	handler.Walk(func(info *handler.ApiInfo) bool {
		apis[info.GetReqCrc()] = new(ApiInfo)
		apis[info.GetRspCrc()] = new(ApiInfo)
		return true
	})
	return &Player{
		Async:   async.NewAsync(),
		uid:     uid,
		conn:    socket.NewSocket(&Frame{}, ws),
		apis:    apis,
		closeCh: closeCh,
		exitCh:  make(chan struct{}, 1),
	}, nil
}

func (d *Player) Login(cb domain.ResultCB) {
	d.Start()       // 开启任务队列协程
	go d.loopRead() // 开启ws监听协程

	newCB := func(rr *domain.Result) {
		cb(rr)
		vv := rr.Response.(*pb.LoginResponse)
		if rr.Error != nil || vv.PacketHead.Code != 0 {
			d.closeCh <- d.uid
		}
	}
	// 查询登录token
	key, err := token.GetLoginToken(d.uid)
	if err != nil {
		newCB(&domain.Result{UID: d.uid, Error: fmt.Errorf("查询登录token错误: %v", err)})
		return
	}
	strToken := util.MD5(fmt.Sprintf("chen%d%d%s", d.uid, d.uid, key))
	// 发送登录请求
	req := &pb.LoginRequest{
		PacketHead:  handler.BuildPacketHead(d.uid, pb.SERVICE_GATE),
		AccountName: fmt.Sprintf("chen%d", d.uid),
		TokenKey:    strToken,
	}
	d.Send(&pb.RpcHead{Id: d.uid, FuncName: "LoginRequest"}, req, newCB)
}

func (d *Player) GetUID() uint64 {
	return d.uid
}

func (d *Player) Close() {
	d.exitCh <- struct{}{}
	d.conn.Close()
	d.Stop()
}

func (d *Player) loopRead() {
	defer func() { d.closeCh <- d.uid }()
	for {
		// 接受请求
		buf, err := d.conn.Read()
		if err != nil {
			plog.Error("数据接受失败, uid: %d, error: %v", d.uid, err)
			return
		}
		// 解包
		id, buf := uint32(binary.BigEndian.Uint32(buf[:4])), buf[4:]
		pac := handler.Get(id)
		if pac == nil {
			plog.Error("协议不支持, uid: %d, protoID: %D", d.uid, id)
			return
		}
		rsp := pac.NewResponse()
		if err := proto.Unmarshal(buf, rsp); err != nil {
			plog.Error("协议解析错误, uid: %d, error: %v", d.uid, err)
			return
		}
		// 应答处理
		switch vv := rsp.(type) {
		case *pb.AllPlayerInfoNotify, *pb.ProtocolNameNotify:
		case *pb.HeartbeatResponse:
			if api, ok := d.apis[id]; ok {
				d.handle(api, rsp)
			}
		case *pb.LoginResponse:
			if vv.PacketHead.Code == 0 {
				go d.keepAlive()
			}
			if api, ok := d.apis[id]; ok {
				d.handle(api, rsp)
			}
		default:
			if api, ok := d.apis[id]; ok {
				d.handle(api, rsp)
			}
		}
	}
}

func (d *Player) handle(api *ApiInfo, rsp proto.Message) {
	atomic.StoreInt32(&api.status, 0)
	atomic.StoreInt64(&api.endTime, util.GetNowUnixMilli())
	if cb, ok := api.callback.Load().(domain.ResultCB); ok && cb != nil {
		// 发送任务队列
		d.Push(func() {
			cb(&domain.Result{
				UID:      d.uid,
				Cost:     atomic.LoadInt64(&api.beginTime) - atomic.LoadInt64(&api.endTime),
				Response: rsp,
			})
		})
	}
}

func (d *Player) keepAlive() {
	tt := time.NewTicker(3 * time.Second)
	defer tt.Stop()
	for {
		select {
		case <-tt.C:
			req := &pb.HeartbeatRequest{
				PacketHead: handler.BuildPacketHead(d.uid, pb.SERVICE_GATE),
				Time:       uint64(util.GetNowUnixSecond()),
			}
			d.Send(&pb.RpcHead{Id: d.uid, FuncName: "HeartbeatRequest"}, req, domain.DefaultResult)
		case <-d.exitCh:
			return
		}
	}
}

// notify之类的请求，一律不予支持
func (d *Player) Send(head *pb.RpcHead, data proto.Message, cb domain.ResultCB) {
	dd := data.(handler.IHead)
	head.Id = d.uid
	packetHead := dd.GetPacketHead()
	packetHead.Ckx = 0x72
	packetHead.Stx = 0x27
	packetHead.DestServerType = pb.SERVICE_GATE
	packetHead.Seqid = head.SeqId
	packetHead.Id = head.Id
	// 设置api信息
	api := d.apis[handler.GetCrc(head.FuncName)]
	if atomic.CompareAndSwapInt32(&api.status, 1, 1) {
		// 上一个请求尚未结束
		return
	}
	// 发送请求
	if _, err := d.conn.Send(handler.Encode(data)); err != nil {
		cb(&domain.Result{UID: d.uid, Error: err})
		return
	}
	// 设置api时间
	atomic.StoreInt64(&api.beginTime, util.GetNowUnixMilli())
	atomic.StoreInt32(&api.status, 1)
	api.callback.Store(cb)
}
