package entity

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/pb"
	"corps/server/game/module/entry/domain"
	"math"
)

type KeyValue struct {
	workTag uint32            // 效果类型
	data    map[uint32]uint32 // 属性-属性值
}

type TwoArrayEntity struct {
	paramType  uint32               // 参数类型
	effectType uint32               // 效果类型
	data       map[uint32]*KeyValue // 生效对象--属性集合
}

func NewTwoArrayEntity(vv *pb.EntryEffect) domain.IEntity {
	ret := make(map[uint32]*KeyValue)
	for _, item := range vv.List {
		if _, ok := ret[item.Object]; !ok {
			ret[item.Object] = &KeyValue{
				workTag: item.Object,
				data:    make(map[uint32]uint32),
			}
		}

		for _, val := range item.Values {
			ret[item.Object].data[val.List[0]] = val.List[1]
		}
	}
	return &TwoArrayEntity{data: ret, effectType: vv.Type, paramType: vv.ParamsType}
}

func (d *TwoArrayEntity) PercentAll(per uint32) {
	for _, work := range d.data {
		for key, val := range work.data {
			work.data[key] += uint32(math.Floor(float64(val*per) / base.MIL_PERCENT))
		}
	}
}

func (d *TwoArrayEntity) AddAll(param uint32) {
	for _, work := range d.data {
		for key := range work.data {
			work.data[key] += param
		}
	}
}

func (d *TwoArrayEntity) GetType() uint32 {
	return d.effectType
}

func (d *TwoArrayEntity) GetParamType() uint32 {
	return d.paramType
}

func (d *TwoArrayEntity) ToProto() *pb.EntryEffect {
	elem := &pb.EntryEffect{Type: d.effectType, ParamsType: d.paramType}
	for _, vals := range d.data {
		ret := &pb.EntryEffectData{Object: vals.workTag}
		for key, val := range vals.data {
			ret.Values = append(ret.Values, &pb.EntryEffectValue{List: []uint32{key, val}})
		}
		elem.List = append(elem.List, ret)
	}
	return elem
}

func (d *TwoArrayEntity) GetWorkTags() (rets [][]uint32) {
	for workTag, vals := range d.data {
		for shuxin, shuxinVal := range vals.data {
			rets = append(rets, []uint32{workTag, shuxin, shuxinVal})
		}
	}
	return
}

func (d *TwoArrayEntity) Add(workTag, times uint32, params ...uint32) {
	if len(params) < 2 {
		return
	}
	// 判断类型是否解锁
	val, ok := d.data[workTag]
	if !ok {
		val = &KeyValue{
			workTag: workTag,
			data:    make(map[uint32]uint32),
		}
		d.data[workTag] = val
	}
	// 累加属性
	val.data[params[0]] += (params[1] * times)
}

// 获取效果
func (d *TwoArrayEntity) Get(workTag uint32) (rets []*pb.EntryEffectValue) {
	nnil := uint32(cfgEnum.EntryWorkTag_None)
	// 从全局读取
	tmps := map[uint32]uint32{}
	if vals, ok := d.data[nnil]; ok {
		for key, val := range vals.data {
			tmps[key] += val
		}
	}
	// 从workTag中读取
	if vals, ok := d.data[workTag]; nnil != workTag && ok {
		for key, val := range vals.data {
			tmps[key] += val
		}
	}
	// 返回
	for key, value := range tmps {
		rets = append(rets, &pb.EntryEffectValue{List: []uint32{key, value}})
	}
	return
}
