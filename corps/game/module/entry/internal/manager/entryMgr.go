package manager

import (
	"corps/pb"
	"corps/server/game/module/entry/domain"
	"corps/server/game/module/entry/internal/entity"
)

type EntityFunc func(*pb.EntryEffect) domain.IEntity

var (
	conds   = make(map[uint32]domain.ICondition)
	entitys = make(map[uint32]EntityFunc)
)

func RegisterCondition(f domain.ICondition, typs ...uint32) {
	for _, typ := range typs {
		if _, ok := conds[typ]; ok {
			panic("repeated register")
		}
		conds[typ] = f
	}
}

func RegisterEntity(f EntityFunc, typs ...uint32) {
	for _, typ := range typs {
		if _, ok := entitys[typ]; ok {
			panic("repeated register")
		}
		entitys[typ] = f
	}
}

// 获取条件
func GetCondition(typ uint32) domain.ICondition {
	if val, ok := conds[typ]; ok {
		return val
	}
	return nil
}

func NewEntity(vv *pb.EntryEffect) domain.IEntity {
	if f, ok := entitys[vv.ParamsType]; ok {
		return f(vv)
	}
	return &entity.EmptyEntity{}
}
