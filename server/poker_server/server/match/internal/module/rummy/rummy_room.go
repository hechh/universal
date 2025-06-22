package rummy

import (
	"poker_server/common/config/repository/rummy_config"
	"poker_server/common/config/repository/rummy_machine_config"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"poker_server/library/util"
	"sort"
)

type MatchRummyRoom struct {
	actor.Actor
	datas map[uint64]*pb.RummyRoomData
}

// NewMatchRummyRoom key => uint64(cfg.GameType)<<32 | uint64(cfg.CoinType)
func NewMatchRummyRoom(id uint64, datas ...*pb.RummyRoomData) *MatchRummyRoom {
	ret := &MatchRummyRoom{datas: make(map[uint64]*pb.RummyRoomData)}
	for _, data := range datas {
		ret.datas[data.RoomId] = data
	}
	ret.Actor.Register(ret)
	ret.Actor.SetId(id)
	ret.Actor.Start()
	return ret
}

// 每个actor的key 按位存放了查询条件
func (d *MatchRummyRoom) GetGameKey(gameType pb.GameType, coinType pb.CoinType) uint64 {
	return uint64(gameType)<<32 | uint64(coinType)
}

// Query 读取房间数据 roomId uint64(gameType)<<32|uint64(coinType)
func (d *MatchRummyRoom) Query(head *pb.Head, roomId uint64) error {
	if tmps := d.check(); len(tmps) > 0 {
		if err := d.build(tmps); err != nil {
			return err
		}
	}

	roomData, ok := d.datas[roomId]
	if !ok {
		return uerror.New(1, pb.ErrorCode_RUMMY_ROOM_NOT_FOUND, "房间不存在 head:%v, roomId:%d", head, roomId)
	}
	return framework.SendResponse(head, roomData)
}

// Without 换桌获取一个同类型可入房间
func (d *MatchRummyRoom) Without(head *pb.Head, req *pb.RummyChangeRoomReq, rsp *pb.RummyChangeRoomRsp) error {
	roomId := req.GetRoomId()
	if tmps := d.check(roomId); len(tmps) > 0 {
		if err := d.build(tmps); err != nil {
			return err
		}
	}

	roomData, ok := d.datas[roomId]
	if !ok {
		return uerror.New(1, pb.ErrorCode_RUMMY_ROOM_NOT_FOUND, "旧房间不存在 head:%v, roomId:%d", head, roomId)
	}

	for _, dataItem := range d.datas {
		if dataItem.RoomId != roomId && roomData.RoomCfg.CoinType == dataItem.RoomCfg.CoinType && roomData.RoomCfg.RoomType == dataItem.RoomCfg.RoomType {
			//人数满足
			if dataItem.Common == nil || dataItem.Common.EmptySeats == nil || int32(len(dataItem.Common.EmptySeats)) >= 0 {
				roomData = dataItem
				break
			}
		}
	}

	req.RoomId = roomData.RoomId
	return framework.Send(framework.SwapToGame(head, head.Uid, "Player", "RummyChangeToRoomReq"), req)
}

// 更新数据
func (d *MatchRummyRoom) Update(head *pb.Head, roomData *pb.RummyRoomData) error {
	if head == nil || roomData == nil {
		return nil
	}

	// 更新数据
	d.datas[roomData.RoomId] = roomData

	// 同步数据
	notify := &pb.UpdateRummyRoomDataNotify{List: []*pb.RummyRoomData{roomData}}
	return actor.SendMsg(&pb.Head{ActorName: "MatchRummyRoomMgr", FuncName: "Collect"}, notify)
}

// 房间列表
func (d *MatchRummyRoom) RoomListReq(head *pb.Head, req *pb.RummyRoomListReq, rsp *pb.RummyRoomListRsp) error {
	if tmps := d.check(); len(tmps) > 0 {
		if err := d.build(tmps); err != nil {
			return err
		}
	}

	for _, item := range d.datas {
		cfg := rummy_config.MGetID(item.GameId)
		rsp.RoomList = append(rsp.RoomList, &pb.RummyRoomPubData{
			RoomId:  item.RoomId,
			RoomCfg: cfg,
			Common:  item.Common,
			Match:   item.Match,
			Stage:   item.Stage,
		})
	}

	// 排序
	sort.Slice(rsp.RoomList, func(i, j int) bool {
		return (rsp.RoomList[i].RoomId) < (rsp.RoomList[j].RoomId)
	})
	return nil
}

// 检查是否需要创建新房间
func (d *MatchRummyRoom) check(blackIds ...uint64) (rets []uint32) {
	nobuilds := map[uint32]struct{}{}

	for roomId, item := range d.datas {
		// 加载配置
		cfg := rummy_config.MGetID(item.GameId)
		if cfg == nil {
			delete(d.datas, roomId)
			mlog.Errorf("Rummy游戏配置不存在: %d", item.RoomId)
			continue
		}

		if len(blackIds) > 0 && util.SliceIsset[uint64](roomId, blackIds) {
			continue
		}

		// 不用新建房间
		if item.Common == nil || item.Common.EmptySeats == nil || int32(len(item.Common.EmptySeats)) >= 0 {
			nobuilds[cfg.ID] = struct{}{}
		}
	}

	id := d.GetId()
	gameType, coinType := pb.GameType(id>>32), pb.CoinType(id&0xFFFFFFFF)
	for _, cfg := range rummy_config.GGetGameTypeCoinType(gameType, coinType) {
		if _, ok := nobuilds[cfg.ID]; !ok {
			nobuilds[cfg.ID] = struct{}{}
			rets = append(rets, cfg.ID)
		}
	}
	return
}

// 创建房间
func (d *MatchRummyRoom) build(tmps []uint32) error {
	ll := len(tmps)
	id := d.Actor.GetId()
	dst := framework.NewBuilderRouter(uint64(pb.GeneratorType_GeneratorTypeRummy), "BuilderRummyGenerator", "GenRoomIdReq")
	head := framework.NewHead(dst, pb.RouterType_RouterTypeGameType, uint64(pb.GeneratorType_GeneratorTypeRummy))
	newReq := &pb.GenRoomIdReq{
		GeneratorType: pb.GeneratorType_GeneratorTypeRummy,
		GameType:      pb.GameType(id >> 32),
		CoinType:      pb.CoinType(uint32(id)),
		Count:         int32(ll),
	}
	newRsp := &pb.GenRoomIdRsp{}
	// 获取房间ID
	if err := framework.Request(head, newReq, newRsp); err != nil {
		return err
	}
	if newRsp.Head != nil {
		return uerror.ToError(newRsp.Head)
	}

	// 生成房间
	notify := &pb.UpdateRummyRoomDataNotify{}
	for i, gameId := range tmps {
		rummyCfg := rummy_config.MGetID(gameId)
		item := &pb.RummyRoomData{
			RoomId:     newRsp.RoomIdList[i],
			GameId:     gameId,
			Stage:      pb.GameState_Rummy_STAGE_INIT,
			RoomCfg:    rummyCfg,
			MachineCfg: rummy_machine_config.MGetGameType(rummyCfg.GameType),
		}
		notify.List = append(notify.List, item)
		d.datas[item.RoomId] = item
	}
	actor.SendMsg(&pb.Head{ActorName: "MatchRummyRoomMgr", FuncName: "Collect"}, notify)
	return nil
}
