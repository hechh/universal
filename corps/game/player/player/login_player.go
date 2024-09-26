package player

import (
	"context"
	"corps/base"
	"corps/base/cfgEnum"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/player/domain"
	"corps/server/game/player/internal/playerFun"

	"google.golang.org/protobuf/proto"
)

/*  登录流程说明
1、db/DbPlayerMgr.LoadPlayerDB()
	db/player.LoadPlayerDB()			从数据库读取数据
	db/player.sendPlayerToGame()		发送db服务缓存数据到game服务
2、game/player.LoadPlayerDBType()：		game服务加载db服务传过来的数据
	game/PlayerFun.LoadSystem()：		加载PlayerSystem中的数据
	game/PlayerFun.Load()：				加载非PlayerSystem数据
3、game/player.LoadPlayerDBFinish()：	玩家数据加载完成
	game/player.NewPlayer()：			新玩家初始化数据(db服务没有新玩家数据)
	game/player.loginSuccess()：
	gate/account.OnGameLoginResponse()：设置玩家登录状态
	game/player.sendAllInfoToClient()：	把PlayerFun数据同步给客户端
	game/playerFun.LoadComplete()：		通知各个系统加载完成
*/

// 玩家登录(向db服务请求数据)
func (this *Player) Login(ctx context.Context) {
	head := this.GetRpcHead(ctx)
	this.SetOffline(0)
	// 通知db加载玩家数据
	if !cluster.SendToDb(head, "DbPlayerMgr", "LoadPlayerDB") {
		plog.Info("玩家登录失败 db找不到 Login id:%d", this.GetId())
		cluster.SendToClient(head, &pb.LoginResponse{
			PacketHead: &pb.IPacket{
				Id:   this.GetId(),
				Code: uint32(cfgEnum.ErrorCode_ServerBusy),
			},
		}, cfgEnum.ErrorCode_Success)
		return
	}
	plog.Info("玩家登录成功 Login id:%d", this.GetId())
}

// 断线重连
func (this *Player) ReLogin(ctx context.Context) {
	// 判单缓存数据是否有效
	if !this.IsValid() {
		this.Login(ctx)
		return
	}
	//发送给网关
	head := this.GetRpcHead(ctx)
	this.loginSuccess(head) // 初始化
	this.UpdateOffline()    // 登录成功，设置离线时间
	plog.Info("玩家重连登录成功 (p *Player) ReLogin id: %d", head.Id)
}

// 接受db服务传送过来的数据，初始化各个游戏模块
func (this *Player) LoadPlayerDBType(ctx context.Context, dataType pb.PlayerDataType, pData []byte) {
	head := this.GetRpcHead(ctx)
	plog.Trace("玩家数据加载成功 LoadPlayerDBType head: %v, PlayerDataType: %v, data: %v", head, dataType, pData)

	// fun加载数据
	switch dataType {
	case pb.PlayerDataType_System:
		data := &pb.PBPlayerSystem{}
		proto.Unmarshal(pData, data)
		this.Walk(func(typ pb.PlayerDataType, fun domain.IPlayerFun) bool {
			switch vv := fun.(type) {
			case domain.IPlayerSystemFun:
				vv.LoadSystem(data)
			}
			return true
		})
	default:
		this.GetIPlayerFun(dataType).Load(pData)
	}
}

// 所有fun加载数据已经完成的回调
func (this *Player) LoadPlayerDBFinish(ctx context.Context) {
	head := this.GetRpcHead(ctx)
	plog.Trace("玩家:%d 数据加载成功 LoadPlayerDBFinish, head: %v", this.GetId(), head)

	//需要通知各个系统加载数据库加载完成
	this.Walk(func(typ pb.PlayerDataType, fun domain.IPlayerFun) bool {
		fun.LoadPlayerDBFinish()
		plog.Trace("uid: %d, type: %s LoadPlayerDBFinish", this.GetId(), typ.String())
		return true
	})

	//判断是否是新系统
	listNewPlayerType := playerFun.GetPlayerBaseFun().GetNewPlayerTypeList()
	for _, fun := range this.GetIPlayerFunList() {
		if base.ArrayContainsValue(listNewPlayerType, uint32(fun.GetPlayerDataType())) {
			continue
		}
		fun.NewPlayer()
		plog.Trace("uid: %d, type: %s NewPlayer", this.GetId(), fun.GetPlayerDataType().String())
		playerFun.GetPlayerBagFun().AddNewPlayerTypeList(uint32(fun.GetPlayerDataType()))
	}

	//发送给网关
	this.loginSuccess(head)

	// 开启定时器
	this.InitTimer()
}

// 登录完成
func (this *Player) loginSuccess(head *pb.RpcHead) {
	plog.Trace("head: %v", head)
	this.SetState(domain.ELS_LoadComplete)

	// 通知网关设置区服ID
	serverID := playerFun.GetPlayerBaseFun().GetServerId()
	cluster.SendToGate(head, "AccountMgr", "OnGameLoginResponse", serverID)

	// 推送数据到客户端
	this.sendAllInfoToClient(head)
	this.SetState(domain.ELS_SendClient)

	//通知各个系统加载完成, 初始化玩家数据
	for _, fun := range this.GetIPlayerFunList() {
		fun.LoadComplete()
	}

	//更新离线数据
	playerFun.GetPlayerSystemOfflineFun().UpdateLoginTime()
	this.UpdateOffline()

	//通知成功
	cluster.SendToClient(head, &pb.LoginResponse{
		PacketHead: &pb.IPacket{
			Id:   this.GetId(),
			Code: uint32(cfgEnum.ErrorCode_Success),
		},
		Time: base.GetNow(),
	}, cfgEnum.ErrorCode_Success)
}

// 同步所有数据到客户端
func (this *Player) sendAllInfoToClient(head *pb.RpcHead) {
	// 无效数据不处理
	if !this.IsValid() {
		return
	}
	//初始化玩家数据
	rsp := &pb.AllPlayerInfoNotify{PacketHead: &pb.IPacket{Id: this.GetId()}}
	this.Walk(func(typ pb.PlayerDataType, fun domain.IPlayerFun) bool {
		if typ == pb.PlayerDataType_Hero {
			return true
		}
		rsp.PlayerData = &pb.PBPlayerData{}
		rsp.Mark = 0
		rsp.Mark = base.SetBit32(rsp.Mark, uint32(typ), true)
		fun.CopyTo(rsp.PlayerData)
		cluster.SendToClient(head, rsp, cfgEnum.ErrorCode_Success)
		plog.Trace("loginSuccess %s: %v", typ.String(), rsp)
		return true
	})
	//英雄放最后
	if fun := this.GetIPlayerFun(pb.PlayerDataType_Hero); fun != nil {
		rsp.PlayerData = &pb.PBPlayerData{}
		rsp.Mark = 0
		rsp.Mark = base.SetBit32(rsp.Mark, uint32(pb.PlayerDataType_Hero), true)
		fun.CopyTo(rsp.PlayerData)
		cluster.SendToClient(head, rsp, cfgEnum.ErrorCode_Success)
		plog.Trace("loginSuccess %s: %v", pb.PlayerDataType_Hero.String(), rsp)
	}
	//完成标记
	rsp.Mark = 0
	rsp.PlayerData = &pb.PBPlayerData{}
	cluster.SendToClient(head, rsp, cfgEnum.ErrorCode_Success)
}
