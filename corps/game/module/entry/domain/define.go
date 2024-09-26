package domain

import (
	"corps/common/cfgData"
	"corps/pb"
)

// 效果数据
type IEntity interface {
	GetParamType() uint32                                // 获取参数效果类型
	GetType() uint32                                     // 获取效果类型
	ToProto() *pb.EntryEffect                            // 转成pb格式存储
	Get(efObject uint32) []*pb.EntryEffectValue          // 获取所有效果
	Add(efObejct uint32, times uint32, params ...uint32) // 累加属性
	AddAll(val uint32)                                   // 全属性增加
	PercentAll(val uint32)                               // 百分比提升
	GetWorkTags() [][]uint32                             // 获取标签效果值
}

// 词条
type ICondition interface {
	Update(pbData *pb.EntryCondition, cfg *cfgData.EntryCfg, times uint32, subTypes ...uint32) uint32
}
