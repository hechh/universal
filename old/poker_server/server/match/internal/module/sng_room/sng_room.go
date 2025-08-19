package sng_room

import (
	"poker_server/common/config/repository/sng_match_config"
	"poker_server/common/config/repository/sng_match_rank_reward_config"
	"poker_server/common/config/repository/texas_config"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"reflect"
	"sort"
	"time"
)

type SngRoomMgr struct {
	actor.Actor
	datas map[int32]*pb.TexasRoomData
}

func NewSngRoomMgr() *SngRoomMgr {
	ret := &SngRoomMgr{datas: make(map[int32]*pb.TexasRoomData)}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Actor.SetId(uint64(pb.DataType_DataTypeSngRoom))
	ret.Actor.Start()
	actor.Register(ret)
	return ret
}

func (d *SngRoomMgr) Close() {
	d.Actor.Stop()
	mlog.Infof("SngRoomMgr关闭成功")
}

func (d *SngRoomMgr) HasRoomReq(head *pb.Head, req *pb.HasRoomReq, rsp *pb.HasRoomRsp) error {
	for _, rr := range d.datas {
		if rr.RoomId == req.RoomId {
			rsp.IsExist = true
			break
		}
	}
	return nil
}

func (d *SngRoomMgr) Query(head *pb.Head, roomId uint64) error {
	if tmps := d.check(); len(tmps) > 0 {
		if err := d.build(tmps...); err != nil {
			return err
		}
	}
	var data *pb.TexasRoomData
	for _, roomData := range d.datas {
		if roomData.RoomId == roomId {
			data = roomData
			break
		}
	}
	if data == nil {
		return uerror.New(1, pb.ErrorCode_TEXAS_ROOM_NOT_FOUND, "SNG房间人数已满，无法在加入 head:%v, roomId:%d", roomId)
	}
	return nil
}

func (d *SngRoomMgr) Update(head *pb.Head, roomData *pb.TexasRoomData) {
	for key, rData := range d.datas {
		if rData.RoomId != roomData.RoomId {
			continue
		}

		cfg := texas_config.MGetID(rData.GameId)
		d.datas[key] = roomData
		if cfg.MinPlayerCount <= int32(roomData.OnlineNumber) {
			delete(d.datas, key)
		}
		break
	}
}

func (d *SngRoomMgr) SngJoinRoomReq(head *pb.Head, req *pb.SngJoinRoomReq, rsp *pb.SngJoinRoomRsp) error {
	if tmps := d.check(); len(tmps) > 0 {
		if err := d.build(tmps...); err != nil {
			return err
		}
	}

	cfg := texas_config.MGetID(req.MatchId)
	if cfg == nil {
		return uerror.New(1, pb.ErrorCode_CONFIG_NOT_FOUND, "req:%v", req)
	}

	sngcfg := sng_match_config.MGetGameId(req.MatchId)
	if sngcfg == nil {
		return uerror.New(1, pb.ErrorCode_CONFIG_NOT_FOUND, "req:%v", req)
	}

	data := d.datas[req.MatchId]
	rsp.RoomId = data.RoomId
	rsp.Chip = sngcfg.Chip
	return nil
}

func (d *SngRoomMgr) RoomListReq(head *pb.Head, req *pb.SngRoomListReq, rsp *pb.SngRoomListRsp) error {
	if tmps := d.check(); len(tmps) > 0 {
		if err := d.build(tmps...); err != nil {
			return err
		}
	}
	for _, item := range d.datas {
		cfg := texas_config.MGetID(item.GameId)
		sngcfg := sng_match_config.MGetGameId(item.GameId)
		total := []*pb.CoinReward{}
		tmps := map[pb.CoinType]*pb.CoinReward{}
		for _, cfg := range sng_match_rank_reward_config.GGetPrizeType(sngcfg.PrizeType) {
			for _, rr := range cfg.Rewards {
				if val, ok := tmps[rr.CoinType]; ok {
					val.Incr += rr.Incr
				} else {
					tmps[rr.CoinType] = &pb.CoinReward{
						CoinType: rr.CoinType,
						Incr:     rr.Incr,
					}
					total = append(total, tmps[rr.CoinType])
				}
			}
		}
		elem := &pb.SngRoomInfo{
			Title:          sngcfg.Title,
			MatchType:      int32(cfg.MatchType),
			GameType:       int32(cfg.GameType),
			Consume:        sngcfg.Consume,
			EntryFee:       sngcfg.EntryFee,
			Chip:           sngcfg.Chip,
			TableCount:     1,
			MinPlayerCount: cfg.MinPlayerCount,
			PlayerCount:    int32(item.OnlineNumber),
			FirstRewards:   sng_match_rank_reward_config.MGetPrizeTypeLevel(sngcfg.PrizeType, 1).Rewards,
			Rewards:        total,
			RoomId:         item.RoomId,
			CreateTime:     item.CreateTime,
			RoomState:      int32(item.RoomState),
			MatchId:        item.GameId,
		}
		rsp.RoomList = append(rsp.RoomList, elem)
	}
	sort.Slice(rsp.RoomList, func(i, j int) bool {
		return rsp.RoomList[i].MinPlayerCount < rsp.RoomList[j].MinPlayerCount
	})
	return nil
}

// 检查是否需要创建新房间
func (d *SngRoomMgr) check() (rets []*pb.TexasConfig) {
	nobuilds := map[int32]struct{}{}
	for gameId, item := range d.datas {
		cfg := texas_config.MGetID(gameId)
		if cfg == nil {
			mlog.Errorf("SNG德州扑克游戏配置不存在: %d", gameId)
			delete(d.datas, gameId)
			continue
		}
		if int32(item.OnlineNumber) < int32(cfg.MinPlayerCount) {
			nobuilds[cfg.ID] = struct{}{}
		}
	}
	for _, cfg := range texas_config.GGetMatchType(pb.MatchType_MatchTypeSNG) {
		if _, ok := nobuilds[cfg.ID]; !ok {
			nobuilds[cfg.ID] = struct{}{}
			rets = append(rets, cfg)
		}
	}
	return
}

// 创建房间
func (d *SngRoomMgr) build(cfgs ...*pb.TexasConfig) error {
	dst := framework.NewBuilderRouter(uint64(pb.GeneratorType_GeneratorTypeTexas), "BuilderTexasGenerator", "GenRoomIdReq")
	for _, cfg := range cfgs {
		req := &pb.GenRoomIdReq{
			GeneratorType: pb.GeneratorType_GeneratorTypeTexas,
			MatchType:     cfg.MatchType,
			GameType:      cfg.GameType,
			CoinType:      cfg.CoinType,
			Count:         1,
		}
		rsp := &pb.GenRoomIdRsp{}
		if err := cluster.Request(&pb.Head{Dst: dst}, req, rsp); err != nil {
			return err
		}
		if rsp.Head != nil {
			return uerror.ToError(rsp.Head)
		}
		d.datas[cfg.ID] = &pb.TexasRoomData{
			RoomId:     rsp.RoomIdList[0],
			GameId:     cfg.ID,
			RoomState:  pb.RoomStatus_RoomStatusWait,
			CreateTime: time.Now().Unix(),
		}
	}
	return nil
}
