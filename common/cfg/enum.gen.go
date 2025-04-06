package cfg

type PropertyType uint32

const (
	PropertyType_Empty   PropertyType = 0 // 道具类型-空
	PropertyType_Coin    PropertyType = 1 // 道具类型-金币
	PropertyType_Diamond PropertyType = 2 // 道具类型-钻石
	PropertyType_Sword   PropertyType = 3 // 装备类型-剑
)

type QualityType uint32

const (
	QualityType_None  QualityType = 0 // 品质类型-空
	QualityType_White QualityType = 1 // 品质类型-白
	QualityType_Blue  QualityType = 2 // 品质类型-蓝
)
