package domain

import (
	"corps/pb"
)

const (
	MAILBOX_TL_TIME = 600
	MAX_DWORD_VALUE = 0xffffffff
)

const (
	ELS_Init         = 0 // 默认
	ELS_LoadComplete = 1 // 加载完成
	ELS_SendClient   = 2 // 发送完成
)

// 数据处理接口
type IPlayerFun interface {
	GetPlayerDataType() pb.PlayerDataType           // 获取类型
	RegisterTimer()                                 // 注册定时器
	Load([]byte) error                              // 加载数据(非system数据)
	LoadPlayerDBFinish()                            // db数据加载完成回调，在LoadComplete之前调用
	NewPlayer() error                               // 初始化新开启的模块数据，在LoadComplete之后调用
	LoadComplete()                                  // 加载完成，在NewPlayer之后调用
	Heat()                                          // 心跳包
	PassDay(isDay, isWeek, isMonth bool)            // 是否跨天
	SetUserTypeInfo([]byte) error                   // 设置缓存数据
	SaveDataToClient(pbData *pb.PBPlayerData) error // 深度拷贝数据
	IsSave() bool                                   // 判断是否存储数据
	//UpdateSave(bSave bool)                          // 设置保存状态
}

// system特殊数据接口
type IPlayerSystemFun interface {
	IPlayerFun
	LoadSystem(pbSystem *pb.PBPlayerSystem) // 加载系统数据
}

// 发奖函数协议 +  loot表调用专用
type RewardFunc func(*pb.RpcHead, pb.EmDoingType, interface{}, ...*pb.PBAddItemData) ([]*pb.PBAddItemData, error)
