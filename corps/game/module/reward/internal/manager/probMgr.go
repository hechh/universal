package manager

import (
	"corps/pb"
	"corps/server/game/module/reward/domain"
	"fmt"
)

var (
	probMgr = make(map[uint32]domain.IReward)
)

func Register(dd domain.IReward, typs ...uint32) {
	for _, typ := range typs {
		if _, ok := probMgr[typ]; ok {
			panic(fmt.Sprintf("IProbability has already registered, type: %d", typ))
		}
		probMgr[typ] = dd
	}
}

func GetIReward(typ uint32) domain.IReward {
	return probMgr[typ]
}

func GetClass(items ...*pb.PBAddItemData) (classType map[uint32][]*pb.PBAddItemData, rets []*pb.PBAddItemData) {
	classType = map[uint32][]*pb.PBAddItemData{}
	for _, item := range items {
		// 判断是否支持
		if _, ok := probMgr[item.Kind]; !ok {
			rets = append(rets, item)
			continue
		}
		// 支持类型
		if _, ok := classType[item.Kind]; !ok {
			classType[item.Kind] = []*pb.PBAddItemData{}
		}
		classType[item.Kind] = append(classType[item.Kind], item)
	}
	return
}
