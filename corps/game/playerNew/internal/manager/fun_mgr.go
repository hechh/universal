package manager

import (
	"corps/base/cfgEnum"
	"corps/pb"
	"corps/server/game/playerNew/domain"
)

type FunMgr struct {
	loots    map[cfgEnum.ESystemType]domain.RewardFunc // 掉落接口
	rewards  map[cfgEnum.ESystemType]domain.RewardFunc // 发奖接口
	funs     map[pb.PlayerDataType]domain.IPlayerFun   // 业务接口
	priority []domain.IPlayerFun                       // 玩家数据初始化顺序
}

func NewFunMgr() *FunMgr {
	return &FunMgr{
		loots:   make(map[cfgEnum.ESystemType]domain.RewardFunc),
		rewards: make(map[cfgEnum.ESystemType]domain.RewardFunc),
		funs:    make(map[pb.PlayerDataType]domain.IPlayerFun),
	}
}

func (d *FunMgr) RegisterIPlayerFun(funs ...domain.IPlayerFun) {
	for _, fun := range funs {
		typ := fun.GetPlayerDataType()
		d.funs[typ] = fun
		d.priority = append(d.priority, fun)
	}
}

func (d *FunMgr) RegisterReward(typ cfgEnum.ESystemType, rs domain.RewardFunc) {
	d.rewards[typ] = rs
}

func (d *FunMgr) GetReward(typ cfgEnum.ESystemType) domain.RewardFunc {
	return d.rewards[typ]
}

func (d *FunMgr) RegisterLoot(typ cfgEnum.ESystemType, rs domain.RewardFunc) {
	d.loots[typ] = rs
}

func (d *FunMgr) GetLoot(typ cfgEnum.ESystemType) domain.RewardFunc {
	return d.loots[typ]
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

// 发送奖励
func (d *FunMgr) SendReward(head *pb.RpcHead, do pb.EmDoingType, extra interface{}, items ...*pb.PBAddItemData) (rets []*pb.PBAddItemData, err error) {

	return
}
