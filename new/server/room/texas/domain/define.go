package domain

import (
	"poker_server/common/pb"

	"github.com/golang/protobuf/proto"
)

const (
	HOLD_TIME = 2000
)

// 状态接口
type IState interface {
	OnEnter(int64, pb.TexasGameState, IRoom)
	OnExit(int64, pb.TexasGameState, IRoom)
	OnTick(int64, pb.TexasGameState, IRoom) pb.TexasGameState
}

// 状态机接口
type IMachine interface {
	Handle(int64, IRoom)      // 执行状态机
	GetStateStartTime() int64 // 获取状态机开始时间
	SetStateStartTime(int64)  // 设置状态机开始时间
}

// 比较牌型大小
type ICompare interface {
	Get(...uint32) (pb.TexasCardType, uint32, []uint32)
}

type ITexas interface {
	GetCompare() ICompare                                                       // 获取比牌算法
	GetCursor() *pb.TexasPlayerData                                             // 获取当前操作玩家
	GetDealer() *pb.TexasPlayerData                                             // 获取庄家
	GetSmall() *pb.TexasPlayerData                                              // 获取小盲注玩家
	GetBig() *pb.TexasPlayerData                                                // 获取大盲注玩家
	Walk(int, func(*pb.TexasPlayerData) bool)                                   // 遍历游戏玩家
	WalkPrev(int, func(*pb.TexasPlayerData) bool)                               // 遍历上一个玩家
	GetPrev(int, pb.TexasGameState, ...pb.TexasOperateType) *pb.TexasPlayerData // 获取上一个玩家
	GetNext(int, pb.TexasGameState, ...pb.TexasOperateType) *pb.TexasPlayerData // 获取下一个玩家
	GetPlayers(pb.TexasGameState) []*pb.TexasPlayerData                         // 获取指定游戏状态的所有玩家
	Operate(*pb.TexasPlayerData, pb.TexasOperateType, int64)                    // 玩家下注操作
	UpdateMain(...*pb.TexasPlayerData)                                          // 更新主池
	UpdateSide(...*pb.TexasPlayerData)                                          // 更新边池
	Reward(int, func(uint64, int64)) int64                                      // 结算
	Shuffle(int)                                                                // 洗牌
	Deal(uint32, func(uint32, uint32))                                          // 发牌
	UpdateBest(*pb.TexasPlayerData, ICompare, []uint32)                         // 更新最大牌型
}

type ISend interface {
	Report(proto.Message) error                                          // 上报数据
	NotifyToClient(pb.TexasEventType, proto.Message) error               // 广播游戏事件
	NotifyToPlayerClient(uint64, pb.TexasEventType, proto.Message) error // 发送数据到指定玩家
	SendToClient(uint64, pb.CMD, proto.Message) error                    // 发送数据到客户端
}

type IRoom interface {
	ITexas
	ISend
	Change()                             // 标记房间数据有变化
	IsChange() bool                      // 是否有变化
	OnTick(int64)                        // 定时调用
	SetRecord(*pb.TexasGameRecord)       // 获取游戏记录
	GetRecord() *pb.TexasGameRecord      // 获取游戏记录
	GetTexasRoomData() *pb.TexasRoomData // 获取房间数据
	GetCurStateTTL() int64               // 获取状态时间
	GetStateStartTime() int64            // 获取状态机开始时间
	SetStateStartTime(int64)             // 设置状态机开始时间
	GetCurState() pb.TexasGameState      // 获取当前状态
	GetNextState() pb.TexasGameState     // 获取下一个状态
}
