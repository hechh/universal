package internal

import (
	"poker_server/common/pb"
	"poker_server/common/room_util"
	"poker_server/framework/actor"
	"poker_server/library/mlog"
	"poker_server/library/snowflake"
	"poker_server/server/game/module/http"
	"reflect"
	"strconv"
)

var (
	propMgr = NewPropMgr()
)

type PropMgr struct {
	actor.ActorPool
}

func Close() {
	propMgr.ActorPool.Stop()
}

func NewPropMgr() *PropMgr {
	ret := &PropMgr{}
	ret.ActorPool.Register(ret, 100)
	ret.ActorPool.ParseFunc(reflect.TypeOf(ret))
	ret.Start()
	actor.Register(ret)
	return ret
}

func (p *PropMgr) TexasFinishNotify(head *pb.Head, event *pb.TexasFinishNotify) error {
	// 如果玩家超时被清理
	_, gameType, coinType := room_util.TexasRoomIdTo(event.RoomId)
	uuid, _ := snowflake.GenUUID()
	param := &pb.TransParam{
		GameSn:   strconv.FormatUint(uuid, 10),
		GameType: gameType,
		CoinType: coinType,
		Incr:     event.Incr,
		Uid:      head.Uid,
	}
	rsp := &pb.HttpTransferOutRsp{}
	if err := http.ChargeTransOutRequest(param, rsp); err != nil {
		mlog.Error(head, "ChargeTransOut http err: %v", err)
		return err
	}
	mlog.Info(head, "PlayerMgr ChargeTransOut uid:%d, param:%v, rsp:%v", head.Uid, param, rsp)
	return nil
}
