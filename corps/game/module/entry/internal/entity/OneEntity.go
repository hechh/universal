package entity

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/pb"
	"corps/server/game/module/entry/domain"
	"math"
)

type OneEntity struct {
	paramType  uint32            // 参数类型
	effectType uint32            // 效果类型
	data       map[uint32]uint32 // 生效对象---数值
}

func NewOneEntity(vv *pb.EntryEffect) domain.IEntity {
	ret := make(map[uint32]uint32)
	for _, item := range vv.List {
		ret[item.Object] = item.Values[0].List[0]
	}
	return &OneEntity{data: ret, effectType: vv.Type, paramType: vv.ParamsType}
}

func (d *OneEntity) GetType() uint32 {
	return d.effectType
}

func (d *OneEntity) GetParamType() uint32 {
	return d.paramType
}

func (d *OneEntity) ToProto() *pb.EntryEffect {
	elem := &pb.EntryEffect{Type: d.effectType, ParamsType: d.paramType}
	for workTag, vals := range d.data {
		elem.List = append(elem.List, &pb.EntryEffectData{
			Object: workTag,
			Values: []*pb.EntryEffectValue{
				{List: []uint32{vals}},
			},
		})
	}
	return elem
}

func (d *OneEntity) GetWorkTags() (rets [][]uint32) {
	for workTag, val := range d.data {
		rets = append(rets, []uint32{workTag, val})
	}
	return
}

func (d *OneEntity) AddAll(param uint32) {
	for key := range d.data {
		d.data[key] += param
	}
}

// 百分比提升
func (d *OneEntity) PercentAll(per uint32) {
	for key, val := range d.data {
		d.data[key] += uint32(math.Floor(float64(val * per / base.MIL_PERCENT)))
	}
}

// 累加
func (d *OneEntity) Add(effectObejct, times uint32, params ...uint32) {
	d.data[effectObejct] += (params[0] * times)
}

func (d *OneEntity) Get(workTag uint32) (rets []*pb.EntryEffectValue) {
	nnil := uint32(cfgEnum.EntryWorkTag_None)
	// 从全局读取
	total := uint32(0)
	if global, ok := d.data[nnil]; ok {
		total += global
	}
	// 从workTag中读取
	if val, ok := d.data[workTag]; workTag != nnil && ok {
		total += val
	}
	// 返回
	if total > 0 {
		rets = append(rets, &pb.EntryEffectValue{List: []uint32{total}})
	}
	return
}
