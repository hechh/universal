package playerfun

import (
	"poker_server/common/config/repository/php_config"
	"poker_server/common/pb"
	"poker_server/common/room_util"
	"poker_server/framework"
	"poker_server/library/httplib"
	"poker_server/library/mlog"
	"poker_server/library/snowflake"
	"poker_server/library/uerror"
	"poker_server/server/game/internal/player/domain"
	"strconv"

	"github.com/golang/protobuf/proto"
)

type PlayerBagFun struct {
	*PlayerFun
	data *pb.PlayerDataBag
}

func NewPlayerBagFun(fun *PlayerFun) domain.IPlayerFun {
	return &PlayerBagFun{PlayerFun: fun}
}

func (p *PlayerBagFun) GetData() *pb.PlayerDataBag {
	return p.data
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

func (d *PlayerBagFun) GetProp(coinType pb.CoinType) int64 {
	if item, ok := d.data.Items[uint32(coinType)]; ok {
		return item.Count
	}
	return 0
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
	d.Change()
}

func (d *PlayerBagFun) SubProp(coinType pb.CoinType, val int64) error {
	if val == 0 {
		return nil
	}
	if val < 0 {
		return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "负值错误")
	}

	item, ok := d.data.Items[uint32(coinType)]
	if !ok || item.Count < val {
		return uerror.New(1, pb.ErrorCode_GAME_PROP_NOT_ENOUGH, "游戏道具不足")
	}

	item.Count -= val
	d.Change()
	return nil
}

func (d *PlayerBagFun) GetBagReq(head *pb.Head, req *pb.GetBagReq, rsp *pb.GetBagRsp) error {
	for _, item := range d.data.Items {
		rsp.List = append(rsp.List, item)
	}
	return nil
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
		return uerror.New(1, pb.ErrorCode_GAME_PROP_NOT_ENOUGH, "道具数量不足: 背包为空")
	}
	tmps := map[uint32]*pb.Reward{}
	for _, rw := range req.RewardList {
		if rw.Incr < 0 {
			return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "消耗的数量不能小于0")
		}
		val, ok := tmps[rw.PropId]
		if !ok {
			val = rw
			tmps[rw.PropId] = rw
		} else {
			val.Incr += rw.Incr
		}
		if items, ok := d.data.Items[rw.PropId]; !ok || items.Count < val.Incr {
			mlog.Info(head, "Consume items : %v, %v", d.data.Items, val.Incr)
			return uerror.New(1, pb.ErrorCode_GAME_PROP_NOT_ENOUGH, "道具数量不足: %v", items)
		}
	}
	for _, rw := range tmps {
		item := d.data.Items[rw.PropId]
		item.Count -= rw.Incr
		rw.Total = item.Count
		rsp.RewardList = append(rsp.RewardList, rw)
		d.Change()
	}
	return nil
}

func (d *PlayerBagFun) ChargeTransIn(head *pb.Head, roomId uint64) error {
	_, gameType, coinType := room_util.TexasRoomIdTo(roomId)
	item, ok := d.data.Items[uint32(coinType)]
	if !ok {
		item = &pb.PbItem{PropId: uint32(coinType)}
		d.data.Items[uint32(coinType)] = item
	}

	uuid, _ := snowflake.GenUUID()
	params := map[string]interface{}{
		"game_sn":   strconv.FormatUint(uuid, 10),
		"game_type": uint32(gameType),
		"coin_type": uint32(coinType),
		"incr":      0,
		"uid":       head.Uid,
	}

	url := php_config.MGetEnvTypeName(framework.GetEnvType(), "TransferOutUrl").Url
	rsp := &pb.HttpTransferInRsp{}
	if err := httplib.POST(url, params, rsp); err != nil {
		mlog.Infof("ChargeTransIn http err: %v, rsp:%v", err, rsp)
		return err
	}
	if rsp.RespMsg.Code != 100 {
		mlog.Infof("ChargeTransIn http rsp:%v", rsp)
		return uerror.New(1, pb.ErrorCode_REQUEST_FAIELD, "返回错误码%v", rsp)
	}

	item.Count += (-rsp.RespData.Change)
	d.Change()
	mlog.Infof("ChargeTransIn uid:%d, param:%v, resp:%v, item:%v", head.Uid, params, rsp, item)
	return nil
}

func (d *PlayerBagFun) ChargeTransOut(head *pb.Head, roomId uint64) error {
	_, gameType, coinType := room_util.TexasRoomIdTo(roomId)
	item, ok := d.data.Items[uint32(coinType)]
	if !ok || item.Count <= 0 {
		return nil
	}

	uuid, _ := snowflake.GenUUID()
	params := map[string]interface{}{
		"game_sn":   strconv.FormatUint(uuid, 10),
		"game_type": uint32(gameType),
		"coin_type": uint32(coinType),
		"incr":      item.Count,
		"uid":       head.Uid,
	}
	url := php_config.MGetEnvTypeName(framework.GetEnvType(), "TransferInUrl").Url
	rsp := &pb.HttpTransferOutRsp{}

	if err := httplib.POST(url, params, rsp); err != nil {
		mlog.Infof("ChargeTransIn http err: %v, rsp:%v", err, rsp)
		return err
	}
	if rsp.RespMsg.Code != 100 {
		mlog.Infof("ChargeTransIn http rsp:%v", rsp)
		return uerror.New(1, pb.ErrorCode_REQUEST_FAIELD, "返回错误码%v", rsp)
	}

	item.Count = 0
	d.Change()
	mlog.Infof("ChargeTransOut uid:%d, param:%v, rsp:%v, item:%v", head.Uid, params, rsp, item)
	return nil
}
