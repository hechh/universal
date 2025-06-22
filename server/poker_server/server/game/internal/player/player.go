package player

import (
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"poker_server/library/util"
	"poker_server/server/game/internal/player/domain"
	"poker_server/server/game/internal/player/factory"
	"poker_server/server/game/internal/player/playerfun"
	"time"
)

const (
	TTL = 15 * 60
)

type Player struct {
	actor.Actor
	*playerfun.PlayerFun
	loginTime  int64 // 登录时间
	updateTime int64 // 更新时间
}

func NewPlayer(uid uint64, data *pb.PlayerData) *Player {
	ret := &Player{PlayerFun: playerfun.NewPlayerFun()}
	ret.Actor.Register(ret)
	ret.Actor.SetId(uid)
	return ret
}

func (p *Player) OnTick() {
	if err := p.Save(); err != nil {
		mlog.Errorf("保存玩家数据失败: %v", err)
	}

	// 剔除玩家
	uid := p.GetId()
	if p.updateTime+TTL < time.Now().Unix() {
		dst := framework.NewGateRouter(uid, "GatePlayerMgr", "Kick")
		framework.Send(framework.NewHead(dst, pb.RouterType_RouterTypeUid, uid))
		actor.SendMsg(&pb.Head{ActorName: "PlayerMgr", FuncName: "Kick", ActorId: uid}, uid)
	}
}

func (p *Player) Save() error {
	if !p.PlayerFun.HasChange() {
		return nil
	}

	playerData := &pb.PlayerData{Uid: p.GetId()}
	p.PlayerFun.Walk(func(tt pb.PlayerDataType, fun domain.IPlayerFun) bool {
		fun.Save(playerData)
		return true
	})

	// 发送db保存数据
	dst := framework.NewDbRouter(uint64(pb.DataType_DataTypePlayerData), "PlayerDataMgr", "Update")
	return framework.Send(framework.NewHead(dst, pb.RouterType_RouterTypeUid, p.GetId()), &pb.UpdatePlayerDataNotify{
		DataType: pb.DataType_DataTypePlayerData,
		Data:     playerData,
	})
}

func (p *Player) Relogin(head *pb.Head, req *pb.GateLoginRequest, rsp *pb.GateLoginResponse) error {
	p.loginTime = time.Now().Unix()
	p.updateTime = p.loginTime
	return framework.Send(framework.SwapToGate(head, head.Uid, "Player", "LoginSuccess"), req)
}

// 登录请求
func (p *Player) Login(head *pb.Head, req *pb.GateLoginRequest, rsp *pb.GateLoginResponse) error {
	// 初始化所有模块
	for tt, f := range factory.FUNCS {
		p.PlayerFun.Set(tt, f(p.PlayerFun))
	}

	// 按照顺序加载模块
	for _, tt := range factory.LoadList {
		fun := p.PlayerFun.Get(tt)
		if err := fun.Load(req.PlayerData); err != nil {
			return err
		}
	}

	// 加载完成回调
	for _, tt := range factory.LoadList {
		fun := p.PlayerFun.Get(tt)
		if err := fun.LoadComplate(); err != nil {
			return err
		}
	}
	p.loginTime = time.Now().Unix()
	p.updateTime = p.loginTime
	p.RegisterTimer(&pb.Head{ActorName: "Player", FuncName: "OnTick", ActorId: head.Uid, Uid: head.Uid}, 5*time.Second, -1)
	return framework.Send(framework.SwapToGate(head, head.Uid, "Player", "LoginSuccess"), rsp)
}

func (p *Player) HeartRequest(head *pb.Head, req *pb.GateHeartRequest, rsp *pb.GateHeartResponse) error {
	now := time.Now().Unix()
	if p.updateTime+TTL <= now {
		dst := framework.NewGateRouter(head.Uid, "GatePlayerMgr", "Kick")
		framework.Send(framework.NewHead(dst, pb.RouterType_RouterTypeUid, head.Uid))

		actor.SendMsg(&pb.Head{ActorName: "PlayerMgr", FuncName: "Kick", ActorId: head.Uid}, head.Uid)
		return uerror.NEW(pb.ErrorCode_TIME_OUT, head, "心跳超时: req:%v", req)
	}

	p.updateTime = now
	rsp.Utc = req.Utc
	rsp.BeginTime = req.BeginTime
	rsp.EndTime = now
	return nil
}

func (p *Player) RewardReq(head *pb.Head, req *pb.RewardReq, rsp *pb.RewardRsp) error {
	return p.GetBagFunc().RewardReq(head, req, rsp)
}

func (p *Player) ConsumeReq(head *pb.Head, req *pb.ConsumeReq, rsp *pb.ConsumeRsp) error {
	return p.GetBagFunc().ConsumeReq(head, req, rsp)
}

func (p *Player) TexasJoinRoomReq(head *pb.Head, req *pb.TexasJoinRoomReq, rsp *pb.TexasJoinRoomRsp) error {
	return p.GetBaseFunc().TexasJoinRoom(head, req, rsp)
}

func (p *Player) TexasQuitRoomReq(head *pb.Head, req *pb.TexasQuitRoomReq, rsp *pb.TexasQuitRoomRsp) error {
	newHead := &pb.Head{
		Dst: framework.NewRoomRouter(req.RoomId, "TexasGameMgr", "QuitRoomReq"),
		Uid: head.Uid,
	}
	if err := framework.Request(newHead, req, rsp); err != nil {
		return err
	}
	p.GetBagFunc().AddProp(rsp.CoinType, rsp.Chip)
	return p.GetBaseFunc().TexasQuitRoom(head, req, rsp)
}

func (p *Player) TexasBuyInReq(head *pb.Head, req *pb.TexasBuyInReq, rsp *pb.TexasBuyInRsp) error {
	// 先扣除玩家道具
	/*
		conReq := reward.ToConsumeRequest(pb.CoinType(req.CoinType), req.Chip)
		conRsp := &pb.ConsumeRsp{}
		if err := p.GetBagFunc().ConsumeReq(head, conReq, conRsp); err != nil {
			return err
		}
	*/
	return framework.Send(framework.SwapToRoom(head, req.RoomId, "TexasGameMgr", "BuyInReq"), req)
}

func (p *Player) RummyJoinRoomReq(head *pb.Head, req *pb.RummyJoinRoomReq, rsp *pb.RummyJoinRoomRsp) error {
	return p.GetBaseFunc().RummyJoinRoom(head, req, rsp)
}

func (p *Player) RummyQuitRoomReq(head *pb.Head, req *pb.RummyQuitRoomReq, rsp *pb.RummyQuitRoomRsp) error {
	return p.GetBaseFunc().RummyQuitRoom(head, req, rsp)
}

func (p *Player) QueryPlayerData(head *pb.Head, req *pb.QueryPlayerDataReq, rsp *pb.QueryPlayerDataRsp) error {
	return p.GetBaseFunc().QueryPlayerData(head, req, rsp)
}

// RummyChangeRoomReq 换桌请求 第一步退桌
func (p *Player) RummyChangeRoomReq(head *pb.Head, req *pb.RummyChangeRoomReq, rsp *pb.RummyChangeRoomRsp) error {
	// 异步清理room服数据
	quitReq := &pb.RummyQuitRoomReq{
		RoomId:   req.RoomId,
		IsChange: true,
	}
	return p.GetBaseFunc().RummyQuitRoom(head, quitReq, nil)
}

// RummyChangeNewRoomReq 退桌成功后 第二步执行匹配新桌
func (p *Player) RummyChangeNewRoomReq(head *pb.Head, req *pb.RummyChangeRoomReq, rsp *pb.RummyChangeRoomRsp) error {
	// 同步获取匹配服数据
	types := util.DestructRoomId(req.RoomId)
	return framework.Send(framework.SwapToMatch(head, util.GenMatchId(types), "MatchRummyRoom", "Without"), req)
}

// RummyChangeToRoomReq 第三步匹配到新桌 加入新桌
func (p *Player) RummyChangeToRoomReq(head *pb.Head, req *pb.RummyChangeRoomReq, rsp *pb.RummyChangeRoomRsp) error {
	joinReq := &pb.RummyJoinRoomReq{
		RoomId:   req.RoomId,
		IsChange: true,
	}
	return p.GetBaseFunc().RummyJoinRoom(head, joinReq, nil)
}
