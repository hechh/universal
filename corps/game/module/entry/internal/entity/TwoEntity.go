package entity

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/pb"
	"corps/server/game/module/entry/domain"
	"math"
)

type PairInfo struct {
	workTag uint32
	key     uint32
	value   uint32
}

type TwoEntity struct {
	paramType  uint32               // 参数类型
	effectType uint32               // 效果类型
	data       map[uint32]*PairInfo // 生效对象--属性集合
}

func NewTwoEntity(vv *pb.EntryEffect) domain.IEntity {
	ret := make(map[uint32]*PairInfo)
	for _, item := range vv.List {
		ret[item.Object] = &PairInfo{item.Object, item.Values[0].List[0], item.Values[0].List[1]}
	}
	return &TwoEntity{data: ret, effectType: vv.Type, paramType: vv.ParamsType}
}

func (d *TwoEntity) PercentAll(per uint32) {
	for _, work := range d.data {
		work.value += uint32(math.Floor(float64(work.value*per) / base.MIL_PERCENT))
	}
}

func (d *TwoEntity) AddAll(param uint32) {
	for _, work := range d.data {
		work.value += param
	}
}

func (d *TwoEntity) GetType() uint32 {
	return d.effectType
}

func (d *TwoEntity) GetParamType() uint32 {
	return d.paramType
}

func (d *TwoEntity) ToProto() *pb.EntryEffect {
	elem := &pb.EntryEffect{Type: d.effectType, ParamsType: d.paramType}
	for _, vals := range d.data {
		elem.List = append(elem.List, &pb.EntryEffectData{
			Object: vals.workTag,
			Values: []*pb.EntryEffectValue{
				{List: []uint32{vals.key, vals.value}},
			},
		})
	}
	return elem
}

func (d *TwoEntity) GetWorkTags() (rets [][]uint32) {
	for workTag, val := range d.data {
		rets = append(rets, []uint32{workTag, val.key, val.value})
	}
	return
}

func (d *TwoEntity) Add(workTag, times uint32, params ...uint32) {
	// 判断类型是否解锁
	val, ok := d.data[workTag]
	if !ok {
		val = &PairInfo{workTag: workTag}
		d.data[workTag] = val
	}
	// 累加属性
	val.key = params[0]
	val.value += (params[1] * times)
}

// 获取效果
func (d *TwoEntity) Get(workTag uint32) (rets []*pb.EntryEffectValue) {
	nnil := uint32(cfgEnum.EntryWorkTag_None)
	// 从全局读取
	tmp := &PairInfo{}
	if global, ok := d.data[nnil]; ok {
		tmp.key = global.key
		tmp.value += global.value
	}
	// 从workTag中读取
	if val, ok := d.data[workTag]; nnil != workTag && ok {
		tmp.key = val.key
		tmp.value += val.value
	}
	// 返回
	if tmp.value > 0 {
		rets = append(rets, &pb.EntryEffectValue{List: []uint32{tmp.key, tmp.value}})
	}
	return
}
