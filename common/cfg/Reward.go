package cfg

type Reward struct {
	PropType PropertyType // 道具类型
	Quality  QualityType  // 品质
	Star     uint32       // 星级
	Add      int64        // 道具数量
}
