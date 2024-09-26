package entry

import (
	"corps/base/cfgEnum"
	"corps/common/cfgData"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/module/achieve"
	"corps/server/game/module/entry/domain"
	"corps/server/game/module/entry/internal/manager"
	"encoding/json"
)

type EntryService struct {
	uid           uint64                          // 玩家uid
	conditions    map[uint32]*pb.EntryCondition   // 条件数据
	entrys        map[uint32][]*pb.EntryCondition // 词条类型--conds
	effects       map[uint32][]*pb.EntryCondition // 效果类型--conds
	condChanges   map[uint32]struct{}             // 条件变更
	effectChanges map[uint32]struct{}             // 计算效果值
	pAchieveBase  *achieve.AchieveBase            // 成就数据
	hooks         func(effectType uint32)         // 通知钩子函数
	mapCrystal    map[uint32]*pb.PBCrystal        // 晶核系统
}

func NewEntryService(uid uint64, ss *pb.PBPlayerCrystal, pAchieveBase *achieve.AchieveBase, hooks func(effectType uint32)) *EntryService {
	ret := &EntryService{
		uid:           uid,
		conditions:    make(map[uint32]*pb.EntryCondition),
		entrys:        make(map[uint32][]*pb.EntryCondition),
		effects:       make(map[uint32][]*pb.EntryCondition),
		condChanges:   make(map[uint32]struct{}),
		effectChanges: make(map[uint32]struct{}),
		pAchieveBase:  pAchieveBase,
		hooks:         hooks,
	}
	for _, cond := range ss.Conditions {
		// 判断配置是否还存在
		cfg := cfgData.GetCfgEntry(cond.CfgID)
		if cfg == nil {
			plog.Warn("EntryConfigCfg(%d) not found", cond.CfgID)
			continue
		}
		// 加载条件
		ret.conditions[cond.CfgID] = cond

		// 词条类型分类
		if _, ok := ret.entrys[cfg.Type]; !ok {
			ret.entrys[cfg.Type] = make([]*pb.EntryCondition, 0)
		}
		ret.entrys[cfg.Type] = append(ret.entrys[cfg.Type], cond)

		// 效果类型分类
		if _, ok := ret.effects[cfg.EffectType]; !ok {
			ret.effects[cfg.EffectType] = make([]*pb.EntryCondition, 0)
		}
		ret.effects[cfg.EffectType] = append(ret.effects[cfg.EffectType], cond)
	}
	return ret
}

func (d *EntryService) SetCrystal(data map[uint32]*pb.PBCrystal) {
	d.mapCrystal = data
}

func (d *EntryService) Print() string {
	data := &pb.PBPlayerCrystal{}
	d.ToProto(true, data)
	buf, _ := json.Marshal(data)
	return string(buf)
}

func (d *EntryService) ToProto(isCalculate bool, data *pb.PBPlayerCrystal) {
	conds := []*pb.EntryCondition{}
	for _, item := range d.conditions {
		conds = append(conds, item)
		data.Conditions = append(data.Conditions, item)
	}
	if isCalculate {
		// 计算效果值
		for _, val := range d.GetIEntity(conds...) {
			data.Effects = append(data.Effects, val.ToProto())
		}
	}
}

func (d *EntryService) GetIEntity(conds ...*pb.EntryCondition) map[uint32]domain.IEntity {
	results := make(map[uint32]domain.IEntity)
	tmps := map[uint32]domain.IEntity{}
	for _, cond := range conds {
		// 判断是否完成条件, 判断配置是否还存在
		cfg := cfgData.GetCfgEntry(cond.CfgID)
		if cond.Times <= 0 || cfg == nil {
			continue
		}
		// 不同词条处理
		switch cfg.EntryClass {
		case uint32(cfgEnum.EEntryClass_Main):
			// 判断晶核是否解锁
			crystal, ok := d.mapCrystal[cfgData.GetCfgEntryToCrystalID(cfg.PassiveSkillID)]
			if !ok {
				continue
			}
			cfgCrystal := cfgData.GetCfgCrystal(crystal.CrystalID)
			if cfgCrystal == nil {
				continue
			}
			// 创建IEntity
			item := manager.NewEntity(&pb.EntryEffect{ParamsType: cfg.EffectParamType, Type: cfg.EffectType})
			tmps[crystal.CrystalID] = item
			for _, workTag := range cfg.WorkTag {
				for _, params := range cfg.EffectUintParam {
					item.Add(workTag, 1, params...)
				}
			}

			// 升级属性增加
			effectParam := crystal.Level * cfgCrystal.EffectParamRate
			item.AddAll(effectParam)

			// 百分比增加
			percent := uint32(0)
			for i := uint32(1); i <= crystal.Star; i++ {
				if qualityCfg := cfgData.GetCfgCrystalQuality(crystal.Quality, i); qualityCfg != nil {
					percent += qualityCfg.RedefinePercent
				}
			}
			item.PercentAll(percent)
		default:
			// 获取IEntity
			val, ok := results[cfg.EffectType]
			if !ok {
				val = manager.NewEntity(&pb.EntryEffect{ParamsType: cfg.EffectParamType, Type: cfg.EffectType})
				results[cfg.EffectType] = val
			}
			// 计算累加值
			for _, workTag := range cfg.WorkTag {
				for _, params := range cfg.EffectUintParam {
					val.Add(workTag, cond.Times, params...)
				}
			}
		}
	}
	for _, item := range tmps {
		effectType := item.GetType()
		if ent, ok := results[effectType]; !ok {
			results[effectType] = item
		} else {
			for _, params := range item.GetWorkTags() {
				ent.Add(params[0], 1, params[1:]...)
			}
		}
	}
	return results
}

// 获取词条效果参数
func (d *EntryService) Get(effectType uint32, workTag uint32) (rets []*pb.EntryEffectValue) {
	if conds, ok := d.effects[effectType]; ok {
		ent := d.GetIEntity(conds...)
		if val, ok := ent[effectType]; ok {
			rets = val.Get(workTag)
		}
	}
	return
}

func (d *EntryService) add(cfg *cfgData.EntryCfg) *pb.EntryCondition {
	cond := &pb.EntryCondition{CfgID: cfg.Id}
	// 词条ID分类
	d.conditions[cfg.Id] = cond

	// 词条类型分类
	if _, ok := d.entrys[cfg.Type]; !ok {
		d.entrys[cfg.Type] = make([]*pb.EntryCondition, 0)
	}
	d.entrys[cfg.Type] = append(d.entrys[cfg.Type], cond)

	// 效果类型分类
	if _, ok := d.effects[cfg.EffectType]; !ok {
		d.effects[cfg.EffectType] = make([]*pb.EntryCondition, 0)
	}
	d.effects[cfg.EffectType] = append(d.effects[cfg.EffectType], cond)
	return cond
}

// 解锁词条
func (d *EntryService) Unlock(head *pb.RpcHead, passiveSkillID uint32) {
	plog.Trace("head: %v, skillID: %d", head, passiveSkillID)
	for _, cfg := range cfgData.GetCfgEntryByPassiveSkillID(passiveSkillID) {
		// 判断词条配置是否变更
		if _, ok := d.conditions[cfg.Id]; ok {
			plog.Debug("----Unlock is failed----- errorCode: %d, skillID: %d", cfgEnum.ErrorCode_NoData, passiveSkillID)
			continue
		}
		// 解锁条件
		cond := d.add(cfg)
		// 解锁生效类型
		if cfg.Type == uint32(cfgEnum.AchieveType_UnlockEffect) || cfg.IsTotal > 0 {
			//是否取成就数据
			times := uint32(1)
			if cfg.IsTotal > 0 {
				times = d.pAchieveBase.GetAchieveValue(cfg.Type, cfg.SubType...)
			}
			// 触发条件执行
			if cc := manager.GetCondition(cfg.CondParamType); cc != nil {
				cc.Update(cond, cfg, times, cfg.SubType...)
			}
			// 重新计算指定效果
			d.effectChanges[cfg.Id] = struct{}{}
		}

		// 变更记录
		d.condChanges[cfg.Id] = struct{}{}
	}
	d.notifyToClient(head)
}

// 触发词条
func (d *EntryService) Trigger(head *pb.RpcHead, entryType uint32, times uint32, subTypes ...uint32) {
	plog.Trace("head: %v, entryType: %d, times: %d, subTypes: %v", head, entryType, times, subTypes)
	// 判断词条类型是否结果
	conds, ok := d.entrys[entryType]
	if !ok {
		return
	}
	for _, cond := range conds {
		cfg := cfgData.GetCfgEntry(cond.CfgID)
		// 词条被删除
		if cfg == nil {
			continue
		}
		//修正当前值 是否取成就
		if cfg.IsTotal > 0 {
			times = d.pAchieveBase.GetAchieveValue(cfg.Type, cfg.SubType...)
		}
		// 获取条件
		if cen := manager.GetCondition(cfg.CondParamType); cen != nil {
			// 执行条件
			if cen.Update(cond, cfg, times, subTypes...) > 0 {
				// 重新计算指定效果
				d.effectChanges[cfg.Id] = struct{}{}
			}
			// 变更记录
			d.condChanges[cfg.Id] = struct{}{}
		}
	}
	d.notifyToClient(head)
}

// 主词条通知
func (d *EntryService) NotifyMainEntry(head *pb.RpcHead, skillID uint32) {
	if cfg := cfgData.GetCfgMainEntryByPassiveSkillID(skillID); cfg != nil {
		for effecType, ent := range d.GetIEntity(d.effects[cfg.EffectType]...) {
			data := &pb.EntryEffectNotify{PacketHead: &pb.IPacket{}, Effect: ent.ToProto()}
			cluster.SendToClient(head, data, cfgEnum.ErrorCode_Success)

			plog.Debug("head: %v, MainEntryEffect: %v", head, ent.ToProto())
			// 通知
			d.hooks(effecType)
		}
	}
}

// 通知客户端变更
func (d *EntryService) notifyToClient(head *pb.RpcHead) {
	if head == nil {
		head = &pb.RpcHead{Id: d.uid}
	}

	for cfgID := range d.condChanges {
		// 更行条件
		data := &pb.EntryConditionNotify{PacketHead: &pb.IPacket{}, Condition: d.conditions[cfgID]}
		cluster.SendToClient(head, data, cfgEnum.ErrorCode_Success)

		plog.Debug("head: %v, EntryCondition: %v", head, d.conditions[cfgID])
		// 删除记录
		delete(d.condChanges, cfgID)
	}
	// 需要计算效果类型的词条
	filter := map[uint32]struct{}{}
	for cfgID := range d.effectChanges {
		// 判断配置是否还存在
		cfg := cfgData.GetCfgEntry(cfgID)
		if cfg == nil {
			plog.Warn("EntryConfigCfg(%d) not found", cfgID)
			continue
		}
		filter[cfg.EffectType] = struct{}{}
		delete(d.effectChanges, cfgID)
	}
	conds := []*pb.EntryCondition{}
	for effectType := range filter {
		conds = append(conds, d.effects[effectType]...)
	}
	for effecType, ent := range d.GetIEntity(conds...) {
		data := &pb.EntryEffectNotify{PacketHead: &pb.IPacket{}, Effect: ent.ToProto()}
		cluster.SendToClient(head, data, cfgEnum.ErrorCode_Success)
		plog.Debug("head: %v, EntryEffect: %v", head, ent.ToProto())
		// 通知
		d.hooks(effecType)
	}
}
