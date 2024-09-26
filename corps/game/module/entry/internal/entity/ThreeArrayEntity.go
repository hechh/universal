package entity

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/pb"
	"corps/server/game/module/entry/domain"
	"math"
)

type ItemInfo struct {
	workTag uint32                     // 对象类型
	data    map[common.ItemIndex]int64 // 道具类型--道具数据
}

type ThreeArrayEntity struct {
	paramType  uint32               // 参数类型
	effectType uint32               // 效果类型
	data       map[uint32]*ItemInfo // 生效对象--系统道具
}

func NewThreeArrayEntity(vv *pb.EntryEffect) domain.IEntity {
	ret := make(map[uint32]*ItemInfo)
	for _, item := range vv.List {
		info, ok := ret[item.Object]
		if !ok {
			info = &ItemInfo{workTag: item.Object, data: make(map[common.ItemIndex]int64)}
			ret[item.Object] = info
		}
		for _, val := range item.Values {
			info.data[common.ItemIndex{Kind: val.List[0], Id: val.List[1]}] = int64(val.List[2])
		}
	}
	return &ThreeArrayEntity{data: ret, effectType: vv.Type, paramType: vv.ParamsType}
}

func (d *ThreeArrayEntity) AddAll(param uint32) {
	for _, work := range d.data {
		for key := range work.data {
			work.data[key] += int64(param)
		}
	}
}

func (d *ThreeArrayEntity) PercentAll(per uint32) {
	for _, work := range d.data {
		for key, val := range work.data {
			work.data[key] += int64(math.Floor(float64(int64(per)*val) / base.MIL_PERCENT))
		}
	}
}

func (d *ThreeArrayEntity) GetType() uint32 {
	return d.effectType
}

func (d *ThreeArrayEntity) GetParamType() uint32 {
	return d.paramType
}

func (d *ThreeArrayEntity) ToProto() *pb.EntryEffect {
	elem := &pb.EntryEffect{Type: d.effectType, ParamsType: d.paramType}
	for _, vals := range d.data {
		ret := &pb.EntryEffectData{Object: vals.workTag}
		for index, val := range vals.data {
			ret.Values = append(ret.Values, &pb.EntryEffectValue{List: []uint32{index.Kind, index.Id, uint32(val)}})
		}
		elem.List = append(elem.List, ret)
	}
	return elem
}

func (d *ThreeArrayEntity) GetWorkTags() (rets [][]uint32) {
	for workTag, vals := range d.data {
		for daoju, daojuVal := range vals.data {
			rets = append(rets, []uint32{workTag, daoju.Kind, daoju.Id, uint32(daojuVal)})
		}
	}
	return
}

func (d *ThreeArrayEntity) Add(workTag, times uint32, params ...uint32) {
	if len(params) < 3 {
		return
	}
	val, ok := d.data[workTag]
	if !ok {
		val = &ItemInfo{workTag: workTag, data: make(map[common.ItemIndex]int64)}
		d.data[workTag] = val
	}
	val.data[common.ItemIndex{Kind: params[0], Id: params[1]}] += int64(params[2] * times)
}

// 获取效果
func (d *ThreeArrayEntity) Get(workTag uint32) (rets []*pb.EntryEffectValue) {
	nnil := uint32(cfgEnum.EntryWorkTag_None)
	// 从workTag中读取
	tmps := map[common.ItemIndex]int64{}
	if val, ok := d.data[workTag]; ok {
		for index, item := range val.data {
			tmps[index] = item
		}
	}
	// 从全局中读取
	if vals, ok := d.data[nnil]; workTag != nnil && ok {
		for index, item := range vals.data {
			tmps[index] += item
		}
	}
	// 返回
	for index, item := range tmps {
		rets = append(rets, &pb.EntryEffectValue{List: []uint32{index.Kind, index.Id, uint32(item)}})
	}
	return
}
