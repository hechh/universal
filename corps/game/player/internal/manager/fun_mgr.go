package manager

import (
	"corps/base/cfgEnum"
	"corps/pb"
	"corps/server/game/player/domain"
)

type FunMgr struct {
	rewards  map[cfgEnum.ESystemType]domain.IReward  // 发奖接口
	funs     map[pb.PlayerDataType]domain.IPlayerFun // 业务接口
	priority []domain.IPlayerFun                     // 玩家数据初始化顺序
}

func NewFunMgr() *FunMgr {
	return &FunMgr{
		rewards: make(map[cfgEnum.ESystemType]domain.IReward),
		funs:    make(map[pb.PlayerDataType]domain.IPlayerFun),
	}
}

func (d *FunMgr) RegisterIPlayerFun(funs ...domain.IPlayerFun) {
	for _, fun := range funs {
		typ := fun.GetPlayerDataType()
		d.funs[typ] = fun
		d.priority = append(d.priority, fun)

		switch vv := fun.(type) {
		case domain.IReward:
			d.rewards[vv.GetSystemType()] = vv
		default:
		}
	}
}

func (d *FunMgr) RegisterIReward(rs ...domain.IReward) {
	for _, r := range rs {
		d.rewards[r.GetSystemType()] = r
	}
}

func (d *FunMgr) GetIReward(typ cfgEnum.ESystemType) domain.IReward {
	return d.rewards[typ]
}

func (d *FunMgr) GetIPlayerFun(typ pb.PlayerDataType) domain.IPlayerFun {
	return d.funs[typ]
}

func (d *FunMgr) GetIPlayerFunList() []domain.IPlayerFun {
	return d.priority
}

func (d *FunMgr) Walk(f func(pb.PlayerDataType, domain.IPlayerFun) bool) {
	for typ, val := range d.funs {
		if !f(typ, val) {
			break
		}
	}
}
