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
	"corps/server/game/player/domain"
	"corps/server/game/playerMgr/playerFun"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
)

type Player struct {
	actor.Actor
	*playerFun.PlayerFun
	offlineTime uint64 // 离线时间
	isValid     bool   // 游戏缓存数是否有效
	loginState  int32  // 登录状态
}

func NewPlayer() (item *Player) {
	item = new(Player)
	item.PlayerFun = playerFun.NewPlayerFun()
	// 注册业务模块
	item.RegisterIPlayerFun(new(playerFun.PlayerBaseFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerBagFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerEquipmentFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerClientFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerHeroFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerMailFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerCrystalFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemCommonFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemProfessionFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemBattleFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemBoxFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemTaskFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemShopFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemDrawFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemChargeFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemGeneFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemOfflineFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemHookTechFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemRepairFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemSevenDayFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemActivityFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemWorldBossFun))
	item.RegisterIPlayerFun(new(playerFun.PlayerSystemChampionshipFun))
	return
}

// 设置登录状态
func (this *Player) SetState(val int32) {
	this.loginState = val
}

// 获取登录状态
func (this *Player) GetState() int32 {
	return this.loginState
}

// 更新离线时间
func (this *Player) SetOffline(val uint64) {
	this.offlineTime = val
}

func (this *Player) UpdateOffline() {
	this.offlineTime = base.GetNow() + domain.MAILBOX_TL_TIME
}

// 是否有效
func (this *Player) IsValid() bool {
	return this.isValid
}

// 设置玩家缓存数据有效状态
func (this *Player) SetValid(flag bool) {
	this.isValid = flag
}

// 初始化actor(实现actor接口)
func (this *Player) Init() {
	this.Actor.Init()
	this.Actor.Start()
	this.FunMgr.Init()
	this.loginState = domain.ELS_Init
}

// 关闭player协程(实现actor接口)
func (this *Player) Stop() {
	// 保存数据
	this.save()
	//通知游戏PlayerMgr
	cluster.SendToGame(&pb.RpcHead{Id: this.GetId()}, "PlayerMgr", "OnStopSon")
	// 退出
	this.Actor.Stop()
	plog.Info(" --->>> Player(%d) stop finished offline(%d)", this.GetId(), this.offlineTime)
}

// 退出player协程
func (p *Player) ShutDown(ctx context.Context) {
	plog.Info("(p *Player) ShutDown [%d] ", p.GetId())
	p.Stop()
}

// 拷贝玩家所有数据
func (this *Player) Copy(types ...pb.PlayerDataType) (result *pb.PBPlayerData) {
	result = &pb.PBPlayerData{}
	if len(types) <= 0 {
		this.Walk(func(typ pb.PlayerDataType, fun domain.IPlayerFun) bool {
			fun.CopyTo(result)
			return true
		})
	} else {
		for _, typ := range types {
			this.GetIPlayerFun(typ).CopyTo(result)
		}
	}
	return
}

// 初始化定时器
func (this *Player) InitTimer() {
	this.RegisterTimer((base.R_GateAccountGate_ExpireTime/2)*time.Second, this.offline) //定时器
	this.RegisterTimer(10*time.Second, this.save)                                       //定时器
	this.RegisterTimer(3*time.Second, this.heart)                                       //定时器
	this.RegisterTimer(1*time.Second, this.passSecond)                                  //定时器
	this.RegisterTimer(10*time.Second, this.timeReadOffline)                            //定时器
	// 初始化业务层定时器
	this.Walk(func(typ pb.PlayerDataType, fun domain.IPlayerFun) bool {
		fun.RegisterTimer()
		return true
	})
}

// 玩家下线(定时器定时执行)
func (this *Player) offline() {
	if base.GetNow() >= this.offlineTime {
		this.save()
		//需要销毁actor
		this.Stop()
	}
}

// 保存玩家数据到mysql数据库(定时器定时执行)
func (this *Player) save() {
	// 无效数据不处理
	if !this.IsValid() {
		return
	}
	//存储数据
	isSave := false
	tmps := []domain.IPlayerSystemFun{}
	this.Walk(func(typ pb.PlayerDataType, fun domain.IPlayerFun) bool {
		switch vv := fun.(type) {
		case domain.IPlayerSystemFun:
			isSave = isSave || vv.IsSave()
			tmps = append(tmps, vv)
		default:
			vv.Save(true)
		}
		return true
	})
	if isSave {
		data := &pb.PBPlayerSystem{}
		for _, fun := range tmps {
			fun.SaveSystem(data)
		}
		buff, _ := proto.Marshal(data)
		cluster.SendToDb(&pb.RpcHead{Id: this.GetId()}, "DbPlayerMgr", "SavePlayerDB", pb.PlayerDataType_System, buff, false)
	}
}

// 心跳处理(定时器定时执行)
func (this *Player) heart() {
	this.Walk(func(typ pb.PlayerDataType, fun domain.IPlayerFun) bool {
		fun.Heat()
		return true
	})
}

// 跨天(定时器定时执行)
func (this *Player) passSecond() {
	if arrs := playerFun.GetPlayerBaseFun().CheckDay(); len(arrs) >= 3 && (arrs[0] || arrs[1] || arrs[2]) {
		this.Walk(func(typ pb.PlayerDataType, fun domain.IPlayerFun) bool {
			fun.PassDay(arrs[0], arrs[1], arrs[2])
			return true
		})

		//登录天数成就
		playerFun.GetPlayerSystemTaskFun().TriggerAchieve(&pb.RpcHead{}, cfgEnum.AchieveType_LoginDay, 1)
	}
}

// (定时器定时执行)
func (this *Player) timeReadOffline() {
	//先读取个人邮件
	redisGame := redis.GetRedisByAccountID(this.GetId())
	if redisGame == nil {
		return
	}
	//遍历邮件
	strKey := fmt.Sprintf("%s%d", base.ERK_GamePlayerOffline, this.GetId())
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
			playerFun.GetPlayerMailFun().AddMail(&pb.RpcHead{Id: this.GetId()}, pbData.Mail)
			break
		case pb.EmPlayerOfflineType_EPOT_Item:
			playerFun.GetPlayerBagFun().AddPbItems(&pb.RpcHead{Id: this.GetId()}, pbData.Item, pbData.DoingType, pbData.Notify)
			//this.getPlayerBagFun().AddItem(&pb.RpcHead{Id: this.GetId()}, pbData.Item.Kind, pbData.Item.Id, pbData.Item.Count, pbData.Item.DoingType, false, pbData.Item.Params...)
			break
		default:
			plog.Info("Player TimeReadOffline Error: %d", pbData.OfflineType)
			break
		}
	}
}
