package rummy

import (
	"math/rand"
	"poker_server/framework"
	"time"
)

type StateDuration uint32

// 状态机时间默认值
const (
	ReadyDuration     StateDuration = 10 //准备倒计时
	HostDuration      StateDuration = 3  // 定庄
	DealDuration      StateDuration = 5  // 发牌
	PlayDuration      StateDuration = 30 //用户正常倒计时
	ShortPlayDuration StateDuration = 30 //用户离线倒计时
	FixDuration       StateDuration = 60 //最后调整牌时间
	SettleDuration    StateDuration = 10 //结束放排行榜倒计时
)

type RummyRoom struct {
	framework.Actor
	//房间配置
	RoomName string //房间名
	//todo RoomType RummyDeskType
	BaseScore  int64 //底分或门票
	Prize      int64 // Prize 奖池金币
	TotalRound int32 // TotalRound 总局数
	// MaxPlayers 最大进入玩家
	MaxPlayers uint32
	// MinStartPlayers 最少开局人数
	MinStartPlayers uint32
	// 是否允许旁观
	AllowWatch bool
	// MinCoin 最小准入金币
	MinCoin int64
	// MaxCoin 最大准入金币
	MaxCoin          int64
	DisableMathLimit bool //限制匹配
	// IsTimeoutContinue 是否超时继续
	IsTimeoutContinue bool

	// todo 状态机配置

	//游戏数据
	Round          int32             // 当前回合
	HostPlayerID   uint64            // 庄
	seats          []uint64          // 座位集合
	PlayerIDs      []uint64          // 玩家id集合
	PlayerMap      map[uint64]string // 玩家数据
	GameStartTime  time.Time         //游戏开始时间
	RoundStartTime time.Time         // RoundStartTime 回合开始时间
	DeskRand       *rand.Rand        //随机数种子
	// todo 游戏记录

}
