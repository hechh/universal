package common

import (
	"corps/pb"
)

const (
	PERCENT = 100
)

type Common struct{}

func (d *Common) Handle(prob int64, items ...*pb.PBAddItemData) (rets []*pb.PBAddItemData) {
	for _, item := range items {
		newItem := *item
		newItem.Params = append(newItem.Params, item.Params...)
		newItem.Count = newItem.Count * prob
		rets = append(rets, &newItem)
	}
	return
}
