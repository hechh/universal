package player

import (
	"context"
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common/orm/redis"
	"corps/framework/actor"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/playerMgr/playerFun"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"

	"github.com/golang/protobuf/proto"
)

const (
	MAILBOX_TL_TIME = 600
)

type ELoginState uint32

const (
	ELS_Init         ELoginState = 0 //默认
	ELS_LoadComplete ELoginState = 1 //加载完成
	ELS_SendClient   ELoginState = 2 //发送完成
)

type (
	Player struct {
		actor.Actor
		playerFun.FunCommon

		/*
			pbBase   pb.PBPlayerBase   //基本信息
			pbSystem pb.PBPlayerSystem //系统信息
			pbBag    pb.PBPlayerBag    //背包信息
			pbClient pb.PBClientData   //客户端信息
			pbHero   pb.PBPlayerHero   //伙伴信息
			pbMail   pb.PBPlayerMail   //邮件信息
			pbFriend pb.PBPlayerFriend //好友信息*/

		offlineTime      uint64      //离线时间
		isInGame         bool        //登录游戏
		isNeedSendClient bool        //是否需要同步给客户端
		ELoginState      ELoginState //登录状态

		mapPlayerData  map[pb.PlayerDataType]playerFun.IPlayerFun //玩家数据
		listPlayerData []playerFun.IPlayerFun                     //玩家数据 用于优先级

		mapRequest map[string]func(head *pb.RpcHead, message proto.Message) //请求协议
	}
)

// 注册请求函数
func (this *Player) registerRequest() {
	this.mapRequest["Test"] = this.getPlayerSystemProfessionFun().Test
	this.mapRequest["MainTaskFinishRequest"] = this.getPlayerSystemTaskFun().MainTaskFinishRequest
}

func (this *Player) RegisterPlayerFun(pbType pb.PlayerDataType, fun playerFun.IPlayerFun) {
	fun.Init(pbType, &this.FunCommon)
	this.listPlayerData = append(this.listPlayerData, fun)
	this.mapPlayerData[pbType] = fun
}

func (p *Player) RegisterTimers() {
	p.RegisterTimer((base.R_GateAccountGate_ExpireTime/2)*time.Second, p.UpdateLease) //定时器
	p.RegisterTimer(10*time.Second, p.SaveDB)                                         //定时器
	p.RegisterTimer(3*time.Second, p.Heat)                                            //定时器
	p.RegisterTimer(1*time.Second, p.PassSecond)                                      //定时器
	p.RegisterTimer(10*time.Second, p.TimeReadOffline)                                //定时器
	p.RegisterTimer(30*time.Second, p.getPlayerSystemClientFun().TimeReadClientData)  //定时器
	p.RegisterTimer(20*time.Second, p.getPlayerSystemChampionshipFun().InitRank)      //定时器
}

func (p *Player) Init() {
	p.Actor.Init()
	p.Actor.Start()

	//初始化数据
	p.mapPlayerData = make(map[pb.PlayerDataType]playerFun.IPlayerFun)
	p.listPlayerData = make([]playerFun.IPlayerFun, 0)
	p.FunCommon.MapFun = &p.mapPlayerData
	p.mapRequest = make(map[string]func(head *pb.RpcHead, message proto.Message))
	p.ELoginState = ELS_Init
	//注册fun
	p.registerPlayerFun()
	p.registerRequest()
}

// 重载接口
func (this *Player) Stop() {
	// 保存数据
	this.SaveDB()

	//通知游戏PlayerMgr
	cluster.SendToGame(&pb.RpcHead{Id: this.AccountId}, "PlayerMgr", "OnStopSon")

	// 退出
	this.Actor.Stop()
	plog.Info(" --->>> Player(%d) stop finished offline(%d)", this.GetId(), this.offlineTime)
}

// 客户端通信
func (this *Player) updateOffline() {
	this.offlineTime = base.GetNow() + MAILBOX_TL_TIME
}

// 客户端通信
func (this *Player) getPlayerFun(emType pb.PlayerDataType) playerFun.IPlayerFun {
	fun, ok := this.mapPlayerData[emType]
	if !ok {
		return nil
	}

	return fun
}

// 客户端通信
func (p *Player) SendToClient(head *pb.RpcHead, msg proto.Message) {
	head.ClusterId = 0
	cluster.SendToClient(head, msg, cfgEnum.ErrorCode_Success)
}

// 更新连接
func (this *Player) UpdateLease() {
	if base.GetNow() >= this.offlineTime {
		this.SaveDB()

		//需要销毁actor
		this.Stop()
	}
}

// 退出
func (this *Player) stop() {
	//需要销毁actor
	this.Actor.Stop()

	//通知游戏PlayerMgr
	cluster.SendToGame(&pb.RpcHead{Id: this.AccountId}, "PlayerMgr", "OnStopSon")

	plog.Info("(this *Player)stop", this.offlineTime, base.GetNow())

	//通知游戏DB退出
	//cluster.SendToGame(&pb.RpcHead{Id: this.AccountId}, "PlayerMgr", "OnStopSon")
}

// 存库
func (p *Player) SaveDB() {
	if !p.isInGame {
		return
	}

	//存储数据
	for pbType, fun := range p.mapPlayerData {
		if pbType == pb.PlayerDataType_System || pbType >= pb.PlayerDataType_SystemCommon {
			continue
		}

		fun.Save(true)
	}

	//系统数据单独存储
	p.SavePlayerSystem()
}

// 心跳处理
func (this *Player) Heat() {
	for _, fun := range this.mapPlayerData {
		fun.Heat()
	}
}

// 跨天
func (this *Player) PassSecond() {
	if arrs := this.getPlayerBaseFun().CheckDay(); len(arrs) >= 3 && (arrs[0] || arrs[1] || arrs[2]) {
		for _, fun := range this.mapPlayerData {
			fun.PassDay(arrs[0], arrs[1], arrs[2])
		}

		//登录天数成就
		this.getPlayerSystemTaskFun().TriggerAchieve(&pb.RpcHead{}, cfgEnum.AchieveType_LoginDay, 1)
	}
}

// 系统存库
func (this *Player) SavePlayerSystem() {

	bSave := false
	for pbType, fun := range this.mapPlayerData {
		if pbType < pb.PlayerDataType_SystemCommon {
			continue
		}
		if fun.IsSave() {
			bSave = true
			break
		}
	}

	if bSave {
		pbSystem := new(pb.PBPlayerSystem)
		for pbType, fun := range this.mapPlayerData {
			if pbType < pb.PlayerDataType_SystemCommon {
				continue
			}

			fun.SaveSystem(pbSystem)
		}

		buff, _ := proto.Marshal(pbSystem)
		cluster.SendToDb(&pb.RpcHead{Id: this.AccountId}, "DbPlayerMgr", "SavePlayerDB", pb.PlayerDataType_System, buff, false)
	}
}

func (p *Player) SaveToPb(pbPlayer *pb.PBPlayerData) {
	//存储数据
	for _, fun := range p.mapPlayerData {
		fun.SaveDataToClient(pbPlayer)
	}
}

func (p *Player) SaveToTypePb(pbPlayer *pb.PBPlayerData, emType pb.PlayerDataType) {
	if emType < pb.PlayerDataType_Max {
		if fun, ok := p.mapPlayerData[emType]; ok {
			fun.SaveDataToClient(pbPlayer)
		}
	} else {
		for i := pb.PlayerDataType_SystemCommon; i < pb.PlayerDataType_SystemMax; i++ {
			if fun, ok := p.mapPlayerData[i]; ok {
				fun.SaveDataToClient(pbPlayer)
			}
		}
	}
}

// 注册player
func (this *Player) registerPlayerFun() {
	for i := pb.PlayerDataType_Crystal; i < pb.PlayerDataType_Max; i++ {
		switch i {
		case pb.PlayerDataType_Base:
			this.RegisterPlayerFun(i, new(playerFun.PlayerBaseFun))
		case pb.PlayerDataType_System:
			for j := pb.PlayerDataType_SystemCommon; j < pb.PlayerDataType_SystemMax; j++ {
				switch j {
				case pb.PlayerDataType_SystemCommon:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemCommonFun))
					break
				case pb.PlayerDataType_SystemProfession:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemProfessionFun))
					break
				case pb.PlayerDataType_SystemBattle:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemBattleFun))
					break
				case pb.PlayerDataType_SystemBox:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemBoxFun))
					break
				case pb.PlayerDataType_SystemTask:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemTaskFun))
					break
				case pb.PlayerDataType_SystemShop:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemShopFun))
					break
				case pb.PlayerDataType_SystemDraw:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemDrawFun))
					break
				case pb.PlayerDataType_SystemCharge:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemChargeFun))
					break
				case pb.PlayerDataType_SystemGene:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemGeneFun))
					break
				case pb.PlayerDataType_SystemOffline:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemOfflineFun))
					break
				case pb.PlayerDataType_SystemHookTech:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemHookTechFun))
					break
				case pb.PlayerDataType_SystemRepair:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemRepairFun))
					break
				case pb.PlayerDataType_SystemSevenDay:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemSevenDayFun))
					break
				case pb.PlayerDataType_SystemActivity:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemActivityFun))
					break
				case pb.PlayerDataType_SystemWorldBoss:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemWorldBossFun))
					break
				case pb.PlayerDataType_SystemChampionship:
					this.RegisterPlayerFun(j, new(playerFun.PlayerSystemChampionshipFun))
					break
				}
			}
		case pb.PlayerDataType_Bag:
			this.RegisterPlayerFun(i, new(playerFun.PlayerBagFun))
		case pb.PlayerDataType_Equipment:
			this.RegisterPlayerFun(i, new(playerFun.PlayerEquipmentFun))
		case pb.PlayerDataType_Client:
			this.RegisterPlayerFun(i, new(playerFun.PlayerClientFun))
		case pb.PlayerDataType_Hero:
			this.RegisterPlayerFun(i, new(playerFun.PlayerHeroFun))
		case pb.PlayerDataType_Mail:
			this.RegisterPlayerFun(i, new(playerFun.PlayerMailFun))
		case pb.PlayerDataType_Crystal:
			this.RegisterPlayerFun(i, new(playerFun.PlayerCrystalFun))
		}
	}
}

func (p *Player) GetPlayerId() uint64 {
	return p.AccountId
}

// 同步所有数据到客户端
func (p *Player) sendAllInfoToClient(head *pb.RpcHead) {
	if !p.isInGame {
		return
	}

	pbResponse := &pb.AllPlayerInfoNotify{
		PacketHead: &pb.IPacket{},
	}

	//初始化玩家数据
	for key, fun := range p.mapPlayerData {
		if key == pb.PlayerDataType_Hero {
			continue
		}

		pbResponse.PlayerData = &pb.PBPlayerData{}
		pbResponse.Mark = 0
		pbResponse.Mark = base.SetBit32(pbResponse.Mark, uint32(key), true)
		fun.SaveDataToClient(pbResponse.PlayerData)
		p.SendToClient(head, pbResponse)
		plog.Trace("loginSuccess %s: %v", key.String(), pbResponse)
	}

	//英雄放最后
	fun, ok := p.mapPlayerData[pb.PlayerDataType_Hero]
	if ok {
		pbResponse.PlayerData = &pb.PBPlayerData{}
		pbResponse.Mark = 0
		pbResponse.Mark = base.SetBit32(pbResponse.Mark, uint32(pb.PlayerDataType_Hero), true)
		fun.SaveDataToClient(pbResponse.PlayerData)
		p.SendToClient(head, pbResponse)
		plog.Trace("loginSuccess %s: %v", pb.PlayerDataType_Hero.String(), pbResponse)
	}

	//完成标记
	pbResponse.Mark = 0
	pbResponse.PlayerData = &pb.PBPlayerData{}
	//plog.Debug("sendAllInfoToClient: %+v", pbResponse.PlayerData)
	p.SendToClient(head, pbResponse)
}

// 玩家断开链接
func (p *Player) Logout(ctx context.Context, playerId uint64) {
	plog.Info("[%d] 断开链接", playerId)
	p.SaveDB()
}

// 登录完成
func (this *Player) loginSuccess(head *pb.RpcHead) {
	this.ELoginState = ELS_LoadComplete
	plog.Trace("head: %v", head)
	//通知网关 设置状态
	serverID := this.getPlayerBaseFun().GetServerId()
	cluster.SendToGate(head, "AccountMgr", "OnGameLoginResponse", serverID)

	//同步所有数据到客户端
	this.sendAllInfoToClient(head)
	this.isNeedSendClient = false
	this.ELoginState = ELS_SendClient
	//通知各个系统加载完成
	//初始化玩家数据
	for _, fun := range this.listPlayerData {
		fun.LoadComplete()
	}

	//更新离线数据
	this.getPlayerSystemOfflineFun().UpdateLoginTime()
	this.updateOffline()

	//通知成功
	this.SendToClient(head, &pb.LoginResponse{
		PacketHead: &pb.IPacket{
			Id:   this.AccountId,
			Code: uint32(cfgEnum.ErrorCode_Success),
		},
		Time: base.GetNow(),
	})

	// 数据上报
	/*
		report.Send(head, &report.ReportLogin{
			Doing:      uint32(pb.EmDoingType_EDT_Login),
			PlayerID:   this.AccountId,
			PlayerName: this.getPlayerBaseFun().GetDisplay().PlayerName,
			Level:      this.getPlayerBaseFun().GetPlayerLevel(),
			Vip:        this.getPlayerBaseFun().GetDisplay().VipLevel,
			LoginTime:  base.GetNow(),
			//LastLoginTime:   this.getPlayerBaseFun().GetPlayerBase().LastModifyTime,
			//LastLogoutTime: this.LastLogoutTime,
		})*/
}

// 设置玩家数据
func (this *Player) DipSetUserTypeInfo(ctx context.Context, emType pb.PlayerDataType, strData string) bool {
	fun := this.getPlayerFun(emType)
	if fun == nil {
		plog.Error("(this *Player) DipSetUserTypeInfo fun acccountid:%d,type:%d", this.AccountId, emType)
		return false
	}
	marshaler := jsonpb.Unmarshaler{}
	pbData := fun.GetProtoPtr()
	if marshaler.Unmarshal(strings.NewReader(strData), pbData) != nil {
		return false
	}

	if !fun.SetUserTypeInfo(pbData) {
		plog.Error("(this *Player) DipSetUserTypeInfo set acccountid:%d,type:%d", this.AccountId, emType)
		return false
	}

	//返回结果
	head := this.GetRpcHead(ctx)
	head.DestServerType = pb.SERVICE_Dip
	cluster.ReplyMsgTo(head, emType)

	return false
}

func (p *Player) ShutDown(ctx context.Context) {
	//立即入库
	p.SaveDB()

	plog.Info("(p *Player) ShutDown [%d] ", p.AccountId)
	p.Stop()
}

// 晶核系统
func (this *Player) getPlayerCrystalFun() *playerFun.PlayerCrystalFun {
	return this.getPlayerFun(pb.PlayerDataType_Crystal).(*playerFun.PlayerCrystalFun)
}

// 玩家基本数据
func (this *Player) getPlayerBaseFun() *playerFun.PlayerBaseFun {
	return this.getPlayerFun(pb.PlayerDataType_Base).(*playerFun.PlayerBaseFun)
}

// 玩家背包数据
func (this *Player) getPlayerBagFun() *playerFun.PlayerBagFun {
	return this.getPlayerFun(pb.PlayerDataType_Bag).(*playerFun.PlayerBagFun)
}

// 邮件数据
func (this *Player) getPlayerMailFun() *playerFun.PlayerMailFun {
	return this.getPlayerFun(pb.PlayerDataType_Mail).(*playerFun.PlayerMailFun)
}

func (this *Player) getPlayerSystemChampionshipFun() *playerFun.PlayerSystemChampionshipFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemChampionship).(*playerFun.PlayerSystemChampionshipFun)
}

// 玩家系统通用数据
func (this *Player) getPlayerSystemCommonFun() *playerFun.PlayerSystemCommonFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemCommon).(*playerFun.PlayerSystemCommonFun)
}

func (this *Player) getPlayerSystemSevenDayFun() *playerFun.PlayerSystemSevenDayFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemSevenDay).(*playerFun.PlayerSystemSevenDayFun)
}
func (this *Player) getPlayerSystemActivityFun() *playerFun.PlayerSystemActivityFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemActivity).(*playerFun.PlayerSystemActivityFun)
}
func (this *Player) getPlayerActivityGrowRoadFun() *playerFun.PlayerActivityGrowRoad {
	return this.getPlayerSystemActivityFun().GetActivityGrowRoadFun()
}
func (this *Player) getPlayerActivityAdventureFun() *playerFun.PlayerActivityAdventure {
	return this.getPlayerSystemActivityFun().GetActivityAdventureFun()
}

func (this *Player) getPlayerSystemWorldBossFun() *playerFun.PlayerSystemWorldBossFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemWorldBoss).(*playerFun.PlayerSystemWorldBossFun)
}

// 商店系统
func (this *Player) getPlayerSystemShopFun() *playerFun.PlayerSystemShopFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemShop).(*playerFun.PlayerSystemShopFun)
}

// 抽奖系统
func (this *Player) getPlayerSystemDrawFun() *playerFun.PlayerSystemDrawFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemDraw).(*playerFun.PlayerSystemDrawFun)
}

// 星源系统
func (this *Player) getPlayerSystemHookTechFun() *playerFun.PlayerSystemHookTechFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemHookTech).(*playerFun.PlayerSystemHookTechFun)
}

// 充值系统
func (this *Player) getPlayerSystemChargeFun() *playerFun.PlayerSystemChargeFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemCharge).(*playerFun.PlayerSystemChargeFun)
}
func (this *Player) GetChargeBPFun() *playerFun.PlayerSystemChargeBP {
	return this.getPlayerSystemChargeFun().GetChargeBPFun()
}
func (this *Player) GetChargeCardFun() *playerFun.PlayerSystemChargeCard {
	return this.getPlayerSystemChargeFun().GetChargeCardFun()
}

// 玩家系统职业数据
func (this *Player) getPlayerSystemProfessionFun() *playerFun.PlayerSystemProfessionFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemProfession).(*playerFun.PlayerSystemProfessionFun)
}
func (this *Player) getPlayerSystemBoxFun() *playerFun.PlayerSystemBoxFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemBox).(*playerFun.PlayerSystemBoxFun)
}

// 玩家系统战斗数据
func (this *Player) getPlayerSystemBattleFun() *playerFun.PlayerSystemBattleFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemBattle).(*playerFun.PlayerSystemBattleFun)
}

// 玩家系统战斗数据
func (this *Player) getPlayerSystemBattleNormalFun() *playerFun.PlayerSystemBattleNormalFun {
	return this.getPlayerSystemBattleFun().GetBattleNoramlFun()
}
func (this *Player) getPlayerSystemBattleHookFun() *playerFun.PlayerSystemBattleHookFun {
	return this.getPlayerSystemBattleFun().GetBattleHookFun()
}

// 玩家系统任务数据
func (this *Player) getPlayerSystemTaskFun() *playerFun.PlayerSystemTaskFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemTask).(*playerFun.PlayerSystemTaskFun)
}
func (this *Player) getPlayerSystemClientFun() *playerFun.PlayerClientFun {
	return this.getPlayerFun(pb.PlayerDataType_Client).(*playerFun.PlayerClientFun)
}

// 玩家系统装备数据
func (this *Player) GetPlayerEquipmentFun() *playerFun.PlayerEquipmentFun {
	return this.getPlayerFun(pb.PlayerDataType_Equipment).(*playerFun.PlayerEquipmentFun)
}

// 玩家伙伴数据
func (this *Player) getPlayerHeroFun() *playerFun.PlayerHeroFun {
	return this.getPlayerFun(pb.PlayerDataType_Hero).(*playerFun.PlayerHeroFun)
}

// 玩家系统宝箱数据
func (this *Player) getPlayerSystemBox() *playerFun.PlayerSystemBoxFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemBox).(*playerFun.PlayerSystemBoxFun)
}

// 获取数据
func (this *Player) DipGetUserInfo(ctx context.Context, dataType pb.PlayerDataType) {
	head := this.GetRpcHead(ctx)
	pbPlayerData := &pb.PBPlayerData{}
	fun := this.getPlayerFun(dataType)
	if fun == nil {
		return
	}

	fun.SaveDataToClient(pbPlayerData)
	if dataType == pb.PlayerDataType_Crystal {
		pbPlayerData.Crystal.Effects = nil
	}

	head.DestServerType = pb.SERVICE_Dip
	cluster.ReplyMsgTo(head, pbPlayerData)
}

// 从redis读取邮件
func (this *Player) TimeReadOffline() {
	//先读取个人邮件
	redisGame := redis.GetRedisByAccountID(this.AccountId)
	if redisGame == nil {
		return
	}

	//遍历邮件
	strKey := fmt.Sprintf("%s%d", base.ERK_GamePlayerOffline, this.AccountId)
	for {
		strData := redisGame.RPop(strKey)
		if strData == "" {
			break
		}

		pbData := &pb.PBOfflineData{}
		err := json.Unmarshal([]byte(strData), pbData)
		if err != nil {
			continue
		}

		switch pbData.OfflineType {
		case pb.EmPlayerOfflineType_EPOT_Mail:
			this.getPlayerMailFun().AddMail(&pb.RpcHead{Id: this.AccountId}, pbData.Mail)
			break
		case pb.EmPlayerOfflineType_EPOT_Item:
			this.getPlayerBagFun().AddPbItems(&pb.RpcHead{Id: this.AccountId}, pbData.Item, pbData.DoingType, pbData.Notify)
			//this.getPlayerBagFun().AddItem(&pb.RpcHead{Id: this.AccountId}, pbData.Item.Kind, pbData.Item.Id, pbData.Item.Count, pbData.Item.DoingType, false, pbData.Item.Params...)
			break
		default:
			plog.Info("Player TimeReadOffline Error: %d", pbData.OfflineType)
			break
		}
	}
}

// 注册函数的 map
func (this *Player) doRequest(strName string, ctx context.Context, pbRequest proto.Message) {
	head := this.GetRpcHead(ctx)
	fun, ok := this.mapRequest[strName]
	if !ok {
		plog.Error("doRequest cannot find id:%d name:%s", this.AccountId, strName)
		return
	}
	fun(head, pbRequest)

	//调用结束 需要同步道具改变
	this.getPlayerBagFun().SendChangeItem(head)
}
