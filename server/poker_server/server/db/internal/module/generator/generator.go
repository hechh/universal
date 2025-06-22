package generator

import (
	"poker_server/common/dao/repository/redis/generator_data"
	"poker_server/common/pb"
	"poker_server/framework/actor"
	"reflect"
	"time"
)

type DbGeneratorMgr struct {
	actor.Actor
	datas    map[pb.GeneratorType]*pb.GeneratorInfo
	isChange bool
}

func NewGeneratorMgr() *DbGeneratorMgr {
	ret := &DbGeneratorMgr{}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Actor.SetId(uint64(pb.DataType_DataTypeGenerator))
	return ret
}

func (m *DbGeneratorMgr) Init() error {
	if data, _, err := generator_data.Get(); err != nil {
		return err
	} else {
		m.datas = make(map[pb.GeneratorType]*pb.GeneratorInfo)
		for _, item := range data.List {
			m.datas[item.Id] = item
		}
	}
	for i := pb.GeneratorType_GeneratorTypeBegin + 1; i < pb.GeneratorType_GeneratorTypeEnd; i++ {
		if _, ok := m.datas[i]; ok {
			continue
		}
		m.datas[i] = &pb.GeneratorInfo{
			Id:   i,
			Incr: 1,
		}
		m.isChange = true
	}
	m.Start()
	actor.Register(m)
	m.RegisterTimer(&pb.Head{FuncName: "OnTick"}, 5*time.Second, -1)
	return nil
}

func (m *DbGeneratorMgr) Close() {
	m.Save()
	m.Actor.Stop()
}

func (m *DbGeneratorMgr) OnTick() {
	if !m.isChange {
		return
	}
	m.SendMsg(&pb.Head{FuncName: "Save"})
}

// 保存数据
func (m *DbGeneratorMgr) Save() error {
	if !m.isChange {
		return nil
	}
	data := &pb.GeneratorData{}
	for _, item := range m.datas {
		data.List = append(data.List, item)
	}
	if err := generator_data.Set(data); err != nil {
		return err
	}
	m.isChange = false
	return nil
}

// 集中所有变更数据
func (m *DbGeneratorMgr) Update(head *pb.Head, req *pb.UpdateGeneratorDataNotify) error {
	if len(req.List) <= 0 {
		return nil
	}
	for _, item := range req.List {
		m.datas[item.Id] = item
		m.isChange = true
	}
	return nil
}

// 加载数据请求(同步和异步请求支持)
func (m *DbGeneratorMgr) Query(head *pb.Head, req *pb.GetGeneratorDataReq, rsp *pb.GetGeneratorDataRsp) error {
	for _, item := range m.datas {
		rsp.List = append(rsp.List, item)
	}
	return nil
}
