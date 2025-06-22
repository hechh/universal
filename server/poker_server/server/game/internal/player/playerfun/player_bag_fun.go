package playerfun

import (
	"poker_server/common/pb"
	"poker_server/library/uerror"
	"poker_server/server/game/internal/player/domain"

	"github.com/golang/protobuf/proto"
)

type PlayerBagFun struct {
	*PlayerFun
	data *pb.PlayerDataBag
}

func NewPlayerBagFun(fun *PlayerFun) domain.IPlayerFun {
	return &PlayerBagFun{PlayerFun: fun}
}

func (d *PlayerBagFun) Load(msg *pb.PlayerData) error {
	if msg == nil {
		return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "玩家数据为空")
	}
	if msg.Bag == nil {
		msg.Bag = &pb.PlayerDataBag{Items: make(map[uint32]*pb.PbItem)}
	}
	d.data = msg.Bag
	return nil
}

func (d *PlayerBagFun) Save(data *pb.PlayerData) error {
	if data == nil {
		return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "玩家数据为空")
	}
	buf, _ := proto.Marshal(d.data)
	newBase := &pb.PlayerDataBag{}
	proto.Unmarshal(buf, newBase)
	data.Bag = newBase
	return nil
}

func (d *PlayerBagFun) LoadComplate() error {
	return nil
}

func (d *PlayerBagFun) AddProp(coinType pb.CoinType, val int64) {
	if item, ok := d.data.Items[uint32(coinType)]; ok {
		item.Count += val
	} else {
		d.data.Items[uint32(coinType)] = &pb.PbItem{
			PropId: uint32(coinType),
			Count:  val,
		}
	}
}

func (d *PlayerBagFun) RewardReq(head *pb.Head, req *pb.RewardReq, rsp *pb.RewardRsp) error {
	if d.data.Items == nil {
		d.data.Items = make(map[uint32]*pb.PbItem)
	}
	tmps := map[uint32]*pb.Reward{}
	for _, rw := range req.RewardList {
		if rw.Incr < 0 {
			return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "奖励增加数量不能小于0")
		}
		if val, ok := tmps[rw.PropId]; ok {
			val.Incr += rw.Incr
		} else {
			tmps[rw.PropId] = rw
		}
	}
	for _, rw := range tmps {
		val, ok := d.data.Items[rw.PropId]
		if !ok {
			val = &pb.PbItem{PropId: rw.PropId, Count: 0}
			d.data.Items[rw.PropId] = val
		}
		val.Count += rw.Incr
		rw.Total = val.Count
		rsp.RewardList = append(rsp.RewardList, rw)
		d.Change()
	}
	return nil
}

func (d *PlayerBagFun) ConsumeReq(head *pb.Head, req *pb.ConsumeReq, rsp *pb.ConsumeRsp) error {
	if d.data.Items == nil {
		d.data.Items = make(map[uint32]*pb.PbItem)
	}
	tmps := map[uint32]*pb.Reward{}
	for _, rw := range req.RewardList {
		if rw.Incr < 0 {
			return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "消耗的数量不能小于0")
		}
		val, ok := tmps[rw.PropId]
		if !ok {
			val = rw
			tmps[rw.PropId] = val
		}
		val.Incr += rw.Incr
		if items, ok := d.data.Items[rw.PropId]; !ok || items.Count < val.Incr {
			return uerror.New(1, pb.ErrorCode_GAME_PROP_NOT_ENOUGH, "道具数量不足: %v", items)
		}
	}
	for _, rw := range tmps {
		val := d.data.Items[rw.PropId]
		val.Count -= rw.Incr
		rw.Total = val.Count
		rsp.RewardList = append(rsp.RewardList, rw)
		d.Change()
	}
	return nil
}
