package texas_room

import (
	"poker_server/common/config/repository/texas_config"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/framework/cluster"
	"poker_server/library/uerror"
	"sort"
	"time"
)

type MatchTexasRoom struct {
	actor.Actor
	datas   map[uint64]*pb.TexasRoomData
	updates []uint64
	deletes []uint64
}

func NewMatchTexasRoom(id uint64, datas ...*pb.TexasRoomData) *MatchTexasRoom {
	ret := &MatchTexasRoom{datas: make(map[uint64]*pb.TexasRoomData)}
	for _, data := range datas {
		ret.datas[data.RoomId] = data
	}
	ret.Actor.Register(ret)
	ret.Actor.SetId(id)
	ret.Actor.Start()
	return ret
}

func (d *MatchTexasRoom) Update(head *pb.Head, event *pb.UpdateTexasRoomDataNotify) error {
	d.datas[event.Data.RoomId] = event.Data
	d.updates = append(d.updates, event.Data.RoomId)

	return cluster.Send(&pb.Head{
		Src: framework.NewSrcRouter(d.GetId(), d.GetActorName()),
		Dst: framework.NewDbRouter(uint64(pb.DataType_DataTypeTexasRoom), "DbTexasRoomMgr", "Update"),
	}, event)
}

func (d *MatchTexasRoom) Delete(head *pb.Head, event *pb.DeleteTexasRoomDataNotify) error {
	delete(d.datas, event.RoomId)
	d.deletes = append(d.deletes, event.RoomId)

	return cluster.Send(&pb.Head{
		Src: framework.NewSrcRouter(d.GetId(), "MatchTexasRoom"),
		Dst: framework.NewDbRouter(uint64(pb.DataType_DataTypeTexasRoom), "DbTexasRoomMgr", "Delete"),
	}, event)
}

func (d *MatchTexasRoom) HasRoomReq(head *pb.Head, req *pb.HasRoomReq, rsp *pb.HasRoomRsp) error {
	_, ok := d.datas[req.RoomId]
	rsp.IsExist = ok
	return nil
}

func (d *MatchTexasRoom) Query(head *pb.Head, req *pb.TexasQueryReq, rsp *pb.TexasQueryRsp) error {
	if item, ok := d.datas[req.RoomId]; ok {
		rsp.Data = item
		return nil
	}
	return uerror.New(1, pb.ErrorCode_TEXAS_ROOM_NOT_FOUND, "房间已经不存在%v", req)
}

func (d *MatchTexasRoom) MatchRoomReq(head *pb.Head, req *pb.TexasMatchRoomReq, rsp *pb.TexasMatchRoomRsp) error {
	cfg := texas_config.MGetID(int32(req.TableId))
	if cfg == nil {
		return uerror.New(1, pb.ErrorCode_CONFIG_NOT_FOUND, "房间配置不存在")
	}
	// 检测是否新建房间
	if tmps := d.check(req.RoomId); len(tmps) > 0 {
		if err := d.build(tmps); err != nil {
			return err
		}
	}
	// 获取房间列表
	arrs := []*pb.TexasRoomData{}
	for _, val := range d.datas {
		if val.GameId != int32(req.TableId) || val.RoomId == req.RoomId {
			continue
		}
		arrs = append(arrs, val)
	}
	sort.Slice(arrs, func(i, j int) bool {
		if len(arrs[i].Table.ChairInfo) == len(arrs[j].Table.ChairInfo) {
			if arrs[i].OnlineNumber == arrs[j].OnlineNumber {
				return arrs[i].CreateTime < arrs[j].CreateTime
			}
			return arrs[i].OnlineNumber > arrs[j].OnlineNumber
		}
		return len(arrs[i].Table.ChairInfo) > len(arrs[j].Table.ChairInfo)
	})
	for _, item := range arrs {
		if len(item.Table.Players) >= int(cfg.MaxRoomPlayerCount) || len(item.Table.ChairInfo) >= int(cfg.MaxPlayerCount) {
			continue
		}
		rsp.TableId = uint32(item.GameId)
		rsp.RoomId = item.RoomId
		break
	}
	return nil
}

func (d *MatchTexasRoom) RoomListReq(head *pb.Head, req *pb.TexasRoomListReq, rsp *pb.TexasRoomListRsp) error {
	if tmps := d.check(); len(tmps) > 0 {
		if err := d.build(tmps); err != nil {
			return err
		}
	}
	tmps := map[int32]int64{}
	arrs := map[int32][]*pb.TexasRoomInfo{}
	for _, item := range d.datas {
		cfg := texas_config.MGetID(item.GameId)
		elem := &pb.TexasRoomInfo{
			RoomId:             item.RoomId,
			TableId:            uint32(item.GameId),
			GameType:           int32(cfg.GameType),
			RoomStage:          int32(cfg.RoomType),
			CoinType:           int32(cfg.CoinType),
			RoomState:          int32(item.RoomState),
			CreateTime:         item.CreateTime,
			FinishTime:         item.CreateTime + cfg.RoomKeepLive*60,
			MinBuyIn:           cfg.MinBuyIn,
			MaxBuyIn:           cfg.MaxBuyIn,
			SmallBlind:         int64(cfg.SmallBlind),
			BigBlind:           int64(cfg.BigBlind),
			MaxPlayerCount:     int32(cfg.MaxPlayerCount),
			PlayerCount:        int32(item.OnlineNumber),
			ChairCount:         int32(len(item.Table.ChairInfo)),
			MaxRoomPlayerCount: int32(cfg.MaxRoomPlayerCount),
		}
		tmps[cfg.ID] += item.OnlineNumber
		arrs[cfg.ID] = append(arrs[cfg.ID], elem)
	}
	for iid, vals := range arrs {
		sort.Slice(vals, func(i, j int) bool {
			if vals[i].ChairCount == vals[j].ChairCount {
				return vals[i].CreateTime < vals[j].CreateTime
			}
			return vals[i].ChairCount > vals[j].ChairCount
		})
		for _, item := range vals {
			if item.MaxRoomPlayerCount <= item.PlayerCount || item.ChairCount >= item.MaxPlayerCount {
				continue
			}
			item.PlayerCount = int32(tmps[iid])
			rsp.RoomList = append(rsp.RoomList, item)
			break
		}
	}
	sort.Slice(rsp.RoomList, func(i, j int) bool {
		return (rsp.RoomList[i].TableId) < (rsp.RoomList[j].TableId)
	})
	return nil
}

// 检查是否需要创建新房间
func (d *MatchTexasRoom) check(roomIds ...uint64) (rets []int32) {
	nobuilds := map[int32]struct{}{}
	now := time.Now().Unix()
	filter := map[uint64]struct{}{}
	for _, roomId := range roomIds {
		filter[roomId] = struct{}{}
	}

	for roomId, item := range d.datas {
		if _, ok := filter[roomId]; ok {
			continue
		}
		cfg := texas_config.MGetID(item.GameId)
		if cfg == nil || item.CreateTime+int64(cfg.RoomKeepLive*60) <= now {
			continue
		}
		if item.Table == nil {
			item.Table = &pb.TexasTableData{CurState: pb.GameState_TEXAS_INIT}
		}
		if item.Table.Players == nil {
			item.Table.Players = make(map[uint64]*pb.TexasPlayerData)
		}
		if item.Table.ChairInfo == nil {
			item.Table.ChairInfo = make(map[uint32]uint64)
		}
		if int32(len(item.Table.ChairInfo)) < int32(cfg.MaxPlayerCount) {
			nobuilds[cfg.ID] = struct{}{}
		}
	}
	id := d.GetId()
	gameType, coinType := pb.GameType(id>>16)&0xFFFF, pb.CoinType(id&0xFFFF)
	for _, cfg := range texas_config.GGetGameTypeMatchTypeCoinType(gameType, pb.MatchType_MatchTypeNone, coinType) {
		if _, ok := nobuilds[cfg.ID]; !ok {
			nobuilds[cfg.ID] = struct{}{}
			rets = append(rets, cfg.ID)
		}
	}
	return
}

// 创建房间
func (d *MatchTexasRoom) build(tmps []int32) error {
	ll := len(tmps)
	id := d.Actor.GetId()
	head := &pb.Head{
		Src: framework.NewSrcRouter(id, "MatchTexasRoom"),
		Dst: framework.NewBuilderRouter(uint64(pb.GeneratorType_GeneratorTypeTexas), "BuilderTexasGenerator", "GenRoomIdReq"),
	}
	newreq := &pb.GenRoomIdReq{
		GeneratorType: pb.GeneratorType_GeneratorTypeTexas,
		GameType:      pb.GameType((id >> 16) & 0xFFFF),
		CoinType:      pb.CoinType(uint32(id & 0xFFFF)),
		Count:         int32(ll),
	}
	newrsp := &pb.GenRoomIdRsp{}

	// 获取房间ID
	if err := cluster.Request(head, newreq, newrsp); err != nil {
		return err
	}
	if newrsp.Head != nil {
		return uerror.ToError(newrsp.Head)
	}

	// 生成房间
	for i, gameId := range tmps {
		item := &pb.TexasRoomData{
			RoomId:     newrsp.RoomIdList[i],
			GameId:     gameId,
			RoomState:  pb.RoomStatus_RoomStatusWait,
			CreateTime: time.Now().Unix(),
			Table: &pb.TexasTableData{
				CurState:  pb.GameState_TEXAS_INIT,
				Players:   make(map[uint64]*pb.TexasPlayerData),
				ChairInfo: make(map[uint32]uint64),
			},
		}
		d.datas[item.RoomId] = item
		d.updates = append(d.updates, item.RoomId)
	}
	return nil
}
