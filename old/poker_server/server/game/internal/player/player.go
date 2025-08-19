package player

import (
	"poker_server/common/pb"
	"poker_server/common/room_util"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"poker_server/library/util"
	"poker_server/server/game/internal/player/domain"
	"poker_server/server/game/internal/player/factory"
	"poker_server/server/game/internal/player/playerfun"
	"time"
)

const (
	TTL      = 15 * 60
	MATCHTTL = int64(15)
)

type Player struct {
	actor.Actor
	*playerfun.PlayerFun
	loginTime   int64 // 登录时间
	updateTime  int64 // 更新时间
	status      pb.PlayerStatus
	closeCliCon bool
}

func NewPlayer(uid uint64, data *pb.PlayerData) *Player {
	ret := &Player{PlayerFun: playerfun.NewPlayerFun()}
	ret.Actor.Register(ret)
	ret.Actor.SetId(uid)
	return ret
}

func (p *Player) Stop() {
	uid := p.GetId()
	p.GetBaseFunc().Close(uid)
	p.Save()
	p.Actor.Stop()
	mlog.Infof("Player关闭成功 uid:%d", uid)
}

func (p *Player) OnTick() {
	now := time.Now().Unix()
	if now-p.updateTime > 3 {
		if err := p.Save(); err != nil {
			mlog.Errorf("保存玩家数据失败: %v", err)
		}
	}

	// 检测匹配超时
	matchInfo := p.GetBaseFunc().GetMatchInfo()
	if matchInfo != nil && now >= matchInfo.TimeOut {
		p.GetBaseFunc().MatchStop()
	}

	// 剔除玩家
	uid := p.GetId()
	if p.updateTime+TTL <= now {
		if p.closeCliCon == false {
			p.closeCliCon = true
			cluster.Send(&pb.Head{
				Src: framework.NewSrcRouter(uid, "Player"),
				Dst: framework.NewGateRouter(uid, "GatePlayerMgr", "Kick"),
			})
		}
		if p.status != pb.PlayerStatus_PlayerStatusJoinGame { //防止重连失效
			actor.SendMsg(&pb.Head{ActorName: "PlayerMgr", FuncName: "Kick", ActorId: uid}, uid)
		}
	}
}

func (p *Player) Save() error {
	if !p.PlayerFun.IsChange() {
		return nil
	}

	playerData := &pb.PlayerData{Uid: p.GetId()}
	p.PlayerFun.Walk(func(tt pb.PlayerDataType, fun domain.IPlayerFun) bool {
		fun.Save(playerData)
		return true
	})

	// 发送db保存数据
	head := &pb.Head{
		Src: framework.NewSrcRouter(p.GetId(), "Player"),
		Dst: framework.NewDbRouter(uint64(pb.DataType_DataTypePlayerData), "PlayerDataMgr", "Update"),
	}
	err := cluster.Send(head, &pb.UpdatePlayerDataNotify{
		DataType: pb.DataType_DataTypePlayerData,
		Data:     playerData,
	})
	if err == nil {
		p.HasSave()
	}
	return err
}

func (p *Player) Relogin(head *pb.Head, req *pb.GateLoginRequest, rsp *pb.GateLoginResponse) error {
	p.loginTime = time.Now().Unix()
	p.updateTime = p.loginTime
	p.closeCliCon = false
	if req.PlayerData != nil && req.PlayerData.Base != nil {
		p.GetBaseFunc().SetPlayerInfo(req.PlayerData.Base.PlayerInfo)
	}

	playerData := &pb.PlayerData{Uid: p.GetId()}
	p.PlayerFun.Walk(func(tt pb.PlayerDataType, fun domain.IPlayerFun) bool {
		fun.Save(playerData)
		return true
	})
	mlog.Info(head, "%d重新登录成功 %v", head.Uid, playerData)
	return cluster.Send(framework.SwapToGate(head, head.Uid, "Player", "LoginSuccess"), rsp)
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
	p.closeCliCon = false
	p.RegisterTimer(&pb.Head{ActorName: "Player", FuncName: "OnTick", ActorId: head.Uid, Uid: head.Uid}, 5*time.Second, -1)
	mlog.Info(head, "%d登录成功 %v", head.Uid, req.PlayerData)
	return cluster.Send(framework.SwapToGate(head, head.Uid, "Player", "LoginSuccess"), rsp)
}

func (p *Player) HeartRequest(head *pb.Head, req *pb.GateHeartRequest, rsp *pb.GateHeartResponse) error {
	now := time.Now().Unix()
	if p.updateTime+TTL <= now {
		if p.closeCliCon == false {
			p.closeCliCon = true

			cluster.Send(&pb.Head{
				Uid: head.Uid,
				Src: framework.NewSrcRouter(head.Uid, "Player"),
				Dst: framework.NewGateRouter(head.Uid, "GatePlayerMgr", "Kick"),
			})
		}
		if p.status != pb.PlayerStatus_PlayerStatusJoinGame { //防止重连失效
			actor.SendMsg(&pb.Head{ActorName: "PlayerMgr", FuncName: "Kick", ActorId: head.Uid}, head.Uid)
		}
		return uerror.New(1, pb.ErrorCode_TIME_OUT, "心跳超时: req:%v", req)
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

func (p *Player) GetBagReq(head *pb.Head, req *pb.GetBagReq, rsp *pb.GetBagRsp) error {
	return p.GetBagFunc().GetBagReq(head, req, rsp)
}

// 加入德州房间
func (p *Player) TexasJoinRoomReq(head *pb.Head, req *pb.TexasJoinRoomReq, rsp *pb.TexasJoinRoomRsp) error {
	if err := p.GetBaseFunc().TexasJoinRoom(head, req, rsp, true); err != nil {
		return err
	}
	// php扣钱
	if req.MatchType == pb.MatchType_MatchTypeNone {
		if err := p.GetBagFunc().ChargeTransIn(head, req.RoomId); err != nil {
			mlog.Errorf("转入钱失败: %v", err)
		}
	}
	return nil
}

// 退出德州房间请求
func (p *Player) TexasQuitRoomReq(head *pb.Head, req *pb.TexasQuitRoomReq, rsp *pb.TexasQuitRoomRsp) error {
	if err := p.GetBaseFunc().TexasQuitRoom(head, req, rsp); err != nil {
		return err
	}
	// 转入php道具
	matchType, _, _ := room_util.TexasRoomIdTo(req.RoomId)
	if matchType == pb.MatchType_MatchTypeNone {
		if err := p.GetBagFunc().ChargeTransOut(head, req.RoomId); err != nil {
			mlog.Error(head, "转出钱失败 error:%v", err)
		}
	}
	return nil
}

// 德州买入请求
func (p *Player) TexasBuyInReq(head *pb.Head, req *pb.TexasBuyInReq, rsp *pb.TexasBuyInRsp) error {
	matchType, _, _ := room_util.TexasRoomIdTo(req.RoomId)
	if matchType == pb.MatchType_MatchTypeNone {
		bagFun := p.GetBagFunc()
		if err := bagFun.SubProp(pb.CoinType(req.CoinType), req.Chip); err != nil {
			return err
		}
		mlog.Info(head, "TexasBuyInReq: %v successed", req)
	}
	return cluster.Send(framework.SwapToRoom(head, req.RoomId, "TexasGameMgr", "BuyInReq"), req)
}

// 房间销毁通知
func (p *Player) TexasFinishNotify(head *pb.Head, event *pb.TexasFinishNotify) error {
	if event == nil {
		return nil
	}
	if event.Incr > 0 {
		p.GetBagFunc().AddProp(pb.CoinType(event.PropId), event.Incr)
	}
	// 转入php道具
	if err := p.GetBagFunc().ChargeTransOut(head, event.RoomId); err != nil {
		mlog.Error(head, "转出钱失败 error:%v", err)
	}
	p.GetBaseFunc().ResetRoomInfo()
	return nil
}

// 德州换房间请求
func (p *Player) TexasChangeReq(head *pb.Head, req *pb.TexasChangeRoomReq, rsp *pb.TexasChangeRoomRsp) error {
	return p.GetBaseFunc().TexasChangeRoom(head, req, rsp)
}

func (p *Player) RummyJoinRoomReq(head *pb.Head, req *pb.RummyJoinRoomReq, rsp *pb.RummyJoinRoomRsp) error {
	//大厅货币转入背包
	err := p.GetBagFunc().ChargeTransIn(head, req.RoomId)
	if err != nil {
		return err
	}

	// 判断是否为断线重连
	baseFun := p.GetBaseFunc()
	if roomId, err := baseFun.GetRummyRealRoomId(head, req.RoomId); err != nil {
		return err
	} else {
		req.RoomId = roomId
	}

	// 加入房间
	err = p.rummyJoinRoomReq(head, req, rsp)
	if err != nil {
		return err
	}
	return nil
}

func (p *Player) rummyJoinRoomReq(head *pb.Head, req *pb.RummyJoinRoomReq, rsp *pb.RummyJoinRoomRsp) error {
	baseFun := p.GetBaseFunc()
	types := util.DestructRoomId(req.RoomId)

	if baseFun.GetRoomInfo() == nil { //新用户检测
		newHead := &pb.Head{
			Src: framework.NewSrcRouter(head.Uid, "Player"),
			Dst: framework.NewMatchRouter(uint64(types.GetGameType())<<32|uint64(types.GetCoinType()), "MatchRummyRoom", "Query"),
		}
		data := &pb.RummyRoomData{}
		if err := cluster.Request(newHead, req.RoomId, data); err != nil {
			return err
		}

		if p.GetBagFunc().GetProp(types.GetCoinType()) < data.RoomCfg.MinBuyIn {

			err := p.GetBagFunc().ChargeTransOut(head, req.RoomId)
			if err != nil {
				return err
			}

			return uerror.New(1, pb.ErrorCode_GAME_PROP_NOT_ENOUGH, "玩家道具不足加入游戏")
		}
	}

	if p.GetBagFunc().GetProp(types.GetCoinType()) > 0 {
		req.Coin = p.GetBagFunc().GetProp(types.GetCoinType())
		p.GetBagFunc().SubProp(types.GetCoinType(), req.Coin)
	}

	// 加入房间
	req.PlayerInfo = baseFun.GetPlayerInfo()
	joinHead := &pb.Head{
		Src: head.Src,
		Dst: framework.NewRoomRouter(req.RoomId, "RummyGameMgr", "JoinRoomReq"),
		Uid: head.Uid,
	}

	if err := cluster.Request(joinHead, req, rsp); err != nil {
		return err
	} else if rsp.Head != nil {
		p.GetBagFunc().AddProp(types.GetCoinType(), req.Coin)
		err = p.GetBagFunc().ChargeTransOut(head, req.RoomId)
		if err != nil {
			return err
		}
		return uerror.ToError(rsp.Head)
	}

	// 缓存信息
	p.GetBaseFunc().RummyJoinRoom(req.RoomId)
	return nil
}

func (p *Player) RummyQuitRoomReq(head *pb.Head, req *pb.RummyQuitRoomReq, rsp *pb.RummyQuitRoomRsp) error {
	// 判断是否为断线重连
	if roomId, err := p.GetBaseFunc().GetRummyRealRoomId(head, req.RoomId); err != nil {
		return err
	} else {
		req.RoomId = roomId
	}

	newHead := &pb.Head{
		Dst: framework.NewRoomRouter(req.RoomId, "RummyGame", "QuitRoomReq"),
		Uid: head.Uid,
	}
	if err := cluster.Request(newHead, req, rsp); err != nil {
		return err
	}
	bagFunc := p.GetBagFunc()
	types := util.DestructRoomId(req.RoomId)
	bagFunc.AddProp(types.GetCoinType(), rsp.Charge)
	err := bagFunc.ChargeTransOut(head, req.RoomId)
	if err != nil {
		mlog.Info(head, "chargeTransOut err:%v", err)
		return err
	}
	p.GetBaseFunc().RummyQuitRoom()
	return nil
}

func (p *Player) RummyGiveUpReq(head *pb.Head, req *pb.RummyGiveUpReq, rsp *pb.RummyGiveUpRsp) error {
	// 判断是否为断线重连
	if roomId, err := p.GetBaseFunc().GetRummyRealRoomId(head, req.RoomId); err != nil {
		return err
	} else {
		req.RoomId = roomId
	}

	newHead := &pb.Head{
		Dst: framework.NewRoomRouter(req.RoomId, "RummyGame", "RummyGiveUpReq"),
		Uid: head.Uid,
	}

	if err := cluster.Request(newHead, req, rsp); err != nil {
		return err
	}
	bagFunc := p.GetBagFunc()
	types := util.DestructRoomId(req.RoomId)
	bagFunc.AddProp(types.GetCoinType(), rsp.Charge)
	err := bagFunc.ChargeTransOut(head, req.RoomId)
	if err != nil {
		mlog.Info(head, "chargeTransOut err:%v", err)
		return err
	}
	p.GetBaseFunc().RummyQuitRoom()
	return nil
}

func (p *Player) QueryPlayerData(head *pb.Head, req *pb.QueryPlayerDataReq, rsp *pb.QueryPlayerDataRsp) error {
	return p.GetBaseFunc().QueryPlayerData(head, req, rsp)
}

// RummyChangeRoomReq 换桌请求
func (p *Player) RummyChangeRoomReq(head *pb.Head, req *pb.RummyChangeRoomReq, rsp *pb.RummyChangeRoomRsp) error {
	// 同步获取匹配服数据
	types := util.DestructRoomId(req.RoomId)

	if err := cluster.Request(framework.SwapToMatch(head, util.GenMatchId(types), "MatchRummyRoom", "Without"), req, rsp); err != nil {
		return err
	} else if rsp.Head != nil {
		return uerror.ToError(rsp.Head)
	}

	// 退出房间
	newHead := &pb.Head{Src: head.Src, Uid: head.Uid}
	newHead.Dst = framework.NewGameRouter(head.Uid, "Player", "RummyQuitRoomReq")
	quitReq := &pb.RummyQuitRoomReq{
		RoomId: req.RoomId,
	}
	quitRsp := &pb.RummyQuitRoomRsp{}
	if err := p.RummyQuitRoomReq(newHead, quitReq, quitRsp); err != nil {
		return err
	} else if quitRsp.Head != nil {
		return uerror.ToError(rsp.Head)
	}

	// 加入新房间
	newHead.Dst = framework.NewGameRouter(head.Uid, "Player", "RummyJoinRoomReq")
	joinReq := &pb.RummyJoinRoomReq{
		RoomId: rsp.RoomId,
	}
	joinRsp := &pb.RummyJoinRoomRsp{}
	if err := p.RummyJoinRoomReq(newHead, joinReq, joinRsp); err != nil { //todo 优化一次货币转化
		return err
	} else if quitRsp.Head != nil {
		return uerror.ToError(rsp.Head)
	}

	// 转发返回
	rsp.RoomInfo = joinRsp.RoomInfo
	rsp.GaveScore = joinRsp.GaveScore
	rsp.IsReconnect = joinRsp.IsReconnect
	return nil
}

// RummyKickPlayerReq 异步移除玩家请求
func (p *Player) RummyKickPlayerReq(head *pb.Head, req *pb.RummyKickPlayerReq) {
	types := util.DestructRoomId(req.RoomId)
	p.GetBagFunc().AddProp(types.GetCoinType(), req.Charge)

	err := p.GetBagFunc().ChargeTransOut(head, req.RoomId)
	if err != nil {
		mlog.Info(head, "chargeTransOut err:%v", err)
		return
	}

	p.GetBaseFunc().RummyQuitRoom()
}

// RummyMatchReq Rummy分支玩法开启匹配请求
func (p *Player) RummyMatchReq(head *pb.Head, req *pb.RummyMatchReq, rsp *pb.RummyMatchRsp) error {
	if p.GetBaseFunc().GetMatchInfo() != nil {
		mlog.Infof("GetMatchInfo :  %v", p.GetBaseFunc().GetMatchInfo())
		return uerror.New(1, pb.ErrorCode_MATCH_PLAYER_IN_QUEUE, "player match already exists")
	}

	err := p.GetBagFunc().ChargeTransIn(head, uint64(req.GameType&0xFF)<<32|uint64(req.CoinType&0xFF)<<24)
	if err != nil {
		return err
	}
	if p.GetBagFunc().GetProp(req.CoinType) > 0 {
		req.Coin = p.GetBagFunc().GetProp(req.CoinType)
	}

	req.PlayerInfo = p.GetBaseFunc().GetPlayerInfo()
	head.Dst = framework.NewMatchRouter(uint64(req.GameType)<<32|uint64(req.CoinType), "MatchRummyRoom", "RummyMatchReq")
	if err = cluster.Request(head, req, rsp); err != nil {
		return err
	} else if rsp.Head != nil { //匹配失败
		err = p.GetBagFunc().ChargeTransOut(head, uint64(req.GameType&0xFF)<<32|uint64(req.CoinType&0xFF)<<24)
		if err != nil {
			return err
		}
		return uerror.ToError(rsp.Head)
	}

	p.GetBaseFunc().MatchStart(time.Now().Unix()+MATCHTTL, req.GameType, req.CoinType)
	return nil
}

// Update 异步更新玩家状态 todo 房间服操控背包
func (p *Player) Update(head *pb.Head, newStatus *pb.SetPlayerStatus) {
	p.status = newStatus.Status
}

// RummyBatchJoinReq rummy导入玩家到房间服
func (p *Player) RummyBatchJoinReq(head *pb.Head, req *pb.RummyJoinRoomReq, rsp *pb.RummyJoinRoomRsp) error {
	p.GetBaseFunc().MatchStop()
	types := util.DestructRoomId(req.RoomId)
	if p.GetBagFunc().GetProp(types.GetCoinType()) > 0 {
		req.Coin = p.GetBagFunc().GetProp(types.GetCoinType())
		p.GetBagFunc().SubProp(types.GetCoinType(), req.Coin)
	}

	head.Dst = framework.NewRoomRouter(head.Uid, "RummyGameMgr", "JoinRoomReq")
	if err := cluster.Request(head, req, rsp); err != nil {
		return err
	} else if rsp.Head != nil { // 匹配货币只流转于背包
		p.GetBagFunc().AddProp(types.GetCoinType(), req.Coin)
		err = p.GetBagFunc().ChargeTransOut(head, req.RoomId)
		if err != nil {
			return err
		}
		return uerror.ToError(rsp.Head)
	}
	p.GetBaseFunc().RummyJoinRoom(req.RoomId)
	return nil
}

func (p *Player) RummyCancelMatchReq(head *pb.Head, req *pb.RummyCancelMatchReq, rsp *pb.RummyCancelMatchRsp) error {
	if p.GetBaseFunc().GetMatchInfo() == nil {
		return uerror.New(1, pb.ErrorCode_MATCH_FINISH, "match have finished")
	}

	head.Dst = framework.NewMatchRouter(uint64(req.GameType)<<32|uint64(req.CoinType), "MatchRummyRoom", "RummyCancelMatchReq")
	if err := cluster.Request(head, req, rsp); err != nil {
		return err
	} else if rsp.Head != nil { //匹配失败
		return uerror.ToError(rsp.Head)
	}

	p.GetBaseFunc().MatchStop()
	err := p.GetBagFunc().ChargeTransOut(head, uint64(req.GameType&0xFF)<<32|uint64(req.CoinType&0xFF)<<24)
	if err != nil {
		return err
	}
	return nil
}
