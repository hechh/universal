package rummy

import (
	"poker_server/common/config/repository/rummy_config"
	"poker_server/common/config/repository/rummy_machine_config"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"poker_server/library/util"
	"sort"
	"time"

	"github.com/golang/protobuf/proto"
)

const TIMEOUT = 10000

type MatchRummyRoom struct {
	actor.Actor
	datas  map[uint64]*pb.RummyRoomData
	queues map[uint32]*RummyMatchQueue
}

type RummyMatchQueue struct {
	topRoom     *pb.RummyRoomData
	cfg         *pb.RummyConfig
	queue       []func()
	uids        []uint64
	playerInfos map[uint64]*pb.PlayerInfo
	timeout     int64 //队列超时时间 取队头玩家
}

// NewMatchRummyRoom key => uint64(cfg.GameType)<<32 | uint64(cfg.CoinType)
func NewMatchRummyRoom(id uint64, datas ...*pb.RummyRoomData) *MatchRummyRoom {
	types := util.DestructMatchId(id)
	mods := rummy_config.GGetGameTypeCoinType(types.GameType, types.CoinType)

	ret := &MatchRummyRoom{datas: make(map[uint64]*pb.RummyRoomData), queues: make(map[uint32]*RummyMatchQueue, len(mods))}
	for _, data := range datas {
		ret.datas[data.RoomId] = data
	}

	for i := range mods {
		ret.queues[mods[i].ID] = &RummyMatchQueue{
			cfg:         mods[i],
			queue:       make([]func(), 0, mods[i].GetMaxPlayerCount()),
			uids:        make([]uint64, 0, mods[i].GetMaxPlayerCount()),
			playerInfos: make(map[uint64]*pb.PlayerInfo, mods[i].GetMaxPlayerCount()),
		}
	}

	ret.Actor.Register(ret)
	ret.Actor.SetId(id)
	ret.Actor.Start()

	// 启动定时器
	head := &pb.Head{SendType: pb.SendType_POINT, ActorName: "MatchRummyRoom", FuncName: "OnTick"}
	err := ret.RegisterTimer(head, 50*time.Millisecond, -1)
	if err != nil {
		mlog.Debug(head, "register timer err: %v", err)
	}
	return ret
}

func (d *MatchRummyRoom) OnTick() {
	nowMs := time.Now().UnixMilli()
	for i := range d.queues {
		if len(d.queues[i].queue) > 0 {
			if nowMs >= d.queues[i].timeout { //超时处理
				d.BatchJoinRoom(d.queues[i])
			}
		}
	}
}

func (d *MatchRummyRoom) RummyMatchReq(head *pb.Head, req *pb.RummyMatchReq, rsp *pb.RummyMatchRsp) error {
	if d.queues[req.CfgId] == nil {
		return uerror.New(1, pb.ErrorCode_CONFIG_NOT_FOUND, "config not found")
	}

	if req.Coin < d.queues[req.CfgId].cfg.MinBuyIn {
		return uerror.New(1, pb.ErrorCode_GAME_PROP_NOT_ENOUGH, "玩家道具不足加入游戏")
	}

	if d.queues[req.CfgId].topRoom == nil {
		for i := range d.datas {
			if d.datas[i].Stage == pb.GameState_RummyExt_STAGE_INIT && d.datas[i].RoomCfg.ID == req.CfgId {
				d.queues[req.CfgId].topRoom = d.datas[i]
				break
			}
		}

		if d.queues[req.CfgId].topRoom == nil {
			err := d.build([]uint32{req.CfgId})
			if err != nil {
				return uerror.New(1, pb.ErrorCode_RUMMY_ROOM_NOT_FOUND, "build new room failed")
			}
		}

		for i := range d.datas {
			if d.datas[i].Stage == pb.GameState_RummyExt_STAGE_INIT && d.datas[i].RoomCfg.ID == req.CfgId {
				d.queues[req.CfgId].topRoom = d.datas[i]
				break
			}
		}

		if d.queues[req.CfgId].topRoom == nil {
			err := d.build([]uint32{req.CfgId})
			if err != nil {
				return uerror.New(1, pb.ErrorCode_RUMMY_ROOM_NOT_FOUND, "match room failed")
			}
		}
	}

	d.queues[req.CfgId].queue = append(d.queues[req.CfgId].queue, func() {
		newHead := &pb.Head{
			Uid: head.Uid,
			Src: head.Src,
		}
		newHead.Cmd = uint32(pb.CMD_RUMMY_JOIN_ROOM_REQ)
		newHead.ActorId = d.queues[req.CfgId].topRoom.RoomId
		newHead.Dst = framework.NewGameRouter(head.Uid, "Player", "RummyBatchJoinReq")
		joinReq := &pb.RummyJoinRoomReq{
			Coin:       req.Coin,
			PlayerInfo: req.PlayerInfo,
			RoomId:     d.queues[req.CfgId].topRoom.RoomId,
			IsReady:    true,
		}
		err := cluster.Send(newHead, joinReq)
		if err != nil {
			mlog.Infof("BatchJoinRoom err: %v", err)
		}
	})
	d.queues[req.CfgId].uids = append(d.queues[req.CfgId].uids, req.PlayerInfo.Uid)
	d.queues[req.CfgId].playerInfos[req.PlayerInfo.Uid] = req.PlayerInfo

	if d.queues[req.CfgId].timeout == 0 {
		d.queues[req.CfgId].timeout = time.Now().UnixMilli() + TIMEOUT
	}

	rsp.CurCount = uint32(len(d.queues[req.CfgId].queue))
	if d.queues[req.CfgId].topRoom.Common != nil && d.queues[req.CfgId].topRoom.Common.PlayerIds != nil {
		rsp.CurCount += uint32(len(d.queues[req.CfgId].topRoom.Common.PlayerIds))
	}

	rsp.Cfg = d.queues[req.CfgId].cfg
	rsp.CfgId = req.CfgId
	rsp.Timeout = d.queues[req.CfgId].timeout

	ntf := &pb.RummyMatchNtf{
		CurCount: rsp.CurCount,
		Timeout:  d.queues[req.CfgId].timeout,
		CfgId:    req.CfgId,
	}

	uids := d.queues[req.CfgId].uids

	if d.queues[req.CfgId].topRoom.Common != nil && d.queues[req.CfgId].topRoom.Common.PlayerIds != nil {
		uids = append(uids, d.queues[req.CfgId].topRoom.Common.PlayerIds...)
	}

	err := d.NotifyToClient(uids, pb.RummyEventType_MatchNtf, ntf)
	mlog.Infof("RummyMatchReq send ntf err: %v", err)

	if rsp.CurCount >= d.queues[req.CfgId].cfg.MaxPlayerCount {
		d.BatchJoinRoom(d.queues[req.CfgId])
	}
	return nil
}

// BatchJoinRoom 匹配成功或者匹配时间结束 批量将用户加入房间
func (d *MatchRummyRoom) BatchJoinRoom(queue *RummyMatchQueue) {
	for i := range queue.queue {
		ff := queue.queue[i]
		ff()
	}

	queue.queue = queue.queue[:0]
	queue.uids = queue.uids[:0]
	queue.timeout = 0
	queue.topRoom = nil
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
		return uerror.New(1, pb.ErrorCode_RUMMY_ROOM_NOT_FOUND, "房间不存在roomId:%v, head:%v ", head, roomId)
	}
	return cluster.SendResponse(head, roomData)
}

func (d *MatchRummyRoom) HasRoomReq(head *pb.Head, req *pb.HasRoomReq, rsp *pb.HasRoomRsp) error {
	_, ok := d.datas[req.RoomId]
	rsp.IsExist = ok
	return nil
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
		if dataItem.RoomId != roomId && roomData.RoomCfg.ID == dataItem.RoomCfg.ID {
			if dataItem.Common == nil || dataItem.Common.EmptySeats == nil || int32(len(dataItem.Common.EmptySeats)) >= 0 {
				roomData = dataItem
				break
			}
		}
	}

	rsp.RoomId = roomData.RoomId
	return nil
}

// Update 更新数据
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

// RoomListReq 房间列表
func (d *MatchRummyRoom) RoomListReq(head *pb.Head, req *pb.RummyRoomListReq, rsp *pb.RummyRoomListRsp) error {
	if tmps := d.check(); len(tmps) > 0 {
		if err := d.build(tmps); err != nil {
			return err
		}
	}

	configs := rummy_config.GGetGameTypeCoinType(pb.GameType(req.GetGameType()), pb.CoinType(req.GetCoinType()))
	retMap := make(map[uint32]any, len(configs))
	for i := range configs {
		retMap[configs[i].ID] = struct{}{}
	}

	for _, item := range d.datas {
		if retMap[item.RoomCfg.ID] != nil && (item.Common == nil || len(item.Common.Seats) < int(item.RoomCfg.MaxPlayerCount)) {
			rsp.RoomList = append(rsp.RoomList, &pb.RummyRoomPubData{
				RoomId:  item.RoomId,
				RoomCfg: item.RoomCfg,
				Common:  item.Common,
				Match:   item.Match,
				Stage:   item.Stage,
			})
			delete(retMap, item.RoomCfg.ID)
		}

		if len(retMap) == 0 {
			break
		}
	}

	// 排序
	sort.Slice(rsp.RoomList, func(i, j int) bool {
		if rsp.RoomList[i].RoomCfg.ID != rsp.RoomList[j].RoomCfg.ID { //对比配置
			return (rsp.RoomList[i].RoomCfg.ID) < (rsp.RoomList[j].RoomCfg.ID)
		}

		if rsp.RoomList[i].Stage != rsp.RoomList[j].Stage {
			return int32(rsp.RoomList[i].Stage) < int32(rsp.RoomList[j].Stage)
		}

		countLeft, countRight := 0, 0
		if rsp.RoomList[i].Common != nil && rsp.RoomList[i].Common.Players != nil {
			countLeft = len(rsp.RoomList[i].Common.Players)
		}
		if rsp.RoomList[j].Common != nil && rsp.RoomList[j].Common.Players != nil {
			countRight = len(rsp.RoomList[j].Common.Players)
		}

		// 人数检查
		return countLeft < countRight
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
		if item.Common == nil || len(item.Common.Seats) < int(item.RoomCfg.MaxPlayerCount) {
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

	head := &pb.Head{
		Src: framework.NewSrcRouter(d.GetId(), d.GetActorName()),
		Dst: framework.NewBuilderRouter(uint64(pb.GeneratorType_GeneratorTypeRummy), "BuilderRummyGenerator", "GenRoomIdReq"),
	}
	newReq := &pb.GenRoomIdReq{
		GeneratorType: pb.GeneratorType_GeneratorTypeRummy,
		GameType:      pb.GameType(id >> 32),
		CoinType:      pb.CoinType(uint32(id)),
		Count:         int32(ll),
	}
	newRsp := &pb.GenRoomIdRsp{}
	// 获取房间ID
	if err := cluster.Request(head, newReq, newRsp); err != nil {
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
			RoomCfg:    rummyCfg,
			Status:     pb.RoomStatus_RoomStatusWait,
			MachineCfg: rummy_machine_config.MGetGameType(rummyCfg.GameType),
		}

		switch rummyCfg.GameType {
		case pb.GameType_GameTypePR:
			item.Stage = pb.GameState_Rummy_STAGE_INIT
		default:
			item.Stage = pb.GameState_RummyExt_STAGE_INIT
		}

		notify.List = append(notify.List, item)
		d.datas[item.RoomId] = item
	}
	actor.SendMsg(&pb.Head{ActorName: "MatchRummyRoomMgr", FuncName: "Collect"}, notify)
	return nil
}

func (d *MatchRummyRoom) NewRummyEventNotify(notifyId pb.RummyEventType, msg proto.Message) *pb.RummyEventNotify {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return nil
	}
	return &pb.RummyEventNotify{Event: notifyId, Content: buf}
}

// NotifyToClient 组播 根据uid进行组播
func (d *MatchRummyRoom) NotifyToClient(uids []uint64, notifyId pb.RummyEventType, msg proto.Message) error {
	ntf := d.NewRummyEventNotify(notifyId, msg)

	newHead := &pb.Head{
		Src: framework.NewSrcRouter(d.GetId(), d.GetActorName()),
		Cmd: uint32(pb.CMD_RUMMY_EVENT_NOTIFY),
	}
	return cluster.SendToClient(newHead, ntf, uids...)
}

// RummyCancelMatchReq 主动退出匹配队列
func (d *MatchRummyRoom) RummyCancelMatchReq(head *pb.Head, req *pb.RummyCancelMatchReq, rsp *pb.RummyCancelMatchRsp) error {
	if d.queues[req.CfgId] == nil {
		return uerror.New(1, pb.ErrorCode_CONFIG_NOT_FOUND, "config not found")
	}

	if d.queues[req.CfgId].topRoom == nil {
		return uerror.New(1, pb.ErrorCode_MATCH_FINISH, "match have finished")
	}

	for i, uid := range d.queues[req.CfgId].uids {
		if uid == head.Uid {
			d.queues[req.CfgId].uids = append(d.queues[req.CfgId].uids[:i], d.queues[req.CfgId].uids[i+1:]...)
			d.queues[req.CfgId].queue = append(d.queues[req.CfgId].queue[:i], d.queues[req.CfgId].queue[i+1:]...)
			delete(d.queues[req.CfgId].playerInfos, uid)

			if len(d.queues[req.CfgId].queue) == 0 {
				d.queues[req.CfgId].timeout = 0
				d.queues[req.CfgId].topRoom = nil
			}

			rsp.CoinType = req.CoinType
			rsp.GameType = req.GameType
			rsp.CfgId = req.CfgId

			return nil
		}
	}

	return uerror.New(1, pb.ErrorCode_PLAYER_NOT_FOUND, "player undefined")
}
