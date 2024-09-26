package domain

import (
	"corps/base/cfgEnum"
	"corps/pb"
)

const (
	MAILBOX_TL_TIME = 600
)

const (
	ELS_Init         = 0 // 默认
	ELS_LoadComplete = 1 // 加载完成
	ELS_SendClient   = 2 // 发送完成
)

// 数据处理接口
type IPlayerFun interface {
	GetPlayerDataType() pb.PlayerDataType // 获取类型
	RegisterTimer()                       // 注册定时器
	Init(uint64, interface{})             // 初始化
	Load([]byte)                          // 加载数据(非system数据)
	IsSave() bool                         // 判断是否存储数据
	Save(bNow bool)                       // 存储数据(非system数据)
	LoadPlayerDBFinish()                  // db数据加载完成回调，在LoadComplete之前调用
	NewPlayer()                           // 初始化新开启的模块数据，在LoadComplete之后调用
	LoadComplete()                        // 加载完成，在NewPlayer之后调用
	UpdateSave(bSave bool)                // 设置保存状态
	Heat()                               // 心跳包
	PassDay(isDay, isWeek, isMonth bool)  // 是否跨天
	SetUserTypeInfo([]byte) error         // 设置缓存数据
	CopyTo(pbData *pb.PBPlayerData)       // 深度拷贝数据
}

// system特殊数据接口
type IPlayerSystemFun interface {
	IPlayerFun
	LoadSystem(pbSystem *pb.PBPlayerSystem)      // 加载系统数据
	SaveSystem(pbSystem *pb.PBPlayerSystem) bool // 存储数据 返回存储标志
}

// 发奖接口
type IReward interface {
	IPlayerFun
	GetSystemType() cfgEnum.ESystemType
	GetPbItem(itemId uint32, itemCount int64, emDoing pb.EmDoingType, params ...uint32) (arrPbItems []*pb.PBAddItemData)
	SendReward(*pb.RpcHead, pb.EmDoingType, interface{}, ...*pb.PBAddItemData) ([]*pb.PBAddItemData, error)
}
