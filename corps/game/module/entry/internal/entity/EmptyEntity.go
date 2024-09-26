package entity

import (
	"corps/base/cfgEnum"
	"corps/pb"
	"corps/server/game/module/entry/domain"
)

type EmptyEntity struct{}

func NewEmptyEntity(vv *pb.EntryEffect) domain.IEntity {
	return &EmptyEntity{}
}

func (d *EmptyEntity) ToProto() (ret *pb.EntryEffect)             { return }
func (d *EmptyEntity) Add(_, _ uint32, params ...uint32)          {}
func (d *EmptyEntity) Get(_ uint32) (rets []*pb.EntryEffectValue) { return }
func (d *EmptyEntity) GetType() uint32                            { return uint32(cfgEnum.EntryEffectType_None) }
func (d *EmptyEntity) GetParamType() uint32                       { return uint32(cfgEnum.EEntryParamTypeEffect_None) }
func (d *EmptyEntity) GetWorkTags() [][]uint32                    { return nil }
func (d *EmptyEntity) AddAll(params uint32)                       {}
func (d *EmptyEntity) PercentAll(uint32)                          {}
