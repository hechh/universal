package player

import (
	"encoding/binary"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/async"
	"universal/framework/handler"
	"universal/framework/plog"
	"universal/framework/socket"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
)

type ApiInfo struct {
	status    int32        // 状态
	StartTime uint64       // 开始时间戳
	EndTime   uint64       // 结束时间戳
	Callback  atomic.Value // 回调
}

type Player struct {
	*async.Async
	uid     uint64              // 玩家uid
	conn    *socket.Socket      // websocket连接
	mapApi  map[string]*ApiInfo // API调用统计
	closeCh chan<- uint64       // 关闭通知
	exitCh  chan struct{}       // 退出通知
}

func NewPlayer(ws *websocket.Conn, uid uint64, closeCh chan<- uint64) (*Player, error) {
	// 启动协程
	task := async.NewAsync()
	task.Start()
	// 初始化api
	calls := make(map[string]*ApiInfo)
	handler.Walk(func(info *handler.ApiInfo) bool {
		if name := info.GetRspName(); strings.HasSuffix(name, "Response") {
			calls[name] = new(ApiInfo)
		}
		return true
	})
	// 返回玩家结构
	ret := &Player{
		Async:   task,
		uid:     uid,
		conn:    socket.NewSocket(&Frame{}, ws),
		mapApi:  calls,
		closeCh: closeCh,
		exitCh:  make(chan struct{}, 1),
	}
	go ret.loopRead()
	return ret, nil
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
			plog.Trace("uid: %d, HeartbeatResponse: %v", d.uid, vv)
			if api, ok := d.mapApi[pac.GetRspName()]; ok {
				d.handle(api, rsp)
			}
		case *pb.LoginResponse:
			if vv.PacketHead.Code == 0 {
				// 开启心跳
				go d.keepAlive()
			}
			if api, ok := d.mapApi[pac.GetRspName()]; ok {
				d.handle(api, rsp)
			}
		default:
			if api, ok := d.mapApi[pac.GetRspName()]; ok {
				d.handle(api, rsp)
			}
		}
	}
}

func (d *Player) handle(api *ApiInfo, rsp proto.Message) {
	atomic.StoreInt32(&api.status, 0)
	atomic.StoreUint64(&api.EndTime, base.GetNowMil())

	if cb, ok := api.Callback.Load().(domain.ApiCallback); ok && cb != nil {
		// 发送任务队列
		d.Push(func() {
			cb(&domain.ApiResult{
				UID:       d.uid,
				StartTime: atomic.LoadUint64(&api.StartTime),
				EndTime:   atomic.LoadUint64(&api.EndTime),
				Rsp:       rsp,
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
				PacketHead: protocol.BuildPacketHead(d.uid, pb.SERVICE_GATE),
				Time:       base.GetNow(),
			}
			d.Send(&pb.RpcHead{Id: d.uid, FuncName: "HeartbeatRequest"}, req, domain.DefaultCB)
		case <-d.exitCh:
			return
		}
	}
}

func (d *Player) Login(cb domain.ApiCallback) {
	newCB := func(rr *domain.ApiResult) {
		vv, ok := rr.Rsp.(*pb.LoginResponse)
		if rr.Error != nil || ok && vv != nil && vv.PacketHead.Code != uint32(cfgEnum.ErrorCode_Success) {
			d.closeCh <- d.uid
		}
		cb(rr)
	}

	// 获取redis
	redisCli := redis.GetRedisByAccountID(d.uid)
	if redisCli == nil {
		newCB(&domain.ApiResult{UID: d.uid, Error: fmt.Errorf("redis数据库无法连接")})
		return
	}

	// 生成token
	playerName := fmt.Sprintf("chen%d", d.uid)
	tokenKey := redisCli.GetString(fmt.Sprintf("%s_%d", base.ERK_LoginToken, d.uid))
	strToken := base.MD5(fmt.Sprintf("%s%d%s", playerName, d.uid, tokenKey))
	req := &pb.LoginRequest{
		PacketHead:  protocol.BuildPacketHead(d.uid, pb.SERVICE_GATE),
		AccountName: playerName,
		TokenKey:    strToken,
	}

	// 发送登录请求
	d.Send(&pb.RpcHead{Id: d.uid, FuncName: "LoginRequest"}, req, newCB)
}

// notify之类的请求，一律不予支持
func (d *Player) Send(head *pb.RpcHead, data proto.Message, cb domain.ApiCallback) {
	dd := data.(domain.IHead)
	head.Id = d.uid
	packetHead := dd.GetPacketHead()
	packetHead.Ckx = 0x72
	packetHead.Stx = 0x27
	packetHead.DestServerType = pb.SERVICE_GATE
	packetHead.Seqid = head.SeqId
	packetHead.Id = head.Id

	// 设置api信息
	rspName := strings.Replace(head.FuncName, "Request", "Response", 1)
	api := d.mapApi[rspName]
	if atomic.CompareAndSwapInt32(&api.status, 1, 1) {
		// 上一个请求尚未结束
		return
	}

	// 发送请求
	now := base.GetNowMil()
	if _, err := d.conn.Send(protocol.Encode(data)); err != nil {
		cb(&domain.ApiResult{UID: d.uid, Error: err})
	}

	// 设置api时间
	atomic.StoreUint64(&api.StartTime, now)
	atomic.StoreInt32(&api.status, 1)
	api.Callback.Store(cb)
}
